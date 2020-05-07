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


type books struct{
	id int
	name string
	author string
	publisher string
	available int
	reason string
}
type records struct{
	book_id int
	student_id int
	time string
	rtime string
	delay int
	re int
}

// AddBook adds a book into the library
func (lib *Library) AddBook(id int,title, auther, ISBN string) error {
	sql:=`
	INSERT INTO book
		VALUES(?,?,?,?,0,"")`
	smt,err:=lib.db.Prepare(sql)
	if err!=nil {
		return err
	}
	_,err=smt.Exec(id,title,auther,ISBN)
	return err
}

//RemoveBook removes a book from the library
func (lib *Library) RemoveBook(id int,reason string) error{
	sql:=`
	UPDATE book
	SET available=1,reason=?
	WHERE id=?`
	smt,err:=lib.db.Prepare(sql)
	if err!=nil{
		return err
	}
	_,err=smt.Exec(reason,id)
	return err
}

//AddStudent adds a student into the student account
func (lib *Library) AddStudent(id int) error{
	sql:=`
	INSERT INTO student
		VALUES(?,0)`
	smt,err:=lib.db.Prepare(sql)
	if err!=nil{
		return err
	}
	_,err=smt.Exec(id)
	return err
}

//Query a book from Library 
func (lib *Library) Query_id(id int) error{
	b:=books{}
	sql:=`
	SELECT *
	FROM book
	WHERE id=?`
	smt,_:=lib.db.Prepare(sql)
	rows,err:=smt.Query(id)
	defer rows.Close()
	flag:=0//flag=1 if result  is not nil
	if err!=nil{
		return err
	}	
	if rows.Next(){
		rows.Scan(&(b.id),&(b.name),&(b.author),&(b.publisher),&(b.available),&(b.reason))
		flag=1
	}
	if flag==0{
		fmt.Println("No such book\n")
		return nil
	}
	fmt.Println(b,"\n")
	return nil
}
func (lib *Library) Query_name(name string) error{
	sql:=`
	SELECT *
	FROM book
	WHERE name=?`
	smt,_:=lib.db.Prepare(sql)
	rows,err:=smt.Query(name)
	defer rows.Close()
	if err!=nil{
		return err
	}
	b:=books{}
	flag:=0
	for i:=0;rows.Next();i++{
		err=rows.Scan(&(b.id),&(b.name),&(b.author),&(b.publisher),&(b.available),&(b.reason))
		if err!=nil{
			return err
		}
		fmt.Println(b)
		flag=1
	}
	if flag==0{
		fmt.Println("No such book")
	}
	fmt.Print("\n")
	return nil
}
func (lib *Library) Query_publisher(publisher string) error{
	sql:=`
	SELECT *
	FROM book
	WHERE publisher=?`
	smt,_:=lib.db.Prepare(sql)
	rows,err:=smt.Query(publisher)
	defer rows.Close()
	if err!=nil{
		return err
	}
	b:=books{}
	flag:=0
	for i:=0;rows.Next();i++{
		err=rows.Scan(&(b.id),&(b.name),&(b.author),&(b.publisher),&(b.available),&(b.reason))
		if err!=nil{
			return err
		}
		fmt.Println(b)
		flag=1
	}
	if flag==0{
		fmt.Println("No such book")
	}
	fmt.Print("\n")
	return nil
}
func (lib *Library) Query_author(author string) error{
	sql:=`
	SELECT *
	FROM book
	WHERE author=?`
	smt,_:=lib.db.Prepare(sql)
	rows,err:=smt.Query(author)
	defer rows.Close()
	if err!=nil{
		return err
	}
	b:=books{}
	flag:=0
	for i:=0;rows.Next();i++{
		err=rows.Scan(&(b.id),&(b.name),&(b.author),&(b.publisher),&(b.available),&(b.reason))
		if err!=nil{
			return err
		}
		fmt.Println(b)
		flag=1
	}
	if flag==0{
		fmt.Println("No such book")
	}
	fmt.Print("\n")
	return nil
}

