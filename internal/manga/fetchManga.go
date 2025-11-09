package manga

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var baseUrl string = "https://api.mangadex.org"

func getFullBuiltUrl(apiUrl *string, params *url.Values) string {
	return *apiUrl + "?" + params.Encode()
}

func makeGETApiRequest(fullApiUrl *string) (map[string]interface{}, error) {
	response, err := http.Get(*fullApiUrl)
	if err != nil {
		return nil, fmt.Errorf("Error making GET request")
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading GET response body")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return result, nil
}

func fetchMangaNames(mangaName *string) ([]string, error) {
	var apiUrl string = baseUrl + "/manga"

	params := url.Values{}
	params.Add("title", *mangaName)

	var fullApiUrl string = getFullBuiltUrl(&apiUrl, &params)

	data, err := makeGETApiRequest(&fullApiUrl)
	if err != nil {
		return nil, err
	}

	fmt.Println(data)
	return []string{"lol"}, nil
}
