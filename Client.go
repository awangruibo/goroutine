package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("客户端 conn err :", err)
		return
	}

	defer conn.Close()

	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := inputReader.ReadString('\n')
		inputinfo := strings.Trim(input, "\r\n")
		_, err := conn.Write([]byte(inputinfo))
		if err != nil {
			return
		}

		var buf [512]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("recv failed, err:", err)
			return
		}

		fmt.Println(string(buf[:n]))

	}

}
