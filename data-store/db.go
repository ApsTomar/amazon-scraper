package data_store

import (
	"github.com/amazon/migrations"
	"github.com/amazon/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"strings"
)

type Storage struct {
	Db *gorm.DB
}

type DataStore interface {
	UpsertProduct(models.Product) error
	AddCategory(models.Category) (uint64, error)
}

func DbConnect() DataStore {
	db, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/amazon?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalf("DB connection not established due to: %v", err)
	}
	err = migrations.InitMySQL(db)
	if err != nil {
		log.Fatalf("error running migrations: %v", err)
	}
	return &Storage{Db: db}
}

func (store *Storage) UpsertProduct(product models.Product) error {
	err := store.Db.Create(&product).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			prod := &models.Product{}
			err := store.Db.First(prod).Where("data_asin=?", product.DataAsin).Error
			if err != nil {
				return err
			}
			prod.CategoryID = product.CategoryID
			prod.Description = product.Description
			prod.Ratings = product.Ratings
			prod.Price = product.Price
			prod.ProductName = product.ProductName
			prod.Manufacturer = product.Manufacturer

			err = store.Db.Save(prod).Error
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (store *Storage) AddCategory(category models.Category) (uint64, error) {
	cat := &models.Category{}
	err := store.Db.Where("category_name=?", category.CategoryName).Find(&cat).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = store.Db.Create(&category).Error
			if err != nil {
				return 0, err
			} else {
				return category.ID, nil
			}
		}
		return 0, err
	}
	return cat.ID, nil
}
