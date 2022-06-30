package parsers

import (
	"errors"
	"hltvapi/internal/models"
	"hltvapi/internal/urlBuilder"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PlayerParser struct {
	builder urlBuilder.UrlBuilder
}

func NewPlayerParser(builder urlBuilder.UrlBuilder) *PlayerParser {
	return &PlayerParser{
		builder: builder,
	}
}

func (p PlayerParser) GetPlayer(id int) (*models.Player, error) {
	p.builder.Clear()
	p.builder.Player()
	p.builder.AddId(id)

	response, err := SendRequest(p.builder.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	player, err := p.parsePlayerResponse(response)
	if err != nil {
		return nil, err
	}

	player.Id = id

	statsParser := StatsParser{}
	stats, err := statsParser.GetStats(p.builder.String())
	if err != nil {
		return nil, err
	}
	player.Stats = *stats

	socialParser := SocialParser{}
	socials, err := socialParser.GetSocials(p.builder.String())
	if err != nil {
		return nil, err
	}
	player.Social = socials

	return player, nil
}

func (p PlayerParser) GetPlayers() ([]models.Player, error) {
	ids, err := p.GetAllPlayersIds()
	if err != nil {
		return nil, err
	}

	result := make([]models.Player, len(ids))
	for i, id := range ids {
		player, err := p.GetPlayer(id)
		if err != nil {
			continue
		}

		result[i] = *player
	}

	return result, nil
}

func (p PlayerParser) GetAllPlayersIds() ([]int, error) {
	p.builder.Clear()
	p.builder.PlayersStats()

	response, err := SendRequest(p.builder.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	matchTags := document.Find(".stats-table").Find("tbody").Find("tr")
	ids := make([]int, matchTags.Length())

	matchTags.Each(func(i int, selection *goquery.Selection) {
		link, ok := selection.Find(".playerCol").Find("a").Attr("href")
		if !ok {
			return
		}

		re := regexp.MustCompile(`\/\w+\/(\d+)`)
		idStr := re.FindStringSubmatch(link)[1]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}

		ids[i] = id
	})

	return ids, nil
}

func (p PlayerParser) parsePlayerResponse(response *http.Response) (*models.Player, error) {
	player, err := p.parse(response)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (p PlayerParser) parse(response *http.Response) (*models.Player, error) {
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

func (p PlayerParser) parseAge(document *goquery.Document) (int, error) {
	ageTag := document.Find(".playerAge").Find("span").Last()

	ageStr := strings.Split(ageTag.Text(), " ")
	age, err := strconv.Atoi(ageStr[0])
	if err != nil {
		return 0, errors.New("age was in incorrect format")
	}

	return age, nil
}

func (p PlayerParser) parseNickname(document *goquery.Document) (string, error) {
	nicknameTag := document.Find(".playerNickname")

	return nicknameTag.Text(), nil
}

func (p PlayerParser) parseTeam(document *goquery.Document) (string, error) {
	teamTag := document.Find(".playerTeam").Find("span").Last()

	team := teamTag.Text()

	return team, nil
}

func (p PlayerParser) parseName(document *goquery.Document) (string, string, error) {
	nameTag := document.Find(".playerRealname").Text()

	nameSplit := strings.Split(nameTag, " ")[1:]

	if len(nameSplit) < 2 {
		return "", "", errors.New("error in parsing full name")
	}

	return nameSplit[0], nameSplit[1], nil
}

func (p PlayerParser) parseCountry(document *goquery.Document) (string, error) {
	flagTag := document.Find(".flag")

	country, ok := flagTag.Attr("alt")
	if !ok {
		return "", errors.New("flag/country was not found on page")
	}

	return country, nil
}
