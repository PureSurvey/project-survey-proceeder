package contracts

import "database/sql"

type IReader interface {
	Connect() error
	GetStoredProcedureResult(storedProcedure string) (*sql.Rows, error)
	CloseConnection()
}
