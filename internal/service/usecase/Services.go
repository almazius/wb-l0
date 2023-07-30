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

// NewServices Create entity service.IServices
func NewServices(conf *config.Config) (service.IServices, error) {
	s := Services{}
	var err error
	s.Cache = make(map[string][]byte)
	s.Log = zerolog.New(os.Stderr)
	s.repository, err = repository.NewPostgres(conf)
	if err != nil {
		return nil, err
	}
	err = s.GetUpCache()
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// SaveModel save model on cache and repository
func (s *Services) SaveModel(id string, jsonModel []byte) error {
	_, exist := s.Cache[id]
	if exist {
		return &service.MyError{
			Message: "record with this key already exists",
			Code:    401,
		}
	}
	s.Cache[id] = jsonModel
	err := s.repository.AddNote(id, jsonModel)
	if err != nil {
		return err
	}

	return nil
}

// GetModel return model from cache
func (s *Services) GetModel(id string) ([]byte, error) {
	el, isExist := s.Cache[id]
	if !isExist {
		s.Log.Error().Timestamp().Str("Service", "Usecase").Msg("Getting an entry by invalid id")
		return nil, &service.MyError{
			Message: "Cant find your model",
			Code:    404,
		}
	}

	return el, nil
}

// GetUpCache restores cache
func (s *Services) GetUpCache() error {
	var err error
	s.Cache, err = s.repository.GetAllNote()
	if err != nil {
		return err
	}

	return nil
}
