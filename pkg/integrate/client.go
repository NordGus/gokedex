package integrate

type IntegrationSecret string

type DatabaseID string

type client struct {
	conn       Connection
	baseUrl    string
	secret     IntegrationSecret
	databaseID DatabaseID
}

func newClient(conn Connection, secret IntegrationSecret, databaseID DatabaseID) client {
	return client{
		conn:       conn,
		baseUrl:    "https://api.notion.com/v1",
		secret:     secret,
		databaseID: databaseID,
	}
}
