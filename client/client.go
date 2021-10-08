package request

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var Client http.Client

//Add proxy support later if needed
func SetupClient() {
	cookieJar, _ := cookiejar.New(nil)

	//Create client
	Client = newClient()

	//Create and set cookiejar in client
	Client.Jar = cookieJar

}

//private
func newClient() http.Client {
	//set loopback for charles
	proxyUrl, err := url.Parse("http://localhost:8888")
	if err != nil {
		log.Fatal("Failed + " + err.Error())
	}

	//return client
	return http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
}

func NewRequest(requestType interface{}) *http.Request {
	switch v := requestType.(type) {
	case POST:
		req, err := http.NewRequest("POST", v.Endpoint, v.Payload)
		if err != nil {
			log.Fatal("Failed + " + err.Error())
		}
		return req

	case POSTUrlEncoded:
		req, err := http.NewRequest("POST", v.Endpoint, strings.NewReader(v.EncodedPayload))
		if err != nil {
			log.Fatal("Failed + " + err.Error())
		}
		return req

	case GET:
		req, err := http.NewRequest("GET", v.Endpoint, nil)
		if err != nil {
			log.Fatal("Failed + " + err.Error())
		}
		return req

	default:
		log.Fatal("Request type was invalid")
		return nil
	}
}

func NewResponse(request *http.Request) ([]byte, *http.Response) {
	request.Close = true

	resp, err := Client.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body, resp
}
