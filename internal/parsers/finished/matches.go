package finished

import (
	"hltvapi/internal/parsers"
	"hltvapi/internal/urlBuilder"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type FinishedMatchParser struct {
	parsers.MatchParser
}

func (p FinishedMatchParser) GetAllMatchesIds() ([]int, error) {
	var ids []int
	offset := 0

	for {
		url := urlBuilder.NewUrlBuilder()
		url.FinishedMatches()
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

		matchesPage := document.Find(".allres")
		matches := matchesPage.Find(".result-con")

		// Condition to break
		if matches.Length() == 0 {
			break
		}

		matches.Each(func(i int, selection *goquery.Selection) {
			aTag := selection.Find("a")
			link, ok := aTag.Attr("href")
			if !ok {
				return
			}

			id, err := strconv.Atoi(strings.Split(link, "/")[2])
			if err != nil {
				return
			}

			ids = append(ids, id)
		})

		offset += 100
	}

	return ids, nil
}
