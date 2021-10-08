package shopify

import (
	client "Golang-Sitescripts/client"
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

//GET Product page and extract bot-key
func ShopifyGetProductPage() bool {

	//Setup our GET request obj
	get := client.GET{
		Endpoint: fmt.Sprintf("%s/collections/mens/products/adidas-originals-pharrell-williams-boost-slides-fy6140", host),
	}
	//Retrieve a configured HTTP Request obj
	request := client.NewRequest(get)
	//Add our headers to the HTTP Request obj
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	//Obtain the response
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		fmt.Println("Page loaded, checking for bot-key")

		//Find the bot-key input field in the form
		botKey = ExtractValue(string(respBytes), "input", "id", "bot-key")

		//Check if botkey has a value now
		if botKey != "" {
			fmt.Println("Successfully extracted the bot-key")
			return true
		} else {
			fmt.Println("There was an issue getting the bot key")
		}
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, get.Endpoint)
	}

	return false
}

//POST JSON data to the standard endpoint used on the browsers 'addToCart' button.
func ShopifyAddToCartStandard() bool {
	payloadBytes, _ := json.Marshal(AddToCartStandardRequest{
		FormType: "product",
		Utf8:     "",
		Properties: struct {
			BotKey string `json:"bot-key"`
		}{BotKey: botKey},
		OptionSize: size,
		Id:         offerId,
		Quantity:   quantity,
	})

	post := client.POST{
		Endpoint: fmt.Sprintf("%s/cart/add.js", host),
		Payload:  bytes.NewReader(payloadBytes),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil, contentType: "json"}, host)
	_, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		fmt.Println("Successfully added item to cart ")
		return true

	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, post.Endpoint)
	}

	return false
}

//GET the checkout form page and extract the AuthId
func LoadCheckoutForm() bool {
	get := client.GET{
		Endpoint: fmt.Sprintf("%s/checkout", host),
	}

	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		fmt.Println("Checkout form loaded, checking for auth token")

		//Find the auth-key input field in the form
		authKey = ExtractValue(string(respBytes), "input", "name", "authenticity_token")

		//globalise the redirected url
		formUrl = resp.Request.URL.String()

		//Check if authKey has a value now
		if authKey != "" {
			fmt.Printf("Successfully extracted the auth-key : %s\n", authKey)
			fmt.Printf("Retrieved the form url : %s\n", formUrl)
			return true
		} else {
			fmt.Println("There was an issue getting the auth key")
		}
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, get.Endpoint)
	}

	return false
}

//POST the profile information
func SubmitCustomerInfo() bool {
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
		// "checkout[shipping_address][province]": province,
		"checkout[shipping_address][zip]":   {postal_code},
		"checkout[shipping_address][phone]": {phone},
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
		fmt.Printf("Successfully posted the customer information for %s\n", email)
		return true

	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, post.Endpoint)
	}

	return false
}

//GET the shipping details for this profile and extract the shipping token
func GetShippingToken() bool {
	get := client.GET{
		Endpoint: fmt.Sprintf("%s/cart/shipping_rates.json?shipping_address[zip]=%s&shipping_address[country]=%s&shipping_address[province]=%s", host, postal_code, country, province),
	}

	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	respBytes, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		//fmt.Println("Shipping token request loaded, extracting shipping token")
		fmt.Println(string(respBytes))

	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, get.Endpoint)
	}

	return false
}
