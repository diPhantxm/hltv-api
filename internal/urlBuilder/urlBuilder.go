package urlBuilder

/*type UrlBuilder interface {
	PlayersStats()
	TeamsStats()
	Match()
	Event()
	Player()
	Team()
	Results()
	FinishedEvents()
	FinishedMatches()
	AddId(id int)
	AddName(name string)
	AddParam(param string, value string)

	String() string
	Clear()
}*/

type UrlBuilder interface {
	AddId(id int)
	AddName(name string)
	AddParam(param string, value string)

	String() string
	Clear()
}

type MatchUrlBuilder interface {
	UrlBuilder
	Match()
	Results()
}

type EventUrlBuilder interface {
	UrlBuilder
	Event()
	FinishedEvents()
}

type TeamUrlBuilder interface {
	UrlBuilder
	Team()
	TeamsStats()
}

type PlayerUrlBuilder interface {
	UrlBuilder
	Player()
	PlayersStats()
}
