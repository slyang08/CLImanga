package manga

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

var baseURL string = "https://api.mangadex.org"

func getFullBuiltUrl(apiURL *string, params *url.Values) string {
	return *apiURL + "?" + params.Encode()
}

func makeGETApiRequest(fullAPIURL *string) (map[string]interface{}, error) {
	response, err := http.Get(*fullAPIURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading GET response body %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return result, nil
}

func FetchMangasByNameSearch(mangaName *string) (map[string]string, error) {
	apiURl := baseURL + "/manga"

	params := url.Values{}
	params.Add("title", *mangaName)

	var fullAPIURL string = getFullBuiltUrl(&apiURl, &params)

	data, err := makeGETApiRequest(&fullAPIURL)
	if err != nil {
		return nil, err
	}

	mangasFound := make(map[string]string)

	dataArr, ok := data["data"].([]interface{}) // json path: data(id)(attributes(title(en,..)))

	if !ok {
		return nil, fmt.Errorf("couldnt get data") // check if field exists
	}

	for _, item := range dataArr {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		mangaID := itemMap["id"].(string) // hashCodeID (like asdad-asdad-asda)

		attributes, ok := itemMap["attributes"].(map[string]interface{})
		if !ok {
			continue
		}

		if title, ok := attributes["title"].(map[string]interface{}); ok {
			if en, ok := title["en"].(string); ok {
				mangasFound[mangaID] = en
			}
		}
	}
	return mangasFound, nil
}

func GetAllChapterListOfManga(mangaID *string) (map[string]string, error) { // https://api.mangadex.org/docs/04-chapter/search/
	fullAPIURL := baseURL + "/manga/" + *mangaID + "/feed"

	data, err := makeGETApiRequest(&fullAPIURL)
	if err != nil {
		return nil, err
	}

	dataArr, ok := data["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("couldnt get data")
	}

	for _, item := range dataArr {
		itemMap, ok := item.(map[string]interface{})

		if !ok {
			continue
		}

		attributes, ok := itemMap["attributes"].(map[string]interface{})

		if !ok {
			return nil, fmt.Errorf("Couldnt get chapter list")
		}

		log.Println(attributes)
	}

	return nil, nil
}
