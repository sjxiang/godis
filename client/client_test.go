package client

import (
	"context"
	"testing"

	"github.com/rs/zerolog/log"
)

func TestSet(t *testing.T) {
	client := New("localhost:5001")
	if err := client.Set(context.Background(), "foo", "bar"); err != nil {
		log.Fatal().Err(err).Msg("client send set command error")
	}

	t.Log("SET 指令，执行结束")
}