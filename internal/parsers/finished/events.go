package finished

import (
	"hltvapi/internal/parsers"
	"hltvapi/internal/urlBuilder"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type FinishedEventParser struct {
	parsers.EventParser
}

func (p FinishedEventParser) GetAllEventsIds() ([]int, error) {
	var ids []int
	offset := 0

	for {
		url := urlBuilder.NewUrlBuilder()
		url.FinishedEvents()
		url.AddParam("offset", strconv.Itoa(offset))

		time.Sleep(time.Duration(400) * time.Millisecond)

		response, err := parsers.SendRequest(url.String())
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			return nil, err
		}

		eventsPage := document.Find(".events-page")
		events := eventsPage.Find(".small-event")

		// Condition to break
		if events.Length() == 0 {
			break
		}

		events.Each(func(i int, selection *goquery.Selection) {
			link, ok := selection.Attr("href")
			if !ok {
				return
			}

			id, err := strconv.Atoi(strings.Split(link, "/")[2])
			if err != nil {
				return
			}

			ids = append(ids, id)
		})

		offset += 50
	}

	return ids, nil
}
