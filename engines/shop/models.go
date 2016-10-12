package shop

import (
	"github.com/itpkg/champak/engines/auth"
	"github.com/itpkg/champak/web"
)

//Product product
type Product struct {
	web.Model

	Sku         string
	Name        string
	Attachments []auth.Attachment
	Summary     string
	Detail      string

	RegularPrice float64
	SpecialPrice float64

	Catalog Catalog
	Factory Factory
	Options []Option
}

//Catalog catalog
type Catalog struct {
	web.Model

	Name        string
	Attachments []auth.Attachment
	Detail      string

	Products []Product
}

//Factory factory
type Factory struct {
	web.Model

	Name        string
	Attachments []auth.Attachment
	Detail      string

	Products []Product
}

//Vendor vendor
type Vendor struct {
	Name    string
	Address string
	Details string
}

//Option option for product
type Option struct {
	Name  string
	Price float64

	Product Product
}

//Shipment shipment
type Shipment struct {
	Name  string
	Price float64
}
