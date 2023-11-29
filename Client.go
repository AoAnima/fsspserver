package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	_ "github.com/jackc/pgconn"
	"golang.org/x/crypto/ssh"
	"gopkg.in/ldap.v2"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

type UserInfoStruct struct {
	FullName string       `json:"FullName"`
	Initials string       `json:"Initials"`
	Info     *UsersStruct `json:"info"`
	Curator  bool         `json:"curator"`
}
type Client struct {
	Host                   string
	Ip                     string
	НайстройкиПользователя map[string]interface{}
	UserInfo               *UserInfoStruct
	Token                  struct {
		Hash     string
		Истекает string
	}
	//FullName string
	//PostId	int
	//PostName string
	//OspCode int
	//OspName string
	Id         int
	Login      string
	Ssh        *ssh.Client  `json:",omitempty"`
	SSHSession *ssh.Session `json:",omitempty"`

	SshEnv map[string]string

	SshIn  io.WriteCloser
	SshLog bytes.Buffer
	SshErr bytes.Buffer

	Ws *websocket.Conn `json:",omitempty"`

	Message chan *Сообщение `json:",omitempty"`
	DoneCh  chan bool       `json:",omitempty"`

	АктивныеДиалоги []map[string]interface{}
	АктивныйДиалог  map[string]interface{}

	mu               sync.Mutex
	ВходящееСобщение map[int]Сообщение
}
type usersStruct struct {
	OspCode      int            `json:"osp_code"`
	Login        string         `json:"Login"`
	Curator      sql.NullString `json:"curator"`
	Zam          sql.NullString `json:"zam"`
	Ltp          sql.NullString `json:"ltp"`
	Initials     sql.NullString `json:"initials"`
	Givenname    sql.NullString `json:"givenname"`
	SecondName   sql.NullString `json:"second_name"`
	OspName      string         `json:"osp_name"`
	Post         sql.NullString `json:"post"`
	PostName     sql.NullString `json:"post_name"`
	Blocked      sql.NullBool   `json:"blocked"`
	DelayFrom    sql.NullString `json:"delay_from"`
	DelayTo      sql.NullString `json:"delay_to"`
	CuratorPost  string         `json:"curator_post"`
	CuratorOrder int            `json:"curator_order"`
}
type UsersStruct struct {
	Ip             string                   `json:"ip"`
	OspCode        int                      `json:"osp_code"`
	Login          string                   `json:"login"`
	Curator        string                   `json:"curator"`
	Zam            string                   `json:"zam"`
	Ltp            string                   `json:"ltp"`
	Initials       string                   `json:"initials"`
	Givenname      string                   `json:"givenname"`
	SecondName     string                   `json:"second_name"`
	OspName        string                   `json:"osp_name"`
	Post           int                      `json:"post"`
	PostName       string                   `json:"post_name"`
	Blocked        bool                     `json:"blocked"`
	DelayFrom      string                   `json:"delay_from"`
	DelayTo        string                   `json:"delay_to"`
	CuratorPost    string                   `json:"curator_post"`
	CuratorOrder   int                      `json:"curator_order"`
	CuratorOsp     []map[string]interface{} `json:"curator_osp"`
	УровеньДоступа []string                 `json:"уровень_доступа"`
}

