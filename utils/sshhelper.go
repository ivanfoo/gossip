package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"golang.org/x/crypto/ssh"
)

func SSHConnect(username string, hostname string, keyPath string) string {
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		log.Fatalf("parse key failed:%v", err)
	}
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}
	hostname = fmt.Sprintf("%s:22", hostname)
	conn, err := ssh.Dial("tcp", hostname, config)
	if err != nil {
		log.Fatalf("dial failed:%v", err)
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("session failed:%v", err)
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run("ls -l")
	if err != nil {
		log.Fatalf("Run failed:%v", err)
	}
	//log.Printf(">%s", stdoutBuf)
	return stdoutBuf.String()
}

func RunCommand(client *ssh.Client, command string) (stdout string, err error) {
	session, err := client.NewSession()
	if err != nil {
		//log.Print(err)
		return
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	if err != nil {
		//log.Print(err)
		return
	}
	stdout = string(buf.Bytes())

	return
}
