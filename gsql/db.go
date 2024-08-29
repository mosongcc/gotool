package gsql

type DB struct {
}

func (db *DB) Update(set, where any) (err error) {

	return
}

func (db *DB) Delete(where any) (err error) {

	return
}

func (db *DB) Select(where SelectBuild) (err error) {
	// db.Select(...).From(table).where(field).Eq("").OrderBy(field,desc).Limit(3)
	return
}
