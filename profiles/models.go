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
	Ccnumber   string
	Expiry     string
	ExpiryYear string
	Cvc        string
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
	IsoCodeShortShippingShipping string
	NameRegionShipping           string
	IsoCodeCountryShipping       string
	NameCountryShipping          string
	IsoCodeCountryBilling        string
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
