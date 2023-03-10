package handlers

import (
	"fmt"
	"net/http"
	cartdto "nis-waybeans/dto/cart"
	dto "nis-waybeans/dto/result"
	"nis-waybeans/models"
	"nis-waybeans/repositories"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerCart struct {
	CartRepository        repositories.CartRepository
	TransactionRepository repositories.TransactionRepository
	ProductRepository     repositories.ProductRepository
}

func HandlerCart(
	CartRepository repositories.CartRepository,
	TransactionRepository repositories.TransactionRepository,
	ProductRepository repositories.ProductRepository,
) *handlerCart {
	return &handlerCart{
		CartRepository,
		TransactionRepository,
		ProductRepository,
	}
}

func convertResponseCart(u models.Cart) cartdto.CartResponse {
	return cartdto.CartResponse{
		ID:            u.ID,
		ProductID:     u.ProductID,
		OrderQuantity: u.OrderQuantity,
		Checkout:      u.Checkout,
	}
}

func (h *handlerCart) FindCarts(c echo.Context) error {
	carts, err := h.CartRepository.FindCarts()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: carts,
	})
}

func (h *handlerCart) FindUnCheckedOutCarts(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userID := userLogin.(jwt.MapClaims)["id"].(float64)
	var transaction models.Transaction
	transaction, _ = h.TransactionRepository.GetUncheckedOutTransaction(int(userID))
	transactionID := transaction.ID

	carts, err := h.CartRepository.FindUnCheckedOutCarts(transactionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: carts,
	})
}

func (h *handlerCart) GetCart(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	cart, err := h.CartRepository.GetCart(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertResponseCart(cart),
	})
}

func (h *handlerCart) CreateCart(c echo.Context) error {

	userLogin := c.Get("userLogin")
	userID := userLogin.(jwt.MapClaims)["id"].(float64)
	userRole := userLogin.(jwt.MapClaims)["role"].(string)
	if userRole == "admin" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: "User anda tidak bisa membuat order.",
		})
	}

	request := new(cartdto.CreateCartRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: "Error bind:" + err.Error(),
		})
	}
	product, err := h.ProductRepository.GetProduct(request.ProductID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: "Error product: " + err.Error(),
		})
	}

	if product.Stock == 0 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: "Maaf, stock barang habis.",
		})
	}

	if product.Stock < request.OrderQuantity {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: "Maaf, sisa stock barang sekarang adalah " + strconv.Itoa(product.Stock) + " Pcs, tidak mencukupi untuk order anda.",
		})
	}

	var transactionIsMatch = false
	var transactionId int
	for !transactionIsMatch {
		transactionId = int(time.Now().Unix()) //1678180770
		transactionData, _ := h.TransactionRepository.GetTransaction(transactionId)
		if transactionData.ID == 0 {
			transactionIsMatch = true
		}
	}
	transaction := models.Transaction{}
	transaction, _ = h.TransactionRepository.GetUncheckedOutTransaction(int(userID))
	fmt.Println(transaction.ID)
	if transaction.ID == 0 {
		transaction = models.Transaction{
			ID:     transactionId,
			UserID: int(userID),
			Status: "Unchecked Out",
		}
		transaction, _ = h.TransactionRepository.CreateTransaction(transaction)
	}

	request.TransactionID = transaction.ID

	validation := validator.New()
	err = validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	cart := models.Cart{
		ProductID:     request.ProductID,
		OrderQuantity: request.OrderQuantity,
		Checkout:      false,
		TransactionID: request.TransactionID,
	}

	checkCart, _ := h.CartRepository.CheckCartProductID(request.ProductID)
	fmt.Println(checkCart.ID)
	if checkCart.ID != 0 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: "Product ini sudah masuk di cart sebelumnya",
		})
	}

	data, err := h.CartRepository.CreateCart(cart)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	carts, err := h.CartRepository.FindUnCheckedOutCarts(transaction.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	var totalQty = 0
	var subTotal = 0
	for _, element := range carts {
		if !element.Checkout {
			totalQty = element.OrderQuantity + totalQty
			product, _ := h.ProductRepository.GetProduct(element.ProductID)
			subTotal = subTotal + (product.Price * element.OrderQuantity)
		}

	}

	h.UpdateTransactionTotal(transaction, subTotal, totalQty)

	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertResponseCart(data),
	})
}

func (h *handlerCart) UpdateCart(c echo.Context) error {
	request := new(cartdto.UpdateCartRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	id, _ := strconv.Atoi(c.Param("id"))
	cart, err := h.CartRepository.GetCart(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if request.OrderQuantity != 0 {
		cart.OrderQuantity = request.OrderQuantity
	}

	if request.Checkout {
		cart.Checkout = request.Checkout
	}

	data, err := h.CartRepository.UpdateCart(cart)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	userLogin := c.Get("userLogin")
	userID := userLogin.(jwt.MapClaims)["id"].(float64)
	var transaction models.Transaction
	transaction, _ = h.TransactionRepository.GetUncheckedOutTransaction(int(userID))
	transactionID := transaction.ID
	carts, err := h.CartRepository.FindUnCheckedOutCarts(transactionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	var totalQty = 0
	var subTotal = 0
	for _, element := range carts {
		if !element.Checkout {
			totalQty = element.OrderQuantity + totalQty
			product, _ := h.ProductRepository.GetProduct(element.ProductID)
			subTotal = subTotal + (product.Price * element.OrderQuantity)
		}

	}

	h.UpdateTransactionTotal(transaction, subTotal, totalQty)
	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertResponseCart(data),
	})
}

func (h *handlerCart) DeleteCart(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	cart, err := h.CartRepository.GetCart(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	data, err := h.CartRepository.DeleteCart(cart, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertResponseCart(data),
	})
}

func (h *handlerCart) UpdateTransactionTotal(transaction models.Transaction, subTotal int, totalQty int) {
	if transaction.TotalQty != totalQty {
		transaction.TotalQty = totalQty
	}

	if transaction.SubTotal != subTotal {
		transaction.SubTotal = subTotal
	}
	h.TransactionRepository.UpdateTransaction(transaction)
}
