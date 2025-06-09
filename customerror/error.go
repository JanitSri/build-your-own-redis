package customerror

import (
	"fmt"
	"reflect"
)

type InvalidNumberOfArgumentsError struct{}

func (e InvalidNumberOfArgumentsError) Error() string {
	return "invalid number of arguments"
}

type InvalidArgumentError struct{}

func (e InvalidArgumentError) Error() string {
	return "invalid argument"
}

type InvalidCharacterError struct{}

func (e InvalidCharacterError) Error() string {
	return "invalid character"
}

type InvalidRDBValueTypeError struct{}

func (e InvalidRDBValueTypeError) Error() string {
	return "invalid value type in RDB file"
}

type InvalidRespDataTypeError struct{}

func (e InvalidRespDataTypeError) Error() string {
	return "invalid redis RESP type"
}

type InvalidRedisCommandError struct{}

func (e InvalidRedisCommandError) Error() string {
	return "invalid redis command"
}

type NoLeaderAvailableError struct{}

func (e NoLeaderAvailableError) Error() string {
	return "no redis leader available"
}

type InvalidServerConfigError struct {
	Name string
}

func (e InvalidServerConfigError) Error() string {
	return fmt.Sprintf("invalid Redis config: %s", e.Name)
}

type InvalidCommandFlagError struct {
	Cmd  string
	Flag string
}

func (e InvalidCommandFlagError) Error() string {
	return fmt.Sprintf("invalid flag for %s command: %s", e.Cmd, e.Flag)
}

type KeysCommandError struct {
	Flag string
}

func (e KeysCommandError) Error() string {
	return fmt.Sprintf("keys command error: %s", e.Flag)
}

type UnsupportedFieldTypeError struct {
	Kind reflect.Kind
}

func (e UnsupportedFieldTypeError) Error() string {
	return fmt.Sprintf("unsupported field type: %s", e.Kind.String())
}
