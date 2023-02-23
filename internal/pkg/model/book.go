package model

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Name      string
	Author    string
	Page      int
	Publisher string
	Score     string
	Price     string
	Intro     string
	Url       string
}
