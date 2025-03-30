package parser

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"
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
		cmd := rs.parseCmd(t)
		if cmd != nil {
			rs.cmdCh <- cmd
		}
	}
	close(rs.cmdCh)
}

func (rs *RedisScanner) parseCmd(s string) Command {
	rt := string(s[0])

	var cmd Command
	switch rt {
	case ARRAY:
		cmd = rs.handleArrays(s)
	case BULK_STRING:
		cmd = rs.handleBulkString(s)
	default:
		cmd = rs.handleCommand(s)
	}

	return cmd
}

func (rs *RedisScanner) handleArrays(s string) Command {
	if len(s) < 2 {
		log.Fatalln(InvalidCharacterError)
	}

	n, err := strconv.Atoi(string(s[1:]))
	if err != nil {
		log.Fatalln(InvalidCharacterError)
	}

	for i := n; i > 0 && rs.scanner.Scan(); i-- {
		t := rs.scanner.Text()
		return rs.parseCmd(t)
	}

	return nil
}

func (rs *RedisScanner) handleBulkString(s string) Command {
	if len(s) < 2 {
		log.Fatalln(InvalidCharacterError)
	}

	n, err := strconv.Atoi(string(s[1:]))
	if err != nil {
		log.Fatalln(InvalidCharacterError)
	}
	if n == 0 {
		log.Println("empty string")
		return nil
	}

	if !rs.scanner.Scan() {
		log.Fatalln(InvalidRedisCommandError)
	}

	t := rs.scanner.Text()
	return rs.parseCmd(t)
}

func (rs *RedisScanner) handleCommand(cmdString string) Command {
	var cmd Command

	switch strings.ToUpper(cmdString) {
	case PING:
		cmd = rs.parsePingCmd()
	case ECHO:
		cmd = rs.parseEchoCmd()
	case SET:
		cmd = rs.parseSetCmd()
	case GET:
		cmd = rs.parseGetCmd()
	default:
		log.Fatalln(InvalidRedisCommandError)
	}

	return cmd
}

func (rs *RedisScanner) skipLen() {
	i := 0
	for i < 2 && rs.scanner.Scan() {
		i++
	}
	if i != 2 {
		log.Fatalln(InvalidNumberOfArguments)
	}
}

func (rs *RedisScanner) parsePingCmd() *PingCommand {
	return NewPingCommand()
}

func (rs *RedisScanner) parseEchoCmd() *EchoCommand {
	rs.skipLen()
	s := rs.scanner.Text()

	return NewEchoCommand([]string{s})
}

func (rs *RedisScanner) parseSetCmd() *SetCommand {
	rs.skipLen()
	k := rs.scanner.Text()

	rs.skipLen()
	v := rs.scanner.Text()

	args := []string{k, v}
	flags := []*Flag{}

	return NewSetCommand(args, flags)
}

func (rs *RedisScanner) parseGetCmd() *GetCommand {
	rs.skipLen()
	k := rs.scanner.Text()

	args := []string{k}
	flags := []*Flag{}

	return NewGetCommand(args, flags)
}
