package main

import (
	"errors"
	"html/template"
	"net/http"

	"deeptown.com/deepsearch/pkg/models"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (app *application) landing(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/index.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
}

func (app *application) searchProduct(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	s, err := app.products.SearchProduct(name)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	t, err_ := app.products.GetVendor(s)
	if err_ != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err_)
		}
		return
	}
	files := []string{
		"./ui/html/search2.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, t)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) showStat(w http.ResponseWriter, r *http.Request) {
	var Subcategory []string
	var Revenue []opts.BarData
	var Vendor_amount []opts.BarData
	var Products_amount []opts.BarData
	var Sold []opts.BarData
	var Avprice []opts.BarData
	name := "Darkfox"
	s, err := app.products.Stat(name)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	for _, element := range s {
		Subcategory = append(Subcategory, element.Subcategory)
		Revenue = append(Revenue, opts.BarData{Value: element.Revenue, Name: element.Subcategory})
		Vendor_amount = append(Vendor_amount, opts.BarData{Value: element.Vendor_amount, Name: element.Subcategory})
		Products_amount = append(Products_amount, opts.BarData{Value: element.Products_amount, Name: element.Subcategory})
		Sold = append(Sold, opts.BarData{Value: element.Sold, Name: element.Subcategory})
		Avprice = append(Avprice, opts.BarData{Value: element.Avprice, Name: element.Subcategory})
	}

	bar := charts.NewBar()
	bar.Renderer = newSnippetRenderer(bar, bar.Validate)
	bar.SetGlobalOptions(
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "900px",
			Height: "700px",
		}),
		charts.WithColorsOpts(opts.Colors{"#4361ee"}),
		charts.WithTitleOpts(opts.Title{Title: "Revenue by subcategories"}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "slider",
			Start: 0,
			End:   20,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "inside",
			Start: 0,
			End:   100,
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "Don't forget to donate!",
				},
			}},
		),
	)

	bar.SetXAxis(Subcategory).
		AddSeries("Revenue", Revenue).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      true,
				Position:  "top",
				Formatter: "{c}$",
			}))
	// bar.XYReversal()
	var htmlSnippet template.HTML = renderToHtml(bar)

	bar_ := charts.NewBar()
	bar_.Renderer = newSnippetRenderer(bar_, bar_.Validate)
	//"#4cc9f0"
	bar_.SetGlobalOptions(charts.WithInitializationOpts(
		opts.Initialization{
			Width:  "900px",
			Height: "700px",
		}),
		charts.WithColorsOpts(opts.Colors{"#4cc9f0"}),
		charts.WithTitleOpts(opts.Title{Title: "Products/Vendors amount by subcategories"}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "slider",
			Start: 0,
			End:   10,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "inside",
			Start: 10,
			End:   50,
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "Don't forget to donate",
				},
			}},
		),
		charts.WithLegendOpts(opts.Legend{Show: true, Right: "25%"}),
	)
	bar_.SetXAxis(Subcategory).
		AddSeries("Vendors", Vendor_amount).
		AddSeries("Products", Products_amount).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "right",
			}))
	// bar_.XYReversal()

	var htmlSnippet_ template.HTML = renderToHtml(bar_)

	tplVars := map[string]interface{}{
		"Html":  htmlSnippet,
		"Html_": htmlSnippet_,
	}
	ts, err := template.ParseFiles("./ui/html/stat.page.tmpl", "./ui/html/base.layout.tmpl", "./ui/html/footer.partial.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, tplVars)
	if err != nil {
		app.serverError(w, err)
	}
}