func (client *Client) ПолучитьКонтактЛист() map[string][]*UsersStruct {
	//sqlString := `Select  x.osp_code, json_agg(x.*) from
	//			(SELECT ou as osp_code,
	//				   fssp_configs.users.Login,
	//				   fssp_configs.curators.Login as curator, zam,ltp,
	//				   initials,givenname,second_name, osp_name,post,post_name,blocked,delay_from,delay_to
	//			FROM fssp_configs.users
	//					 JOIN fssp_configs.osp_address ON fssp_configs.osp_address.osp_code=fssp_configs.users.ou
	//					 JOIN fssp_configs.posts ON fssp_configs.posts.id= fssp_configs.users.post and fssp_configs.users.post > 0
	//					 left JOIN fssp_configs.curators  ON fssp_configs.curators.Login = fssp_configs.users.Login
	//				AND fssp_configs.curators.osp_code = (SELECT ou FROM fssp_configs.users WHERE Login = $1)
	//			where blocked is null or blocked != TRUE) x  group by x.osp_code`

	//Logins :=[][]byte{}
	//usersOnline := string
	//for login,_:=range server.Clients{
	//	//Logins=append(Logins, []byte(login))
	//}

	sqlString := `select * from fssp_configs.users
        LEFT JOIN fssp_configs.osp_address ON fssp_configs.osp_address.osp_code=fssp_configs.users.ou
        LEFT JOIN fssp_configs.posts ON fssp_configs.posts.id= fssp_configs.users.post
        LEFT JOIN fssp_configs.curators ON fssp_configs.curators.osp_code = (SELECT ou FROM fssp_configs.users WHERE login = $1) AND fssp_configs.users.login= $1
		WHERE fssp_configs.users.login IN ('i.zinchenko@r26','a.stryukova@r26','semenixina@r26') OR (ou=26911 AND (blocked != true OR blocked is null))`

	sqlQuery := sqlStruct{
		Name: "users",
		Sql:  sqlString,
		Values: [][]byte{
			[]byte(client.Login),
		},
		DBSchema: "fssp_configs",
	}

	//CaptionRow, _ := ВыполнитьPgSQL(sqlQuery)
	//Ошибка("CaptionRow %+v\n", CaptionRow)
	//	result := pgconn.ResultReader{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Result, err := sqlQuery.PgSqlResultReader(ctx)
	if err != nil {
		Ошибка("\n !! ERR %+v\n", err)
		return nil
	}
	OSPUsersList := Result.Rows
	//ResultReader:=result.Read()

	//for ResultReader.NextRow() {
	//v := ResultReader.Values()
	//for i,d := range ResultReader.Read(){
	//
	//	Ошибка("i %+v d %+v\n", i,string(d))
	//}

	//}

	ContactList := map[string][]*UsersStruct{}
	for _, osp := range OSPUsersList {
		var UserStruct []*UsersStruct
		//Ошибка("osp %+v\n", string(osp[1]))
		err := json.Unmarshal(osp[1], &UserStruct)
		if err != nil {
			Ошибка("err	 %+v\n", err)
		}

		ContactList[string(osp[0])] = UserStruct

	}

	return ContactList
}

func (client *Client) ПолучитьУровеньДоступа(login string) {

}

/*
ПолучитьДанныеПользователя
Данные о пользователи из локальной бд для бота
*/
func (client *Client) ПолучитьДанныеПользователя(login string) *UserInfoStruct {
	Инфо("ПолучитьДанныеПользователя client %+v login %+v", client, login)
	sqlString := `select jsonb_agg(t.*) from
   (SELECT ou as osp_code,
           jsonb_agg(fssp_configs.curators.*) as curator_osp,
        fssp_configs.users.Login,
           initials,givenname,second_name, osp_name,post,post_name
    FROM fssp_configs.users
             left JOIN fssp_configs.osp_address ON fssp_configs.osp_address.osp_code=fssp_configs.users.ou
             left JOIN fssp_configs.posts ON fssp_configs.posts.id= fssp_configs.users.post and fssp_configs.users.post <> 0
             left JOIN fssp_configs.curators  ON fssp_configs.curators.Login = fssp_configs.users.Login
    where  fssp_configs.users.Login = $1 group by fssp_configs.users.Login, fssp_configs.users.ou, initials,givenname,second_name, osp_name,post,post_name) t`

	UserLogin := login
	if UserLogin == "" {
		UserLogin = client.Login
	}
	ctx, cancel := context.WithCancel(context.Background())
	sqlQuery := sqlStruct{
		Name: "users",
		Sql:  sqlString,
		Values: [][]byte{
			[]byte(UserLogin),
		},
		DBSchema: "fssp_configs",
	}
	defer cancel()

	Result, err := sqlQuery.PgSqlResultReader(ctx)

	if err != nil {
		Ошибка(">>>> ERROR \n %+v \n\n", err)
	}
	//Ошибка(">>>> ДанныеПользователя %+v\n", ДанныеПользователя)
	UserInfo := Result.Rows
	//ДанныеПользователя [map[givenname:Александр initials:Олегович login:maksimchuk@r26 osp_code:26911 osp_name:Отдел информатизации и обеспечения информационной безопасности post:44 post_name:Ведущий специалист-эксперт second_name:Максимчук курируемые_отделы:[нет] уровни_доступа:[admin architect curators user архитектор]]
	var UserStruct []*UsersStruct

	if len(UserInfo[0][0]) > 0 {

		err := json.Unmarshal(UserInfo[0][0], &UserStruct)
		if err != nil {
			Ошибка(" \n err	 %+v UserInfo[0][0] %s  \n \n", err, string(UserInfo[0][0]))
			//Ошибка("\n \n UserInfo %+v \n \n", err, UserInfo)
		}

	} else {
		//Запросим данные в АИС
		Инфо("ПОЛЬЗОВАТЕЛЬ НЕ НАЙДЕН В ПАМЯТИ БОТА, НУЖНО ОБНОВИТЬ ДАННЫЕ ИЗ ЛДАП АИС %+v\n", client.Login)
		Инфо("UserInfo %+v\n", UserInfo)
		Инфо("UserLogin %+v\n", UserLogin)
		Инфо("sqlQuery %+v\n", sqlQuery)

		АктуализироватьДанныеПользователя(client.Login)

		return client.ПолучитьДанныеПользователя("")

	}

	if login != "" {
		if server.Clients[UserLogin] != nil {
			UserStruct[0].Ip = server.Clients[UserLogin].Ip
		}
	} else {
		UserStruct[0].Ip = client.Ip
	}

	UserData := &UserInfoStruct{
		FullName: UserStruct[0].SecondName + " " + UserStruct[0].Givenname + " " + UserStruct[0].Initials,
		Initials: CreatInitails(UserStruct[0].Givenname, UserStruct[0].Initials),
		Info:     UserStruct[0],
		Curator:  false,
	}
	Инфо(" UserData %+v UserStruct[0] %+v", UserData, UserStruct[0])

	if len(UserStruct[0].CuratorOsp) > 0 {
		UserData.Curator = true
	}

	if login == "" { // получаем данные текущего пользователя
		//Инфо(" UserData %+v", client.UserInfo)
		client.UserInfo = UserData
	}
	Инфо("Данные пользователя %+v\n", UserData)
	Инфо(" client %+v", client.UserInfo)
	//	Инфо(" >>>>>>>>>>>>> ДАННЫЕ ОПЛЬЗОВАТЕЛЯ !!!!! UserData %+v  \n", client)
	return UserData
}

