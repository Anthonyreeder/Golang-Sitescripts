package shopify

import (
	client "Golang-Sitescripts/client"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/anaskhan96/soup"
)

///TODO:
//Test all USE cases for 'Size'
//Use Regex to check the potential input that the user passes (Checking first with bossman)
//If it is NOT the offer ID passed in then we must also extract this from the original 'GET' request

//These are hard coded values which should come from the UI
var host = "https://limitededt.com"
var size = "7"
var quantity = "1"

//These are global variables placed here just to ease the burden of passing data between methods, would otherwise be handled by the supporting framework/task system
var botKey = ""
var offerId = "32521243820103" //Either passed in frrom UI or extracted. Depending on what the UI passes to the task.

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

		//Parse HTML using soup
		responseBody := soup.HTMLParse(string(respBytes))

		//Find the bot-key input field in the form
		botKeyElement := responseBody.Find("input", "id", "bot-key").Pointer.Attr
		for _, v := range botKeyElement {
			if v.Key == "value" {
				//Locate the botKey attribute value within this node
				botKey = v.Val
			}
		}

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
		OptionSize: size, //Selected size
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
	respBytes, _ := client.NewResponse(client.Client, request)

	fmt.Println(string(respBytes))

	return false
}
