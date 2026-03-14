package model

type User struct {
	ID   int64  `db:"user_id"`
	Name string `db:""`
}
