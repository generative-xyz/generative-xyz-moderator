package connections

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type postgresCN struct {
	Cnn *gorm.DB
}

func NewPostgres(dsn string) (*postgresCN, error) {
	p := new(postgresCN)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {

		return nil, err
	}
	p.Cnn = db
	return p, nil
}

func (s *postgresCN) Connect() interface{} {
	return s.Cnn
}

func (s *postgresCN) Disconnect() error {
	return nil
}

func (s *postgresCN) GetType() interface{} {
	return *s.Cnn
}
