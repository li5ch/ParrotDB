package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var host string
var port int

func init() {
	flag.StringVar(&host, "h", "localhost", "host")
	flag.IntVar(&port, "p", 18888, "port")
}

func main() {

	flag.Parse()
	tcpAddr := &net.TCPAddr{IP: net.ParseIP(host), Port: port}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	// 关闭连接
	defer conn.Close()
	// 键入数据
	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s:%d>", host, port)
		input, _ := inputReader.ReadString('\n')
		input += "\r\n"
		_, err = conn.Write([]byte(input))
		if err != nil {
			return
		}
		buf := [512]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("conn.Read error : ", err)
			return
		}
		fmt.Println(string(buf[:n]))
	}
}
