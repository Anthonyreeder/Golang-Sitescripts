package shopify

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

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
	} else if header.contentType == "multipart" {
		x.Set("content-type", "multipart/form-data; boundary=----WebKitFormBoundary45pI4iftSbnzXGQ1")

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
	el := _body.Find(elementType, targetType, typeValue)
	if el.Error == nil {
		element := el.Pointer.Attr
		for _, v := range element {
			if v.Key == value {
				//Locate the authKey attribute value within this node
				val = v.Val
			}
		}
	}

	return val
}

//Task helpers, to loop functions and log failures to console
//In future change this so FunctionToRun and Name are in a 'function' object WIll probably build on this in future.
func startTask(functionToRun func(), name string, waitForResult ...bool) {
	taskComplete = false
	fmt.Printf("Running task %s\n", name)
	if len(waitForResult) > 0 {
		go func() {
			loopFunction(functionToRun, name)
		}()
	} else {
		loopFunction(functionToRun, name)
	}
}

func loopFunction(functionToRun func(), name string) {
	for {
		functionToRun()
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(2000)
		time.Sleep(time.Duration(r) * time.Millisecond)
		if taskComplete {
			break
		} else {
			fmt.Printf("%s task failed - retrying\n", name)
		}
	}
}

func startTaskInt(functionToRun func(int), name string, val int) {
	taskComplete = false
	//fmt.Printf("Running task %s\n", name)
	for {
		functionToRun(val)
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(2000)
		time.Sleep(time.Duration(r) * time.Millisecond)
		if taskComplete {
			break
		} else {
			fmt.Printf("%s task failed - retrying\n", name)
		}
	}
}

//Gets the product and checks if its in stock
func GetProductInStock(p Products, sku string) ProductData {
	for _, product := range p {
		for _, variant := range product.Variants {
			if fmt.Sprint(variant.Id) == sku {
				if variant.Available {
					return product
				}
			}
		}
	}
	return ProductData{Title: ""}
}
