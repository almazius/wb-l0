package utils

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"time"
	"wb-l0/config"
	"wb-l0/internal/service"
)

// TryConnectNats try connection on NATS
func TryConnectNats(conf *config.Config, count int, Log *zerolog.Logger) (nc *nats.Conn, err error) {
	for i := 1; i < count+1; i++ {
		nc, err = nats.Connect(conf.Stan.Url)
		if err != nil {
			Log.Error().Err(err).Timestamp().Msgf("Failed connection to NATS. Attempt %d", i)
		} else {
			Log.Info().Timestamp().Msg("Success connection to NATS")
			return nc, err
		}
		time.Sleep(1 * time.Second)
	}
	return nil, &service.MyError{Code: 502, Message: err.Error()}
}
