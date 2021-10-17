package shopify

import (
	client "Golang-Sitescripts/client"
	profiles "Golang-Sitescripts/profiles"
)

///TODO:
//Sites failing shipping are likely 100% working but don't accept the country were using OR we need to use the US endpoint (In which case probably also needs a USA address)
//Ones that are success but fail checkpayment, i need to check the HTML and see whats different. Hopefully they are at least uniform (or some are) in structure.
//Complete fails are likely anti bot or incorrect urls. Need to use charles on these.
/*
1. atmos - https://www.atmosusa.com/ -> 					100% Success
2. bdgastore - https://bdgastore.com/ -> 					100% Success
3. bb branded - https://www.bbbranded.com/ -> 				100% Success
4. blends US - https://www.blendsus.com/ -> 				100% Success
5. coporategotem -> https://corporategotem.com -> 			100% Success
6. dopefactory - https://www.dope-factory.com/ -> 			100% Success
7. deadstock - https://www.deadstockofficial.com/ -> 		100% Success
8. Limited Edition - https://limitededt.com/ -> 			100% Success
9. goodhood - https://goodhoodstore.com/ -> 				100% Success
10. sneakerboxshop.ca - https://sneakerboxshop.ca/ -> 		100% success
11. Xhibition - https://www.xhibition.co/ -> 				100% success
12. juicestore - https://juicestore.com/ ->  				100% Success
13. Public school NY - https://www.publicschoolnyc.com/ -> 	100% Success
14. saintalfred - https://www.saintalfred.com/ -> 			100% Success
15. Sneaker Politics - https://sneakerpolitics.com/ -> 		100% Success
16. apbstore - https://www.apbstore.com/ -> 				100% Success
17. social status - https://www.socialstatuspgh.com/ -> 	100% Success
18. culturekings -> https://www.culturekings.com -> 		100% Success
19. feature - https://feature.com/ -> 					100% Success but could use optimising
20. Closet Inc - https://www.theclosetinc.com/ -> 		100% Success but could use optimising
21. shoepalace - https://www.shoepalace.com/ -> 		100% success
22. dtlr - https://www.dtlr.com/ -> 					100% success
23. ficegallery - https://www.ficegallery.com/ -> 		100% success
24. just don - https://justdon.com/ -> 					100% success
25. solefiness - https://www.solefiness.com/ -> 		100% success

smets.lu - https://smets.lu/ -> 					Paypal only
burn rubber - https://burnrubbersneakers.com/ -> 	Paypal only
sole steal - https://www.solesteals.com/ -> 		Paypal only
	lustmexico - https://www.lustmexico.com/ -> 	Paypal only

Noirfonce EU - https://www.noirfonce.eu/ -> 		3DS (Confirm payment in revolut)
hanon shop - https://www.hanon-shop.com/ -> 		3DS (Confirm payment in revolut)

undefeated - https://undefeated.com/ -> 			Requires login
A ma maniere - https://www.a-ma-maniere.com/ -> 	Requires login
bouncewear - https://bouncewear.com/ -> 			Requires login
exclusity life - https://shop.exclucitylife.com/ -> Requires login
jimmyjazz - https://www.jimmyjazz.com -> 			Requires login
kith - https://kith.com/ -> 						Requires login

concepts - https://cncpts.com/ -> 								Antibot/Password or offline
Funko - https://www.funko.com/ -> 								Antibot/Password or offline
haven - https://havenshop.com/ -> 								Antibot/Password or offline
	Palace SB - https://www.palaceskateboards.com ->			Antibot/Password or offline
travis scott - https://www.t qravisscott.com/ -> 				Antibot/Password or offline
cpfm - https://cactusplantfleamarket.com/ -> 					Antibot/Password or offline
DDT - https://ddtstore.com/password -> 							Antibot/Password or offline
kawsone - https://kawsone.com/ -> 								Antibot/Password or offline
NIce Kick - https://www.nicekicks.com/ -> 						Antibot/Password or offline
ronniefieg - https://shop.ronniefieg.com/ -> 					Antibot/Password or offline
Dover street market - https://london.doverstreetmarket.com/ -> 	Antibot/Password or offline

	ovo - https://uk.octobersveryown.com/ -> 			Fails on get shipping ID for USA canada and singapour
	oneblock down - https://eu.oneblockdown.it/ -> 		Fails on get shipping ID for USA canada and singapour
	Packershoes - https://packershoes.com/ -> 			Fails on get shipping ID for USA canada and singapour
	Suede Store - https://suede-store.com/ -> 			Fails on get shipping ID for USA canada and singapour

Trophy Room - https://www.trophyroomstore.com/ -> 	Fails submitting customer information

//3ds example
//https://hooks.stripe.com/3d_secure_2/hosted?merchant=acct_1C83vfGrEnQm9AGb&payment_intent=pi_3JlXbXGrEnQm9AGb0nu3fZPn&payment_intent_client_secret=pi_3JlXbXGrEnQm9AGb0nu3fZPn_secret_AdnH8n9tl8YVOQBoLw35zxSOz&publishable_key=pk_live_514ke48A5jVMTWATl9JUuGyfmMtf4ldHKrk2CqII9VWHasHtbxEiJHUpkZaOq8TXIc1dAGgs3zh1zkKOMmZRheQ0I00MBg2eNi8&source=src_1JlXbYGrEnQm9AGbXkaJDNI8&stripe_account=acct_1C83vfGrEnQm9AGb


*/
//These are hard coded values which should come from the UI
var host = "https://www.solefiness.com"
var link = "https://www.solefiness.com"
var size = "7"
var quantity = "1"
var offerId = "39503380250695" //AKA Variant ID

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

