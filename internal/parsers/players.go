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
}

func (p PlayerParser) GetPlayer(id int) (*models.Player, error) {
	url := urlBuilder.NewUrlBuilder()
	url.Player()
	url.AddId(id)

	response, err := SendRequest(url.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return p.parsePlayerResponseAndId(url.String(), response)
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
	url := urlBuilder.NewUrlBuilder()
	url.PlayersStats()

	response, err := SendRequest(url.String())
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

		idStr := strings.Split(link, "/")[3]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}

		ids[i] = id
	})

	return ids, nil
}

func (p PlayerParser) parsePlayerResponseAndId(url string, response *http.Response) (*models.Player, error) {
	player, err := p.parse(response)
	if err != nil {
		return nil, err
	}

	player.Id, err = p.parseId(url)
	if err != nil {
		return nil, err
	}

	statsParser := StatsParser{}
	stats, err := statsParser.GetStats(url)
	if err != nil {
		return nil, err
	}
	player.Stats = *stats

	socialParser := SocialParser{}
	socials, err := socialParser.GetSocials(url)
	if err != nil {
		return nil, err
	}
	player.Social = socials

	return player, nil
}

func (p PlayerParser) parseId(url string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	idStr := re.FindString(url)

	id, _ := strconv.Atoi(idStr)

	return id, nil
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
