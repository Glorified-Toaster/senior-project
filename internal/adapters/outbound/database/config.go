package database

import "time"

// Config holds PostgreSQL connection configuration
type DBConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	DBName          string
	SSLMode         string
	MaxConns        int
	MinConns        int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// DefaultConfig returns default configuration
func DefaultConfig() DBConfig {
	return DBConfig{
		Host:            "localhost",
		Port:            5432,
		SSLMode:         "disable",
		MaxConns:        25,
		MinConns:        5,
		MaxConnLifetime: 1 * time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}
}
