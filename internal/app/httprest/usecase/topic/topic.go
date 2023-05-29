package topic

import (
	"be-idx-tsg/internal/app/httprest/model"
	tp "be-idx-tsg/internal/app/httprest/repository/topic"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAll(keyword string) ([]*model.Topic, error)
	GetByID(topicID, keyword string) (*model.Topic, error)
	UpdateHandler(topic model.UpdateTopicHandler, c *gin.Context) (int64, error)
	CreateTopicWithMessage(topic model.CreateTopicWithMessage, c *gin.Context) (int64, error)
	CreateMessage(message model.CreateMessage, c *gin.Context) (int64, error)
	DeleteTopic(topicID string, c *gin.Context) (int64, error)
	ArchiveTopicToFAQ(topic model.ArchiveTopicToFAQ, c *gin.Context) (int64, error)
}

type usecase struct {
	tpRepo tp.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		tp.NewRepository(),
	}
}

func (m *usecase) GetAll(keyword string) ([]*model.Topic, error) {
	return m.tpRepo.GetAll(keyword)
}

func (m *usecase) GetByID(topicID, keyword string) (*model.Topic, error) {
	return m.tpRepo.GetByID(topicID, keyword)
}

func (m *usecase) UpdateHandler(topic model.UpdateTopicHandler, c *gin.Context) (int64, error) {
	return m.tpRepo.UpdateHandler(topic, c)
}

func (m *usecase) CreateTopicWithMessage(topic model.CreateTopicWithMessage, c *gin.Context) (int64, error) {
	return m.tpRepo.CreateTopicWithMessage(topic, c)
}

func (m *usecase) CreateMessage(message model.CreateMessage, c *gin.Context) (int64, error) {
	return m.tpRepo.CreateMessage(message, c)
}

func (m *usecase) DeleteTopic(topicID string, c *gin.Context) (int64, error) {
	return m.tpRepo.DeleteTopic(topicID, c)
}

func (m *usecase) ArchiveTopicToFAQ(topic model.ArchiveTopicToFAQ, c *gin.Context) (int64, error) {
	return m.tpRepo.ArchiveTopicToFAQ(topic, c)
}
