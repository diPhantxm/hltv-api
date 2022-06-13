package parsers

import (
	"errors"
	"hltvapi/internal/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetEvent(url string) (*models.Event, error) {
	response, err := sendRequest(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return parseEventResponseAndId(url, response)
}

func GetUpcomingEventsIds(url string) ([]int, error) {
	response, err := sendRequest(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	eventTags := document.Find(".events-holder").Find(".events-month")
	ids := make([]int, eventTags.Length())

	eventTags.Each(func(i int, selection *goquery.Selection) {
		link, ok := selection.Find(".a-reset").Attr("href")
		if !ok {
			return
		}

		idStr := strings.Split(link, "/")[2]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}

		ids[i] = id
	})

	return ids, nil
}

func parseEventResponseAndId(url string, response *http.Response) (*models.Event, error) {
	p := eventParser{}

	event, err := p.parse(response)
	if err != nil {
		return nil, err
	}

	event.Id, err = p.parseId(url)
	if err != nil {
		return nil, err
	}

	return event, nil
}

type eventParser struct {
}

func (p eventParser) parse(response *http.Response) (*models.Event, error) {
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

func (p eventParser) parseId(url string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	idStr := re.FindString(url)

	id, _ := strconv.Atoi(idStr)

	return id, nil
}

func (p eventParser) parseName(document *goquery.Document) (string, error) {
	nameTag := document.Find(".event-hub-title")

	return nameTag.Text(), nil
}

func (p eventParser) parseStartDate(document *goquery.Document) (time.Time, error) {
	return p.parseDate(document, 0)
}

func (p eventParser) parseEndDate(document *goquery.Document) (time.Time, error) {
	return p.parseDate(document, 2)
}

func (p eventParser) parseDate(document *goquery.Document, index int) (time.Time, error) {
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

func (p eventParser) parsePrizePool(document *goquery.Document) (string, error) {
	prizePoolTag := document.Find(".prizepool").Eq(1)

	return prizePoolTag.Text(), nil
}

func (p eventParser) parseTeams(document *goquery.Document) ([]string, error) {
	teamsAttendingTag := document.Find(".teams-attending")

	teamsTags := teamsAttendingTag.Find(".team-name")

	teams := make([]string, teamsTags.Length())
	teamsTags.Each(func(i int, selection *goquery.Selection) {
		teamName := selection.Find(".text")
		teams[i] = teamName.Text()
	})

	return teams, nil
}

func (p eventParser) parseLocation(document *goquery.Document) (string, error) {
	locationTag := document.Find(".location").Eq(1).Find(".text-ellipsis")

	return locationTag.Text(), nil
}
