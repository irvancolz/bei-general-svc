package pkp

import (
	"be-idx-tsg/internal/app/httprest/model"
	pkp "be-idx-tsg/internal/app/httprest/repository/pkp"
	"time"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAllPKuser(c *gin.Context) ([]*model.PKuser, error)
	CreatePKuser(pkp model.CreatePKuser, c *gin.Context) (int64, error)
	UpdatePKuser(pkp model.UpdatePKuser, c *gin.Context) (int64, error)
	Delete(id string, c *gin.Context) (int64, error)
	GetAllWithFilter(keyword []string) ([]*model.PKuser, error)
	GetAllWithSearch(Code string, Name string, QuestionDate time.Time, Question string, Answers string, answered_by string, AnsweredAt time.Time) ([]*model.PKuser, error)
}

type usecase struct {
	pkpRepo pkp.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		pkp.NewRepository(),
	}
}

func (uc *usecase) GetAllPKuser(c *gin.Context) ([]*model.PKuser, error) {
	return uc.pkpRepo.GetAllPKuser(c)
}

func (uc *usecase) CreatePKuser(pkp model.CreatePKuser, c *gin.Context) (int64, error) {
	return uc.pkpRepo.CreatePKuser(pkp, c)
}

func (uc *usecase) UpdatePKuser(pkp model.UpdatePKuser, c *gin.Context) (int64, error) {
	return uc.pkpRepo.UpdatePKuser(pkp, c)
}

func (uc *usecase) Delete(id string, c *gin.Context) (int64, error) {
	return uc.pkpRepo.Delete(id, c)
}

func (uc *usecase) GetAllWithFilter(keyword []string) ([]*model.PKuser, error) {
	return uc.pkpRepo.GetAllWithFilter(keyword)
}

func (uc *usecase) GetAllWithSearch(Code string, Name string, QuestionDate time.Time, Question string, Answers string, answered_by string, AnsweredAt time.Time) ([]*model.PKuser, error) {
	return uc.pkpRepo.GetAllWithSearch(Code, Name, QuestionDate, Question, Answers, answered_by, AnsweredAt)
}
