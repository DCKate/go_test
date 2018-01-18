package kkaws

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
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

/*
AWSAuthRequ for aws authorization use
*/
type AWSAuthRequ struct {
	HTTPMethod  string
	CanURIBody  string
	CanQueryStr string
	SignHeader  map[string]string
	Payload     string
	Region      string
	Service     string
	Createtime  time.Time
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
	req.Createtime = time.Now()

}

func getCanonicalURI(canURL string) string {
	if len(canURL) == 0 {
		canURL = "/"
	}
	return canURL
}
func getCanonicalQueryString(canQueryStr string) string {
	if len(canQueryStr) != 0 {
		return url.QueryEscape(canQueryStr)
	}
	return ""
}

/*
The x-amz-content-sha256 header is required for all AWS Signature Version 4 requests.
It provides a hash of the request payload.
If there is no payload, you must provide the hash of an empty string.
*/
func getCanonicalHeaders(signHeader map[string]string) string {
	rhdr := ""
	for kk, vv := range signHeader {
		rhdr += strings.ToLower(kk) + ":" + strings.TrimSpace(vv) + "\n"
	}
	return rhdr
}
func getSignedHeaders(signHeader map[string]string) string {
	shdr := ""
	for kk := range signHeader {
		shdr += ";" + strings.ToLower(kk)
	}

	return strings.TrimPrefix(shdr, ";")
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
func CanonicalRequest(requ AWSAuthRequ) string {
	const SHA256HEADER = "x-amz-content-sha256"
	rcanURI := getCanonicalURI(requ.CanURIBody)
	enqrystr := getCanonicalQueryString(requ.CanQueryStr)
	hashpyd := getHexEncodeSHA256(requ.Payload)
	requ.SignHeader[SHA256HEADER] = hashpyd
	rcanhdr := getCanonicalHeaders(requ.SignHeader)
	shdr := getSignedHeaders(requ.SignHeader)

	canReq := requ.HTTPMethod + "\n" +
		rcanURI + "\n" +
		enqrystr + "\n" +
		rcanhdr + "\n" +
		shdr + "\n" +
		hashpyd
	log.Println(canReq)

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
func StringToSign(requ AWSAuthRequ) string {
	// t := time.Now()
	// fmt.Println(t.Format("20060102T150405Z"))

	credScope := getScope(requ.Createtime, requ.Region, requ.Service)
	hashRequ := CanonicalRequest(requ)
	signstr := AWS4HMACSHA256 + "\n" +
		requ.Createtime.Format("20060102T150405Z") + "\n" +
		credScope + "\n" + hashRequ
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
func CalculateSignature(secretkey string, requ AWSAuthRequ) string {
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
func GetAuthorizationHeader(accesskey, secretkey string, requ AWSAuthRequ) string {
	headers := getSignedHeaders(requ.SignHeader)
	credentialScope := getScope(requ.Createtime, requ.Region, requ.Service)
	signature := CalculateSignature(secretkey, requ)
	authstr := fmt.Sprintf("%v Credential=%v/%v, SignedHeaders=%v, Signature=%v", AWS4HMACSHA256, accesskey, credentialScope, headers, signature)
	return authstr
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
