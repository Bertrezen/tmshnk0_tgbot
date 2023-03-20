package bd

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"log"
)

type Info struct {
	Time_add pgtype.Timestamp `json:"time_add"`
	Price    int              `json:"price"`
	Category string           `json:"category"`
}

func WriteBD(db *pgx.Conn, ChatID, Nick, category string, price int) {

	if category != "" {
		s := `INSERT INTO info (chat_id, nickname, category, price, time_add)
		  VALUES ($1, $2, $3, $4, now())`
		_, err := db.Exec(context.Background(), s, ChatID, Nick, category, price)
		if err != nil {
			log.Panic(err, " - trouble with writeBD")
		}
	} else {
		s := `INSERT INTO info (chat_id, nickname, price, time_add)
		  VALUES ($1, $2, $3, now())`
		_, err := db.Exec(context.Background(), s, ChatID, Nick, price)
		if err != nil {
			log.Panic(err, " - trouble with writeBD(not category)")
		}
	}
}

func Last(db *pgx.Conn, ChatID string, N int) ([]Info, error) {
	rows, err := db.Query(context.Background(), "SELECT time_add, category, price FROM info WHERE chat_id = $1 order by time_add desc LIMIT $2", ChatID, N)
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

func Week(db *pgx.Conn, ChatID string) ([]Info, error) {
	rows, err := db.Query(context.Background(), "SELECT category, sum(price) FROM info WHERE chat_id = $1 AND (time_add = now() or time_add >= NOW()::DATE - 7) group by category ORDER BY sum(price) desc", ChatID)
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

func Month(db *pgx.Conn, ChatID string) ([]Info, error) {
	rows, err := db.Query(context.Background(), "SELECT category, sum(price) FROM info WHERE chat_id = $1  AND (time_add between now() - interval '1 months' and now())  group by category ORDER BY sum(price) desc", ChatID)
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
