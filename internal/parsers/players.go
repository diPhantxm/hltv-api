package parsers

import (
	"errors"
	"hltvapi/internal/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetPlayer(url string) (*models.Player, error) {
	response, err := sendRequest(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return parsePlayerResponseAndId(url, response)
}

func parsePlayerResponseAndId(url string, response *http.Response) (*models.Player, error) {
	p := playerParser{}

	player, err := p.parse(response)
	if err != nil {
		return nil, err
	}

	player.Id, err = p.parseId(url)
	if err != nil {
		return nil, err
	}

	stats, err := GetStats(url)
	if err != nil {
		return nil, err
	}
	player.Stats = *stats

	socials, err := GetSocials(url)
	if err != nil {
		return nil, err
	}
	player.Social = socials

	return player, nil
}

type playerParser struct {
}

func (p playerParser) parseId(url string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	idStr := re.FindString(url)

	id, _ := strconv.Atoi(idStr)

	return id, nil
}

func (p playerParser) parse(response *http.Response) (*models.Player, error) {
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	player := &models.Player{}

	player.Age, err = p.parseAge(document)
	if err != nil {
		return nil, err
	}

	player.Nickname, err = p.parseNickname(document)
	if err != nil {
		return nil, err
	}

	player.Team, err = p.parseTeam(document)
	if err != nil {
		return nil, err
	}

	player.FirstName, player.LastName, err = p.parseName(document)
	if err != nil {
		return nil, err
	}

	player.Country, err = p.parseCountry(document)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (p playerParser) parseAge(document *goquery.Document) (int, error) {
	ageTag := document.Find(".playerAge").Find("span").Last()

	ageStr := strings.Split(ageTag.Text(), " ")
	age, err := strconv.Atoi(ageStr[0])
	if err != nil {
		return 0, errors.New("age was in incorrect format")
	}

	return age, nil
}

func (p playerParser) parseNickname(document *goquery.Document) (string, error) {
	nicknameTag := document.Find(".playerNickname")

	return nicknameTag.Text(), nil
}

func (p playerParser) parseTeam(document *goquery.Document) (string, error) {
	teamTag := document.Find(".playerTeam").Find("span").Last()

	team := teamTag.Text()

	return team, nil
}

func (p playerParser) parseName(document *goquery.Document) (string, string, error) {
	nameTag := document.Find(".playerRealname").Text()

	nameSplit := strings.Split(nameTag, " ")[1:]

	if len(nameSplit) < 2 {
		return "", "", errors.New("error in parsing full name")
	}

	return nameSplit[0], nameSplit[1], nil
}

func (p playerParser) parseCountry(document *goquery.Document) (string, error) {
	flagTag := document.Find(".flag")

	country, ok := flagTag.Attr("alt")
	if !ok {
		return "", errors.New("flag/country was not found on page")
	}

	return country, nil
}
