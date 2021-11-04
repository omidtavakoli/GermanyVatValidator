package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

func (h *Handler) VatValidator(c *gin.Context) {
	id := c.Param("id")
	vatId, err := regexp.Compile("^DE\\d{9}$") // 9 digits, e.g. DE999999999
	if err != nil {
		c.JSON(http.StatusInternalServerError, "regex compile error")
		return
	}
	validId := vatId.MatchString(id)
	if !validId {
		c.JSON(http.StatusExpectationFailed, "not valid German VAT ID by structure")
		return
	}

	valid, err := h.ValidatorService.VatValidator(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if !valid {
		c.JSON(http.StatusNotAcceptable, "not valid German VAT ID by EU/VIES")
		return
	}
	c.JSON(http.StatusOK, "VALID")
}
