package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: подходящих продуктов не найдено")

type Product struct {
	UUID        string
	Name        string
	Price       float64
	Sold        int
	Category    string
	Subcategory string
	G_from      string
	G_to        string
	G_details   string
	Vendor_id   string
	Image       string
}

type Vendor struct {
	UUID        string
	Name        string
	All_sales   int
	Link        string
	Trust_level int
	Market      string
}

type Information struct {
	Product_UUID string
	Product_Name string
	Price        float64
	Sold         int
	Category     string
	Subcategory  string
	G_from       string
	G_to         string
	G_details    string
	Image        string
	Vendor_UUID  string
	Vendor_Name  string
	All_sales    int
	Link         string
	Trust_level  int
	Market       string
}

type Statistics struct {
	Date            string
	Category        string
	Subcategory     string
	Revenue         float64
	Vendor_amount   int
	Products_amount int
	Sold            int
	Avprice         float64
}
