package shopify

import (
	client "Golang-Sitescripts/client"
	"fmt"
)

///TODO:
//Test all USE cases for 'Size'
//Use Regex to check the potential input that the user passes (Checking first with bossman)
//If it is NOT the offer ID passed in then we must also extract this from the original 'GET' request
//AddToCart add all response USE cases such as OOS, Sever not responding, Added to cart, Check quantity is correct

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
var offerId = "32521243820103" //AKA Variant ID
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
