package integrate

import "net/http"

type Connection interface {
	Do(*http.Request) (*http.Response, error)
	CloseIdleConnections()
}
