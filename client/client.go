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
func SetupClient(details Proxy) {
	cookieJar, _ := cookiejar.New(nil)

	//Create client
	Client = newClient(details)

	//Create and set cookiejar in client
	Client.Jar = cookieJar

}

//private
func newClient(details Proxy) http.Client {
	proxyUrl, err := url.Parse("http://localhost:8888")
	if err != nil {
		log.Fatal("Failed + " + err.Error())
	}

	return http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

	//return client
	//return http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	return http.Client{
		/*Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				User:   url.UserPassword(details.ProxyUser, details.ProxyPass),
				Host:   details.Proxy + ":" + details.Port,
			}),
		},
		*/
	}
	//proxies.guap.io:31112
	//user_3109_dover_61
	//bBMr9T3U19SAsoTc_country-UnitedStates_session-DOVqvmvP

	// http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
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

	case PUT:
		if v.Payload == nil {
			req, err := http.NewRequest("PUT", v.Endpoint, nil)
			if err != nil {
				log.Fatal("Failed + " + err.Error())
			}
			return req
		}
		req, err := http.NewRequest("PUT", v.Endpoint, v.Payload)
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
