package reader

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"log"
	"project-survey-proceeder/internal/configuration"
	"time"
)

type SqlReader struct {
	config *configuration.DbCacheConfiguration
	db     *sql.DB
}

func NewSqlReader(config *configuration.DbCacheConfiguration) *SqlReader {
	return &SqlReader{config: config}
}

func (s *SqlReader) Connect() error {
	var err error
	s.db, err = sql.Open("sqlserver", s.config.ConnectionString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
		return err
	}

	ctx := context.Background()
	for i := 0; i < s.config.ConnectionRetryCount; i++ {
		err = s.db.PingContext(ctx)
		if err != nil {
			log.Println("Error when connecting to DB:", err.Error(), "Try", i+1, "of", s.config.ConnectionRetryCount)
			time.Sleep(time.Duration(s.config.ConnectionRetrySleepTime) * time.Second)
		}
	}
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func (s *SqlReader) GetStoredProcedureResult(storedProcedure string) (*sql.Rows, error) {
	return s.db.Query(fmt.Sprintf("exec %v", storedProcedure))
}

func (s *SqlReader) CloseConnection() {
	if s.db != nil {
		s.db.Close()
	}
}
