package controller

import (
	"errors"
	"net/http"

	"github.com/deshortone/ledger-system/internal/transfer/domain"
	"github.com/deshortone/ledger-system/pkg/failure"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	transferApplication     domain.TransferApplication
	transferReadonlyService domain.TransferReadOnlyService
}

func NewHandler(
	transferApplication domain.TransferApplication,
	transferReadonlyService domain.TransferReadOnlyService,
) Handler {
	return Handler{
		transferApplication:     transferApplication,
		transferReadonlyService: transferReadonlyService,
	}
}

func (h *Handler) RegisterRoutes(c *gin.RouterGroup) {
	c.POST("/deposit", h.depositMoney)
	c.POST("/getTransferStatus", h.getTransferStatus)
	c.POST("/transfer", h.transferMoney)
}

// depositMoney godoc
//
//	@Summary		Deposits money into an account
//	@Description	Deposits money into an account
//	@Tags			transfer
//	@Accept       	json
//	@Produce		json
//	@Param			"Deposit money request"	body	DepositMoneyRequest		true	"Deposits money request desc"
//	@Success		200	{object}	string
//	@Failure		400	{object}	string
//	@Failure		500	{object}	string
//	@Router			/deposit		[post]
func (h *Handler) depositMoney(c *gin.Context) {
	var daRequest DepositMoneyRequest
	if err := c.BindJSON(&daRequest); err != nil {
		c.Errors = append(c.Errors, c.Error(failure.NewFailure(failure.ConversionError, failure.Validation, errors.New("ensure body is created correctly"), "")))
		return
	}

	toAccountId, err := uuid.Parse(daRequest.ToAccountId)
	if err != nil {
		c.Errors = append(c.Errors, c.Error(failure.NewFailure(failure.ConversionError, failure.Validation, errors.New("ensure toAccountId has been correctly setup"), "")))
		return
	}

	transferId, err := h.transferApplication.TransferMoney(c.Request.Context(), uuid.MustParse("6724081d-6f50-4172-92c1-9c5d571f051c"), toAccountId, daRequest.Amount)
	if err != nil {
		c.Errors = append(c.Errors, c.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"transfer status": "success", "transfer id": transferId.String()})
}

// getTransferStatus godoc
//
//	@Summary		Get status about transfer
//	@Description	Get status about transfer
//	@Tags			transfer
//	@Accept       	json
//	@Produce		json
//	@Param			"Transfer money status request"	body	TransferMoneyStatusRequest		true	"Transfer money status request desc"
//	@Success		200	{object}	string
//	@Failure		400	{object}	string
//	@Failure		500	{object}	string
//	@Router			/getTransferStatus		[post]
func (h *Handler) getTransferStatus(c *gin.Context) {
	var daRequest TransferMoneyStatusRequest
	if err := c.BindJSON(&daRequest); err != nil {
		c.Errors = append(c.Errors, c.Error(failure.NewFailure(failure.ConversionError, failure.Validation, errors.New("ensure body is created correctly"), "")))
		return
	}

	transferId, err := uuid.Parse(daRequest.TransferId)
	if err != nil {
		c.Errors = append(c.Errors, c.Error(failure.NewFailure(failure.ConversionError, failure.Validation, errors.New("ensure transferId has been correctly setup"), "")))
		return
	}

	transferStatus, err := h.transferReadonlyService.GetTransferStatus(c.Request.Context(), transferId)
	if err != nil {
		c.Errors = append(c.Errors, c.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"isTransferSuccessful": transferStatus.IsSuccessful, "reason": transferStatus.ReasonForNotSuccess})
}

// transferMoney godoc
//
//	@Summary		Transfers money from one account to another
//	@Description	Transfers money from one account to another
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
		c.Errors = append(c.Errors, c.Error(failure.NewFailure(failure.ConversionError, failure.Validation, errors.New("ensure body is created correctly"), "")))
		return
	}

	fromAccountId, err := uuid.Parse(daRequest.FromAccountId)
	if err != nil {
		c.Errors = append(c.Errors, c.Error(failure.NewFailure(failure.ConversionError, failure.Validation, errors.New("ensure fromAccountId has been correctly setup"), "")))
		return
	}

	toAccountId, err := uuid.Parse(daRequest.ToAccountId)
	if err != nil {
		c.Errors = append(c.Errors, c.Error(failure.NewFailure(failure.ConversionError, failure.Validation, errors.New("ensure toAccountId has been correctly setup"), "")))
		return
	}

	transferId, err := h.transferApplication.TransferMoney(c.Request.Context(), fromAccountId, toAccountId, daRequest.Amount)
	if err != nil {
		c.Errors = append(c.Errors, c.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"transfer status": "success", "transfer id": transferId})
}
