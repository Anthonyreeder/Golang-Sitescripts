package profiles

func CanadaProfile() ShopifyProfiles {
	profiles := ShopifyProfiles{
		Email:     "JohnSmith5318008@gmail.com",
		FirstName: "John",
		LastName:  "Smith",
		Company:   "",
		Address1:  "2029 Old Spallumcheen Rd",
		Address2:  "",
		City:      "Keremeos",
		Country:   "Canada",
		PostCode:  "V0X1N0",
		Phone:     "250-499-1435",
		Province:  "BC",
	}

	return profiles
}

func UsaProfile() ShopifyProfiles {
	profiles := ShopifyProfiles{
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

	return profiles
}

func SingaporeProfile() ShopifyProfiles {
	profiles := ShopifyProfiles{
		Email:     "JohnSmith5318008@gmail.com",
		FirstName: "John",
		LastName:  "Smith",
		Company:   "",
		Address1:  "37 Shenton Way",
		Address2:  "",
		City:      "",
		Country:   "Singapore",
		PostCode:  "068811",
		Phone:     "68246580",
		Province:  "",
	}

	return profiles
}

func MexicoProfile() ShopifyProfiles {
	profiles := ShopifyProfiles{
		Email:     "JohnSmith5318008@gmail.com",
		FirstName: "John",
		LastName:  "Smith",
		Company:   "",
		Address1:  "RIVA PALACIO 2106 S/N",
		Address2:  "ZONA CENTRO",
		City:      "Chihuahua",
		Country:   "Mexico",
		PostCode:  "31000",
		Phone:     "614 414 1442",
		Province:  "CHIH",
	}

	return profiles
}

func FootsitesSessionInfoUk() FootSitesSessionInfo {

	session := FootSitesSessionInfo{
		RequestAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
		AexOffset:       "-60",
		Browser:         "Netscape",
		Version:         "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
		OsName:          "Win32",
		Appname:         "Netscape",
		AppPlatform:     "Win32",
		Height:          "864",
		Width:           "1536",
		AllPlugins:      "internal-pdf-viewer;mhjfbmdgcfjbbpaeojofohoefgiehjai;internal-nacl-plugin;",
		Referer:         "",
		IntLoc:          "https://www.footlocker.com/checkout",
		GetOffset:       "-60",
		RequestLanguage: "#en-US,en;q=0.9,sr;q=0.8",
	}

	return session
}

func FootSitesProfileUk() FootSitesProfile {
	profile := FootSitesProfile{
		Person: Person{
			Email: "",
			Phone: "",
		},
		Shipping: Shipping{
			FirstNameShipping:            "",
			LastNameShipping:             "",
			Line1Shipping:                "",
			Line2Shipping:                "",
			PostalCodeShipping:           "",
			RecordTypeShipping:           "",
			TownShipping:                 "",
			CountryIsoRegionShipping:     "",
			IsoCodeRegionShipping:        "",
			Isocodeshortshippingshipping: "",
			Nameregionshipping:           "",
			Isocodecountryshipping:       "",
			Namecountryshipping:          "",
			Isocodecountrybilling:        "",
		},
		Billing: Billing{
			Namecountrybilling:        "",
			FirstNamebilling:          "",
			LastNamebilling:           "",
			Line1billing:              "",
			Line2billing:              "",
			Postalcodebilling:         "",
			Recordtypebilling:         "",
			Townbilling:               "",
			Countryisoregionbilling:   "",
			Isocoderegionbilling:      "",
			Isocodeshortregionbilling: "",
			Nameregionbilling:         "",
		},
		CardDetails: CardDetails{
			Ccnumber:   0,
			Expiry:     0,
			ExpiryYear: 0,
			Cvc:        0,
		},
	}

	return profile
}
