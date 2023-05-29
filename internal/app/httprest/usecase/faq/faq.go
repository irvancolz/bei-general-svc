package faq

import (
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/httprest/repository/faq"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAll(keyword string) ([]*model.FAQ, error)
	CreateFAQ(faq model.CreateFAQ, c *gin.Context) (int64, error)
	DeleteFAQ(faqID string, c *gin.Context) (int64, error)
}

type usecase struct {
	faqRepo faq.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		faq.NewRepository(),
	}
}

func (m *usecase) GetAll(keyword string) ([]*model.FAQ, error) {
	return m.faqRepo.GetAll(keyword)
}

func (m *usecase) CreateFAQ(faq model.CreateFAQ, c *gin.Context) (int64, error) {
	return m.faqRepo.CreateFAQ(faq, c)
}

func (m *usecase) DeleteFAQ(faqID string, c *gin.Context) (int64, error) {
	return m.faqRepo.DeleteFAQ(faqID, c)
}
