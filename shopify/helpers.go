package shopify

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func AddHeadersTest(header Header, host string) http.Header {
	var x = http.Header{
		"origin":                    {host},
		"sec-ch-ua":                 {"\"Chromium\";v=\"92\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"92\""},
		"sec-ch-ua-mobile":          {"?0"},
		"Upgrade-Insecure-Requests": {"1"},
		"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36"},
		"Sec-Fetch-Site":            {"same-origin"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-User":            {"?1"},
		"Sec-Fetch-Dest":            {"document"},
		"sec-ch-ua-platform":        {"Windows"},
		"accept-language":           {"en-GB,en-US;q=0.9,en;q=0.8"},
		"cache-control":             {"max-age=0"},
		"referer":                   {"https://feature.com/"},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
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

func PaypalCheckout() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.google.com/`),
		chromedp.Sleep(5000),
	); err != nil {

	}

}

func StartPaypalPayment(ppUrl string) {
	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("window-size", "250,600"),
		chromedp.UserDataDir(dir),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// create a timeout
	//taskCtx, cancel = context.WithTimeout(taskCtx, 10*time.Second)
	//defer cancel()

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}

	// listen network event
	listenForNetworkEvent(taskCtx, cancel)

	chromedp.Run(taskCtx,
		network.Enable(),
		chromedp.Navigate(ppUrl),
		chromedp.WaitVisible(`poop`, chromedp.BySearch),
	)

}

func listenForNetworkEvent(ctx context.Context, cancel context.CancelFunc) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {

		case *network.EventResponseReceived:
			resp := ev.Response
			if strings.Contains(resp.URL, "ThePaypalConfirmUrl") {
				log.Printf("Confirmed paypal payment has been made")
				cancel()
			} else {
				fmt.Println(resp.URL)
			}

		}
		// other needed network Event
	})
}
