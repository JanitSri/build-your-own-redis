package parser

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/JanitSri/codecrafters-build-your-own-redis/customerror"
)

type RedisScanner struct {
	scanner bufio.Scanner
	cmdCh   chan<- Command
}

func NewRedisScanner(rw io.ReadWriter, cmdCh chan<- Command) *RedisScanner {
	return &RedisScanner{
		scanner: *bufio.NewScanner(rw),
		cmdCh:   cmdCh,
	}
}

func (rs *RedisScanner) Scan() {
	for rs.scanner.Scan() {
		t := rs.scanner.Text()
		cmd := rs.parseCmd(t, 0)
		if cmd != nil {
			rs.cmdCh <- cmd
		}
	}
	close(rs.cmdCh)
}

func (rs *RedisScanner) parseCmd(s string, np int) Command {
	rt := string(s[0])

	var cmd Command
	switch rt {
	case ARRAY:
		cmd = rs.handleArrays(s)
	case BULK_STRING:
		cmd = rs.handleBulkString(s, np)
	default:
		cmd = rs.handleCommand(s, np)
	}

	return cmd
}

func (rs *RedisScanner) handleArrays(s string) Command {
	if len(s) < 2 {
		return NewErrorCommand(customerror.InvalidCharacterError{})
	}

	n, err := strconv.Atoi(string(s[1:]))
	if err != nil {
		return NewErrorCommand(customerror.InvalidCharacterError{})
	}

	if !rs.scanner.Scan() {
		return NewErrorCommand(customerror.InvalidRedisCommandError{})
	}

	t := rs.scanner.Text()
	return rs.parseCmd(t, n)
}

func (rs *RedisScanner) handleBulkString(s string, np int) Command {
	if len(s) < 2 {
		return NewErrorCommand(customerror.InvalidCharacterError{})
	}

	n, err := strconv.Atoi(string(s[1:]))
	if err != nil {
		return NewErrorCommand(customerror.InvalidCharacterError{})
	}
	if n == 0 {
		log.Println("empty string")
		return nil
	}

	if !rs.scanner.Scan() {
		return NewErrorCommand(customerror.InvalidRedisCommandError{})
	}

	t := rs.scanner.Text()
	return rs.parseCmd(t, np-1)
}

func (rs *RedisScanner) handleCommand(cmdString string, np int) Command {
	var cmd Command

	switch strings.ToUpper(cmdString) {
	case PING:
		cmd = rs.parsePingCmd()
	case ECHO:
		cmd = rs.parseEchoCmd(np)
	case SET:
		cmd = rs.parseSetCmd(np)
	case GET:
		cmd = rs.parseGetCmd(np)
	case CONFIG:
		cmd = rs.parseConfigCmd(np)
	case KEYS:
		cmd = rs.parseKeysCmd(np)
	case INFO:
		cmd = rs.parseInfoCmd(np)
	default:
		return NewErrorCommand(customerror.InvalidRedisCommandError{})
	}

	return cmd
}

func (rs *RedisScanner) skipLen() error {
	i := 0
	for i < 2 && rs.scanner.Scan() {
		i++
	}
	if i != 2 {
		return customerror.InvalidNumberOfArgumentsError{}
	}

	return nil
}

func (rs *RedisScanner) parsePingCmd() Command {
	return NewPingCommand()
}

func (rs *RedisScanner) parseEchoCmd(np int) Command {
	if np < 1 {
		return NewErrorCommand(customerror.InvalidNumberOfArgumentsError{})
	}

	err := rs.skipLen()
	if err != nil {
		return NewErrorCommand(err)
	}
	s := rs.scanner.Text()

	return NewEchoCommand([]string{s})
}

func (rs *RedisScanner) parseSetCmd(np int) Command {
	if np < 2 {
		return NewErrorCommand(customerror.InvalidNumberOfArgumentsError{})
	}

	err := rs.skipLen()
	if err != nil {
		return NewErrorCommand(err)
	}
	k := rs.scanner.Text()
	np -= 1

	err = rs.skipLen()
	if err != nil {
		return NewErrorCommand(err)
	}
	v := rs.scanner.Text()
	np -= 1

	args := []string{k, v}
	flags := []*Flag{}

	for i := np; i > 0; i-- {
		err := rs.skipLen()
		if err != nil {
			return NewErrorCommand(err)
		}
		f := rs.scanner.Text()
		switch strings.ToUpper(f) {
		case PX:
			err := rs.skipLen()
			if err != nil {
				return NewErrorCommand(err)
			}
			flags = append(flags, NewFlag(f, rs.scanner.Text()))
			i -= 1
		default:
			return NewErrorCommand(customerror.InvalidCommandFlagError{Cmd: SET, Flag: f})
		}
	}

	return NewSetCommand(args, flags)
}

func (rs *RedisScanner) parseGetCmd(np int) Command {
	if np < 1 {
		return NewErrorCommand(customerror.InvalidNumberOfArgumentsError{})
	}

	err := rs.skipLen()
	if err != nil {
		return NewErrorCommand(err)
	}
	k := rs.scanner.Text()

	args := []string{k}
	flags := []*Flag{}

	return NewGetCommand(args, flags)
}

func (rs *RedisScanner) parseConfigCmd(np int) Command {
	if np < 1 {
		return NewErrorCommand(customerror.InvalidNumberOfArgumentsError{})
	}

	args := []string{}
	flags := []*Flag{}

	for i := np; i > 0; i-- {
		err := rs.skipLen()
		if err != nil {
			return NewErrorCommand(err)
		}
		f := rs.scanner.Text()
		i -= 1
		switch strings.ToUpper(f) {
		case GET:
			for i > 0 {
				err := rs.skipLen()
				if err != nil {
					return NewErrorCommand(err)
				}
				flags = append(flags, NewFlag(f, rs.scanner.Text()))
				i -= 1
			}
		default:
			NewErrorCommand(customerror.InvalidCommandFlagError{Cmd: CONFIG, Flag: f})
		}
	}

	return NewConfigCommand(args, flags)
}

func (rs *RedisScanner) parseKeysCmd(np int) Command {
	if np < 1 {
		return NewErrorCommand(customerror.InvalidNumberOfArgumentsError{})
	}

	err := rs.skipLen()
	if err != nil {
		return NewErrorCommand(err)
	}
	k := rs.scanner.Text()

	args := []string{k}
	flags := []*Flag{}

	return NewKeysCommand(args, flags)
}

func (rs *RedisScanner) parseInfoCmd(np int) Command {
	if np < 1 {
		return NewErrorCommand(customerror.InvalidNumberOfArgumentsError{})
	}

	args := []string{}
	flags := []*Flag{}

	if np == 1 {
		err := rs.skipLen()
		if err != nil {
			return NewErrorCommand(err)
		}
		a := rs.scanner.Text()
		np -= 1
		args = append(args, a)
	}

	return NewInfoCommand(args, flags)
}
