package parsers

import (
	"bytes"
	"fmt"
	"hltvapi/internal/models"
	"io/ioutil"
	"strings"
	"testing"
)

func TestGetPlayer(t *testing.T) {
	tests := []struct {
		Url    string
		Result models.Player
	}{
		{"https://www.hltv.org/player/922/snappi", models.Player{
			Id:        922,
			Age:       32,
			Team:      "ENCE",
			Nickname:  "Snappi",
			FirstName: "Marco",
			LastName:  "Pfeiffer",
			Country:   "Denmark",
			Stats: models.Statistics{
				Rating:            0.98,
				KillsPerRound:     0.61,
				Headshots:         46.3,
				MapsPlayed:        84,
				DeathsPerRound:    0.66,
				RoundsContributed: 68.9,
			},
			Social: []models.Social{
				{Name: "twitter", Link: "https://www.twitter.com/SnappiCSGO"},
			},
		}},
		{"https://www.hltv.org/player/9618/nexa", models.Player{
			Id:        9618,
			Age:       25,
			Team:      "OG",
			Nickname:  "nexa",
			FirstName: "Nemanja",
			LastName:  "Isaković",
			Country:   "Serbia",
			Stats: models.Statistics{
				Rating:            1.04,
				KillsPerRound:     0.66,
				Headshots:         48.6,
				MapsPlayed:        12,
				DeathsPerRound:    0.65,
				RoundsContributed: 70.4,
			},
			Social: []models.Social{
				{Name: "twitter", Link: "https://www.twitter.com/nexaOG"},
			},
		}},
	}

	for _, test := range tests {
		player, err := GetPlayer(test.Url)

		if err != nil {
			t.Errorf("Parse Players %s ended with error. Error: %s\n", test.Url, err.Error())
		}

		if ok, field := arePlayersEqual(test.Result, *player); !ok {
			t.Errorf("Parse Players %s ended with error. Field: %s\n", test.Url, field)
		}
	}
}

func BenchmarkParsePlayer(b *testing.B) {
	const url = "https://www.hltv.org/player/8528/hobbit"
	response, _ := sendRequest(url)
	body, _ := ioutil.ReadAll(response.Body)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		b.StartTimer()

		_, err := parsePlayerResponseAndId(url, response)

		if err != nil {
			b.Fatalf("Error during benchmarkParsePlayer. Error: %v\n", err.Error())
		}
	}
}

func arePlayersEqual(x models.Player, y models.Player) (bool, string) {
	if x.Id != y.Id {
		return false, fmt.Sprintf("Field: id. Values: %v and %v", x.Id, y.Id)
	}

	if x.Age != y.Age {
		return false, fmt.Sprintf("Field: age. Values: %v and %v", x.Age, y.Age)
	}

	if !strings.EqualFold(x.Country, y.Country) {
		return false, fmt.Sprintf("Field: country. Values: %v and %v", x.Country, y.Country)
	}

	if !strings.EqualFold(x.FirstName, y.FirstName) {
		return false, fmt.Sprintf("Field: first name. Values: %v and %v", x.FirstName, y.FirstName)
	}

	if !strings.EqualFold(x.LastName, y.LastName) {
		return false, fmt.Sprintf("Field: last name. Values: %v and %v", x.LastName, y.LastName)
	}

	if !strings.EqualFold(x.Nickname, y.Nickname) {
		return false, fmt.Sprintf("Field: nickname. Values: %v and %v", x.Nickname, y.Nickname)
	}

	if len(x.Social) != len(y.Social) {
		return false, fmt.Sprintf("Field: social. Values: %v and %v", x.Social, y.Social)
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

	if ok, field := areStatsEqual(x.Stats, y.Stats); !ok {
		return false, fmt.Sprintf("Field: %v. Values: %v and %v", field, x.Stats, y.Stats)
	}

	return true, ""
}

func areStatsEqual(x models.Statistics, y models.Statistics) (bool, string) {
	if x.DeathsPerRound != y.DeathsPerRound {
		return false, "statistics: deaths per round"
	}

	if x.Headshots != y.Headshots {
		return false, "statistics: headshot"
	}

	if x.KillsPerRound != y.KillsPerRound {
		return false, "statistics: kills per round"
	}

	if x.MapsPlayed != y.MapsPlayed {
		return false, "statistics: maps player"
	}

	if x.Rating != y.Rating {
		return false, "statistics: rating"
	}

	if x.RoundsContributed != y.RoundsContributed {
		return false, "statistics: rounds contributed"
	}

	return true, ""
}
