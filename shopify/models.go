package shopify

import "bytes"

type Header struct {
	cookie      []string
	content     *bytes.Reader
	contentType string
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
	Name  string `json:"name"`
	Price string `json:"price"`
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
