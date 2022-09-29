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

// Note: Tests might contain reduced amount of entities.
func TestGetAllTeams(t *testing.T) {
	tests := []models.Team{
		{
			Id:           4608,
			Ranking:      2,
			WeeksInTop30: 127,
			Name:         "Natus Vincere",
			Country:      "Ukraine",
			Roaster: models.Roaster{
				{
					Nickname:   "s1mple",
					Status:     "STARTER",
					TimeOnTeam: "5 years 10 months",
					MapsPlayed: 1085,
					Rating:     1.30,
				},
				{
					Nickname:   "electroNic",
					Status:     "STARTER",
					TimeOnTeam: "4 years 7 months",
					MapsPlayed: 884,
					Rating:     1.16,
				},
				{
					Nickname:   "Perfecto",
					Status:     "STARTER",
					TimeOnTeam: "2 years 4 months",
					MapsPlayed: 477,
					Rating:     1.01,
				},
				{
					Nickname:   "b1t",
					Status:     "STARTER",
					TimeOnTeam: "1 year 5 months",
					MapsPlayed: 229,
					Rating:     1.12,
				},
				{
					Nickname:   "sdy",
					Status:     "STARTER",
					TimeOnTeam: "24 days",
					MapsPlayed: 12,
					Rating:     1.02,
				},
				{
					Nickname:   "Boombl4",
					Status:     "BENCHED",
					TimeOnTeam: "2 years 12 months",
					MapsPlayed: 531,
					Rating:     0.99,
				},
			},
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
			Roaster:      models.Roaster{},
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
			Roaster: models.Roaster{
				{
					Nickname:   "dennyslaw",
					Status:     "STARTER",
					TimeOnTeam: "2 years 9 months",
					MapsPlayed: 1081,
					Rating:     1.12,
				},
				{
					Nickname:   "Rainwaker",
					Status:     "STARTER",
					TimeOnTeam: "2 years 9 months",
					MapsPlayed: 1077,
					Rating:     1.11,
				},
				{
					Nickname:   "SHiPZ",
					Status:     "STARTER",
					TimeOnTeam: "1 year 5 months",
					MapsPlayed: 606,
					Rating:     1.11,
				},
				{
					Nickname:   "bubble",
					Status:     "STARTER",
					TimeOnTeam: "7 months",
					MapsPlayed: 188,
					Rating:     0.93,
				},
				{
					Nickname:   "dream3r",
					Status:     "STARTER",
					TimeOnTeam: "5 months",
					MapsPlayed: 178,
					Rating:     1.1,
				},
			},
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
			Roaster: models.Roaster{
				{
					Nickname:   "lunAtic",
					Status:     "STARTER",
					TimeOnTeam: "4 years 4 months",
					MapsPlayed: 840,
					Rating:     0.95,
				},
				{
					Nickname:   "bnox",
					Status:     "STARTER",
					TimeOnTeam: "1 year 1 month",
					MapsPlayed: 189,
					Rating:     0.99,
				},
				{
					Nickname:   "reatz",
					Status:     "STARTER",
					TimeOnTeam: "10 months",
					MapsPlayed: 123,
					Rating:     1.00,
				},
				{
					Nickname:   "SAYN",
					Status:     "STARTER",
					TimeOnTeam: "3 months",
					MapsPlayed: 34,
					Rating:     1.00,
				},
				{
					Nickname:   "TOAO",
					Status:     "STARTER",
					TimeOnTeam: "3 months",
					MapsPlayed: 27,
					Rating:     0.93,
				},
			},
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
			Roaster: models.Roaster{
				{
					Nickname:   "chelo",
					Status:     "STARTER",
					TimeOnTeam: "1 year 5 months",
					MapsPlayed: 219,
					Rating:     1.16,
				},
				{
					Nickname:   "exit",
					Status:     "STARTER",
					TimeOnTeam: "1 year 3 months",
					MapsPlayed: 217,
					Rating:     1.07,
				},
				{
					Nickname:   "Tuurtle",
					Status:     "STARTER",
					TimeOnTeam: "8 months",
					MapsPlayed: 155,
					Rating:     1.07,
				},
				{
					Nickname:   "JOTA",
					Status:     "STARTER",
					TimeOnTeam: "8 months",
					MapsPlayed: 139,
					Rating:     1.19,
				},
				{
					Nickname:   "brnz4n",
					Status:     "STARTER",
					TimeOnTeam: "1 month",
					MapsPlayed: 19,
					Rating:     1.02,
				},
				{
					Nickname:   "WOOD7",
					Status:     "BENCHED",
					TimeOnTeam: "8 months",
					MapsPlayed: 128,
					Rating:     1.02,
				},
			},
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
			Roaster: models.Roaster{
				{
					Nickname:   "REZ",
					Status:     "STARTER",
					TimeOnTeam: "4 years 12 months",
					MapsPlayed: 961,
					Rating:     1.08,
				},
				{
					Nickname:   "Plopski",
					Status:     "STARTER",
					TimeOnTeam: "2 years 10 months",
					MapsPlayed: 527,
					Rating:     1.04,
				},
				{
					Nickname:   "hampus",
					Status:     "STARTER",
					TimeOnTeam: "2 years 1 month",
					MapsPlayed: 363,
					Rating:     1.03,
				},
				{
					Nickname:   "es3tag",
					Status:     "STARTER",
					TimeOnTeam: "7 months",
					MapsPlayed: 109,
					Rating:     0.95,
				},
				{
					Nickname:   "Brollan",
					Status:     "STARTER",
					TimeOnTeam: "3 months",
					MapsPlayed: 37,
					Rating:     1.15,
				},
				{
					Nickname:   "device",
					Status:     "BENCHED",
					TimeOnTeam: "1 year 2 months",
					MapsPlayed: 129,
					Rating:     1.13,
				},
			},
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
			Roaster: models.Roaster{
				{
					Nickname:   "JACKZ",
					Status:     "STARTER",
					TimeOnTeam: "3 years 6 months",
					MapsPlayed: 681,
					Rating:     1.00,
				},
				{
					Nickname:   "huNter-",
					Status:     "STARTER",
					TimeOnTeam: "2 years 9 months",
					MapsPlayed: 537,
					Rating:     1.13,
				},
				{
					Nickname:   "NiKo",
					Status:     "STARTER",
					TimeOnTeam: "1 year 7 months",
					MapsPlayed: 313,
					Rating:     1.21,
				},
				{
					Nickname:   "m0NESY",
					Status:     "STARTER",
					TimeOnTeam: "5 months",
					MapsPlayed: 78,
					Rating:     1.13,
				},
				{
					Nickname:   "Aleksib",
					Status:     "STARTER",
					TimeOnTeam: "5 months",
					MapsPlayed: 71,
					Rating:     0.93,
				},
				{
					Nickname:   "kennyS",
					Status:     "BENCHED",
					TimeOnTeam: "5 years 4 months",
					MapsPlayed: 879,
					Rating:     1.11,
				},
				{
					Nickname:   "AMANEK",
					Status:     "BENCHED",
					TimeOnTeam: "3 years 3 months",
					MapsPlayed: 589,
					Rating:     1.02,
				},
			},
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

	if len(x.Roaster) != len(y.Roaster) {
		return false, fmt.Sprintf("Field: roaster. Values %v and %v", x.Roaster, y.Roaster)
	}

	for i := 0; i < len(x.Roaster); i++ {
		if !strings.EqualFold(x.Roaster[i].Nickname, y.Roaster[i].Nickname) ||
			!strings.EqualFold(x.Roaster[i].TimeOnTeam, y.Roaster[i].TimeOnTeam) ||
			x.Roaster[i].MapsPlayed != y.Roaster[i].MapsPlayed ||
			x.Roaster[i].Rating != y.Roaster[i].Rating ||
			!strings.EqualFold(x.Roaster[i].Status, y.Roaster[i].Status) {
			return false, fmt.Sprintf("Field: roaster. Values: %v and %v", x.Roaster[i], y.Roaster[i])
		}
	}

	return true, ""
}
