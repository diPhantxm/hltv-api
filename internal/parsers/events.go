package parsers

import (
	"errors"
	"hltvapi/internal/models"
	"hltvapi/internal/urlBuilder"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type EventParser struct {
	builder urlBuilder.UrlBuilder
}

func NewEventParser(builder urlBuilder.UrlBuilder) *EventParser {
	return &EventParser{
		builder: builder,
	}
}

func (p EventParser) GetEvent(id int) (*models.Event, error) {
	p.builder.Clear()
	p.builder.Event()
	p.builder.AddId(id)

	response, err := SendRequest(p.builder.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	event, err := p.parseEventResponse(response)
	if err != nil {
		return nil, err
	}

	event.Id = id

	return event, nil
}

func (p EventParser) GetEvents() ([]models.Event, error) {
	ids, err := p.GetUpcomingEventsIds()
	if err != nil {
		return nil, err
	}
	result := make([]models.Event, len(ids))

	for i, id := range ids {
		event, err := p.GetEvent(id)
		if err != nil {
			continue
		}

		result[i] = *event
	}

	return result, nil
}

func (p EventParser) GetUpcomingEventsIds() ([]int, error) {
	p.builder.Clear()
	p.builder.Event()

	response, err := SendRequest(p.builder.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	eventTags := document.Find(".events-holder").Find(".a-reset")
	var ids []int

	eventTags.Each(func(i int, selection *goquery.Selection) {
		link, ok := selection.Attr("href")
		if !ok {
			return
		}

		re := regexp.MustCompile(`\/\w+\/(\d+)`)
		idStr := re.FindStringSubmatch(link)[1]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}

		ids = append(ids, id)
	})

	return ids, nil
}

func (p EventParser) parseEventResponse(response *http.Response) (*models.Event, error) {
	event, err := p.parse(response)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (p EventParser) parse(response *http.Response) (*models.Event, error) {
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	event := &models.Event{}

	event.Name, err = p.parseName(document)
	if err != nil {
		return nil, err
	}

	event.StartDate, err = p.parseStartDate(document)
	if err != nil {
		return nil, err
	}

	event.EndDate, err = p.parseEndDate(document)
	if err != nil {
		return nil, err
	}

	event.PrizePool, err = p.parsePrizePool(document)
	if err != nil {
		return nil, err
	}

	event.Teams, err = p.parseTeams(document)
	if err != nil {
		return nil, err
	}

	event.Location, err = p.parseLocation(document)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (p EventParser) parseName(document *goquery.Document) (string, error) {
	nameTag := document.Find(".event-hub-title")

	return nameTag.Text(), nil
}

func (p EventParser) parseStartDate(document *goquery.Document) (time.Time, error) {
	return p.parseDate(document, 0)
}

func (p EventParser) parseEndDate(document *goquery.Document) (time.Time, error) {
	return p.parseDate(document, 2)
}

func (p EventParser) parseDate(document *goquery.Document, index int) (time.Time, error) {
	dateTag := document.Find(".eventdate")

	startTag := dateTag.Eq(1).Find("span").Eq(index)
	unixTimeStr, ok := startTag.Attr("data-unix")

	if !ok {
		return time.Now(), errors.New("start/end date was not found on page")
	}

	unixTime, err := strconv.ParseInt(unixTimeStr, 10, 64)
	if err != nil {
		return time.Now(), errors.New("unix time was in incorrect format")
	}

	return time.UnixMilli(unixTime), nil
}

func (p EventParser) parsePrizePool(document *goquery.Document) (string, error) {
	prizePoolTag := document.Find(".prizepool").Eq(1)

	return prizePoolTag.Text(), nil
}

func (p EventParser) parseTeams(document *goquery.Document) ([]string, error) {
	teamsAttendingTag := document.Find(".teams-attending")

	teamsTags := teamsAttendingTag.Find(".team-name")

	teams := make([]string, teamsTags.Length())
	teamsTags.Each(func(i int, selection *goquery.Selection) {
		teamName := selection.Find(".text")
		teams[i] = teamName.Text()
	})

	return teams, nil
}

func (p EventParser) parseLocation(document *goquery.Document) (string, error) {
	locationTag := document.Find(".location").Eq(1).Find(".text-ellipsis")

	return locationTag.Text(), nil
}