//BorrowBook
func (lib *Library) BorrowBook(bid,sid int) error{	
	//check if the book_id exists
	sql:=`
	SELECT id
	FROM book
	WHERE id=? AND available =0`
	smt,_:=lib.db.Prepare(sql)
	rows,err:=smt.Query(bid)
	defer rows.Close()
	flag:=0
	if err != nil{
		return err
	}
	if rows.Next(){
		flag=1
	}
	if flag==0{
		fmt.Println("No such book\n")
		return nil
	}
	
	fmt.Println("\nbook found")
	
	flag=0
	//check if this account is available
	sql=`
	SELECT id
	FROM student
	WHERE id=? AND available = 0`
	smt,_=lib.db.Prepare(sql)
	rows,err=smt.Query(sid)	
	defer rows.Close()
	if err != nil{
		return err
	}
	if rows.Next(){
		flag=1
	}
	if flag==0{
		fmt.Println("No such student account or permission denied\n")
		return nil	
	}

	fmt.Println("student found")
	
	//add record
	sql=`
	INSERT INTO record
		VALUES(?,?,CURDATE(),DATE_ADD(CURDATE(),INTERVAL 1 MONTH),0,0)`
	smt,_=lib.db.Prepare(sql)
	_,err=smt.Exec(bid,sid)
	if err!=nil{
		return err
	}
	sql=`
	UPDATE book
	SET available = 1
	WHERE id=?`
	smt,_=lib.db.Prepare(sql)
	_,err=smt.Exec(bid)
	if err!=nil{
		return err
	}
	fmt.Println("Success!\n")
	return nil
}

//QueryrRecord
func (lib *Library) QueryRecord(id int) error{
	sql:=`
	SELECT *
	FROM record
	WHERE student_id=?`
	smt,_:=lib.db.Prepare(sql)
	rows,err:=smt.Query(id)
	defer rows.Close()
	if err!=nil {
		return err
	}
	r:= records{}
	flag:=0
	for i:=0;rows.Next();i++{
		err=rows.Scan(&(r.book_id),&(r.student_id),&(r.time),&(r.rtime),&(r.delay),&(r.re))
		if err!=nil{
			return err
		}
		fmt.Println(r)
		flag=1
	}
	if flag==0{
		fmt.Println("No such record!")
	}
	fmt.Print("\n")
	return nil
}

//ReturnBook
func (lib *Library) ReturnBook(bid int,sid int) error{
	sql:=`
	SELECT book_id
	FROM record
	WHERE book_id=? AND student_id=? AND re = 0`
	smt,_:=lib.db.Prepare(sql)
	rows,err:=smt.Query(bid,sid)
	defer rows.Close()
	if err!=nil {
		return err
	}
	flag:=0
	if rows.Next(){
		flag=1
	}
	if flag==0{
		fmt.Println("No such record!\n")
		return nil
	}
	
	sql=`
	UPDATE record
	SET re=1
	WHERE book_id=? AND student_id=? AND re = 0`
	smt,_=lib.db.Prepare(sql)
	_,err=smt.Exec(bid,sid)
	if err!=nil{
		return err
	}
	sql=`
	UPDATE book
	SET available = 0
	WHERE id=?`
	smt,_=lib.db.Prepare(sql)
	_,err=smt.Exec(bid)
	if err!=nil{
		return err
	}
	fmt.Println("Success!\n")
	return nil
	
}

//Query_notReturn queries the books borrowed by a student which have not been returned
func (lib *Library) Query_notReturn(id int) error{
	sql:=`
	SELECT book_id
	FROM record
	WHERE student_id=? AND re=0`
	smt,_:=lib.db.Prepare(sql)
	rows,err:=smt.Query(id)
	defer rows.Close()
	if err!=nil {
		return err
	}
	var bid int
	flag:=0
	for rows.Next(){
		err=rows.Scan(&bid)
		if err!=nil{
			return err
		}
		fmt.Println("book id: ")
		fmt.Println(bid)
		flag=1
	}
	if flag==0{
		fmt.Println("No record find!")
	}
	fmt.Print("\n")
	return nil
}

