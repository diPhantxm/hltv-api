package parsers

import (
	"bytes"
	"fmt"
	"hltvapi/internal/models"
	"hltvapi/internal/urlBuilder/fileUrlBuilder"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestGetUpcomingMatches(t *testing.T) {
	tests := []models.Match{
		{
			Id:               2356965,
			TeamA:            "NAVI Javelins",
			TeamB:            "BIG EQUIPA",
			StartTime:        time.UnixMilli(1656670200000),
			Maps:             []string{"TBA"},
			Viewers:          0,
			PlayerOfTheMatch: "",
		},
		{
			Id:               2357054,
			TeamA:            "Eternal Fire",
			TeamB:            "ex-MAD Lions",
			StartTime:        time.UnixMilli(1656676800000),
			Maps:             []string{"TBA", "TBA", "TBA"},
			Viewers:          6649,
			PlayerOfTheMatch: "",
		},
		{
			Id:               2357055,
			TeamA:            "Complexity",
			TeamB:            "SINNERS",
			StartTime:        time.UnixMilli(1656687600000),
			Maps:             []string{"TBA", "TBA", "TBA"},
			Viewers:          6649,
			PlayerOfTheMatch: "",
		},
		{
			Id:               2357173,
			TeamA:            "UNGENTIUM",
			TeamB:            "FLET",
			StartTime:        time.UnixMilli(1656684000000),
			Maps:             []string{"TBA", "TBA", "TBA"},
			Viewers:          88,
			PlayerOfTheMatch: "",
		},
	}

	p := NewMatchParser(fileUrlBuilder.NewFileUrlBuilder())
	matches, err := p.GetMatches()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}

	matchesMap := make(map[int]models.Match)
	for _, match := range matches {
		matchesMap[match.Id] = match
	}

	if len(tests) > len(matches) {
		t.Errorf("Didn't get all matches. Length of tests is bigger than length of parsed matches.")
	}

	for _, test := range tests {
		if match, ok := matchesMap[test.Id]; !ok {
			t.Errorf("Missing match %d\n", test.Id)
		} else {
			if ok, field := areMatchesEqual(match, test); !ok {
				t.Errorf("Matches with id %d are not equal. %s", match.Id, field)
			}
		}
	}
}

func TestGetByDateMatches(t *testing.T) {
	tests := []struct {
		Date    string
		Matches []models.Match
	}{
		{
			Date: "2021-12-26",
			Matches: []models.Match{
				{
					Id:               2353847,
					TeamA:            "DNG",
					TeamB:            "INDE IRAE",
					StartTime:        time.UnixMilli(1640534700000),
					Maps:             []string{"Vertigo", "Inferno", "Overpass", "Mirage", "Nuke"},
					Viewers:          0,
					PlayerOfTheMatch: "Gospadarov",
				},
				{
					Id:               2353856,
					TeamA:            "B8",
					TeamB:            "IKIGAI",
					StartTime:        time.UnixMilli(1640538000000),
					Maps:             []string{"Nuke", "Ancient", "Mirage"},
					Viewers:          0,
					PlayerOfTheMatch: "cptkurtka023",
				},
				{
					Id:               2353853,
					TeamA:            "Team Shoke",
					TeamB:            "Team StRoGo",
					StartTime:        time.UnixMilli(1640539800000),
					Maps:             []string{"Ancient", "Dust2", "Mirage"},
					Viewers:          0,
					PlayerOfTheMatch: "Patsi",
				},
				{
					Id:               2353852,
					TeamA:            "Team StRoGo",
					TeamB:            "Team Evelone192",
					StartTime:        time.UnixMilli(1640523600000),
					Maps:             []string{"Dust2", "Ancient", "Mirage"},
					Viewers:          0,
					PlayerOfTheMatch: "deko",
				},
			},
		},
	}

	p := NewMatchParser(fileUrlBuilder.NewFileUrlBuilder())

	for _, test := range tests {
		matches, err := p.GetMatchesByDate(test.Date)
		if err != nil {
			t.Fatalf("Error: %s", err.Error())
		}

		matchesMap := make(map[int]models.Match)
		for _, match := range matches {
			matchesMap[match.Id] = match
		}

		for _, testMatch := range test.Matches {
			if match, ok := matchesMap[testMatch.Id]; !ok {
				t.Errorf("Missing match %d\n", testMatch.Id)
			} else {
				if ok, field := areMatchesEqual(match, testMatch); !ok {
					t.Errorf("Matches with id %d are not equal. %s", match.Id, field)
				}
			}
		}
	}
}

func BenchmarkParseMatch(b *testing.B) {
	const url = "https://www.hltv.org/matches/2356673/nip-vs-nasr-global-esports-tour-dubai-2022"
	response, _ := SendRequest(url)
	body, _ := ioutil.ReadAll(response.Body)

	p := MatchParser{}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		b.StartTimer()

		_, err := p.parseMatchResponse(response)

		if err != nil {
			b.Fatalf("Error during benchmarkParseMatch. Error: %v\n", err.Error())
		}
	}
}

func areMatchesEqual(x models.Match, y models.Match) (bool, string) {
	if x.Id != y.Id {
		return false, fmt.Sprintf("Field: id. Values: %d and %d", x.Id, y.Id)
	}

	if !strings.EqualFold(x.TeamA, y.TeamA) ||
		!strings.EqualFold(x.TeamB, y.TeamB) {
		return false, fmt.Sprintf("Field: team names. Values: %s, %s and %s, %s", x.TeamA, x.TeamB, y.TeamA, y.TeamB)
	}

	if !x.StartTime.Equal(y.StartTime) {
		return false, fmt.Sprintf("Field: start time. Values: %v and %v", x.StartTime, y.StartTime)
	}

	if len(x.Maps) != len(y.Maps) {
		return false, fmt.Sprintf("Field: maps. Values: %v and %v", x.Maps, y.Maps)
	}

	for i := 0; i < len(x.Maps); i++ {
		if !strings.EqualFold(x.Maps[i], y.Maps[i]) {
			return false, fmt.Sprintf("Field: maps. Values: %v and %v", x.Maps, y.Maps)
		}
	}

	if x.Viewers != y.Viewers {
		return false, fmt.Sprintf("Field: viewers. Values: %d and %d", x.Viewers, y.Viewers)
	}

	if !strings.EqualFold(x.PlayerOfTheMatch, y.PlayerOfTheMatch) {
		return false, fmt.Sprintf("Field: player of the match. Values: %s and %s", x.PlayerOfTheMatch, y.PlayerOfTheMatch)
	}

	if x.IsOver != y.IsOver {
		return false, fmt.Sprintf("Field: isOver. Values: %v and %v", x.IsOver, y.IsOver)
	}

	return true, ""
}
