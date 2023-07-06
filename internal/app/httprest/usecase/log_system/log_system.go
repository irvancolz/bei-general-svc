package log_system

import (
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/httprest/repository/log_system"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAll() ([]*model.LogSystem, error)
	CreateLogSystem(log model.CreateLogSystem, c *gin.Context) (int64, error)
}

type usecase struct {
	logSystemRepo log_system.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		log_system.NewRepository(),
	}
}

func (m *usecase) GetAll() ([]*model.LogSystem, error) {
	return m.logSystemRepo.GetAll()
}

func (m *usecase) CreateLogSystem(log model.CreateLogSystem, c *gin.Context) (int64, error) {
	return m.logSystemRepo.CreateLogSystem(log, c)
}
