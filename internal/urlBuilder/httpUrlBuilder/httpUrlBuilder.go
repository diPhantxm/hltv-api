package httpUrlBuilder

import (
	"fmt"
	"strings"
)

type HttpUrlBuilder struct {
	url map[string]string
}

var (
	format = []string{
		"baseUrl",
		"group",
		"id",
		"name",
		"params",
	}

	baseUrl         = `https://hltv.org/`
	matchUrl        = `matches/`
	eventUrl        = `events/`
	playerUrl       = `player/`
	teamUrl         = `team/`
	playersStats    = `stats/players/`
	teamsStats      = `stats/teams/`
	resultsUrl      = `results/`
	finishedEvents  = `events/archive/`
	finishedMatches = `results/`
)

func NewHttpUrlBuilder() *HttpUrlBuilder {
	return &HttpUrlBuilder{
		url: make(map[string]string),
	}
}

func (b *HttpUrlBuilder) PlayersStats() {
	b.url["group"] = playersStats
}

func (b *HttpUrlBuilder) TeamsStats() {
	b.url["group"] = teamsStats
}

func (b *HttpUrlBuilder) Match() {
	b.url["group"] = matchUrl
}

func (b *HttpUrlBuilder) Event() {
	b.url["group"] = eventUrl
}

func (b *HttpUrlBuilder) Player() {
	b.url["group"] = playerUrl
}

func (b *HttpUrlBuilder) Team() {
	b.url["group"] = teamUrl
}

func (b *HttpUrlBuilder) FinishedEvents() {
	b.url["group"] = finishedEvents
}

func (b *HttpUrlBuilder) FinishedMatches() {
	b.url["group"] = finishedMatches
}

func (b *HttpUrlBuilder) AddId(id int) {
	b.url["id"] = fmt.Sprintf("%d/", id)
}

func (b *HttpUrlBuilder) AddName(name string) {
	b.url["name"] = fmt.Sprintf("%s?", name)
}

func (b *HttpUrlBuilder) Results() {
	b.url["group"] = resultsUrl
}

func (b *HttpUrlBuilder) AddParam(param string, value string) {
	if b.url["params"] == "" {
		b.url["params"] += "?"
	}
	b.url["params"] += fmt.Sprintf("%s=%s&", param, value)
}

func (b *HttpUrlBuilder) String() string {
	group := b.url["group"]
	if group == playerUrl || group == matchUrl || group == eventUrl || group == teamUrl {
		if b.url["id"] != "" && b.url["name"] == "" {
			b.url["name"] = "_"
		}
	}

	url := strings.Builder{}
	url.WriteString(baseUrl)

	for _, field := range format {
		if value, ok := b.url[field]; ok {
			url.WriteString(value)
		}
	}

	return url.String()
}
