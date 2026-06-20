package ipc

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

type URLUpdate struct {
	URL       string `json:"url"`
	Timestamp string `json:"timestamp"`
}

func StartSocketServer(path string, onUpdate func(URLUpdate)) error {
	// Remove existing socket if any
	_ = os.Remove(path)

	l, err := net.Listen("unix", path)
	if err != nil {
		return err
	}

	go func() {
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Printf("Socket accept error: %v", err)
				continue
			}

			go handleConnection(conn, onUpdate)
		}
	}()

	return nil
}

func handleConnection(conn net.Conn, onUpdate func(URLUpdate)) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	var update URLUpdate
	if err := decoder.Decode(&update); err != nil {
		log.Printf("Socket decode error: %v", err)
		return
	}
	onUpdate(update)
}
