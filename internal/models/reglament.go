package models

import "time"

//Reglament ...
type Reglament struct {
	ID                 int       `db:"id"`
	Code               string    `db:"code"`
	Block              string    `db:"block"`
	LastUpdated        time.Time `db:"last_updated"`
	Active             bool      `db:"active"`
	ServiceBlock       bool      `db:"service_block"`
	InReglament        bool      `db:"in_reglament"`
	ReglamentDayTime   int       `db:"reglament_day_time"`
	ReglamentNightTime int       `db:"reglament_night_time"`
	DayHour            int       `db:"day_hour"`
	NightHour          int       `db:"night_hour"`
}

func (r *Reglament) GetBlockQuery() string {
	return `SELECT * FROM reglament WHERE block=?`
}

func (r *Reglament) GetBlockByCodeQuery() string {
	return `SELECT * FROM reglament WHERE code=?`
}
func (r *Reglament) GetBlocksQuery() string {
	return `SELECT * FROM reglament WHERE active and not service_block`
}

func (r *Reglament) GetUpdateBlockQuery() string {
	return `UPDATE reglament SET last_updated=:last_updated, in_reglament=:in_reglament WHERE id=:id`
}

func (r *Reglament) GetCreateBlockQuery() string {
	return `
	INSERT INTO reglament 
	(code, block, last_updated, active, in_reglament, reglament_day_time, reglament_night_time, 'day_hour', night_hour, service_block) 
	VALUES (:code, :block, :last_updated, :active, :in_reglament, :reglament_day_time, :reglament_night_time, :day_hour, :night_hour, :service_block)`
}
