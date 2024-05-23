package main

import (
	"Ex3_Week6/constants"
	"bufio"
	"errors"
	"github.com/fatih/color"
	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestHandleRequest_ErrorReadingUsername(t *testing.T) {
	green := color.New(color.FgGreen).PrintfFunc()
	green("Starting TestHandleRequest_ErrorReadingUsername\n")

	conn, _ := net.Pipe()
	conn.Close()

	err := HandleRequest(conn)
	if !errors.Is(err, io.EOF) {
		t.Errorf("Expected error while reading username, got: %v", err)
		t.FailNow()
	}

	yellow := color.New(color.FgYellow).PrintfFunc()
	yellow("TestHandleRequest_ErrorReadingUsername completed\n")
}

func TestHandleRequest(t *testing.T) {
	log.AddHandler(console.New(true), log.AllLevels...)
	green := color.New(color.FgGreen).PrintfFunc()
	green("Starting TestHandleRequest\n")

	conn, serverConn := net.Pipe()

	go HandleRequest(conn)

	username := "testUser\n"
	_, err := serverConn.Write([]byte(username))
	if err != nil {
		log.Errorf("Error writing username: %s", err)
		t.Fatalf("Error writing username: %s", err)
	}

	message := "Test message\n"
	_, err = serverConn.Write([]byte(message))
	if err != nil {
		log.Errorf("Error writing message: %s", err)
		t.Fatalf("Error writing message: %s", err)
	}
	time.Sleep(100 * time.Millisecond)
	yellow := color.New(color.FgYellow).PrintfFunc()
	yellow("TestHandleRequest completed\n")

}

func TestHandleRequest_MultipleClients(t *testing.T) { // 1
	log.AddHandler(console.New(true), log.AllLevels...)

	green := color.New(color.FgGreen).PrintfFunc()
	green("Starting TestHandleRequest_MultipleClients\n")

	done := make(chan bool)
	messageSent := make(chan bool)

	conn1, serverConn1 := net.Pipe()
	conn2, serverConn2 := net.Pipe()

	defer conn1.Close()
	defer conn2.Close()

	go func() {
		HandleRequest(conn1)
		done <- true
	}()

	go func() {
		username1 := "user1\n"
		_, err := serverConn1.Write([]byte(username1))
		if err != nil {
			log.Errorf("Error writing username for client 1: %s", err)
			t.Fatalf("Error writing username for client 1: %s", err)
		}

		message1 := "Message from user1\n"
		_, err = serverConn1.Write([]byte(message1))
		if err != nil {
			log.Errorf("Error writing message for client 1: %s", err)
			t.Fatalf("Error writing message for client 1: %s", err)
		}

		messageSent <- true
	}()

	<-messageSent

	go HandleRequest(conn2)

	username2 := "user2\n"
	_, err := serverConn2.Write([]byte(username2))
	if err != nil {
		log.Errorf("Error writing username for client 2: %s", err)
		t.Fatalf("Error writing username for client 2: %s", err)
	}

	message2 := "Message from user2\n"
	_, err = serverConn2.Write([]byte(message2))
	if err != nil {
		log.Errorf("Error writing message for client 2: %s", err)
		t.Fatalf("Error writing message for client 2: %s", err)
	}

	time.Sleep(100 * time.Millisecond)
	yellow := color.New(color.FgYellow).PrintfFunc()
	yellow("TestHandleRequest_MultipleClients completed\n")
} // 1
func TestHandleRequest_ClientExiting(t *testing.T) { // 2
	log.AddHandler(console.New(true), log.AllLevels...)

	green := color.New(color.FgGreen).PrintfFunc()
	green("Starting TestHandleRequest_ClientExiting\n")

	conn, serverConn := net.Pipe()

	go HandleRequest(conn)

	username := "testUser\n"
	_, err := serverConn.Write([]byte(username))
	if err != nil {
		log.Errorf("Error writing username: %s", err)
		t.Fatalf("Error writing username: %s", err)
	}

	conn.Close()

	time.Sleep(100 * time.Millisecond)

	yellow := color.New(color.FgYellow).PrintfFunc()
	yellow("TestHandleRequest_ClientExiting completed\n")
} // 2
func TestHandleRequest_LargeMessage(t *testing.T) { // 3
	log.AddHandler(console.New(true), log.AllLevels...)
	green := color.New(color.FgGreen).PrintfFunc()
	green("Starting TestHandleRequest_LargeMessage\n")

	conn, serverConn := net.Pipe()

	go HandleRequest(conn)

	username := "testUser\n"
	_, err := serverConn.Write([]byte(username))
	if err != nil {
		log.Errorf("Error writing username: %s", err)
		t.Fatalf("Error writing message: %s", err)
	}

	largeMessage := strings.Repeat("This is a large message part ", 1000)

	_, err = serverConn.Write([]byte(largeMessage + "\n"))
	if err != nil {
		log.Errorf("Error writing message: %s", err)
		t.Fatalf("Error writing message: %s", err)
	}

	time.Sleep(500 * time.Millisecond)
	yellow := color.New(color.FgYellow).PrintfFunc()
	yellow("TestHandleRequest_LargeMessage completed\n")
} // 3
func TestHandleRequest_InvalidPort(t *testing.T) { // 4
	log.AddHandler(console.New(true), log.AllLevels...)
	green := color.New(color.FgGreen).PrintfFunc()
	green("Starting TestHandleRequest_InvalidPort\n")

	_, err := net.Listen(constants.TYPE, constants.HOST+":0")
	if err == nil {
		t.Errorf("Expected error when listening on invalid port, got nil")
	}

	yellow := color.New(color.FgYellow).PrintfFunc()
	yellow("TestHandleRequest_InvalidPort completed\n")
} // 4
func TestHandleRequest_EmptyUsername(t *testing.T) { // 5
	log.AddHandler(console.New(true), log.AllLevels...)
	green := color.New(color.FgGreen).PrintfFunc()
	green("Starting TestHandleRequest_EmptyUsername\n")

	conn, serverConn := net.Pipe()
	defer conn.Close()
	defer serverConn.Close()

	go HandleRequest(conn)

	_, err := serverConn.Write([]byte(""))
	if err != nil {
		log.Errorf("Error writing username: %s", err)
		t.Fatalf("Error writing username: %s", err)
	}

	message := "Test message\n"
	_, err = serverConn.Write([]byte(message))
	if err != nil {
		log.Errorf("Error writing message: %s", err)
		t.Fatalf("Error writing message: %s", err)
	}

	time.Sleep(100 * time.Millisecond)
	yellow := color.New(color.FgYellow).PrintfFunc()
	yellow("TestHandleRequest_EmptyUsername completed\n")
} // 5
func TestHandleRequest_ConnectionTimeout(t *testing.T) {
	log.AddHandler(console.New(true), log.AllLevels...)
	green := color.New(color.FgGreen).PrintfFunc()
	green("Starting TestHandleRequest_ConnectionTimeout\n")

	conn, _ := net.DialTimeout("tcp", constants.HOST+":"+constants.PORT, 10*time.Millisecond)
	defer conn.Close()

	time.Sleep(100 * time.Millisecond)
	err := conn.Close()
	if err != nil {
		t.Errorf("Unexpected error closing connection: %v", err)
	}

	yellow := color.New(color.FgYellow).PrintfFunc()
	yellow("TestHandleRequest_ConnectionTimeout completed\n")
} // 6
func TestHandleRequest_LargeUsername(t *testing.T) { // 8
	log.AddHandler(console.New(true), log.AllLevels...)
	green := color.New(color.FgGreen).PrintfFunc()
	green("Starting TestHandleRequest_LargeUsername\n")

	conn, serverConn := net.Pipe()
	defer conn.Close()
	defer serverConn.Close()

	go HandleRequest(conn)

	largeUsername := strings.Repeat("x", constants.MaxUsernameLength+1) + "\n"

	_, err := serverConn.Write([]byte(largeUsername))
	if err == nil {
		t.Errorf("Expected error when writing large username, got nil")
	} else if !strings.Contains(err.Error(), "username too large") {
		t.Errorf("Unexpected error: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	yellow := color.New(color.FgYellow).PrintfFunc()
	yellow("TestHandleRequest_LargeUsername completed\n")
} // 7
func TestHandleRequest_DisconnectedServer(t *testing.T) {
	log.AddHandler(console.New(true), log.AllLevels...)
	green := color.New(color.FgGreen).PrintfFunc()
	green("Starting TestHandleRequest_DisconnectedServer\n")

	conn, serverConn := net.Pipe()
	defer conn.Close()

	go HandleRequest(conn)

	username := "testUser\n"
	_, err := serverConn.Write([]byte(username))
	if err != nil {
		log.Errorf("Error writing username: %s", err)
		t.Fatalf("Error writing username: %s", err)
	}

	serverConn.Close()

	message := "Test message\n"
	_, err = conn.Write([]byte(message))
	if err == nil {
		t.Errorf("Expected error writing to disconnected server, got nil")
	}

	if !strings.Contains(err.Error(), "write: connection closed") && !strings.Contains(err.Error(), "use of closed network connection") {
		t.Errorf("Unexpected error: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	yellow := color.New(color.FgYellow).PrintfFunc()
	yellow("TestHandleRequest_DisconnectedServer completed\n")
} // 8
func TestHandleRequest_ForbiddenWords(t *testing.T) {
	log.AddHandler(console.New(true), log.AllLevels...)
	green := color.New(color.FgGreen).PrintfFunc()
	green("Starting TestHandleRequest_ForbiddenWords\n")

	conn, serverConn := net.Pipe()
	defer conn.Close()
	defer serverConn.Close()

	go HandleRequest(conn)

	username := "testUser\n"
	_, err := serverConn.Write([]byte(username))
	if err != nil {
		log.Errorf("Error writing username: %s", err)
		t.Fatalf("Error writing username: %s", err)
	}

	message := "This message contains alcohol\n"
	_, err = serverConn.Write([]byte(message))
	if err != nil {
		log.Errorf("Error writing message: %s", err)
		t.Fatalf("Error writing message: %s", err)
	}

	response, err := bufio.NewReader(serverConn).ReadString('\n')
	if err != nil {
		t.Fatalf("Error reading response: %s", err)
	}

	if !strings.Contains(response, "Error: Message contains forbidden words") {
		t.Errorf("Expected forbidden word error message, got: %s", response)
	}

	time.Sleep(100 * time.Millisecond)
	yellow := color.New(color.FgYellow).PrintfFunc()
	yellow("TestHandleRequest_ForbiddenWords completed\n")
} // 9
func TestHandleRequest_InvalidMessageFormat(t *testing.T) {
	log.AddHandler(console.New(true), log.AllLevels...)
	green := color.New(color.FgGreen).PrintfFunc()
	green("Starting TestHandleRequest_InvalidMessageFormat\n")

	conn, serverConn := net.Pipe()
	defer conn.Close()
	defer serverConn.Close()

	go HandleRequest(conn)

	username := "testUser\n"
	_, err := serverConn.Write([]byte(username))
	if err != nil {
		log.Errorf("Error writing username: %s", err)
		t.Fatalf("Error writing username: %s", err)
	}

	invalidMessage := "This is an incomplete message"
	_, err = serverConn.Write([]byte(invalidMessage)) // Message without newline
	if err != nil {
		log.Errorf("Error writing invalid message: %s", err)
		t.Fatalf("Error writing invalid message: %s", err)
	}

	response, err := bufio.NewReader(serverConn).ReadString('\n')
	if err != nil {
		t.Fatalf("Error reading response: %s", err)
	}

	if !strings.Contains(response, "Error: Invalid message format") {
		t.Errorf("Expected invalid message format error, got: %s", response)
	}

	time.Sleep(100 * time.Millisecond)
	yellow := color.New(color.FgYellow).PrintfFunc()
	yellow("TestHandleRequest_InvalidMessageFormat completed\n")
} //10
