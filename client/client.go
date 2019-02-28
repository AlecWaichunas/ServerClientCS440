package main

import "fmt"
import "net"

func main(){
	_,err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Could not connect to server")
		return
	}
	fmt.Println("Connected to server")
}