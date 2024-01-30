package db

import "lynx-backend/models"

type DB interface {
	Disconnect() error

	CreateMetric(m *models.Metric) error
	Metrics() ([]models.Metric, error)
}
