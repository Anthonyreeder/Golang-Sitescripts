package shopify

import "bytes"

type Header struct {
	cookie      []string
	content     *bytes.Reader
	contentType string
}

type AddToCartStandardRequest struct {
	FormType   string `json:"form_type"`
	Utf8       string `json:"utf8"`
	Properties struct {
		BotKey string `json:"bot-key"`
	} `json:"properties"`
	OptionSize string `json:"option-Size"`
	Id         string `json:"id"`
	Quantity   string `json:"quantity"`
}
