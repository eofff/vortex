package repository

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	_ "github.com/lib/pq"
)

type AlgorithmStatusRepository struct {
	Id       int64 `json:"id"`
	ClientId int64 `json:"client_id"`
	VWAP     bool  `json:"VWAP"`
	TWAP     bool  `json:"TWAP"`
	HFT      bool  `json:"HFT"`
}

var db *sql.DB

func Init(dbp *sql.DB) {
	db = dbp
}

func (a *AlgorithmStatusRepository) GetAll() ([](*AlgorithmStatusRepository), error) {
	result := make([](*AlgorithmStatusRepository), 0)

	rows, err := db.Query(`SELECT "id", "client_id", "VWAP", "TWAP", "HFT" FROM AlgorithmStatus`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		algorithmStatus := new(AlgorithmStatusRepository)

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

func (a *AlgorithmStatusRepository) GetById(id int) (*AlgorithmStatusRepository, error) {
	rows, err := db.Query(
		`SELECT "id", "client_id", "VWAP", "TWAP", "HFT" FROM AlgorithmStatus WHERE indexcolumn = $1`,
		id,
	)

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		algorithmStatus := new(AlgorithmStatusRepository)

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

		return algorithmStatus, nil
	} else {
		return nil, errors.New("Algorithm status with id = " + strconv.Itoa(id) + " not exist")
	}
}

func (a *AlgorithmStatusRepository) Insert(algorithmStatus *AlgorithmStatusRepository) error {
	_, err := db.Query(
		"INSERT INTO AlgorithmStatus (client_id, VWAP, TWAP, HFT) VALUES ($1, $2, $3, $4);",
		algorithmStatus.ClientId,
		algorithmStatus.VWAP,
		algorithmStatus.TWAP,
		algorithmStatus.HFT,
	)

	if err == nil {
		log.Printf("Algorithm status with id = %d inserted into db", algorithmStatus.Id)
	}

	return err
}

func (a *AlgorithmStatusRepository) Update(algorithmStatus *AlgorithmStatusRepository) error {
	_, err := db.Query(
		"UPDATE AlgorithmStatus SET (client_id, VWAP, TWAP, HFT) = ($2, $3, $4, $5) WHERE id = $1;",
		algorithmStatus.Id,
		algorithmStatus.ClientId,
		algorithmStatus.VWAP,
		algorithmStatus.TWAP,
		algorithmStatus.HFT,
	)

	if err == nil {
		log.Printf("Algorithm status with id = %d updated in db", algorithmStatus.Id)
	}

	return err
}

func (a *AlgorithmStatusRepository) Delete(algorithmStatus *AlgorithmStatusRepository) error {
	_, err := db.Query(
		"DELETE FROM AlgorithmStatus WHERE id = $1;",
		algorithmStatus.Id,
	)

	if err == nil {
		log.Printf("Algorithm status with id = %d deleted from db", algorithmStatus.Id)
	}

	return err
}
