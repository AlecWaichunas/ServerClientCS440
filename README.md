### CS440 Lab 1

Alec Waichunas  
March 7th, 2019  
Prof. Bidyut Gupta  

#### Introduction
The Client Server program requires a server and a client to connect to the server.
The server waits for a client to connect to the it's specified address and port number.
Once a client asks to connect, the server will accept the connection.
The server then waits for the client to send a command over. 
Once it receives a message the server will try and execute the command.
The output of the command, on the server is reutrned.
If the command fails an error message is returned.

For this program, I have decided to use the programming language Go.
Go has a special keyword 'go' that allows easy concurrency within the program.
This allows multiple clients to connect with ease. 
It also has very readable syntax, making the code very straightforward in what the objective is.
Go also compiles into a single binary file, making it executable on many machines just by transfering the compiled file over.

#### Execution
The server program is executed in a CLI. 
Using Bash the execution of the server program looks like below.
The command will only work only if you are in the same directory as the server file.
````bash
./server &
````
The ampersand is to start the program in the background.

The client program is execute also ina CLI.
Using Bash the executtion of the client program looks like below.
````bash
./client localhost:12641
````
localhost represents the server address, this can be replaced by any valid IP address that the server is running on.
The port number is followed by the server address.
If the client succeeds then any line typed will be sent to the server for execution.
The client then waits for a response from the server and prints the output of the command that was ran on the server.  
Outputs of the program:
````bash
ls
client
godef
go-outline
gopkgs
hello
server
````
````bash
pwd
/home/alec/go/bin
````
````bash
who
alec     :0           2019-03-06 21:33 (:0)
````
````bash
ifconfig
eno1: flags=4099<UP,BROADCAST,MULTICAST>  mtu 1500
        ether 98:e7:f4:d9:3d:08  txqueuelen 1000  (Ethernet)
        RX packets 0  bytes 0 (0.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 0  bytes 0 (0.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536
        inet 127.0.0.1  netmask 255.0.0.0
        inet6 ::1  prefixlen 128  scopeid 0x10<host>
        loop  txqueuelen 1000  (Local Loopback)
        RX packets 65794  bytes 5669646 (5.6 MB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 65794  bytes 5669646 (5.6 MB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

wlo1: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 10.88.159.144  netmask 255.255.252.0  broadcast 10.88.159.255
        inet6 fe80::4f16:5721:bd64:57f  prefixlen 64  scopeid 0x20<link>
        ether 84:ef:18:70:62:00  txqueuelen 1000  (Ethernet)
        RX packets 190377  bytes 120341257 (120.3 MB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 108831  bytes 18501414 (18.5 MB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
````
````bash
thisisnotacommand
Could not run command
````


#### The Server Program
The Server Program consists of 3 functions. The main function, handleconnection function and runcommand function.
The main function is the entry point into the program. 
It consists of creating a listening socket and a loop that accepts client connection and calls the handleconnection function.
The handleconnection function is called with a goroutine running the function in a new channel.
````go
for {
 conn, err := ln.Accept()
 if err != nil {
  fmt.Println("Error connecting to client")
 }else{
  fmt.Println("Connected to client")
  go handleconnection(conn)
 }
}
````
The connection is established and the handlconnection function is called, running in a new channel.
The handleconnection function takes the Conn intereface as its only parameter.
The server will now wait until it receives a message from the client. Once it does it will try and run the comand.
````go
msg := make([]byte, 4096)
n, err := bufio.NewReader(c).Read(msg)
.
.
.
readmsg := msg[:n]
cmderr := false;
stdout, stderr := runcommand(string(readmsg))
````
The message was read. Before the command is ran, there are a few errors that are checked and then the runcommand function is called.
Once the runcommand function is finished, the standard output and standard error are returned.
The handleconnection function goes on and checks if the standard output is valid, if not then use the standard error output.
````go
n, err = bufio.NewReader(stdout).Read(msg)
if err != nil && err == io.EOF{
 //some error occured, test stderr
 n, err = bufio.NewReader(stderr).Read(msg)
 if err != nil {
  fmt.Printf(err.Error())
  cmderr = true;
 }
}
````
Once the command was read and the message is received, the extra bytes are sliced off the msg array and sent back to the client.
````go
//send to server
readmsg = msg[:n]
n, err = c.Write(readmsg)
````
The runcommand function is really simple. 
It takes a single string as the parameter and returns the standard output and standard error of the command.
It uses go's exec package to run the command.
````go
stringcmds := strings.Split(s, " ")
cmd := exec.Command(stringcmds[0], stringcmds[1:]...)
stdout, err := cmd.StdoutPipe()
````

#### The Client Program
The Client Program also consists of 3 functions, the main function, writetoserver and readfromserver.
The main function is the entry point of the program.
It reads the first argument listed with command and then tries to connect to the server with that address.
````go
addr := os.Args[1]
conn ,err := net.Dial("tcp", addr)
````
 If connected the writetoserver function is called and takes over the program.  
The wrtietoserver function it takes the go's Conn interface as a parameter.
Once called it goes into a loop until the connection is closed, or the program is exited.
The client program will now halt until a line is read from the standard input. 
If read and no errors occured it is sent to the server.
````go
msg, _, err := bufio.NewReader(os.Stdin).ReadLine()
.
.
.
n, err := conn.Write(msg);
````
Once the message is sent to the server, the readfromserver is function is called.  
The readfromserver function also takes go's Conn intereface as a parameter.
It will wait to receive a message from the server.
When it does it simplies prints it out to the standard output.
````go
msg := make([]byte, 4098)
n, err := bufio.NewReader(conn).Read(msg)
.
.
.
fmt.Printf("%v", string(msg[:n]))
````
#### Error detection within the programs
There are many places within the program that could cause errors. 
Go handles most of the read and write errors by itself.
Making the program very easily configureable.
Go also handles the errors taken by the execution of the commands.
The only extra the program has to do is check if an error did occur and handle it properly.
Handling the error consists of either ending the program or moving on.

#### Known Bugs
There is a bug that was not caught when creating the program. 
When the client is waiting on the server for the response there is no timeout.
This could leave the client in a hanging state. 
To fix this a simple timeout function could be made with a goroutine to timeout after x amount of time.
 
#### Packages
There are a few go packages that were used to create this program.
The net package was used to let the server listen on a port, and for the client to connect to the server.
The os package was used to get the arguments from the program and to get the standard input.
The exec package was used to execute the command.
The bufio package was used to make reading from all inputs very easy.

#### Possible Improvements
When creating the program I thought of it very linearally.
I thought of another way to improve the client, and that would be to make reading from the server in a different channel.
This would make reading and writing concurrent. 
It will also improve the responsiveness and incase if a message will never be received.
Creating these 2 channels would have shared resources of how many messages sent and how many received.
Go would be a good language to do this in as one of the goals of the language was to improve concurrency and defeat race conditions.
 
#### Conclusion
This lab was a good start for me learning go lang. 
I enjoyed looking into all the language has to offer and all of the defects it still has.
The language made connecting a client to a server very easy and to have the server run multiple clients even easier.
