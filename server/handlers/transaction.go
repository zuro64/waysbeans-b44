package handlers

import (
	"fmt"
	"net/http"
	dto "nis-waybeans/dto/result"
	transactiondto "nis-waybeans/dto/transaction"
	"nis-waybeans/models"
	"nis-waybeans/repositories"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
	CartRepository        repositories.CartRepository
	ProductRepository     repositories.ProductRepository
}

func HandlerTransaction(
	TransactionRepository repositories.TransactionRepository,
	CartRepository repositories.CartRepository,
	ProductRepository repositories.ProductRepository) *handlerTransaction {

	return &handlerTransaction{
		TransactionRepository,
		CartRepository,
		ProductRepository,
	}
}
func convertResponseTransaction(u models.Transaction) transactiondto.TransactionResponse {
	return transactiondto.TransactionResponse{
		ID:         u.ID,
		User:       u.User,
		Name:       u.Name,
		Email:      u.Email,
		Phone:      u.Phone,
		Address:    u.Address,
		Attachment: u.Attachment,
		Status:     u.Status,
		Cart:       u.Cart,
	}

}

func convertResponseTransactionUnfinished(u models.Transaction) transactiondto.TransactionResponse {
	return transactiondto.TransactionResponse{
		SubTotal: u.SubTotal,
		TotalQty: u.TotalQty,
	}
}

func (h *handlerTransaction) FindTransactions(c echo.Context) error {
	transactions, err := h.TransactionRepository.FindTransactions()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: transactions,
	})
}

func (h *handlerTransaction) FindTransactionsByDate(c echo.Context) error {
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	userLogin := c.Get("userLogin")
	userID := userLogin.(jwt.MapClaims)["id"].(float64)
	dateFormat := "2006-01-02"
	startDateFormatted, _ := time.Parse(dateFormat, startDate)
	endDateFormatted, _ := time.Parse(dateFormat, endDate)
	fmt.Println("Start date", startDateFormatted)
	fmt.Println("End date", endDateFormatted)
	transactions, err := h.TransactionRepository.FindTransactionsByDate(int(userID), startDateFormatted, endDateFormatted)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: transactions,
	})

}

func (h *handlerTransaction) FindTransactionsByProductID(c echo.Context) error {
	productID := c.QueryParam("product_id")

	userLogin := c.Get("userLogin")
	userID := userLogin.(jwt.MapClaims)["id"].(float64)
	productIDConv, _ := strconv.Atoi(productID)
	transactions, err := h.TransactionRepository.FindTransactionsByProductID(int(userID), productIDConv)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: transactions,
	})

}

func (h *handlerTransaction) GetTransaction(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	transaction, err := h.TransactionRepository.GetTransactionWithCart(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertResponseTransaction(transaction),
	})
}

func (h *handlerTransaction) GetUncheckedOutTransaction(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userID := userLogin.(jwt.MapClaims)["id"].(float64)
	transaction, err := h.TransactionRepository.GetUncheckedOutTransactionByUserID(int(userID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertResponseTransactionUnfinished(transaction),
	})
}

func (h *handlerTransaction) CreateTransaction(c echo.Context) error {
	request := new(transactiondto.CreateTransactionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	userLogin := c.Get("userLogin")
	userID := userLogin.(jwt.MapClaims)["id"].(float64)
	var transaction models.Transaction

	transaction, _ = h.TransactionRepository.GetUncheckedOutTransaction(int(userID))

	var transactionIsMatch = false
	var transactionId int
	for !transactionIsMatch {
		transactionId = int(time.Now().Unix()) //1678180770
		transactionData, _ := h.TransactionRepository.GetTransaction(transactionId)
		if transactionData.ID == 0 {
			transactionIsMatch = true
		}
	}
	if transaction.ID == 0 {
		transaction = models.Transaction{
			UserID:     int(userID),
			ID:         transactionId,
			Name:       request.Name,
			Email:      request.Email,
			Phone:      request.Phone,
			PostCode:   request.PostCode,
			Address:    request.Address,
			Attachment: request.Attachment,
			Status:     request.Status,
		}

		data, err := h.TransactionRepository.CreateTransaction(transaction)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResult{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusOK, dto.SuccessResult{
			Code: http.StatusOK,
			Data: convertResponseTransaction(data),
		})
	}
	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertResponseTransaction(transaction),
	})
}

func (h *handlerTransaction) UpdateTransaction(c echo.Context) error {
	dataFile := c.Get("dataFile").(string)
	fmt.Println("THIS IS DATA FILE !", dataFile)
	userLogin := c.Get("userLogin")
	userID := userLogin.(jwt.MapClaims)["id"].(float64)
	request := transactiondto.CreateTransactionRequest{

		Name:       c.FormValue("name"),
		Email:      c.FormValue("email"),
		Phone:      c.FormValue("phone"),
		PostCode:   c.FormValue("post_code"),
		Address:    c.FormValue("address"),
		Attachment: dataFile,
		Status:     "Waiting For Verification",
	}

	transaction, err := h.TransactionRepository.GetUncheckedOutTransactionByUserID(int(userID))

	if transaction.ID == 0 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: "User ini belum memiliki transaksi",
		})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	h.CheckOutCart(int(userID))

	if request.Name != "" {
		transaction.Name = request.Name
	}

	if request.Email != "" {
		transaction.Email = request.Email
	}

	if request.Phone != "" {
		transaction.Phone = request.Phone
	}

	if request.PostCode != "" {
		transaction.PostCode = request.PostCode
	}

	if request.Address != "" {
		transaction.Address = request.Address
	}

	if request.Attachment != "" {
		transaction.Attachment = request.Attachment
	}

	if request.Status != "" {
		transaction.Status = request.Status
	}

	_, err = h.TransactionRepository.UpdateTransaction(transaction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	data, _ := h.TransactionRepository.GetTransactionWithCart(int(userID))
	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: data,
	})
}

func (h *handlerTransaction) DeleteTransaction(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	transaction, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	data, err := h.TransactionRepository.DeleteTransaction(transaction, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertResponseTransaction(data),
	})
}

func (h *handlerTransaction) CheckOutCart(userID int) {
	var transaction models.Transaction
	transaction, _ = h.TransactionRepository.GetUncheckedOutTransaction(int(userID))
	transactionID := transaction.ID

	carts, _ := h.CartRepository.FindUnCheckedOutCarts(transactionID)

	for _, element := range carts {
		id := element.ID
		if !element.Checkout {
			cart, _ := h.CartRepository.GetCart(id)
			cart.Checkout = true
			h.CartRepository.UpdateCart(cart)
		}

	}
	h.UpdateStockProduct(transaction.ID)

}

func (h *handlerTransaction) UpdateStockProduct(TransactionID int) {
	carts, err := h.CartRepository.FindCartsByTransactionID(TransactionID)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, element := range carts {
		product, _ := h.ProductRepository.GetProduct(element.ProductID)
		product.Stock = product.Stock - element.OrderQuantity
		h.ProductRepository.UpdateProduct(product)
	}

}
