package parsers

import (
	"errors"
	"hltvapi/internal/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetMatch(url string) (*models.Match, error) {
	response, err := sendRequest(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return parseMatchResponseAndId(url, response)
}

func parseMatchResponseAndId(url string, response *http.Response) (*models.Match, error) {
	p := matchParser{}
	match, err := p.parse(response)
	if err != nil {
		return nil, err
	}

	match.Id = p.parseId(url)

	return match, nil
}

type matchParser struct {
}

func (p matchParser) parse(response *http.Response) (*models.Match, error) {
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	match := &models.Match{}

	teamA, teamB, err := p.parseTeamNames(document)
	if err != nil {
		return nil, err
	}

	match.TeamA = teamA
	match.TeamB = teamB

	time, err := p.parseStartTime(document)
	if err != nil {
		return nil, err
	}

	match.StartTime = time

	maps, err := p.parseMaps(document)
	if err != nil {
		return nil, err
	}
	match.Maps = maps

	viewers, err := p.parseViewers(document)
	if err != nil {
		return nil, err
	}
	match.Viewers = viewers

	potm, err := p.parsePotm(document)
	if err != nil {
		return nil, err
	}
	match.PlayerOfTheMatch = potm

	return match, nil
}

func (p matchParser) parseId(url string) int {
	re := regexp.MustCompile(`\d+`)
	idStr := re.FindString(url)

	id, _ := strconv.Atoi(idStr)

	return id
}

func (p matchParser) parseTeamNames(document *goquery.Document) (teamA, teamB string, err error) {
	teamNames := document.Find(".teamName")
	if teamNames == nil {
		return "", "", errors.New("teams were not found on page")
	}

	teamA = strings.ToLower(teamNames.First().Text())
	teamB = strings.ToLower(teamNames.Eq(1).Text())

	return teamA, teamB, nil
}

func (p matchParser) parseStartTime(document *goquery.Document) (time.Time, error) {
	unixTimeTag := document.Find(".time")

	if unixTimeTag == nil {
		return time.Now(), errors.New("start time was not found on page")
	}

	unixTimeStr, ok := unixTimeTag.Attr("data-unix")
	if !ok {
		return time.Now(), errors.New("data unix attribute was not found in tag")
	}

	unixTime, err := strconv.ParseInt(unixTimeStr, 10, 64)
	if err != nil {
		return time.Now(), err
	}

	t := time.UnixMilli(unixTime)
	return t, nil
}

func (p matchParser) parseMaps(document *goquery.Document) ([]string, error) {
	maps := []string{}
	mapTags := document.Find(".mapname")

	if mapTags == nil {
		return nil, errors.New("no maps were found on page")
	}

	mapTags.Each(func(t int, selection *goquery.Selection) {
		maps = append(maps, strings.ToLower(selection.Text()))
	})

	return maps, nil
}

func (p matchParser) parseViewers(document *goquery.Document) (int, error) {
	viewersTags := document.Find(".viewers")
	if viewersTags == nil {
		return 0, nil
	}

	viewers := 0

	viewersTags.Each(func(i int, selection *goquery.Selection) {
		selectionViewers, err := strconv.Atoi(selection.Text())
		if err == nil {
			viewers += selectionViewers
		}
	})

	return viewers, nil
}

// Parse Player of The Match
func (p matchParser) parsePotm(document *goquery.Document) (string, error) {
	potmTag := document.Find(".potm-container")

	if potmTag == nil {
		return "", errors.New("potm container was not found on page")
	}

	potmNicknameTag := potmTag.Find(".player-nick")

	nickname := strings.ToLower(potmNicknameTag.Text())

	return nickname, nil
}
