package provider

import (
	"fmt"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

type testServer struct {
	waitGroup sync.WaitGroup
}

func (s *testServer) serve(port net.Listener, cfg *ssh.ServerConfig) {
	defer s.waitGroup.Wait()

	for {
		tcpConn, err := port.Accept()

		if err != nil {
			fmt.Printf("Accept failed: %s\n", err)

			return
		}

		go s.handleConn(tcpConn, cfg)
	}
}

func (s *testServer) handleConn(tcpConn net.Conn, cfg *ssh.ServerConfig) {
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()
	defer tcpConn.Close()

	sshConn, _, _, err := ssh.NewServerConn(tcpConn, cfg)

	if err != nil {
		fmt.Printf("Handshake result: %s\n", err)

		return
	}

	fmt.Println("Server connection succeeded ..?")
	sshConn.Close()
}
