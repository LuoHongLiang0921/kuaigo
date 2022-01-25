package gormv2

import (
	"errors"

	"gorm.io/gorm"
)

type (
	DB      = gorm.DB
	Config  = gorm.Config
	Session = gorm.Session
)

var (
	errSlowCommand    = errors.New("mysql slow command")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

var (
	Open = gorm.Open
)
