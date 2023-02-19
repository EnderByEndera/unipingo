package services

import (
	"melodie-site/server/db"
	"melodie-site/server/models"
)

type TagsService struct {
}

var tagsService *TagsService

func (service *TagsService) CreateTag(tag *models.Tag) error {
	conn := db.GetDBConn()
	result := conn.Create(tag)
	return result.Error
}

func (service *TagsService) GetTag(tagID int) (*models.Tag, error) {
	conn := db.GetDBConn()
	tag := &models.Tag{}
	tag.ID = tagID
	conn.First(tag)
	return tag, nil
}

func (service *TagsService) UpdateTag(tag *models.Tag) error {
	conn := db.GetDBConn()
	conn.Save(tag)
	return nil
}

func (service *TagsService) GetAllTags() ([]*models.Tag, error) {
	conn := db.GetDBConn()
	tags := []*models.Tag{}
	err := conn.Find(&tags).Error
	return tags, err
}

func GetTagsService() *TagsService {
	if tagsService == nil {
		tagsService = &TagsService{}
	}
	return tagsService
}
