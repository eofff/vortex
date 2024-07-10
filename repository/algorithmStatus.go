package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type AlgorithmStatusRepository struct {
	db *sql.DB
}

type AlgorithmStatus struct {
	Id       int64 `json:"id"`
	ClientId int64 `json:"client_id"`
	VWAP     bool  `json:"VWAP"`
	TWAP     bool  `json:"TWAP"`
	HFT      bool  `json:"HFT"`
}

func (a *AlgorithmStatusRepository) Init(db *sql.DB) {
	a.db = db
}

func (a *AlgorithmStatusRepository) GetAll() ([](*AlgorithmStatus), error) {
	result := make([](*AlgorithmStatus), 0)

	rows, err := a.db.Query(`SELECT "id", "client_id", "VWAP", "TWAP", "HFT" FROM AlgorithmStatus`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		algorithmStatus := new(AlgorithmStatus)

		err = rows.Scan(
			&algorithmStatus.Id,
			&algorithmStatus.ClientId,
			&algorithmStatus.VWAP,
			&algorithmStatus.TWAP,
			&algorithmStatus.HFT,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, algorithmStatus)
	}

	return result, nil
}
