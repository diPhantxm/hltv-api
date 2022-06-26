package parsers

import (
	"hltvapi/internal/models"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

type SocialParser struct {
}

func (p SocialParser) GetSocials(url string) ([]models.Social, error) {
	response, err := SendRequest(url)
	if err != nil {
		return nil, err
	}

	return p.parse(response)
}

func (p SocialParser) parse(response *http.Response) ([]models.Social, error) {
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
