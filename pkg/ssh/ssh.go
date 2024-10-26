package ssh

import (
	"fmt"
	"os"

	"github.com/prudhvideep/fnd/pkg/config"
	"golang.org/x/crypto/ssh"
)

func InitializeConnection(credentials *config.Credentials) (*ssh.Client, error) {
	hostName := credentials.Host + ":" + credentials.Port

	key, err := loadPrivateKey(credentials.KeyPath)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: credentials.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", hostName, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH connection: %v", err)
	}

	return client, nil

}

func loadPrivateKey(path string) (ssh.Signer, error) {
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return signer, nil

}
