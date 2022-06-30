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

func TestGetAllPlayers(t *testing.T) {
	tests := []models.Player{
		{
			Id:        922,
			Age:       32,
			Team:      "ENCE",
			Nickname:  "Snappi",
			FirstName: "Marco",
			LastName:  "Pfeiffer",
			Country:   "Denmark",
			Stats: models.Statistics{
				Rating:            0.95,
				KillsPerRound:     0.60,
				Headshots:         45.0,
				MapsPlayed:        73,
				DeathsPerRound:    0.68,
				RoundsContributed: 68.1,
			},
			Social: []models.Social{
				{Name: "twitter", Link: "https://www.twitter.com/SnappiCSGO"},
			},
		},
		{
			Id:        9618,
			Age:       25,
			Team:      "OG",
			Nickname:  "nexa",
			FirstName: "Nemanja",
			LastName:  "Isaković",
			Country:   "Serbia",
			Stats: models.Statistics{
				Rating:            1.01,
				KillsPerRound:     0.66,
				Headshots:         49.1,
				MapsPlayed:        29,
				DeathsPerRound:    0.68,
				RoundsContributed: 69.0,
			},
			Social: []models.Social{
				{Name: "twitter", Link: "https://www.twitter.com/nexaOG"},
			},
		},
		{
			Id:        11893,
			Age:       21,
			Team:      "Vitality",
			Nickname:  "ZyWOo",
			FirstName: "Mathieu",
			LastName:  "Herbaut",
			Country:   "France",
			Stats: models.Statistics{
				Rating:            1.26,
				KillsPerRound:     0.80,
				Headshots:         38.7,
				MapsPlayed:        43,
				DeathsPerRound:    0.59,
				RoundsContributed: 75.5,
			},
			Social: []models.Social{
				{Name: "twitter", Link: "https://www.twitter.com/zywoo"},
				{Name: "twitch", Link: "https://www.twitch.tv/cs_zywoo"},
				{Name: "instagram", Link: "https://www.instagram.com/cs_zywoo"},
			},
		},
		{
			Id:        7998,
			Age:       24,
			Team:      "Natus Vincere",
			Nickname:  "s1mple",
			FirstName: "Aleksandr",
			LastName:  "Kostyliev",
			Country:   "Ukraine",
			Stats: models.Statistics{
				Rating:            1.35,
				KillsPerRound:     0.88,
				Headshots:         38.7,
				MapsPlayed:        46,
				DeathsPerRound:    0.59,
				RoundsContributed: 76.1,
			},
			Social: []models.Social{
				{Name: "twitter", Link: "https://www.twitter.com/s1mpleO"},
				{Name: "twitch", Link: "https://www.twitch.tv/s1mple"},
				{Name: "instagram", Link: "https://www.instagram.com/s1mpleo"},
			},
		},
		{
			Id:        17861,
			Age:       24,
			Team:      "MIBR",
			Nickname:  "JOTA",
			FirstName: "Jhonatan",
			LastName:  "Willian",
			Country:   "Brazil",
			Stats: models.Statistics{
				Rating:            1.15,
				KillsPerRound:     0.74,
				Headshots:         42.5,
				MapsPlayed:        63,
				DeathsPerRound:    0.64,
				RoundsContributed: 72.2,
			},
			Social: []models.Social{
				{Name: "twitter", Link: "https://www.twitter.com/jotafps"},
				{Name: "twitch", Link: "https://www.twitch.tv/jotaacs"},
				{Name: "instagram", Link: "https://www.instagram.com/jotaacs"},
			},
		},
		{
			Id:        13042,
			Age:       24,
			Team:      "No team",
			Nickname:  "LOGAN",
			FirstName: "Logan",
			LastName:  "Corti",
			Country:   "France",
			Stats: models.Statistics{
				Rating:            0.00,
				KillsPerRound:     0.00,
				Headshots:         0.00,
				MapsPlayed:        0,
				DeathsPerRound:    0.00,
				RoundsContributed: 0.00,
			},
			Social: []models.Social{},
		},
	}

	p := NewPlayerParser(fileUrlBuilder.NewFileUrlBuilder())
	players, err := p.GetPlayers()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}

	playersMap := make(map[int]models.Player)
	for _, player := range players {
		playersMap[player.Id] = player
	}

	if len(tests) > len(players) {
		t.Errorf("Didn't get all players. Length of tests is bigger than length of parsed players.")
	}

	for _, test := range tests {
		if player, ok := playersMap[test.Id]; !ok {
			t.Errorf("Missing player %d\n", test.Id)
		} else {
			if ok, field := arePlayersEqual(player, test); !ok {
				t.Errorf("Players with id %d are not equal. Field: %s", player.Id, field)
			}
		}
	}
}

func BenchmarkParsePlayer(b *testing.B) {
	const url = "https://www.hltv.org/player/8528/hobbit"
	response, _ := SendRequest(url)
	body, _ := ioutil.ReadAll(response.Body)

	p := PlayerParser{}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		b.StartTimer()

		_, err := p.parsePlayerResponse(response)

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
