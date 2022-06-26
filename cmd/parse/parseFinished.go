package main

import (
	"fmt"
	s "hltvapi/internal/controllers/storedImpl/store"
	"hltvapi/internal/controllers/storedImpl/taskScheduler"
	"hltvapi/internal/parsers"
	"hltvapi/internal/parsers/finished"
	"time"
)

var store *s.Store

func main() {
	dbConfig := s.Config{
		ConnectionString: "user= password= dbname= sslmode=disable",
		Driver:           "postgres",
	}

	var err error

	store, err = s.NewStore(dbConfig)
	if err != nil {
		return
	}

	//parseFinishedMatches()
	parseFinishedEvents()

	scheduler := taskScheduler.GetStraightScheduler()
	timeMs := int64(scheduler.Interval.Milliseconds()) * int64(scheduler.Length())
	timeStr := time.UnixMilli(timeMs).Format("03:04:05")

	fmt.Printf("Estimated time: %s\n", timeStr)

	for !taskScheduler.GetStraightScheduler().IsEmpty() {
	}
}

func parseFinishedEvents() {
	p := finished.FinishedEventParser{}
	ids, err := p.GetAllEventsIds()

	if err != nil {
		return
	}

	for _, id := range ids {
		taskScheduler.GetStraightScheduler().Add(pollEvent, id)
	}
}

func parseFinishedMatches() {
	p := finished.FinishedMatchParser{}
	ids, err := p.GetAllMatchesIds()

	if err != nil {
		return
	}

	for _, id := range ids {
		taskScheduler.GetStraightScheduler().Add(pollMatch, id)
	}
}

func pollEvent(params ...interface{}) {
	id := params[0].([]interface{})[0].(int)

	p := parsers.EventParser{}

	event, err := p.GetEvent(id)
	if err != nil {
		return
	}

	store.Events.AddOrEdit(*event)
}

func pollMatch(params ...interface{}) {
	id := params[0].([]interface{})[0].(int)

	p := parsers.MatchParser{}

	match, err := p.GetMatch(id)
	if err != nil {
		return
	}

	store.Matches.AddOrEdit(*match)
}
