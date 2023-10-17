package host

import (
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Host struct {
	Addr         string
	username     string
	clientConfig *ssh.ClientConfig

	// TODO: separate out to HostConnection
	Connection *ssh.Client
}

func NewHost(connectionStr string, privateKeyPath string) (*Host, error) {
	spl := strings.Split(connectionStr, "@")
	if len(spl) < 2 {
		return nil, errors.New("you must provide a username and an address for the host")
	}
	hostPort := spl[1]

	spl2 := strings.Split(hostPort, ":")

	host := spl2[0]

	port := "22"
	if len(spl2) > 1 {
		port = spl2[1]
	}

	username := spl[0]

	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParseRawPrivateKey(keyBytes)
	if err != nil {
		log.Fatalf("unable to parse private key: %s", err)
	}

	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		return nil, err
	}

	return &Host{
		Addr:     net.JoinHostPort(host, port),
		username: username,
		clientConfig: &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         5 * time.Second,
		},
	}, nil
}

func (h *Host) Connect() (*ssh.Client, error) {
	client, err := ssh.Dial("tcp", h.Addr, h.clientConfig)
	if err != nil {
		return nil, err
	}

	h.Connection = client

	return client, nil
}

func (h *Host) Close() error {
	return h.Connection.Close()
}
