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

func TestGetAllEvents(t *testing.T) {
	tests := []models.Event{
		{
			Id:        6317,
			Name:      "ESL Challenger Valencia 2022",
			StartDate: time.UnixMilli(1656669600000),
			EndDate:   time.UnixMilli(1656842400000),
			PrizePool: "$100,000",
			Teams: []string{
				"FURIA", "Outsiders", "Movistar Riders", "MIBR", "Sprout", "HUMMER", "Rare Atom", "00NATION",
			},
			Location: "Valencia, Spain",
		},
		{
			Id:        6140,
			Name:      "IEM Cologne 2022",
			StartDate: time.UnixMilli(1657188000000),
			EndDate:   time.UnixMilli(1658052000000),
			PrizePool: "$972,000",
			Teams: []string{
				"FaZe", "Natus Vincere", "ENCE", "Cloud9", "G2", "FURIA", "NIP", "Liquid", "", "", "", "", "", "", "", "",
			},
			Location: "Cologne, Germany",
		},
		{
			Id:        6640,
			Name:      "ESL Challenger Melbourne 2022 North America Open Qualifier 1",
			StartDate: time.UnixMilli(1657188000000),
			EndDate:   time.UnixMilli(1657274400000),
			PrizePool: "Spot in Closed Qualifier",
			Teams:     []string{},
			Location:  "North America (Online)",
		},
	}

	p := NewEventParser(fileUrlBuilder.NewFileUrlBuilder())
	events, err := p.GetEvents()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}

	eventsMap := make(map[int]models.Event)
	for _, event := range events {
		eventsMap[event.Id] = event
	}

	if len(tests) > len(events) {
		t.Errorf("Didn't get all events. Length of tests is bigger than length of parsed events.")
	}

	for _, test := range tests {
		if event, ok := eventsMap[test.Id]; !ok {
			t.Errorf("Missing event %d\n", test.Id)
		} else {
			if ok, field := areEventsEqual(event, test); !ok {
				t.Errorf("Events with id %d are not equal. %s", event.Id, field)
			}
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
