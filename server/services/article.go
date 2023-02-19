package services

import (
	"melodie-site/server/db"
	"melodie-site/server/models"
)

type ArticleService struct {
}

var articleService *ArticleService

func (service *ArticleService) GetAllArticles() ([]*models.Article, error) {
	conn := db.GetDBConn()
	articles := make([]*models.Article, 0)
	result := conn.Model(&models.Article{}).Preload("Tags").Find(&articles)
	if result.Error != nil {
		return nil, result.Error
	}
	return articles, nil
}

func (service *ArticleService) CreateArticle(article *models.Article) error {
	conn := db.GetDBConn()
	conn.Create(article)
	return nil
}

func (service *ArticleService) GetArticle(articleID int) (*models.Article, error) {
	conn := db.GetDBConn()
	article := &models.Article{}
	article.ID = articleID
	err := conn.Model(&models.Article{}).Preload("Tags").First(article).Error
	return article, err
}

func (service *ArticleService) UpdateArticle(article *models.Article) error {
	conn := db.GetDBConn()
	conn.Save(article)
	return nil
}

func GetArticleService() *ArticleService {
	if articleService == nil {
		articleService = &ArticleService{}
	}
	return articleService
}
