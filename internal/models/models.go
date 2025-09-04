package models

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/google/uuid"
)

type Subscribe struct {
	Id          int       `json:"id,omitempty"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date,omitzero"`
	Date        Date      `json:"-"`
}

type Date struct {
	StartDate time.Time
	EndDate   time.Time
}

type Summary struct {
	Symmary int `json:"summary"`
}

type Id struct {
	Id uint `json:"id"`
}

func ParseResponse(b io.Reader) (Subscribe, error) {
	var parsVar Subscribe
	if err := json.NewDecoder(b).Decode(&parsVar); err != nil {
		return parsVar, err
	}
	var err error
	parsVar.Date.StartDate, err = time.Parse("02.01.2006", "01."+parsVar.StartDate)
	if err != nil {
		return parsVar, err
	}
	if parsVar.EndDate != "" {
		parsVar.Date.EndDate, err = time.Parse("02.01.2006", "01."+parsVar.EndDate)
		if err != nil {
			return parsVar, err
		}
		if parsVar.Date.EndDate.After(time.Now()) {
			return parsVar, err
		}
	}
	if parsVar.Date.StartDate.Before(time.Now()) {
		return parsVar, err
	}
	return parsVar, nil
}

func ParseJson(s []Subscribe) ([]byte, error) {
	byteJson, err := json.Marshal(s)
	if err != nil {
		return byteJson, err
	}
	return byteJson, nil
}

type Filters struct {
	ServiceName string    `json:"service_name,omitempty"`
	UserId      uuid.UUID `json:"user_id,omitempty"`
	StartDate   string    `json:"start_date,omitempty"`
	EndDate     string    `json:"end_date,omitzero"`
	Date        Date      `json:"-"`
	Summary     bool      `json:"summary"`
}

func (f Filters) ParseFilters() (Filters, error) {
	var err error

	if err := f.Check(); err != nil {
		return f, err
	}
	if f.StartDate != "" {
		f.Date.StartDate, err = time.Parse("02.01.2006", "01."+f.StartDate)
		if err != nil {
			return f, err
		}
		if f.Date.StartDate.Before(time.Now()) {
			return f, err
		}
	}

	if f.EndDate != "" {
		f.Date.EndDate, err = time.Parse("02.01.2006", "01."+f.EndDate)
		if err != nil {
			return f, err
		}
		if f.Date.EndDate.After(time.Now()) {
			return f, err
		}
	}
	return f, nil
}

func (f *Filters) Check() error {
	if f.ServiceName == "" && f.EndDate == "" && f.StartDate == "" && f.UserId == uuid.Nil {
		return errors.New("need chose filter")
	}
	return nil
}
