package parsers

import (
	"errors"
	"hltvapi/internal/models"
	"hltvapi/internal/urlBuilder"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type MatchParser struct {
}

func (p MatchParser) GetMatch(id int) (*models.Match, error) {
	url := urlBuilder.NewUrlBuilder()
	url.Match()
	url.AddId(id)

	response, err := SendRequest(url.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return p.parseMatchResponseAndId(url.String(), response)
}

func (p MatchParser) GetMatches() ([]models.Match, error) {
	ids, err := p.GetUpcomingMatchesIds()
	if err != nil {
		return nil, err
	}

	result := make([]models.Match, len(ids))
	for i, id := range ids {
		match, err := p.GetMatch(id)
		if err != nil {
			continue
		}

		result[i] = *match
	}

	return result, nil
}

func (p MatchParser) GetMatchesByDate(date string) ([]models.Match, error) {
	url := urlBuilder.NewUrlBuilder()
	url.Results()
	url.AddParam("startDate", date)
	url.AddParam("endDate", date)

	ids, err := p.getMatchesIdsByDate(url.String())
	if err != nil {
		return nil, err
	}

	result := make([]models.Match, len(ids))
	for i, id := range ids {
		match, err := p.GetMatch(id)
		if err != nil {
			continue
		}

		result[i] = *match
	}

	return result, nil
}

func (p MatchParser) GetUpcomingMatchesIds() ([]int, error) {
	url := urlBuilder.NewUrlBuilder()
	url.Match()

	response, err := SendRequest(url.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	matchTags := document.Find(".upcomingMatchesSection")
	ids := make([]int, matchTags.Length())

	matchTags.Each(func(i int, selection *goquery.Selection) {
		link, ok := selection.Find(".match").Attr("href")
		if !ok {
			return
		}

		idStr := strings.Split(link, "/")[2]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}

		ids[i] = id
	})

	return ids, nil
}

func (p MatchParser) getMatchesIdsByDate(url string) ([]int, error) {
	response, err := SendRequest(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	idTags := document.Find(".results-all").Find("a")
	ids := make([]int, idTags.Length())

	idTags.Each(func(i int, selection *goquery.Selection) {
		link, ok := selection.Attr("href")
		if !ok {
			return
		}

		id, err := strconv.Atoi(strings.Split(link, "/")[2])
		if err != nil {
			return
		}

		ids[i] = id
	})

	return ids, nil
}

func (p MatchParser) parseMatchResponseAndId(url string, response *http.Response) (*models.Match, error) {
	match, err := p.parse(response)
	if err != nil {
		return nil, err
	}

	match.Id = p.parseId(url)

	return match, nil
}

func (p MatchParser) parse(response *http.Response) (*models.Match, error) {
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

func (p MatchParser) parseId(url string) int {
	re := regexp.MustCompile(`\d+`)
	idStr := re.FindString(url)

	id, _ := strconv.Atoi(idStr)

	return id
}

func (p MatchParser) parseTeamNames(document *goquery.Document) (teamA, teamB string, err error) {
	teamNames := document.Find(".teamName")
	if teamNames == nil {
		return "", "", errors.New("teams were not found on page")
	}

	teamA = strings.ToLower(teamNames.First().Text())
	teamB = strings.ToLower(teamNames.Eq(1).Text())

	return teamA, teamB, nil
}

func (p MatchParser) parseStartTime(document *goquery.Document) (time.Time, error) {
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

func (p MatchParser) parseMaps(document *goquery.Document) ([]string, error) {
	maps := []string{}
	mapTags := document.Find(".mapholder").Find(".mapname")

	if mapTags == nil {
		return nil, errors.New("no maps were found on page")
	}

	mapTags.Each(func(t int, selection *goquery.Selection) {
		maps = append(maps, strings.ToLower(selection.Text()))
	})

	return maps, nil
}

func (p MatchParser) parseViewers(document *goquery.Document) (int, error) {
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
func (p MatchParser) parsePotm(document *goquery.Document) (string, error) {
	potmTag := document.Find(".potm-container")

	if potmTag == nil {
		return "", errors.New("potm container was not found on page")
	}

	potmNicknameTag := potmTag.Find(".player-nick")

	nickname := strings.ToLower(potmNicknameTag.Text())

	return nickname, nil
}
