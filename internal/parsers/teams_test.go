package parsers

import (
	"bytes"
	"fmt"
	"hltvapi/internal/models"
	"hltvapi/internal/urlBuilder/fileUrlBuilder"
	"io/ioutil"
	"strings"
	"testing"
)

func TestGetAllTeams(t *testing.T) {
	// Note: Tests might contain reduced amount of entities.
	tests := []models.Team{
		{
			Id:           4608,
			Ranking:      2,
			WeeksInTop30: 127,
			Name:         "Natus Vincere",
			Country:      "Ukraine",
			Social: []models.Social{
				{
					Name: "Twitter",
					Link: "https://www.twitter.com/natusvincere",
				},
				{
					Name: "Instagram",
					Link: "https://www.instagram.com/natus_vincere_official",
				},
			},
			Achievements: []models.Achievement{
				{
					Name:      "BLAST Premier Spring Final 2022",
					Placement: "1st",
				},
				{
					Name:      "PGL Major Antwerp 2022",
					Placement: "2nd",
				},
				{
					Name:      "IEM Katowice 2022",
					Placement: "3-4th",
				},
			},
			AverageAge: 23.2,
		},
		{
			Id:           6673,
			Ranking:      0,
			WeeksInTop30: 0,
			Name:         "NRG",
			Country:      "United States",
			Social: []models.Social{
				{
					Name: "Twitter",
					Link: "https://www.twitter.com/NRGgg",
				},
			},
			Achievements: []models.Achievement{
				{
					Name:      "StarLadder Major Berlin 2019",
					Placement: "3-4th",
				},
				{
					Name:      "IEM Katowice 2019",
					Placement: "Group stage",
				},
			},
			AverageAge: 0,
		},
		{
			Id:           10386,
			Ranking:      30,
			WeeksInTop30: 55,
			Name:         "SKADE",
			Country:      "Bulgaria",
			Social: []models.Social{
				{
					Name: "Twitter",
					Link: "https://www.twitter.com/skadegg",
				},
				{
					Name: "Twitch",
					Link: "https://www.twitch.tv/skadegg",
				},
				{
					Name: "Instagram",
					Link: "https://www.instagram.com/skadegg",
				},
			},
			Achievements: []models.Achievement{},
			AverageAge:   24.4,
		},
		{
			Id:           8248,
			Ranking:      111,
			WeeksInTop30: 0,
			Name:         "PACT",
			Country:      "Poland",
			Social: []models.Social{
				{
					Name: "Twitter",
					Link: "https://www.twitter.com/PACT_gg",
				},
			},
			Achievements: []models.Achievement{},
			AverageAge:   25.5,
		},
		{
			Id:           9215,
			Ranking:      20,
			WeeksInTop30: 15,
			Name:         "MIBR",
			Country:      "Brazil",
			Social: []models.Social{
				{
					Name: "Twitter",
					Link: "https://www.twitter.com/MIBR",
				},
			},
			Achievements: []models.Achievement{
				{
					Name:      "StarLadder Major Berlin 2019",
					Placement: "Group stage",
				},
				{
					Name:      "IEM Katowice 2019",
					Placement: "3-4th",
				},
				{
					Name:      "FACEIT Major 2018",
					Placement: "3-4th",
				},
			},
			AverageAge: 23.4,
		},
		{
			Id:           4411,
			Ranking:      8,
			WeeksInTop30: 110,
			Name:         "NIP",
			Country:      "Sweden",
			Social: []models.Social{
				{
					Name: "Twitter",
					Link: "https://www.twitter.com/nipcs",
				},
				{
					Name: "Instagram",
					Link: "https://www.instagram.com/nipgaming",
				},
			},
			Achievements: []models.Achievement{
				{
					Name:      "PGL Major Antwerp 2022",
					Placement: "1/4 final",
				},
				{
					Name:      "PGL Major Stockholm 2021",
					Placement: "1/4 final",
				},
				{
					Name:      "StarLadder Major Berlin 2019",
					Placement: "Group stage",
				},
			},
			AverageAge: 23.0,
		},
		{
			Id:           5995,
			Ranking:      6,
			WeeksInTop30: 72,
			Name:         "G2",
			Country:      "Europe",
			Social: []models.Social{
				{
					Name: "Twitter",
					Link: "https://www.twitter.com/G2esports",
				},
				{
					Name: "Instagram",
					Link: "https://www.instagram.com/g2esports",
				},
			},
			Achievements: []models.Achievement{
				{
					Name:      "PGL Major Antwerp 2022",
					Placement: "Group stage",
				},
				{
					Name:      "PGL Major Stockholm 2021",
					Placement: "2nd",
				},
				{
					Name:      "StarLadder Major Berlin 2019",
					Placement: "Group stage",
				},
			},
			AverageAge: 24.8,
		},
	}

	p := NewTeamParser(fileUrlBuilder.NewFileUrlBuilder())
	teams, err := p.GetTeams()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}

	teamsMap := make(map[int]models.Team)
	for _, team := range teams {
		teamsMap[team.Id] = team
	}

	if len(tests) > len(teams) {
		t.Errorf("Didn't get all teams. Length of tests is bigger than length of parsed teams.")
	}

	for _, test := range tests {
		if team, ok := teamsMap[test.Id]; !ok {
			t.Errorf("Missing team %d\n", test.Id)
		} else {
			if ok, field := areTeamsEqual(team, test); !ok {
				t.Errorf("Teams with id %d are not equal. Field: %s", team.Id, field)
			}
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

		_, err := p.parseTeamResponse(response)

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
			return false, fmt.Sprintf("Field: achievements. Values: %v and %v", x.Achievements[i], y.Achievements[i])
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
