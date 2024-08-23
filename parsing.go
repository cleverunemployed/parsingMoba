package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sclevine/agouti"
)

type Catalog struct {
	Title, Link string
}

type Good struct {
	Name string
}

func parseCatalogs() []Catalog {
	var url string = "https://moba.ru/catalog/"

	var catalogs []Catalog

	c := colly.NewCollector()

	// Set custom headers using the OnRequest callback
	c.OnRequest(func(req *colly.Request) {
		req.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Headers.Set("accept-language", "ru,en;q=0.9,en-GB;q=0.8,en-US;q=0.7")
		req.Headers.Set("cache-control", "max-age=0")
		req.Headers.Set("cookie", "BITRIX_SM_GUEST_ID=82618509; BITRIX_SM_CITY_ID=9706; BITRIX_CONVERSION_CONTEXT_s1=%7B%22ID%22%3A7%2C%22EXPIRE%22%3A1724432340%2C%22UNIQUE%22%3A%5B%22conversion_visit_day%22%5D%7D; PHPSESSID=m4ShdRR95pUG2aCqvlsNEVw6LbMuMqkn; BITRIX_SM_LAST_VISIT=23.08.2024%2022%3A42%3A54; BITRIX_SM_mb_h_=ab38a5547e8503ed20aae2204d1c5410; BITRIX_SM_mb_key_=e5a16d493ccf89cd971357c29bb0c505; BITRIX_SM_mb_key2_=b8d02a9171b224067f02887f13a47f2a; _ym_debug=null")
		req.Headers.Set("priority", "u=0, i")
		req.Headers.Set("referer", "https://moba.ru/catalog/")
		req.Headers.Set("sec-ch-ua", `"Chromium";v="128", "Not;A=Brand";v="24", "Microsoft Edge";v="128"`)
		req.Headers.Set("sec-ch-ua-mobile", "?0")
		req.Headers.Set("sec-ch-ua-platform", `"Windows"`)
		req.Headers.Set("sec-fetch-dest", "document")
		req.Headers.Set("sec-fetch-mode", "navigate")
		req.Headers.Set("sec-fetch-site", "same-origin")
		req.Headers.Set("sec-fetch-user", "?1")
		req.Headers.Set("upgrade-insecure-requests", "1")
	})

	// Define a callback for the collected data
	c.OnHTML("li.name", func(e *colly.HTMLElement) {
		link, _ := e.DOM.Find("a.dark_link").Attr("href")
		catalogs = append(catalogs, Catalog{Title: strings.Replace(e.Text, "\n", "", -1), Link: "https://moba.ru" + link})
	})

	// Define a callback for errors
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping
	err := c.Visit(url)
	if err != nil {
		fmt.Println("Visit failed:", err)
	}

	return catalogs
}

func saveCatalogs(catalogs []Catalog) {
	file, _ := json.MarshalIndent(catalogs, "", "\t")

	_ = ioutil.WriteFile("catalogs.json", file, 0644)
}

func parseGoods(catalog Catalog) {

	// driver := agouti.PhantomJS()
	// driver := agouti.ChromeDriver()
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{"--headless", "--disable-gpu", "--no-sandbox"}),
	)

	if err := driver.Start(); err != nil {
		log.Fatal("Failed to start driver:", err)
	}

	page, err := driver.NewPage()
	if err != nil {
		log.Fatal("Failed to open page:", err)
	}

	if err := page.Navigate("https://agouti.org/"); err != nil {
		log.Fatal("Failed to navigate:", err)
	}

	sectionTitle, err := page.FindByID(`getting-agouti`).Text()
	log.Println(sectionTitle)

	if err := driver.Stop(); err != nil {
		log.Fatal("Failed to close pages and stop WebDriver:", err)
	}

}

func main() {
	parseGoods(Catalog{})
}
