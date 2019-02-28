package main

import (
	"fmt"
	"net"
	"io"
	"bufio"
	"os/exec"
	"strings"
)

func main(){
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error creating server")
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error connecting to client")
		}else{
			fmt.Println("Connected to client")
			go handleconnection(conn)
		}
	}
}

func handleconnection(c net.Conn){
	readbytes := make([]byte, 4096)
	n,err := bufio.NewReader(c).Read(readbytes)
	if err != nil || n >= 4096 {
		//handle error or data too big for a page
	}
	fmt.Printf("%q\n", readbytes[:n])
	//stdout, stderr := runcommand("ls -a ~/go")

	c.Close()
}

func runcommand(s string) (stdout io.ReadCloser, errout io.ReadCloser){
	stringcmds := strings.Split(s, " ")
	cmd := exec.Command(stringcmds[0], stringcmds[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error getting standard output pipe for exec")
	}
	errout, err = cmd.StderrPipe()
	if err != nil {
		fmt.Print("Error getting standard error pipe for exec")
	}
	if err = cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %s\n", s)
	}
	
	return stdout, errout
}