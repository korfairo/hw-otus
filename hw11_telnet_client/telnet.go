package main

import (
	"bufio"
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
)

var ErrNotConnected = errors.New("telnet client not connected to host")

type TelnetClient struct {
	address string
	timeout time.Duration

	in  *bufio.Reader
	out io.Writer

	conn       net.Conn
	connReader *bufio.Reader
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return TelnetClient{
		address: address,
		timeout: timeout,
		out:     out,
		in:      bufio.NewReader(in),
	}
}

func (t *TelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return errors.Wrap(err, "failed to connect")
	}
	t.conn = conn
	t.connReader = bufio.NewReader(conn)
	return nil
}

func (t *TelnetClient) Send() error {
	if connected := t.checkConnection(); !connected {
		return ErrNotConnected
	}
	err := passString(t.in, t.conn)
	return errors.Wrap(err, "failed to send message")
}

func (t *TelnetClient) Receive() error {
	if connected := t.checkConnection(); !connected {
		return ErrNotConnected
	}
	err := passString(t.connReader, t.out)
	return errors.Wrap(err, "failed to receive message")
}

func (t *TelnetClient) Close() error {
	err := t.conn.Close()
	return errors.Wrap(err, "failed to close connection")
}

func (t *TelnetClient) checkConnection() bool {
	return t.conn != nil
}

func passString(from *bufio.Reader, to io.Writer) error {
	bytes, err := from.ReadBytes('\n')
	if err != nil {
		return errors.Wrap(err, "failed to read bytes")
	}
	_, err = to.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "failed to write bytes")
	}
	return nil
}
