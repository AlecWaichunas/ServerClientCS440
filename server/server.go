package main

/*
 package main creates a server that will accept client
 connection. Once a connection is accepted a goroutine
 is made that will listen to the client for a command,
 run the command on the server and send a response back.
*/

import (
	"fmt"
	"net"
	"io"
	"bufio"
	"os/exec"
	"strings"
)

// func main creates a server to listen for clients and
// accepts them on command. It will then run the func
// handleconnection
func main(){
	ln, err := net.Listen("tcp", ":12641")
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

// handleconnection takes a Conn argument from the package
// net. It will begin to listen to the client for a message
// to be read. Once the message is read, it will call func
// runcommand. Then sends the stdout or stderr messages
// to the client
func handleconnection(c net.Conn){
	defer c.Close()
	defer fmt.Printf("Closed Connection\n");

	for c != nil {
		msg := make([]byte, 4096)
		n, err := bufio.NewReader(c).Read(msg)
		if err != nil{
			//End of connection
			if err == io.EOF { break }
			fmt.Printf(err.Error() + "\n")
		}
		//something was read
		if n > 0 {
			readmsg := msg[:n]
			cmderr := false;
			stdout, stderr := runcommand(string(readmsg))
			n := 0;
			if stdout == nil && stderr == nil {
				cmderr = true;
			}else{
				n, err = bufio.NewReader(stdout).Read(msg)
				if err != nil && err == io.EOF{
					//some error occured, test stderr
					n, err = bufio.NewReader(stderr).Read(msg)
					if err != nil {
						fmt.Printf(err.Error())
						cmderr = true;
					}
				}
			}
			if cmderr {
				errmsg := "Could not run command\n"
				msg = []byte(errmsg)
				n = len(errmsg)
			}
			//send to server
			readmsg = msg[:n]
			n, err = c.Write(readmsg)
			if err != nil {
				fmt.Printf("Error writing data to client, %d bytes written", n)
			}
		}
	}
}

// func runcommand takes a string as an argument.
// the string is then fed through the func Command 
// from the exec package. The stdout and errout is
// then returned.
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
		return nil, nil;
	}
	
	return stdout, errout
}