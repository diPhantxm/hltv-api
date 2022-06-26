package parsers

import (
	"hltvapi/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type StatsParser struct {
}

func (p StatsParser) GetStats(url string) (*models.Statistics, error) {
	response, err := SendRequest(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return p.parse(response)
}

func (p StatsParser) parse(response *http.Response) (*models.Statistics, error) {
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	stats := &models.Statistics{}

	statsTags := document.Find(".player-stat")

	statsTags.Each(func(i int, selection *goquery.Selection) {
		name := strings.ToLower(selection.Find("b").Text())

		valueStr := strings.ReplaceAll(selection.Find(".statsVal").Text(), "%", "")
		value, err := strconv.ParseFloat(valueStr, 32)
		if err != nil {
			return
		}

		if name == "rating 2.0" {
			stats.Rating = float32(value)
		}
		if name == "kills per round" {
			stats.KillsPerRound = float32(value)
		}
		if name == "headshots" {
			stats.Headshots = float32(value)
		}
		if name == "maps played" {
			stats.MapsPlayed = int(value)
		}
		if name == "deaths per round" {
			stats.DeathsPerRound = float32(value)
		}
		if name == "rounds contributed" {
			stats.RoundsContributed = float32(value)
		}
	})

	return stats, nil
}
