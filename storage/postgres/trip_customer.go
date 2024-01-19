package postgres

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"main.go/api/models"
	"main.go/storage"
)

type tripCustomerRepo struct {
	db *sql.DB
}

func NewTripCustomerRepo(db *sql.DB) storage.ITripCustomerRepo {
	return &tripCustomerRepo{
		db: db,
	}
}

func (c *tripCustomerRepo) Create(req models.CreateTripCustomer) (string, error) {

	uid := uuid.New()
	if _, err := c.db.Exec(`insert into trip_customers values ($1,$2,$3)`,
		uid,
		req.TripID,
		req.CustomerID,
	); err != nil {
		fmt.Println("error while inserting", err.Error())
		return "", err
	}

	return uid.String(), nil
}

func (c *tripCustomerRepo) Get(id models.PrimaryKey) (models.TripCustomer, error) {
	tripcustomer := models.TripCustomer{}
	query := `select id, trip_id,customer_id, created_at from trip_customers where id=$1`
	if err := c.db.QueryRow(query, id.ID).Scan(
		&tripcustomer.ID,
		&tripcustomer.TripID,
		&tripcustomer.CustomerID,
		&tripcustomer.CreatedAt,
	); err != nil {
		fmt.Println(err.Error())
		return models.TripCustomer{}, err
	}
	return tripcustomer, nil
}

func (c *tripCustomerRepo) GetList(req models.GetListRequest) (models.TripCustomersResponse, error) {
	var (
		tripcustomers     = []models.TripCustomer{}
		count             = 0
		countQuery, query string
		page              = req.Page
		offset            = (page - 1) * req.Limit
	)

	countQuery = `select count(1) from trip_customers`
	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while counting", err.Error())
		return models.TripCustomersResponse{}, err
	}
	query = `SELECT
    t.id,
    t.trip_id,
    t.customer_id,
    t.created_at,
    c.id AS customer_id,
    c.fullname,
    c.phone,
    c.email
FROM
    trip_customers AS t
LEFT JOIN
    customers AS c ON t.customer_id = c.id
GROUP BY
    t.id, t.trip_id, t.customer_id, t.created_at, c.id, c.full_name, c.phone`

	query += ` LIMIT $1 OFFSET $2`
	rows, err := c.db.Query(query, req.Limit, offset)
	if err != nil {
		return models.TripCustomersResponse{}, err
	}
	for rows.Next() {
		tc := models.TripCustomer{}
		if err = rows.Scan(
			&tc.ID,
			&tc.TripID,
			&tc.CustomerID,
			&tc.CreatedAt,
			&tc.CustomerData.ID,
			&tc.CustomerData.FullName,
			&tc.CustomerData.Phone,
			&tc.CustomerData.Email,
		); err != nil {
			return models.TripCustomersResponse{}, err
		}
		tripcustomers = append(tripcustomers, tc)
	}

	return models.TripCustomersResponse{
		TripCustomers: tripcustomers,
		Count:         count,
	}, nil
}

func (c *tripCustomerRepo) Update(req models.TripCustomer) (string, error) {
	query := `update trip_customers set customer_id=$1 where id=$2`
	if _, err := c.db.Exec(query, req.CustomerID, req.ID); err != nil {
		return "", err
	}
	return req.ID, nil
}

func (c *tripCustomerRepo) Delete(id models.PrimaryKey) error {
	query := `delete from trip_customers where id=$1`
	if _, err := c.db.Exec(query, id.ID); err != nil {
		fmt.Println(err.Error())
		return err

	}
	return nil
}
