package usecase

import (
	"github.com/rs/zerolog"
	"os"
	"wb-l0/internal/service"
)

type Services struct {
	Cache      map[string][]byte
	repository service.Repository
	Log        zerolog.Logger
}

func NewServices() (service.IServices, error) {
	s := Services{}
	s.Log = zerolog.New(os.Stderr)
	err := s.GetUpCache()
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Services) SaveModel(id string, jmodel []byte) error {
	s.Cache[id] = jmodel
	err := s.repository.AddNote(id, jmodel)
	if err != nil {
		return err
	}

	return nil
}

func (s *Services) GetModel(id string) ([]byte, error) {
	el, isExist := s.Cache[id]
	if !isExist {
		s.Log.Error().Timestamp().Msg("Перевести: Получение записи по невалидному id")
		return nil, &service.MyError{
			Message: "Cant find your model",
			Code:    404,
		}
	}

	return el, nil
}

func (s *Services) GetUpCache() error {

}
