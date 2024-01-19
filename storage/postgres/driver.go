package postgres

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"main.go/api/models"
	"main.go/storage"
)

type driverRepo struct {
	DB *sql.DB
}

func NewDriverRepo(db *sql.DB) storage.IDriverRepo {
	return driverRepo{
		DB: db,
	}
}

func (d driverRepo) Create(driver models.CreateDriver) (string, error) {
	uid := uuid.New()
	if _, err := d.DB.Exec(`insert into drivers values ($1,$2,$3,$4,$5)`,
		uid,
		driver.FullName,
		driver.Phone,
		driver.FromCityID,
		driver.ToCityID,
	); err != nil {

		fmt.Println("error", err.Error())
		return "", err
	}
	return uid.String(), nil
}

func (d driverRepo) Get(id models.PrimaryKey) (models.Driver, error) {
	driver := models.Driver{}
	query := `
	select id , full_name , phone ,from_city_id,to_city_id,created_at from drivers where id=$1`
	if err := d.DB.QueryRow(query, id.ID).Scan(
		&driver.ID,
		&driver.FullName,
		&driver.Phone,
		&driver.FromCityID,
		&driver.ToCityID,
		&driver.CreatedAt,
	); err != nil {
		fmt.Println(err.Error())
		return models.Driver{}, err
	}
	return driver, nil
}

func (d driverRepo) GetList(req models.GetListRequest) (models.DriversResponse, error) {
	var (
		drivers           = []models.Driver{}
		count             = 0
		countQuery, query string
		page              = req.Page
		offset            = (page - 1) * req.Limit
		search            = req.Search
	)
	countQuery = `select count(1) from  drivers`
	if search != "" {
		countQuery += fmt.Sprintf(` and (full_name ilike '%%%s%%')`, search)
	}
	if err := d.DB.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println(err.Error())
		return models.DriversResponse{}, err
	}
	query = `SELECT d.id, d.full_name, d.phone, d.from_city_id, d.to_city_id, d.created_at,
    c.id as city_id, c.name as city_name, c.created_at as city_created_at
FROM drivers as d
LEFT JOIN cities as c ON d.from_city_id = c.id
GROUP BY d.id, d.full_name, d.phone, d.from_city_id, d.to_city_id, d.created_at, c.id, c.name, c.created_at

`
	if search != "" {
		query += fmt.Sprintf(` and (full_name ilike '%%%s%%')`, search)
	}
	query += ` LIMIT $1 OFFSET $2`
	rows, err := d.DB.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println(err.Error())
		return models.DriversResponse{}, err
	}

	for rows.Next() {
		driver := models.Driver{}
		if err = rows.Scan(
			&driver.ID,
			&driver.FullName,
			&driver.Phone,
			&driver.FromCityID,
			&driver.FromCityData.Name,
			&driver.FromCityData.CreatedAt,
			&driver.ToCityID,
			&driver.ToCityData.ID,
			&driver.ToCityData.Name,
			&driver.ToCityData.CreatedAt,
			&driver.CreatedAt,
		); err != nil {
			fmt.Println(err.Error())
			return models.DriversResponse{}, err
		}
		drivers = append(drivers, driver)
	}

	return models.DriversResponse{
		Drivers: drivers,
		Count:   count,
	}, nil
}

func (d driverRepo) Update(driver models.UpdateDriver) (string, error) {
	query := `update drivers set full_name=$1, phone=$2 where id=$3`
	if _, err := d.DB.Exec(query, driver.FullName, driver.Phone, driver.ID); err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return driver.ID, nil
}

func (d driverRepo) Delete(id models.PrimaryKey) error {
	query2 := `delete from cars where driver_id=$1`
	if _, err := d.DB.Exec(query2, id.ID); err != nil {
		fmt.Println(err.Error())
		return err

	}
	query := `delete from drivers where id=$1`
	if _, err := d.DB.Exec(query, id.ID); err != nil {
		fmt.Println(err.Error())
		return err

	}
	return nil
}
