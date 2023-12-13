package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

var (
	timeout time.Duration
	host    string
	port    string
)

func parseFlags() error {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() != 2 {
		return errors.New("some arguments are missing, pass host and port to connect")
	}
	host = flag.Args()[0]
	port = flag.Args()[1]

	return nil
}

func main() {
	logger := log.New(os.Stderr, "...", 0)

	if err := parseFlags(); err != nil {
		logger.Println("Failed to parse arguments:", err)
		return
	}

	address := net.JoinHostPort(host, port)
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		logger.Println("Couldn't connect to host:", err)
		return
	}
	defer client.Close()
	logger.Println("Connected to", address)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := client.Send()
				if err != nil && errors.Is(err, io.EOF) {
					logger.Fatal("EOF")
					return
				}

			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := client.Receive()
				if err != nil && errors.Is(err, io.EOF) {
					logger.Fatal("Connection was closed by peer")
					return
				}
			}
		}
	}()

	<-ctx.Done()
}
