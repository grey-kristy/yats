package client

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

const maxRetry = 10
const dataSize = 4*1024

func Client(port string, n int) {
	conn, err := net.Dial("tcp", port)
	defer conn.Close()

	if err != nil {
		fmt.Println("failed to connect to server ", port, ", err:", err)
		os.Exit(1)
	}
	for i := 0; i < n; i++ {
		data := cookData(dataSize)
		sureSend(conn, data)
	}
}

func sureSend(conn net.Conn, data []byte) {
	for i := 0; i < maxRetry; i++ {
		if sendData(conn, data) {
			return
		}
	}
}

func sendData(conn net.Conn, data []byte) bool {
	fmt.Printf("sending %d bytes of data: %s...\n", len(data), string(data[0:15]))
	fmt.Fprintf(conn, "%d\n", len(data))
	conn.Write(data)
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("read error: %s\n", err)
		return false
	}
	if msg == "ok\n" {
		fmt.Println("sended ok")
		return true
	}
	fmt.Println("respond error: ", msg)
	return false
}

func cookData(n int) []byte {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Int63()%int64(len(letterBytes))]
	}
	return b
}
