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
	"strings"
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
	tasks.SetShipping()
	tasks.SetBilling()
	tasks.GetOriginKey()
	tasks.AydenEncrypt1()
	tasks.AydenEncrypt2()
	tasks.EncryptCardWithAdyen()
	tasks.Order()
	//tasks.PutDeliveryMode()
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
	//8. putDeliveryMode -> I dont think we need this
	//9. getO	riginKey
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
		Endpoint: fmt.Sprintf("%s/api/v5/session?timestamp=%d", t.Host, dateTimeStamp),
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
		Payload:  nil,
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

func (t *Task) SetShipping() bool {
	dateTimeStamp := time.Now().UTC().UnixNano()

	payloadBytes, _ := json.Marshal(SetShipping{
		ShippingAddress: ShippingAddress{
			SetAsDefaultBilling:  false,
			SetAsDefaultShipping: false,
			FirstName:            profileInfo.Shipping.FirstNameShipping,
			LastName:             profileInfo.Shipping.LastNameShipping,
			Email:                profileInfo.Person.Email,
			Phone:                profileInfo.Person.Phone,
			Country: Country{
				IsoCodeCountryShipping: profileInfo.Shipping.IsoCodeCountryShipping,
				NameCountryShipping:    profileInfo.Shipping.NameCountryShipping,
			},
			Id:                "",
			SetAsBilling:      true,
			SaveInAddressBook: true,
			Type:              "default",
			LoqateSearch:      "",
			Line1:             profileInfo.Shipping.Line1Shipping,
			Line2:             profileInfo.Shipping.Line2Shipping,
			PostalCode:        profileInfo.Shipping.PostalCodeShipping,
			Town:              profileInfo.Shipping.TownShipping,
			ShippingAddress:   true,
		},
	})

	post := client.POST{
		Endpoint: fmt.Sprintf("%s/api/users/carts/current/addresses/shipping?timestamp=%d", t.Host, dateTimeStamp),
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
	case 201:
		//Read the 'decision' to see if its correct
		return true
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, t.CurrentTaskTemplate.name)
	}

	return false
}

func (t *Task) SetBilling() bool {
	dateTimeStamp := time.Now().UTC().UnixNano()

	payloadBytes, _ := json.Marshal(BillingAddress{
		SetAsDefaultBilling:  false,
		SetAsDefaultShipping: false,
		FirstName:            profileInfo.Billing.FirstNamebilling,
		LastName:             profileInfo.Billing.LastNamebilling,
		Email:                profileInfo.Person.Email,
		Phone:                profileInfo.Person.Phone,
		Country: Country{
			IsoCodeCountryShipping: profileInfo.Billing.IsoCodeCountryBilling,
			NameCountryShipping:    profileInfo.Billing.Namecountrybilling,
		},
		Id:                "",
		SetAsBilling:      true,
		SaveInAddressBook: true,
		Type:              "default",
		LoqateSearch:      "",
		Line1:             profileInfo.Billing.Line1billing,
		Line2:             profileInfo.Billing.Line2billing,
		PostalCode:        profileInfo.Billing.Postalcodebilling,
		Town:              profileInfo.Billing.Townbilling,
		ShippingAddress:   true,
	})

	post := client.POST{
		Endpoint: fmt.Sprintf("%s/api/users/carts/current/set-billing?timestamp=%d", t.Host, dateTimeStamp),
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

//this may not be needed
func (t *Task) PutDeliveryMode() bool {
	dateTimeStamp := time.Now().UTC().UnixNano()

	payloadBytes, _ := json.Marshal(DeliveryMode{
		DeliveryModeId: "fl-standard",
	})

	put := client.PUT{
		Endpoint: fmt.Sprintf("%s/api/users/carts/current/deliverymode?timestamp=%d", t.Host, dateTimeStamp),
		Payload:  bytes.NewReader(payloadBytes),
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

func (t *Task) GetOriginKey() bool {
	dateTimeStamp := time.Now().UTC().UnixNano()

	get := client.GET{
		Endpoint: fmt.Sprintf("%s/apigate/payment/origin-key?timestamp=%d", t.Host, dateTimeStamp),
	}

	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil, contentType: "json"}, "localhost")
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		for _, c := range request.Cookies() {
			if strings.Contains(c.Name, "cart-guid") {
				t.CartId = c.Value
			}
		}

		OriginResponse := OriginResponse{}
		json.Unmarshal(respBytes, &OriginResponse)
		t.OKey = OriginResponse.OKey
		return true
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, resp.Request.URL)
	}

	return false
}

