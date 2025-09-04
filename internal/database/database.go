package database

import (
	"context"
	"fmt"
	"main/internal/config"
	"main/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseInterface interface {
	Subscribe(sub models.Subscribe) (int, error)
	DeleteSubscribe(id uint) error
	UpdateSubscribe(sub models.Subscribe) error
	GetSubscribe(f models.Filters) ([]models.Subscribe, error)
}

type Database struct {
	Pool *pgxpool.Pool
	Env  config.Env
}

func InitBD(e config.Env) (Database, error) {
	pool, err := pgxpool.New(context.Background(), e.Postgresql)
	if err != nil {
		return Database{}, err
	}
	return Database{Pool: pool}, nil
}

func (db *Database) Subscribe(Sub models.Subscribe) (int, error) {
	var id int
	err := db.Pool.QueryRow(context.Background(), `INSERT INTO subscriptions (
	serviceName,
	price,
	userID,
	startDate,
	endDate) VALUES ($1,$2,$3,$4,$5) RETURNING id;`, Sub.ServiceName, Sub.Price, Sub.UserId, Sub.Date.StartDate, Sub.Date.EndDate).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (db *Database) DeleteSubscribe(id uint) error {
	_, err := db.Pool.Exec(context.Background(), `DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateSubscribe(sub models.Subscribe) error {
	_, err := db.Pool.Exec(context.Background(), `UPDATE subscriptions SET (serviceName, price, userID, startDate, endDate) = ($1,$2,$3,$4,$5) WHERE id = $6`,
		sub.ServiceName,
		sub.Price,
		sub.UserId,
		sub.Date.StartDate,
		sub.Date.EndDate,
		sub.Id)
	if err != nil {
		return err
	}
	return nil
}
func (db *Database) GetSubscribe(f models.Filters) ([]models.Subscribe, error) {
	var subs []models.Subscribe

	query, args := makeQuery(f)

	rows, err := db.Pool.Query(context.Background(), query, args...)

	if err != nil {
		return subs, err
	}
	defer rows.Close()
	for rows.Next() {
		var sub models.Subscribe
		if err := rows.Scan(
			&sub.ServiceName,
			&sub.Price,
			&sub.UserId,
			&sub.Date.StartDate,
			&sub.Date.EndDate); err != nil {
			return subs, err
		}
		if !sub.Date.EndDate.IsZero() {
			sub.EndDate = sub.Date.EndDate.Format("01.2006")
		}
		if !f.Date.EndDate.IsZero() && sub.Date.EndDate.IsZero() {
			continue
		}
		sub.StartDate = sub.Date.StartDate.Format("01.2006")
		subs = append(subs, sub)
	}
	return subs, nil
}

func makeQuery(f models.Filters) (string, []any) {
	query := `SELECT serviceName, price, userID, startDate, endDate FROM subscriptions WHERE 1=1`
	args := []any{}
	i := 1

	if f.ServiceName != "" {
		query += fmt.Sprintf(" AND serviceName = $%d", i)
		args = append(args, f.ServiceName)
		i++
	}

	if f.UserId != uuid.Nil {
		query += fmt.Sprintf(" AND userID = $%d", i)
		args = append(args, f.UserId)
		i++
	}

	if f.StartDate != "" {
		query += fmt.Sprintf(" AND startDate >= $%d", i)
		args = append(args, f.Date.StartDate)
		i++
	}

	if f.EndDate != "" {
		query += fmt.Sprintf(" AND endDate <= $%d", i)
		args = append(args, f.Date.EndDate)
		i++
	}
	return query, args
}
