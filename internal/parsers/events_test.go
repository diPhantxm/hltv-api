package parsers

import (
	"bytes"
	"hltvapi/internal/models"
	"hltvapi/internal/urlBuilder/fileUrlBuilder"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestGetEvent(t *testing.T) {
	tests := []struct {
		Id     int
		Result models.Event
	}{
		{6138, models.Event{
			Id:        6138,
			Name:      "IEM Dallas 2022",
			StartDate: time.UnixMilli(1653904800000),
			EndDate:   time.UnixMilli(1654423200000),
			PrizePool: "$250,000",
			Teams: []string{
				"FaZe", "ENCE", "NIP", "G2", "FURIA", "Cloud9", "Vitality", "Astralis",
				"Liquid", "Imperial", "BIG", "Movistar Riders", "MOUZ", "MIBR", "Complexity", "Encore",
			},
			Location: "Dallas, TX, United States",
		}},
		{6372, models.Event{
			Id:        6372,
			Name:      "PGL Major Antwerp 2022",
			StartDate: time.UnixMilli(1652522400000),
			EndDate:   time.UnixMilli(1653213600000),
			PrizePool: "$1,000,000",
			Teams: []string{
				"FaZe", "Natus Vincere", "ENCE", "Cloud9", "Heroic", "G2", "FURIA", "NIP",
				"Outsiders", "Vitality", "Liquid", "Copenhagen Flames", "BIG", "Spirit",
				"Bad News Eagles", "Imperial",
			},
			Location: "Antwerp, Belgium",
		}},
		{6344, models.Event{
			Id:        6344,
			Name:      "Blast Premier Spring Showdown 2022 Europe",
			StartDate: time.UnixMilli(1651053600000),
			EndDate:   time.UnixMilli(1651399200000),
			PrizePool: "$67,500",
			Teams: []string{
				"Heroic", "ENCE", "NIP", "Astralis", "Movistar Riders", "Copenhagen Flames", "Bad News Eagles", "NKT",
			},
			Location: "Europe (Online)",
		}},
		{6137, models.Event{
			Id:        6137,
			Name:      "ESL Pro League Season 15",
			StartDate: time.UnixMilli(1646823600000),
			EndDate:   time.UnixMilli(1649584800000),
			PrizePool: "$823,000",
			Teams: []string{
				"Natus Vincere", "Players", "G2", "FaZe", "Vitality", "Outsiders", "Heroic", "NIP", "Astralis", "BIG",
				"Entropiq", "MOUZ", "FURIA", "ENCE", "GODSENT", "Sprout", "Movistar Riders",
				"Complexity", "Liquid", "Party Astronauts", "Evil Geniuses", "Looking For Org", "AGO", "fnatic",
			},
			Location: "Düsseldorf, Germany",
		}},
		{6317, models.Event{
			Id:        6317,
			Name:      "ESL Challenger Valencia 2022",
			StartDate: time.UnixMilli(1656669600000),
			EndDate:   time.UnixMilli(1656842400000),
			PrizePool: "$100,000",
			Teams: []string{
				"FURIA", "Outsiders", "Movistar Riders", "MIBR", "Sprout", "HUMMER", "Rare Atom", "00NATION",
			},
			Location: "Valencia, Spain",
		}},
		{6140, models.Event{
			Id:        6140,
			Name:      "IEM Cologne 2022",
			StartDate: time.UnixMilli(1657188000000),
			EndDate:   time.UnixMilli(1658052000000),
			PrizePool: "$972,000",
			Teams: []string{
				"FaZe", "Natus Vincere", "ENCE", "Cloud9", "G2", "FURIA", "NIP", "Liquid", "", "", "", "", "", "", "", "",
			},
			Location: "Cologne, Germany",
		}},
	}

	p := NewEventParser(fileUrlBuilder.NewFileUrlBuilder())

	for _, test := range tests {
		event, err := p.GetEvent(test.Id)
		if err != nil {
			t.Errorf("Parse Event Error. Error: %s\n", err.Error())
		}
		event.Id = test.Result.Id

		if ok, field := areEventsEqual(test.Result, *event); !ok {
			t.Errorf("Parse Event %d ended with error. Field: %s\n", test.Result.Id, field)
		}
	}
}

func BenchmarkParseEvent(b *testing.B) {
	const url = "https://www.hltv.org/events/6587/global-esports-tour-dubai-2022"
	response, _ := SendRequest(url)
	body, _ := ioutil.ReadAll(response.Body)

	p := EventParser{}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		b.StartTimer()

		_, err := p.parseEventResponse(response)

		if err != nil {
			b.Fatalf("Error during benchmarkParseEvent. Error: %v\n", err.Error())
		}
	}
}

func areEventsEqual(x models.Event, y models.Event) (bool, string) {
	if x.Id != y.Id {
		return false, "id"
	}

	if !strings.EqualFold(x.Name, y.Name) {
		return false, "name"
	}

	if !x.StartDate.Equal(y.StartDate) {
		return false, "start date"
	}

	if !x.EndDate.Equal(y.EndDate) {
		return false, "end date"
	}

	if !strings.EqualFold(x.PrizePool, y.PrizePool) {
		return false, "prize pool"
	}

	if !strings.EqualFold(x.Location, y.Location) {
		return false, "location"
	}

	if len(x.Teams) != len(y.Teams) {
		return false, "teams"
	}

	for i := range x.Teams {
		if !strings.EqualFold(x.Teams[i], y.Teams[i]) {
			return false, "teams"
		}
	}

	return true, ""
}
