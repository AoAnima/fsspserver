package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/jackc/pgconn"
	"gopkg.in/ldap.v2"
	"time"
)

//func DBconn() *sql.DB {
//	ecpProp := fmt.Sprintf("user=%s password=%s port=%d dbname=%s host=%s sslmode=disable search_path=%s", "postgres","prx318#0", 5432, "postgres", "10.26.4.27", "fssp_chat")
//	conn, err := sql.Open("postgres", ecpProp)
//
//	if err != nil {
//		log.Printf("err %+v\n", err)
//		return nil;
//	}
//	return conn
//}

type DBConfigStruct struct {
	Host       string // Имя основного хоста по умолчаю пока что zoo.ru
	DBname     string // Имя основной базы rout
	DBuser     string
	DBpassword string
	DBhost     string // ip или url адрес БД
	DBport     int64
	SearchPath string
}

var DBConfig = DBConfigStruct{
	//DBname:     "fsspsk",
	DBname:     "postgres",
	DBuser:     "postgres",
	DBpassword: "prx318#0",
	//DBhost:     "10.26.6.15",
	DBhost:     "10.26.4.27",
	DBport:     5432,
	SearchPath: "fssp_chat",
}
var IoDBConfig = DBConfigStruct{
	//DBname:     "fsspsk",
	DBname:     "io",
	DBuser:     "io",
	DBpassword: "anima@iO",
	//DBpassword: "io@Bot",
	//DBhost:     "10.26.6.15",
	//DBhost:     "10.26.6.25",
	DBhost:     "127.0.0.1",
	DBport:     5432,
	SearchPath: "iobot",
}

/*
PGConnect
Нжно передать строку с именем схемы в БД или пустую строку, если пустая строка то используеться IoDBConfig.SearchPath = iobot
*/
func PGConnect(schema string, ctx context.Context) (*pgconn.PgConn, error) {
	var cancel context.CancelFunc
	if ctx == nil {
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
	}
	if schema == "" {
		schema = IoDBConfig.SearchPath
	}

	connString := fmt.Sprintf("user=%s password=%s port=%d dbname=%s host=%s sslmode=disable search_path=%s", IoDBConfig.DBuser, IoDBConfig.DBpassword, IoDBConfig.DBport, IoDBConfig.DBname, IoDBConfig.DBhost, schema)

	PgConn, err := pgconn.Connect(ctx, connString)
	if err != nil {
		Ошибка("pgconn Ошибка подключения к базе:", err)
		connString := fmt.Sprintf("user=%s password=%s port=%d dbname=%s host=%s sslmode=disable search_path=%s", IoDBConfig.DBuser, IoDBConfig.DBpassword, IoDBConfig.DBport, IoDBConfig.DBname, "10.26.6.25", schema)
		PgConn, err = pgconn.Connect(ctx, connString)
		if err != nil {
			Ошибка(" %+v ", err)
		}
	}

	return PgConn, err
}

func MVVPGConnect(schema string, ctx context.Context) (*pgconn.PgConn, error) {
	var cancel context.CancelFunc
	if ctx == nil {
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
	}
	if schema == "" {
		schema = DBConfig.SearchPath
	}
	ИмяБД := DBConfig.DBname
	if schema == "ecp" {
		ИмяБД = "ecp"
	}

	connString := fmt.Sprintf("user=%s password=%s port=%d dbname=%s host=%s sslmode=disable", DBConfig.DBuser, DBConfig.DBpassword, DBConfig.DBport, ИмяБД, DBConfig.DBhost)
	Инфо("%+v", connString)
	var err error

	//ctx := context.Background()
	PgConn, err := pgconn.Connect(ctx, connString)
	if err != nil {

		Ошибка("pgconn Ошибка подключения к базе:", err)
	}

	return PgConn, err
}

func FSSPconnect() *sql.DB {
	dbconfig := "sysdba:512613@10.26.4.243:3050/ncore-fssp?lc_ctype=WIN1251"

	conn, err := sql.Open("firebirdsql", dbconfig)

	//log.Printf(" conn.Ping() %+v\n", conn.Stats())
	//log.Printf(" conn.Ping() %+v\n",  conn.Ping())

	if err := conn.Ping(); err != nil {
		Ошибка(" conn.Ping() %+v\n", conn.Ping())
		return nil
	}
	if err != nil {
		Ошибка("Ошибка подключения к firebirdsql %+v, Ваша конфигурация %+v\n", err, dbconfig)
		return nil
	}
	return conn
}

