package shopify

import (
	client "Golang-Sitescripts/client"
)

///TODO:
//Sites failing shipping are likely 100% working but don't accept the country were using OR we need to use the US endpoint (In which case probably also needs a USA address)
//Ones that are success but fail checkpayment, i need to check the HTML and see whats different. Hopefully they are at least uniform (or some are) in structure.
//Complete fails are likely anti bot or incorrect urls. Need to use charles on these.
/*
A ma maniere - https://www.a-ma-maniere.com/ -> Fails on Submitting the shipping method.
apbstore - https://www.apbstore.com/ -> Success but 'CheckPaymentStatus' throws Soup error (html parse)
atmos - https://www.atmosusa.com/ -> Fails on 'GetShippingRates' -> Likely because its USA and I think they use another URL (I have this somewhere)
bdgastore - https://bdgastore.com/ -> 100% Success
bb branded - https://www.bbbranded.com/ -> Fails on 'GetShippingRates' -> Country not supported
burn rubber - https://burnrubbersneakers.com/ -> Success but 'CheckPaymentStatus' throws Soup error
blends US - https://www.blendsus.com/ -> 100% Success
bouncewear - https://bouncewear.com/ -> Fails on LoadCheckoutForm (Require login maybe?)
concepts - https://cncpts.com/ -> complete fail won't get product (antibot?)
Closet Inc - https://www.theclosetinc.com/ -> Fails on GetShippingRates
cpfm - https://cactusplantfleamarket.com/ -> Site seems to be offline
culturekings -> https://www.culturekings.com -> Fails on GetShippingRates
coporategotem -> https://corporategotem.com/ -> Fails on GetShippingRates
Dover street market (US/UK/JP/SG) - https://london.doverstreetmarket.com/ -> Complete fail cant find products
deadstock - https://www.deadstockofficial.com/ -> 100% Success
dtlr - https://www.dtlr.com/ -> Fails on GetShippingRates
DDT - https://ddtstore.com/password -> Fails to get product
dopefactory - https://www.dope-factory.com/ -> 100% Success
exclusity life - https://shop.exclucitylife.com/ -> Fails on get shipping ID
Funko - https://www.funko.com/ -> Fails to get product (anti bot?)
ficegallery - https://www.ficegallery.com/ -> Fails on get shipping ID
feature - https://feature.com/ -> Success but 'CheckPaymentStatus' throws soup error
goodhood - https://goodhoodstore.com/ -> Fails on get shipping ID
haven - https://havenshop.com/ -> Fails, I thnk this URL is wrong
hanon shop - https://www.hanon-shop.com/ -> Success but 'CheckPaymentStatus' throws soup error
jimmyjazz - https://www.jimmyjazz.com/ -> Fails on get shipping ID
just don - https://justdon.com/ -> Success but 'CheckPaymentStatus' throws soup error
juicestore - https://juicestore.com/ -> 100% Success
kawsone - https://kawsone.com/ -> Site is offline
kith - https://kith.com/ -> Failed on Submit customer info
Limited Edition - https://limitededt.com/ -> 100% Success
lustmexico - https://www.lustmexico.com/ -> Fails on get shipping ID
Noirfonce EU - https://www.noirfonce.eu/ -> Success but 'CheckPaymentStatus' throws soup error
NIce Kick - https://www.nicekicks.com/ -> Can't get product
ovo - https://uk.octobersveryown.com/ -> Fails on get shipping ID
oneblock down - https://eu.oneblockdown.it/ -> Fails on get shipping ID
Public school NY - https://www.publicschoolnyc.com/ -> 100% Success
Packershoes - https://packershoes.com/ -> Fails on get shipping ID
	Palace SB - https://www.palaceskateboards.com -> cannot get product 404 (URL is wrong)
ronniefieg - https://shop.ronniefieg.com/ -> cannot get product
Sneaker Politics - https://sneakerpolitics.com/ -> Fails on get shipping ID
Suede Store - https://suede-store.com/ -> Fails on get shipping ID
social status PGH - https://www.socialstatuspgh.com/ -> Success but 'CheckPaymentStatus' throws soup error
sole steal - https://www.solesteals.com/ -> Fails on get shipping ID
sneakerboxshop.ca - https://sneakerboxshop.ca/ -> 100% success
solefiness - https://www.solefiness.com/ -> Success but 'CheckPaymentStatus' throws soup error
saintalfred - https://www.saintalfred.com/ -> Fails on get shipping ID
shoepalace - https://www.shoepalace.com/ -> Fails on get shipping ID
smets.lu - https://smets.lu/ -> Success but 'CheckPaymentStatus' throws soup error
travis scott - https://www.travisscott.com/ -> Cant get product
Trophy Room - https://www.trophyroomstore.com/ -> Fails on get shipping ID
undefeated - https://undefeated.com/ -> Fails on get shipping ID
Xhibition - https://www.xhibition.co/ -> 100% success
*/
//These are hard coded values which should come from the UI
var host = "https://limitededt.com/"
var link = "https://limitededt.com/"
var size = "7"
var quantity = "1"
var offerId = "39503380250695" //AKA Variant ID

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

//Login details
var loginEmail = "anthonyreeder123@gmail.com"
var password = "Test123"

//Entry point for Shopify Demo
func Shopify() {
	//Setup
	client.SetupClient()

	host = "https://limitededt.com/"
	link = "https://limitededt.com/"

	offerId = GetRandomId(true)
	FrontEndPreCartDemo()

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
