package main

import (
	"bytes"
	"context"
	_ "encoding/json"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	_ "github.com/nakagami/firebirdsql"
	"io"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Аргументы struct {
	Название string `json:"название"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var server = Server{
	Messages:  []*Сообщение{},
	Clients:   map[string]*Client{},
	addCh:     nil,
	delCh:     nil,
	sendAllCh: nil,
	doneCh:    nil,
	errCh:     nil,
}

type Server struct {

	//pattern   string
	Messages  []*Сообщение
	mx        sync.RWMutex
	Clients   map[string]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
}

//var logFile *os.File
var РабочаяПапка string
var IoLoger *IoLog

func init() {
	Инфо("init ИО %+v", time.Now())
	var err error
	РабочаяПапка, err = filepath.Abs(filepath.Dir("" + os.Args[0]))
	Инфо(" РабочаяПапка %+v", РабочаяПапка)
	if err != nil {
		Ошибка(" Не удалось определить рабочую папку %+v", err)
	}
	//log.SetFlags(log.Llongfile)
	//log.SetFlags(log.Ltime|log.Lshortfile)
	//ЛогерВсехЗапросов()
	ActionsInit()
	Инфо("Инициализация %+v", time.Now())

}

func main() {

	http.HandleFunc("/debug/", pprof.Index) // http://10.26.6.25:8081/debug/pprof/

	http.HandleFunc("/static/", StaticHandler)
	http.HandleFunc("/uploads/", StaticHandler)
	//flag.Parse()

	// рендерим рабочий стол
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Инфо("*************Входящий запрос INDEX ***********\n", r)
		index(w, r)
	})

	http.HandleFunc("/wsconnect", func(w http.ResponseWriter, r *http.Request) {
		Инфо("*************Входящий запрос wsconnect *********** RemoteAddr: %+v; \n", r.RemoteAddr)
		connected(w, r)
	})

	serverError := http.ListenAndServe(":8080", nil)

	Инфо("Запуск ИО %+v", time.Now())

	if serverError != nil {
		Ошибка(">>>> serverError ERROR \n %+v \n\n", serverError)
	}
}

func ЛогерВсехЗапросов() {
	IoLoger = &IoLog{}
	/*
		Логируем все rintf в базу
	*/
	IoLoger.DB, _ = PGConnect("logs", nil)
	IoLoger.QueryName = "main_log"
	IoLoger.ctx, IoLoger.cancel = context.WithCancel(context.Background())
	_, err := IoLoger.DB.Prepare(IoLoger.ctx, IoLoger.QueryName, "INSERT INTO wsserver_log (time,log) VALUES ($1,$2)", nil)
	if err != nil {
		Ошибка("IoLoger.DB.Prepare err %+v IoLoger.QueryName %+v\n", err, IoLoger.QueryName)
	}

	//mw := io.MultiWriter(os.Stdout, IoLoger)
	//log.SetOutput(mw)
}
func (server *Server) УведомитьВсехОПодключении(clientOnline *Client, канал chan string) {
	for sendLogin, client := range server.Clients {

		Инфо("sendLogin %+v \n", sendLogin)

		if sendLogin != clientOnline.Login {
			responseMes := Сообщение{
				Id:     0,
				От:     "io",
				Кому:   sendLogin,
				Online: clientOnline.Login,
				Content: struct {
					Target     string      `json:"target"`
					Data       interface{} `json:"data"`
					Html       string      `json:"html"`
					Обработчик string      `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
				}{
					Data: clientOnline.UserInfo,
				},
			}
			client.Message <- &responseMes
		}
	}
	канал <- "конец УведомитьВсехОПодключении " + clientOnline.Login

}

func (server *Server) УведомитьВсехОбОтключении(code int, text string, login string, канал chan string) {

	for sendLogin, client := range server.Clients {
		//Инфо("sendLogin  %+v", sendLogin)
		responseMes := Сообщение{
			Id:      0,
			От:      "io",
			Кому:    sendLogin,
			Offline: login,
		}
		client.Message <- &responseMes
	}
	//rintf("УведомитьВсехОбОтключении code %+v text  %+v\n", code, text,)
	//rintf("УведомитьВсехОбОтключении server %+v\n", server.Clients)

	канал <- "конец УведомитьВсехОбОтключении " + login

}

