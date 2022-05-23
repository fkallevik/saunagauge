package sauna

import "github.com/jmoiron/sqlx"

type Sauna struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Sauna {
	return &Sauna{
		db: db,
	}
}

func (sauna *Sauna) GetTemperature() float32 {
	return 85.5
}

func (sauna *Sauna) GetHumidity() float32 {
	return 14.4
}
