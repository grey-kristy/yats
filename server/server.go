package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const dbFile = "db.log"

func Server(port string) {
	db := newDB(dbFile)
	go db.run()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("failed to create listener, err:", err)
		os.Exit(1)
	}
	fmt.Printf("listening on %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("failed to accept connection, err:", err)
			continue
		}
		go handleConnection(conn, db)
	}
}

func handleConnection(conn net.Conn, db *DB) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		s, err := reader.ReadString(byte('\n'))
		if err != nil {
			if err != io.EOF {
				fmt.Println("failed to read data, err:", err)
			}
			return
		}
		num, err := strconv.ParseInt(strings.TrimSuffix(s, "\n"), 10, 32)
		if err != nil {
			fmt.Println("ParseInt err:", err)
			conn.Write([]byte(fmt.Sprintf("error %s\n", err)))
			fmt.Println("send to connect")
			continue
		}
		buf := make([]byte, num)
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			fmt.Println("ReadFull err:", err)
			conn.Write([]byte(fmt.Sprintf("error %s\n", err)))
			continue
		}
		fmt.Printf("server recived %d bytes: %s...\n", num, string(buf[0:15]))
		res := db.save(buf)
		conn.Write([]byte(fmt.Sprintf("%s\n", res)))
	}
}
