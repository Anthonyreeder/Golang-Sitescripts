package shopify

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/anaskhan96/soup"
)

//Default headers with functionality to set the host, content type and add 1-off hard-coded cookies.
func AddHeaders(header Header, host string) http.Header {
	var x = http.Header{
		"Host":                      {host},
		"sec-ch-ua":                 {"\"Chromium\";v=\"92\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"92\""},
		"sec-ch-ua-mobile":          {"?0"},
		"Upgrade-Insecure-Requests": {"1"},
		"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"Sec-Fetch-Site":            {"none"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-User":            {"?1"},
		"Sec-Fetch-Dest":            {"document"},
		"Accept-Language":           {"en-GB,en-US;q=0.9,en;q=0.8"},
	}

	if header.content != nil {
		x.Set("Content-Length", fmt.Sprint(header.content.Size()))
	}

	if header.contentType == "json" {
		x.Set("content-type", "application/json")
	} else {
		x.Set("content-type", "application/x-www-form-urlencoded")
	}

	if len(header.cookie) > 0 {
		buildString := ""
		for i := 0; i < len(header.cookie); i++ {
			buildString += header.cookie[i] + "; "
		}
		x.Set("Cookie", buildString+strings.Join(x.Values("Cookie"), "; "))
	}

	return x
}

//Used in multiple methods to extract key values
func ExtractValue(body, elementType, targetType, typeValue string, optionalAttribute ...string) string {
	var val = ""
	var value = "value"
	if len(optionalAttribute) > 0 {
		value = optionalAttribute[0]
	}
	_body := soup.HTMLParse(body)
	element := _body.Find(elementType, targetType, typeValue).Pointer.Attr
	for _, v := range element {
		if v.Key == value {
			//Locate the authKey attribute value within this node
			val = v.Val
		}
	}

	return val
}
