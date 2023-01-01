package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type dElement struct {
	Title         string
	Votes         string
	Answers       string
	Views         string
	Views_per_day string
	Link          string
}

func get(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func ParsePage(url string) []dElement {
	var data []dElement
	body := get(url)

	// class = "s-post-summary js-post-summary"
	body.Find(".s-post-summary").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".s-post-summary--content-title").Find("a").Text()
		statsEl := s.Find(".s-post-summary--stats.js-post-summary-stats")
		link, _ := s.Find(".s-post-summary--content-title").Find("a").Attr("href")
		link = "https://stackoverflow.com" + link

		var stats []string
		statsEl.Find(".s-post-summary--stats-item").Each(func(i int, r *goquery.Selection) {
			stat := r.Find(".s-post-summary--stats-item-number").Text()
			stats = append(stats, stat)
		})
		votes := stats[0]
		answers := stats[1]
		views := stats[2]

		if views[len(views)-1:] == "k" {
			views = strconv.Itoa(int(convertViewsToNumber(views[:len(views)-1]) * 1000))
		} else if views[len(views)-1:] == "m" {
			views = strconv.Itoa(int(convertViewsToNumber(views[:len(views)-1]) * 1000000))
		}

		date := convertDateToUnix(s.Find(".relativetime").AttrOr("title", "no date"))
		now := time.Now().Unix()
		views_per_day := strconv.Itoa(convertToInt(views) / int((now-date)/86400))
		data = append(data, dElement{
			Title:         title,
			Votes:         votes,
			Answers:       answers,
			Views:         views,
			Views_per_day: views_per_day,
			Link:          link,
		})
	})

	return data

}

func get_all_pages(tag string) []dElement {
	var data []dElement
	nUrl := fmt.Sprintf("https://stackoverflow.com/questions/tagged/%s?tab=frequent&page=1&pagesize=50", tag)

	body := get(nUrl)

	links := gen_links(convertToInt(body.Find(".s-pagination").Find("a").Last().Prev().Text()), tag)
	// first 3 pages
	links = links[:3]

	dataChan := make(chan []dElement, len(links))
	var wg sync.WaitGroup
	wg.Add(len(links))

	for _, link := range links {
		go func(link string) {
			defer wg.Done()

			dataChan <- ParsePage(link)
			println(dataChan)
		}(link)
	}
	wg.Wait()
	close(dataChan)
	var mdata []dElement
	for n := range dataChan {
		mdata = append(mdata, n...)
	}
	fmt.Println(len(mdata))
	var finalArr []dElement
	for _, v := range mdata {
		if v.Title != "" {
			finalArr = append(finalArr, v)
		}
	}

	to_csv(finalArr, tag)
	return data

}

func to_csv(data []dElement, tag string) {

	file, err := os.Create(tag + "_data.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	file.WriteString("Title,Votes,Answers,Views,Views_per_day,Link \n")

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		line := []string{value.Title, value.Votes, value.Answers, value.Views, value.Views_per_day, value.Link}
		err := writer.Write(line)
		checkError("Cannot write to file", err)
	}

}

func gen_links(pages int, tag string) []string {
	var links []string
	for i := 1; i <= pages; i++ {
		tem_url := fmt.Sprintf("https://stackoverflow.com/questions/tagged/%s?tab=frequent&page=%d&pagesize=50", tag, i)
		links = append(links, tem_url)
	}
	return links
}
