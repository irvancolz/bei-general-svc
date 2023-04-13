package guidances

import (
	repo "be-idx-tsg/internal/app/httprest/repository/guidances"
)

type guidancesUsecase struct {
	Repository repo.GuidancesRepoInterface
}

type GuidancesRegulationAndFileUsecaseInterface interface {
	GuidancesUsecaseInterface
	RegulationUsecaseInterface
	FilesUsecaseInterface
}

func NewGuidanceUsecase() GuidancesRegulationAndFileUsecaseInterface {
	return &guidancesUsecase{
		Repository: repo.NewGuidancesRepository(),
	}
}
