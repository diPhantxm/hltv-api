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
		Id     int
		Result models.Match
	}{
		{2356303, models.Match{
			Id:               2356303,
			TeamA:            "ENCE",
			TeamB:            "Cloud9",
			StartTime:        time.UnixMilli(1654444800000),
			Maps:             []string{"Mirage", "Overpass", "Ancient", "Vertigo", "Dust2"},
			Viewers:          0,
			PlayerOfTheMatch: "sh1ro",
		}},
		{2356302, models.Match{
			Id:               2356302,
			TeamA:            "BIG",
			TeamB:            "Cloud9",
			StartTime:        time.UnixMilli(1654371000000),
			Maps:             []string{"Overpass", "Dust2", "Ancient"},
			Viewers:          0,
			PlayerOfTheMatch: "HObbit",
		}},
		{2356301, models.Match{
			Id:               2356301,
			TeamA:            "ENCE",
			TeamB:            "FURIA",
			StartTime:        time.UnixMilli(1654358400000),
			Maps:             []string{"Vertigo", "Nuke", "Mirage"},
			Viewers:          0,
			PlayerOfTheMatch: "Snax",
		}},
		{2356300, models.Match{
			Id:               2356300,
			TeamA:            "FaZe",
			TeamB:            "Cloud9",
			StartTime:        time.UnixMilli(1654296900000),
			Maps:             []string{"Overpass", "Inferno", "Mirage"},
			Viewers:          0,
			PlayerOfTheMatch: "Ax1le",
		}},
		{2356298, models.Match{
			Id:               2356298,
			TeamA:            "Liquid",
			TeamB:            "Cloud9",
			StartTime:        time.UnixMilli(1654124400000),
			Maps:             []string{"Ancient", "Vertigo", "Overpass"},
			Viewers:          0,
			PlayerOfTheMatch: "Ax1le",
		}},
		{2356296, models.Match{
			Id:               2356296,
			TeamA:            "ENCE",
			TeamB:            "FaZe",
			StartTime:        time.UnixMilli(1654111800000),
			Maps:             []string{"Mirage", "Ancient", "Nuke"},
			Viewers:          0,
			PlayerOfTheMatch: "Snax",
		}},
	}

	p := MatchParser{}

	for _, test := range tests {
		actual, err := p.GetMatch(test.Id)

		if err != nil {
			t.Errorf("Parse Match %d ended with error. Error: %s\n", test.Id, err.Error())
		}

		if ok, field := areMatchesEqual(*actual, test.Result); !ok {
			t.Errorf("Parse Match %d was incorrect. Field: %s\n", test.Id, field)
		}
	}
}

func TestGetIdsByDate(t *testing.T) {
	tests := []struct {
		Url    string
		Result []int
	}{
		{"https://www.hltv.org/results?startDate=2022-06-15&endDate=2022-06-15", []int{2356632}},
		{"https://www.hltv.org/results/?startDate=2022-06-14&endDate=2022-06-14&", []int{2356788, 2356787}},
		{"https://www.hltv.org/results/?startDate=2022-06-13&endDate=2022-06-13&", []int{2356786, 2356785, 2356834, 2356826, 2356819}},
		{"https://www.hltv.org/results?startDate=2022-06-13&endDate=2022-06-13", []int{2356786, 2356785, 2356834, 2356826, 2356819}},
	}

	p := MatchParser{}

	for _, test := range tests {
		ids, err := p.getMatchesIdsByDate(test.Url)
		if err != nil {
			t.Fatalf(err.Error())
		}

		idsMap := make(map[int]struct{})

		for _, id := range ids {
			idsMap[id] = struct{}{}
		}

		for _, id := range test.Result {
			if _, ok := idsMap[id]; !ok {
				t.Errorf("There is no id %d in result\n", id)
			}
		}
	}
}

func BenchmarkParseMatch(b *testing.B) {
	const url = "https://www.hltv.org/matches/2356673/nip-vs-nasr-global-esports-tour-dubai-2022"
	response, _ := sendRequest(url)
	body, _ := ioutil.ReadAll(response.Body)

	p := MatchParser{}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		b.StartTimer()

		_, err := p.parseMatchResponseAndId(url, response)

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
