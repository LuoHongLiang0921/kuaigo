// @description gorm.io/gorm/clause clause 别名

package clause

import "gorm.io/gorm/clause"

type (
	Clause        = clause.Clause
	OnConflict    = clause.OnConflict
	Delete        = clause.Delete
	Table         = clause.Table
	Update        = clause.Update
	Insert        = clause.Insert
	Column        = clause.Column
	Gt            = clause.Gt
	Gte           = clause.Gte
	GroupBy       = clause.GroupBy
	Where         = clause.Where
	Expr          = clause.Expr
	Eq            = clause.Eq
	Select        = clause.Select
	From          = clause.From
	Like          = clause.Like
	Limit         = clause.Limit
	OrConditions  = clause.OrConditions
	OrderBy       = clause.OrderBy
	OrderByColumn = clause.OrderByColumn
	AndConditions = clause.AndConditions
	Assignment    = clause.Assignment

	Set       = clause.Set
	Values    = clause.Values
	IN        = clause.IN
	Lt        = clause.Lt
	Lte       = clause.Lte
	NamedExpr = clause.NamedExpr
	Neq       = clause.Neq
	Join      = clause.Join
	Locking   = clause.Locking
)

var (
	Assignments       = clause.Assignments
	AssignmentColumns = clause.AssignmentColumns
)