func RDBconnect() *sql.DB { // Региональная база данных со всеми производствами с края, обновление данных два раза в сутки
	dbconfig := "sysdba:512613@10.26.4.244:3050/ncore-rbd?lc_ctype=WIN1251"
	conn, err := sql.Open("firebirdsql", dbconfig)
	if err != nil {
		Ошибка("Ошибка подключения к firebirdsql %+v, Ваша конфигурация %+v\n", err, dbconfig)
		return nil
	}
	return conn
}

func LdapConnect() (*ldap.Conn, error) {
	//addr :=fmt.Sprintf("%s:%d", "10.26.4.240", 389)
	ldapConn, err := ldap.Dial("tcp", "10.26.4.240:389")
	if err != nil {
		Ошибка("%+v\n", err)

		return nil, err
	}
	err = ldapConn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		Ошибка("%+v\n", err)

		return nil, err
	}
	err = ldapConn.Bind("cn=root,dc=fssprus,dc=ru", "123456")
	if err != nil {
		Ошибка("%+v\n", err)

		return nil, err
	}
	return ldapConn, err
}

//func LdapConnect()( *ldap.Conn, error) {
//	//addr :=fmt.Sprintf("%s:%d", "10.26.4.240", 389)
//	ldap.DefaultTimeout = 60
//	ldapConn, err := ldap.Dial("tcp", "10.26.4.240:389")
//	if err != nil {
//		Ошибка("%+v\n", err)
//		return nil, err
//	}
//	err = ldapConn.StartTLS(&tls.Config{InsecureSkipVerify: true})
//	if err != nil {
//		Ошибка("%+v\n", err)
//
//		return nil, err
//	}
//	err = ldapConn.Bind("cn=root,dc=fssprus,dc=ru", "123456")
//	if err != nil {
//		Ошибка("%+v\n", err)
//
//		return nil, err
//	}
//	return ldapConn, err
//}

func OSPconnect(pwd string, ip string) *sql.DB {

	dbconfig := "sysdba:" + pwd + "@" + ip + ":3050/ncore-fssp?lc_ctype=WIN1251"
	var conn *sql.DB
	//var err error
	//to := time.After(3 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//Инфо(" %+v", ctx)
	//select {
	//case <-time.After(3*time.Second):
	//	Инфо(" время вышло %+v", "3 сек")
	//	return nil
	//	break
	//default:
	//	conn, err := sql.Open("firebirdsql", dbconfig)
	//	if	err!=nil{
	//		Ошибка("Ошибка подключения к firebirdsql %+v, Ваша конфигурация %+v\n", err, dbconfig)
	//		return nil
	//	}
	//	return conn
	//}
	//done := make(chan bool, 1)
	//go func() {
	//	defer fmt.Println("Выход из горутины")
	//	for {
	//		select {
	//		case  <-to:
	//			fmt.Println("Время истекло")
	//			done <- false
	//			break
	//			return
	//		default:
	//			//Инфо("Устанавлиаем соединение %+v", dbconfig)
	//			//	conn, err = sql.Open("firebirdsql", dbconfig)
	//			//	if	err!=nil{
	//			//		Ошибка("Ошибка подключения к firebirdsql %+v, Ваша конфигурация %+v\n", err, dbconfig)
	//			//
	//			//		done <- false
	//			//		break
	//			//		return
	//			//	}
	//			//	Инфо("Удачное подключение conn %+v dbconfig %+v", conn, dbconfig)
	//			//done <- true
	//			//break
	//			//return
	//		}
	//	}
	//}()
	//go func(){
	//	for {
	//		рез := <- done
	//
	//		fmt.Println("Время истекло", рез)
	//
	//	}
	//	Инфо(" %+v", )
	//}()
	if Ping(ip) {
		conn = Соединится(ctx, dbconfig)
		//Инфо("Соединится %+v", conn)
	} else {
		return nil
	}

	return conn

}

func Соединится(ctx context.Context, dbconfig string) *sql.DB {
	conn, err := sql.Open("firebirdsql", dbconfig)
	//<- done
	if err := conn.Ping(); err != nil {
		Ошибка(" conn.Ping() %+v\n", conn.Ping())
		return nil
	}
	if err != nil {
		Ошибка("Ошибка подключения к firebirdsql %+v, Ваша конфигурация %+v\n", err, dbconfig)
	} else {
		Инфо("соединение установленно conn %+v dbconfig %+v", conn, dbconfig)
	}
	return conn
}
