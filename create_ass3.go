package main
import(
	"fmt"
	_ "github.com/go-sql-driver/mysql" // this is the driver for mysql
	sqlx "github.com/jmoiron/sqlx"      // this is the connector, both package are external packages that you need to use `go install` to install before use
)


const (
	User     = "root"
	Password = ""
	DBName   = "ass3"
)

type Library struct {
	db *sqlx.DB
}

func (lib *Library) ConnectDB() {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	if err != nil {
		panic(err)
	}
	lib.db = db
}

// CreateTables created the tables in MySQL
func (lib *Library) CreateTables() error {
	dsql1:=`
	DROP TABLE IF EXISTS student`
	dsql2:=`
	DROP TABLE IF EXISTS book`
	dsql3:=`
	DROP TABLE IF EXISTS record`
	sql1:=`
	CREATE TABLE student
	  	(id INTEGER NOT NULL,
		available SMALLINT DEFAULT 0,
  		PRIMARY KEY(id))`
	sql2:=`
	CREATE TABLE book
  		(id INTEGER NOT NULL,
  		name VARCHAR(32),
  		author VARCHAR(32),
  		publisher VARCHAR(32),
		available SMALLINT DEFAULT 0,
		reason VARCHAR(80),
  		PRIMARY KEY(id))`
	sql3:=`
	CREATE TABLE record
  		(book_id INTEGER,
  		student_id INTEGER,
  		time DATE,
		rtime DATE,
		delay SMALLINT CHECK (delay >= 0 AND delay <= 3),
		re SMALLINT DEFAULT 0,
  		PRIMARY KEY(book_id,student_id,time),
  		FOREIGN KEY(book_id)REFERENCES book(id),
  		FOREIGN KEY(student_id)REFERENCES student(id))`

	sql := []string{dsql1,dsql2,dsql3,sql1,sql2,sql3}
	for i:=0;i<6;i++{
		smt,err:=lib.db.Prepare(sql[i])
		if err!=nil{
			return err
		}
		smt.Exec()
	}
	return nil
}

func main() {
	var lib *Library
	lib = new (Library)
	lib.ConnectDB()
	defer lib.db.Close()
	err := lib.CreateTables()
	if err!=nil {
		panic(err)
	}
}
