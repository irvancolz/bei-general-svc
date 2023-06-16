package faq

import (
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/httprest/repository/faq"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAll(keyword, userId string) ([]*model.FAQ, error)
	CreateFAQ(faq model.CreateFAQ, c *gin.Context, isDraft bool) (int64, error)
	DeleteFAQ(faqID string, c *gin.Context) (int64, error)
	UpdateStatusFAQ(faq model.UpdateFAQStatus, c *gin.Context) (int64, error)
}

type usecase struct {
	faqRepo faq.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		faq.NewRepository(),
	}
}

func (m *usecase) GetAll(keyword, userId string) ([]*model.FAQ, error) {
	return m.faqRepo.GetAll(keyword, userId)
}

func (m *usecase) CreateFAQ(faq model.CreateFAQ, c *gin.Context, isDraft bool) (int64, error) {
	return m.faqRepo.CreateFAQ(faq, c, isDraft)
}

func (m *usecase) DeleteFAQ(faqID string, c *gin.Context) (int64, error) {
	return m.faqRepo.DeleteFAQ(faqID, c)
}

func (m *usecase) UpdateStatusFAQ(faq model.UpdateFAQStatus, c *gin.Context) (int64, error) {
	return m.faqRepo.UpdateStatusFAQ(faq, c)
}
