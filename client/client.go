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
	}

	fmt.Println("Connected to server")
}