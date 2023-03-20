package rest_api

import (
	"context"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

type Info struct {
	Time_add pgtype.Timestamp `json:"time_add"`
	Price    int              `json:"price"`
	Category string           `json:"category"`
}

func Last(db *pgx.Conn, N int) ([]Info, error) {
	rows, err := db.Query(context.Background(), "SELECT time_add, category, price FROM info ORDER BY id DESC LIMIT $1", N)
	if err != nil {
		log.Panic(err, " - trouble with func last")
	}

	var body []Info

	for rows.Next() {
		s := Info{}
		err = rows.Scan(&s.Time_add, &s.Category, &s.Price)
		if err != nil {
			return nil, err
		}
		body = append(body, s)
	}
	return body, nil
}

func Week(db *pgx.Conn) ([]Info, error) {
	rows, err := db.Query(context.Background(), "SELECT category, sum(price) FROM info WHERE (time_add = now() or time_add >= NOW()::DATE - 7) group by category ORDER BY sum(price)")
	if err != nil {
		log.Panic(err, " - trouble with func week")
	}

	var body []Info

	for rows.Next() {
		s := Info{}
		err = rows.Scan(&s.Category, &s.Price)
		if err != nil {
			return nil, err
		}
		body = append(body, s)
	}
	return body, nil
}

func Month(db *pgx.Conn) ([]Info, error) {
	rows, err := db.Query(context.Background(), "SELECT category, sum(price) FROM info WHERE time_add between now() - interval '1 months' and now()  group by category ORDER BY sum(price)")
	if err != nil {
		log.Panic(err, " - trouble with func month")
	}

	var body []Info

	for rows.Next() {
		s := Info{}
		err = rows.Scan(&s.Category, &s.Price)
		if err != nil {
			return nil, err
		}
		body = append(body, s)
	}
	return body, nil
}

func WeekF(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {

		WeekL, err := Week(db)
		fmt.Println("", WeekL)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return errors.Wrap(c.JSON(http.StatusOK, WeekL), "error with week/json")
	}
}

func MonthF(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {

		MonthL, err := Month(db)
		fmt.Println("", MonthL)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return errors.Wrap(c.JSON(http.StatusOK, MonthL), "error with month/json")
	}
}

func LastF(db *pgx.Conn, N int) echo.HandlerFunc {
	return func(c echo.Context) error {
		LastL, err := Last(db, N)
		fmt.Println("", LastL)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return errors.Wrap(c.JSON(http.StatusOK, LastL), "error with last/json")
	}
}
