package db

import (
	"context"
	"errors"
	"fmt"
	"slurm/go-on-practice-2/http_06/internals/app/models"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type CarsStorage struct {
	databasePool *pgxpool.Pool
}

type userCar struct {
	UserId  int64 `db:"userid"`
	Name    string
	Rank    string
	CarId   int64 `db:"carid"`
	Colour  string
	Brand   string
	License string
}

func convertJoinedQueryToCar(input userCar) models.Car {
	return models.Car{
		Id:      input.CarId,
		Colour:  input.Colour,
		Brand:   input.Brand,
		License: input.License,
		Owner: models.User{
			Id:   input.UserId,
			Name: input.Name,
			Rank: input.Rank,
		},
	}
}

func NewCarsStorage(pool *pgxpool.Pool) *CarsStorage {
	storage := new(CarsStorage)
	storage.databasePool = pool
	return storage
}

func (storage *CarsStorage) CreateCar(car models.Car) error {

	ctx := context.Background()
	tx, err := storage.databasePool.Begin(ctx)
	defer func() {
		err = tx.Rollback(context.Background())
		if err != nil {
			logrus.Errorln(err)
		}
	}()

	query := "SELECT id FROM users WHERE id = $1"

	id := -1

	err = pgxscan.Get(ctx, tx, &id, query, car.Owner.Id)
	if err != nil {
		logrus.Errorln(err)
		err = tx.Rollback(context.Background())
		if err != nil {
			logrus.Errorln(err)
		}
		return err
	}

	if id == -1 {
		return errors.New("user not found")
	}

	insertQuery := "INSERT INTO cars(user_id, colour, brand, license) VALUES ($1, $2, $3, $4)"

	_, err = tx.Exec(context.Background(), insertQuery, car.Owner.Id, car.Colour, car.Brand, car.License)

	if err != nil {
		logrus.Errorln(err)
		err = tx.Rollback(context.Background())
		if err != nil {
			logrus.Errorln(err)
		}
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	return nil
}

func (storage *CarsStorage) GetCarById(id int64) models.Car {
	query := `
		SELECT 
		users.id AS userid,	users.name,	users.rank,
		c.id AS carid, c.colour, c.brand, c.license
		FROM users JOIN cars c on users.id = c.user_id
		WHERE c.id = $1`

	var result userCar

	err := pgxscan.Get(context.Background(), storage.databasePool, &result, query, id)
	if err != nil {
		logrus.Errorln(err)
	}

	return convertJoinedQueryToCar(result)
}

func (storage *CarsStorage) GetCarsList(userIdFilter int64, brandFilter, colourFilter, licenseFilter string) []models.Car {
	query := `
		SELECT 
		users.id AS userid,	users.name,	users.rank,
		c.id AS carid, c.colour, c.brand, c.license
		FROM users JOIN cars c on users.id = c.user_id
		WHERE 1=1`

	placeholderNum := 1
	args := make([]any, 0)

	if userIdFilter != 0 {
		query += fmt.Sprintf(" AND users.id = $%d", placeholderNum)
		args = append(args, userIdFilter)
		placeholderNum++
	}
	if brandFilter != "" {
		query += fmt.Sprintf(" AND c.brand ILIKE $%d", placeholderNum)
		args = append(args, fmt.Sprintf("%%%s%%", brandFilter))
		placeholderNum++
	}
	if colourFilter != "" {
		query += fmt.Sprintf(" AND c.colour ILIKE $%d", placeholderNum)
		args = append(args, fmt.Sprintf("%%%s%%", colourFilter))
		placeholderNum++
	}
	if licenseFilter != "" {
		query += fmt.Sprintf(" AND c.license ILIKE $%d", placeholderNum)
		args = append(args, fmt.Sprintf("%%%s%%", licenseFilter))
		// placeholderNum++
	}

	var dbResult []userCar

	err := pgxscan.Select(context.Background(), storage.databasePool, &dbResult, query, args...)
	if err != nil {
		logrus.Errorln(err)
	}

	result := make([]models.Car, len(dbResult))

	for idx, dbEntity := range dbResult {
		result[idx] = convertJoinedQueryToCar(dbEntity)
	}

	return result
}
