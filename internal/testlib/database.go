package testlib

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func FlushDB(tx *sqlx.DB) error {
	var tables []string
	err := tx.Select(&tables, "SHOW TABLES")
	if err != nil {
		return err
	}

	for _, table := range tables {
		tx, err := tx.Begin()
		if err != nil {
			return err
		}

		_, err = tx.Exec("SET FOREIGN_KEY_CHECKS=0;")
		if err != nil {
			return err
		}
		q := fmt.Sprintf("TRUNCATE TABLE %s", table)
		_, err = tx.Exec(q)
		if err != nil {
			return err
		}
		_, err = tx.Exec("SET FOREIGN_KEY_CHECKS=1;")
		if err != nil {
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func DropTable(tx *sqlx.DB) error {
	var tables []string
	err := tx.Select(&tables, "SHOW TABLES")
	if err != nil {
		return err
	}

	for _, table := range tables {
		tx, err := tx.Begin()
		if err != nil {
			return err
		}

		_, err = tx.Exec("SET FOREIGN_KEY_CHECKS=0;")
		if err != nil {
			return err
		}
		q := fmt.Sprintf("DROP TABLE %s", table)
		_, err = tx.Exec(q)
		if err != nil {
			return err
		}
		_, err = tx.Exec("SET FOREIGN_KEY_CHECKS=1;")
		if err != nil {
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	fmt.Sprintf("Database table dropped")
	return nil
}

func FlushDBAndDropTable(tx *sqlx.DB) error {
	err := FlushDB(tx)
	if err != nil {
		return err
	}
	err = DropTable(tx)
	return err
}
