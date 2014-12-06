package main

import "time"

// Notes:
// If you're reading this on my GitHub, this schema wasn't designed by me
// The "Number" fields should be identifiers but that's how it is in the original schema

type Customer struct {
	Number int
	Name   string
	Address
}

type Address struct {
	Street string
	City   string
	State  string
}

type Inventory struct {
	Sku         string
	SupplierNum int
	Description string
	Qoh         int
	Cost        float32
	UnitPrice   float32
	UnitWeight  float64
}

type Invoice struct {
	OrderNumber    int
	CustomerNumber int
	OrderDate      time.Time
	Status         string
	lineItems      *[]InvoiceLineItem
}

// Gets the line items associated with this Invoice
func (i *Invoice) GetLineItems() []InvoiceLineItem {
	if i.lineItems == nil {
		rows, err := db.Query("SELECT li.* FROM InvoiceLineItem li WHERE li.OrderNum = ?", i.OrderNumber)
		// Errors shouldn't happen here....
		if err != nil {
			Log.Error(err)
			return nil
		}

		// This should probably be pointers to the InvoiceLineItem, but there aren't many so copying
		// doesn't hurt *that* bad
		items := make([]InvoiceLineItem, 0)

		for rows.Next() {
			lineItem := InvoiceLineItem{}
			err := rows.Scan(&lineItem.OrderNumber, &lineItem.LineNumber, &lineItem.Sku, &lineItem.Quantity)

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
}

type Supplier struct {
	Number      int
	CompanyName string
	Contact     string
	Phone       string
	Address
}
