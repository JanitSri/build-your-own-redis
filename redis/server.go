package redis

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var InvalidCharacterError = errors.New("invalid error type")

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
	cmap sync.Map
}

func NewRedisServer(config ServerConfig) *RedisServer {
	return &RedisServer{
		ServerConfig: config,
	}
}

func (rs *RedisServer) Run() {
	ln, err := net.Listen(rs.network, fmt.Sprintf("%s:%s", rs.host, rs.port))
	if err != nil {
		log.Fatalln("Failed to bind", ln.Addr().String())
	}
	log.Println("Listening on", ln.Addr().String())

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
	c := make(chan Command)
	sc := newRedisScanner(conn, c)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(sc *RedisScanner) {
		defer wg.Done()
		sc.scan()
	}(sc)

	go func(c <-chan Command) {
		defer wg.Done()
		for cmd := range c {
			b := cmd.Execute(rs)
			conn.Write(b)
		}
	}(c)

	wg.Wait()
	log.Println("Closing connection for", conn.RemoteAddr().String())
}
