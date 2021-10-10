package shopify

import (
	client "Golang-Sitescripts/client"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

//GET Product page and extract bot-key
func ShopifyGetProductPage() {
	//Setup our GET request obj
	get := client.GET{
		Endpoint: fmt.Sprintf("%s/collections/mens/products/adidas-originals-pharrell-williams-boost-slides-fy6140", link),
	}
	//Retrieve a configured HTTP Request obj
	request := client.NewRequest(get)
	//Add our headers to the HTTP Request obj
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	//Obtain the response
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:

		//Find the bot-key input field in the form
		botKey = ExtractValue(string(respBytes), "input", "id", "bot-key")

		//Check if botkey has a value now
		if botKey != "" {
			taskComplete = true
			return
		} else {
			fmt.Println("There was an issue getting the bot key")
		}
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, get.Endpoint)
	}

	taskComplete = false
}

//POST JSON data to the standard endpoint used on the browsers 'addToCart' button.
func ShopifyAddToCartStandard() {
	var addToCartId string

	if _offerid != "" {
		addToCartId = _offerid
		_offerid = ""
	} else {
		addToCartId = offerId
	}

	payloadBytes, _ := json.Marshal(AddToCartStandardRequest{
		Id:       addToCartId,
		Quantity: quantity,
	})

	post := client.POST{
		Endpoint: fmt.Sprintf("%s/cart/add.js", link),
		Payload:  bytes.NewReader(payloadBytes),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil, contentType: "json"}, host)
	_, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		taskComplete = true
		return
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, post.Endpoint)
	}

	taskComplete = false
}

//https://shopify.dev/api/admin-rest/2021-10/resources/payment#[post]https://elb.deposit.shopifycs.com/sessions
func CreatePaymentSession() {
	payloadBytes, _ := json.Marshal(PaymentSessionRequest{
		CreditCard: CreditCard{
			Number:             cardNumber,
			Name:               name,
			Month:              month,
			Year:               year,
			Verification_value: ccv,
		},
	})

	post := client.POST{
		Endpoint: "https://elb.deposit.shopifycs.com/sessions",
		Payload:  bytes.NewReader(payloadBytes),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil, contentType: "json"}, "elb.deposit.shopifycs.com")
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:

		paymentSessionResponse := PaymentSessionResponse{}
		json.Unmarshal(respBytes, &paymentSessionResponse)
		payment_token = paymentSessionResponse.Id

		taskComplete = true
		return
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, post.Endpoint)
	}

	taskComplete = false
}

//GET the checkout form page and extract the AuthId
func LoadCheckoutForm() {
	get := client.GET{
		Endpoint: fmt.Sprintf("%s/checkout", link),
	}

	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:

		//Find the auth-key input field in the form
		authKey = ExtractValue(string(respBytes), "input", "name", "authenticity_token")

		//globalise the redirected url
		formUrl = resp.Request.URL.String()

		//Check if authKey has a value now
		if authKey != "" {
			taskComplete = true
			return
		} else {
			fmt.Println("There was an issue getting the auth key")
		}
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, get.Endpoint)
	}

	taskComplete = false
}

