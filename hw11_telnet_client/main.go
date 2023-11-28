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

const updateInterval = time.Millisecond * 10

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	senderErrCh := doWithInterval(ctx, updateInterval, client.Send)
	receiverErrCh := doWithInterval(ctx, updateInterval, client.Receive)

	go processErrors(ctx, logger, senderErrCh, receiverErrCh)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	<-quitCh
}

func doWithInterval(ctx context.Context, interval time.Duration, f func() error) chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := f(); err != nil {
					errCh <- err
				}
			}
		}
	}()

	return errCh
}

func processErrors(ctx context.Context, logger *log.Logger, senderErrCh, receiverErrCh chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-senderErrCh:
			if errors.Is(err, io.EOF) {
				logger.Fatal("EOF")
				return
			}
		case err := <-receiverErrCh:
			if errors.Is(err, io.EOF) {
				logger.Fatal("Connection was closed by peer")
				return
			}
		}
	}
}
