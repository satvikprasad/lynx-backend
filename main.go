package main

import (
	"fmt"
	"log"
	"lynx-backend/db"
	"lynx-backend/models"
	"lynx-backend/server"
	"os"

	"github.com/joho/godotenv"
)

type MetricRequest struct {
	Time float64 `json:"time"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Didn't find a .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	db, err := db.CreateMongoDB()
	if err != nil {
		panic(err)
	}

	defer db.Disconnect()

	s := server.CreateServer(db, port)

	s.Get("/", func(c *server.Context) error {
		c.WriteHTML(200, "home/index.tmpl", nil)
		return nil
	})

	s.Get("/download", func(c *server.Context) error {
		c.WriteHTML(200, "home/index.tmpl", nil)
		return nil
	})

	s.Get("/home-page", func(c *server.Context) error {
		c.WriteHTML(200, "home/home.tmpl", nil)
		return nil
	})

	s.Get("/download-page", func(c *server.Context) error {
		c.WriteHTML(200, "home/download.tmpl", nil)
		return nil
	})

	s.Get("/api/v1/health", func(c *server.Context) error {
		c.WriteJSON(200, map[string]string{"status": "ok"})
		return nil
	})

	s.Get("/api/v1/metrics", func(c *server.Context) error {
		metrics, err := c.Database.Metrics()
		if err != nil {
			return err
		}

		c.WriteJSON(200, metrics)
		return nil
	})

	s.Get("/api/v1/metrics/total-time", func(c *server.Context) error {
		metrics, err := c.Database.Metrics()
		if err != nil {
			return err
		}

		var totalTime float64 = 0
		for _, metric := range metrics {
			totalTime += metric.Time
		}

		c.WriteJSON(200, map[string]float64{"time": totalTime})
		return nil
	})

	s.Post("/api/v1/add-metric", func(c *server.Context) error {
		r := MetricRequest{}

		if err := c.BindJSON(&r); err != nil {
			return err
		}

		if err := c.Database.CreateMetric(&models.Metric{Time: r.Time}); err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})

	s.Listen()
}
