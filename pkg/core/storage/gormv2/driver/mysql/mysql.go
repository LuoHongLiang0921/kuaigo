package mysql

import "gorm.io/driver/mysql"

type (
	Config    = mysql.Config
	Column    = mysql.Column
	Dialector = mysql.Dialector
	Migrator  = mysql.Migrator
)

var (
	Open   = mysql.Open
	New    = mysql.New
	Update = mysql.Update
)
