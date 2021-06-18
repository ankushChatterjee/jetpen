package models

import (
	"time"

	"github.com/lib/pq"
)

type Newsletter struct {
	Id          string
	Name        string
	Description string
	Owner       string
	CreatedAt   time.Time
}

type Letter struct {
	Id          string
	Subject     string
	Owner       string
	Nid         string
	Content     string
	CreatedAt   time.Time
	PublishedAt pq.NullTime
	IsPublished bool
}