package tracker

import (
	"time"
)

type Session struct {
	AppName    string
	Title      string
	URL        string
	CategoryID string
	StartedAt  time.Time
	EndedAt    time.Time
}

type Tracker struct {
	CurrentSession *Session
	OnSessionEnd   func(s *Session)
}

func NewTracker(onSessionEnd func(s *Session)) *Tracker {
	return &Tracker{
		OnSessionEnd: onSessionEnd,
	}
}

func (t *Tracker) Update(appName, title, url string) {
	now := time.Now()

	// If no current session, start one
	if t.CurrentSession == nil {
		t.CurrentSession = &Session{
			AppName:   appName,
			Title:     title,
			URL:       url,
			StartedAt: now,
		}
		return
	}

	// If window/app changed, close current and start new
	if t.CurrentSession.AppName != appName || t.CurrentSession.Title != title || t.CurrentSession.URL != url {
		t.CurrentSession.EndedAt = now
		if t.OnSessionEnd != nil {
			t.OnSessionEnd(t.CurrentSession)
		}

		t.CurrentSession = &Session{
			AppName:   appName,
			Title:     title,
			URL:       url,
			StartedAt: now,
		}
	}
}

func (t *Tracker) Pause() {
	if t.CurrentSession != nil {
		t.CurrentSession.EndedAt = time.Now()
		if t.OnSessionEnd != nil {
			t.OnSessionEnd(t.CurrentSession)
		}
		t.CurrentSession = nil
	}
}
