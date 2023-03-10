package repositories

import (
	"fmt"
	"nis-waybeans/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	FindCarts() ([]models.Cart, error)
	FindUnCheckedOutCarts(transactionID int) ([]models.Cart, error)
	GetCart(ID int) (models.Cart, error)
	CreateCart(cart models.Cart) (models.Cart, error)
	UpdateCart(cart models.Cart) (models.Cart, error)
	DeleteCart(cart models.Cart, ID int) (models.Cart, error)
	FindCartsByTransactionID(transactionID int) ([]models.Cart, error)
	CheckCartProductID(productID int) (models.Cart, error)
}

func RepositoryCart(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindCarts() ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.Preload("Product").Find(&carts).Error
	return carts, err
}

func (r *repository) FindCartsByTransactionID(transactionID int) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.Where("transaction_id = ? and checkout=?", transactionID, true).Preload("Product").Find(&carts).Error
	return carts, err
}

func (r *repository) FindUnCheckedOutCarts(transactionID int) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.Where("checkout = ?", false).Preload("Product").Find(&carts).Error
	return carts, err
}

func (r *repository) GetCart(ID int) (models.Cart, error) {
	var cart models.Cart
	err := r.db.First(&cart, ID).Error
	return cart, err
}

func (r *repository) CheckCartProductID(productID int) (models.Cart, error) {
	var cart models.Cart
	err := r.db.Where("product_id = ? and checkout = ?", productID, false).First(&cart).Error
	return cart, err
}

func (r *repository) CreateCart(cart models.Cart) (models.Cart, error) {
	err := r.db.Create(&cart).Error
	return cart, err
}

func (r *repository) UpdateCart(cart models.Cart) (models.Cart, error) {
	err := r.db.Save(&cart).Error
	return cart, err
}

func (r *repository) DeleteCart(cart models.Cart, ID int) (models.Cart, error) {
	err := r.db.Delete(&cart, ID).Error
	return cart, err
}

func (r *repository) GetSummaryProductTransaction(TransactionID int) ([]models.CartProductSummaryTransaction, error) {
	var cartSummary []models.CartProductSummaryTransaction

	//err := r.db.Raw("SELECT product_id,sum(order_quantity) as order_quantity FROM carts WHERE transaction_id = 1 AND checkout = 1 group by product_id").Scan(&cartSummary).Error
	err := r.db.Table("carts").Select("product_id,sum(order_quantity) as order_quantity").Group("product_id").Scan(&cartSummary).Error
	fmt.Println("Summary repo", cartSummary)
	return cartSummary, err
}
