package usecase

import (
	"github.com/rs/zerolog"
	"os"
	"wb-l0/config"
	"wb-l0/internal/service"
	"wb-l0/internal/service/repository"
)

type Services struct {
	Cache      map[string][]byte
	repository service.Repository
	Log        zerolog.Logger
}

func NewServices(conf *config.Config) (service.IServices, error) {
	s := Services{}
	var err error
	s.Cache = make(map[string][]byte)
	s.Log = zerolog.New(os.Stderr)
	s.repository, err = repository.NewPostgres(conf)
	err = s.GetUpCache()
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
	var err error
	s.Cache, err = s.repository.GetAllNote()
	if err != nil {
		return err
	}

	return nil
}
