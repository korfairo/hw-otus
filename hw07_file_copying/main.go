package main

import (
	"flag"
	"fmt"
	"log"
)

const serviceName = "go-cp"

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	l := log.Logger{}
	l.SetPrefix(fmt.Sprintf("%s:", serviceName))

	if from == "" {
		l.Println("missing --from operand")
		return
	}

	if to == "" {
		l.Println("missing --to operand")
	}

	if err := Copy(from, to, offset, limit); err != nil {
		l.Println("failed to copy:", err)
	}
}
