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
	var totalSales, totalProfit float32
	var totalWeight float64

	for _, invoice := range invoices {
		for _, lineItem := range invoice.GetLineItems() {
			totalItems += lineItem.Quantity
			totalWeight += lineItem.ExtendedWeight()
		}

		// At this point, this is actually the "total cost"
		totalProfit += invoice.OrderCost()
		totalSales += invoice.OrderTotal()
	}

	totalProfit = totalSales - totalProfit

	template := template.Must(template.ParseFiles("./views/index.html"))
	data := struct {
		Invoices           []*Invoice
		TotalOrders        int
		TotalSales         float32
		TotalItems         int
		TotalWeight        float64
		TotalProfit        float32
		AverageOrderAmount float32
		AverageOrderProfit float32
	}{
		invoices,
		totalOrders,
		totalSales,
		totalItems,
		totalWeight,
		totalProfit,
		totalSales / float32(len(invoices)),
		totalProfit / float32(len(invoices)),
	}
	template.Execute(w, data)
}
