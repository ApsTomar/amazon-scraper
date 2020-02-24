package models

import (
	"time"
)

type Product struct {
	ID           uint64 `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	DataAsin     string `json:"dataAsin"`
	ProductName  string `json:"productName"`
	CategoryID   uint64 `json:"category_id"`
	Manufacturer string `json:"manufacturer"`
	Price        string `json:"price"`
	Ratings      string `json:"ratings"`
	Description  string `json:"description"`
}

func (Product) TableName() string {
	return "product"
}

type Category struct {
	ID           uint64 `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	CategoryName string `json:"categoryName"`
}

func (Category) TableName() string {
	return "category"
}

type ViewProduct struct {
	DataAsin     string `json:"dataAsin"`
	ProductName  string `json:"productName"`
	Manufacturer string `json:"manufacturer"`
	Price        string `json:"price"`
	Ratings      string `json:"ratings"`
	Description  string `json:"description"`
}