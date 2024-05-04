package reader

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"log"
)

type SqlReader struct {
	connectionString string
	db               *sql.DB
}

func NewSqlReader(connectionString string) *SqlReader {
	return &SqlReader{connectionString: connectionString}
}

func (s *SqlReader) Connect() error {
	var err error
	s.db, err = sql.Open("sqlserver", s.connectionString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
		return err
	}

	ctx := context.Background()
	err = s.db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func (s *SqlReader) GetStoredProcedureResult(storedProcedure string) (*sql.Rows, error) {
	return s.db.Query(fmt.Sprintf("exec %v", storedProcedure))
}
