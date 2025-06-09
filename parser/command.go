package parser

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/JanitSri/codecrafters-build-your-own-redis/customerror"
	"github.com/JanitSri/codecrafters-build-your-own-redis/data"
	"github.com/JanitSri/codecrafters-build-your-own-redis/util"
)

func writeBulkString(s string) []byte {
	l := strconv.Itoa(len(s))
	var buf bytes.Buffer
	buf.WriteString(BULK_STRING)
	buf.WriteString(l)
	buf.WriteString(REDIS_TERMINATOR)
	buf.WriteString(s)
	buf.WriteString(REDIS_TERMINATOR)
	return buf.Bytes()
}

func writeSimpleError(err error) []byte {
	en := reflect.TypeOf(err).Name()
	em := fmt.Sprintf("%s %s", en, err.Error())

	var buf bytes.Buffer
	buf.WriteString(SIMPLE_ERROR)
	buf.WriteString(em)
	buf.WriteString(REDIS_TERMINATOR)
	return buf.Bytes()
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
	Execute(*data.RedisContext) []byte
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

func (pc *PingCommand) Execute(rc *data.RedisContext) []byte {
	log.Println("ponging...")

	var buf bytes.Buffer
	buf.WriteString(PONG)
	return buf.Bytes()
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

func (ec *EchoCommand) Execute(rc *data.RedisContext) []byte {
	log.Println("echoing...")

	if len(ec.args) != 1 {
		return writeSimpleError(customerror.InvalidNumberOfArgumentsError{})
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

func (sc *SetCommand) Execute(rc *data.RedisContext) []byte {
	log.Println("setting...")

	if len(sc.args) != 2 {
		return writeSimpleError(customerror.InvalidNumberOfArgumentsError{})
	}

	v := data.NewRedisValue(sc.args[1], time.Time{})
	for _, f := range sc.flags {
		switch strings.ToUpper(f.name) {
		case PX:
			ms, _ := strconv.Atoi(f.value)
			v.SetExpiry(time.Now().Add(time.Duration(ms) * time.Millisecond))
		default:
			return writeSimpleError(customerror.InvalidCommandFlagError{Cmd: SET, Flag: f.name})
		}
	}

	rc.DataStore.Set(sc.args[0], v)

	var buf bytes.Buffer
	buf.WriteString(OK)

	return buf.Bytes()
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

func (gc *GetCommand) Execute(rc *data.RedisContext) []byte {
	log.Println("getting...")

	if len(gc.args) != 1 {
		return writeSimpleError(customerror.InvalidNumberOfArgumentsError{})
	}

	v, ok := rc.DataStore.Get(gc.args[0])
	var buf bytes.Buffer
	if !ok {
		buf.WriteString(NULL_BULK_STRING)
		return buf.Bytes()
	}

	res := v.(*data.RedisValue)
	if res.IsExpired() {
		buf.WriteString(NULL_BULK_STRING)
		return buf.Bytes()
	}

	vs := res.Value().(string)
	return writeBulkString(vs)
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

func (cc *ConfigCommand) Execute(rc *data.RedisContext) []byte {
	log.Println("configuring...")

	var buf bytes.Buffer
	for _, f := range cc.flags {
		switch strings.ToUpper(f.name) {
		case GET:
			// todo: update to handle multiple CONFIG GET params
			// https://redis.io/docs/latest/commands/config-get/
			cn := f.value
			cv := rc.DataStore.GetConfig(cn)
			buf.WriteString(ARRAY + strconv.Itoa(2) + REDIS_TERMINATOR)
			buf.Write(writeBulkString(cn))
			buf.Write(writeBulkString(cv))
		default:
			return writeSimpleError(customerror.InvalidCommandFlagError{Cmd: CONFIG, Flag: f.name})
		}
	}

	return buf.Bytes()
}

type KeysCommand struct {
	BaseCommand
}

func NewKeysCommand(args []string, flags []*Flag) *KeysCommand {
	return &KeysCommand{
		BaseCommand{
			args,
			flags,
		},
	}
}

func (kc *KeysCommand) Execute(rc *data.RedisContext) []byte {
	log.Println("Getting Keys...")

	var buf bytes.Buffer

	if len(kc.args) == 0 {
		return writeSimpleError(customerror.InvalidNumberOfArgumentsError{})
	}

	p := kc.args[0]
	var tempBuf bytes.Buffer
	l := 0
	ks := rc.DataStore.Keys()
	for _, k := range ks {
		ks := k.(string)
		if p == "*" {
			tempBuf.Write(writeBulkString(ks))
			l++
		}
	}

	buf.WriteString(ARRAY)
	buf.WriteString(strconv.Itoa(l))
	buf.WriteString(REDIS_TERMINATOR)
	tempBuf.WriteTo(&buf)

	return buf.Bytes()
}

type InfoCommand struct {
	BaseCommand
}

func NewInfoCommand(args []string, flags []*Flag) *InfoCommand {
	return &InfoCommand{
		BaseCommand{
			args,
			flags,
		},
	}
}

func (ic *InfoCommand) Execute(rc *data.RedisContext) []byte {
	log.Println("info...")

	var arg string
	if len(ic.args) >= 1 {
		arg = ic.args[0]
	}

	var buf bytes.Buffer
	switch strings.ToUpper(arg) {
	case REPLICATION:
		s, err := util.SerializeSection(*rc.RedisInfo.Replication)
		if err != nil {
			return writeSimpleError(err)
		}
		sb := writeBulkString(s)
		buf.Write(sb)
	default:
		return writeSimpleError(customerror.InvalidNumberOfArgumentsError{})
	}

	return buf.Bytes()
}

type ErrorCommand struct {
	err error
}

func NewErrorCommand(err error) *ErrorCommand {
	return &ErrorCommand{
		err,
	}
}

func (ec *ErrorCommand) Execute(_ *data.RedisContext) []byte {
	return writeSimpleError(ec.err)
}
