package integrate

import "net/http"

type Connection interface {
	Get(string) (*http.Response, error)
	CloseIdleConnections()
}
