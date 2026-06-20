package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"focusflow-agent/config"
	"focusflow-agent/internal/cache"
	"focusflow-agent/internal/ipc"
	"focusflow-agent/internal/supabase"
	"focusflow-agent/internal/tracker"
	"focusflow-agent/internal/watcher"
)

const (
	PollInterval  = 5 * time.Second
	IdleThreshold = 60 * time.Second
	SyncInterval  = 30 * time.Second
)

func main() {
	fmt.Println("FocusFlow Agent starting...")

	// Find config path: default to ~/.config/focusflow/config.json
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "focusflow", "config.json")
	if envPath := os.Getenv("FOCUSFLOW_CONFIG"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("Warning: Could not load config from %s: %v. Running in log-only mode.", configPath, err)
	}

	// Initialize Supabase Client
	var sbClient *supabase.Client
	if cfg != nil {
		sbClient = supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseKey)
	}

	// Initialize SQLite Store
	dbPath := filepath.Join(home, ".config", "focusflow", "buffer.db")
	_ = os.MkdirAll(filepath.Dir(dbPath), 0755)
	store, err := cache.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatalf("Critical: Could not initialize SQLite store: %v", err)
	}

	// Initialize tracker
	var currentUrl string
	t := tracker.NewTracker(func(s *tracker.Session) {
		fmt.Printf("[SESSION END] %s | %s | %s | Duration: %v | URL: %s\n",
			s.AppName, s.Title, s.StartedAt.Format("15:04:05"), s.EndedAt.Sub(s.StartedAt), s.URL)

		// Save to local cache first
		if err := store.SaveLog(s.AppName, s.Title, s.URL, s.StartedAt, s.EndedAt); err != nil {
			log.Printf("Error saving to local cache: %v", err)
		}
	})

	// Start IPC Socket Server for Browser URL updates
	socketPath := filepath.Join(home, ".config", "focusflow", "agent.sock")
	err = ipc.StartSocketServer(socketPath, func(update ipc.URLUpdate) {
		fmt.Printf("[IPC] Received URL update: %s\n", update.URL)
		currentUrl = update.URL
	})
	if err != nil {
		log.Printf("Warning: Could not start IPC server: %v", err)
	}

	// Background Sync Process
	if sbClient != nil {
		go func() {
			syncTicker := time.NewTicker(SyncInterval)
			for range syncTicker.C {
				logs, err := store.GetPendingLogs()
				if err != nil {
					log.Printf("Error fetching pending logs: %v", err)
					continue
				}

				if len(logs) > 0 {
					fmt.Printf("[SYNC] Attempting to sync %d logs to Supabase...\n", len(logs))
				}

				for _, l := range logs {
					err := sbClient.InsertLog(supabase.ActivityLog{
						UserID:    "", // Supabase handles this via auth.uid() or we can pass if needed
						DeviceID:  cfg.DeviceId,
						AppName:   l.AppName,
						Title:     l.Title,
						URL:       l.URL,
						StartedAt: l.StartedAt,
						EndedAt:   &l.EndedAt,
					})
					if err == nil {
						_ = store.DeleteLog(l.ID)
					} else {
						log.Printf("Sync error for log %d: %v", l.ID, err)
						break // Stop syncing for now if we lose connection
					}
				}
			}
		}()
	}

	// Setup signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	fmt.Println("Polling active window every", PollInterval)

	for {
		select {
		case <-sigs:
			fmt.Println("\nShutting down...")
			t.Pause()
			return
		case <-ticker.C:
			// Check idle
			idle, err := watcher.IsIdle(IdleThreshold)
			if err != nil {
				log.Printf("Idle detection error: %v", err)
			}

			if idle {
				if t.CurrentSession != nil {
					fmt.Println("[IDLE] Pausing tracking")
					t.Pause()
				}
				continue
			}

			// Get active window
			win, err := watcher.GetActiveWindowX11()
			if err != nil {
				// Don't log error every 5s if it's just "no active window"
				continue
			}

			// Update tracker
			t.Update(win.AppName, win.Title, currentUrl)
		}
	}
}
