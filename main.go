package main

import (
	"Golang-Sitescripts/pokemoncenter"
	"Golang-Sitescripts/profiles"
	"Golang-Sitescripts/shopify"
	"flag"
	"fmt"
)

//global "Golang-Sitescripts/global"
var profileInfo profiles.ShopifyProfiles
var paymentDetails shopify.PaymentDetails
var productInfo shopify.ProductInfo
var loginDetails shopify.LoginDetails

var UseArgs bool

func main() {
	UseArgs = false
	ParseArgs()
	pokemoncenter.PokemonCENTER()
	//shopify.Shopify(profileInfo, productInfo, paymentDetails, loginDetails)

}

func ParseArgs() {
	if UseArgs {
		//Load profile into from flags into struct
		flag.StringVar(&profileInfo.Email, "Email", "", "Profile information")
		flag.StringVar(&profileInfo.FirstName, "FirstName", "", "Profile information")
		flag.StringVar(&profileInfo.LastName, "LastName", "", "Profile information")
		flag.StringVar(&profileInfo.Company, "Company", "", "Profile information")
		flag.StringVar(&profileInfo.Address1, "Address1", "", "Profile information")
		flag.StringVar(&profileInfo.Address2, "Address2", "", "Profile information")
		flag.StringVar(&profileInfo.City, "City", "", "Profile information")
		flag.StringVar(&profileInfo.Country, "Country", "", "Profile information")
		flag.StringVar(&profileInfo.PostCode, "PostCode", "", "Profile information")
		flag.StringVar(&profileInfo.Phone, "Phone", "", "Profile information")
		flag.StringVar(&profileInfo.Province, "Province", "", "Profile information")

		flag.StringVar(&paymentDetails.CardNumber, "CardNumber", "", "Profile information")
		flag.StringVar(&paymentDetails.Month, "Month", "", "Profile information")
		flag.StringVar(&paymentDetails.Year, "Year", "", "Profile information")
		flag.StringVar(&paymentDetails.CardName, "CardName", "", "Profile information")
		flag.StringVar(&paymentDetails.Ccv, "Ccv", "", "Profile information")

		flag.StringVar(&productInfo.Host, "Host", "", "Profile information")
		flag.StringVar(&productInfo.Link, "Link", "", "Profile information")
		flag.StringVar(&productInfo.OfferId, "OfferId", "", "Profile information")
		flag.StringVar(&productInfo.Quantity, "Quantity", "", "Profile information")

		flag.StringVar(&loginDetails.LoginEmail, "LoginEmail", "", "Profile information")
		flag.StringVar(&loginDetails.Password, "Password", "", "Profile information")
		flag.StringVar(&loginDetails.Profile, "Profile", "", "Profile information")

		flag.StringVar(&loginDetails.Proxy, "Proxy", "", "Profile information")
		flag.StringVar(&loginDetails.Port, "Port", "", "Profile information")
		flag.StringVar(&loginDetails.ProxyUser, "ProxyUser", "", "Profile information")
		flag.StringVar(&loginDetails.ProxyPass, "ProxyPass", "", "Profile information")
		flag.Parse()
		fmt.Println("Profile information:")
		fmt.Println(profileInfo)
		fmt.Println("Payment details:")
		fmt.Println(paymentDetails)
		fmt.Println("Product info:")
		fmt.Println(productInfo)
		fmt.Println("Login details")
		fmt.Println(loginDetails)
	} else {
		profileInfo = profiles.ShopifyProfiles{
			Email:     "JohnSmith5318008@gmail.com",
			FirstName: "John",
			LastName:  "Smith",
			Company:   "",
			Address1:  "2382 Hickory Ridge Drive",
			Address2:  "",
			City:      "Las Vegas",
			Country:   "US",
			PostCode:  "89108",
			Phone:     "702-645-3077",
			Province:  "NV",
		}
		paymentDetails = shopify.PaymentDetails{
			CardNumber: "5354568000637394",
			Month:      "03",
			Year:       "2026",
			CardName:   "Mr John Smith",
			Ccv:        "960",
		}
		productInfo = shopify.ProductInfo{
			Host:     "https://limitededt.com",
			Link:     "https://limitededt.com",
			OfferId:  "39515671396423",
			Quantity: "1",
		}
		loginDetails = shopify.LoginDetails{
			LoginEmail: "LoginEmail",
			Password:   "Password",
			Profile:    "Profile",
			Proxy:      "103.7.206.68",

			Port:      "58262",
			ProxyUser: "run",
			ProxyPass: "iF97ZAGT",
		}
	}
}
