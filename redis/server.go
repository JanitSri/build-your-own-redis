package redis

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"

	"github.com/JanitSri/codecrafters-build-your-own-redis/data"
	"github.com/JanitSri/codecrafters-build-your-own-redis/parser"
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
	ServerConfig
	redisContext *data.RedisContext
}

func NewRedisServer(sc ServerConfig, rc data.RedisConfig, ri *data.RedisInfo) *RedisServer {
	rs := data.NewRedisStore(rc)

	return &RedisServer{
		ServerConfig: sc,
		redisContext: data.NewRedisContext(ri, rs),
	}
}

func (rs *RedisServer) Run() {

	ln, err := net.Listen(rs.network, fmt.Sprintf("%s:%s", rs.host, rs.port))
	if err != nil {
		log.Fatalln("Failed to bind", ln.Addr().String())
	}
	rs.startupTasks()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	doneChan := make(chan interface{})

	go func(ln net.Listener) {
		defer close(doneChan)

		for {
			conn, err := ln.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					log.Println("server shutting down...")
					return
				} else {
					log.Println("error accepting connection from", conn.RemoteAddr().String())
					continue
				}
			}

			go rs.handleConnections(conn)
		}
	}(ln)

	sig := <-sigChan
	log.Println("Shutting down with", sig)

	ln.Close()

	<-doneChan
}

func (rs *RedisServer) handleConnections(conn net.Conn) {
	defer conn.Close()
	log.Println("handling connection from", conn.RemoteAddr().String())
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
			b := cmd.Execute(rs.redisContext)
			conn.Write(b)
		}
	}()

	wg.Wait()
	log.Println("Closing connection for", conn.RemoteAddr().String())
}

func (rs *RedisServer) startupTasks() {
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

	dir := rs.redisContext.DataStore.GetConfig("dir")
	fn := rs.redisContext.DataStore.GetConfig("dbfilename")
	p := path.Join(strings.TrimSpace(dir), strings.TrimSpace(fn))

	if p == "" {
		return
	}

	rb, err := os.ReadFile(p)
	if err != nil {
		return
	}

	pairs := parser.ParseRBDFile(rb)

	for k, v := range pairs {
		rs.redisContext.DataStore.Set(k, v)
	}
}
