package interfun

import (
    "fmt"
    "crypto/tls"
    "bytes"
    "net/http"
    "io/ioutil"
)

func MakeJsonPost(url string,token string,urlbody []byte)(string,[]byte){

    var jsonStr = []byte(urlbody)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("authorization", token)
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
    }
    // fmt.Printf("Request : %v \n", req)
    client := &http.Client{Transport:tr}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("error:", err)
        // panic(err)
    }
    defer resp.Body.Close()

    // fmt.Println("response Status:", resp.Status)
    // fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)

    return resp.Status,body
}

func MakePost(url string,headers map[string]string,urlbody []byte)(string,[]byte){

    var jsonStr = []byte(urlbody)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    for kk,vv:=range headers {
        req.Header.Set(kk, vv)
    }
    // req.Header.Set("authorization", token)
    // req.Header.Set("Accept", "application/json")
    // req.Header.Set("Content-Type", "application/json")
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
    }
    // fmt.Printf("Request : %v \n", req)
    client := &http.Client{Transport:tr}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("error:", err)
        // panic(err)
    }
    defer resp.Body.Close()

    // fmt.Println("response Status:", resp.Status)
    // fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)

    return resp.Status,body
}

func MakeGet(urlpath string,headers map[string]string)(string,[]byte){
	
	req, err := http.NewRequest("urlpath", urlpath, nil)
    for kk,vv:=range headers {
        req.Header.Add(kk, vv)
    }
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
    }
    // fmt.Printf("Request : %v \n", req)
    client := &http.Client{Transport:tr}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("error:", err)
        // panic(err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    return resp.Status,body
}

func MakeSimpleGet(urlpath string)(string,[]byte) {
	resp, err := http.Get(urlpath)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
        // panic(err)
	}

	// fmt.Println(string(body))
	 return resp.Status,body
}

func httpPostForm(urlpath string) {
	resp, err := http.PostForm(urlpath,
		url.Values{"key": {"Value"}, "id": {"123"}})

	if err != nil {
		mt.Println("error:", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mt.Println("error:", err)
	}

	fmt.Println(string(body))

}