package model
import "time"

type Flat struct {
	Ref string
	Publisher string
	Address string
	Telephone string
	Price float32
	Agency bool
	Description string
	Neighbourhood string
	District string
	City string
	Province string
	State string
	Country string
	Pictures []ImgNode
	Area Area
	Rooms int
	Bathrooms int
	Floor int
	Age time.Time
	Maintenance string
	Certify string
	ComFees PriceRange
}

type Area struct {
	Built float64
	Use float64
}

type Price struct {
	Currency string
	Amount float64
}

type PriceRange struct {
	From Price
	To Price
}

type FlatReact func(*Flat) error

type FlatRepo interface {
	Save(Flat) error
	FindAllByAddress(string, FlatReact) error
	FindByAddress(string, int, int, FlatReact) error
}