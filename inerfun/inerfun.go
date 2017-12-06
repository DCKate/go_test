package interfun

import (
    "fmt"
    "crypto/tls"
    "bytes"
    "net/http"
    "net/url"
    "io/ioutil"
)

func MakePost(url string,headers map[string]string,urlbody []byte)(int,[]byte){

    var jsonStr = []byte(urlbody)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    for kk,vv:=range headers {
        // req.Header.Set("authorization", token)
        // req.Header.Set("Accept", "application/json")
        // req.Header.Set("Content-Type", "application/json")
        req.Header.Set(kk, vv)
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

    // fmt.Println("response Status:", resp.Status)
    // fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)

    return resp.StatusCode,body
}

func MakeGet(urlpath string,headers map[string]string,para map[string]string)(int,[]byte){
    
    req, err := http.NewRequest("GET", urlpath, nil)
    for kk,vv:=range headers {
        req.Header.Set(kk, vv)
    }

    q := req.URL.Query()
    for pk,pv:=range para {
        q.Set(pk,pv)    
    }
    req.URL.RawQuery = q.Encode()
    // fmt.Println(req.URL.String())

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

    return resp.StatusCode,body
}

func MakeSimpleGetToByte(urlpath string)(int,[]byte) {
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
  return resp.StatusCode,body
}

func MakeSimpleGetToStr(urlpath string)(int,string) {
    resp, err := http.Get(urlpath)
    if err != nil {
        // handle error
    }

    defer resp.Body.Close()
    buf := new(bytes.Buffer)
    buf.ReadFrom(resp.Body)
    s := buf.String() // Does a complete copy of the bytes in the buffer.
    return resp.StatusCode,s
}

func httpPostForm(urlpath string) {
    resp, err := http.PostForm(urlpath,
        url.Values{"key": {"Value"}, "id": {"123"}})

    if err != nil {
        fmt.Println("error:", err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("error:", err)
    }

    fmt.Println(string(body))

}