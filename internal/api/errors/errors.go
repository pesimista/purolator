package errors

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pesimista/purolator-api/internal/api/openapi"
)

// type defaultMessage struct {
// 	code    int
// 	message string
// }

func JSON(c *gin.Context, operation, message string, err error, httpCode int) {
	msg := message
	if len(strings.TrimSpace(msg)) == 0 {
		msg = fmt.Sprintf("%s", err)
	}

	fmt.Printf("%s: %s: %s\n", operation, message, err)
	c.JSON(httpCode, openapi.Error{
		Code:    httpCode,
		Message: msg,
	})
}
