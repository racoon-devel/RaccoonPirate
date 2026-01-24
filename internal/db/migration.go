package db

import (
	"fmt"
)

type migrateFn func(databaseInternal) error

var migrations = map[uint]migrateFn{
	1: migrateV1toV2,
}

func migrateDatabase(dbase databaseInternal, version uint) error {
	if version > currentDatabaseVersion {
		return fmt.Errorf("database created by the newest version: %d", version)
	}
	if version < 1 {
		return fmt.Errorf("database created by the unknown version: %d", version)
	}

	for nextVersion := version; version < currentDatabaseVersion; version++ {
		if err := migrations[nextVersion](dbase); err != nil {
			return fmt.Errorf("migrate from %d to %d failed: %w", nextVersion, nextVersion+1, err)
		}
	}

	return nil
}

func migrateV1toV2(dbase databaseInternal) error {
	// All necessary stuff have done on the application update stage

	return dbase.SetDatabaseVersion(2)
}
