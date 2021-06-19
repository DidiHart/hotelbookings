package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/DidiHart/hotelbookings/internal/config"
	"github.com/DidiHart/hotelbookings/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"humanDate": HumanDate,
}
var app *config.AppConfig
var pathToTemplates = "./templates"

//NewTemplates set config for the template pkg
func NewRenderer(a *config.AppConfig) {
	app = a
}

//returns time in human readable format
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

//Adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	//use any data you want to add to every page
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

//Template renders templates using html/template
func Template(w http.ResponseWriter, r *http.Request, html string, td *models.TemplateData) {
	// _, err := CacheTemplate(w)
	// if err != nil {
	// 	fmt.Println("Error getting template cache", err)
	// }

	//get template from app config
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache

	} else {
		tc, _ = CacheTemplate()
	}

	t, ok := tc[html]
	if !ok {
		log.Fatal("couldn't get template from template cache")
	}

	b := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(b, td)

	_, err := b.WriteTo(w)

	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}

	// parsedTemplate, _ := template.ParseFiles("./templates/" + html)
	// err := parsedTemplate.Execute(w, nil)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}

//create a template cacher as a map
func CacheTemplate() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// fmt.Println("Pages is currently", page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matched, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))

		if err != nil {
			return myCache, err
		}

		if len(matched) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
