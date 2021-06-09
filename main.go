package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// init() executes before the main program
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

// Represents the token response received from Github
type githubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// declared global as it will be required throughout the application
var (
	ghresp githubAccessTokenResponse
	user   githubUser
)

func main() {

	// routes
	// Simply returns a link to the login route
	http.HandleFunc("/", rootHandler)

	// Login route
	http.HandleFunc("/github/login", githubLoginHandler)

	// Github callback
	http.HandleFunc("/github/callback", githubCallbackHandler)

	http.HandleFunc("/github/repo", githubRepoHandler)

	// Route where the authenticated user is redirected to
	http.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		loggedinHandler(w, r, "")
	})

	// Listen and serve on port 3000
	// Intentinally kept 3000 port hardcoded as it is configured in the github settings
	fmt.Println("Server is listning on PORT 3000 ")
	log.Panic(
		http.ListenAndServe(":3000", nil),
	)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="/github/login">Login with GitHub</a>`)
}

func loggedinHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" && ghresp.AccessToken == "" {
		// Unauthorized users get an unauthorized message
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	w.Header().Set("Content-type", "text/html")

	json.Unmarshal([]byte(githubData), &user)

	list := (githubRepos{}).listRepos().html()

	// Return the prettified JSON as a string
	fmt.Fprintf(w, list)
}

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := getGithubClientID()

	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", githubClientID, "http://localhost:3000/github/callback")

	http.Redirect(w, r, redirectURL, 301)
}

func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	githubAccessToken := getGithubAccessToken(code)
	url := `/user`

	rr := Req{
		Method: "GET",
		URL:    url,
	}

	githubData := rr.getGithubData(githubAccessToken)

	if githubData == "" {
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	loggedinHandler(w, r, githubData)
}

func getGithubAccessToken(code string) string {

	clientID := getGithubClientID()
	clientSecret := getGithubClientSecret()

	requestBodyMap := map[string]string{"client_id": clientID, "client_secret": clientSecret, "code": code}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, reqerr := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(requestJSON))
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respbody, &ghresp)

	log.Println(ghresp)
	return ghresp.AccessToken
}

func getGithubClientID() string {

	githubClientID, exists := os.LookupEnv("CLIENT_ID")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}

	return githubClientID
}

func getGithubClientSecret() string {

	githubClientSecret, exists := os.LookupEnv("CLIENT_SECRET")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}

	return githubClientSecret
}
