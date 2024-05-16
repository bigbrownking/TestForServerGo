package client

import (
	"Ex3_Week6/constants"
	"bufio"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
)

var (
	clients        = make(map[net.Conn]bool)
	mu             sync.Mutex
	forbiddenWords map[string]struct{}
)

func init() {
	forbiddenWords = map[string]struct{}{
		"alcohol":   {},
		"war":       {},
		"crime":     {},
		"terrorism": {},
		"death":     {},
	}
}

func main() {
	log.AddHandler(console.New(true), log.AllLevels...)

	listen, err := net.Listen(constants.TYPE, constants.HOST+":"+constants.PORT)
	if err != nil {
		log.WithError(err).Error("error starting server")
		os.Exit(1)
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.WithError(err).Error("error accepting connection")
			os.Exit(1)
		}
		mu.Lock()
		clients[conn] = true
		mu.Unlock()
		go HandleRequest(conn)
	}
}

func HandleRequest(conn net.Conn) error {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	username, err := reader.ReadString('\n')
	if err != nil {
		log.WithError(err).Error("error reading from client")
		return err
	}
	if len(username) > constants.MaxUsernameLength+1 {
		log.Errorf("Username '%s' is too long (max %d characters)", username, constants.MaxUsernameLength)
		_, err = conn.Write([]byte("Error: Username too long\n"))
		if err != nil {
			log.WithError(err).Error("Error sending error message to client")
		}
		return err
	}

	log.Infof("New client connected: %s", username)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Infof("%s disconnected.", username)
			break
		}

		if isMessageForbidden(message) {
			log.Infof("Message from %s contains forbidden word: %s", username, message)
			_, err = conn.Write([]byte("Error: Message contains forbidden words\n"))
			if err != nil {
				log.WithError(err).Error("Error sending error message to client")
			}
			return err
		}
		log.Infof("Received message from %s: %s", username, message)

		_, err = conn.Write([]byte("Server received the message\n"))
		if err != nil {
			log.WithError(err).Error("Error sending confirmation message to client")
			return err
		}
	}
	return nil
}
func isMessageForbidden(message string) bool {
	for word := range forbiddenWords {
		if strings.Contains(strings.ToLower(message), word) {
			return true
		}
	}
	return false
}
