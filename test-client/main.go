package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/quic-go/quic-go"
)

const message = "foobar"

func main() {
	var addr string
	flag.StringVar(&addr, "a", "", "server address")
	flag.Parse()

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"sfkit"},
	}
	conn, err := quic.DialAddr(context.Background(), addr, tlsConf, nil)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Client: Sending '%s'\n", message)
	_, err = stream.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
	stream.Close() // close for writing (send EOF)

	buf := make([]byte, len(message))
	_, err = io.ReadFull(stream, buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Client: Got '%s'\n", buf)
}
