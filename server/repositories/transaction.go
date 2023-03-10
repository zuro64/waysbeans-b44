package repositories

import (
	"nis-waybeans/models"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindTransactions() ([]models.Transaction, error)
	GetTransaction(ID int) (models.Transaction, error)
	GetUncheckedOutTransactionByUserID(userID int) (models.Transaction, error)
	GetTransactionWithCart(userID int) (models.Transaction, error)
	CreateTransaction(transaction models.Transaction) (models.Transaction, error)
	UpdateTransaction(transaction models.Transaction) (models.Transaction, error)
	DeleteTransaction(transaction models.Transaction, ID int) (models.Transaction, error)
	GetUncheckedOutTransaction(UserID int) (models.Transaction, error)
	FindTransactionsByDate(userID int, startData time.Time, endDate time.Time) ([]models.Transaction, error)
	FindTransactionsByProductID(userID int, productID int) ([]models.Transaction, error)
}

func RepositoryTransaction(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("User").Preload("Cart.Product").Find(&transactions).Error
	return transactions, err
}

func (r *repository) GetTransaction(ID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").First(&transaction, ID).Error
	return transaction, err
}

func (r *repository) GetUncheckedOutTransactionByUserID(userID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Where("user_id = ? and status = 'Unchecked Out'", userID).Preload("User").First(&transaction).Error
	return transaction, err
}

func (r *repository) GetTransactionWithCart(userID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Where("user_id =?", userID).Preload("User").Preload("Cart.Product").First(&transaction).Error
	return transaction, err
}

func (r *repository) CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Create(&transaction).Error
	return transaction, err
}

func (r *repository) UpdateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Save(&transaction).Error
	return transaction, err
}

func (r *repository) DeleteTransaction(transaction models.Transaction, ID int) (models.Transaction, error) {
	err := r.db.Delete(&transaction).Error
	return transaction, err
}

func (r *repository) GetUncheckedOutTransaction(userID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.First(&transaction, "status = ? and user_id = ?", "Unchecked Out", userID).Error
	return transaction, err
}

func (r *repository) FindTransactionsByDate(userID int, startDate time.Time, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("user_id = ? AND CAST(updated_at AS DATE) BETWEEN ? AND ?", userID, startDate, endDate).Preload("User").Preload("Cart.Product").Find(&transactions).Error
	return transactions, err
}

func (r *repository) FindTransactionsByProductID(userID int, productID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	transactionsId := r.db.Select("transaction_id").Where("product_id = ?", productID).Table("carts")
	err := r.db.Where("user_id = ? and id in (?)", userID, transactionsId).Preload("User").Preload("Cart", "product_id = ?", productID).Find(&transactions).Error

	return transactions, err
}
