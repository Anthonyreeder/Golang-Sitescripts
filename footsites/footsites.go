package footsites

import (
	client "Golang-Sitescripts/client"
	"Golang-Sitescripts/profiles"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

var profileInfo = profiles.FootSitesProfile{}
var sessionInfo = profiles.FootSitesSessionInfo{}
var offerId = "SIZE_1235010" //AKA Variant ID
var datadomeCookie = "M7sTRtf1PAbvORz0hOCT.r-75jwAccJiuL5iHbb7ImffTFFY0-54snnMvCSLIIrI7D.tJOWOhVJU_qT1KR65fStUwUMydKgHG.XNb8AElo"

func Footsites() {
	client.SetupClient()

	//Set Profile
	sessionInfo = profiles.FootsitesSessionInfoUk()
	profileInfo = profiles.FootSitesProfileUk()
	//Setup
	client.SetupClient()
	tasks := Task{Host: "https://www.footlocker.co.uk", Link: "https://www.footlocker.co.uk"}

	tasks.GetSnare()
	tasks.EncryptSessionId()
	tasks.GetCSRF()
	tasks.GetCart()
	tasks.VerificationAddress()
	tasks.PutMailFromVerificationAddress()
	//
	//Eastbay
	//Champsports
	//Footlocker US
	//Footlocker Kids
	//Footaction
	//https://www.footlocker.co.uk/
	//https://www.footlocker.co.uk/

	//1. getSnare
	//2. getCSRF
	//3. getChart
	//4. verificationAddress
	//5. putMailFromVerificationAddress
	//6. setShipping
	//7. setBilling
	//8. putDeliveryMode
	//9. getOriginKey
	//10. getAdyenSecured1
	//11. getAdyenSecured2
	//12. order
	//

	//tasks.TaskTemplates = append(tasks.TaskTemplates, TaskTemplate{functionToRun: tasks.ShopifyGetProductPageB, name: "1"})
	//tasks.TaskTemplates = append(tasks.TaskTemplates, TaskTemplate{functionToRun: tasks.ShopifyGetProductPageB, name: "2"})
	//tasks.TaskTemplates = append(tasks.TaskTemplates, TaskTemplate{functionToRun: tasks.ShopifyGetProductPageB, name: "3"})
	//tasks.TaskTemplates = append(tasks.TaskTemplates, TaskTemplate{functionToRun: tasks.ShopifyGetProductPageB, name: "4"})
	//tasks.TaskTemplates = append(tasks.TaskTemplates, TaskTemplate{functionToRun: tasks.GetSnare, name: "1"})

	//runTasks(tasks)

}

func MakeRequest(request Request) ([]byte, *http.Response) {
	//First we create an empty payload
	//This can be a GET, Post or posturlencoded
	var payload interface{}

	switch v := request.Request.(type) {
	case client.POST:
		payload = client.POST{
			Endpoint: v.Endpoint,
			Payload:  v.Payload,
		}

	case client.POSTUrlEncoded:
		payload = client.POSTUrlEncoded{
			Endpoint:       v.Endpoint,
			EncodedPayload: v.EncodedPayload,
		}

	case client.GET:
		payload = client.GET{
			Endpoint: v.Endpoint,
		}

	default:
		log.Fatal("Request type was invalid")
		//fatal error
	}

	//Setup the request with paylaod
	requestToMake := client.NewRequest(payload)
	requestToMake.Header = AddHeaders(Header{cookie: []string{}, content: nil}, request.Host)
	respBytes, resp := client.NewResponse(requestToMake)

	return respBytes, resp
}

func (t *Task) GetSnare() bool {
	get := client.GET{
		Endpoint: "https://mpsnare.iesnare.com/snare.js",
	}

	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil, contentType: "json"}, "mpsnare.iesnare.com")
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:

		s := string(respBytes)
		valueCaptcha := `\"JSSRC\"\,_i_o.__if_ap\("([^"]*).*?JSTOKEN","([^"]*).*?var _i_fd=decodeURIComponent\("([^"]*).*?HACCLNG",decodeURIComponent\("([^"]*).*?SVRTIME","([^"]*).*?IGGY","([^"]*)`

		re := regexp.MustCompile(valueCaptcha)
		result := re.FindStringSubmatch(s)

		sessionInfo := SessionInfo{
			OutJSSRC:        result[1],
			JSTOKEN:         result[2],
			AgentEncoded:    result[3],
			LanguageEncoded: result[4],
			SVRTIME:         result[5],
			IGGY:            result[6],
			RequestAgent:    sessionInfo.RequestAgent,
			AexOffset:       sessionInfo.AexOffset,
			Browser:         sessionInfo.Browser,
			Version:         sessionInfo.Version,
			OsName:          sessionInfo.OsName,
			Appname:         sessionInfo.Appname,
			AppPlatform:     sessionInfo.AppPlatform,
			Height:          sessionInfo.Height,
			Width:           sessionInfo.Width,
			AllPlugins:      sessionInfo.AllPlugins,
			Referer:         sessionInfo.Referer,
			IntLoc:          sessionInfo.IntLoc,
			GetOffset:       sessionInfo.GetOffset,
		}

		t.SessionInfo = sessionInfo
		fmt.Println(result)

		return true
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, resp.Request.URL)
	}

	return false
}

func (t *Task) EncryptSessionId() bool {
	payloadBytes, _ := json.Marshal(t.SessionInfo)

	post := client.POST{
		Endpoint: "http://localhost:3000/api/sessionIdGen",
		Payload:  bytes.NewReader(payloadBytes),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{cookie: []string{}, content: bytes.NewReader(payloadBytes), contentType: "json"}, "localhost")
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		fmt.Println(string(respBytes))
		t.GenDeviceId = string(respBytes)
		return true
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, t.CurrentTaskTemplate.name)
	}

	return false
}

