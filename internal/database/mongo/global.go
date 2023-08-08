package mongo

import (
	"app/internal/database"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Client
var dbMap = map[string]*mongo.Client{}

func Db(name ...string) *mongo.Client {
	if len(name) > 0 {
		return dbMap[name[0]]
	}
	return db
}

func Init(opts database.ConnectionOpts, name ...string) error {
	if len(name) > 0 {
		if _, ok := dbMap[name[0]]; ok {
			return fmt.Errorf("failed to replace mongo connection[%s]", name[0])
		}
	}
	conn, err := NewConnection(opts)
	if err != nil {
		return fmt.Errorf("failed to initialise mongo instance: %s", err)
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
	var conn *mongo.Client
	if len(name) > 0 {
		var ok bool
		conn, ok = dbMap[name[0]]
		if !ok {
			return fmt.Errorf("failed to get mongo connection[%s]", name[0])
		}
		delete(dbMap, name[0])
		connName = name[0]
	} else {
		conn = db
	}
	disconnectCtx, cancelDisconnect := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelDisconnect()
	if err := conn.Disconnect(disconnectCtx); err != nil {
		return fmt.Errorf("failed to disconnect mongo connection[%s]", connName)
	}
	return nil
}