//POST the profile information
func SubmitCustomerInfo() {
	payload := url.Values{
		"utf8":                                   {`\u2713`},
		"_method":                                {"patch"},
		"authenticity_token":                     {authKey},
		"previous_step":                          {"contact_information"},
		"step":                                   {"shipping_method"},
		"checkout[email]":                        {email},
		"checkout[buyer_accepts_marketing]":      {"1"},
		"checkout[pickup_in_store][selected]":    {"false"},
		"checkout[shipping_address][first_name]": {fname},
		"checkout[shipping_address][last_name]":  {lname},
		"checkout[shipping_address][company]":    {company},
		"checkout[shipping_address][address1]":   {addy1},
		"checkout[shipping_address][address2]":   {addy2},
		"checkout[shipping_address][city]":       {city},
		"checkout[shipping_address][country]":    {country},
		"checkout[shipping_address][province]":   {province},
		"checkout[shipping_address][zip]":        {postal_code},
		"checkout[shipping_address][phone]":      {phone},
		// "g-recaptcha-response": captcha_token,
		"checkout[client_details][browser_width]":      {"1029"},
		"checkout[client_details][browser_height]":     {"937"},
		"checkout[client_details][javascript_enabled]": {"1"},
		"checkout[client_details][color_depth]":        {"24"},
		"checkout[client_details][java_enabled]":       {"false"},
		"checkout[client_details][browser_tz]":         {"300"},
	}

	post := client.POSTUrlEncoded{
		Endpoint:       formUrl,
		EncodedPayload: payload.Encode(),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	_, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		taskComplete = true
		return
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, post.Endpoint)
	}

	taskComplete = false
}

//GET the shipping rates for this profile and extract the shipping id
//There is an async POST method which may be quicker
func ExtractShippingRates() {
	get := client.GET{
		Endpoint: fmt.Sprintf("%s/cart/shipping_rates.json?shipping_address[zip]=%s&shipping_address[country]=%s&shipping_address[province]=%s", link, postal_code, country, province),
	}

	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:

		//Decode the response into a json struct
		shippingMethodResponse := ShippingMethodResponse{}
		json.Unmarshal(respBytes, &shippingMethodResponse)

		//extract the name and price
		name := strings.Replace(shippingMethodResponse.ShippingRates[0].Name, " ", "%20", -1)
		price := shippingMethodResponse.ShippingRates[0].Price

		//# Generate the shipping id to submit with checkout
		shipping_option = "shopify-" + name + "-" + price
		if shipping_option != "" {
			taskComplete = true
			return
		}

	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, get.Endpoint)
	}

	taskComplete = false
}

//GET the shipping rates for this profile and extract the shipping id
//There is an async POST method which may be quicker
func POSTExtractShippingRates() {
	payload := url.Values{
		"shipping_address[zip]":      {postal_code},
		"shipping_address[country]":  {country},
		"shipping_address[province]": {province},
	}

	post := client.POSTUrlEncoded{
		Endpoint:       fmt.Sprintf("%s/cart/prepare_shipping_rates.json", link),
		EncodedPayload: payload.Encode(),
	}
	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:

		//Decode the response into a json struct
		shippingMethodResponse := ShippingMethodResponse{}
		json.Unmarshal(respBytes, &shippingMethodResponse)

		//extract the name and price
		name := strings.Replace(shippingMethodResponse.ShippingRates[0].Name, " ", "%20", -1)
		price := shippingMethodResponse.ShippingRates[0].Price

		//# Generate the shipping id to submit with checkout
		shipping_option = "shopify-" + name + "-" + price
		if shipping_option != "" {
			taskComplete = true
			return
		}

	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, post.Endpoint)
	}

	taskComplete = false
}

//GET the next-step in shipping to extrac the shipping token
func ExtractShippingToken() {
	authKey = ""

	get := client.GET{
		Endpoint: fmt.Sprintf("%s?step=shipping_method", formUrl),
	}

	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:

		//Find the auth-key input field in the form
		authKey = ExtractValue(string(respBytes), "input", "name", "authenticity_token")

		//Check if authKey has a value now
		if authKey != "" {
			taskComplete = true
			return
		} else {
			fmt.Println("There was an issue getting the auth key")
		}
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, get.Endpoint)
	}

	taskComplete = false
}