func (t *Task) AydenEncrypt1() bool {
	get := client.GET{
		Endpoint: fmt.Sprintf("https://checkoutshopper-live.adyen.com/checkoutshopper/assets/js/%s/securedFields.1.5.5.min.js", t.OKey),
	}

	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil, contentType: "json"}, "localhost")
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		valueCaptcha := `genTime = [^:]*:([^"]*)`
		re := regexp.MustCompile(valueCaptcha)
		result := re.FindStringSubmatch(string(respBytes))
		t.GenerateTimeMain = result[1][:len(result[1])]
		return true
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, resp.Request.URL)
	}

	return false
}

func (t *Task) AydenEncrypt2() bool {
	var midNumb = strings.Split(t.OKey, ".")[2]
	get := client.GET{
		Endpoint: fmt.Sprintf("https://live.adyen.com/hpp/cse/js/%s.shtml", midNumb),
	}

	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil, contentType: "json"}, "localhost")
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:

		valueCaptcha := `var key = "([^"]*)`
		re := regexp.MustCompile(valueCaptcha)
		result := re.FindStringSubmatch(string(respBytes))
		t.EncryptionKey = result[1][:len(result[1])-1]
		return true
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, resp.Request.URL)
	}

	return false
}

func (t *Task) EncryptCardWithAdyen() bool {
	AydenEncrpytLocal := AydenEncrpytLocal{
		EncryptionKey: t.EncryptionKey,
		DfValueMain:   sessionInfo.DfValue,
		CreditCardNUmber: CreditCardNUmber{
			Activate:            "1",
			Generationtime:      t.GenerateTimeMain,
			InitializeCount:     "1",
			LuhnCount:           "1",
			LuhnOkCount:         "1",
			LuhnSameLengthCount: "1",
			Number:              profileInfo.CardDetails.Ccnumber,
		},
		MonthExpiry: MonthExpiry{
			Activate:        "1",
			ExpiryMonth:     profileInfo.CardDetails.Expiry,
			Generationtime:  t.GenerateTimeMain,
			InitializeCount: "1",
		},
		YearExpiry: YearExpiry{
			Activate:        "1",
			ExpiryYear:      profileInfo.CardDetails.ExpiryYear,
			Generationtime:  t.GenerateTimeMain,
			InitializeCount: "1",
		},
		CvcNumb: CvcNumb{
			Activate:        "1",
			Cvc:             profileInfo.CardDetails.Cvc,
			Generationtime:  t.GenerateTimeMain,
			InitializeCount: "1",
		},
	}

	//Post to local node JS
	payloadBytes, _ := json.Marshal(AydenEncrpytLocal)

	post := client.POST{
		Endpoint: "http://localhost:3000/api/encryptCard",
		Payload:  bytes.NewReader(payloadBytes),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{cookie: []string{}, content: bytes.NewReader(payloadBytes), contentType: "json"}, "localhost")
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		EncryptedCard := EncryptedCard{}
		json.Unmarshal(respBytes, &EncryptedCard)
		t.EncryptedCard = EncryptedCard

		fmt.Println(string(respBytes))
		return true
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, t.CurrentTaskTemplate.name)
	}

	return false

}

func (t *Task) Order() bool {
	dateTimeStamp := time.Now().UTC().UnixNano()

	payloadBytes, _ := json.Marshal(Order{
		OptIn:                 false,
		PreferredLanguage:     "en",
		TermsAndCondition:     true,
		DeviceId:              t.GenDeviceId,
		CartId:                t.CartId,
		EncryptedCardNumber:   t.EncryptedCard.Card,
		EncryptedExpiryMonth:  t.EncryptedCard.Month,
		EncryptedExpiryYear:   t.EncryptedCard.Year,
		EncryptedSecurityCode: t.EncryptedCard.Cvc,
		PaymentMethod:         "CREDITCARD",
		ReturnUrl:             fmt.Sprintf("%s/adyen/checkout", t.Host),
		BrowserInfo: BrowserInfo{
			ScreenWidth:    sessionInfo.Width,
			ScreenHeight:   sessionInfo.Height,
			ColorDepth:     sessionInfo.StoreColorDepth,
			UserAgent:      sessionInfo.UserAgent,
			TimeZoneOffset: sessionInfo.GetOffset,
			Language:       "en-US",
			JavaEnabled:    false,
		},
	})

	post := client.POST{
		Endpoint: fmt.Sprintf("%s/api/users/orders/adyen?timestamp=%d", t.Host, dateTimeStamp),
		Payload:  bytes.NewReader(payloadBytes),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{additionalHeaders: []additionalHeaders{{key: "x-csrf-token", value: t.CsrfToken}}, cookie: []string{}, content: bytes.NewReader(payloadBytes), contentType: "json"}, "localhost")
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
