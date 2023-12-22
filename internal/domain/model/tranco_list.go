package model

import "time"

type TrancoList struct {
	ID        string
	CreatedOn time.Time `db:"created_on"`
}
