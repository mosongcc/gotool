package gsql

import (
	"database/sql"
	"testing"
)

func TestDB_Tx(t *testing.T) {
	db, err := Open("", "")
	if err != nil {
		return
	}
	_, err = db.Tx(func(tx *sql.Tx) (any, error) {

		//Update("", "").Where().And().And().OrderBy().Limit().ExecTx(tx)

		return nil, nil
	})
}
