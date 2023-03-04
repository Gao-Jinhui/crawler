package master

import (
	"context"
	proto "crawler/internal/pkg/proto/crawler"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/golang/protobuf/ptypes/empty"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Master struct {
	ID         string
	ready      int32
	leaderID   string
	workNodes  map[string]*NodeSpec
	resources  map[string]*ResourceSpec
	IDGen      *snowflake.Node
	etcdCli    *clientv3.Client
	forwardCli proto.CrawlerMasterService
	rlock      sync.Mutex
	options
}

func New(id string, opts ...Option) (*Master, error) {
	m := &Master{}

	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	m.options = options
	m.resources = make(map[string]*ResourceSpec)

	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}
	m.IDGen = node
	ipv4, err := getLocalIP()
	if err != nil {
		return nil, err
	}
	m.ID = genMasterID(id, ipv4, m.GRPCAddress)
	m.logger.Sugar().Debugln("master_id:", m.ID)
	endpoints := []string{m.registryURL}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		return nil, err
	}
	m.etcdCli = cli

	m.updateWorkNodes()
	m.AddSeed()
	go m.Campaign()
	return m, nil
}

func (m *Master) SetForwardCli(forwardCli proto.CrawlerMasterService) {
	m.forwardCli = forwardCli
}

func genMasterID(id string, ipv4 string, GRPCAddress string) string {
	return "master" + id + "-" + ipv4 + GRPCAddress
}

func (m *Master) IsLeader() bool {
	return atomic.LoadInt32(&m.ready) != 0
}

func (m *Master) Campaign() {
	s, err := concurrency.NewSession(m.etcdCli, concurrency.WithTTL(5))
	if err != nil {
		fmt.Println("NewSession", "error", "err", err)
	}
	defer s.Close()

	// 创建一个新的etcd选举election
	e := concurrency.NewElection(s, "/crawler/election")
	leaderCh := make(chan error)
	go m.elect(e, leaderCh)
	leaderChange := e.Observe(context.Background())
	select {
	case resp := <-leaderChange:
		m.logger.Info("watch leader change", zap.String("leader:", string(resp.Kvs[0].Value)))
		m.leaderID = string(resp.Kvs[0].Value)
	}
	workerNodeChange := m.WatchWorker()

	for {
		select {
		case err := <-leaderCh:
			if err != nil {
				m.logger.Error("leader elect failed", zap.Error(err))
				go m.elect(e, leaderCh)
			} else {
				m.logger.Info("master start change to leader")
				m.leaderID = m.ID
				if !m.IsLeader() {
					if err := m.BecomeLeader(); err != nil {
						m.logger.Error("BecomeLeader failed", zap.Error(err))
					}
				}
			}
		case resp := <-leaderChange:
			if len(resp.Kvs) > 0 {
				m.logger.Info("watch leader change", zap.String("leader:", string(resp.Kvs[0].Value)))
				m.leaderID = string(resp.Kvs[0].Value)
				if m.ID != string(resp.Kvs[0].Value) {
					//当前已不再是leader
					atomic.StoreInt32(&m.ready, 0)
				}
			}
		case resp := <-workerNodeChange:
			m.logger.Info("watch worker change", zap.Any("worker:", resp))
			m.updateWorkNodes()
			if err := m.loadResource(); err != nil {
				m.logger.Error("loadResource failed:%w", zap.Error(err))
			}
			m.reAssign()
		case <-time.After(20 * time.Second):
			rsp, err := e.Leader(context.Background())
			if err != nil {
				m.logger.Info("get Leader failed", zap.Error(err))
				if errors.Is(err, concurrency.ErrElectionNoLeader) {
					go m.elect(e, leaderCh)
				}
			}
			if rsp != nil && len(rsp.Kvs) > 0 {
				m.logger.Info("get Leader", zap.String("value", string(rsp.Kvs[0].Value)))
				m.leaderID = string(rsp.Kvs[0].Value)
				if m.IsLeader() && m.ID != string(rsp.Kvs[0].Value) {
					//当前已不再是leader
					atomic.StoreInt32(&m.ready, 0)
				}
			}
		}
	}
}

func (m *Master) elect(e *concurrency.Election, ch chan error) {
	// 堵塞直到选取成功
	err := e.Campaign(context.Background(), m.ID)
	ch <- err
}

func (m *Master) WatchWorker() chan *registry.Result {
	watch, err := m.registry.Watch(registry.WatchService(WorkerServiceName))
	if err != nil {
		panic(err)
	}
	ch := make(chan *registry.Result)
	go func() {
		for {
			res, err := watch.Next()
			if err != nil {
				m.logger.Info("watch worker service failed", zap.Error(err))
				continue
			}
			ch <- res
		}
	}()
	return ch

}
func (m *Master) BecomeLeader() error {
	m.updateWorkNodes()
	if err := m.loadResource(); err != nil {
		return fmt.Errorf("loadResource failed:%w", err)
	}
	m.reAssign()
	atomic.StoreInt32(&m.ready, 1)
	return nil
}

func (m *Master) updateWorkNodes() {
	services, err := m.registry.GetService(WorkerServiceName)
	if err != nil {
		m.logger.Error("get service", zap.Error(err))
	}
	m.rlock.Lock()
	defer m.rlock.Unlock()
	nodes := make(map[string]*NodeSpec)
	if len(services) > 0 {
		for _, spec := range services[0].Nodes {
			nodes[spec.Id] = &NodeSpec{
				Node: spec,
			}
		}
	}
	added, deleted, changed := workNodeDiff(m.workNodes, nodes)
	m.logger.Sugar().Info("worker joined: ", added, ", leaved: ", deleted, ", changed: ", changed)
	m.workNodes = nodes
}

