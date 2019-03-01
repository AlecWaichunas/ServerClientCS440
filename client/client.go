package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
)

func main(){
	conn ,err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Could not connect to server")
		return
	}

	writetoserver(conn)
}

func readfromserver(conn net.Conn){
	msg := make([]byte, 4098)
	n, err := bufio.NewReader(conn).Read(msg)
	if err != nil {
		fmt.Printf("Error Reading from server\n")
	}

	if n > 0 {
		fmt.Printf("%v", string(msg[:n]))
	}
}

func writetoserver(conn net.Conn){
	var isfull bool = true
	for isfull || conn != nil {
		msg, isfull, err := bufio.NewReader(os.Stdin).ReadLine()
		//write more
		if isfull && err != nil{
			//do something
		}
		n, err := conn.Write(msg);
		if err != nil || n == 0 {
			fmt.Printf("Could not write to server\n")
		}
		readfromserver(conn)
	}
}