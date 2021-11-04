package validator

func (s service) VatValidator(vatNumber string) (bool, error) {
	vatResponse, err := s.app.CheckVAT(vatNumber)
	if err != nil {
		return false, err
	}
	return vatResponse.Valid, nil
}
