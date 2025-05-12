package parser

import (
	"bytes"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/JanitSri/codecrafters-build-your-own-redis/data"
)

var ErrInvalidNumberOfArguments = errors.New("invalid number of arguments")
var ErrInvalidArgument = errors.New("invalid argument")
var ErrInvalidCharacterError = errors.New("invalid error type")

func writeBulkString(s string) []byte {
	l := strconv.Itoa(len(s))
	var b bytes.Buffer
	b.WriteString(BULK_STRING)
	b.WriteString(l)
	b.WriteString(REDIS_TERMINATOR)
	b.WriteString(s)
	b.WriteString(REDIS_TERMINATOR)
	return b.Bytes()
}

type Flag struct {
	name  string
	value string
}

func NewFlag(name, value string) *Flag {
	return &Flag{
		name,
		value,
	}
}

type Command interface {
	Execute(data.DataStore) []byte
}

type BaseCommand struct {
	args  []string
	flags []*Flag
}

type PingCommand struct {
}

func NewPingCommand() *PingCommand {
	return &PingCommand{}
}

func (pc *PingCommand) Execute(ds data.DataStore) []byte {
	log.Println("ponging...")

	var b bytes.Buffer
	b.WriteString(PONG)
	return b.Bytes()
}

type EchoCommand struct {
	BaseCommand
}

func NewEchoCommand(args []string) *EchoCommand {
	return &EchoCommand{
		BaseCommand{
			args,
			nil,
		},
	}
}

func (ec *EchoCommand) Execute(ds data.DataStore) []byte {
	log.Println("echoing...")

	if len(ec.args) != 1 {
		log.Fatal(ErrInvalidNumberOfArguments)
	}

	return writeBulkString(ec.args[0])
}

type SetCommand struct {
	BaseCommand
}

func NewSetCommand(args []string, flags []*Flag) *SetCommand {
	return &SetCommand{
		BaseCommand{
			args,
			flags,
		},
	}
}

func (sc *SetCommand) Execute(ds data.DataStore) []byte {
	log.Println("setting...")

	if len(sc.args) != 2 {
		log.Fatal(ErrInvalidNumberOfArguments)
	}

	v := data.NewRedisValue(sc.args[1], time.Time{})
	for _, f := range sc.flags {
		switch strings.ToUpper(f.name) {
		case PX:
			ms, _ := strconv.Atoi(f.value)
			v.SetExpiry(time.Now().Add(time.Duration(ms) * time.Millisecond))
		default:
			log.Fatalln(InvalidSetCommandFlag(f.name))
		}
	}

	ds.Set(sc.args[0], v)

	var b bytes.Buffer
	b.WriteString(OK)

	return b.Bytes()
}

type GetCommand struct {
	BaseCommand
}

func NewGetCommand(args []string, flags []*Flag) *GetCommand {
	return &GetCommand{
		BaseCommand{
			args,
			flags,
		},
	}
}

func (gc *GetCommand) Execute(ds data.DataStore) []byte {
	log.Println("getting...")

	if len(gc.args) != 1 {
		log.Fatal(ErrInvalidNumberOfArguments)
	}

	v, ok := ds.Get(gc.args[0])
	var b bytes.Buffer
	if !ok {
		b.WriteString(NULL_BULK_STRING)
		return b.Bytes()
	}

	res := v.(*data.RedisValue)
	if res.IsExpired() {
		b.WriteString(NULL_BULK_STRING)
		return b.Bytes()
	}

	return writeBulkString(res.Value())
}

type ConfigCommand struct {
	BaseCommand
}

func NewConfigCommand(args []string, flags []*Flag) *ConfigCommand {
	return &ConfigCommand{
		BaseCommand{
			args,
			flags,
		},
	}
}

func (cc *ConfigCommand) Execute(ds data.DataStore) []byte {
	log.Println("configuring...")

	var res bytes.Buffer
	for _, f := range cc.flags {
		switch strings.ToUpper(f.name) {
		case GET:
			// todo: update to handle multiple CONFIG GET params
			// https://redis.io/docs/latest/commands/config-get/
			cn := f.value
			cv := ds.GetConfig(cn)
			res.WriteString(ARRAY + strconv.Itoa(2) + REDIS_TERMINATOR)
			res.Write(writeBulkString(cn))
			res.Write(writeBulkString(cv))
		default:
			log.Fatalln(InvalidSetCommandFlag(f.name))
		}
	}

	return res.Bytes()
}
