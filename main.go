package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-oss/avro-bq-schema/schema"
)

const (
	name = "avro-bq-schema"
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

	cmd := flag.NewFlagSet("", flag.ExitOnError)
	defaultUsage := cmd.Usage
	cmd.Usage = func() {
		defaultUsage()
		fmt.Fprintf(cmd.Output(), "  %s [file]\n", name)
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

	j, err := bqSchema.ToJSONFields()
	if err != nil {
		log.Println("bqSchema.ToJSONFields:", err)
		return
	}

	fmt.Println(string(j))
}
