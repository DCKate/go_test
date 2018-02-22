package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
)

//GetResourceURI after get the oauth token, use the token to query the user info
const GetResourceURI = ""

//AutherOAuthURI start oauth, use to get the grant code
const AutherOAuthURI = ""

//GetOAuthTokenURI use the grant code to get token
const GetOAuthTokenURI = ""

type info struct {
	Title  string
	Domain string
	Cid    string
	Csrt   string
}

func (ii *info) combine() string {
	reURL := "http://localhost:55555/oauth_back/"
	rr := ii.Domain + AutherOAuthURI + "?client_id=" + ii.Cid + "&redirect_uri=" + reURL + "&response_type=code&scope=read"
	return rr
}

var ginfo info

func parseRepbyKey(repdata []byte, key string) string {
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(repdata, &objmap)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	var str string
	err = json.Unmarshal(*objmap[key], &str)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	return str
}

func getResource(tok string, urlpath string) []byte {
	autok := fmt.Sprintf("Bearer %v", tok)
	req, err := http.NewRequest("GET", urlpath, nil)
	req.Header.Add("Authorization", autok)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response : %v \n", string(body))
	return body
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *info) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func oauthbkHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.PostForm(ginfo.Domain+GetOAuthTokenURI,
		url.Values{"client_id": {ginfo.Cid},
			"client_secret": {ginfo.Csrt},
			"code":          {r.FormValue("code")},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {"http://localhost:55555/oauth_back/"}})

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
	tok := parseRepbyKey(body, "access_token")
	rep := getResource(tok, GetResourceURI)
	user := parseRepbyKey(rep, "username")

	fmt.Fprintf(w, "<h1>Hello %s</h1>", user)

}

func knockHandler(w http.ResponseWriter, r *http.Request) {
	bk := info{
		Title:  "",
		Domain: r.FormValue("domain"),
		Cid:    r.FormValue("id"),
		Csrt:   r.FormValue("secret"),
	}
	aurl := bk.combine()
	ginfo = bk
	http.Redirect(w, r, aurl, http.StatusSeeOther)
}

func doorHandler(w http.ResponseWriter, r *http.Request) {
	data := &info{
		Title:  "Enter the App info",
		Domain: "where",
		Cid:    "who",
		Csrt:   "secret",
	}
	renderTemplate(w, "knock", data)
}

func main() {
	http.HandleFunc("/oauth_back/", oauthbkHandler)
	http.HandleFunc("/knock/", doorHandler)
	http.HandleFunc("/goknock/", knockHandler)
	http.ListenAndServe(":55555", nil)
}
