package kkaws

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"go_test/inerfun"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*
AWSREQUEST aws sign const
*/
const AWSREQUEST = "aws4_request"

/*
AWS4HMACSHA256 AWS sign  algorithm
*/
const AWS4HMACSHA256 = "AWS4-HMAC-SHA256"
const AMZ_SHA256HEADER = "x-amz-content-sha256"
const UNSIGNED_PAYLOAD = "UNSIGNED-PAYLOAD"
const AMZ_DATE = "x-amz-date"

/*
AWSAuthRequ for aws authorization use
*/
type AWSAuthRequ struct {
	Region              string
	Service             string
	HTTPMethod          string
	CanURIBody          string
	CanQueryStr         string
	SignHeader          map[string]string
	SignHeaderStr       string
	CanonicalHeadersStr string
	Unsignpyd           bool
	Payload             []byte
	Createtime          time.Time
}

func getScope(t time.Time, region, service string) string {
	scope := t.Format("20060102") + "/" + region + "/" + service + "/" + AWSREQUEST
	return scope
}

func getCanonicalURI(canURL string) string {
	if len(canURL) == 0 {
		canURL = "/"
	}
	return canURL
}
func getCanonicalQueryString(canQueryStr string) string {
	if len(canQueryStr) != 0 {
		str := url.QueryEscape(canQueryStr)
		return str
	}
	return ""
}

func getHexEncodeSHA256(payloadbyte []byte) string {
	h := sha256.New()
	h.Write(payloadbyte)
	encodedStr := hex.EncodeToString(h.Sum(nil))
	return encodedStr
}

func getHMACSHA256(key []byte, message string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)
	return expectedMAC
}

/*
CreatAWSAuthRequ initial AWS Authorization request
*/
func (req *AWSAuthRequ) CreatAWSAuthRequ(
	HTTPMethod, CanURIBody, CanQueryStr, Region, Service string,
	UnPayload bool, Payload []byte,
	SignHeader map[string]string) {
	req.HTTPMethod = HTTPMethod
	req.CanURIBody = CanURIBody
	req.CanQueryStr = CanQueryStr
	req.SignHeader = SignHeader
	req.Payload = Payload
	req.Unsignpyd = UnPayload
	req.Region = Region
	req.Service = Service
	// layout := "20060102T150405Z"
	// str := "20180119T110113Z"
	// t, err := time.Parse(layout, str)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// req.Createtime = t
	// fmt.Printf("x-amz-date:%v\n", req.Createtime.Format("20060102T150405Z"))
	req.Createtime = time.Now().UTC()
	// req.SignHeader["x-amz-date"] = req.Createtime.Format("20060102T150405Z")

}

/*
The x-amz-content-sha256 header is required for all AWS Signature Version 4 requests.
It provides a hash of the request payload.
If there is no payload, you must provide the hash of an empty string.
*/
func (req *AWSAuthRequ) getCanonicalHeaders() {
	rhdr := ""
	shdr := ""
	// it seems that the order of header is dependent
	if vv, ok := req.SignHeader["Content-Length"]; ok {
		rhdr += strings.ToLower("Content-Length") + ":" + strings.TrimSpace(vv) + "\n"
		shdr += ";" + strings.ToLower("Content-Length")
	}
	if vv, ok := req.SignHeader["Host"]; ok {
		rhdr += strings.ToLower("Host") + ":" + strings.TrimSpace(vv) + "\n"
		shdr += ";" + strings.ToLower("Host")
	}
	if vv, ok := req.SignHeader[AMZ_SHA256HEADER]; ok {
		rhdr += strings.ToLower(AMZ_SHA256HEADER) + ":" + strings.TrimSpace(vv) + "\n"
		shdr += ";" + strings.ToLower(AMZ_SHA256HEADER)
	}
	if vv, ok := req.SignHeader[AMZ_DATE]; ok {
		rhdr += strings.ToLower(AMZ_DATE) + ":" + strings.TrimSpace(vv) + "\n"
		shdr += ";" + strings.ToLower(AMZ_DATE)
	}

	// for kk, vv := range req.SignHeader {
	// 	rhdr += strings.ToLower(kk) + ":" + strings.TrimSpace(vv) + "\n"
	// 	shdr += ";" + strings.ToLower(kk)
	// }
	req.CanonicalHeadersStr = rhdr
	req.SignHeaderStr = strings.TrimPrefix(shdr, ";")
}

func (req *AWSAuthRequ) setSighHeaders(key string, vv string) {
	req.SignHeader[key] = vv
}

/*
CanonicalRequest =
  HTTPRequestMethod + '\n' +
  CanonicalURI + '\n' +
  CanonicalQueryString + '\n' +
  CanonicalHeaders + '\n' +
  SignedHeaders + '\n' +
  HexEncode(Hash(RequestPayload))
*/
func (req *AWSAuthRequ) CanonicalRequest() string {
	rcanURI := getCanonicalURI(req.CanURIBody)
	enqrystr := getCanonicalQueryString(req.CanQueryStr)
	hashpyd := UNSIGNED_PAYLOAD
	if !req.Unsignpyd {
		hashpyd = getHexEncodeSHA256(req.Payload)
	}
	req.SignHeader[AMZ_SHA256HEADER] = hashpyd
	req.SignHeader[AMZ_DATE] = req.Createtime.Format("20060102T150405Z")
	req.getCanonicalHeaders()
	canReq := req.HTTPMethod + "\n" +
		rcanURI + "\n" +
		enqrystr + "\n" +
		req.CanonicalHeadersStr + "\n" +
		req.SignHeaderStr + "\n" +
		hashpyd
	// fmt.Printf("Canonrequest :\n%v\n", canReq)
	return canReq
}

