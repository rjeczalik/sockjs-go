package sockjs

import "net/http"

// Session represents a connection between server and client.
type Session interface {
	// Id returns a session id
	ID() string

	// Request returns the first http request
	Request() *http.Request

	// Recv reads one text frame from session
	Recv() (string, error)

	// Send sends one text frame to session
	Send(string) error

	// Close closes the session with provided code and reason.
	Close(status uint32, reason string) error

	// GetSessionState gives the state of the session.
	//
	// SessionOpening/SessionActive/SessionClosing/SessionClosed;
	GetSessionState() SessionState

	// Done returns channel which gets closed when session is closed.
	//
	// If it returns nil channel, it means no lifetime info is available.
	Done() <-chan struct{}

	// Err gives non-nil error if session was closed due to an error.
	Err() error
}
