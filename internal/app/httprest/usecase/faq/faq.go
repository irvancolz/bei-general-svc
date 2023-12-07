package faq

import (
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/httprest/repository/faq"
	"be-idx-tsg/internal/app/utilities"
	"be-idx-tsg/internal/pkg/email"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAll(keyword, userId string) ([]*model.FAQ, error)
	CreateFAQ(faq model.CreateFAQ, c *gin.Context, isDraft bool) (int64, error)
	DeleteFAQ(faqID string, c *gin.Context) (int64, error)
	UpdateStatusFAQ(faq model.UpdateFAQStatus, c *gin.Context) (int64, error)
	UpdateFAQ(faq model.UpdateFAQ, c *gin.Context) (int64, error)
	UpdateOrderFAQ(faqs []model.UpdateFAQOrder, c *gin.Context) (int64, error)
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
	data, err := m.faqRepo.CreateFAQ(faq, c, isDraft)
	if err != nil {
		return 0, err
	}

	if !isDraft {
		emailSubject := "Aktivitas Baru Di Menu FAQ"
		notifType := "FAQ"
		notifMsg := fmt.Sprintf("User %s menambahkan FAQ baru", c.GetString("name_user"))

		utilities.CreateNotifForUserAng(c, notifType, notifMsg)
		utilities.CreateNotifForInternalBursa(c, notifType, notifMsg)
		email.SendEmailForUserAng(c, emailSubject, notifMsg)
		email.SendEmailForUserInternalBursa(c, emailSubject, notifMsg)
	}

	return data, nil
}

func (m *usecase) DeleteFAQ(faqID string, c *gin.Context) (int64, error) {
	data, err := m.faqRepo.DeleteFAQ(faqID, c)
	if err != nil {
		return 0, err
	}

	emailSubject := "Aktivitas Baru Di Menu FAQ"
	notifType := "FAQ"
	notifMsg := fmt.Sprintf("User %s menghapus FAQ", c.GetString("name_user"))

	utilities.CreateNotifForUserAng(c, notifType, notifMsg)
	utilities.CreateNotifForInternalBursa(c, notifType, notifMsg)
	email.SendEmailForUserAng(c, emailSubject, notifMsg)
	email.SendEmailForUserInternalBursa(c, emailSubject, notifMsg)

	return data, nil
}

func (m *usecase) UpdateStatusFAQ(faq model.UpdateFAQStatus, c *gin.Context) (int64, error) {
	data, err := m.faqRepo.UpdateStatusFAQ(faq, c)
	if err != nil {
		return 0, err
	}

	if faq.Status == model.PublishedFAQ {
		emailSubject := "Aktivitas Baru Di Menu FAQ"
		notifType := "FAQ"
		notifMsg := fmt.Sprintf("User %s menambahkan FAQ baru", c.GetString("name_user"))

		utilities.CreateNotifForUserAng(c, notifType, notifMsg)
		utilities.CreateNotifForInternalBursa(c, notifType, notifMsg)
		email.SendEmailForUserAng(c, emailSubject, notifMsg)
		email.SendEmailForUserInternalBursa(c, emailSubject, notifMsg)
	}

	return data, nil
}

func (m *usecase) UpdateFAQ(faq model.UpdateFAQ, c *gin.Context) (int64, error) {
	return m.faqRepo.UpdateFAQ(faq, c)
}

func (m *usecase) UpdateOrderFAQ(faqs []model.UpdateFAQOrder, c *gin.Context) (int64, error) {
	return m.faqRepo.UpdateOrderFAQ(faqs, c)
}
