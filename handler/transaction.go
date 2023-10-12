package handler

import (
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHanlder(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
// parameter di uri
// tangkap paramter mapping input struct
// didalam handler panggil service, input sebagai parameter
// dalam service, berbekal campaign id bisa panggil repo
// repo mencari data transaction suatu campaign

	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign`s transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	transactions, err := h.service.GetTransactionsByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign`s transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign`s transactions", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
// GetUserTransactions 
// handler : ambil niali user dari jwt/middleware 
// service: repository ambil data trasactions (preload data campaign)
 currentUser := c.MustGet("currentUser").(user.User)
 userID := currentUser.ID

 transactions, err := h.service.GetTransactionsByUserID(userID)
 if err != nil {
	response := helper.APIResponse("Failed to get user`s transactions", http.StatusBadRequest, "error", nil)
	c.JSON(http.StatusBadRequest, response)
	return
}

response := helper.APIResponse("User`s transactions", http.StatusOK, "success",transaction.FormatUserTransactions(transactions))
c.JSON(http.StatusOK, response)

}