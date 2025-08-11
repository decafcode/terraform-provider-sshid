package provider

import (
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

const testPublicKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIPu1TqZTnE2574YCR5bNiw03wd0vbsaribTbz+LM4pdd"
const testServerHost = "localhost"
const testServerPort = 54321

type testFramework struct {
	port   net.Listener
	server testServer
}

func newFramework() (*testFramework, error) {
	privateKey, err := ssh.ParsePrivateKey([]byte(`
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACD7tU6mU5xNue+GAkeWzYsNN8HdL27Gq4m028/izOKXXQAAAIjgCQxY4AkM
WAAAAAtzc2gtZWQyNTUxOQAAACD7tU6mU5xNue+GAkeWzYsNN8HdL27Gq4m028/izOKXXQ
AAAEB8TxrsvYPLq0Adr0G/q+ttEWA4Bsraj8xrqBaMnqBarvu1TqZTnE2574YCR5bNiw03
wd0vbsaribTbz+LM4pddAAAAAAECAwQF
-----END OPENSSH PRIVATE KEY-----
	`))

	if err != nil {
		return nil, err
	}

	cfg := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			return nil, fmt.Errorf("no access")
		},
	}

	cfg.AddHostKey(privateKey)

	addr := net.JoinHostPort(testServerHost, fmt.Sprintf("%d", testServerPort))
	port, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	fw := &testFramework{
		port:   port,
		server: testServer{},
	}

	go fw.server.serve(port, cfg)

	return fw, nil
}

func (fw *testFramework) Close() error {
	return fw.port.Close()
}

func (fw *testFramework) Host() string {
	return testServerHost
}

func (fw *testFramework) Port() uint16 {
	return testServerPort
}

func (fw *testFramework) PublicKey() string {
	return testPublicKey
}
