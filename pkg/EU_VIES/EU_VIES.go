package EU_VIES

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"text/template"
	"time"
)

const envelope = `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:v1="http://schemas.conversesolutions.com/xsd/dmticta/v1">
<soapenv:Header/>
<soapenv:Body>
 <checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
   <countryCode>{{.countryCode}}</countryCode>
   <vatNumber>{{.vatNumber}}</vatNumber>
 </checkVat>
</soapenv:Body>
</soapenv:Envelope>
`

var (
	ErrVATnumberNotValid     = errors.New("vat number is not valid")
	ErrVATserviceUnreachable = errors.New("vat number validation service is offline")
	ErrVATserviceError       = "VAT number validation service returns an error : "
)

func NewClient(config *Config) (*client, error) {
	httpClient := &http.Client{
		Timeout: 10000 * time.Millisecond,
	}
	return &client{
		HttpClient: httpClient,
		Config:     config,
	}, nil
}

func (c client) CheckVAT(vatNumber string) (*VATresponse, error) {
	vatNumber = sanitizeVatNumber(vatNumber)

	e, err := getEnvelope(vatNumber)
	if err != nil {
		return nil, err
	}
	eb := bytes.NewBufferString(e)

	res, err := c.HttpClient.Post(c.Config.Url, "text/xml;charset=UTF-8", eb)
	if err != nil {
		return nil, ErrVATserviceUnreachable
	}
	defer res.Body.Close()
	xmlRes, err := ioutil.ReadAll(res.Body)

	if bytes.Contains(xmlRes, []byte("INVALID_INPUT")) {
		return nil, ErrVATnumberNotValid
	}

	var response Response

	if err := xml.Unmarshal(xmlRes, &response); err != nil {
		return nil, err
	}

	if response.Soap.SoapFault.Message != "" {
		return nil, errors.New(ErrVATserviceError + response.Soap.SoapFault.Message)
	}

	if response.Soap.Soap.RequestDate == "" {
		return nil, errors.New("service returned invalid request date")
	}

	pDate, err := time.Parse("2006-01-02-07:00", response.Soap.Soap.RequestDate)
	if err != nil {
		return nil, err
	}

	r := &VATresponse{
		CountryCode: response.Soap.Soap.CountryCode,
		VATnumber:   response.Soap.Soap.VATnumber,
		RequestDate: pDate,
		Valid:       response.Soap.Soap.Valid,
		Name:        response.Soap.Soap.Name,
		Address:     response.Soap.Soap.Address,
	}

	return r, nil
}

// sanitizeVatNumber removes white space
func sanitizeVatNumber(vatNumber string) string {
	vatNumber = strings.TrimSpace(vatNumber)
	return regexp.MustCompile(" ").ReplaceAllString(vatNumber, "")
}

// getEnvelope parses envelope template
func getEnvelope(vatNumber string) (string, error) {
	if len(vatNumber) < 9 {
		return "", errors.New("vat number is too short for Germany")
	}

	t, err := template.New("envelope").Parse(envelope)
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	if err := t.Execute(&result, map[string]string{
		"countryCode": strings.ToUpper(vatNumber[0:2]),
		"vatNumber":   vatNumber[2:],
	}); err != nil {
		return "", err
	}
	return result.String(), nil
}
