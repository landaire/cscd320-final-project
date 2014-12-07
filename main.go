package main

import (
	"html/template"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	nlogrus "github.com/meatballhat/negroni-logrus"
)

var Log *logrus.Logger

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	logger := nlogrus.NewMiddleware()
	Log = logger.Logger

	n := negroni.Classic()
	n.Use(logger)
	n.UseHandler(mux)
	n.Run(":8000")
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	invoices := GetAllInvoices()

	totalOrders := len(invoices)
	var totalItems int
	var totalSales float32
	var totalWeight float64

	for _, invoice := range invoices {
		for _, lineItem := range invoice.GetLineItems() {
			totalItems += lineItem.Quantity
			totalWeight += lineItem.ExtendedWeight()
		}

		totalSales += invoice.OrderTotal()
	}

	template := template.Must(template.ParseFiles("./views/index.html"))
	data := struct {
		Invoices    []*Invoice
		TotalOrders int
		TotalSales  float32
		TotalItems  int
		TotalWeight float64
	}{
		invoices,
		totalOrders,
		totalSales,
		totalItems,
		totalWeight,
	}
	template.Execute(w, data)
	//	invoice := Invoice{}
	//	invoice.OrderNumber = 1920
	//
	//	Log.Printf("%#v\n", invoice.GetLineItems())
}
