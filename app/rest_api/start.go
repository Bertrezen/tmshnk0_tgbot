package rest_api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func (s *Server) Start() {
	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/week", WeekF(s.DB))
	e.GET("/month", MonthF(s.DB))
	e.GET("/last", LastF(s.DB, 5))

	if err := e.Start("localhost:8000"); err != nil {
		log.Fatal(err)
	}
}
