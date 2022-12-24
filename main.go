package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/42LoCo42/go-zeolite"
)

func main() {
	genKeys := flag.Bool("gen", false, "Generate a new key pair")
	genSaveTo := flag.String("out", "", "Where to save a new key pair (used with -gen)")

	connectTo := flag.Bool("connect", false, "Connect to a server")
	serve := flag.Bool("serve", false, "Serve the terminal to clients")
	identFile := flag.String("ident", "", "The identity to assume")
	addr := flag.String("addr", "tcp://127.0.0.1:1234", "The port to serve on or connect to (eg. tcp://127.0.0.1:1234")

	flag.Parse()

	if *genKeys {
		newIdentity := Identity{}
		err := newIdentity.Generate()

		if err != nil {
			panic(err)
		}

		if len(*genSaveTo) > 0 {
			err = SaveIdentity(*genSaveTo, &newIdentity)

			if err != nil {
				panic(err)
			}
		}

		newIdentity.Print()
	} else if len(*identFile) > 0 {
		identFileBytes, err := ReadFile(*identFile)

		if err != nil {
			panic(err)
		}

		identity := Identity{}
		err = identity.LoadIdentity(identFileBytes)

		if err != nil {
			panic(err)
		}

		identity.Print()

		proto, ipAndPort, err := ParseAddr(*addr)

		if err != nil {
			panic(err)
		}

		if *serve {
			fmt.Println("Serve mode")

			server := &Server{identity: &identity}
			err = server.Listen(proto, ipAndPort)

			if err != nil {
				panic(err)
			}
		} else if *connectTo {
			fmt.Println("Client mode")

			conn, err := net.Dial(proto, ipAndPort)

			if err != nil {
				panic(err)
			}

			stream, err := identity.NewStream(conn, shouldTrustClient)

			if err != nil {
				panic(err)
			}

			go func() {
				io.Copy(stream, os.Stdin)
				os.Stdin.Close()
			}()

			zeolite.BlockCopy(os.Stdout, stream)
			os.Stdout.Close()
		}
	}

	// if err := spawnTerminal(); err != nil {
	// 	log.Fatal(err)
	// }
}
