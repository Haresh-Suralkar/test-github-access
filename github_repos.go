package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type githubUser struct {
	Login    string `json:"login"`
	ReposURL string `json:"repos_url"`
	// GlobalMap
}

type GlobalMap map[string]interface{}

type githubRepo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	// GlobalMap
}

type githubRepos []githubRepo

func (gr githubRepos) listRepos() githubRepos {
	err := gr.getRepos()

	if err != nil {
		log.Println("error in get repos: ", err.Error())
		return nil
	}

	return gr
}

func (gr *githubRepos) getRepos() error {

	rr := Req{
		Method: "GET",
		URL:    strings.Replace(user.ReposURL, "https://api.github.com", "", 1),
		Data:   nil,
	}
	resp := rr.getGithubData(ghresp.AccessToken)

	json.Unmarshal([]byte(resp), gr)
	return nil
}

func (gr githubRepos) html() string {
	var li []string

	for _, r := range gr {
		li = append(li, tagLi(r.URL, r.Name))
	}

	return tagUl(li)
}

func githubRepoHandler(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query().Get("u")
	uu, _ := url.QueryUnescape(u)
	uu += "/git/refs"

	var gg []GlobalMap
	rr := Req{
		Method: "GET",
		URL:    uu,
	}
	resp := rr.getGithubData(ghresp.AccessToken)

	json.Unmarshal([]byte(resp), &gg)
	// log.Println(rr, gg)

	// create new branch from head
	sha := ""
	for k, v := range (gg[0]["object"]).(map[string]interface{}) {
		if k == "sha" {
			sha = v.(string)
			break
		}
	}
	b := createBranch(sha, uu)

	log.Println("branch name ==> ", b)
	fmt.Fprintf(w, resp)
}

func createBranch(sha, u string) string {
	newBranch := "test-branch-" + fmt.Sprint(time.Now().Unix())
	payload := GlobalMap{
		// "ref": ("refs/heads/" + newBranch),
		"ref": ("refs/heads/aaa"),
		"sha": sha,
	}

	p, _ := json.Marshal(payload)
	ps := string(p)
	// log.Println(string(p), u)

	var gg []GlobalMap
	rr := Req{
		Method: "POST",
		URL:    u,
		Data:   &ps,
	}
	resp := rr.getGithubData(ghresp.AccessToken)

	json.Unmarshal([]byte(resp), &gg)
	log.Println(resp)

	return newBranch
}