func (t *Task) GetCSRF() bool {
	dateTimeStamp := time.Now().UTC().UnixNano()

	get := client.GET{
		Endpoint: fmt.Sprintf("%s/api/v3/session?timestamp=%d", t.Host, dateTimeStamp),
	}

	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil, contentType: "json"}, "mpsnare.iesnare.com")
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		CsrfResponse := CsrfResponse{}
		json.Unmarshal(respBytes, &CsrfResponse)
		t.CsrfToken = CsrfResponse.Data.CsrfToken
		return true
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, resp.Request.URL)
	}

	return false
}

func (t *Task) GetCart() bool {
	dateTimeStamp := time.Now().UTC().UnixNano()

	payloadBytes, _ := json.Marshal(GetCartRequest{
		ProductQuantity: 1,
		ProductId:       offerId,
	})

	post := client.POST{
		Endpoint: fmt.Sprintf("%s/api/users/carts/current/entries?timestamp=%d", t.Host, dateTimeStamp),
		Payload:  bytes.NewReader(payloadBytes),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{additionalHeaders: []additionalHeaders{
		{key: "x-csrf-token", value: t.CsrfToken},
		{key: "x-fl-productid", value: offerId},
		{key: "origin", value: "https://www.footlocker.co.uk"},
		{key: "referer", value: fmt.Sprintf("%s/en/product/nike-goadome-men-boots/314626041504.html", t.Host)},
		{key: "x-api-lang", value: "en-gb"},
	},
		cookie:  []string{"datadome=3X3o.gD~Y~m6PGoV.-Tu-sT210ZkST.T1Hh6E2FJsio4wkjOJbvrWNFQb.j4zbiPzw0mC8V2UsqxUzZOa8.iLlDAVWDFebpHpgXTVVm4iy"},
		content: bytes.NewReader(payloadBytes), contentType: "json"}, "localhost")
	_, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		return true
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, t.CurrentTaskTemplate.name)
	}

	return false
}

//Read the 'decision' to see if its correct
//Accepted = good
func (t *Task) VerificationAddress() bool {
	dateTimeStamp := time.Now().UTC().UnixNano()

	payloadBytes, _ := json.Marshal(VerificationAddressPayload{
		Country: Country{
			IsoCodeCountryShipping: profileInfo.Shipping.IsoCodeCountryShipping,
			NameCountryShipping:    profileInfo.Shipping.NameCountryShipping,
		},
		Line1Shipping:      profileInfo.Shipping.Line1Shipping,
		Line2Shipping:      profileInfo.Shipping.Line2Shipping,
		PostalCodeShipping: profileInfo.Shipping.PostalCodeShipping,
		TownShipping:       profileInfo.Shipping.TownShipping,
		Region: Region{
			CountryIsoRegionShipping: profileInfo.Shipping.CountryIsoRegionShipping,
			IsoCodeRegionShipping:    profileInfo.Shipping.IsoCodeRegionShipping,
			IsoCodeShortShipping:     profileInfo.Shipping.IsoCodeShortShippingShipping,
			NameRegionShipping:       profileInfo.Shipping.NameRegionShipping,
		},
	})

	post := client.POST{
		Endpoint: fmt.Sprintf("%s/api/v3/users/addresses/verification?timestamp=%d", t.Host, dateTimeStamp),
		Payload:  bytes.NewReader(payloadBytes),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{additionalHeaders: []additionalHeaders{
		{key: "x-csrf-token", value: t.CsrfToken},
		//{key: "x-fl-productid", value: offerId},
		{key: "origin", value: "https://www.footlocker.co.uk"},
		{key: "referer", value: fmt.Sprintf("%s/checkout", t.Host)},
		{key: "x-api-lang", value: "en-gb"},
	},
		cookie:  []string{"datadome=3X3o.gD~Y~m6PGoV.-Tu-sT210ZkST.T1Hh6E2FJsio4wkjOJbvrWNFQb.j4zbiPzw0mC8V2UsqxUzZOa8.iLlDAVWDFebpHpgXTVVm4iy"},
		content: bytes.NewReader(payloadBytes), contentType: "json"}, "localhost")
	_, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		//Read the 'decision' to see if its correct
		return true
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, t.CurrentTaskTemplate.name)
	}

	return false
}

//Dumb name, change this later
//WE are getting 200 but empty body? Not sure if this is correct. Assume its correct for now.
func (t *Task) PutMailFromVerificationAddress() bool {
	dateTimeStamp := time.Now().UTC().UnixNano()

	put := client.PUT{
		Endpoint: fmt.Sprintf("%s/api/users/carts/current/email/%s?timestamp=%d", t.Host, profileInfo.Person.Email, dateTimeStamp),
	}

	request := client.NewRequest(put)
	request.Header = AddHeaders(Header{additionalHeaders: []additionalHeaders{
		{key: "x-csrf-token", value: t.CsrfToken},
		//{key: "x-fl-productid", value: offerId},
		{key: "origin", value: "https://www.footlocker.co.uk"},
		{key: "referer", value: fmt.Sprintf("%s/checkout", t.Host)},
		{key: "x-api-lang", value: "en-gb"},
	},
		cookie:      []string{"datadome=3X3o.gD~Y~m6PGoV.-Tu-sT210ZkST.T1Hh6E2FJsio4wkjOJbvrWNFQb.j4zbiPzw0mC8V2UsqxUzZOa8.iLlDAVWDFebpHpgXTVVm4iy"},
		contentType: "json"}, "localhost")
	_, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		//Read the 'decision' to see if its correct
		return true
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, t.CurrentTaskTemplate.name)
	}

	return false
}
