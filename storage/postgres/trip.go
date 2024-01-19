package postgres

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"main.go/api/models"
	"main.go/storage"
)

type tripRepo struct {
	db                *sql.DB
	tripNumberCounter int
}

func NewTripRepo(db *sql.DB) storage.ITripRepo {
	return &tripRepo{
		db: db,
	}
}
func getMaxTripNumber() (int, error) {
	maxTripNumber := 0
	return maxTripNumber, nil
}
func (c *tripRepo) Create(req models.CreateTrip) (string, error) {
	uid := uuid.New()
	if c.tripNumberCounter == 0 {
		maxTripNumber, err := getMaxTripNumber()
		if err != nil {
			fmt.Println("error while maxtripnumber", err.Error())
			return "", err
		}

		c.tripNumberCounter = maxTripNumber + 1
	}

	tripNumberID := fmt.Sprintf("T-%d", c.tripNumberCounter)

	if _, err := c.db.Exec(`insert into trips values($1,$2,$3,$4,$5,$6)`,
		uid,
		tripNumberID,
		req.FromCityID,
		req.ToCityID,
		req.DriverID,
		req.Price,
	); err != nil {
		fmt.Println("error is here", err.Error())
		return "", err
	}
	c.tripNumberCounter++
	return uid.String(), nil
}

func (c *tripRepo) Get(id models.PrimaryKey) (models.Trip, error) {
	trip := models.Trip{}
	query := `select id, trip_number_id, from_city_id,to_city_id,driver_id,price,created_at from trips where id=$1`
	if err := c.db.QueryRow(query, id.ID).Scan(&trip.ID, &trip.TripNumberID, &trip.FromCityID, &trip.ToCityID, &trip.DriverID, &trip.Price, &trip.CreatedAt); err != nil {
		return models.Trip{}, err
	}
	return trip, nil
}

func (c *tripRepo) GetList(req models.GetListRequest) (models.TripsResponse, error) {
	var (
		trips             = []models.Trip{}
		count             = 0
		countQuery, query string
		page              = req.Page
		offset            = (page - 1) * req.Limit
		search            string
	)
	countQuery = `select count(1) from trips`
	if search != "" {
		countQuery += fmt.Sprintf(` and ( price ilike '%%%s%%')`, search)
	}
	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		return models.TripsResponse{}, err
	}

	query = `SELECT 
    t.id AS trip_id,
    t.trip_number_id,
    t.from_city_id,
    from_city.id,
    from_city.name AS from_city_name,
    from_city.created_at AS from_city_created_at,
    t.to_city_id,
    to_city.id AS to_city_id,
    to_city.name AS to_city_name,
    to_city.created_at AS to_city_created_at,
    t.driver_id,
    driver.id AS driver_id,
    driver.full_name AS driver_full_name,
    driver.phone AS driver_phone,
    driver.from_city_id AS driver_from_city_id,
    driver_from_city.id AS driver_from_city_id,
    driver_from_city.name AS driver_from_city_name,
    driver_from_city.created_at AS driver_from_city_created_at,
    driver.to_city_id AS driver_to_city_id,
    driver_to_city.id AS driver_to_city_id,
    driver_to_city.name AS driver_to_city_name,
    driver_to_city.created_at AS driver_to_city_created_at,
    t.price,
    t.created_at AS trip_created_at
FROM 
    trips t
LEFT JOIN
    cities from_city ON t.from_city_id = from_city.id
LEFT JOIN
    cities to_city ON t.to_city_id = to_city.id
LEFT JOIN
    drivers driver ON t.driver_id = driver.id
LEFT JOIN
    cities driver_from_city ON driver.from_city_id = driver_from_city.id
LEFT JOIN
    cities driver_to_city ON driver.to_city_id = driver_to_city.id  `
	if search != "" {
		query += fmt.Sprintf(` and (full_name ilike '%%%s%%')`, search)
	}
	query += ` LIMIT $1 OFFSET $2 `
	rows, err := c.db.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error:->", err.Error())
		return models.TripsResponse{}, err
	}

	for rows.Next() {

		trip := models.Trip{}
		if err = rows.Scan(
			&trip.ID,
			&trip.TripNumberID,
			&trip.FromCityID,
			&trip.FromCityData.ID,
			&trip.FromCityData.Name,
			&trip.FromCityData.CreatedAt,
			&trip.ToCityID,
			&trip.ToCityData.ID,
			&trip.ToCityData.Name,
			&trip.ToCityData.CreatedAt,
			&trip.DriverID,
			&trip.DriverData.ID,
			&trip.DriverData.FullName,
			&trip.DriverData.Phone,
			&trip.DriverData.FromCityID,
			&trip.DriverData.FromCityData.ID, 
			&trip.DriverData.FromCityData.Name,
			&trip.DriverData.FromCityData.CreatedAt,
			&trip.DriverData.ToCityID,
			&trip.DriverData.ToCityData.ID,
			&trip.DriverData.ToCityData.Name,
			&trip.DriverData.ToCityData.CreatedAt,
			&trip.Price,
			&trip.CreatedAt,
		); err != nil {
			return models.TripsResponse{}, err
		}
		trips = append(trips, trip)
	}

	return models.TripsResponse{
		Trips: trips,
		Count: count,
	}, nil
}

func (c *tripRepo) Update(req models.Trip) (string, error) {
	query := `update trips set price=$1 where id=$2`
	if _, err := c.db.Exec(query, req.Price, req.ID); err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return req.ID, nil
}

func (c *tripRepo) Delete(id models.PrimaryKey) error {
	query2 := `delete from trip_customers where id=$1`
	if _, err := c.db.Exec(query2, id.ID); err != nil {
		return err
	}
	query := `delete from trips where id=$1`
	if _, err := c.db.Exec(query, id.ID); err != nil {
		return err
	}
	return nil
}
