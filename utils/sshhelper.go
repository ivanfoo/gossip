package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"golang.org/x/crypto/ssh"
)	

func SSHConnect(username string, hostname string, keyPath string) *ssh.Client {
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		log.Fatalf("parse key failed:%v", err)
	}
	config := &ssh.ClientConfig {
		User: username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}
	hostname = fmt.Sprintf("%s:22", hostname)
	conn, err := ssh.Dial("tcp", hostname, config)
	if err != nil {
		log.Fatalf("dial failed:%v", err)
	}

	return conn
}

func RunCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		log.Print(err)
		
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	if err != nil {
		log.Print(err)
	}

	stdout := string(buf.Bytes())

	return stdout, err
}
