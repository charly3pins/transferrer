package transferrer

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Transfer transport the required information to make a transfer between two users
type Transfer struct {
	OriginUser        string  `json:"originUser" binding:"required"`
	OriginNumber      string  `json:"originNumber" binding:"required"`
	DestinationUser   string  `json:"destinationUser" binding:"required"`
	DestinationNumber string  `json:"destinationNumber" binding:"required"`
	Amount            float64 `json:"amount" binding:"required"`
}

// Handler is the object that handle the request from the API
type Handler struct {
	Store Store
}

// Balance is the handler method that returns the balance of a given user
func (h Handler) Balance(c *gin.Context) {
	user, ok := c.Get(emailContextKey)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	account, err := h.Store.Account(user.(string))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": account.Balance, "currency": account.Currency})
}

// Transfer is the handler method that makes a transfer of money between two users
func (h Handler) Transfer(c *gin.Context) {
	email, ok := c.Get(emailContextKey)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	transfer := Transfer{}
	if err := c.BindJSON(&transfer); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	account, err := h.Store.Account(email.(string))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if account.Balance < transfer.Amount {
		c.JSON(http.StatusInternalServerError, "Balance insufficient for this operation.")
		return
	}

	err = h.Store.Move(transfer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
