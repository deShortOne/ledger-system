package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/deshortone/ledger-system/pkg/failure"
	"github.com/gin-gonic/gin"
)

func ErrorResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var failureObj *failure.Failure

		if !errors.As(err, &failureObj) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
				"misc":  "error was not fully captured",
			})

			log.Printf("Error: not fully captured: user message: %s - technical message: %s\n", err.Error(), failureObj.GetTechnicalMessage())
			return
		}

		log.Printf("Error: user message: %s - technical message: %s\n", err.Error(), failureObj.GetTechnicalMessage())
		switch failureObj.GetFailureType() {
		case failure.GeneralFailure:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":         failureObj.GetCode(),
				"user message": err.Error(),
			})
		case failure.NotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"code":         failureObj.GetCode(),
				"user message": err.Error(),
			})
		case failure.Validation:
			c.JSON(http.StatusBadRequest, gin.H{
				"code":         failureObj.GetCode(),
				"user message": err.Error(),
			})
		case failure.AccessUnauthorised:
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":         failureObj.GetCode(),
				"user message": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":         failureObj.GetCode(),
				"user message": err.Error(),
			})
			log.Printf("Error: unhandled internal code!: %v user message: %s - technical message: %s\n", failureObj.GetFailureType(), err.Error(), failureObj.GetTechnicalMessage())
		}
	}
}