//CheckAccount checks if the student's account is available
func (lib *Library) CheckAccount(id int) error{
	sql:=`
	SELECT COUNT(*)
	FROM record
	WHERE student_id=? AND re=0 AND delay>0`
	smt,_:=lib.db.Prepare(sql)
	rows,err:=smt.Query(id)
	defer rows.Close()
	if err!=nil {
		return err
	}
	flag:=0
	var c int
	for rows.Next(){
		err=rows.Scan(&c)
		if err!=nil{
			return err
		}
		flag=1
	}
	if flag==0{
		return nil
	}
	if c>3{
		sql=`
		UPDATE student
		SET available=1
		WHERE id=?`
		smt,_:=lib.db.Prepare(sql)
		_,err = smt.Exec(id)
		if err!=nil{
			return err
		}
	}else{
		sql=`
		UPDATE student
		SET available=0
		WHERE id=?`
		smt,_=lib.db.Prepare(sql)
		_,err = smt.Exec(id)
		if err!=nil{
			return err
		}
	}
	return nil
	
}

//CheckOverdue
func (lib *Library) CheckOverdue(id int) error{
	sql:=`
	SELECT book_id
	FROM record
	WHERE student_id=? AND re = 0 AND DATEDIFF(rtime,CURDATE())<0`
	smt,_:=lib.db.Prepare(sql)
	rows,err:=smt.Query(id)
	defer rows.Close()
	if err!=nil {
		return err
	}
	var bid int
	flag:=0
	for rows.Next(){
		err=rows.Scan(&bid)
		if err!=nil{
			return err
		}
		fmt.Println(bid)
		flag=1
	}
	if flag==0{
		fmt.Println("No overdue record")
	}
	fmt.Print("\n")
	return nil
}

//CheckDeadline
func (lib *Library) CheckDeadline(id int) error{
	sql:=`
	SELECT rtime,DATEDIFF(rtime,CURDATE()) days
	From record
	WHERE book_id=? AND re=0`
	var t string
	var d int
	smt,_:=lib.db.Prepare(sql)
	rows,err:=smt.Query(id)
	defer rows.Close()
	if err!=nil {
		return err
	}
	flag:=0
	if rows.Next(){
		err=rows.Scan(&t,&d)
		if err!=nil{
			return err
		}
		flag=1
	}
	if flag==0{
		fmt.Println("No such record")
		return nil
	}
	fmt.Println(t)
	fmt.Println("Days to due: ",d,"\n")
	return nil
}

//Delay operation
func (lib *Library) Delay(bid,sid int) error{
	sql:=`
	SELECT delay
	FROM record
	WHERE book_id=? AND student_id=? AND re=0`
	var d int
	smt,_:=lib.db.Prepare(sql)
	rows,err:=smt.Query(bid,sid)
	defer rows.Close()
	if err!=nil{
		return err
	}
	flag:=0
	if rows.Next(){
		err=rows.Scan(&d)
		if err!=nil{
			return err
		}
		flag=1
	}
	if flag==0{
		fmt.Println("No such record")
		return nil
	}
	if d >= 3{
		fmt.Println("Times limitted, cannot delay\n")
		return nil
	}
	sql=`
	UPDATE record
	SET delay = delay+1, rtime=DATE_ADD(rtime,INTERVAL 1 MONTH)
	WHERE book_id=? AND student_id = ?`
	smt,_=lib.db.Prepare(sql)
	_,err=smt.Exec(bid,sid)
	if err!=nil{
		return err
	}
	fmt.Println("Success!\n")
	return nil
}

//LookUp searches for an account if it exists
func (lib *Library) LookUp(id int) (bool,error){
	sql:=`
	SELECT id
	FROM student
	WHERE id=?`
	smt,_:=lib.db.Prepare(sql)
	rows,err:=smt.Query(id)
	defer rows.Close()
	if err!=nil {
		return false,err
	}
	if rows.Next(){
		return true,nil
	}
	return false,nil
}

