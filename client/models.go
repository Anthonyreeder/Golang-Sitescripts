package request

import "bytes"

type GET struct {
	Endpoint string
}

type POST struct {
	Endpoint string
	Payload  *bytes.Reader
}

type PUT struct {
	Endpoint string
	Payload  *bytes.Reader
}

type POSTUrlEncoded struct {
	Endpoint       string
	EncodedPayload string
}

type Proxy struct {
	Proxy     string
	Port      string
	ProxyUser string
	ProxyPass string
}
