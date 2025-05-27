package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func fetch(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

func EnrichPerson(name string) (age *int, gender *string, nationality *string) {
	type ageResp struct {
		Age int `json:"age"`
	}
	type genderResp struct {
		Gender string `json:"gender"`
	}
	type natResp struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}

	var a ageResp
	var g genderResp
	var n natResp

	fetch(fmt.Sprintf("https://api.agify.io/?name=%s", name), &a)
	fetch(fmt.Sprintf("https://api.genderize.io/?name=%s", name), &g)
	fetch(fmt.Sprintf("https://api.nationalize.io/?name=%s", name), &n)

	return &a.Age, &g.Gender, func() *string {
		if len(n.Country) > 0 {
			return &n.Country[0].CountryID
		}
		return nil
	}()
}
