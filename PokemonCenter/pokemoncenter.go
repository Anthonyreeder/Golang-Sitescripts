package pokemoncenter

import (
	client "Golang-Sitescripts/client"
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

//Entry point for Shopify Demo
func PokemonCENTER() {
	proxy := client.Proxy{}

	client.SetupClient(proxy)
	PokemonCenterLogin()
	//	FroneEndDemo()
}

type AddToCartRequest struct {
	Configuration string `json:"configuration"`
	ProductUri    string `json:"productURI"`
	Quantity      int    `json:"quantity"`
}

func PokemonCenterLogin() bool {
	params := url.Values{}
	params.Add("username", "anthonyreeder123@gmail.com")
	params.Add("password", "thekid225")
	params.Add("grant_type", "password")
	params.Add("role", "REGISTERED")
	params.Add("scope", "pokemon")
	fmt.Println(params.Encode())

	post := client.POST{
		Endpoint: "https://www.pokemoncenter.com/tpci-ecommweb-api/auth?format=zoom.nodatalinks",
		Payload:  bytes.NewReader([]byte(params.Encode())),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{cookie: []string{"datadome=1P4Dt6zlnWd5.eYx2jmbsGBlvsP1YxuUFvJUadelRLmaAkqfsqfJuo~vETHB3UB9R3d8fZPD3bpiaKFpb2FJAYTFACD5TqEZd3z.quyte5"}, content: nil, contentType: "json"}, "www.pokemoncenter.com")
	_, resp := client.NewResponse(request)
	switch resp.StatusCode {
	case 200:
		return true
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, post.Endpoint)
	}
	fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, post.Endpoint)
	return false
}

func AddToCart() bool {

	payloadBytes, _ := json.Marshal(AddToCartRequest{
		Configuration: "",
		ProductUri:    "/carts/items/pokemon/qgqvhkjxgazs2mbvgqydg=/form",
		Quantity:      1,
	})

	post := client.POST{
		Endpoint: "https://www.pokemoncenter.com/tpci-ecommweb-api/cart?type=product&format=zoom.nodatalinks",
		Payload:  bytes.NewReader(payloadBytes),
	}

	request := client.NewRequest(post)
	request.Header = AddHeaders(Header{cookie: []string{}, content: nil, contentType: "json"}, "")
	_, resp := client.NewResponse(request)

	switch resp.StatusCode {
	case 200:
		return true
	default:
		fmt.Printf("unexpected status code %v when requesting : %s", resp.StatusCode, post.Endpoint)
	}

	return false
}
