package urlBuilder

type UrlBuilder interface {
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
}
