package model

type Flat struct {
	Ref string
	Publisher string
	Address string
	Telephone string
	Price float32
	Agency bool
	Description string
	District string
	City string
	Province string
	State string
	Country string
	Pictures []ImgNode
}

type FlatReact func(*Flat) error

type FlatRepo interface {
	Save(Flat) error
	FindAllByAddress(string, FlatReact) error
	FindByAddress(string, int, int, FlatReact) error
}