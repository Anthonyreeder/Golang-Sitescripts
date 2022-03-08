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
	//	client.SetupClient()

	//Set Profile
	sessionInfo = profiles.FootsitesSessionInfoUk()
	profileInfo = profiles.FootSitesProfileUk()
	//Setup
	//client.SetupClient()
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
		valueCaptcha := `genTime = "([^"]*)`
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
		DfValueMain:   sessionInfo.DfValue,
		EncryptionKey: t.EncryptionKey,

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
	fmt.Println("Test")
	fmt.Println(AydenEncrpytLocal.CreditCardNUmber)
	fmt.Println(AydenEncrpytLocal.MonthExpiry)
	fmt.Println(AydenEncrpytLocal.YearExpiry)
	fmt.Println(AydenEncrpytLocal.CvcNumb)
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
		fmt.Println("length")
		fmt.Println(len(EncryptedCard.Card))
		fmt.Println(len(EncryptedCard.Month))
		fmt.Println(len(EncryptedCard.Year))
		fmt.Println(len(EncryptedCard.Cvc))

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
		PreferredLanguage:     "en",
		TermsAndCondition:     true,
		DeviceId:              t.GenDeviceId,
		CartId:                t.CartId,
		EncryptedCardNumber:   "adyenjs_0_1_18$lVVV2mqti3NJNABloIoTjYhSzBUQIklPgCBuprAKA/pc4oXhC/7WaJk1OhTShsbjqLU05EcPtxCiUmwpxLyYpEU/XhZ8mNdVT5pAKNpD2rW/qsGl/GE9ITRfTHo5QhH/HMcMsNvhoSGSJX4pPUOksItqDcQxiZ7dpojJtGiX4AQbYNvOKxOqJAgfZLCN2J82rIAhKrVNwQg993Y1jePkFihbNthMjH6V1UWGRodoyBPn036SHhpJG89vSMXoAAMNsWce29K6g1y944SVbD0KOMgqZ94axog3Rb2j6ISB4Eo568DSfjbXFcqRWxYrI2yKd8be31WpomfIMulH4lCppg==$LPPlHPbNQbjBqJxE7DT40QFMmVRzrtXLKgnO6j5XOsw6HS/pCKDT1swKNoY2mHBW6SLnW/owsg7FnjrrF4KhPsyyr3rQX/pc0uvZ08tU2sGlIHhktpuRmTzh8witV0WEAM7h3lSSIGcLSO+BWK1jhkn3e+2+g47KPNWGZmJo4MOi6LDVFQeNWyyX4hQjCSl3IsEMowVkI2Z7tO3MUR+CHwK465lfuI09Bnb8oD25S0XuuKJSWBXdopGuEV85UAc0khuUsZAyPA+tInM8E9kO4fKHMnV+DDQUQVgtmgIdAnGuzJY9OrjhKjXVCS4hDKDirlNT7Z/UWa/QeXvqjFFQv574LUSqxmvch+DeZ5i9w+6nWzzWxQnzZ+VfmOZX+uzbfzn2YFDQxk2uD9sgSpYc3sanDdWoH4naIVwk9pmP03kRXAFnE+Ve0Q4=",
		EncryptedExpiryMonth:  "adyenjs_0_1_18$I7n+NtkSi3NvEVCURkPSHYOHGSHUhQumbmQS6K7YmakJN+sp0t0W2590FAVOtayzimL1pYjYq9mtQwqnM/OhoSay/g5edSsdwOLc8udQ6Zta63q58VJFh1w0y49roTllQycpwOx428WhiMSGaJJiXp/W/Mbe6y+l3MXsj3Q4yrW0gwDzsTr3sIm6oTrwYp7sDd6UADBtPykNxqpQCcrLezI3rcXygjh2IvM4WRyGreXZ6EP4xU5nnK84y7lY0WEjMRF4N0h6ov7wG3LEw1OxmcPmqTV8Z67ZjfGoTJ9yTZ8PshM1EFUjuh+CdW9aMeziZfzoxc+E/sVrShpYaKX5zA==$RChDi74GS08UZGImbC3C8P/MwajnDIFNwLicdY88ZInsriaJUCKatJejn8TjTDl1GJUgWPjsrXnuVqQw1LBBke+1lld8A4dpiDJAj2ZBDoe3g+F+gSUfzePDt6rTCU0oNZf3gKHQQO0x13r14zCsSPvpLMFVkvyT+gpU+5J7eHiSuUyRx9wmjhfglfCaDH3r0I7oWEP1BKm808s4m60bVSLh7ysv25mC07UAfePS2P1XKkh64uQx9Ip64GFCpSbKacCwsIbutAPJVyDK0RqsPwJxzHdGkUx+BwyA3aNNj0ZlaCqYG63eS6rgBl+uLDxkuC6fXNk0osP3tfheoh/+r7A=",
		EncryptedExpiryYear:   "adyenjs_0_1_18$XgAzPCRszMcq5+Oar3S6RdJD3+2FrSo1R2/9+lpjMLZwp8MxO8utuJOddePoX2S9eOgOBbaQpYPkrBttxgnwqa1rEYVGIziiZlS9cDmGvsgRE5YRxH1g74UHebUSTk67r/EZTyhIrdqa5naUl+E+hRmAGI4uBZegrqIQUVE76AQLlwMx1M90dtEhqoI41sdOzewRl1QkjK7/sG6doKwK07KJeOJaMRwG88pwjYpAKYinGtAkDH3DytINAe9GxLsbHnVbQnYwt/CSFKUxK429Jq++4lBe78ru+SLvdcD9xxHcS3wDJgbYL2Qokup9N4bvdOeKD767PPFsFsDiIIR25Q==$PTRm0/puTMnP7788Mhj9JOxsf0hgiEQeNl0emQEuKyCePuSkYMQwrSD8OPg8lPnWhjs6PZ0VYC6yGG5nqcvzOvDuaWSvwHrJfOS5YBqyPHwyIjJ5S4TX2Dhz5HTUm/M5+9O0VkrLxzCd8G6ZmnyPo62bhDuJ810WpD52diE5CEGIhUxfdRveN8j1KHQ4hCzsLSqLUKs0DFNsOmP6Vomcmevd/DnClIOZJ725EZICZ4Wt5up4c+4Kl1CNzuZ2J+vplRHvDgZEx9Bsh8zI3kPhoxoE2J3RT2UECmSe1Rj4y0x5TCg/kjJcapX/xtdWBxfjtJ57EwE1mRsAauA472VPnJRr",
		EncryptedSecurityCode: "adyenjs_0_1_18$P7hcyGfSf3Ndhf8GDremTbbSE/hXVFqN/nS+qaeN324lQcQcJHSb23/p7riQXzqzoyRYgw4yoGa6Knp2Uw7Jloarm+id5daCqfBU951DhY+EF8JRlWU77JM1AzRCXAF8Irs3fzGdaja4Ie0Cg2GojPI3FnDQdGawPe3/9iAeohxSl/TGvVIJS0jxHb+ZJxieMmEDQ8EUe0DQqB+lEY5mpEMoEIn27RoBIJM+9pPrumIGGlcFqQH+I0P/NH6qS4f+Xvy3oOVPs3I2Y7yhVbZCWyIKBCS4FvJd9PWXHBtM4RhAaK/zw7fZT+vdiLN/fBZ6SRinkYUY0CQ47/KYRQQRoA==$fj+BgpzYEBvk51CdLChWHpSysGPm4H+0mIbq63s5BnrgsSnsHC7zP6gMmAQPTyNAWxMXk6f1O372RcqZDFU9do2W/nNJPE4i5EWiAR5h0+av7lwnSEz+i+e6b7V97hEYfeutgQMKzWTMMStIg7aKlPcsdn7XZM86uk+XzRmeMdL9rjWB3Lnjk4RxMoMmwzqRfQXJYcPvXjgLXHcqazxEUMuAPR5dbN+X8sma7uNlA5lyJxdHhOgvTjwxrEA8CMnd1p17IzC996I/SbUtJQZs+S3Z81hOECgVAdLuSJ484ARDzcUpjUw7+wTuAegA+a0lHatL6ThvIdUs6w==",
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
		Endpoint: fmt.Sprintf("%s/api/v2/users/orders?timestamp=%d", t.Host, dateTimeStamp),
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
