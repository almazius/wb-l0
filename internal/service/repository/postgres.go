package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"os"
	"time"
	"wb-l0/config"
	"wb-l0/internal/service"
	"wb-l0/utils"
)

type Postgres struct {
	Connection *sqlx.DB
	Log        *zerolog.Logger
}

func NewPostgres(config *config.Config) (service.Repository, error) {
	Log := zerolog.New(os.Stderr)
	Log.Info().Time("key", time.Now()).Msg("Postgres logg active")
	conn, err := utils.GetConn(config)
	if err != nil {
		Log.Fatal().Timestamp().Err(err)
		return nil, err
	}

	return &Postgres{
		Connection: conn,
		Log:        &Log,
	}, nil
}

func (p *Postgres) AddNote(id string, jmodel []byte) error {
	_, err := p.Connection.Exec(`insert into models values ($1, $2)`, id, jmodel)
	if err != nil {
		p.Log.Err(err).Timestamp().Send()
		return err
	}
	return nil
}

func (p *Postgres) GetAllNote() (map[string][]byte, error) {
	models := make(map[string][]byte)
	objects, err := p.Connection.Queryx("select * from models")
	if err != nil {
		p.Log.Err(err).Timestamp().Send()
		return nil, err
	}

	for objects.Next() {
		var (
			id   string
			json []byte
		)
		err = objects.Scan(&id, &json)
		if err != nil {
			p.Log.Error().Timestamp().Err(err)
			return nil, &service.MyError{
				Message: err.Error(),
				Code:    500,
			}
		}
		models[id] = json
	}
	return models, nil
}
