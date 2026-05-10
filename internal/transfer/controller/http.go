package controller

import (
	"net/http"

	"github.com/deshortone/ledger-system/internal/transfer/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	transferApplication domain.TransferApplication
}

func NewHandler(
	transferApplication domain.TransferApplication,
) Handler {
	return Handler{
		transferApplication: transferApplication,
	}
}

func (h *Handler) RegisterRoutes(c *gin.RouterGroup) {
	c.POST("/transfer", h.transferMoney)
}

// transferMoney godoc
//
//	@Summary		Transfers money from one account to another
//	@Description	Transfers money from one account to another a user
//	@Tags			transfer
//	@Accept       	json
//	@Produce		json
//	@Param			"Transfer money request"	body	TransferMoneyRequest	true	"Transfer money request desc"
//	@Success		200	{object}	string
//	@Failure		400	{object}	string
//	@Failure		500	{object}	string
//	@Router			/transfer		[post]
func (h *Handler) transferMoney(c *gin.Context) {
	var daRequest TransferMoneyRequest
	if err := c.BindJSON(&daRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ensure body is created correctly"})
		return
	}

	fromAccountId, err := uuid.Parse(daRequest.FromAccountId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ensure fromAccountId has been correctly setup"})
		return
	}

	toAccountId, err := uuid.Parse(daRequest.ToAccountId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ensure toAccountId has been correctly setup"})
		return
	}

	err = h.transferApplication.TransferMoney(c.Request.Context(), fromAccountId, toAccountId, daRequest.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown error", "Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transfer status": "success"})
}
