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
	defer c.Close()
	defer fmt.Printf("Closed Connection\n");
	msg := make([]byte, 4096)

	for c != nil {
		n, err := bufio.NewReader(c).Read(msg)
		if err != nil{
			//End of connection
			if err == io.EOF { break }
			fmt.Printf(err.Error() + "\n")
		}
		//something was read
		if n > 0 {
			readmsg := msg[:n]
			stdout, stderr := runcommand(string(readmsg))
			n, err := bufio.NewReader(stdout).Read(msg)
			if err != nil && err == io.EOF{
				//some error occured, test stderr
				n, err = bufio.NewReader(stderr).Read(msg)
				if err != nil {
					//handle other errors
				}
			}
			if n > 0 {
				//send to server
				readmsg = msg[:n]
				n, err = c.Write(msg[:n])
				if err != nil {
					fmt.Printf("Error writing data to client, %d bytes written", n)
				}
			}
		}
	}
	//stdout, stderr := runcommand("ls -a ~/go")

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
		//command does not exist
		fmt.Printf("Error starting command: %s\n", s)
	}
	
	return stdout, errout
}