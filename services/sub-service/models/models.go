package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type Subscription struct {
	Id        string
	Email     string
	Nid       string
	CreatedAt *time.Time
	SubToken  string
}

func (sub *Subscription) GenerateSubToken()  {
	sub.SubToken = uuid.NewV4().String()
}

func (sub *Subscription) GenerateID() {
	sub.Id = uuid.NewV4().String()
}