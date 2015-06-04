package web

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"reflect"
	"time"
)

var (
	DB     *sql.DB
	DBName string = "./db.db"
)

func typeToSqlite3Type(i interface{}) string {
	switch i.(type) {
	case string:
		return "text"
	case bool:
		return "boolean"
	case float64:
		return "float"
	case []byte:
		return "blob"
		/*
			case t == uint || t == int || t == uint8 || t == int8 || t == uint16 || t == int16 || t == uint32 || t == int32:
				return "integer"
		*/
	case uint64:
		return "integer"
	case time.Time:
		return "timestamp"
	default:
		log.Fatalf("Unknown golang type %q", reflect.TypeOf(i))
	}
	return ""
}

func structToSQLCreate(i interface{}) string {
	t := reflect.Indirect(reflect.ValueOf(i)).Type()
	n := t.NumField()
	s := "create table " + t.Name() + " ("
	for i := 0; i < n; i++ {
		r := t.Field(i)
		s += r.Name
		s += " " + typeToSqlite3Type(reflect.Zero(r.Type).Interface())
		s += " " + r.Tag.Get("sql")
		if i != n-1 {
			s += ", "
		}
	}
	s += ");"
	return s
}

func DBInit(tables ...interface{}) {
	exists := FileExists(DBName)
	var err error
	DB, err = sql.Open("sqlite3", DBName)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		for _, t := range tables {
			s := structToSQLCreate(t)
			_, err = DB.Exec(s)
			if err != nil {
				os.Remove(DBName)
				log.Fatal(err)
			}
		}
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

//func DBGet(table i
