package profiles

type ShopifyProfiles struct {
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

type FootSitesSessionInfo struct {
	RequestAgent    string
	AexOffset       string
	Browser         string
	Version         string
	OsName          string
	Appname         string
	AppPlatform     string
	Height          string
	Width           string
	AllPlugins      string
	Referer         string
	IntLoc          string
	GetOffset       string
	RequestLanguage string
}

type FootSitesProfile struct {
	Person      Person
	Shipping    Shipping
	Billing     Billing
	CardDetails CardDetails
}

type CardDetails struct {
	Ccnumber   int
	Expiry     int
	ExpiryYear int
	Cvc        int
}
type Person struct {
	Email string
	Phone string
}
type Shipping struct {
	FirstNameShipping            string
	LastNameShipping             string
	Line1Shipping                string
	Line2Shipping                string
	PostalCodeShipping           string
	RecordTypeShipping           string
	TownShipping                 string
	CountryIsoRegionShipping     string
	IsoCodeRegionShipping        string
	Isocodeshortshippingshipping string
	Nameregionshipping           string
	Isocodecountryshipping       string
	Namecountryshipping          string
	Isocodecountrybilling        string
}

type Billing struct {
	Namecountrybilling        string
	FirstNamebilling          string
	LastNamebilling           string
	Line1billing              string
	Line2billing              string
	Postalcodebilling         string
	Recordtypebilling         string
	Townbilling               string
	Countryisoregionbilling   string
	Isocoderegionbilling      string
	Isocodeshortregionbilling string
	Nameregionbilling         string
}
