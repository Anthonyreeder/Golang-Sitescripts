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
var company = ""
var addy1 = "37 Shenton Way"
var addy2 = ""
var city = ""
var country = "Singapore"
var postal_code = "068811"
var phone = "68246580"
var province = ""

//Card details
var cardNumber = "5354568000637394"
var name = "JohnSmith"
var month = "03"
var year = "2026"
var ccv = "960"

//Scoped-access variables to ease the burden of passing data between methods, would otherwise be handled by the supporting framework/task system
var authKey = ""
var botKey = ""
var gatewayKey = ""
var offerId = "32521243820103" //AKA Variant ID
var formUrl = ""
var shipping_option = ""
var total_amount = ""
var payment_token = ""
var process_url = ""

//Entry point for Shopify Demo
func Shopify() {
	//Setup
	client.SetupClient()

	//Get the shopify page to set the bot-key used in addtoCart
	fmt.Println("Getting shopify page")
	if !ShopifyGetProductPage() {
		fmt.Println("Failed to get page")
	}

	//Get the payment session, we will need this later for the payment
	if !CreatePaymentSession() {
		fmt.Println("Failed to add to cart")
	}

	//Now the bot-key is set we add the product to cart
	fmt.Println("Adding product to cart")
	if !ShopifyAddToCartStandard() {
		fmt.Println("Failed to add to cart")
	}

	//Get the shipping ID we need to submit shipping
	fmt.Println("Grabbing the shipping id")
	if !ExtractShippingRates() {
		fmt.Println("Failed to get the shipping id")
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

	//Get the shipping token we need to submit shipping
	fmt.Println("Extacting the shipping token")
	if !ExtractShippingToken() {
		fmt.Println("Failed to extact the shipping token")
	}

	//Submit the shipping
	fmt.Println("Submitting the shipping method details")
	if !SubmitShippingMethodDetails() {
		fmt.Println("Failed to submit the shipping method details")
	}

	//Get the payment gateway ID we need to process payment
	fmt.Println("Extracting payment gateway Id")
	if !ExtractPaymentGatewayId() {
		fmt.Println("Failed to extract the payment gateway Id")
	}

	//Submit payment details
	fmt.Println("Submitting payment details")
	if !SubmitPayment() {
		fmt.Println("Failed to submit payment details")
	}

	//Check if processing has finished and read result
	fmt.Println("CheckPaymentProcess")
	if !CheckPaymentProcess() {
		fmt.Println("Payment isn't finished")
	}
}