//POST the shipping token and shipping ID
func SubmitShippingMethodDetails() {
	payload := url.Values{
		"utf8":                        {`\u2713`},
		"_method":                     {"patch"},
		"authenticity_token":          {authKey},
		"previous_step":               {"contact_information"},
		"step":                        {"payment_method"},
		"checkout[shipping_rate][id]": {shipping_option},
		"checkout[client_details][browser_width]":      {"1029"},
		"checkout[client_details][browser_height]":     {"937"},
		"checkout[client_details][javascript_enabled]": {"1"},
		"checkout[client_details][color_depth]":        {"24"},
		"checkout[client_details][java_enabled]":       {"false"},
		"checkout[client_details][browser_tz]":         {"300"},
	}

	post := client.POSTUrlEncoded{
		Endpoint:       formUrl,
		EncodedPayload: payload.Encode(),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	_, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		taskComplete = true
		return
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, post.Endpoint)
	}

	taskComplete = false
}

//GET the payment_method values needed to submit a payment
func ExtractPaymentGatewayId() {
	authKey = ""

	get := client.GET{
		Endpoint: fmt.Sprintf("%s?previous_step=shipping_method&step=payment_method", formUrl),
	}

	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:

		//Find the gateway, auth and total amount values.
		gatewayKey = ExtractValue(string(respBytes), "input", "name", "checkout[payment_gateway]")
		authKey = ExtractValue(string(respBytes), "input", "name", "authenticity_token")
		total_amount = ExtractValue(string(respBytes), "span", "class", "payment-due__price", "data-checkout-payment-due-target")

		//Check if authKey has a value now
		if authKey != "" && gatewayKey != "" && total_amount != "" {
			taskComplete = true
			return
		} else {
			fmt.Println("There was an issue getting the auth key")
		}
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, get.Endpoint)
	}

	taskComplete = false
}

func SubmitPayment() {
	payload := url.Values{
		"_method":                             {"patch"},
		"authenticity_token":                  {authKey},
		"previous_step":                       {"payment_method"},
		"step":                                {""},
		"s":                                   {payment_token},
		"checkout[payment_gateway]":           {gatewayKey},
		"checkout[credit_card][vault]":        {"false"},
		"checkout[different_billing_address]": {"false"}, //This should be set to true then we use profile billing. Look at later.
		"checkout[vault_phone]":               {""},
		"checkout[total_price]":               {total_amount},
		"complete":                            {"1"},
		"checkout[client_details][browser_width]":      {strconv.Itoa(rand.Intn(2000-1000) + 1000)}, //I dont like this, look at this later.
		"checkout[client_details][browser_height]":     {strconv.Itoa(rand.Intn(2000-1000) + 1000)},
		"checkout[client_details][javascript_enabled]": {"1"},
		"checkout[client_details][color_depth]":        {"24"},
		"checkout[client_details][java_enabled]":       {"false"},
		"checkout[client_details][browser_tz]":         {"300"},
	}

	post := client.POSTUrlEncoded{
		Endpoint:       formUrl,
		EncodedPayload: payload.Encode(),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	_, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		process_url = resp.Request.URL.String()
		taskComplete = true
		return

	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, post.Endpoint)
	}

	taskComplete = false
}

func CheckPaymentProcess() {
	get := client.GET{
		Endpoint: process_url,
	}

	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		if !strings.Contains(resp.Request.URL.String(), "processing") {
			_body := soup.HTMLParse(string(respBytes))
			messageResponse := _body.Find("p", "class", "notice__text").Text()
			fmt.Printf("Checkout response: %s", messageResponse)
			taskComplete = true
			return
		}

	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, get.Endpoint)
	}

	taskComplete = false
}

//More methods
//POST JSON data to the standard endpoint used on the browsers 'addToCart' button.
func ShopifyChangeCart(q int) {
	payloadBytes, _ := json.Marshal(ChangeCartStandardRequest{
		Id:       "32521243820103",
		Quantity: q,
	})

	post := client.POST{
		Endpoint: fmt.Sprintf("%s/cart/change.js", link),
		Payload:  bytes.NewReader(payloadBytes),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil, contentType: "json"}, host)
	_, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		taskComplete = true
		return

	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, post.Endpoint)
	}

	taskComplete = false
}
