package main

import "fmt"
import "net"

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
			conn.Close()
		}
	}
}