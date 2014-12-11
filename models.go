package main

import (
	"errors"
	"time"
)

// Notes:
// If you're reading this on my GitHub, this schema wasn't designed by me
// The "Number" fields should be identifiers but that's how it is in the original schema
// Error handling is pretty poor since this is just running a static report

type Customer struct {
	Number int
	Name   string
	Address
}

func GetCustomer(id int) (*Customer, error) {
	rows, err := db.Query("SELECT * FROM Customer WHERE CustNum = ?", id)

	customer := Customer{}
	if !rows.Next() {
		return nil, errors.New("No such customer")
	}

	err = rows.Scan(&customer.Number, &customer.Name, &customer.Street, &customer.City, &customer.State, &customer.Zip)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

type Product struct {
	Sku         string
	SupplierNum int
	Description string
	Qoh         int
	Cost        float32
	UnitPrice   float32
	UnitWeight  float64
}

func GetProduct(sku string) (*Product, error) {
	rows, err := db.Query("SELECT * FROM Inventory WHERE SKU = ?", sku)

	if err != nil {
		return nil, err
	}

	product := Product{}

	if !rows.Next() {
		return nil, errors.New("No such product")
	}

	err = rows.Scan(&product.Sku, &product.SupplierNum, &product.Description, &product.Qoh, &product.Cost, &product.UnitPrice, &product.UnitWeight)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

type Invoice struct {
	OrderNumber    int
	CustomerNumber int
	OrderDate      time.Time
	Status         string
	lineItems      *[]*InvoiceLineItem
	customer       *Customer
}

func (i *Invoice) GetCustomer() *Customer {
	if i.customer == nil {
		customer, err := GetCustomer(i.CustomerNumber)
		if err != nil {
			Log.Error(err)
			return nil
		}

		i.customer = customer
	}

	return i.customer
}

func (i *Invoice) OrderTotal() float32 {
	var total float32
	for _, lineItem := range i.GetLineItems() {
		total += lineItem.ExtendedPrice()
	}

	return total
}

func (i *Invoice) OrderCost() float32 {
	var cost float32
	for _, lineItem := range i.GetLineItems() {
		cost += lineItem.ExtendedCost()
	}

	return cost
}

func (i *Invoice) OrderProfit() float32 {
	return i.OrderTotal() - i.OrderCost()
}

func GetAllInvoices() []*Invoice {
	rows, err := db.Query("SELECT COUNT(*) FROM Invoice")
	if err != nil {
		Log.Error(err)
		return []*Invoice{}
	}

	var count int
	rows.Scan(&count)

	rows, err = db.Query("SELECT i.*, c.* FROM Invoice i NATURAL JOIN Customer c")
	if err != nil {
		Log.Error(err)
		return []*Invoice{}
	}

	invoices := make([]*Invoice, 0, count)
	for rows.Next() {
		invoice := &Invoice{}
		customer := &Customer{}
		err := rows.Scan(&invoice.OrderNumber, &invoice.CustomerNumber, &invoice.OrderDate, &invoice.Status,
			&customer.Number, &customer.Name, &customer.Street, &customer.City, &customer.State, &customer.Zip)

		invoice.customer = customer

		if err != nil {
			Log.Error(err)
			return []*Invoice{}
		}

		invoices = append(invoices, invoice)
	}

	return invoices
}

// Gets the line items associated with this Invoice
func (i *Invoice) GetLineItems() []*InvoiceLineItem {
	if i.lineItems == nil {
		rows, err := db.Query("SELECT li.*, p.* FROM InvoiceLineItem li NATURAL JOIN Inventory p WHERE li.OrderNum = ?", i.OrderNumber)
		// Errors shouldn't happen here....
		if err != nil {
			Log.Error(err)
			return nil
		}

		// This should probably be pointers to the InvoiceLineItem, but there aren't many so copying
		// doesn't hurt *that* bad
		items := make([]*InvoiceLineItem, 0)

		for rows.Next() {
			lineItem := &InvoiceLineItem{}
			product := &Product{}
			err := rows.Scan(&lineItem.OrderNumber, &lineItem.LineNumber, &lineItem.Sku, &lineItem.Quantity,
				&product.Sku, &product.SupplierNum, &product.Description, &product.Qoh, &product.Cost,
				&product.UnitPrice, &product.UnitWeight)

			lineItem.product = product

			if err != nil {
				Log.Error(err)
			}

			items = append(items, lineItem)
		}

		i.lineItems = &items
	}

	return *(i.lineItems)
}

type InvoiceLineItem struct {
	OrderNumber int
	LineNumber  int
	Sku         string
	Quantity    int
	product     *Product
}

func (i *InvoiceLineItem) GetProduct() *Product {
	if i.product == nil {
		product, err := GetProduct(i.Sku)
		if err != nil {
			Log.Error(err)
			return nil
		}

		i.product = product
	}

	return i.product
}

// Use a pointer in case the product hasn't yet been looked up
func (i *InvoiceLineItem) ExtendedCost() float32 {
	return i.GetProduct().Cost * float32(i.Quantity)
}

func (i *InvoiceLineItem) ExtendedProfit() float32 {
	return i.ExtendedPrice() - i.ExtendedCost()
}

func (i *InvoiceLineItem) ExtendedPrice() float32 {
	return i.GetProduct().UnitPrice * float32(i.Quantity)
}

func (i *InvoiceLineItem) ExtendedWeight() float64 {
	return i.GetProduct().UnitWeight * float64(i.Quantity)
}

type Supplier struct {
	Number      int
	CompanyName string
	Contact     string
	Phone       string
	Address
}
