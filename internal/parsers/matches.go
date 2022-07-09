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
	builder urlBuilder.MatchUrlBuilder
}

func NewMatchParser(builder urlBuilder.MatchUrlBuilder) *MatchParser {
	return &MatchParser{
		builder: builder,
	}
}

func (p MatchParser) GetMatch(id int) (*models.Match, error) {
	p.builder.Clear()
	p.builder.Match()
	p.builder.AddId(id)

	response, err := SendRequest(p.builder.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	event, err := p.parseMatchResponse(response)
	if err != nil {
		return nil, err
	}

	event.Id = id

	return event, nil
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
	p.builder.Clear()
	p.builder.Results()
	p.builder.AddParam("startDate", date)
	p.builder.AddParam("endDate", date)

	ids, err := p.getMatchesIdsByDate(p.builder.String())
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
	p.builder.Clear()
	p.builder.Match()

	response, err := SendRequest(p.builder.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	matchTags := document.Find(".upcomingMatchesSection")
	var ids []int

	matchTags.Each(func(i int, selection *goquery.Selection) {
		links := selection.Find(".match")
		links.Each(func(i int, link *goquery.Selection) {
			href, ok := link.Attr("href")
			if !ok {
				return
			}

			re := regexp.MustCompile(`\/\w+\/(\d+)`)
			idStr := re.FindStringSubmatch(href)[1]

			id, err := strconv.Atoi(idStr)
			if err != nil {
				return
			}

			ids = append(ids, id)
		})
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

	idTags := document.Find(".results-holder").Find(".results-all").Find("a")
	ids := make([]int, idTags.Length())

	idTags.Each(func(i int, selection *goquery.Selection) {
		link, ok := selection.Attr("href")
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

func (p MatchParser) parseMatchResponse(response *http.Response) (*models.Match, error) {
	match, err := p.parse(response)
	if err != nil {
		return nil, err
	}

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

	isOver, err := p.parseIsOver(document)
	if err != nil {
		return nil, err
	}
	match.IsOver = isOver

	return match, nil
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

func (p MatchParser) parseMaps(document *goquery.Document) ([]models.Map, error) {
	mapTags := document.Find(".mapholder")
	maps := make([]models.Map, mapTags.Length())

	if mapTags == nil {
		return nil, errors.New("no maps were found on page")
	}

	mapTags.Each(func(i int, selection *goquery.Selection) {
		mapName := selection.Find(".mapname")
		maps[i].Name = mapName.Text()
	})

	scores := document.Find(".results-team-score")
	scores.Each(func(i int, score *goquery.Selection) {
		if i%2 == 0 {
			maps[i/2].TeamAScore, _ = strconv.Atoi(score.Text())
		} else {
			maps[i/2].TeamBScore, _ = strconv.Atoi(score.Text())
		}
	})

	return maps, nil
}

func (p MatchParser) parseViewers(document *goquery.Document) (int, error) {
	viewersTags := document.Find(".left-right-padding")
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

func (p MatchParser) parseIsOver(document *goquery.Document) (bool, error) {
	progressTag := document.Find(".countdown")

	if strings.ToLower(progressTag.Text()) == "match over" {
		return true, nil
	}

	return false, nil
}
