package inerfun

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func MakePost(url string, headers map[string]string, urlbody []byte) (int, []byte) {

	var jsonStr = []byte(urlbody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	for kk, vv := range headers {
		// req.Header.Set("authorization", token)
		// req.Header.Set("Accept", "application/json")
		// req.Header.Set("Content-Type", "application/json")
		req.Header.Set(kk, vv)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	// fmt.Printf("Request : %v \n", req)
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return -1, nil
		// panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	return resp.StatusCode, body
}

func MakePostForm(apiurl string, headers map[string]string, postform map[string]string) (int, []byte) {
	data := url.Values{}
	for kk, vv := range postform {
		// data.Set(kk, vv)
		data.Add(kk, vv)
	}

	req, _ := http.NewRequest("POST", apiurl, strings.NewReader(data.Encode())) // URL-encoded payload
	for kk, vv := range headers {
		req.Header.Set(kk, vv)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	// fmt.Printf("Request : %v \n", req)
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return -1, nil
		// panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	return resp.StatusCode, body
}

func MakePutFile(apiurl string, headers map[string]string, data io.Reader) (int, []byte) {

	cl, _ := strconv.ParseInt(headers["Content-Length"], 10, 64)
	req, err := http.NewRequest(http.MethodPut, apiurl, data)
	req.ContentLength = cl
	for kk, vv := range headers {
		if kk != "Content-Length" {
			req.Header.Add(kk, vv)
		}
	}
	// fmt.Println(req.Header["Content-Length"])
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	// dump, _ := httputil.DumpRequest(req, true)
	// fmt.Println(string(dump))
	// fmt.Printf("Request : %v \n", req)
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return -1, nil
		// panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response : %v \n", string(body))
	return resp.StatusCode, body
}

func MakeGet(urlpath string, headers map[string]string, para map[string]string) (int, []byte) {

	req, err := http.NewRequest("GET", urlpath, nil)
	for kk, vv := range headers {
		req.Header.Set(kk, vv)
	}

	q := req.URL.Query()
	for pk, pv := range para {
		q.Set(pk, pv)
	}
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	// fmt.Printf("Request : %v \n", req)
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return -1, nil
		// panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Printf("Response : %v \n", string(body))
	return resp.StatusCode, body
}

func MakeSimpleGetToByte(urlpath string) (int, []byte) {
	resp, err := http.Get(urlpath)
	if err != nil {
		fmt.Println("http get error:", err)
		return -1, nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("http io error:", err)
		return -1, nil
		// panic(err)
	}
	// fmt.Println(string(body))
	return resp.StatusCode, body
}

func MakeSimpleGetToStr(urlpath string) (int, string) {
	resp, err := http.Get(urlpath)
	if err != nil {
		return -1, ""
		// handle error
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String() // Does a complete copy of the bytes in the buffer.
	return resp.StatusCode, s
}

func httpPostForm(urlpath string) {
	resp, err := http.PostForm(urlpath,
		url.Values{"key": {"Value"}, "id": {"123"}})
	// resp, err := http.PostForm(urlpath,
	// 	url.Values{"key": {"Value"}, "id": {"123"}})

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
