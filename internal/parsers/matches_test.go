package parsers

import (
	"bytes"
	"hltvapi/internal/models"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestGetFinishedMatch(t *testing.T) {
	tests := []struct {
		Url    string
		Result models.Match
	}{
		{"https://www.hltv.org/matches/2356303/ence-vs-cloud9-iem-dallas-2022", models.Match{
			Id:               2356303,
			TeamA:            "ENCE",
			TeamB:            "Cloud9",
			StartTime:        time.UnixMilli(1654444800000),
			Maps:             []string{"Mirage", "Overpass", "Ancient", "Vertigo", "Dust2"},
			Viewers:          0,
			PlayerOfTheMatch: "sh1ro",
		}},
		{"https://www.hltv.org/matches/2356302/big-vs-cloud9-iem-dallas-2022", models.Match{
			Id:               2356302,
			TeamA:            "BIG",
			TeamB:            "Cloud9",
			StartTime:        time.UnixMilli(1654371000000),
			Maps:             []string{"Overpass", "Dust2", "Ancient"},
			Viewers:          0,
			PlayerOfTheMatch: "HObbit",
		}},
		{"https://www.hltv.org/matches/2356301/ence-vs-furia-iem-dallas-2022", models.Match{
			Id:               2356301,
			TeamA:            "ENCE",
			TeamB:            "FURIA",
			StartTime:        time.UnixMilli(1654358400000),
			Maps:             []string{"Vertigo", "Nuke", "Mirage"},
			Viewers:          0,
			PlayerOfTheMatch: "Snax",
		}},
		{"https://www.hltv.org/matches/2356300/faze-vs-cloud9-iem-dallas-2022", models.Match{
			Id:               2356300,
			TeamA:            "FaZe",
			TeamB:            "Cloud9",
			StartTime:        time.UnixMilli(1654296900000),
			Maps:             []string{"Overpass", "Inferno", "Mirage"},
			Viewers:          0,
			PlayerOfTheMatch: "Ax1le",
		}},
		{"https://www.hltv.org/matches/2356298/liquid-vs-cloud9-iem-dallas-2022", models.Match{
			Id:               2356298,
			TeamA:            "Liquid",
			TeamB:            "Cloud9",
			StartTime:        time.UnixMilli(1654124400000),
			Maps:             []string{"Ancient", "Vertigo", "Overpass"},
			Viewers:          0,
			PlayerOfTheMatch: "Ax1le",
		}},
		{"https://www.hltv.org/matches/2356296/ence-vs-faze-iem-dallas-2022", models.Match{
			Id:               2356296,
			TeamA:            "ENCE",
			TeamB:            "FaZe",
			StartTime:        time.UnixMilli(1654111800000),
			Maps:             []string{"Mirage", "Ancient", "Nuke"},
			Viewers:          0,
			PlayerOfTheMatch: "Snax",
		}},
	}

	for _, test := range tests {
		actual, err := GetMatch(test.Url)

		if err != nil {
			t.Errorf("Parse Match %s ended with error. Error: %s\n", test.Url, err.Error())
		}

		if ok, field := areMatchesEqual(*actual, test.Result); !ok {
			t.Errorf("Parse Match %s was incorrect. Field: %s\n", test.Url, field)
		}
	}
}

func BenchmarkParseMatch(b *testing.B) {
	const url = "https://www.hltv.org/matches/2356673/nip-vs-nasr-global-esports-tour-dubai-2022"
	response, _ := sendRequest(url)
	body, _ := ioutil.ReadAll(response.Body)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		b.StartTimer()

		_, err := parseMatchResponseAndId(url, response)

		if err != nil {
			b.Fatalf("Error during benchmarkParseMatch. Error: %v\n", err.Error())
		}
	}
}

func areMatchesEqual(x models.Match, y models.Match) (bool, string) {
	if x.Id != y.Id {
		return false, "id"
	}

	if !strings.EqualFold(x.TeamA, y.TeamA) ||
		!strings.EqualFold(x.TeamB, y.TeamB) {
		return false, "team names"
	}

	if !x.StartTime.Equal(y.StartTime) {
		return false, "start time"
	}

	if len(x.Maps) != len(y.Maps) {
		return false, "maps"
	}

	for i := 0; i < len(x.Maps); i++ {
		if !strings.EqualFold(x.Maps[i], y.Maps[i]) {
			return false, "maps"
		}
	}

	if x.Viewers != y.Viewers {
		return false, "viewers"
	}

	if !strings.EqualFold(x.PlayerOfTheMatch, y.PlayerOfTheMatch) {
		return false, "player of the match"
	}

	return true, ""
}