func (m *Master) AddResources(rs []*ResourceSpec) {
	for _, r := range rs {
		m.addResources(r)
	}
}

func (m *Master) addResources(r *ResourceSpec) (*NodeSpec, error) {
	r.ID = m.IDGen.Generate().String()
	ns, err := m.Assign(r)
	if err != nil {
		m.logger.Error("assign failed", zap.Error(err))
		return nil, err
	}

	if ns.Node == nil {
		m.logger.Error("no node to assgin")
		return nil, err
	}

	r.AssignedNode = ns.Node.Id + "|" + ns.Node.Address
	r.CreationTime = time.Now().UnixNano()
	m.logger.Debug("add resource", zap.Any("specs", r))

	_, err = m.etcdCli.Put(context.Background(), getResourcePath(r.Name), Encode(r))
	if err != nil {
		m.logger.Error("put etcd failed", zap.Error(err))
		return nil, err
	}

	m.resources[r.Name] = r
	ns.Payload++
	return ns, nil
}

func (m *Master) Assign(r *ResourceSpec) (*NodeSpec, error) {
	candidates := make([]*NodeSpec, 0, len(m.workNodes))

	for _, node := range m.workNodes {
		candidates = append(candidates, node)
	}

	//  找到最低的负载
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Payload < candidates[j].Payload
	})

	if len(candidates) > 0 {
		return candidates[0], nil
	}

	return nil, errors.New("no worker nodes")
}

func (m *Master) AddSeed() {
	rs := make([]*ResourceSpec, 0, len(m.Seeds))
	for _, seed := range m.Seeds {
		if seed == nil {
			continue
		}
		resp, err := m.etcdCli.Get(context.Background(), getResourcePath(seed.Name), clientv3.WithPrefix(), clientv3.WithSerializable())
		if err != nil {
			m.logger.Error("etcd get faiiled", zap.Error(err))
			continue
		}
		if len(resp.Kvs) == 0 {
			r := &ResourceSpec{
				Name: seed.Name,
			}
			rs = append(rs, r)
		}
	}

	m.AddResources(rs)
}

func (m *Master) loadResource() error {
	resp, err := m.etcdCli.Get(context.Background(), RESOURCEPATH, clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return fmt.Errorf("etcd get failed")
	}

	resources := make(map[string]*ResourceSpec)
	for _, kv := range resp.Kvs {
		r, err := Decode(kv.Value)
		if err == nil && r != nil {
			resources[r.Name] = r
		}
	}
	m.logger.Info("leader init load resource", zap.Int("lenth", len(m.resources)))
	m.rlock.Lock()
	defer m.rlock.Unlock()
	m.resources = resources

	for _, r := range m.resources {
		if r.AssignedNode != "" {
			id, err := getNodeID(r.AssignedNode)
			if err != nil {
				m.logger.Error("getNodeID failed", zap.Error(err))
			}
			node, ok := m.workNodes[id]
			if ok {
				node.Payload++
			}
		}
	}
	return nil
}

func (m *Master) reAssign() {
	rs := make([]*ResourceSpec, 0, len(m.resources))
	m.rlock.Lock()
	defer m.rlock.Unlock()
	for _, r := range m.resources {
		if r.AssignedNode == "" {
			rs = append(rs, r)
			continue
		}

		id, err := getNodeID(r.AssignedNode)

		if err != nil {
			m.logger.Error("get nodeid failed", zap.Error(err))
		}

		if _, ok := m.workNodes[id]; !ok {
			rs = append(rs, r)
		}
	}
	m.AddResources(rs)
}

func getNodeID(assigned string) (string, error) {
	node := strings.Split(assigned, "|")
	if len(node) < 2 {
		return "", errors.New("")
	}
	id := node[0]
	return id, nil
}

func (m *Master) DeleteResource(ctx context.Context, spec *proto.ResourceSpec, empty *empty.Empty) error {
	if !m.IsLeader() && m.leaderID != "" && m.leaderID != m.ID {
		addr := getLeaderAddress(m.leaderID)
		_, err := m.forwardCli.DeleteResource(ctx, spec, client.WithAddress(addr))
		return err
	}

	m.rlock.Lock()
	defer m.rlock.Unlock()
	r, ok := m.resources[spec.Name]

	if !ok {
		return errors.New("no such task")
	}

	if _, err := m.etcdCli.Delete(context.Background(), getResourcePath(spec.Name)); err != nil {
		return err
	}
	delete(m.resources, spec.Name)
	if r.AssignedNode != "" {
		nodeID, err := getNodeID(r.AssignedNode)
		if err != nil {
			return err
		}

		if ns, ok := m.workNodes[nodeID]; ok {
			ns.Payload -= 1
		}
	}
	return nil
}

func (m *Master) AddResource(ctx context.Context, req *proto.ResourceSpec, resp *proto.NodeSpec) error {
	if !m.IsLeader() && m.leaderID != "" && m.leaderID != m.ID {
		addr := getLeaderAddress(m.leaderID)
		fmt.Println(">>>>>>>>>>leader addr: ", addr)
		nodeSpec, err := m.forwardCli.AddResource(ctx, req, client.WithAddress(addr))
		fmt.Println("error :>>>>>>>>>>", err.Error())
		resp.Id = nodeSpec.Id
		resp.Address = nodeSpec.Address
		return err
	}
	m.rlock.Lock()
	defer m.rlock.Unlock()
	nodeSpec, err := m.addResources(&ResourceSpec{Name: req.Name})
	if nodeSpec != nil {
		resp.Id = nodeSpec.Node.Id
		resp.Address = nodeSpec.Node.Address
	}
	return err
}

func getLeaderAddress(address string) string {
	s := strings.Split(address, "-")
	if len(s) < 2 {
		return ""
	}
	return s[1]
}
