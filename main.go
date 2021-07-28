// port_fwd -l <file> -s <src_port> -d <dst_port>
package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	fname := flag.String("l", "port_fwd.log", "log file")
	srcPort := flag.String("s", ":80", "source port")
	dest := flag.String("d", "127.0.0.1:3000", "destination")
	flag.Parse()

	fp, err := os.Create(*fname)
	if err != nil {
		panic(err)
	}

	logger := log.New(fp, "", log.LUTC|log.Lshortfile)

	ln, err := net.Listen("tcp", *srcPort)
	if err != nil {
		logger.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Fatal(err)
		}

		go func(conn net.Conn) {
			proxy, err := net.Dial("tcp", *dest)
			if err != nil {
				logger.Fatal(err)
			}

			go copy(conn, proxy)
			go copy(proxy, conn)
		}(conn)
	}
}

func copy(src, dst net.Conn) {
	defer src.Close()
	defer dst.Close()
	io.Copy(src, dst)
}
