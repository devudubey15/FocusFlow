package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

type Msg struct {
	URL       string `json:"url"`
	Timestamp string `json:"timestamp"`
}

func main() {
	// Log to a file for debugging NMH (since it runs in background via browser)
	home, _ := os.UserHomeDir()
	logPath := filepath.Join(home, ".config", "focusflow", "nmh.log")
	_ = os.MkdirAll(filepath.Dir(logPath), 0755)
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		log.SetOutput(f)
		defer f.Close()
	}

	log.Println("NMH started")

	for {
		// Read message length (4 bytes)
		var length uint32
		err := binary.Read(os.Stdin, binary.LittleEndian, &length)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading length: %v", err)
			}
			break
		}

		// Read message
		msgBytes := make([]byte, length)
		_, err = io.ReadFull(os.Stdin, msgBytes)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received: %s", string(msgBytes))

		// Forward to main agent via Unix socket
		forwardToAgent(msgBytes)
	}
}

func forwardToAgent(data []byte) {
	home, _ := os.UserHomeDir()
	socketPath := filepath.Join(home, ".config", "focusflow", "agent.sock")

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Printf("Error connecting to agent socket: %v", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		log.Printf("Error writing to agent socket: %v", err)
	}
}
