package models

import (
	"time"

	"github.com/jinzhu/now"
)

type reports string

func (r reports) String() string {
	switch r {
	case pastMonth:
		return "#rep1"
	case currMonth:
		return "#rep2"
	case yesterday:
		return "#rep3"
	case today:
		return "#rep4"
	}
	return ""
}

const (
	pastMonth reports = "#rep1"
	currMonth reports = "#rep2"
	yesterday reports = "#rep3"
	today     reports = "#rep4"
)

//Statistic model
type Statistic struct {
	ID                int       `db:"id"`
	Block             string    `db:"block_name"`
	OutOfReglament    time.Time `db:"out_of_reglament"`
	ReturnInReglament time.Time `db:"return_in_reglament"`
	Downtime          int       `db:"downtime"`
}

//GetReportQuery -return query str for get data by selected reps
func (s *Statistic) GetReportQuery(rep string) (string, time.Time, time.Time) {
	var res string = `SELECT block_name as block, CAST(SUM(downtime)/60 as SIGNED) as downtime, COUNT(downtime) as counts FROM statistic s 
						WHERE out_of_reglament BETWEEN ? AND ?
						GROUP BY block_name`
	param1 := now.BeginningOfDay()
	param2 := now.EndOfDay()
	switch rep {
	case pastMonth.String():
		param1 = now.BeginningOfMonth().AddDate(0, -1, 0)
		param2 = now.New(now.EndOfMonth().AddDate(0, -1, 0)).EndOfMonth()
	case currMonth.String():
		param1 = now.BeginningOfMonth()
		param2 = now.EndOfMonth()
	case yesterday.String():
		param1 = param1.AddDate(0, 0, -1)
		param2 = param2.AddDate(0, 0, -1)
	case today.String():
		// Дата за сегодня выставленна по умолчанию
	default:
		// Дата за сегодня выставленна по умолчанию
	}
	return res, param1, param2
}

// GetReportTitle - return report title
func (s *Statistic) GetReportTitle(rep string) string {
	switch rep {
	case pastMonth.String():
		return "Статистика простоя за прошлый месяц:\n"
	case currMonth.String():
		return "Статистика простоя за текущий месяц:\n"
	case yesterday.String():
		return "Статистика простоя за вчера:\n"
	case today.String():
		return "Статистика простоя за сегодня:\n"
	default:
		return "Статистика простоя за сегодня:\n"
	}
}

//CreateNewStatiscitQuery ...
func (s *Statistic) CreateNewStatiscitQuery() string {
	return `INSERT INTO statistic (block_name, out_of_reglament, return_in_reglament, downtime) 
	VALUES (:block_name, :out_of_reglament, :return_in_reglament, :downtime)`
}
