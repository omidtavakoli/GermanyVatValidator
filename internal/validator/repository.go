package validator

import "VatIdValidator/pkg/EU_VIES"

type AppRepository interface {
	CheckVAT(vatNumber string) (*EU_VIES.VATresponse, error)
}