func index(writer http.ResponseWriter, request *http.Request) {
	Инфо("index  %+v", time.Now())

	Token, err := request.Cookie("Token")
	//v, e := url.ParseQuery(request.URL.RawQuery)
	//if e != nil {
	//	Ошибка("  %+v", e)
	//}
	//Инфо("request  %+v", request.Cookies() )

	if err != nil {
		Ошибка("err Token не найден, рендерим форму авторизации %+v\n", err)
		//tplName = "auth"
	}

	tplName := "index"
	var data interface{}
	var IP string
	Инфо("request.RemoteAddr %+v", request.RemoteAddr)

	if strings.Contains(request.RemoteAddr, "[::1]") {
		IP = "10.26.6.13"
	} else {
		IP = strings.Split(request.RemoteAddr, ":")[0]
	}
	Инфо("IP %+v", IP)
	if Token != nil {
		//rintf(" Token %+v\n", Token.Value)
		Инфо("  %+v Token.Value %+v", "ПолучитьДанныеСессии", Token.Value)
		ДанныеСессии, err := sqlStruct{
			Name: "user_session",
			Sql:  "select * from iobot.user_session where token = $1",
			Values: [][]byte{
				[]byte(Token.Value),
			},
		}.Выполнить(nil)

		if err != nil {
			Ошибка(">>>> ERROR \n %+v \n\n", err)
		}
		if len(ДанныеСессии) > 0 {
			ДанныеСессии := ДанныеСессии[0]
			Инфо("ДанныеСессии %+v\n", ДанныеСессии)

			data = map[string]interface{}{
				"ContentData": map[string]interface{}{
					"tplName": "dashboard",
					"tplData": map[string]string{"login": ДанныеСессии["uid"].(string)},
				},
			}

		} else {
			Инфо("Данные авторизации не найдены в базе данных, токен не найден %+v\n", Token)

			//ДанныеСессии, err := sqlStruct{
			//	Name:   "user_session",
			//	Sql:    "select * from iobot.user_session where token = $1",
			//	Values: [][]byte{
			//		[]byte(Token.Value),
			//	},
			//}.Выполнить(nil)

			data = map[string]interface{}{
				"ContentData": map[string]interface{}{
					"tplName": "auth",
					"tplData": ОпределитьПользователяПоIp(IP),
				},
			}
			//ContentHtml = render("auth", nil)
		}
	} else {
		data = map[string]interface{}{
			"ContentData": map[string]interface{}{
				"tplName": "auth",
				"tplData": ОпределитьПользователяПоIp(IP),
			},
		}
		//ContentHtml = render("auth", nil)
	}
	Инфо(" tplName %+v, data %+v", tplName, data)

	_, errWrite := writer.Write(render(tplName, data))
	//Инфо("КоличествоБайт  %+v", КоличествоБайт)
	if errWrite != nil {
		Ошибка("render:", errWrite)
	}
}

func connected(w http.ResponseWriter, r *http.Request) {

	Инфо("  %+v", r.RemoteAddr)

	канал := make(chan string, 100)
	//if strings.Split(r.RemoteAddr,":")[0] !="10.26.6.13" {
	//	return
	//}
	go wsConnector(w, r, канал)

	//Инфо("поток   %+v", "wsConnector")

	for {
		result := <-канал
		Инфо("connected result: %+s \n", result)
	}

}