type OspStruct struct {
	OspCode    int    `json:"osp_code"`
	OspAddress string `json:"osp_address"`
	OspName    string `json:"osp_name"`
	Ip         string `json:"ip_ais"`
	Pwd        string `json:"pwd"`
	Login      string `json:"uid"`
}

func ПолучитьСписокОСП() map[int]*OspStruct {
	sqlString := `SELECT * FROM fssp_configs.osp_address WHERE site is not null and disabled = false ORDER BY osp_code`
	sqlQuery, err := sqlStruct{
		Name:     "users",
		Sql:      sqlString,
		Values:   nil,
		DBSchema: "fssp_configs",
	}.Выполнить(nil)
	if err != nil {
		Ошибка(">>>> ERROR \n %+v \n\n", err)
	}
	OspData := map[int]*OspStruct{}
	for _, Osp := range sqlQuery {
		//КодОсп, err := strconv.Atoi(Osp["osp_code"].(string))
		//if err != nil {
		//	Ошибка(">>>> ERROR \n %+v \n\n", err)
		//}
		//Ошибка("OSP %+v\n", Osp)

		OspData[Osp["osp_code"].(int)] = &OspStruct{
			OspCode:    Osp["osp_code"].(int),
			OspAddress: Osp["osp_address"].(string),
			OspName:    Osp["osp_name"].(string),
			Ip:         Osp["ip_ais"].(string),
			Pwd:        Osp["pwd"].(string),
			Login:      Osp["uid"].(string),
		}
	}
	return OspData
}

