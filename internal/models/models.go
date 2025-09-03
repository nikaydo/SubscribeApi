package models

import (
	"encoding/json"
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

func ParseFilters(b io.Reader) (Filters, error) {
	var parsVar Filters
	if err := json.NewDecoder(b).Decode(&parsVar); err != nil {
		return parsVar, err
	}
	var err error
	if parsVar.StartDate != "" {
		parsVar.Date.StartDate, err = time.Parse("02.01.2006", "01."+parsVar.StartDate)
		if err != nil {
			return parsVar, err
		}
		if parsVar.Date.StartDate.Before(time.Now()) {
			return parsVar, err
		}
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

	return parsVar, nil
}
