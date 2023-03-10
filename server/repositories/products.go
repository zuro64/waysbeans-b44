package repositories

import (
	"nis-waybeans/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	FindProducts() ([]models.Product, error)
	GetProduct(ID int) (models.Product, error)
	CreateProduct(product models.Product) (models.Product, error)
	UpdateProduct(product models.Product) (models.Product, error)
	DeleteProduct(product models.Product, ID int) (models.Product, error)
	SearchProduct(name string) ([]models.Product, error)
	FindTopProducts() ([]models.Product, error)
}

func RepositoryProduct(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *repository) GetProduct(ID int) (models.Product, error) {
	var product models.Product
	err := r.db.First(&product, ID).Error
	return product, err
}

func (r *repository) FindTopProducts() ([]models.Product, error) {
	var products []models.Product
	//err := r.db.Where("checkout = ?", true).Preload("Product").Find(&carts).Error
	err := r.db.Raw("SELECT * FROM ( SELECT t2.*, sum(order_quantity) order_quantity FROM `carts` t1 JOIN `products` t2 on t1.product_id = t2.id WHERE checkout = 1 group by product_id UNION SELECT *,0 FROM PRODUCTS WHERE id NOT IN (SELECT product_id FROM CARTS WHERE checkout = 1)) T1 ORDER BY order_quantity DESC;").Scan(&products).Error
	return products, err
}

func (r *repository) SearchProduct(name string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("name LIKE ?", "%"+name+"%").Find(&products).Error
	return products, err
}

func (r *repository) CreateProduct(product models.Product) (models.Product, error) {
	err := r.db.Create(&product).Error
	return product, err
}

func (r *repository) UpdateProduct(product models.Product) (models.Product, error) {
	err := r.db.Save(&product).Error
	return product, err
}

func (r *repository) DeleteProduct(product models.Product, ID int) (models.Product, error) {
	err := r.db.Delete(&product, ID).Error
	return product, err
}
