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

type TeamParser struct {
	builder urlBuilder.TeamUrlBuilder
}

func NewTeamParser(builder urlBuilder.TeamUrlBuilder) *TeamParser {
	return &TeamParser{
		builder: builder,
	}
}

func (p TeamParser) GetTeam(id int) (*models.Team, error) {
	p.builder.Clear()
	p.builder.Team()
	p.builder.AddId(id)

	response, err := SendRequest(p.builder.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	team, err := p.parseTeamResponse(response)
	if err != nil {
		return nil, err
	}

	team.Id = id

	socialParser := SocialParser{}
	socials, err := socialParser.parse(response)
	if err != nil {
		return nil, err
	}
	team.Social = socials

	return team, nil
}

func (p TeamParser) GetTeams() ([]models.Team, error) {
	ids, err := p.GetAllTeamsIds()
	if err != nil {
		return nil, err
	}

	result := make([]models.Team, len(ids))
	for i, id := range ids {
		team, err := p.GetTeam(id)
		if err != nil {
			continue
		}

		result[i] = *team
	}

	return result, nil
}

func (p TeamParser) GetAllTeamsIds() ([]int, error) {
	p.builder.Clear()
	p.builder.TeamsStats()

	response, err := SendRequest(p.builder.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	teamTags := document.Find(".stats-table").Find("tbody").Find("tr")
	ids := make([]int, teamTags.Length())

	teamTags.Each(func(i int, selection *goquery.Selection) {
		link, ok := selection.Find(".teamCol-teams-overview").Find("a").Attr("href")
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

func (p TeamParser) parseTeamResponse(response *http.Response) (*models.Team, error) {
	team, err := p.parse(response)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (p TeamParser) parse(response *http.Response) (*models.Team, error) {
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	team := &models.Team{}

	name, err := p.parseName(document)
	if err != nil {
		return nil, err
	}
	team.Name = name

	country, err := p.parseCountry(document)
	if err != nil {
		return nil, err
	}
	team.Country = country

	ranking, weeksInTop30, averageAge := p.parseProfileTeamStats(document)
	team.Ranking = ranking
	team.WeeksInTop30 = weeksInTop30
	team.AverageAge = averageAge

	team.Achievements = p.parseAchievements(document)

	team.Roaster = p.parseRoaster(document)

	return team, nil
}

func (p TeamParser) parseName(document *goquery.Document) (string, error) {
	nameTag := document.Find(".profile-team-name")
	name := nameTag.Text()

	if name == "" {
		return "", errors.New("team name was not found on page")
	}

	return name, nil
}

func (p TeamParser) parseCountry(document *goquery.Document) (string, error) {
	countryTag := document.Find(".team-country").Find("img.flag")

	country, ok := countryTag.Attr("alt")
	if !ok {
		return "", errors.New("team flag was not found on page")
	}

	return country, nil
}

func (p TeamParser) parseProfileTeamStats(document *goquery.Document) (int, int, float32) {
	ranking := 0
	weeksInTop30 := 0
	var averageAge float32 = 0

	statsTag := document.Find(".profile-team-stat")
	statsTag.Each(func(i int, selection *goquery.Selection) {
		statName := selection.Find("b").Text()
		statValue := selection.Find("span").Text()

		if statName == "Weeks in top30 for core" {
			valueInt, err := strconv.ParseInt(statValue, 10, 32)
			if err != nil {
				return
			}
			weeksInTop30 = int(valueInt)
		}
		if statName == "World ranking" {
			statValue = strings.ReplaceAll(statValue, "#", "")
			valueInt, err := strconv.ParseInt(statValue, 10, 32)
			if err != nil {
				return
			}
			ranking = int(valueInt)
		}
		if statName == "Average player age" {
			valueInt, err := strconv.ParseFloat(statValue, 32)
			if err != nil {
				return
			}
			averageAge = float32(valueInt)
		}
	})

	return ranking, weeksInTop30, averageAge
}

func (p TeamParser) parseAchievements(document *goquery.Document) []models.Achievement {
	table := document.Find(".achievement-table")
	rows := table.Find("t-body").Find("tr")

	achievements := make([]models.Achievement, rows.Length())
	rows.Each(func(i int, selection *goquery.Selection) {
		achievements[i].Placement = strings.Trim(selection.Find(".placement-cell").Text(), "\n ")
		achievements[i].Name = selection.Find(".tournament-name-cell").Text()
	})

	return achievements
}

func (p TeamParser) parseRoaster(document *goquery.Document) models.Roaster {
	rows := document.Find(".players-table").Find("tbody").Find("tr")

	roaster := make(models.Roaster, rows.Length())
	rows.Each(func(i int, selection *goquery.Selection) {
		data := selection.Find(".players-cell")
		player := models.TeamPlayer{}

		player.Nickname = strings.Trim(data.Eq(0).Text(), "\n ")
		player.Status = strings.Trim(data.Eq(1).Text(), "\n ")

		timeOntTeam, _ := data.Eq(2).Html()
		player.TimeOnTeam = strings.ReplaceAll(timeOntTeam, "<br/>", " ")

		player.MapsPlayed, _ = strconv.Atoi(data.Eq(3).Text())

		rating, _ := strconv.ParseFloat(data.Eq(4).Text(), 32)
		player.Rating = float32(rating)

		roaster[i] = player
	})

	return roaster
}
