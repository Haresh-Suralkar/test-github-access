package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Req struct {
	Method string
	URL    string
	Data   *string
}

func (r Req) getGithubData(accessToken string) string {

	var (
		req    *http.Request
		reqerr error
	)

	if r.Data != nil {
		log.Println(*r.Data)
		req, reqerr = http.NewRequest(r.Method, "https://api.github.com"+r.URL, strings.NewReader(*r.Data))
	} else {
		req, reqerr = http.NewRequest(r.Method, "https://api.github.com"+r.URL, nil)
	}

	if reqerr != nil {
		log.Panic("API Request creation failed", reqerr.Error())
		return ""
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if accessToken != "" {
		authorizationHeaderValue := fmt.Sprintf("Bearer %s", accessToken)
		req.Header.Set("Authorization", authorizationHeaderValue)
	}

	log.Printf("%+v\n", req)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed: ", resperr.Error())
		return ""
	}

	respbody, _ := ioutil.ReadAll(resp.Body)
	return string(respbody)
}
