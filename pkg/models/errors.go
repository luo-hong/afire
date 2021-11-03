package models

import "errors"

var (
	ErrNoFound        = errors.New("no found")
	ErrNoField        = errors.New("no field in struct")
	ErrEmptyField     = errors.New("empty select field")
	ErrNoPrimarykey   = errors.New("no primary key")
	ErrNoRowsAffected = errors.New("no rows affected")
	ErrNoDate         = errors.New("date is zero")
	ErrNoFoundMachine = errors.New("no found machinecode")
	ErrNoFoundSubPro  = errors.New("no found subprojectcode")
)
