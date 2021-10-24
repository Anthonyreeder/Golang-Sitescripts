package footsites

import "bytes"

type Request struct {
	Request interface{}
	Host    string
}

type Header struct {
	cookie            []string
	additionalHeaders []additionalHeaders
	content           *bytes.Reader
	contentType       string
}

type additionalHeaders struct {
	key   string
	value string
}

type AddToCartStandardRequest struct {
	Id       string `json:"id"`
	Quantity string `json:"quantity"`
	FormType string `json:"form_type"` // "product",
}

type ChangeCartStandardRequest struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type ShippingRate struct {
	Name   string `json:"name"`
	Price  string `json:"price"`
	Source string `json:"source"`
	Code   string `json:"code"`
	Markup string `json:"markup"`
}

type ShippingMethodResponse struct {
	ShippingRates []ShippingRate `json:"shipping_rates"`
}

type PaymentSessionRequest struct {
	CreditCard CreditCard `json:"credit_card"`
}

type CreditCard struct {
	Number             string `json:"number"`
	Name               string `json:"name"`
	Month              string `json:"month"`
	Year               string `json:"year"`
	Verification_value string `json:"verification_value"`
}

type PaymentSessionResponse struct {
	Id string `json:"id"`
}

type Product struct {
	Products []ProductData `json:"products"`
}
type ProductData struct {
	Id       int       `json:"id"`
	Title    string    `json:"title"`
	Handle   string    `json:"handle"`
	Variants []Variant `json:"variants"`
}

type Variant struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Sku       string `json:"sku"`
	Available bool   `json:"available"`
}

type Products []ProductData

type Profiles struct {
	Email     string
	FirstName string
	LastName  string
	Company   string
	Address1  string
	Address2  string
	City      string
	Country   string
	PostCode  string
	Phone     string
	Province  string
}

type CsrfResponse struct {
	Data struct {
		CsrfToken string `json:"csrfToken"`
	}
}

type GetCartRequest struct {
	ProductQuantity int    `json:"productQuantity"`
	ProductId       string `json:"productId"`
}

type VerificationAddressPayload struct {
	Country            Country `json:"country"`
	Line1Shipping      string  `json:"line1"`
	Line2Shipping      string  `json:"line2"`
	PostalCodeShipping string  `json:"postalCode"`
	TownShipping       string  `json:"town"`
	Region             Region  `json:"region"`
}
type Country struct {
	IsoCodeCountryShipping string `json:"isocode"`
	NameCountryShipping    string `json:"name"`
}

type Region struct {
	CountryIsoRegionShipping string `json:"countryIso"`
	IsoCodeRegionShipping    string `json:"isocode"`
	IsoCodeShortShipping     string `json:"isocodeShort"`
	NameRegionShipping       string `json:"name"`
}
