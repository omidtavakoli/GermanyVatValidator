package EU_VIES

import (
	"encoding/xml"
	"net/http"
	"time"
)

type Config struct {
	Url string `yaml:"URL"`
}

type client struct {
	HttpClient *http.Client
	Ready      bool
	Config     *Config
}

type VATresponse struct {
	CountryCode string
	VATnumber   string
	RequestDate time.Time
	Valid       bool
	Name        string
	Address     string
}

type Response struct {
	XMLName xml.Name `xml:"Envelope"`
	Soap    struct {
		XMLName xml.Name `xml:"Body"`
		Soap    struct {
			XMLName     xml.Name `xml:"checkVatResponse"`
			CountryCode string   `xml:"countryCode"`
			VATnumber   string   `xml:"vatNumber"`
			RequestDate string   `xml:"requestDate"` // 2015-03-06+01:00
			Valid       bool     `xml:"valid"`
			Name        string   `xml:"name"`
			Address     string   `xml:"address"`
		}
		SoapFault struct {
			XMLName string `xml:"Fault"`
			Code    string `xml:"faultcode"`
			Message string `xml:"faultstring"`
		}
	}
}
