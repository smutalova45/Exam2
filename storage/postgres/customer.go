package postgres

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"main.go/api/models"
	"main.go/storage"
)

type customerRepo struct {
	db *sql.DB
}

func NewCustomerRepo(db *sql.DB) storage.ICustomerRepo {
	return customerRepo{
		db,
	}
}

func (c customerRepo) Create(customer models.CreateCustomer) (string, error) {
	uid := uuid.New()
	
	if _, err := c.db.Exec(`insert into customers values ($1,$2,$3,$4)`, uid, customer.FullName, customer.Phone, customer.Email); err != nil {
		fmt.Println("error while inserting customer data", err.Error())
		return "", err
	}
	return uid.String(), nil
}

func (c customerRepo) Get(id models.PrimaryKey) (models.Customer, error) {
	customer := models.Customer{}
	query := `
	select id ,full_name, phone,email, created_at from customers where id =$1`
	if err := c.db.QueryRow(query, id.ID).Scan(
		&customer.ID,
		&customer.FullName,
		&customer.Phone,
		&customer.Email,
		&customer.CreatedAt,
	); err != nil {
		fmt.Println("error", err.Error())
		return models.Customer{}, err
	}
	return customer, nil
}

func (c customerRepo) GetList(req models.GetListRequest) (models.CustomersResponse, error) {

	var (
		customers         = []models.Customer{}
		count             = 0
		countQuery, query string
		page              = req.Page
		offset            = (page - 1) * req.Limit
		search            = req.Search
	)
	countQuery = `
	SELECT count(1) from customers
	`
	if search != "" {
		countQuery += fmt.Sprintf(` and (full_name ilike '%%%s%%' or email ilike '%%%s%%')`, search, search)
	}
	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println(err.Error())
		return models.CustomersResponse{}, err
	}
	query = `
	select id, full_name ,phone,email ,created_at from customers`
	if search != "" {
		query += fmt.Sprintf(` and (full_name ilike '%%%s%%' or email ilike '%%%s%%')`, search, search)
	}
	query += ` LIMIT $1 OFFSET $2`
	fmt.Println("here")
	rows, err := c.db.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error is here",err.Error())
		return models.CustomersResponse{}, err
	}
	for rows.Next() {
		customer := models.Customer{}
		if err = rows.Scan(
			&customer.ID,
			&customer.FullName,
			&customer.Phone,
			&customer.Email,
			&customer.CreatedAt,
		); err != nil {
			fmt.Println(err.Error())
			return models.CustomersResponse{}, err
		}
		customers = append(customers, customer)
	}
	return models.CustomersResponse{
		Customers: customers,
		Count:     count,
	}, nil
}

func (c customerRepo) Update(customer models.UpdateCustomer) (string, error) {
	query := `update customers set full_name=$1, phone=$2 where id=$4`
	if _, err := c.db.Exec(query, customer.FullName, customer.Phone, customer.ID); err != nil {
		fmt.Println("error while updating customers", err.Error())
		return "", err
	}
	return customer.ID, nil
}

func (c customerRepo) Delete(id models.PrimaryKey) error {
	
	query2 := `delete from trip_customers where customer_id=$1`
	if _, err := c.db.Exec(query2, id.ID); err != nil {
		fmt.Println(err.Error())
		return err
	}

	query1 := `
	 delete from customers where id=$1`

	if _, err := c.db.Exec(query1, id.ID); err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
