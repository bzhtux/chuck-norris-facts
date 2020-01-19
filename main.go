package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// API const is the API url ;-)
const API = "https://api.chucknorris.io/jokes/random"

// Facts struct represent a chuck norris fact ;-)
type Facts struct {
	ID      string
	IconURL string
	URL     string
	Value   string
	Updated string
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render func help to render html page
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func (f *Facts) getOneFact() {
	resp, err := http.Get(API)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	byt := []byte(string(body))
	e := json.Unmarshal(byt, &f)
	if e != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	f := Facts{}
	e := echo.New()
	e.Use(middleware.Logger())
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("template.html")),
	}
	e.Renderer = renderer
	// Named route "foobar"
	e.GET("/", func(c echo.Context) error {
		f.getOneFact()
		return c.Render(http.StatusOK, "template.html", map[string]interface{}{
			"Fact": f.Value,
			"ID":   f.ID,
			"URL":  f.URL,
		})
	}).Name = "home"

	e.Logger.Fatal(e.Start(":9000"))
}
