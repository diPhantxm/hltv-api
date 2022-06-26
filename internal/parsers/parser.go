package parsers

import (
	"errors"
	"net/http"
)

func SendRequest(url string) (*http.Response, error) {
	client := &http.Client{}

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "text/html")
	request.Header.Add("User-Agent", "Chrome/23.0.1271.64")
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("request to url was not successful")
	}

	return response, nil
}
