package database

import (
	"github.com/jinzhu/gorm"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreateUpload will create an upload
func (db *Database) CreateUpload(source, version, maintainer, changedBy string) (*models.Upload, error) {
	upload := &models.Upload{
		Source:     source,
		Version:    version,
		Maintainer: maintainer,
		ChangedBy:  changedBy,
	}

	if err := db.gormDB.Create(upload).Error; err != nil {
		return nil, err
	}

	return upload, nil
}

// GetAllUploads returns all uploads
func (db *Database) GetAllUploads() ([]*models.Upload, error) {
	var uploads []*models.Upload

	if err := db.gormDB.Model(&models.Upload{}).Find(&uploads).Error; err != nil {
		return nil, err
	}

	return uploads, nil
}

// GetUpload returns the Upload with the given id
func (db *Database) GetUpload(id uint) (*models.Upload, error) {
	var upload models.Upload

	query := db.gormDB.Where(
		&models.Upload{
			ID: id,
		},
	)

	err := query.First(&upload).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &upload, nil
}
