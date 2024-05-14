package usecase

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/sight"
)

type SightUseCaseInterface interface {
	GetSightByID(sightID int) (entities.Sight, error)
	GetCommentsBySightID(commentID int) ([]entities.Comment, error)
	GetCommentsByUserID(userID int) ([]entities.Comment, error)
	GetSightsList() ([]entities.Sight, error)
	SearchSights(str string) (entities.Sights, error)
}

type SightUseCase struct {
	SightStorage *sight.SightStorage
}

func NewSightUseCase(storage *sight.SightStorage) *SightUseCase {
	return &SightUseCase{
		SightStorage: storage,
	}
}

func (su *SightUseCase) GetSightByID(sightID int) (entities.Sight, error) {
	return su.SightStorage.GetSight(sightID)
}

func (su *SightUseCase) GetCommentsBySightID(commentID int) ([]entities.Comment, error) {
	return su.SightStorage.GetCommentsBySightID(commentID)
}

func (su *SightUseCase) GetCommentsByUserID(userID int) ([]entities.Comment, error) {
	return su.SightStorage.GetCommentsByUserID(userID)
}

func (su *SightUseCase) GetSightsList() ([]entities.Sight, error) {
	return su.SightStorage.GetSightsList()
}

func (su *SightUseCase) SearchSights(str string) (entities.Sights, error) {
	return su.SightStorage.SearchSights(str)
}
