package main

import (
	"gorm.io/gorm"
)

type product struct {
	ID    uint    `gorm:"primary_key" json:"id" uri:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p *product) getProduct(db *gorm.DB) error {
	if err := db.First(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *product) updateProduct(db *gorm.DB) error {
	if err := db.Save(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *product) deleteProduct(db *gorm.DB) error {
	if err := db.Delete(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *product) createProduct(db *gorm.DB) error {
	if err := db.Save(p).Error; err != nil {
		return err
	}
	return nil
}

func getProducts(db *gorm.DB, start, count int) ([]product, error) {
	var p []product
	if err := db.Offset(start).Limit(count).Find(&p).Error; err != nil {
		return nil, err
	}
	return p, nil
}
