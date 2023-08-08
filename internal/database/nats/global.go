package nats

import (
	"app/internal/database"
	"fmt"

	"github.com/nats-io/nats.go"
)

var db *nats.Conn
var dbMap = map[string]*nats.Conn{}

func Db(name ...string) *nats.Conn {
	if len(name) > 0 {
		return dbMap[name[0]]
	}
	return db
}

func Init(opts database.ConnectionOpts, name ...string) error {
	if len(name) > 0 {
		if _, ok := dbMap[name[0]]; ok {
			return fmt.Errorf("failed to replace nats connection[%s]", name[0])
		}
	}
	conn, err := NewConnection(opts)
	if err != nil {
		return fmt.Errorf("failed to initialise nats instance: %s", err)
	}
	if len(name) > 0 {
		dbMap[name[0]] = conn
	} else {
		db = conn
	}
	return nil
}

func Close(name ...string) error {
	connName := "default"
	var conn *nats.Conn
	if len(name) > 0 {
		var ok bool
		conn, ok = dbMap[name[0]]
		if !ok {
			return fmt.Errorf("failed to get nats connection[%s]", name[0])
		}
		delete(dbMap, name[0])
		connName = name[0]
	} else {
		conn = db
	}
	if conn.IsConnected() && !conn.IsDraining() && !conn.IsClosed() {
		if err := conn.Drain(); err != nil {
			return fmt.Errorf("failed to drain nats connection[%s]", connName)
		}
		conn.Close()
	}
	return nil
}
