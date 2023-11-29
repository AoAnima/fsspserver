package main

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)
type AuthDate struct {
	Login string
	Password string
}

func (client *Client) ЛокальнаяАвторизация (ДанныеАвторизации AuthDate){

	UIDПользователя ,ОшибкаЗапроса := sqlStruct{
			Name:   `Авторизация`,
			Sql:    `SELECT uid FROM fssp_configs.users WHERE login = $1 AND pwd = $2`,
			Values: [][]byte{
				[]byte(ДанныеАвторизации.Login),
				[]byte(ДанныеАвторизации.Password),
			},
			DBSchema:`fssp_configs`,
		}.Выполнить(nil)
	if ОшибкаЗапроса != nil{
	 Ошибка(">>>> Ошибка SQL запроса: %+v \n\n", ОшибкаЗапроса)
	} else {
		if len(UIDПользователя) >0 {
			UID := UIDПользователя[0]
			if UID != nil{
				client.Login = ДанныеАвторизации.Login
				 client.ПолучитьДанныеПользователя("")
			}
		}
	}
}

func (client *Client) Авторизация (вопрос Сообщение){

	Аргументы := вопрос.Выполнить.Действие["Авторизация"]

	Инфо("Авторизация Аргументы %+v\n", Аргументы)

	Логин, ЕстьЛогин := Аргументы["login"]
	Пароль, ЕстьПароль := Аргументы["password"]




	ДанныеАвторизации := AuthDate{}
	if (ЕстьЛогин && Логин != "") && (ЕстьПароль && Пароль != ""){
		if !strings.Contains(Логин.(string), "@r26") {
			Логин = Логин.(string)+"@r26"
		}
		ДанныеАвторизации = AuthDate{
			Login:   Логин.(string),
			Password:  Пароль.(string),
		}
		

	} else {
		Инфо("Нет логина или пароля: %+v\n", Аргументы)
		return
	}
	_ ,ОшибкаЗапроса := sqlStruct{
		Name:   `Авторизация`,
		Sql:    `UPDATE fssp_configs.users SET pwd = $2 where login = $1`,
		Values: [][]byte{
			[]byte(ДанныеАвторизации.Login),
			[]byte(ДанныеАвторизации.Password),
		},
		DBSchema:`fssp_configs`,
	}.Выполнить(nil)

	if ОшибкаЗапроса != nil {
		Ошибка(" %+v ", ОшибкаЗапроса)
	}
	
	
	LdapConn, LdapErr := LdapConnect()
	//var LdapConn *ldap.Conn

	if LdapConn == nil || LdapErr != nil { //|| err != nil
		Инфо("Ldap не отвечает, пробуем авторизироватся локально %+v", ДанныеАвторизации)

		client.ЛокальнаяАвторизация (ДанныеАвторизации)

		if client.UserInfo != nil{
			Инфо("  %+v", "СозадтьИСохранитьТокен")
			client.СозадтьИСохранитьТокен()
			client.СоздатьМессенджер(вопрос)
		} else {

			Ошибка("НЕ удаётся пройти авторизацию локально. ДанныеАвторизации %+v client %+v", ДанныеАвторизации, client)

			СообщениеКлиенту:= &Сообщение{
				Текст:   "Сервер авторизации не отвечает. Пройти авторизацию локально не удаётся.",
				От: "io",
				Кому:client.Login,
				MessageType: []string{"float_message","error"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
		}
		return
	} else {

	defer LdapConn.Close()
	Инфо("  %+v", ДанныеАвторизации)
	// Подключаемся под пользователем, если выдал ошибку то значит не праивльный логин или пароль
	errBind := LdapConn.Bind("uid="+ДанныеАвторизации.Login+",ou=r26,ou=users,dc=fssprus,dc=ru", ДанныеАвторизации.Password)

	if errBind != nil {
		//LDAP Result Code 49 "Invalid Credentials": не верный пароль
		Ошибка("%+v\n", errBind)
		СообщениеКлиенту:= &Сообщение{
			Текст:   "Лоин и/или пароль не свопадают",
			От: "io",
			Кому:client.Login,
			MessageType: []string{"float_message","error"},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)


	} else {
		//ldap авторизация прошла успешно
		_, ЛогинНайден := server.Clients[ДанныеАвторизации.Login]

		if ЛогинНайден {

			Инфо(" >>>>> \n Странно логина не должно быть. или был вход с портала %+v ДанныеАвторизации.Login %+v <<<<< \n\n", server.Clients[ДанныеАвторизации.Login], ДанныеАвторизации.Login)
			client.Login = ДанныеАвторизации.Login
			client.СозадтьИСохранитьТокен()

		} else {
			//заменим данные WS подключения в карте сервера

			_, IpНайден := server.Clients[client.Ip]
			client.Login = ДанныеАвторизации.Login
			if IpНайден{


				Инфо("\n >>>>>> Авторизация server.Clients[client.Ip] %+v client.Ip %+v\n", server.Clients[client.Ip], client.Ip)
				Инфо("\n server.Clients[ДанныеАвторизации.Login] %+v \n", server.Clients[ДанныеАвторизации.Login])

				//client.СозадтьИСохранитьТокен()

				if server.Clients[ДанныеАвторизации.Login] != nil {
					Инфо("В Карте сервера есть запись с логином ДанныеАвторизации.Login  %+v  %+v", ДанныеАвторизации.Login, server.Clients[ДанныеАвторизации.Login])
				}
				server.mx.Lock()
				server.Clients[ДанныеАвторизации.Login] = client
				server.mx.Unlock()

				delete(server.Clients, client.Ip)

				//log.Printf("\n server.Clients[client.Ip] %+v\n", server.Clients[client.Ip])
				//log.Printf("\n server.Clients[ДанныеАвторизации.Login] %+v\n", server.Clients[ДанныеАвторизации.Login])

				Инфо("\n >>>>>> Авторизация client %+v\n", client)
				// создадим, токен, и сгенерируем время действия токена +5 дней с текущей даты
			}
		}


		АктуализироватьДанныеПользователя(ДанныеАвторизации.Login)
		Инфо(" Авторизация получаем данные пользователя  %+v", ДанныеАвторизации.Login)
		//client.ПолучитьДанныеПользователя(ДанныеАвторизации.Login)

		client.ПолучитьДанныеПользователя("")
		client.ПолучитьНaстройки()

		client.СозадтьИСохранитьТокен()
		client.СоздатьМессенджер(вопрос)
		client.ПолучитьМеню()
		Инфо("  %+v", "СоздатьРабочийСтол")
		client.СоздатьРабочийСтол(вопрос)
	}
	}

}

func (client *Client)СозадтьИСохранитьТокен() {
	Инфо("  %+v", "СозадтьИСохранитьТокен")
	время := time.Now()
	времяДействия := время.Add(time.Hour * 24 *5)

	utcTime := времяДействия.Format(time.RFC1123)
	символы := client.Login+client.Ip+strconv.FormatInt(времяДействия.Unix(), 10)
	Инфо(" символы %+v", символы)
	token := make([]byte, 30)

	for i := 0; i < 30; i++ {
		token[i] = символы[rand.Intn(len(символы))]
	}
	Инфо(" token %+s", token)
	client.Token = struct {
		Hash     string
		Истекает string
	}{
		Hash: string(token),
		Истекает: utcTime,
	}
	//Инфо(" СозадтьИСохранитьТокен client.Token %+v\n", client.Token)
	_ ,err:= sqlStruct{
		Name:   "user_session",
		Sql:    "INSERT INTO iobot.user_session (uid, token, ip, date_auth) VALUES ($1,$2,$3, now()) ON CONFLICT (uid) DO UPDATE SET token = $2, date_auth = NOW()",
		Values: [][]byte{
			[]byte(client.Login),
			[]byte(client.Token.Hash),
			[]byte(client.Ip),
		},
	}.Выполнить(nil)

	if err != nil{
		Ошибка(">>>> Ошибка SQL запроса: %+v \n\n",err)
	}
}
