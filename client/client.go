package main

import (
	"Ex3_Week6/constants"
	"bufio"
	"fmt"
	"github.com/go-playground/log"
	"net"
	"os"
	"strings"
	"time"
)

type Client struct {
	conn       *net.TCPConn
	reader     *bufio.Reader
	username   string
	historyLog []string
}

func NewClient(conn *net.TCPConn, reader *bufio.Reader, username string) *Client {
	return &Client{
		conn:       conn,
		reader:     reader,
		username:   username,
		historyLog: make([]string, 0),
	}
}

func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	if err != nil {
		return err
	}
	c.historyLog = append(c.historyLog, message)
	return nil
}

func (c *Client) DisplayHistory() {
	fmt.Println("=== Message History ===")
	for _, msg := range c.historyLog {
		fmt.Println(msg)
	}
	fmt.Println("======================")
}

func main() {
	tcpServer, err := net.ResolveTCPAddr(constants.TYPE, constants.HOST+":"+constants.PORT)
	if err != nil {
		fmt.Println("ResolveTCPAddr failed:", err)
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		fmt.Println("Dial failed:", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')

	client := NewClient(conn, reader, username)

	err = client.Send(username)
	if err != nil {
		fmt.Println("Write data failed:", err)
		os.Exit(1)
	}

	for {
		fmt.Print("Enter message: ")
		text, _ := reader.ReadString('\n')

		if strings.TrimSpace(text) == "HISTORY" {
			client.DisplayHistory()
			continue
		}

		err = client.Send(text)
		if err != nil {
			fmt.Println("Write data failed:", err)
			os.Exit(1)
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.WithError(err).Error("Read response failed")
			os.Exit(1)
		}
		fmt.Print(response)

		if strings.TrimSpace(text) == "EXIT" {
			fmt.Println("Exiting...")
			conn.Close()
			os.Exit(0)
		}

		time.Sleep(1 * time.Second)
	}
}
