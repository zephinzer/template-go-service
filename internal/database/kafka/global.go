package kafka

import (
	"app/internal/database"
	"fmt"

	"github.com/segmentio/kafka-go"
)

var db *kafka.Conn
var dbMap = map[string]*kafka.Conn{}

func Db(name ...string) *kafka.Conn {
	if len(name) > 0 {
		return dbMap[name[0]]
	}
	return db
}

func Init(opts database.ConnectionOpts, name ...string) error {
	if len(name) > 0 {
		if _, ok := dbMap[name[0]]; ok {
			return fmt.Errorf("failed to replace kafka connection[%s]", name[0])
		}
	}
	conn, err := NewConnection(opts)
	if err != nil {
		return fmt.Errorf("failed to initialise kafka instance: %s", err)
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
	var conn *kafka.Conn
	if len(name) > 0 {
		var ok bool
		conn, ok = dbMap[name[0]]
		if !ok {
			return fmt.Errorf("failed to get kafka connection[%s]", name[0])
		}
		delete(dbMap, name[0])
		connName = name[0]
	} else {
		conn = db
	}
	if err := conn.Close(); err != nil {
		return fmt.Errorf("failed to close kafka connection[%s]", connName)
	}
	return nil
}
