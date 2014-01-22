package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
)

type ModelQueries struct {
	Insert string
}

type Model struct {
	schema  interface{}
	table   string
	typeOf  reflect.Type
	ptrOf   reflect.Type
	queries ModelQueries
}

type Query struct {
	*Model
	*sqlx.DB
}

func New(schema interface{}, table string, mq ModelQueries) *Model {
	typeOf := reflect.TypeOf(schema)
	return &Model{
		schema:  schema,
		table:   table,
		typeOf:  typeOf,
		ptrOf:   reflect.PtrTo(typeOf),
		queries: mq,
	}
}

func (m *Model) Table() string {
	return m.table
}

func (m *Model) Q(db *sqlx.DB) *Query {
	return &Query{
		Model: m,
		DB:    db,
	}
}

func (q *Query) ById(id int) (interface{}, error) {
	result := reflect.New(reflect.TypeOf(q.schema))
	query := fmt.Sprintf("select * from %s where id = ?", q.table)
	err := q.Get(result.Interface(), query, id)
	if err != nil {
		return nil, err
	}
	return result.Interface(), nil
}

func (m *Model) T(query string) string {
	return fmt.Sprint(query, m.table)
}

func (q *Query) Insert(item interface{}) error {
	t := reflect.TypeOf(item)
	switch {
	case t == q.typeOf:
	case t == q.ptrOf:
	default:
		return fmt.Errorf("Wrong type for model")
	}

	_, err := q.NamedExec(q.queries.Insert, item)
	return err
}
