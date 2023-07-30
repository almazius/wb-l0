package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"os"
	"wb-l0/config"
	"wb-l0/internal/service"
	"wb-l0/utils"
)

type Postgres struct {
	Connection *sqlx.DB
	Log        zerolog.Logger
}

// NewPostgres Create entity service.Repository
func NewPostgres(config *config.Config) (service.Repository, error) {
	Log := zerolog.New(os.Stderr)
	conn, err := utils.GetConn(config)
	if err != nil {
		Log.Fatal().Timestamp().Err(err)
		return nil, &service.MyError{
			Message: err.Error(),
			Code:    500,
		}
	}

	return &Postgres{
		Connection: conn,
		Log:        Log,
	}, nil
}

// AddNote add note in repository
func (p *Postgres) AddNote(id string, jsonModel []byte) error {
	_, err := p.Connection.Exec(`insert into models values ($1, $2)`, id, jsonModel)
	if err != nil {
		p.Log.Error().Timestamp().Err(err).Send()
		return &service.MyError{
			Message: err.Error(),
			Code:    500,
		}
	}
	return nil
}

// GetAllNote return all notes from repository
func (p *Postgres) GetAllNote() (map[string][]byte, error) {
	models := make(map[string][]byte)
	objects, err := p.Connection.Queryx("select * from models")
	if err != nil {
		p.Log.Error().Timestamp().Str("Service", "Repository").Err(err).Send()
		return nil, &service.MyError{
			Message: err.Error(),
			Code:    500,
		}
	}

	for objects.Next() {
		var (
			id   string
			json []byte
		)
		err = objects.Scan(&id, &json)
		if err != nil {
			p.Log.Error().Timestamp().Str("Service", "Repository").Err(err).Send()
			return nil, &service.MyError{
				Message: err.Error(),
				Code:    500,
			}
		}
		models[id] = json
	}
	return models, nil
}
