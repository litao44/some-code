package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	num  = 10000
	age  = 10
	name = "test-name"
)

var (
	db  *sql.DB
	orm *gorm.DB
	dsn = "root:@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local&timeout=3s"
)

type Test struct {
	ID   int64
	Name string
	Age  int64
}

func init() {
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	orm, err = gorm.Open("mysql", db)
	if err != nil {
		panic(err)
	}
}

func main() {
	orm.DropTableIfExists(&Test{})
	orm.CreateTable(&Test{})
	ormStart := time.Now().UnixNano()
	for i := 0; i < num; i++ {
		orm.Create(&Test{Name: name, Age: age})
	}
	ormEnd := time.Now().UnixNano()

	fmt.Printf("orm create %d obj time %d\n", num, ormEnd-ormStart)

	orm.DropTableIfExists(&Test{})
	orm.CreateTable(&Test{})
	ormStart = time.Now().UnixNano()
	for i := 0; i < num; i++ {
		_, err := db.Exec("insert into tests (name,age) value (?,?)", name, age)
		if err != nil {
			panic(err)
		}
	}
	ormEnd = time.Now().UnixNano()
	fmt.Printf("sql create %d obj time %d\n", num, ormEnd-ormStart)

	orm.DropTableIfExists(&Test{})
	orm.CreateTable(&Test{})
	st, _ := db.Prepare("insert into tests (name,age) value (?,?)")
	ormStart = time.Now().UnixNano()
	for i := 0; i < num; i++ {
		_, err := st.Exec(name, age)
		if err != nil {
			panic(err)
		}
	}
	ormEnd = time.Now().UnixNano()
	fmt.Printf("prepare create %d obj time %d\n", num, ormEnd-ormStart)

	orm.DropTableIfExists(&Test{})
	orm.CreateTable(&Test{})
	ormStart = time.Now().UnixNano()
	sql := "insert into tests (name,age) values "
	args := []interface{}{}
	for i := 0; i < num; i++ {
		if i != num-1 {
			sql += "(?,?),"
		} else {
			sql += "(?,?)"
		}
		args = append(args, name, age)
	}

	_, err := db.Exec(sql, args...)
	if err != nil {
		panic(err)
	}
	ormEnd = time.Now().UnixNano()
	fmt.Printf("batch create %d obj time %d\n", num, ormEnd-ormStart)
}
