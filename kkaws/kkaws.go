package kkaws

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"inerfun"
	"log"
	"net/url"
	"os"
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
const SHA256HEADER = "x-amz-content-sha256"

/*
AWSAuthRequ for aws authorization use
*/
type AWSAuthRequ struct {
	HTTPMethod  string
	CanURIBody  string
	CanQueryStr string
	SignHeader  map[string]string
	// SignHeaderStrCurl   string
	SignHeaderStr       string
	CanonicalHeadersStr string
	Payload             string
	Region              string
	Service             string
	Createtime          time.Time
}

/*
CreatAWSAuthRequ initial AWS Authorization request
*/
func (req *AWSAuthRequ) CreatAWSAuthRequ(
	HTTPMethod, CanURIBody, CanQueryStr, Payload, Region, Service string,
	SignHeader map[string]string) {
	req.HTTPMethod = HTTPMethod
	req.CanURIBody = CanURIBody
	req.CanQueryStr = CanQueryStr
	req.SignHeader = SignHeader
	req.Payload = Payload
	req.Region = Region
	req.Service = Service
	layout := "20060102T150405Z"
	str := "20180119T110113Z"
	t, err := time.Parse(layout, str)
	if err != nil {
		fmt.Println(err)
	}
	req.Createtime = t
	// fmt.Printf("x-amz-date:%v\n", req.Createtime.Format("20060102T150405Z"))
	// req.Createtime = time.Now().UTC()
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
	// hhdr := ""
	for kk, vv := range req.SignHeader {
		rhdr += strings.ToLower(kk) + ":" + strings.TrimSpace(vv) + "\n"
		shdr += ";" + strings.ToLower(kk)
		// hhdr += ";" + kk
	}
	req.CanonicalHeadersStr = rhdr
	req.SignHeaderStr = strings.TrimPrefix(shdr, ";")
	// req.SignHeaderStrCurl = strings.TrimPrefix(hhdr, ";") + ";Authorization"
}

func (req *AWSAuthRequ) setSighHeaders(key string, vv string) {
	req.SignHeader[key] = vv
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

func getHexEncodeSHA256(payload string) string {
	h := sha256.New()
	h.Write([]byte(payload))
	encodedStr := hex.EncodeToString(h.Sum(nil))
	return encodedStr
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
func CanonicalRequest(requ *AWSAuthRequ) string {
	rcanURI := getCanonicalURI(requ.CanURIBody)
	enqrystr := getCanonicalQueryString(requ.CanQueryStr)
	hashpyd := getHexEncodeSHA256(requ.Payload)
	requ.setSighHeaders(SHA256HEADER, hashpyd)
	requ.setSighHeaders("x-amz-date", requ.Createtime.Format("20060102T150405Z"))
	requ.getCanonicalHeaders()
	canReq := requ.HTTPMethod + "\n" +
		rcanURI + "\n" +
		enqrystr + "\n" +
		requ.CanonicalHeadersStr + "\n" +
		requ.SignHeaderStr + "\n" +
		hashpyd
	fmt.Printf("Canonrequest :\n%v\n", canReq)
	return canReq
}

func getScope(t time.Time, region, service string) string {
	scope := t.Format("20060102") + "/" + region + "/" + service + "/" + AWSREQUEST
	return scope
}

/*
StringToSign =
    Algorithm + \n +
    RequestDateTime + \n +
    CredentialScope + \n +
    HashedCanonicalRequest
*/
func StringToSign(requ *AWSAuthRequ) string {
	credScope := getScope(requ.Createtime, requ.Region, requ.Service)
	hashRequ := getHexEncodeSHA256(CanonicalRequest(requ))
	signstr := AWS4HMACSHA256 + "\n" +
		requ.Createtime.Format("20060102T150405Z") + "\n" +
		credScope + "\n" + hashRequ
	fmt.Printf("StringToSign :\n%v\n", signstr)
	return signstr

}
func getHMACSHA256(key []byte, message string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)
	return expectedMAC
}

/*
CalculateSignature = HexEncode(HMAC(derived signing key, string to sign))
kSecret = your secret access key
kDate = HMAC("AWS4" + kSecret, Date)
kRegion = HMAC(kDate, Region)
kService = HMAC(kRegion, Service)
kSigning = HMAC(kService, "aws4_request")
*/
func CalculateSignature(secretkey string, requ *AWSAuthRequ) string {
	const AWS4 = "AWS4"
	secret := []byte(AWS4 + secretkey)
	date := getHMACSHA256(secret, requ.Createtime.Format("20060102"))
	region := getHMACSHA256(date, requ.Region)
	service := getHMACSHA256(region, requ.Service)
	signingkey := getHMACSHA256(service, AWSREQUEST)
	signingstr := hex.EncodeToString(getHMACSHA256(signingkey, StringToSign(requ)))
	return signingstr
}

/*
GetAuthorizationHeader the authorization header
Authorization: algorithm Credential=access key ID/credential scope, SignedHeaders=SignedHeaders, Signature=signature
*/
func GetAuthorizationHeader(accesskey, secretkey string, requ *AWSAuthRequ) string {
	credentialScope := getScope(requ.Createtime, requ.Region, requ.Service)
	signature := CalculateSignature(secretkey, requ)
	authstr := fmt.Sprintf("%v Credential=%v/%v, SignedHeaders=%v, Signature=%v", AWS4HMACSHA256, accesskey, credentialScope, requ.SignHeaderStr, signature)
	fmt.Println(authstr)
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
	req.CreatAWSAuthRequ("GET", filename, "", "", Region, "s3", header)
	auth := GetAuthorizationHeader(accesskey, secretkey, req)
	header["Authorization"] = auth
	var para map[string]string
	return inerfun.MakeGet(urlpath+filename, header, para)
}

func PutS3Object(accesskey, secretkey, urlpath, Region, s3keyname, localfilename string, header map[string]string) (int, []byte) {
	file, err := os.Open(localfilename)
	if err != nil {
		log.Fatal(err)
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	ll := strconv.FormatInt(fi.Size(), 10)
	var req = &AWSAuthRequ{}
	req.CreatAWSAuthRequ("PUT", s3keyname, "", "UNSIGNED-PAYLOAD", Region, "s3", header)
	// req.setSighHeaders("x-amz-acl", "public-read")
	req.setSighHeaders("Content-Length", ll)
	auth := GetAuthorizationHeader(accesskey, secretkey, req)
	header[SHA256HEADER] = req.SignHeader[SHA256HEADER]
	// header["Content-Length"] = ll
	header["Authorization"] = auth
	return inerfun.MakePutFile(urlpath+s3keyname, header, file)
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
