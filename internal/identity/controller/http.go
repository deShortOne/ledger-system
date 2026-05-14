package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/deshortone/ledger-system/internal/identity/domain"
	"github.com/deshortone/ledger-system/pkg/failure"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	accountService domain.AccountService
	accountCreator domain.AccountCreator
	userService    domain.UserService
}

func NewHandler(
	accountService domain.AccountService,
	accountCreator domain.AccountCreator,
	userService domain.UserService,
) Handler {
	return Handler{
		accountService: accountService,
		accountCreator: accountCreator,
		userService:    userService,
	}
}

func (h *Handler) RegisterRoutes(c *gin.RouterGroup) {
	c.POST("/create-user", h.createUser)
	c.POST("/:userId/create-account", h.createAccount)
	c.GET("/:userId/accounts", h.getAccountsForUser)
}

// createUser godoc
//
//	@Summary		Creates a user
//	@Description	Creates a user
//	@Tags			identity
//	@Produce		json
//	@Success		202	{object}	string
//	@Failure		400	{object}	string
//	@Failure		500	{object}	string
//	@Router			/create-user	[post]
func (h *Handler) createUser(c *gin.Context) {
	var daRequest NewUserRequst
	if err := c.BindJSON(&daRequest); err != nil {
		c.Errors = append(c.Errors, c.Error(failure.NewFailure(failure.ConversionError, failure.Validation, errors.New("ensure body is created correctly"), "")))
		return
	}

	user, err := h.userService.CreateNewUser(c.Request.Context(), daRequest.FirstName, daRequest.LastName)
	if err != nil {
		c.Errors = append(c.Errors, c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": user.Identifier})
}

// createAccount godoc
//
//	@Summary		Creates an account attached to a user
//	@Description	Creates an account attached to a user
//	@Tags			identity
//	@Produce		json
//	@Success		202
//	@Failure		400	{object}			string
//	@Failure		500	{object}			string
//	@Router			/:userId/create-account [post]
func (h *Handler) createAccount(c *gin.Context) {
	userIdString := c.Param("userId")
	if userIdString == "" {
		c.Errors = append(c.Errors, c.Error(failure.NewFailure(failure.ConversionError, failure.Validation, errors.New("user id is missing"), "")))
		return
	}
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		c.Errors = append(c.Errors, c.Error(failure.NewFailure(failure.ConversionError, failure.Validation, errors.New("user id is incorrect"), err.Error())))
		return
	}

	var daRequest NewAccountRequest
	if err := c.BindJSON(&daRequest); err != nil {
		c.Errors = append(c.Errors, c.Error(failure.NewFailure(failure.ConversionError, failure.Validation, errors.New("ensure body is created correctly"), err.Error())))
		return
	}

	if _, err := h.accountCreator.AddAccountToUser(c.Request.Context(), userId, daRequest.AccountType, daRequest.Currency); err != nil {
		c.Errors = append(c.Errors, c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, nil)
}

// getAccountsForUser godoc
//
//	@Summary		Creates an account attached to a user
//	@Description	Creates an account attached to a user
//	@Tags			identity
//	@Produce		json
//	@Success		200 {object}		GetAccountsForUserResponse
//	@Failure		400	{object}		string
//	@Failure		500	{object}		string
//	@Router			/:userId/accounts 	[get]
func (h *Handler) getAccountsForUser(c *gin.Context) {
	userIdString := c.Param("userId")
	if userIdString == "" {
		c.Errors = append(c.Errors, c.Error(failure.NewFailure(failure.ConversionError, failure.Validation, errors.New("user id is missing"), "")))
		return
	}
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		c.Errors = append(c.Errors, c.Error(failure.NewFailure(failure.ConversionError, failure.Validation, errors.New("user id is incorrect"), err.Error())))
		return
	}

	accounts, err := h.accountService.GetAccountsOwnedByUser(c.Request.Context(), userId)
	if err != nil {
		c.Errors = append(c.Errors, c.Error(err))
		return
	}

	accountsToReturn := make([]AccountResponse, 0, len(accounts))
	for _, account := range accounts {
		accountsToReturn = append(accountsToReturn, AccountResponse{
			Identifier:  account.Identifier.String(),
			AccountType: account.AccountType,
			Currency:    account.Currency,
			Status:      account.Status,
			DateCreated: account.CreatedAt.Time.Format(time.DateOnly),
		})
	}

	c.JSON(http.StatusCreated, GetAccountsForUserResponse{
		Accounts: accountsToReturn,
	})
}
