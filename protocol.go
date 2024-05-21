package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/resp"
)

const (
	CommandSet = "SET"
	CommandGet = "GET"

)

type Command interface {
	Execute()
}

type SetCommand struct {
	key, value string
}
func (c *SetCommand) Execute() {

}

type GetCommand struct {
	key string
}
func (c *GetCommand) Execute() {

}



func parseCommand(raw string) (Command, error) {
	
	rd := resp.NewReader(bytes.NewBufferString(raw))
	
	for {
		v, _, err := rd.ReadValue() 
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("resp read value error")
		}

		if v.Type() == resp.Array {
			
			for _, val := range v.Array() {
			
				switch val.String() {
				// SET
				case CommandSet:  
					return extractSetCmd(v) 
				// GET
				case CommandGet:
					return extractGetCmd(v)
				// ??

				}

			} 
		}	
	}

	return nil, fmt.Errorf("invalid or unknown command received: %s", raw)
}

func extractSetCmd(value resp.Value) (*SetCommand, error) {

	if len(value.Array()) != 3 {
		return nil, fmt.Errorf("invalid number of variables for SET command")
	}
	cmd := SetCommand{
		key:   value.Array()[1].String(),
		value: value.Array()[2].String(),
	}

	return &cmd, nil
}

func extractGetCmd(value resp.Value) (*GetCommand, error) {

	if len(value.Array()) != 2 {
		return nil, fmt.Errorf("invalid number of variables for GET command")
	}
	cmd := GetCommand{
		key:   value.Array()[1].String(),
	}

	return &cmd, nil
}