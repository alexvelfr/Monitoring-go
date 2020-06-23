package models

import "time"

//Reglament ...
type Reglament struct {
	ID          int       `db:"id"`
	Block       string    `db:"block"`
	LastUpdated time.Time `db:"last_updated"`
	InReglament bool      `db:"in_reglament"`
}

func (r *Reglament) GetBlockQuery() string {
	return `SELECT * FROM reglament WHERE block=?`
}

func (r *Reglament) GetUpdateBlockQuery() string {
	return `UPDATE reglament SET last_updated=:last_updated, in_reglament=:in_reglament WHERE block=:block`
}

func (r *Reglament) GetCreateBlockQuery() string {
	return `INSERT INTO reglament (block, last_updated, in_reglament) VALUES (:block, :last_updated, :in_reglament)`
}
