package shopify

import (
	client "Golang-Sitescripts/client"
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

///TODO:
//Test all USE cases for 'Size'
//Use Regex to check the potential input that the user passes (Checking first with bossman)
//If it is NOT the offer ID passed in then we must also extract this from the original 'GET' request

//These are hard coded values which should come from the UI
var host = "https://limitededt.com"
var size = "7"
var quantity = "1"

//Profile information
var email = "JohnSmith5318008@gmail.com"
var fname = "John"
var lname = "Smith"
var company = "Torchwood"
var addy1 = "1 Castle St"
var addy2 = ""
var city = "Cardiff"
var country = "United kingdom"
var postal_code = "CF10 3RB"
var phone = "01763220883"
var province = "Cardiff"

//Global variables to ease the burden of passing data between methods, would otherwise be handled by the supporting framework/task system
var authKey = ""
var botKey = ""
var offerId = "32521243820103" //Either passed in frrom UI or extracted. Depending on what the UI passes to the task.
var formUrl = ""

//Entry point for Shopify Demo
func Shopify() {
	//Setup
	client.SetupClient()

	//Get the shopify page to set the bot-key used in addtoCart
	fmt.Println("Getting shopify page")
	if !ShopifyGetProductPage() {
		fmt.Println("Failed to get page")
	}

	//Now the bot-key is set we add the product to cart
	fmt.Println("Adding product to cart")
	if !ShopifyAddToCartStandard() {
		fmt.Println("Failed to add to cart")
	}

	//Load the checkout form so we can extract the AuthId
	fmt.Println("Loading the checkoutForm")
	if !LoadCheckoutForm() {
		fmt.Println("Failed to load checkoutForm")
	}

	//Submit the profile information
	fmt.Println("Submitting customer information")
	if !SubmitCustomerInfo() {
		fmt.Println("Failed to submit customer information")
	}

	GetShippingToken()

}

//GET Product page and extract bot-key
func ShopifyGetProductPage() bool {
	get := client.GET{
		Endpoint: fmt.Sprintf("%s/collections/mens/products/adidas-originals-pharrell-williams-boost-slides-fy6140", host),
	}
	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	respBytes, resp := client.NewResponse(client.Client, request)

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
		fmt.Printf("GET request to %s returned the following unexpected response code: %v", get.Endpoint, resp.StatusCode)
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
	_, resp := client.NewResponse(client.Client, request)

	switch resp.StatusCode {
	case 200:
		fmt.Println("Successfully added item to cart ")
		return true

	default:
		fmt.Printf("request to %s returned the following unexpected response code: %v", post.Endpoint, resp.StatusCode)
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
	respBytes, resp := client.NewResponse(client.Client, request)

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
		fmt.Printf("GET request to %s returned the following unexpected response code: %v", get.Endpoint, resp.StatusCode)
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
	_, resp := client.NewResponse(client.Client, request)

	switch resp.StatusCode {
	case 200:
		fmt.Printf("Successfully posted the customer information for %s\n", email)
		return true

	default:
		fmt.Printf("request to %s returned the following unexpected response code: %v", post.Endpoint, resp.StatusCode)
	}

	return false
}

func GetShippingToken() bool {
	get := client.GET{
		Endpoint: fmt.Sprintf("%s/cart/shipping_rates.json?shipping_address[zip]=%s&shipping_address[country]=%s&shipping_address[province]=%s", host, postal_code, country, province),
	}

	request := client.NewRequest(get)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil}, host)
	respBytes, resp := client.NewResponse(client.Client, request)

	switch resp.StatusCode {
	case 200:
		//fmt.Println("Shipping token request loaded, extracting shipping token")
		fmt.Println(string(respBytes))

	default:
		fmt.Printf("GET request to %s returned the following unexpected response code: %v", get.Endpoint, resp.StatusCode)
	}

	return false
}
