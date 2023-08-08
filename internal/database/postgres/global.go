package postgres

import (
	"app/internal/database"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn
var dbMap = map[string]*pgx.Conn{}

func Db(name ...string) *pgx.Conn {
	if len(name) > 0 {
		return dbMap[name[0]]
	}
	return db
}

func Init(opts database.ConnectionOpts, name ...string) error {
	if len(name) > 0 {
		if _, ok := dbMap[name[0]]; ok {
			return fmt.Errorf("failed to replace postgres connection[%s]", name[0])
		}
	}
	conn, err := NewConnection(opts)
	if err != nil {
		return fmt.Errorf("failed to initialise postgres instance: %s", err)
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
	var conn *pgx.Conn
	if len(name) > 0 {
		var ok bool
		conn, ok = dbMap[name[0]]
		if !ok {
			return fmt.Errorf("failed to get postgres connection[%s]", name[0])
		}
		delete(dbMap, name[0])
		connName = name[0]
	} else {
		conn = db
	}
	closeCtx, cancelClose := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelClose()
	if err := conn.Close(closeCtx); err != nil {
		return fmt.Errorf("failed to close postgres connection[%s]", connName)
	}
	return nil
}
