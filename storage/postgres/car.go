package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"main.go/api/models"
	"main.go/storage"
)

type carRepo struct {
	db *sql.DB
}

func NewCarRepo(db *sql.DB) storage.ICarRepo {
	return carRepo{
		db,
	}
}

func (c carRepo) Create(car models.CreateCar) (string, error) {
	uid := uuid.New()
	if _, err := c.db.Exec(`insert into cars values($1,$2,$3,$4,$5,$6)`, uid, car.Model, car.Brand, car.Number,car.Status,car.DriverID); err != nil {
		fmt.Println("error ", err.Error())
		return "", err
	}
	return uid.String(), nil

}

func (c carRepo) Get(id models.PrimaryKey) (models.Car, error) {
	car := models.Car{}
	query := `
	select id, model,brand , number , driver_id, created_at where id=$1`
	if err := c.db.QueryRow(query, id.ID).Scan(&car.ID, &car.Model, &car.Brand, &car.Number,&car.DriverID); err != nil {
		fmt.Println(err.Error())
		return models.Car{}, err
	}
	return car, nil
}

func (c carRepo) GetList(req models.GetListRequest) (models.CarsResponse, error) {
	var (
		cars              = []models.Car{}
		count             = 0
		countQuery, query string
		page              = req.Page
		offset            = (page - 1) * req.Limit
		search            = req.Search
	)
	countQuery = `
	select count(1) from cars
	`
	if search != "" {
		countQuery += fmt.Sprintf(` and (number ilike '%%%s%%' )`, search)
	}
	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println(err.Error())
		return models.CarsResponse{}, err
	}
	query = `
	SELECT c.id, c.model, c.brand, c.number, c.driver_id, c.created_at, 
    d.id as driver_id, d.full_name, d.phone, d.from_city_id, d.to_city_id
FROM cars as c
LEFT JOIN drivers as d ON d.id = c.driver_id
GROUP BY c.id, c.model, c.brand, c.number, c.driver_id, c.created_at, d.id, d.full_name, d.phone, d.from_city_id, d.to_city_id

	`
	if search != "" {
		query += fmt.Sprintf(` and (number ilike '%%%s%%' )`, search)
	}
	query += ` LIMIT $1 OFFSET $2`
	rows, err := c.db.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.CarsResponse{}, err
	}
	for rows.Next() {
		car := models.Car{}
		if err = rows.Scan(
			&car.ID,
			&car.Model,
			&car.Brand,
			&car.Number,
			&car.DriverID,
			&car.DriverData.ID,
			&car.DriverData.FullName,
			&car.DriverData.Phone,
		); err != nil {
			fmt.Println(err.Error())
			return models.CarsResponse{}, err
		}
		cars = append(cars, car)
	}
	return models.CarsResponse{
		Cars:  cars,
		Count: count,
	}, nil
}

func (c carRepo) Update(car models.Car) (string, error) {
	query := `update cars set model=$1 where id=$2`
	if _, err := c.db.Exec(query, car.Model, car.ID); err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return car.ID, nil
}

func (c carRepo) Delete(id models.PrimaryKey) error {
	query := `delete from cars where driver_id=$1`
	if _, err := c.db.Exec(query, id.ID); err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (c carRepo) UpdateCarRoute(updateCarRoute models.UpdateCarRoute) error {

	stmt, err := c.db.Prepare("UPDATE cars SET departure_time=$1, from_city_id=$2, to_city_id=$3 WHERE car_id=$4")
	if err != nil {
		return err
	}
	defer stmt.Close()
	updateCarRoute.DepartureTime = time.Now()
	_, err = stmt.Exec(updateCarRoute.DepartureTime, updateCarRoute.FromCityID, updateCarRoute.ToCityID, updateCarRoute.CarID)
	if err != nil {
		return err
	}

	return nil
}
func (c carRepo) UpdateCarStatus(updateCarStatus models.UpdateCarStatus) error {
	stmt, err := c.db.Prepare("UPDATE cars SET status=$1 WHERE id=$2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(updateCarStatus.Status, updateCarStatus.ID)
	if err != nil {
		return err
	}
	return nil

}