func wsConnector(w http.ResponseWriter, r *http.Request, канал chan string) {
	Инфо("wsConnector  %+v", r.RemoteAddr)
	//Инфо(" r.Header.Get %+s \n", r.Header.Get("Origin"))
	queryArgs, _ := url.ParseQuery(r.URL.RawQuery)

	var ЛогинСПортала string

	if _, ЕстьЛогинССайта := queryArgs["login"]; ЕстьЛогинССайта {
		ЛогинСПортала = queryArgs["login"][0]
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		Ошибка("upgrade:", err)
		return
	}
	var IP string
	if strings.Contains(r.RemoteAddr, "[::1]") {
		IP = "10.26.6.13"
	} else {
		IP = strings.Split(r.RemoteAddr, ":")[0]
	}

	client := &Client{
		Host:    r.Host,
		Ip:      IP,
		Ws:      conn,
		Message: make(chan *Сообщение),
	}
	Инфо(" client  %+v ЛогинСПортала %+v", client, ЛогинСПортала)

	if client.Login == "" {

		if ЛогинСПортала == "auth" || ЛогинСПортала == "" {
			client.Login = IP
		} else {
			client.Login = ЛогинСПортала
			client.ПолучитьДанныеПользователя("")
			client.ПолучитьНaстройки()
		}
	} else {
		Инфо("client.Login %+v", client.Login)
		client.ПолучитьДанныеПользователя("")
		client.ПолучитьНaстройки()
	}

	client.Ws.SetCloseHandler(func(code int, text string) error {

		if server.Clients[client.Login] != nil && server.Clients[client.Login].Ssh == nil {
			// тут возможно ислкючение когда струткура Ssh ещё пустая но ведёться попытка  оединиться с клиентом, тогда может выпасть ошибка
			// note: возможно нужно сделать задержку перед удалением из памяти, и добавить ключ онлайн оффлан, чтобы не писать/удалять постоянно клитента если он просто перешёл на новую страницу или обновил страницу...
			server.mx.Lock()
			delete(server.Clients, client.Login)
			server.mx.Unlock()
		}
		go server.УведомитьВсехОбОтключении(code, text, client.Login, канал)
		return nil
	})

	if server.Clients[client.Login] != nil {
		Ошибка("В Карте сервера уже есть запись с логином %+v; %+v", client.Login, server.Clients[client.Login])
		Инфо("старый RemoteAddr()  %+v", server.Clients[client.Login].Ws.RemoteAddr())
		Инфо("  %+v", client.Ws.RemoteAddr())
	}

	Инфо("RemoteAddr()  %+v", client.Ws.RemoteAddr())

	server.mx.Lock()
	server.Clients[client.Login] = client
	server.mx.Unlock()

	go client.SendMessage()
	go client.ReadMessage()

	Token, errToken := r.Cookie("Token")
	Инфо("r.Cookies()  %+v", r.Cookies())
	if errToken != nil {
		Ошибка("Token не найден %+v\n", err)
	}

	//if queryArgs["reconect"] != nil && queryArgs["reconect"][0]=="true"{
	//
	//} else {
	if ЛогинСПортала != "auth" && Token != nil {
		Инфо("ЛогинСПортала  %+v", ЛогинСПортала)
		Инфо("Token  %+v", Token)
		client.Token = struct {
			Hash     string
			Истекает string
		}{
			Hash:     Token.Value,
			Истекает: Token.Expires.String(),
		}
		client.СозадтьИСохранитьТокен()
		go client.СоздатьМессенджер(Сообщение{}) //server.Clients[userName]

		OriginS := strings.Split(r.Header.Get("Origin"), ":")
		Origin := strings.Replace(OriginS[1], "//", "", -1)

		if Origin == "10.26.6.25" {
			go client.СоздатьМессенджер(Сообщение{})
			client.ПолучитьМеню()
			client.СоздатьРабочийСтол(Сообщение{})
		}

	} else if Token == nil && (ЛогинСПортала != "auth" && ЛогинСПортала != "") && r.Header.Get("Origin") == "http://10.26.4.20" {

		Инфо(" Token %+v", Token)
		Инфо(" ЛогинСПортала %+v", ЛогинСПортала)
		Инфо(" r.Header.Get(Origin) %+v", r.Header.Get("Origin"))
		//client.СозадтьИСохранитьТокен()
		go client.СоздатьМессенджер(Сообщение{})
	}

	Инфо(" r.Header.Get(Origin) %+v r.Header %+v", r.Header.Get("Origin"), r.Header)

	//}

	if ЛогинСПортала != "auth" && Token != nil {
		go server.УведомитьВсехОПодключении(client, канал)
	}

	if client.НайстройкиПользователя != nil {
		Инфо("client.НайстройкиПользователя %+v", client.НайстройкиПользователя) // client.Setting
	}

	//Инфо(" %+v", "Закончили создавать WS подключение")  {"action": "getChatLog", "arg": {"Login": login}}

	//if r.Header.Get("Origin") == "http://10.26.6.25" || r.Header.Get("Origin") == "http://10.26.6.25:8081" {
	////if client.Host == "localhost" || client.Host == "10.26.6.25" || client.Host == "10.26.6.25:8081" || client.Host == "10.26.6.30"{
	//	Инфо("   %+v client.Host %+v", "ПолучитьМеню", client.Host)
	//	client.ПолучитьМеню()
	//}
	канал <- "Закончили создавать WS подключение"
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	var static_file string
	//Инфо("StaticHandler  %+v", StaticHandler)
	//rintf("РабочаяПапка %+v\n", РабочаяПапка)

	//rintf("http.Dir(static) %+v\n", http.Dir(РабочаяПапка+"/static/"))

	//static_file = req.URL.Path[len("/static/"):]
	static_file = req.URL.Path

	if len(static_file) != 0 {

		if static_file == "/static/js/tpl.js" {
			//rintf("static_file %+v\n", static_file)
			jsFile := renderJS("tplsJs", nil)
			fileBytes := bytes.NewReader(jsFile)
			content := io.ReadSeeker(fileBytes)

			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}

		f, err := http.Dir(РабочаяПапка).Open(static_file)

		//Инфо("f  %+v", f)

		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		} else {
			Ошибка("%+v\n", err)
		}
	}
	http.NotFound(w, req)

	//if strings.Contains(req.URL.Path, "/uploads/") {
	//
	//}

}
