package rest_api

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

type Server struct {
	DB *pgx.Conn
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	s.DB = conn
}

func (s *Server) Close() {
	if err := s.DB.Close(context.Background()); err != nil {
		log.Fatal(err)
	}
}
