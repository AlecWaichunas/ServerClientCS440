package main

/*
 package main to start the client program. The client connects
 to a given server in the command e.g. client localhost. The
 client will connect to the running server on the addressed
 machine. Once connected, typing commands through the standard
 input will send them to the server and reply back with a
 response.
*/

import (
	"fmt"
	"net"
	"os"
	"bufio"
)
// func main connect to given ip address argument and starts 
// the process to write to the server

func main(){
	conn ,err := net.Dial("tcp", "localhost:12641")
	if err != nil {
		fmt.Println("Could not connect to server")
		return
	}
	writetoserver(conn)
}

// func readfromserver takes a Conn argument from the package
// net. It reads the message that the server has to send.
func readfromserver(conn net.Conn) (exit bool){
	msg := make([]byte, 4098)
	n, err := bufio.NewReader(conn).Read(msg)
	if err != nil {
		fmt.Printf("Error Reading from server\n")
		exit = true;
	}
	if n > 0 {
		fmt.Printf("%v", string(msg[:n]))
		exit = false;
	}
	return exit
}

// func writetoserver takes a Conn argument from the package
// net. It will listen to standard input and send the line to
// the server. The func then waits to read from the server.
func writetoserver(conn net.Conn){
	for conn != nil {
		msg, _, err := bufio.NewReader(os.Stdin).ReadLine()
		if err != nil{
			//error reading from console
			fmt.Printf("Error reading command")
		}else{
			//write and read from server
			n, err := conn.Write(msg);
			if err != nil || n == 0 {
				fmt.Printf("Could not write to server\n")
				break;
			}
			if readfromserver(conn){
				break;
			}
		}
	}
}