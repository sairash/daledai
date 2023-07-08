package main

import (
	// "encoding/json"
	"fmt"

	// "io/ioutil"

	"net/http"

	"github.com/foolin/goview/supports/echoview"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/recoilme/pudge"
)

type lists struct {
	Proxies []proxies `json:"proxies"`
}

type proxies struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
	Type string `json:"type"`
}

const defaultAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"

// port := fmt.Sprintf("%d", list.Proxies[1].Port)
var search_opt = SearchOptions{
	CountryCode: "np",
	// ProxyAddr:    list.Proxies[1].Type + "://" + list.Proxies[1].Ip + ":" + port,
	UserAgent:    defaultAgent,
	LanguageCode: "en",
	Start:        0,
}

func main() {
	defer pudge.CloseAll()

	// jsonFile, err := os.Open("http.json")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Successfully Opened http.json")
	// byteValue, _ := ioutil.ReadAll(jsonFile)

	// var list lists

	// json.Unmarshal(byteValue, &list)
	// jsonFile.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = echoview.Default()
	e.Static("/static", "assets")
	e.File("/favicon.ico", "assets/favicon.ico")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", echo.Map{
			"title": "Dale Dai",
			"add": func(a int, b int) int {
				return a + b
			},
		})
	})

	e.GET("/search", func(c echo.Context) error {
		q := c.Request().URL.Query().Get("q")

		result := []Result{}
		next_page_url := ""
		err := pudge.Get("/home/sairash/search/db", q, &result)

		if err != nil {

			result, next_page_url, err = Search(q, search_opt)
			if err != nil {
				fmt.Println(err)
			}
			err = pudge.Set("/home/sairash/search/db", q, result)
			if err != nil {
				fmt.Println(err)
			}
		}

		return c.Render(http.StatusOK, "search", echo.Map{
			"title":         "Dale Search",
			"result":        result,
			"next_page_url": next_page_url,
			"q":             q})
	})

	e.Logger.Fatal(e.Start(":9090"))

}
