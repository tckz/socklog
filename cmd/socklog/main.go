package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/tckz/socklog"
)

func main() {
	ec := run()
	os.Exit(ec)
}

func run() int {
	log.SetOutput(os.Stderr)
	bind := flag.String("bind", "0.0.0.0:8100", "Address and tcp port to bind")
	dest := flag.String("dest", "", "Address and tcp port where connect to")
	mask := flag.Bool("mask", true, "Mask control code(0x00-1f, except 0x09,0a,0c,0d)")
	out := flag.String("out", "", "/path/to/outfile (default stdout)")
	flag.Parse()

	if *dest == "" {
		log.Printf("--dest must be specified")
		return 1
	}

	log.Printf("Listen: %s", *bind)
	listen, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Printf("*** Failed to listen %s: %s", *bind, err)
		return 1
	}
	defer listen.Close()

	var wOut io.Writer
	wOut = os.Stdout
	if *out != "" {
		log.Printf("Out: %s", *out)
		f, err := os.Create(*out)
		if err != nil {
			log.Printf("*** Failed to open %s: %s", *out, err)
			return 1
		}
		defer f.Close()
		wOut = f
	}

	if *mask {
		wOut = socklog.NewMaskingWriter(wOut)
	}

	for {
		srcConn, err := listen.Accept()
		if err != nil {
			log.Printf("*** Failed to accept: %s", err)
			return 1
		}

		go func() {
			defer srcConn.Close()

			destConn, err := net.Dial("tcp", *dest)
			if err != nil {
				log.Printf("*** Failed to connect %s: %s", *dest, err)
				return
			}
			defer destConn.Close()

			wg := new(sync.WaitGroup)
			wg.Add(1)
			go func() {
				defer wg.Done()
				upStream := io.MultiWriter(wOut, destConn)
				if _, err := io.Copy(upStream, srcConn); err != nil {
					log.Printf("*** Failed to Copy to upstream: %s", err)
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				downStream := io.MultiWriter(wOut, srcConn)
				if _, err := io.Copy(downStream, destConn); err != nil {
					log.Printf("*** Failed to Copy to downstream: %s", err)
				}
			}()

			wg.Wait()
		}()
	}

	return 0
}
