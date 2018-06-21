package anondb

import (
	"fmt"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var databaseInitScripts = [...]func(*mgo.Database) error{
	createDatasetsCollection,
}

// SetupDatabase initializes the database (creates collections, etc.)
func SetupDatabase() error {
	return setupDatabase(globalSession)
}

func setupDatabase(session *mgo.Session) error {
	db := session.DB("anondb")

	system := db.C("system")
	dbVersion := getCurrentDatabaseVersion(system)
	currentVersion := len(databaseInitScripts)

	if dbVersion > currentVersion {
		return fmt.Errorf("The database version (%v) is higher than the current version (%v)", dbVersion, currentVersion)
	} else if dbVersion == currentVersion {
		return nil
	} else {
		log.Printf("The database version (%v) is lower than the current verison (%v). Running init scripts...", dbVersion, currentVersion)
		return runDatabaseInitScripts(db, system, dbVersion)
	}
}

func getCurrentDatabaseVersion(system *mgo.Collection) int {
	var versionInfo bson.M
	if err := system.FindId(1).One(&versionInfo); err != nil {
		return 0
	}

	return versionInfo["version"].(int)
}

func runDatabaseInitScripts(db *mgo.Database, system *mgo.Collection, dbVersion int) error {
	for ix, script := range databaseInitScripts {
		// skip already run scripts
		if ix < dbVersion {
			continue
		}

		// try to run the upgrade script
		if err := script(db); err != nil {
			return err
		}

		// set the version number in the database accordingly
		if _, err := system.UpsertId(1, bson.M{"version": ix + 1}); err != nil {
			return err
		}

		log.Printf("Database successfully updated to version %v", ix+1)
	}

	// every script run successfully
	return nil
}