func main() {
	fmt.Println("Welcome to the Library Management System!")
	var lib *Library
	lib = new (Library)
	lib.ConnectDB()
	defer lib.db.Close()
	var err error
	fmt.Println("Please sign in")
	fmt.Println("user name:(Please use '10000' if you are administrator)")
	flag:=0//flag if user is 10000
	var user int
	fmt.Scan(&user)
	if user==10000 {
		fmt.Println("password:")
		var pwd string
		fmt.Scanln(&pwd)
		if pwd=="root"{
			flag=1
		}
	} else {
		t_f,err:=lib.LookUp(user);
		if err!=nil{
			panic(err)
		}
		if t_f==false{
			fmt.Println("No such account")
			return
		}
	}
	var code int
	b:=books{}
	if flag==1{
		//admin view
		for{
			fmt.Println("Please type a code to choose an operation")
			fmt.Println("0 for exiting")
			fmt.Println("1 for adding a student account")
			fmt.Println("2 for adding a book into library")
			fmt.Println("3 for removing a book from library")
			fmt.Println("4 for querying a book")
			fmt.Println("5 for querying the borrowing history of a student account")
			fmt.Println("6 for querying the books a student has borrowed and not returned yet")
			fmt.Println("7 for checking the deadline of returning a borrowed book")
			fmt.Println("8 for extending the deadline of returning a book, at most 3 times")
			fmt.Println("9 for checking if a student has any overdue books that needs to be returned")
			fmt.Scanln(&code)
			if code==0{
				break;
			}
			var sid int
			var err error
			switch (code){
				case 1:
					fmt.Println("Please type 0 for return")
					fmt.Println("Please type the student id (only numbers allowed)")
					fmt.Print("student id: ")
					fmt.Scanln(&sid)
					err=lib.AddStudent(sid)
					if err!=nil{
						panic(err)
					}
				case 2:
					fmt.Println("Please type the information of this book")
					fmt.Print("id: ")
					fmt.Scanln(&(b.id))
					fmt.Print("name: ")
					fmt.Scanln(&(b.name))
					fmt.Print("author: ")
					fmt.Scanln(&(b.author))
					fmt.Print("publisher: ")
					fmt.Scanln(&(b.publisher))
					err=lib.AddBook(b.id,b.name,b.author,b.publisher)
					if err!=nil{
						panic(err)
					}
				case 3:
					fmt.Println("Please type id of the book you want to remove")
					fmt.Print("book id: ")
					fmt.Scanln(&(b.id))
					fmt.Print("reason: ")
					fmt.Scanln(&(b.reason))
					err=lib.RemoveBook(b.id,b.reason)
					if err!=nil{
						panic(err)
					}
				case 4:
					fmt.Println("Please type a code to choose a operation")
					fmt.Println("1 for querying id")
					fmt.Println("2 for querying name")
					fmt.Println("3 for querying author")
					fmt.Println("4 for querying publisher")
					fmt.Scanln(&code)
					switch (code){
						case 1:
							fmt.Println("Please type the id")
							fmt.Print("book id: ")
							fmt.Scanln(&(b.id))
							err=lib.Query_id(b.id)
							if err!=nil{
								panic(err)
							}
						case 2:
							fmt.Println("Please type the name")
							fmt.Print("name: ")
							fmt.Scanln(&(b.name))
							err=lib.Query_name(b.name)
							if err!=nil{
								panic(err)
							}
						case 3:
							fmt.Println("Please type the author")
							fmt.Print("author: ")
							fmt.Scanln(&(b.author))
							err=lib.Query_author(b.author)
							if err!=nil{
								panic(err)
							}
						case 4:
							fmt.Println("Please type the publisher")
							fmt.Print("Publisher: ")
							fmt.Scanln(&(b.publisher))
							err=lib.Query_publisher(b.publisher)
							if err!=nil{
								panic(err)
							}
						default:
							fmt.Println("No such operation\n")
					}
				case 5:
					fmt.Println("Please type the student id")
					fmt.Print("id: ")
					fmt.Scan(&sid)
					err=lib.QueryRecord(sid)
					if err!=nil{
						panic(err)
					}
				case 6:
					fmt.Println("Please type the student id")
					fmt.Print("id: ")
					fmt.Scan(&sid)
					err=lib.Query_notReturn(sid)
					if err!=nil{
						panic(err)
					}
				case 7:
					fmt.Println("Please type the book id")
					fmt.Print("id: ")
					fmt.Scanln(&(b.id))
					err=lib.CheckDeadline(b.id)
					if err!=nil{
						panic(err)
					}
				case 8:
					fmt.Println("Please type the book id and the student id")
					fmt.Print("book id: ")
					fmt.Scanln(&(b.id))
					fmt.Print("student id: ")
					fmt.Scanln(&sid)
					err=lib.Delay(b.id,sid)
					if err!=nil{
						panic(err)
					}
				case 9:
					fmt.Println("Please type the student id")
					fmt.Print("id: ")
					fmt.Scan(&sid)
					err=lib.CheckOverdue(sid)
					if err!=nil{
						panic(err)
					}
				default:
					fmt.Println("No such operations\n")
			}
		}
	} else {
		//user view
		for{
			fmt.Println("Please type a code to choose an operation")
			fmt.Println("0 for exiting")
			fmt.Println("1 for borrowing book")
			fmt.Println("2 for returning book")
			fmt.Println("3 for querying the borrowing history of your account")
			fmt.Println("4 for querying a book")
			fmt.Println("5 for querying the books you have borrowed and not returned yet")
			fmt.Println("6 for checking the deadline of returning a borrowed book")
			fmt.Println("7 for extending the deadline of returning a book, at most 3 times")
			fmt.Println("8 for checking if you have any overdue books that needs to be returned")
			fmt.Scanln(&code)
			if code==0{
				break
			}
			switch (code){
				case 1:
					fmt.Println("Please type the book id")
					fmt.Print("id: ")
					fmt.Scanln(&(b.id))
					err=lib.BorrowBook(b.id,user)
					if err!=nil{
						panic(err)
					}
					err=lib.CheckAccount(user)
					if err!=nil{
						panic(err)
					}
				case 2:
					fmt.Println("Please type the book id")
					fmt.Print("id: ")
					fmt.Scanln(&(b.id))
					err=lib.ReturnBook(b.id,user)
					if err!=nil{
						panic(err)
					}
				case 3:
					err=lib.QueryRecord(user)
					if err!=nil{
						panic(err)
					}
				case 4:
					fmt.Println("Please type a code to choose a operation")
					fmt.Println("1 for querying id")
					fmt.Println("2 for querying name")
					fmt.Println("3 for querying author")
					fmt.Println("4 for querying publisher")
					fmt.Scanln(&code)
					switch (code){
						case 1:
							fmt.Println("Please type the id")
							fmt.Print("book id: ")
							fmt.Scanln(&(b.id))
							err=lib.Query_id(b.id)
							if err!=nil{
								panic(err)
							}
						case 2:
							fmt.Println("Please type the name")
							fmt.Print("name: ")
							fmt.Scanln(&(b.name))
							err=lib.Query_name(b.name)
							if err!=nil{
								panic(err)
							}
						case 3:
							fmt.Println("Please type the author")
							fmt.Print("author: ")
							fmt.Scanln(&(b.author))
							err=lib.Query_author(b.author)
							if err!=nil{
								panic(err)
							}
						case 4:
							fmt.Println("Please type the publisher")
							fmt.Print("Publisher: ")
							fmt.Scanln(&(b.publisher))
							err=lib.Query_publisher(b.publisher)
							if err!=nil{
								panic(err)
							}
						default:
							fmt.Println("No such operation\n")
					}
				case 5:
					err=lib.Query_notReturn(user)
					if err!=nil{
						panic(err)
					}
				case 6:
					fmt.Println("Please type the book id")
					fmt.Print("id: ")
					fmt.Scanln(&(b.id))
					err=lib.CheckDeadline(b.id)
					if err!=nil{
						panic(err)
					}
				case 7:
					fmt.Println("Please type the book id")
					fmt.Print("book id: ")
					fmt.Scanln(&(b.id))
					err=lib.Delay(b.id,user)
					if err!=nil{
						panic(err)
					}
				case 8:
					err=lib.CheckOverdue(user)
					if err!=nil{
						panic(err)
					}
				default:
					fmt.Println("No such operations\n")
			}
		}
	}
}
