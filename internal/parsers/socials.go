package parsers

import (
	"hltvapi/internal/models"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func GetSocials(url string) ([]models.Social, error) {
	response, err := sendRequest(url)
	if err != nil {
		return nil, err
	}

	p := socialParser{}
	return p.parse(response)
}

type socialParser struct {
}

// TODO: somehow get social names
func (p socialParser) parse(response *http.Response) ([]models.Social, error) {
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	socialsTag := document.Find(".socialMediaButtons")

	links := socialsTag.Find("a")
	socials := make([]models.Social, links.Length())
	namePattern := regexp.MustCompile(`www\.(\w+)\.`)
	links.Each(func(i int, selection *goquery.Selection) {
		socials[i].Link, _ = selection.Attr("href")
		socials[i].Name = namePattern.FindStringSubmatch(socials[i].Link)[1]
	})

	return socials, nil
}
