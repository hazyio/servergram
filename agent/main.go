package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func main() {
	println("Hello world")
	connectWithSSh([]byte(`-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACBkvPjfVB02Al0Jd7yjC2QPXwGG8fWkm5OFiTooZ583TgAAAKDocCPX6HAj
1wAAAAtzc2gtZWQyNTUxOQAAACBkvPjfVB02Al0Jd7yjC2QPXwGG8fWkm5OFiTooZ583Tg
AAAEA0/M0LGqQHRigxDkyvOcvFyNUyG6HEzDJYKJkKRM1ggGS8+N9UHTYCXQl3vKMLZA9f
AYbx9aSbk4WJOihnnzdOAAAAG2RhbmllbEBkYW5pZWwtVGhpbmtQYWQtWDI2MAEC
-----END OPENSSH PRIVATE KEY-----
`))
}

func connectWithSSh(keyBytes []byte) {
	home, _ := os.UserHomeDir()
 
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		log.Fatal(err) // for encrypted keys use ssh.ParsePrivateKeyWithPassphrase
	}
	// Host key verification via known_hosts
	hostKeyCallback, err := knownhosts.New(filepath.Join(home, ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}

	cfg := &ssh.ClientConfig{
		User:            "daniel",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: hostKeyCallback,
		Timeout:         10 * time.Second,
	}
	client, err := ssh.Dial("tcp", "35.208.98.158:22", cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// One session per command
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	out, err := session.CombinedOutput("uname -a")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))

}