func (client *Client) ПолучитьСписокОтделов() map[string]*OspStruct {
	//sqlString :=`SELECT osp_code, osp_address, osp_name FROM osp_address`
	sqlString := `select t.osp_code, jsonb_agg(t.*) from (SELECT osp_code, osp_address, osp_name FROM fssp_configs.osp_address where disabled = false) t group by t.osp_code`
	sqlQuery := sqlStruct{
		Name:     "users",
		Sql:      sqlString,
		Values:   nil,
		DBSchema: "fssp_configs",
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Result, err := sqlQuery.PgSqlResultReader(ctx)
	if err != nil {
		Ошибка("\n !! ERR %+v\n", err)
		return nil
	}
	OspList := Result.Rows
	OspData := map[string]*OspStruct{}
	if len(OspList) > 0 {
		for _, ospData := range OspList {
			var OSPStruct []*OspStruct
			err := json.Unmarshal(ospData[1], &OSPStruct)
			if err != nil {
				Ошибка("err	 %+v\n", err)
			}
			OspData[string(ospData[0])] = OSPStruct[0]
		}
	}

	return OspData

}

//func ПолучитьСписокОтделов(client *Client) map[int]OspStruct{
//	Ошибка("ПолучитьСписокОтделов %+v\n", )
//	conn := DBconn()
//
//	defer conn.Close()
//	sqlStr :=`SELECT osp_code, osp_address, osp_name FROM osp_address`
//	selectOsp,err := conn.Query(sqlStr) //, client.Login
//	if err !=nil{
//		Ошибка("err %+v\n", err)
//	}
//
//	ospData := OspStruct{}
//	OspList := map[int]OspStruct{}
//	for selectOsp.Next() {
//		err:=selectOsp.Scan(&ospData.OspCode, &ospData.OspAddress, &ospData.OspName)
//		if err!=nil{
//			Ошибка("err %+v\n", err)
//		}
//
//		OspList[ospData.OspCode]=ospData
//	}
//	//Ошибка("OspList %+v\n",OspList )
//	return OspList
//}

func ПолучитьДолжность(login string) (int, string) {
	conn := FSSPconnect()
	if conn == nil {
		Ошибка("Ошибка соединения с РДБ %+v\n")
		return 0, ""
	}
	defer conn.Close()
	//
	sql := "select spi_post_id, SUSER_DATE_ORDER_DISMISS from spi left join sys_users su on spi.suser_id = su.suser_id where suser_name = upper('" + login + "')"
	Ошибка("sql %+v\n", sql)
	selectQ, _ := conn.Query(sql)
	var post_id int
	var date_order string

	for selectQ.Next() {
		selectQ.Scan(&post_id, &date_order)
	}
	//
	return post_id, date_order
}

//func ПолучитьИдДолжностиИзАИС(conn *sql.DB, login string) map[string]interface{} {
func ПолучитьИдДолжностиИзАИС(conn *sql.DB, login string) int {
	//sql := "select spi_post_id, SUSER_DATE_ORDER_DISMISS, SUSER_BLOCKED,SUSER_START_BLOCK ,SUSER_END_BLOCK from spi left join sys_users su on spi.suser_id = su.suser_id where suser_name = upper('"+login+"')"
	//Ошибка("sql %+v\n", sql)

	//sql := "select * from spi  where SPI_EMAIL = '"+login+".fssprus.ru'"
	sql := "select SPI_POST_ID from spi WHERE SPI_EMAIL = '" + login + ".fssprus.ru' AND STRUCTURE_DEPT_NAME IS NOT NULL"

	selectQ, err := conn.Query(sql)
	if err != nil {
		Ошибка("err %+v\n", err)
		return -1
	}
	var post_id int
	//var date_order string
	//var blocked int
	//var start_block string
	//var end_block string
	for selectQ.Next() {
		selectQ.Scan(&post_id) //, &date_order, &blocked, &start_block, &end_block
	}
	//

	//result := map[string]interface{}{
	//	"post_id":post_id,
	//	//"date_order":date_order,
	//	//"blocked":blocked,
	//	//"start_block":start_block,
	//	//"end_block":end_block,
	//}

	return post_id
}
func ПолучитьДолжностьИзАИС(conn *sql.DB, login string) string {
	//sql := "select SPI_POST_NAME from spi WHERE (SPI_EMAIL = '"+login+".fssprus.ru' OR ) AND STRUCTURE_DEPT_NAME IS NOT NULL"
	sql := "select SPI_POST_NAME from spi WHERE SPI_EMAIL LIKE '" + login + "%' AND STRUCTURE_DEPT_NAME IS NOT NULL"

	selectQ, err := conn.Query(sql)
	if err != nil {
		Ошибка("err %+v\n", err)
		return ""
	}

	var post string
	for selectQ.Next() {
		selectQ.Scan(&post)
	}
	Инфо("post %+v login %+v sql %+v", post, login, sql)
	return post
}

func CreatInitails(Fname string, Mname string) string {
	//
	//if Fname == "ИО"{
	//	return Fname
	//}
	FnameRune := []rune(Fname)
	MnameRune := []rune(Mname)
	//Ошибка("string(FnameRune[0]) + string(MnameRune[0]) %+v\n", string(FnameRune[0]) + string(MnameRune[0]))
	return string(FnameRune[0]) + string(MnameRune[0])
}

/*
1. Смотрите атрибут rdbActive.

     Если он TRUE - учетка действует.

     FALSE - заблокирована/приостановлена.

2. Если атрибута rdbActive нет, смотрите значение userPassword. Значение есть - учетка действует. Нет - заблокирована/приостановлена.

3. Если учетка заблокирована/приостановлена - смотрите значения DELAY_FORM и DELAY_TO.

Значение нет - заблокирована (т.е. скорее всего уволена).

Значения есть и захватывают текущую дату - приостановлена (отпуск/больничный).

И + есть поле userDateOrderDismiss - Дата увольнения.
*/
func (client *Client) СинхронизироватьДанныеПользователей(mes Сообщение) map[string]interface{} {
	Инфо(" %+v", "СинхронизироватьДанныеПользователей")
	ldapConn, err := LdapConnect()
	if ldapConn == nil || err != nil {
		if err != nil {
			Ошибка(" Нет подключения к ldap %+v ", err)
		}
		return nil //"Ошибка соединения с сервером авторизации, попытка локальной авторизации"
	}
	defer ldapConn.Close()
	filter := "(objectClass=fsspData)"
	fields := []string{"uidNumber", "sn", "givenName", "initials", "ou", "uid", "cn", "userDateOrderDismiss", "description", "delayFrom", "delayTo", "rdbActive", "userPassword"}
	searchResult, err := ldapConn.Search(ldap.NewSearchRequest(
		"dc=fssprus,dc=ru",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		fields,
		nil,
	))
	FSSPconn := FSSPconnect()
	defer FSSPconn.Close()

	for _, Entri := range searchResult.Entries {

		//for i, у := range Entri.Attributes{
		//	Инфо("  %+v %+v",  i, у)
		//
		//}
		// cn фио
		// sn - фамилия,
		// initials - отчество
		// ou подразделение
		ФИО := Entri.GetAttributeValue("cn")
		Логин := Entri.GetAttributeValue("uid")
		ОСП := Entri.GetAttributeValue("ou")
		UidNumber := Entri.GetAttributeValue("uidNumber")
		ЗаблокированС := Entri.GetAttributeValue("delayFrom")
		ЗаблокированПо := Entri.GetAttributeValue("delayTo")
		Активен := Entri.GetAttributeValue("rdbActive")
		ПарольПользователя := Entri.GetAttributeValue("userPassword")
		ДатаУвольнения := Entri.GetAttributeValue("userDateOrderDismiss")
		Описание := Entri.GetAttributeValue("description")

		var Должность string
		if FSSPconn == nil {
			Ошибка("Ошибка соединения с РДБ %+v\n")

		} else {
			Должность = ПолучитьДолжностьИзАИС(FSSPconn, Логин)
		}

		Заблокирован := false
		//Инфо(" Активен %+v", Активен)
		if Активен != "" {
			//Инфо(" Активен %+v", Активен)
			if Активен == "FALSE" {
				Заблокирован = true
			} else if Активен == "TRUE" {
				Заблокирован = false
			}
		} else {
			if ДатаУвольнения != "" {
				//Инфо(" ДатаУвольнения %+v", ДатаУвольнения)
				Заблокирован = true
			}
			if ПарольПользователя == "" {
				Заблокирован = true
			} else {

				if ЗаблокированС != "" && ЗаблокированПо != "" {
					ЗаблокированДо, err := time.Parse("02.01.2006", ЗаблокированПо)
					if err != nil {
						Ошибка("  %+v", err)
					}
					//Инфо(" Заблокирован До  %+v time.Now().Before(ЗаблокированДо) %+v", ЗаблокированДо, time.Now().Before(ЗаблокированДо))
					if time.Now().Before(ЗаблокированДо) {
						Заблокирован = true
					}
				}
				if strings.Contains(Описание, "rdbPassword") {
					//Инфо("Описание  %+v", Описание)
					Заблокирован = true
				}
			}
		}
		var ПериодБлокировки string
		Отпуск := false
		Декрет := false
		if ЗаблокированС != "" && ЗаблокированПо != "" {

			ЗаблокированДо, err := time.Parse("02.01.2006", ЗаблокированПо)
			ЗаблокированНчаинаяС, err := time.Parse("02.01.2006", ЗаблокированС)
			if err != nil {
				Ошибка("  %+v", err)
			}
			КоличествоДнейБлокироки := ЗаблокированДо.Sub(ЗаблокированНчаинаяС).Hours() / 24
			if КоличествоДнейБлокироки <= 45 {
				Отпуск = true
			} else {
				Декрет = true
			}
			Заблокирован = true
			//Инфо(" Отпуск %+v, Декрет %+v, Заблокирован  %+v", Отпуск, Декрет, Заблокирован)
			//
			//Инфо("ЗаблокированНчаинаяС %+v  Заблокирован До  %+v time.Now().Before(ЗаблокированДо) %+v На %+v",ЗаблокированНчаинаяС,  ЗаблокированДо, time.Now().Before(ЗаблокированДо), КоличествоДнейБлокироки)

			ПериодБлокировки = `{"С":"` + ЗаблокированС + `", "По":"` + ЗаблокированПо + `"}`
		}
		if ДатаУвольнения != "" {
			ПериодБлокировки = `{"Уволен":"` + ДатаУвольнения + `"}`
			Заблокирован = true
		}

		Значения := make([][]byte, 9)
		Значения = [][]byte{
			[]byte(ФИО),
			[]byte(ОСП),
			[]byte(Логин),
			[]byte(Должность),
			[]byte(strconv.FormatBool(Заблокирован)),
			[]byte(UidNumber),
			[]byte{},
			[]byte(strconv.FormatBool(Отпуск)),
			[]byte(strconv.FormatBool(Декрет)),
		}

		if ПериодБлокировки != "" {
			Значения[6] = []byte(ПериодБлокировки)
		} else {
			Значения[6] = nil
		}
		if Должность == "" {
			Значения[3] = nil
		}
		//Инфо("Логин %+v UidNumber  %+v %+v",Логин, UidNumber, UidNumber == "")
		if UidNumber == "" {
			Значения[5] = nil
		}

		//if Логин == "kormilceva@r26" {
		//	Инфо("Значения  %+s Заблокирован %+v", Значения, Заблокирован)
		//}
		//

		Обновление, err := sqlStruct{
			Name:     "Обнволение данных УЗ",
			Sql:      "INSERT INTO skzi.sczi (fio, osp, login, post, blocked, uid_number, delay, otpusk, decret, sync) VALUES ($1, $2, $3, $4,$5,$6,$7,$8,$9, now()) ON CONFLICT (login) DO UPDATE SET post = EXCLUDED.post, blocked = EXCLUDED.blocked, delay=EXCLUDED.delay , otpusk = EXCLUDED.otpusk, decret = EXCLUDED.decret, sync = now(), uid_number = EXCLUDED.uid_number returning *",
			Values:   Значения,
			DBSchema: "fssp_data",
		}.Выполнить(nil)
		if err != nil {
			Ошибка("  %+v", err)
		}
		Инфо("  %+v", Обновление)
	}

	Результат := map[string]interface{}{
		"ОтветКлиенту": "Учётные записи синхронихированны",
	}

	return Результат
}
func АктуализироватьДанныеПользователя(uid string) {
	GetLdapUsersPg(uid)
}
func GetLdapUsersPg(uid string) {
	ldapConn, err := LdapConnect()
	Инфо(" GetLdapUsersPg %+v", ldapConn, uid)
	if ldapConn == nil || err != nil {
		if err != nil {
			Ошибка(" Нет подключения к ldap %+v ", err)
		}
		return //"Ошибка соединения с сервером авторизации, попытка локальной авторизации"
	}
	defer ldapConn.Close()
	var filter string
	if uid != "" {
		filter = fmt.Sprintf("(&(objectClass=fsspData)(uid=%s))", uid)
	} else {
		filter = "(objectClass=fsspData)"
	}
	fields := []string{"uidNumber", "sn", "givenName", "initials", "ou", "uid", "userDateOrderDismiss", "description", "delayFrom", "delayTo"}
	searchResult, err := ldapConn.Search(ldap.NewSearchRequest(
		"dc=fssprus,dc=ru",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		fields,
		nil,
	))

	//(|(&(objectClass=fsspRole)(memberUid=maksimchuk@r26))(&(objectClass=fsspData)(uid=maksimchuk@r26)))
	// cn фио
	// sn - фамилия,
	// initials - отчество
	// ou подразхделение

	sqlString := "INSERT INTO users (uid, second_name, givenname, initials, ou, login, blocked, delay_from, delay_to, post) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT (login) DO UPDATE SET  blocked = $7, delay_from=$8, delay_to=$9, post=$10"

	sqlQuery := sqlStruct{
		Name:     "users",
		Sql:      sqlString,
		Values:   [][]byte{},
		DBSchema: "iobot",
	}

	PgConn, err := PGConnect("fssp_configs", nil)
	if err != nil {
		Ошибка("err	 %+v\n", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err = PgConn.Prepare(ctx, sqlQuery.Name, sqlQuery.Sql, nil)
	if err != nil {
		Ошибка("GetLdapUsersPg Prepare err %+v sqlData.Name %+v\n", err, sqlQuery.Name)
	}

	//resultat := PgConn.ExecPrepared(ctx, sqlData.Name, sqlData.Values, nil, nil)

	if len(searchResult.Entries) > 0 {

		FSSPconn := FSSPconnect()
		if FSSPconn == nil {
			Ошибка("Ошибка соединения с РДБ %+v\n")
			return
		}
		defer FSSPconn.Close()

		for _, Entri := range searchResult.Entries {

			//post_id, date_order := ПолучитьДолжность(Entri.GetAttributeValue("uid"))
			//if Entri.GetAttributeValue("userDateOrderDismiss") != "" {
			//	//Ошибка("Дата приказа об увольнении field userDateOrderDismiss  %+v\n", Entri.GetAttributeValue("userDateOrderDismiss"))
			//	continue
			//}
			userPostId := ПолучитьИдДолжностиИзАИС(FSSPconn, Entri.GetAttributeValue("uid"))
			if userPostId == -1 {
				Ошибка("Ошибка userPostId %+v\n", userPostId)
				return
			}
			//user_info := map[string]interface{}{
			//	"post_id":post_id,
			//	"date_order":date_order,
			//	"blocked":blocked,
			//	"start_block":start_block,
			//	"end_block":end_block,
			//}
			// проверим все даты блокироровок
			//if user_info["date_order"]=="" && user_info["start_block"]=="" && user_info["end_block"]=="" && user_info["blocked"]==1 && Entri.GetAttributeValue("description") !=""{
			//	//Ошибка("Уволен  %+v\n",Entri.GetAttributeValue("uid"))
			//	continue
			//}
			//if Entri.GetAttributeValue("ou") == "26911"{
			//	Ошибка(" %+v - description - %+v\n", Entri.GetAttributeValue("uid"),	Entri.GetAttributeValue("description"))
			//	Ошибка("user_info %+v\n", user_info)
			//}

			userDate := make([][]byte, 0, 6)
			for _, field := range fields {
				var inValue interface{}
				if field == "userDateOrderDismiss" {
					continue
				}
				if Entri.GetAttributeValue(field) != "" {
					if field == "description" {
						if Entri.GetAttributeValue("description")[0:9] != "statForms" {
							inValue = "true"
						} else {
							inValue = "false"
						}
					} else {
						inValue = Entri.GetAttributeValue(field)
					}
				} else {
					inValue = nil
				}

				if inValue != nil {
					userDate = append(userDate, []byte(inValue.(string)))
				} else {
					userDate = append(userDate, nil)
				}

			}

			//Ошибка("userDate %+v\n", userDate)

			if userPostId != -1 || userPostId != 0 {
				userDate = append(userDate, []byte(strconv.Itoa(userPostId)))
			} else {
				userDate = append(userDate, nil)
			}

			resultat := PgConn.ExecPrepared(ctx, sqlQuery.Name, userDate, nil, nil)
			if resultat.Read().Err != nil {
				Ошибка("resultat Err %+v\n", resultat.Read().Err)
			}

			//if resultat == nil {
			//	//Ошибка("userDate %+v, res %+v, err %+v\n",userDate, res, err)
			//	continue
			//}
			//err := PgConn.Close(ctx)
			//Ошибка("err %+v\n", err)
		}
	}
}

//func (client *Client) ПолучитьНaстройки(){
//	client.GetSetting()
//}
func (client *Client) СоздатьМессенджер(вопрос Сообщение) {
	Инфо("  %+v %+v %+v", "Создать мессенджер", client, client.UserInfo)

	//client.СозадтьИСохранитьТокен()
	//Инфо(" client.Token %+v", client.Token)
	data := map[string]interface{}{
		"client": client,
		//"usersList":client.ПолучитьКонтактЛист(),
		"UserInfo": client.UserInfo,
		"OUList":   client.ПолучитьСписокОтделов(),
		"Online":   server.Clients,
		"BotMenu":  client.ПолучитьМенюБота(),
		"Dialogs":  client.ПолучитьБыстрыеДиалоги(),
	}
	//if client.UserInfo.Curator ||  client.UserInfo.Info.OspCode== 26911{
	//	data["AdminMenu"]=client.ПолучитьМенюКуратора()
	//}
	//Ошибка("data.Dialogs %+v\n",data["Dialogs"] )
	//Ошибка("client.Token %+v\n", client.Token)
	//Инфо(" client %+v",  client.UserInfo.Info)

	СообщениеКлиенту := Сообщение{
		Token: struct {
			Hash     string
			Истекает string
		}{
			Hash:     client.Token.Hash,
			Истекает: client.Token.Истекает,
		},
		Ip:    client.Ip,
		Текст: "",
		От:    "io",
		Кому:  client.Login,
		Id:    вопрос.Id, // -3 html шаблон
		//MessageType:[]string{"InitMessenger"},
		Content: struct {
			Target     string      `json:"target"`
			Data       interface{} `json:"data"`
			Html       string      `json:"html"`
			Обработчик string      `json:"обработчик"`
		}{
			Target: "body",
			Html:   string(render("widget", data)),
		},
		UserInfo: struct {
			Uid        string `json:"uid"`
			Initials   string `json:"Initials"`
			FullName   string `json:"FullName"`
			FirstName  string `json:"FirstName"`
			LastName   string `json:"LastName"`
			MiddleName string `json:"MiddleName"`
			OspName    string `json:"OspName"`
			OspNum     int    `json:"OspNum"`
			PostName   string `json:"PostName"`
			Инициалы   string `json:"Инициалы"`
			ПолноеИмя  string `json:"ПолноеИмя"`
			Фамилия    string `json:"Фамилия"`
			Имя        string `json:"Имя"`
			Отчество   string `json:"Отчество"`
			ОСП        string `json:"ОСП"`
			КодОСП     int    `json:"КодОСП"`
			Должность  string `json:"Должность"`
		}{
			Uid:       client.Login,
			Initials:  data["UserInfo"].(*UserInfoStruct).Initials,
			FullName:  data["UserInfo"].(*UserInfoStruct).FullName,
			Инициалы:  data["UserInfo"].(*UserInfoStruct).Initials,
			ПолноеИмя: data["UserInfo"].(*UserInfoStruct).FullName,
			Имя:       data["UserInfo"].(*UserInfoStruct).Info.Givenname,
			Отчество:  data["UserInfo"].(*UserInfoStruct).Info.Initials,
			Фамилия:   data["UserInfo"].(*UserInfoStruct).Info.SecondName,
			ОСП:       data["UserInfo"].(*UserInfoStruct).Info.OspName,
			КодОСП:    data["UserInfo"].(*UserInfoStruct).Info.OspCode,
			Должность: data["UserInfo"].(*UserInfoStruct).Info.PostName,
		},
		AdminMenu: client.ПолучитьМенюКуратора(),
	}

	СообщениеКлиенту.СохранитьИОтправить(client)

	responseMes := Сообщение{
		Id:    0,
		Ip:    client.Ip,
		От:    "io",
		Кому:  client.Login,
		Текст: "приветствие",
	}
	responseMes.Выполнить.Arg.Login = "io"
	client.ПолучитьЛогПереписки(responseMes)

}
func (client *Client) ОпределитьПользователяПоIp(IP string) {
	Инфо("ОпределитьПользователяПоIp %+v", IP)

	Результат, err := sqlStruct{
		Name: "статичные_ip",
		Sql:  "SELECT uid FROM iobot.статичные_ip WHERE ip = $1",
		Values: [][]byte{
			[]byte(IP),
		},
	}.Выполнить(nil)
	if err != nil {
		Ошибка(">>>> Ошибка SQL запроса: %+v \n\n", err)
	}
	if len(Результат) > 0 {
		client.Login = Результат[0]["uid"].(string)
	}
}

func (client *Client) ПолучитьНaстройки() {
	sqlSting := "SELECT настройки FROM настройки_пользователя WHERE логин = $1"
	ПолучитьНайстройкиПользователя, err := sqlStruct{
		Name: "user_setting",
		Sql:  sqlSting,
		Values: [][]byte{
			[]byte(client.Login),
		},
	}.Выполнить(nil)
	if err != nil {
		//Ошибка(" ПолучитьНaстройки  \n %+v \n", err)
	}
	if len(ПолучитьНайстройкиПользователя) > 0 {
		client.НайстройкиПользователя = ПолучитьНайстройкиПользователя[0]["настройки"].(map[string]interface{})
	}

	//CaptionRow, _ := ВыполнитьPgSQL(queryCaption)
	////Ошибка("ПолучитьНaстройки CaptionRow %+v\n", CaptionRow)
	//if CaptionRow != nil{
	//	var setting map[string]interface{}
	//	err := json.Unmarshal([]byte(CaptionRow[0]["setting"].(string)), &setting)
	//	if err != nil {
	//		Ошибка("err	 %+v\n", err)
	//	}
	//	client.Setting = setting
	//}
}
