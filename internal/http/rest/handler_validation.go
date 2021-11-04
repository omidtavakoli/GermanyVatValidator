package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

func (h *Handler) VatValidator(c *gin.Context) {
	param := c.Param("id")
	vatId, err := regexp.Compile("^DE\\d{9}$") // 9 digits, e.g. DE999999999
	if err != nil {
		c.JSON(http.StatusInternalServerError, "regex compile error")
		return
	}
	validId := vatId.MatchString(param)
	if !validId {
		c.JSON(http.StatusExpectationFailed, "not valid German VAT ID")
		return
	}
	c.JSON(http.StatusOK, "VALID")
}
