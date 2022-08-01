package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.landing)
	//mux.HandleFunc("/home", app.home)
	mux.HandleFunc("/search", app.searchProduct)
	mux.HandleFunc("/stat", app.showStat)
	//mux.HandleFunc("/product", app.showProduct)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
