package parsers

import (
	"errors"
	"net/http"
)

func SendRequest(url string) (*http.Response, error) {
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("./../../test pages/")))

	client := &http.Client{Transport: t}

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "text/plain")
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
