package main

import (
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func parse(baseURL string) ([]Live, error) {
	collector := colly.NewCollector()

	var lives []Live
	var band string

	collector.OnHTML(".artist-name a", func(e *colly.HTMLElement) {
		band = e.Text
	})

	collector.OnHTML(".live-list li", func(e *colly.HTMLElement) {
		dateStr := e.ChildText(".live-date")
		date, err := parseDate(dateStr)
		if err != nil {
			return
		}

		liveURL, err := resolveURL(baseURL, e.ChildAttr("a", "href"))
		if err != nil {
			return
		}

		lives = append(lives, Live{
			Band:      band,
			Title:     e.ChildText(".live-ttl"),
			EventDate: date,
			Venue:     e.ChildText(".live-venue"),
			URL:       liveURL,
		})
	})

	if err := collector.Visit(baseURL); err != nil {
		return nil, err
	}

	return lives, nil
}

func parseDate(dateStr string) (time.Time, error) {
	layout := "2006/01/02"
	parsedDateStr := strings.Split(dateStr, " ")[0]
	return time.Parse(layout, parsedDateStr)
}

func resolveURL(base, ref string) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	refURL, err := url.Parse(ref)
	if err != nil {
		return "", err
	}

	return baseURL.ResolveReference(refURL).String(), nil
}
