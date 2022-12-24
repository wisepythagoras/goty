package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/42LoCo42/go-zeolite"
	"github.com/creack/pty"
	"golang.org/x/term"
)

func shouldTrustClient(otherPK zeolite.SignPK) (bool, error) {
	// b64 := zeolite.Base64Enc(otherPK[:])

	// for _, id := range trustList {
	// 	if id == b64 {
	// 		return true, nil
	// 	}
	// }

	return true, nil
}

type Server struct {
	identity *Identity
}

// Listen announces on the local network address.
// The network must be "tcp", "tcp4", "tcp6" or "unix".
func (s *Server) Listen(proto, addr string) error {
	listener, err := net.Listen(proto, addr)

	if err != nil {
		return err
	}

	for {
		client, err := listener.Accept()

		if err != nil {
			return err
		}

		stream, err := s.identity.NewStream(client, shouldTrustClient)

		go handler(s, stream)

		fmt.Println(client.RemoteAddr())
	}

	return nil
}

// https://docs.cossacklabs.com/themis/languages/go/installation/
// https://github.com/cossacklabs/themis/blob/master/docs/examples/go/secure_session_server.go
// https://stackoverflow.com/questions/51535047/golang-read-byte-array-from-to-and-display-result#51535926
func spawnTerminal(stream *zeolite.Stream) error {
	// Create arbitrary command.
	c := exec.Command("bash")

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}

	// Make sure to close the pty at the end.
	defer func() {
		ptmx.Close()
	}()

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)

	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()

	// Signal for the initial resizing.
	ch <- syscall.SIGWINCH

	defer func() {
		// Cleanup signals when the session is done.
		signal.Stop(ch)
		close(ch)
	}()

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		panic(err)
	}

	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	go func() {
		// _, _ = io.Copy(ptmx, os.Stdin)
		zeolite.BlockCopy(ptmx, stream)
	}()

	// For use with Themis. ---
	// go WriteToPty(ptmx)

	fmt.Printf("Connected!\n\r")

	// Same system pty.
	_, _ = io.Copy(stream, ptmx) // Or os.Stdout instead of stream.

	// For use with Themis. ---
	// ReadFromPty(ptmx)

	return nil
}

func handler(srv *Server, stream *zeolite.Stream) {
	err := spawnTerminal(stream)

	if err != nil {
		panic(err)
	}
}
