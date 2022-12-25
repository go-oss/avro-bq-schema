package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-oss/avro-bq-schema/schema"
)

const (
	name          = "avro-bq-schema"
	defaultIndent = 2
)

var (
	errInvalidArgument = errors.New("invalid argument")
)

func main() {
	var err error
	defer func() {
		if err != nil {
			os.Exit(1)
		}
	}()

	log.SetFlags(0)

	var indent int
	cmd := flag.NewFlagSet("", flag.ExitOnError)
	cmd.IntVar(&indent, "indent", defaultIndent, "output JSON indent size")
	cmd.Usage = func() {
		b := new(strings.Builder)
		fmt.Fprintln(b, "Usage:")
		fmt.Fprintf(b, "  %s [file]\n", name)
		b.WriteRune('\n')
		fmt.Fprintln(b, "Flags:")
		//nolint:errcheck
		io.WriteString(cmd.Output(), b.String())
		cmd.PrintDefaults()
	}
	err = cmd.Parse(os.Args[1:])
	if err != nil {
		log.Println("cmd.Parse:", err)
		return
	}
	files := cmd.Args()

	if len(files) == 0 {
		err = errInvalidArgument
		log.Println(err)
		return
	}

	f, err := os.Open(files[0])
	if err != nil {
		log.Println("os.Open:", err)
		return
	}
	defer f.Close()

	d, err := io.ReadAll(f)
	if err != nil {
		log.Println("io.ReadAll:", err)
		return
	}

	bqSchema, err := schema.Convert(d)
	if err != nil {
		log.Println("schema.Convert:", err)
		return
	}

	dst, err := schema.ToJSON(bqSchema, indent)
	if err != nil {
		log.Println("schema.ToJSON:", err)
		return
	}

	_, err = io.WriteString(os.Stdout, string(dst))
	if err != nil {
		log.Println("io.WriteString:", err)
		return
	}
}