/*
StringToSign =
    Algorithm + \n +
    RequestDateTime + \n +
    CredentialScope + \n +
    HashedCanonicalRequest
*/
func (req *AWSAuthRequ) StringToSign() string {
	credScope := getScope(req.Createtime, req.Region, req.Service)
	hashRequ := getHexEncodeSHA256([]byte(req.CanonicalRequest()))
	signstr := AWS4HMACSHA256 + "\n" +
		req.Createtime.Format("20060102T150405Z") + "\n" +
		credScope + "\n" + hashRequ
	// fmt.Printf("StringToSign :\n%v\n", signstr)
	return signstr

}

/*
CalculateSignature = HexEncode(HMAC(derived signing key, string to sign))
kSecret = your secret access key
kDate = HMAC("AWS4" + kSecret, Date)
kRegion = HMAC(kDate, Region)
kService = HMAC(kRegion, Service)
kSigning = HMAC(kService, "aws4_request")
*/
func (req *AWSAuthRequ) CalculateSignature(secretkey string) string {
	const AWS4 = "AWS4"
	secret := []byte(AWS4 + secretkey)
	date := getHMACSHA256(secret, req.Createtime.Format("20060102"))
	region := getHMACSHA256(date, req.Region)
	service := getHMACSHA256(region, req.Service)
	signingkey := getHMACSHA256(service, AWSREQUEST)
	signingstr := hex.EncodeToString(getHMACSHA256(signingkey, req.StringToSign()))
	return signingstr
}

/*
GetAuthorizationHeader the authorization header
Authorization: algorithm Credential=access key ID/credential scope, SignedHeaders=SignedHeaders, Signature=signature
*/
func (req *AWSAuthRequ) GetAuthorizationHeader(accesskey, secretkey string) string {
	credentialScope := getScope(req.Createtime, req.Region, req.Service)
	signature := req.CalculateSignature(secretkey)
	authstr := fmt.Sprintf("%v Credential=%v/%v, SignedHeaders=%v, Signature=%v", AWS4HMACSHA256, accesskey, credentialScope, req.SignHeaderStr, signature)
	// fmt.Println(authstr)
	return authstr
}

/*

	HTTPMethod  string
	CanURIBody  string
	CanQueryStr string
	SignHeader  map[string]string
	Payload     string
	Region      string
	Service     string
*/
func GetS3Object(accesskey, secretkey, urlpath, Region, filename string, header map[string]string) (int, []byte) {
	var req = &AWSAuthRequ{}
	req.CreatAWSAuthRequ("GET", filename, "", Region, "s3", false, []byte(""), header)
	auth := req.GetAuthorizationHeader(accesskey, secretkey)
	header["Authorization"] = auth
	var para map[string]string
	return inerfun.MakeGet(urlpath+filename, header, para)
}

func PutS3Object(accesskey, secretkey, urlpath, Region, s3keyname, localfilename string, header map[string]string) (int, []byte) {
	content, err := ioutil.ReadFile(localfilename)
	if err != nil {
		log.Fatal(err)
	}
	// ll := strconv.FormatInt(len(content), 10)
	ll := strconv.Itoa(len(content))
	var req = &AWSAuthRequ{}
	header["Content-Length"] = ll
	req.CreatAWSAuthRequ("PUT", s3keyname, "", Region, "s3", false, content, header)
	// req.setSighHeaders("x-amz-acl", "public-read")
	auth := req.GetAuthorizationHeader(accesskey, secretkey)
	header["Authorization"] = auth
	return inerfun.MakePutFile(urlpath+s3keyname, header, bytes.NewBuffer(content))
}

func getHMACSHA1(key []byte, message string) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)
	return expectedMAC
}

func getRESTSignString(secretkey, bucket, filename string, timestamp int64) string {
	rawstr := fmt.Sprintf("GET\n\n\n%v\n/%v/%v", timestamp, bucket, filename)
	b64str := base64.StdEncoding.EncodeToString(getHMACSHA1([]byte(secretkey), rawstr))
	enurl := url.QueryEscape(b64str)
	return enurl

}

/*
GetS3SignedURL get signed url
*/
func GetS3SignedURL(accesskey, secretkey, bucket, filename string) string {
	timestamp := time.Now().UTC().Unix()
	timestamp += 3600
	signature := getRESTSignString(secretkey, bucket, filename, timestamp)
	fullurl := fmt.Sprintf("http://%v.s3.amazonaws.com/%v?AWSAccessKeyId=%v&Expires=%v&Signature=%v", bucket, filename, accesskey, timestamp, signature)
	return fullurl
}
