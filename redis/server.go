package redis

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/JanitSri/codecrafters-build-your-own-redis/data"
	"github.com/JanitSri/codecrafters-build-your-own-redis/parser"
	"github.com/google/uuid"
)

type ServerConfig struct {
	network string
	host    string
	port    string
}

func NewServerConfig(network, host, port string) *ServerConfig {
	return &ServerConfig{
		network,
		host,
		port,
	}
}

type RedisServer struct {
	id string
	ServerConfig
	RedisContext *data.RedisContext
}

func NewRedisServer(sc ServerConfig, rc data.RedisConfig, ri *data.RedisInfo) *RedisServer {
	rs := data.NewRedisStore(rc)
	id := fmt.Sprintf("%s-%s", ri.Replication.Role, uuid.New().String())

	return &RedisServer{
		ServerConfig: sc,
		RedisContext: data.NewRedisContext(ri, rs),
		id:           id,
	}
}

func (rs *RedisServer) Run(ctx context.Context) {
	ln, err := net.Listen(rs.network, fmt.Sprintf("%s:%s", rs.host, rs.port))
	if err != nil {
		log.Fatalf("%s failed to bind %s\n", rs.id, ln.Addr().String())
	}

	doneChan := make(chan any)

	go func(ln net.Listener) {
		defer close(doneChan)

		for {
			conn, err := ln.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					log.Printf("%s server shutting down...\n", rs.id)
					return
				} else {
					log.Printf("%s error accepting connection from %s\n", rs.id, conn.RemoteAddr().String())
					continue
				}
			}

			go rs.handleConnections(conn)
		}
	}(ln)

	<-ctx.Done()
	log.Printf("%s server received shutdown signal\n", rs.id)

	ln.Close()

	<-doneChan
}

func (rs *RedisServer) handleConnections(conn net.Conn) {
	defer conn.Close()
	log.Printf("%s handling connection from %s\n", rs.id, conn.RemoteAddr().String())
	c := make(chan parser.Command, 10)
	sc := parser.NewRedisScanner(conn, c)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		sc.Scan()
	}()

	go func() {
		defer wg.Done()
		for cmd := range c {
			b := cmd.Execute(rs.RedisContext)
			conn.Write(b)
		}
	}()

	wg.Wait()
	log.Printf("%s closing connection for %s\n", rs.id, conn.RemoteAddr().String())
}

func (rs *RedisServer) StartupTasks() {
	rs.displayBanner()
	rs.loadRDBFile()
}

func (rs *RedisServer) displayBanner() {
	addr := fmt.Sprintf("%s://%s:%s", rs.ServerConfig.network, rs.ServerConfig.host, rs.ServerConfig.port)

	var buf bytes.Buffer
	buf.WriteString("	        _._\n")
	buf.WriteString("           _.-``__ ''-._\n")
	buf.WriteString("      _.-``    `.  `_.  ''-._            Janit's Redis Server Implementation\n")
	buf.WriteString(fmt.Sprintf("  .-`` .-```.  ```\\/    _.,_ ''-._       Listening On: %s\n", addr))
	buf.WriteString(" (    '      ,       .-`  | `,    )\n")
	buf.WriteString(" |`-._`-...-` __...-.``-._|'` _.-'|\n")
	buf.WriteString(" |    `-._   `._    /     _.-'    |\n")
	buf.WriteString("  `-._    `-._  `-./  _.-'    _.-'\n")
	buf.WriteString(" |`-._`-._    `-.__.-'    _.-'_.-'|\n")
	buf.WriteString(" |    `-._`-._        _.-'_.-'    |\n")
	buf.WriteString("  `-._    `-._`-.__.-'_.-'    _.-'\n")
	buf.WriteString(" |`-._`-._    `-.__.-'    _.-'_.-'|\n")
	buf.WriteString(" |    `-._`-._        _.-'_.-'    |\n")
	buf.WriteString("  `-._    `-._`-.__.-'_.-'    _.-'\n")
	buf.WriteString("      `-._    `-.__.-'    _.-'\n")
	buf.WriteString("          `-._        _.-'\n")
	buf.WriteString("              `-.__.-'      \n\n")

	io.Copy(os.Stdout, &buf)
}

func (rs *RedisServer) loadRDBFile() {
	log.Println("Loading RDB file...")

	dir := strings.TrimSpace(rs.RedisContext.DataStore.GetConfig("dir"))
	fn := strings.TrimSpace(rs.RedisContext.DataStore.GetConfig("dbfilename"))

	fd, err := os.OpenRoot(dir)
	if err != nil {
		return
	}
	defer fd.Close()

	f, err := fd.Open(fn)
	if err != nil {
		return
	}
	defer f.Close()

	bd, err := io.ReadAll(f)
	if err != nil {
		return
	}

	pairs := parser.ParseRBDFile(bd)

	for k, v := range pairs {
		rs.RedisContext.DataStore.Set(k, v)
	}
}