//Login details
var loginEmail = "anthonyreeder123@gmail.com"
var password = "Test123"
var profile = profiles.Profiles{}

//Entry point for Shopify Demo
func Shopify() {
	//Set Profile
	profile = profiles.CanadaProfile()

	//Setup
	client.SetupClient()

	offerId = GetRandomId(true)
	FroneEndDemo()

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
	offerId = GetRandomId(false) //get an outofstock item
	startTask(ShopifyGetProductPageB, "Checking for product")
	startTask(ShopifyAddToCartStandard, "AddToCartRealId")
	startTask(ExtractPaymentGatewayId, "GetPaymentGatewayId", true)
	startTask(SubmitPayment, "SubmitThePaymentUrl")
	startTask(CheckPaymentProcess, "CheckPaymentStatus")
}

func FroneEndDemo() {
	startTask(ShopifyGetProductPageB, "Checking for product")
	startTask(CreatePaymentSession, "PaymentSession")   //Dont wait as we dont us this until the end.
	startTask(ShopifyAddToCartStandard, "AddToCartId")  //we must wait for this as 1. Its basically the monitor and 2. It wont submit checkout form without this being complete
	startTask(ExtractShippingRates, "GetShippingRates") //If we are here then ATC was succcess, run async as then if it fails it'll just fail later in the task anyway and its unrecoverable, no point in waiting.
	startTask(LoadCheckoutForm, "LoadCheckoutForm")
	startTask(SubmitCustomerInfo, "SubmitTheCustomerInfo")                //If we skip waiting here then we can start extracting the token now
	startTask(ExtractShippingToken, "ExtractTheShippingToken")            //must wait for token
	startTask(SubmitShippingMethodDetails, "SubmitShippingMethodDetails") //If we skip waiting here then we can start getting the paymentGateID NOW
	startTask(ExtractPaymentGatewayId, "GetPaymentGatewayId")             //We actually need the paymentGateway to continue, but somehow by NOT waiting for the response we shave off about 5-6s
	startTask(SubmitPayment, "SubmitThePaymentUrl")                       //WE don't want to check the paymentStatus UNTIL this is sent so wait
	startTask(CheckPaymentProcess, "CheckPaymentStatus")                  //Loop checkStatus
}
