package urlBuilder

import (
	"strings"
	"testing"
)

func TestUrlBuilder(t *testing.T) {
	tests := []struct {
		action func() string
		Result string
	}{
		{func() string {
			builder := NewUrlBuilder()
			builder.Match()
			builder.AddId(123)
			return builder.String()
		}, "https://hltv.org/matches/123/_"},
		{func() string {
			builder := NewUrlBuilder()
			builder.Event()
			builder.AddId(123)
			return builder.String()
		}, "https://hltv.org/events/123/_"},
		{func() string {
			builder := NewUrlBuilder()
			builder.Team()
			builder.AddId(123)
			return builder.String()
		}, "https://hltv.org/team/123/_"},
		{func() string {
			builder := NewUrlBuilder()
			builder.Player()
			builder.AddId(123)
			return builder.String()
		}, "https://hltv.org/player/123/_"},
		{func() string {
			builder := NewUrlBuilder()
			builder.Match()
			builder.AddId(123)
			builder.AddName("fnatic-vs-natus-vincere")
			return builder.String()
		}, "https://hltv.org/matches/123/fnatic-vs-natus-vincere?"},
		{func() string {
			builder := NewUrlBuilder()
			builder.Team()
			builder.AddId(12546)
			builder.AddName("natus-vincere")
			return builder.String()
		}, "https://hltv.org/team/12546/natus-vincere?"},
		{func() string {
			builder := NewUrlBuilder()
			builder.Event()
			builder.AddId(1)
			builder.AddName("esl-pro-league-season-15")
			return builder.String()
		}, "https://hltv.org/events/1/esl-pro-league-season-15?"},
		{func() string {
			builder := NewUrlBuilder()
			builder.Event()
			return builder.String()
		}, "https://hltv.org/events/"},
	}

	for _, test := range tests {
		got := test.action()
		if !strings.EqualFold(got, test.Result) {
			t.Errorf("Got: %s, Actual: %s", got, test.Result)
		}
	}
}
