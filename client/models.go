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
}

type POSTUrlEncoded struct {
	Endpoint       string
	EncodedPayload string
}
