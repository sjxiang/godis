package client

import (
	"context"
	"testing"

	"github.com/rs/zerolog/log"
)


func TestGet(t *testing.T) {
	client := New("localhost:5001")
	if err := client.Get(context.Background(), "foo"); err != nil {
		log.Fatal().Err(err).Msg("client send get command error")
	}

	t.Log("GET 指令，执行结束")
}


/*

	telnet 127.0.0.1 5001
	SET foo bar
	GET foo

	"*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	
 */
