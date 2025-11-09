package manga

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var baseURL string = "https://api.mangadex.org"

func getFullBuiltURL(apiURL *string, params *url.Values) string {
	return *apiURL + "?" + params.Encode()
}

func makeGETApiRequest(fullAPIURL *string) (map[string]any, error) {
	response, err := http.Get(*fullAPIURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading GET response body %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return result, nil
}

func FetchMangasByNameSearch(mangaName *string) (map[string]string, error) {
	apiURl := baseURL + "/manga"

	params := url.Values{}
	params.Add("title", *mangaName)

	fullAPIURL := getFullBuiltURL(&apiURl, &params)

	data, err := makeGETApiRequest(&fullAPIURL)
	if err != nil {
		return nil, err
	}

	mangasFound := make(map[string]string)

	dataArr, ok := data["data"].([]any) // json path: data(id)(attributes(title(en,..)))

	if !ok {
		return nil, fmt.Errorf("couldnt get data") // check if field exists
	}

	for _, item := range dataArr {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		mangaID := itemMap["id"].(string) // hashCodeID (like asdad-asdad-asda)

		attributes, ok := itemMap["attributes"].(map[string]any)
		if !ok {
			continue
		}

		if title, ok := attributes["title"].(map[string]any); ok {
			if en, ok := title["en"].(string); ok {
				mangasFound[mangaID] = en
			}
		}
	}
	return mangasFound, nil
}

func GetAllChapterListOfManga(mangaID *string) (map[string]map[string]any, error) { // https://api.mangadex.org/docs/04-chapter/search/
	fullAPIURL := baseURL + "/manga/" + *mangaID + "/feed"

	data, err := makeGETApiRequest(&fullAPIURL)
	if err != nil {
		return nil, err
	}

	dataArr, ok := data["data"].([]any)
	if !ok {
		return nil, fmt.Errorf("couldnt get data")
	}

	chapterList := make(map[string]map[string]any) // path: data(id)(attributes)

	for _, item := range dataArr {
		itemMap, ok := item.(map[string]any)

		if !ok {
			continue
		}

		chapterID := itemMap["id"].(string)

		attributes, ok := itemMap["attributes"].(map[string]any)

		if !ok {
			return nil, fmt.Errorf("couldnt get chapter list")
		}

		chapterList[chapterID] = attributes

	}

	return chapterList, nil
}
