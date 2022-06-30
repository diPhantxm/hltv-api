package fileUrlBuilder

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	format = []string{
		"baseUrl",
		"group",
		"id",
		"name",
		"params",
	}

	baseUrl         = `file:///`
	matchUrl        = `matches/`
	eventUrl        = `events/`
	playerUrl       = `players/`
	teamUrl         = `teams/`
	playersStats    = `players/`
	teamsStats      = `teams/`
	resultsUrl      = `matches/`
	finishedEvents  = `events/archive/`
	finishedMatches = `matches/`
)

type FileUrlBuilder struct {
	url map[string]string
}

func NewFileUrlBuilder() *FileUrlBuilder {
	return &FileUrlBuilder{
		url: make(map[string]string),
	}
}

func (b *FileUrlBuilder) PlayersStats() {
	b.url["group"] = playersStats
}

func (b *FileUrlBuilder) TeamsStats() {
	b.url["group"] = teamsStats
}

func (b *FileUrlBuilder) Match() {
	b.url["group"] = matchUrl
}

func (b *FileUrlBuilder) Event() {
	b.url["group"] = eventUrl
}

func (b *FileUrlBuilder) Player() {
	b.url["group"] = playerUrl
}

func (b *FileUrlBuilder) Team() {
	b.url["group"] = teamUrl
}

func (b *FileUrlBuilder) Results() {
	b.url["group"] = resultsUrl
}

func (b *FileUrlBuilder) FinishedEvents() {
	b.url["group"] = finishedEvents
}

func (b *FileUrlBuilder) FinishedMatches() {
	b.url["group"] = finishedMatches
}

func (b *FileUrlBuilder) AddId(id int) {
	b.url["id"] = strconv.Itoa(id)
}

func (b *FileUrlBuilder) AddName(name string) {
	b.url["name"] = name
}

func (b *FileUrlBuilder) AddParam(param string, value string) {
	b.url["params"] += fmt.Sprintf("%s=%s&", param, value)
}

func (b *FileUrlBuilder) String() string {
	url := strings.Builder{}
	url.WriteString(baseUrl)

	for _, field := range format {
		if value, ok := b.url[field]; ok {
			url.WriteString(value)
		}
	}

	_, idOk := b.url["id"]
	_, paramOk := b.url["params"]
	if !(idOk || paramOk) {
		url.WriteString("all")
	}

	url.WriteString(".html")

	return url.String()
}

func (b *FileUrlBuilder) Clear() {
	b.url = make(map[string]string)
}
