package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// API const is the API url ;-)
const API = "https://api.chucknorris.io/jokes/random"

var (
	redisHost = "localhost"
	redisPort = "6379"
	redisUp   = false
	appPort   = "9000"
)

// Facts struct represent a chuck norris fact ;-)
type Facts struct {
	ID    string
	URL   string
	Value string
}

type redisConf struct {
	Host string
	Port string
	URL  string
	Up   bool
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

func (rc *redisConf) redisPing() {
	// fmt.Println("redisPing::Redis::URL => " + rc.URL)
	conn, err := redis.Dial("tcp", rc.URL)
	if err != nil {
		rc.Up = false
	} else {
		val, err := conn.Do("PING")
		if err != nil {
			rc.Up = false
		}
		defer conn.Close()
		if val == "PONG" {
			rc.Up = true
		}
	}
}

func (rc *redisConf) redisConfig() {
	rh := os.Getenv("REDIS_HOST")
	if rh != "" {
		redisHost = rh
	}
	rp := os.Getenv("REDIS_PORT")
	if rp != "" {
		redisPort = rp
	}
	rc.Host = redisHost
	rc.Port = redisPort
	rc.URL = redisHost + ":" + redisPort
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

func (rc *redisConf) redisRecord(f Facts) {
	rc.redisPing()
	// f = Facts{}
	if rc.Up {
		conn, err := redis.Dial("tcp", rc.URL)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("redisRecord::ID => ", f.ID)
		fmt.Println("redisRecord::Value => ", f.Value)
		fmt.Println("redisRecord::URL => ", f.URL)
		conn.Send("HMSET", f.ID, "fact", f.Value, "url", f.URL)
		conn.Flush()
		conn.Receive()
	}
}

func main() {
	ap := os.Getenv("APP_PORT")
	if ap != "" {
		appPort = ap
	}
	// f := Facts{}
	rc := redisConf{}
	rc.redisConfig()
	e := echo.New()
	e.Debug = false
	e.HideBanner = true
	e.Server.ReadTimeout = 3 * time.Second
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer
	e.File("/favicon.ico", "templates/favicon.ico")
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Pong")
	})
	e.GET("/", func(c echo.Context) error {
		f := Facts{}
		f.getOneFact()
		return c.Render(http.StatusOK, "template.html", map[string]interface{}{
			"Fact": f.Value,
			"ID":   f.ID,
			"URL":  f.URL,
		})
	}).Name = "home"

	e.POST("/record", func(c echo.Context) error {
		f := Facts{}
		rc.redisRecord(f)
		fmt.Println("Fact:" + f.Value)
		return c.Render(http.StatusOK, "record.html", map[string]interface{}{
			"Fact": f.Value,
			"ID":   f.ID,
			"URL":  f.URL,
			"Up":   rc.Up,
		})
	}).Name = "record"
	// data, err := json.MarshalIndent(e.Routes(), "", "  ")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// fmt.Println(string(data))
	e.Logger.Fatal(e.Start(":" + appPort))
}
