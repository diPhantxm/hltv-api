package urlBuilder

import (
	"fmt"
	"strings"
)

var (
	format = map[string]int{
		"baseUrl": 0,
		"group":   1,
		"id":      2,
		"name":    3,
		"params":  4,
	}

	baseUrl   = `https://hltv.org/`
	matchUrl  = `matches/`
	eventUrl  = `events/`
	playerUrl = `players/`
	teamUrl   = `teams/`
)

type UrlBuilder struct {
	url string
}

func NewUrlBuilder() *UrlBuilder {
	return &UrlBuilder{
		url: baseUrl,
	}
}

func (b *UrlBuilder) Match() {
	b.url += matchUrl
}

func (b *UrlBuilder) Event() {
	b.url += eventUrl
}

func (b *UrlBuilder) Player() {
	b.url += playerUrl
}

func (b *UrlBuilder) Team() {
	b.url += teamUrl
}

func (b *UrlBuilder) AddId(id int) {
	b.url += fmt.Sprintf("%d/", id)
}

func (b *UrlBuilder) AddName(name string) {
	b.url += name
}

func (b *UrlBuilder) String() string {
	fieldsUrl := b.url[len(baseUrl):]
	fields := strings.Split(fieldsUrl, "/")

	if len(fields) <= format["name"] || fields[format["name"]-1] == "" {
		b.url += "_"
	}

	return b.url
}
