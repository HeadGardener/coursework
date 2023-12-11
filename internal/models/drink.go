package models

type Drink struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	Type   string `db:"type"`
	Bottle int    `db:"bottle"`
	Cost   int    `db:"cost"`
	Soft   bool   `db:"is_soft"`
}
