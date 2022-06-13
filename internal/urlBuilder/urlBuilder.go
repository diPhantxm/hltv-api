package urlBuilder

import (
	"fmt"
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

	baseUrl      = `https://hltv.org/`
	matchUrl     = `matches/`
	eventUrl     = `events/`
	playerUrl    = `player/`
	teamUrl      = `team/`
	playersStats = `stats/players/`
	teamsStats   = `stats/teams/`
)

type UrlBuilder struct {
	url map[string]string
}

func NewUrlBuilder() *UrlBuilder {
	return &UrlBuilder{
		url: make(map[string]string),
	}
}

func (b *UrlBuilder) PlayersStats() {
	b.url["group"] = playersStats
}

func (b *UrlBuilder) TeamsStats() {
	b.url["group"] = teamsStats
}

func (b *UrlBuilder) Match() {
	b.url["group"] = matchUrl
}

func (b *UrlBuilder) Event() {
	b.url["group"] = eventUrl
}

func (b *UrlBuilder) Player() {
	b.url["group"] = playerUrl
}

func (b *UrlBuilder) Team() {
	b.url["group"] = teamUrl
}

func (b *UrlBuilder) AddId(id int) {
	b.url["id"] = fmt.Sprintf("%d/", id)
}

func (b *UrlBuilder) AddName(name string) {
	b.url["name"] = fmt.Sprintf("%s?", name)
}

func (b *UrlBuilder) String() string {
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
