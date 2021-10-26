package footsites

type AydenEncrpytLocal struct {
	DfValueMain      string
	EncryptionKey    string
	CreditCardNUmber CreditCardNUmber
	MonthExpiry      MonthExpiry
	YearExpiry       YearExpiry
	CvcNumb          CvcNumb
}
type CreditCardNUmber struct {
	Activate            string `json:"activate"`
	Generationtime      string `json:"generationtime"`
	InitializeCount     string `json:"initializeCount"`
	LuhnCount           string `json:"luhnCount"`
	LuhnOkCount         string `json:"luhnOkCount"`
	LuhnSameLengthCount string `json:"luhnSameLengthCount"`
	Number              string `json:"number"`
}
type MonthExpiry struct {
	Activate        string `json:"activate"`
	ExpiryMonth     string `json:"expiryMonth"`
	Generationtime  string `json:"generationtime"`
	InitializeCount string `json:"initializeCount"`
}
type YearExpiry struct {
	Activate        string `json:"activate"`
	ExpiryYear      string `json:"expiryYear"`
	Generationtime  string `json:"generationtime"`
	InitializeCount string `json:"initializeCount"`
}
type CvcNumb struct {
	Activate        string `json:"activate"`
	Cvc             string `json:"cvc"`
	Generationtime  string `json:"generationtime"`
	InitializeCount string `json:"initializeCount"`
}
