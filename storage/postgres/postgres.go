package postgres

import (
	"database/sql"
	"fmt"

	"main.go/config"
	"main.go/storage"
)

type Store struct {
	db *sql.DB
}

func New(cfg config.Config) (storage.IStorage, error) {
	url := fmt.Sprintf(`host = %s port = %s user = %s password = %s database = %s sslmode=disable`,
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return Store{}, err
	}

	return Store{
		db: db,
	}, nil
}

func (s Store) CloseDB() {
	s.db.Close()
}

func (s Store) City() storage.ICityRepo {
	newcityrepo:=NewCityRepo(s.db)
	return newcityrepo
}

func (s Store) Customer() storage.ICustomerRepo {
	newcustomerrepo:=NewCustomerRepo(s.db)
	return newcustomerrepo
}


func (s Store) Driver() storage.IDriverRepo {
	newdriver:=NewDriverRepo(s.db)
	return newdriver
}

func (s Store) Car() storage.ICarRepo {
	newcarrepo:=NewCarRepo(s.db)
	return newcarrepo
}

func (s Store) Trip() storage.ITripRepo {
	newtriprepo:=NewTripRepo(s.db)
	return newtriprepo
}
func (s Store) TripCustomer() storage.ITripCustomerRepo {
	newtripcustomer:=NewTripCustomerRepo(s.db)
	return newtripcustomer
}
