package shopify

import (
	client "Golang-Sitescripts/client"
	"fmt"
	"time"
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
var offerId = "32521243820103" //AKA Variant ID

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
	start := time.Now()
	FroneEndDemo()
	elapsed := time.Since(start)
	fmt.Println(elapsed)

	//15.8571866s
	//Checkout time -> 20.3942374s
	//Implemented Goroutine handling of tasks
	//Checkout time -> 9.6317531s

	/*
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			FroneEndDemo()
		}()
		go func() {
			defer wg.Done()
			FrontEndPreCartDemo()
		}()

		wg.Wait()
	*/

}

var taskComplete = false

//PreCart demo is a concept ideea not fully tested.
//The ideea being that in order to create a shipping session an item is required in cart, once this has been created we are able to submit shipping, customer and create a payment agreement.
//Changing the item quantity to 0 removes it frrom cart but does not remove out shipping session or payment agreement. As soon as we add to cart again we instantl try and submit our payment agreement and checkout.
func FrontEndPreCartDemo() {
	startTask(CreatePaymentSession, "PaymentSession")
	startTask(ShopifyAddToCartStandard, "AddToCartFakeId")
	startTask(ExtractShippingRates, "GetShippingRates")
	startTask(LoadCheckoutForm, "LoadCheckoutForm")
	startTask(SubmitCustomerInfo, "SubmitTheCustomerInfo")
	startTask(ExtractShippingToken, "ExtractTheShippingToken")
	startTask(SubmitShippingMethodDetails, "SubmitShippingMethodDetails")

	//The monitoring starts here. Now that the session has been created and payment agreement has been made. We remove the fakePID
	startTaskInt(ShopifyChangeCart, "RemoveFakeProductFromCart", 0)

	//Now we set our offerID to the REAL pid and attempt to add to cart
	offerId = "39499044323399"
	startTask(ShopifyAddToCartStandard, "AddToCartRealId")
	startTask(ExtractPaymentGatewayId, "GetPaymentGatewayId", true)
	startTask(SubmitPayment, "SubmitThePaymentUrl")
	startTask(CheckPaymentProcess, "CheckPaymentStatus")
}

func FroneEndDemo() {
	startTask(CreatePaymentSession, "PaymentSession", true)   //Dont wait as we dont us this until the end.
	startTask(ShopifyAddToCartStandard, "AddToCartId")        //we must wait for this as 1. Its basically the monitor and 2. It wont submit checkout form without this being complete
	startTask(ExtractShippingRates, "GetShippingRates", true) //If we are here then ATC was succcess, run async as then if it fails it'll just fail later in the task anyway and its unrecoverable, no point in waiting.
	startTask(LoadCheckoutForm, "LoadCheckoutForm")
	startTask(SubmitCustomerInfo, "SubmitTheCustomerInfo", true)                //If we skip waiting here then we can start extracting the token now
	startTask(ExtractShippingToken, "ExtractTheShippingToken")                  //must wait for token
	startTask(SubmitShippingMethodDetails, "SubmitShippingMethodDetails", true) //If we skip waiting here then we can start getting the paymentGateID NOW
	startTask(ExtractPaymentGatewayId, "GetPaymentGatewayId", true)             //We actually need the paymentGateway to continue, but somehow by NOT waiting for the response we shave off about 5-6s
	startTask(SubmitPayment, "SubmitThePaymentUrl")                             //WE don't want to check the paymentStatus UNTIL this is sent so wait
	startTask(CheckPaymentProcess, "CheckPaymentStatus")                        //Loop checkStatus
}
