package shopify

import (
	client "Golang-Sitescripts/client"
)

///TODO:
//Test all USE cases for 'Size'
//Use Regex to check the potential input that the user passes (Checking first with bossman)
//If it is NOT the offer ID passed in then we must also extract this from the original 'GET' request
//AddToCart add all response USE cases such as OOS, Sever not responding, Added to cart, Check quantity is correct

//These are hard coded values which should come from the UI
var host = "https://limitededt.com"
var link = "https://limitededt.com"
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
var _offerid = ""

//Entry point for Shopify Demo
func Shopify() {
	//Setup
	client.SetupClient()

	//PlayFunctions(ShopifyGetProductPage())
	FroneEndDemo()
}

var taskComplete = false

func FrontEndPreCartDemo() {
	startTask(CreatePaymentSession, "PaymentSession")
	_offerid = "32521243820103"
	startTask(ShopifyAddToCartStandard, "AddToCartFakeId")
	startTask(ExtractShippingRates, "GetShippingRates")
	startTask(LoadCheckoutForm, "LoadCheckoutForm")
	startTask(SubmitCustomerInfo, "SubmitTheCustomerInfo")
	startTask(ExtractShippingToken, "ExtractTheShippingToken")
	startTask(SubmitShippingMethodDetails, "SubmitShippingMethodDetails")
	startTaskInt(ShopifyChangeCart, "RemoveFakeProductFromCart", 0)
	_offerid = "39499044323399"
	startTask(ShopifyAddToCartStandard, "AddToCartRealId")
	startTask(ExtractPaymentGatewayId, "GetPaymentGatewayId")
	startTask(SubmitPayment, "SubmitThePaymentUrl")
	startTask(CheckPaymentProcess, "CheckPaymentStatus")
}

func FroneEndDemo() {
	//_offerid = "39488656244795"
	startTask(CreatePaymentSession, "PaymentSession")
	startTask(ShopifyAddToCartStandard, "AddToCartId")
	startTask(ExtractShippingRates, "GetShippingRates")
	startTask(LoadCheckoutForm, "LoadCheckoutForm")
	startTask(SubmitCustomerInfo, "SubmitTheCustomerInfo")
	startTask(ExtractShippingToken, "ExtractTheShippingToken")
	startTask(SubmitShippingMethodDetails, "SubmitShippingMethodDetails")
	startTask(ExtractPaymentGatewayId, "GetPaymentGatewayId")
	startTask(SubmitPayment, "SubmitThePaymentUrl")
	startTask(CheckPaymentProcess, "CheckPaymentStatus")
}
