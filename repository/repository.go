// repository package
package repository

import "database/sql"

type IAlgorithmStatusRepository interface {
	Init(db *sql.DB)
	GetAll() ([](*AlgorithmStatus), error)
}
