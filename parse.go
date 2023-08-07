package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func parseByGoquery(resp *http.Response) ([]LiveInfo, error) {
	body := resp.Body
	requestURL := *resp.Request.URL

	var infos []LiveInfo

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("get document error: %w", err)
	}

	bandName := doc.Find(".artist-name a").Text()

	li := doc.Find(".live-list li")
	if li.Size() == 0 {
		return nil, nil
	}

	li.Each(func(_ int, s *goquery.Selection) {
		info := LiveInfo{}

		info.BandName = bandName

		info.Title = s.Find(".live-ttl").Text()

		layout := "2006/01/02"
		str := s.Find(".live-date").Text()
		str = strings.Split(str, " ")[0]
		t, _ := time.Parse(layout, str)
		info.DateTime = t

		info.Venue = s.Find(".live-venue").First().Text()

		infoURL, exists := s.Find("a").First().Attr("href")
		infoURL = strings.TrimPrefix(infoURL, "live/")
		refURL, parseErr := url.Parse(infoURL)

		if exists && parseErr == nil {
			info.URL = (*requestURL.ResolveReference(refURL)).String()
		}

		if info.Title != "" {
			infos = append(infos, info)
		}
	})

	return infos, nil
}

func parseByColly(baseURL string) ([]LiveInfo, error) {
	c := colly.NewCollector()

	var infos []LiveInfo
	var bandName string

	c.OnHTML(".artist-name a", func(e *colly.HTMLElement) {
		bandName = e.Text
	})

	c.OnHTML(".live-list li", func(e *colly.HTMLElement) {
		info := LiveInfo{}

		info.BandName = bandName

		info.Title = e.ChildText(".live-ttl")

		layout := "2006/01/02"
		str := e.ChildText(".live-date")
		str = strings.Split(str, " ")[0]
		t, _ := time.Parse(layout, str)
		info.DateTime = t

		info.Venue = e.ChildText(".live-venue")

		base, parseBaseErr := url.Parse(baseURL)
		ref, parseRefErr := url.Parse(e.ChildAttr("a", "href"))
		if ref != nil && parseBaseErr == nil && parseRefErr == nil {
			info.URL = base.ResolveReference(ref).String()
		}

		if info.Title != "" {
			infos = append(infos, info)
		}
	})

	err := c.Visit(baseURL)
	if err != nil {
		return nil, err
	}

	return infos, nil
}
