package parsers

import (
	"bytes"
	"fmt"
	"hltvapi/internal/models"
	"io/ioutil"
	"strings"
	"testing"
)

func TestGetTeam(t *testing.T) {
	tests := []struct {
		Id     int
		Result models.Team
	}{
		{4869, models.Team{
			Id:           4869,
			Ranking:      2,
			WeeksInTop30: 64,
			AverageAge:   24.9,
			Name:         "ence",
			Country:      "europe",
			Social: []models.Social{
				{Name: "twitter", Link: "https://www.twitter.com/ENCE"},
				{Name: "instagram", Link: "https://www.instagram.com/enceesports"},
			},
			Achievements: []models.Achievement{
				{Name: "pgl major antwerp 2022", Placement: "3-4th"},
				{Name: "pgl major stockholm 2021", Placement: "Group stage"},
				{Name: "starLadder major berlin 2019", Placement: "1/4 final"},
				{Name: "iem katowice 2019", Placement: "2nd"},
			},
		}},
	}

	p := TeamParser{}

	for _, test := range tests {
		team, err := p.GetTeam(test.Id)

		if err != nil {
			t.Errorf("Parse %d ended with error. Error: %s\n", test.Id, err.Error())
		}

		if ok, field := areTeamsEqual(test.Result, *team); !ok {
			t.Errorf("Parse Team %d ended with error. %s\n", test.Id, field)
		}
	}
}

func BenchmarkParseTeam(b *testing.B) {
	const url = "https://www.hltv.org/team/5752/cloud9"
	response, _ := SendRequest(url)
	body, _ := ioutil.ReadAll(response.Body)

	p := TeamParser{}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		b.StartTimer()

		_, err := p.parseTeamResponseAndId(url, response)

		if err != nil {
			b.Fatalf("Error during benchmarkParseTeam. Error: %v\n", err.Error())
		}
	}
}

func areTeamsEqual(x models.Team, y models.Team) (bool, string) {
	if x.Id != y.Id {
		return false, fmt.Sprintf("Field: id. Values: %v and %v", x.Id, y.Id)
	}

	if !strings.EqualFold(x.Country, y.Country) {
		return false, fmt.Sprintf("Field: country. Values: %v and %v", x.Country, y.Country)
	}

	if !strings.EqualFold(x.Name, y.Name) {
		return false, fmt.Sprintf("Field: name. Values: %v and %v", x.Name, y.Name)
	}

	if x.AverageAge != y.AverageAge {
		return false, fmt.Sprintf("Field: average age. Values: %v and %v", x.AverageAge, y.AverageAge)
	}

	if x.Ranking != y.Ranking {
		return false, fmt.Sprintf("Field: ranking. Values: %v and %v", x.Ranking, y.Ranking)
	}

	if x.WeeksInTop30 != y.WeeksInTop30 {
		return false, fmt.Sprintf("Field: Weeks in Top 30. Values: %v and %v", x.WeeksInTop30, y.WeeksInTop30)
	}

	minAchievements := len(x.Achievements)
	if minAchievements > len(y.Achievements) {
		minAchievements = len(y.Achievements)
	}

	for i := 0; i < minAchievements; i++ {
		if !strings.EqualFold(x.Achievements[i].Name, y.Achievements[i].Name) ||
			!strings.EqualFold(x.Achievements[i].Placement, y.Achievements[i].Placement) {
			return false, fmt.Sprintf("Field: achievements. Values: %v and %v", x.Achievements, y.Achievements)
		}
	}

	minSocial := len(x.Social)
	if minSocial > len(y.Social) {
		minSocial = len(y.Social)
	}

	for i := 0; i < minSocial; i++ {
		if !strings.EqualFold(x.Social[i].Name, y.Social[i].Name) ||
			!strings.EqualFold(x.Social[i].Link, y.Social[i].Link) {
			return false, fmt.Sprintf("Field: social. Values: %v and %v", x.Social, y.Social)
		}
	}

	return true, ""
}
