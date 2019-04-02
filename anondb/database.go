package anondb

import (
	"time"

	"github.com/globalsign/mgo"
)

var globalSession *mgo.Session

// InitConnection initializes the database connection
func InitConnection() (err error) {
	dialInfo := &mgo.DialInfo{
		Addrs: []string{"localhost:27017"},
		// Addrs:    []string{"anonymization_database:27017"},
		Timeout:  time.Second * 30,
		Database: "anondb",
		AppName:  "anonymization_server",
	}

	globalSession, err = mgo.DialWithInfo(dialInfo)

	return
}

// CloseConnection closes the connection to the database
func CloseConnection() {
	globalSession.Close()
}

// Ping checks the connection to the database
func Ping() error {
	session := globalSession.Clone()
	defer session.Close()

	session.SetSocketTimeout(time.Second * 5)
	session.SetSyncTimeout(time.Second * 5)
	return session.Ping()
}
