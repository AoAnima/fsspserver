package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"html/template"
	"log"
	"reflect"
	"strconv"
	"strings"
	textTpl "text/template"
	"time"
)

type ДанныеОтвета struct {
	СпособВставки string      `json:"способ_вставки"`
	Контейнер     string      `json:"контейнер"`
	Данные        interface{} `json:"данные"`
	HTML          string      `json:"html"`
	Обработчик    string      `json:"обработчик"` //JS функция или объект/класс/плагин для обработки данных (table..)
}
type Сообщение struct {
	F5    bool `json:"f5"`
	Token struct {
		Hash     string
		Истекает string
	} `json:"Token"`
	Id                int                      `json:"Id"`
	Ip                string                   `json:"ip"`
	От                string                   `json:"От"`
	Кому              string                   `json:"Кому"`
	Текст             string                   `json:"Текст"`
	MessageType       []string                 `json:"MessageType"`
	Время             string                   `json:"Время"`
	ОтветНа           string                   `json:"ОтветНа"`
	Файлы             []string                 `json:"Файлы"`
	Offline           string                   `json:"Offline"`
	Online            string                   `json:"Online"`
	AdminMenu         []map[string]interface{} `json:"AdminMenu"` // []BotMenuStruct
	ВходящиеАргументы map[string]interface{}   `json:"ВходящиеАргументы"`

	ОбратныйВызов string `json:"ОбратныйВызов"`
	Выполнить     struct {
		Action   string                            `json:"action"`
		Действие map[string]map[string]interface{} `json:"Действие"` // "НазваниеДействия" :{"имяАргумента_1":"значение аргумента_1"... }
		Skill    int                               `json:"skill"`
		Навык    json.Number                       `json:"Навык"`
		Cmd      string                            `json:"cmd"`
		Комманду string                            `json:"Комманду"`
		Arg      struct {
			Module string                 `json:"module"`
			Tables []string               `json:"tables"`
			Login  string                 `json:"login"`
			Other  map[string]interface{} `json:"other"`
		} `json:"Arg"`
	} `json:"Выполнить"`
	Контэнт *ДанныеОтвета `json:"Контэнт"`
	Content struct {
		Target     string      `json:"target"`
		Data       interface{} `json:"data"`
		Html       string      `json:"html"`
		Обработчик string      `json:"обработчик"` //JS функция или объект/класс/плагин для обработки данных (table..)
	} `json:"Content"`
	UserInfo struct {
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
	} `json:"UserInfo"`
}

type Message struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	Text     string   `json:"text"`
	Time     string   `json:"time"`
	ReaplyTo string   `json:"reaply_to"`
	Files    []string `json:"files"`
}

type messageRowSql struct {
	Id                  int            `json:"id"`
	Autor               string         `json:"autor"`
	Recipient           string         `json:"recipient"`
	ChatRoom            sql.NullString `json:"chat_room"`
	ReaplyTo            sql.NullString `json:"reaply_to"`
	Date                string         `json:"date"`
	Files               sql.NullString `json:"files"`
	Text                sql.NullString `json:"text"`
	TextHtml            template.HTML  `json:"text_html"`
	AutorName           string         `json:"autor_name"`
	AutorMiddlename     string         `json:"autor_middlename"`
	RecipientName       string         `json:"recipient_name"`
	RecipientMiddlename string         `json:"recipient_middlename"`
}

type messageRow struct {
	Id                  int           `json:"id, mes_order"`
	Autor               string        `json:"autor"`
	Recipient           string        `json:"recipient"`
	ChatRoom            string        `json:"chat_room"`
	ReaplyTo            string        `json:"reaply_to"`
	Date                string        `json:"mes_date"`
	Files               []string      `json:"files"`
	Text                string        `json:"text"`
	TextHtml            template.HTML `json:"text_html"`
	AutorName           string        `json:"autor_name"`
	AutorMiddlename     string        `json:"autor_middlename"`
	RecipientName       string        `json:"recipient_name"`
	RecipientMiddlename string        `json:"recipient_middlename"`
}

func (mes Сообщение) ИзменитьСообщение() int {

	return 0
}

func (mes *Сообщение) СохранитьЛогСообщения() {
	//log.Printf("СохранитьСообщение mes %+v\n", mes)

	columns := ""
	countColumns := 3
	//if mes.Время == ""{
	ТекущеееВремя := time.Now()
	mes.Время = ТекущеееВремя.Format("2006-01-02T15:04:05.999999")
	//}
	sqlArgStr := []string{
		mes.От,
		mes.Кому,
		mes.Время,
	}
	sqlArgs := [][]byte{
		[]byte(mes.От),
		[]byte(mes.Кому),
		[]byte(mes.Время),
	}
	//log.Printf("\n mes.Время %+v\n",mes.Время )
	if mes.ОтветНа != "" {
		columns = columns + ", reaply_to"
		countColumns++
		sqlArgs = append(sqlArgs, []byte(mes.ОтветНа))
	}

	if mes.Файлы != nil {
		columns = columns + ", files"
		countColumns++

		FilesString, err := json.Marshal(mes.Файлы)
		if err != nil {
			Ошибка("err	 %+v\n", err)
		}
		sqlArgs = append(sqlArgs, FilesString)
	}

	columns = columns + ", text"
	countColumns++
	if mes.Текст != "" {
		sqlArgs = append(sqlArgs, []byte(mes.Текст))
		sqlArgStr = append(sqlArgStr, mes.Текст)
	} else {
		sqlArgs = append(sqlArgs, []byte(nil))
		sqlArgStr = append(sqlArgStr, mes.Текст)
	}

	if mes.Выполнить.Action != "" || mes.Выполнить.Cmd != "" || mes.Выполнить.Skill != 0 || Contains(mes.MessageType, "io_action") {
		columns = columns + ", type"
		countColumns++
		sqlArgs = append(sqlArgs, []byte(`["io_action"]`))
		sqlArgStr = append(sqlArgStr, "io_action")

		//
		КомандаБоту := map[string]string{}
		if mes.Выполнить.Action != "" {
			КомандаБоту["Action"] = mes.Выполнить.Action
		}
		if mes.Выполнить.Cmd != "" {
			КомандаБоту["Cmd"] = mes.Выполнить.Cmd
		}
		if mes.Выполнить.Skill != 0 {
			КомандаБоту["Skill"] = strconv.Itoa(mes.Выполнить.Skill)
		}
		byteString, err := json.Marshal(КомандаБоту)
		if err != nil {
			Ошибка("err	 %+v\n", err)
		}
		if len(КомандаБоту) > 0 {
			columns = columns + ", comand_to_io"
			countColumns++
			sqlArgStr = append(sqlArgStr, string(byteString))
			sqlArgs = append(sqlArgs, byteString)
		}

	}
	//if mes.Id != 0{
	columns = columns + ", входящий_номер"
	countColumns++
	sqlArgs = append(sqlArgs, []byte(strconv.Itoa(mes.Id)))
	//}

	valuesPlaceholder := ""
	if countColumns > 3 {
		for i := 4; i <= countColumns; i++ {
			valuesPlaceholder = valuesPlaceholder + ", $" + strconv.Itoa(i)
		}
	}

	//log.Printf("columns %+v valuesPlaceholder %+v\n",columns,  valuesPlaceholder)

	sqlString := `INSERT INTO messages (autor, recipient, mes_date ` + columns + `) VALUES ($1,$2, $3 ` + valuesPlaceholder + `) RETURNING message_id`
	// AND date >= CURRENT_DATE - INTERVAL '1 day'

	sqlQuery := sqlStruct{
		Name:     "messages",
		Sql:      sqlString,
		Values:   sqlArgs,
		DBSchema: "iobot",
	}
	//log.Printf("\n\nsqlArgStr >> %+v\n \n", sqlArgStr)

	//messagesArrayMap, _ := ВыполнитьPgSQL(sqlQuery)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Result, err := sqlQuery.PgSqlResultReader(ctx)
	//log.Printf("\nResultReader %+v\n", ResultReader)

	if err != nil {
		Ошибка("\n\n !!! Ошибка сохранения сообщения %+v\n", err)
		return
	}
	rows := Result.Rows

	var id int
	for _, row := range rows {
		for _, cell := range row {
			id, _ = strconv.Atoi(string(cell))
		}
	}

	//Инфо("mes.Id %+v mes.Id == 0 %+v", mes.Id, mes.Id == 0)

	//Инфо("Id %+v", id)
	if mes.Id == 0 {
		mes.Id = id
	}
	//Инфо("mes.Id %+v Id %+v", mes.Id,id)

}

func (client *Client) ПолучитьЛогТерминала(mes Сообщение) {

}

func (client *Client) ПолучитьЛогПереписки(mes Сообщение) {

	//@language=PostgresSQL
	sqlString := `select t.message_id, row_to_json(t.*) from (
SELECT EXTRACT(EPOCH FROM mes_date) as mes_order , mes_date, iobot.messages.*, user_autor.givenname autor_name, user_autor.initials autor_middlename,
      user_recipient.givenname recipient_name,user_recipient.initials recipient_middlename
FROM iobot.messages
        LEFT JOIN fssp_configs.users user_autor on user_autor.Login = autor
        LEFT JOIN fssp_configs.users user_recipient on user_recipient.Login =recipient
WHERE ((text IS NOT NULL AND text<>'') OR  files IS NOT NULL) AND ((autor = $1 AND recipient = $2) OR (recipient = $1 AND autor = $2)) ORDER BY message_id DESC LIMIT 30) t group by t.message_id, t.*`

	// пара автор, получатель
	sqlQuery := sqlStruct{
		Name: "Лог переписки",
		Sql:  sqlString,
		Values: [][]byte{
			[]byte(client.Login),
			[]byte(mes.Выполнить.Arg.Login),
		},
		DBSchema: "iobot",
	}

	//messagesArrayMap, _ := ВыполнитьPgSQL(sqlQuery)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Result, err := sqlQuery.PgSqlResultReader(ctx)
	if err != nil {
		Ошибка("\n !! ERR %+v\n", err)
	}
	MessagesLog := map[int]messageRow{}

	messagesRows := Result.Rows

	for _, messageByte := range messagesRows {
		//log.Printf("string(messageByte[0] id%+v\n", id)
		//log.Printf("string(messageByte[0] %+v\n", string(messageByte[0]))
		//log.Printf("string(messageByte[1] %+v\n", string(messageByte[1]))
		messageStruct := messageRow{}

		//message.TextHtml = template.HTML(message.Text.String)
		err := json.Unmarshal(messageByte[1], &messageStruct)

		if err != nil {
			Ошибка("Unmarshal messageByte[1] err  %+v\n", err, string(messageByte[1]))
		}

		mesId, err := strconv.Atoi(string(messageByte[0]))
		if err != nil {
			Ошибка("err	 %+v\n", err)
		}

		//text := strings.Replace(messageStruct.Text, "\n", "<br>", -1)
		messageStruct.TextHtml = template.HTML(messageStruct.Text)
		MessagesLog[mesId] = messageStruct
	}
	//UserInfo := &UsersStruct{}
	//if mes.Выполнить.Arg.Login != "io"{
	UserInfo := client.ПолучитьДанныеПользователя(mes.Выполнить.Arg.Login)
	//}

	data := map[string]interface{}{
		"client":   client,
		"uid":      client.Login,
		"log":      MessagesLog,
		"UserInfo": UserInfo,
		"BotMenu":  client.ПолучитьМенюБота(),
		"Dialogs":  client.ПолучитьБыстрыеДиалоги(),
	}

	var Data map[string]string
	if server.Clients[mes.Выполнить.Arg.Login] != nil {
		Data = map[string]string{"ClientIp": server.Clients[mes.Выполнить.Arg.Login].Ip}
	}

	responseMes := Сообщение{
		Id:          0,
		Ip:          client.Ip,
		От:          "io",
		Кому:        client.Login,
		MessageType: []string{"log"},
		Content: struct {
			Target     string      `json:"target"`
			Data       interface{} `json:"data"`
			Html       string      `json:"html"`
			Обработчик string      `json:"обработчик"`
		}{
			Target: "log_wrapper_" + mes.Выполнить.Arg.Login,
			Html:   string(render("messageLog", data)),
			Data:   Data,
		},
	}

	client.Message <- &responseMes

	//Инфо("\n\n !!! >>>>. ПолучитьСообщениеИО mes.Выполнить.Arg.Login %+v\n", mes.Выполнить.Arg.Login)

	if mes.Выполнить.Arg.Login == "io" {

		ТекстСообщения := client.НайтиОтветНаСообщение(&mes)

		Инфо(" mes.Выполнить.Arg.Login== io ТекстСообщения %+v", ТекстСообщения)

		ЕстьОтвет := ТекстСообщения["ответ"]
		if ЕстьОтвет != nil {
			утверждение := ЕстьОтвет.(map[string]interface{})["утверждение"]
			if утверждение != nil {

			}
		}
		//СообщениеКлиенту:= &Сообщение{
		//			Текст:   ТекстСообщения,
		//			От: "io",
		//			Кому:client.Login,
		//		}
		//		//СообщениеКлиенту.СохранитьЛогСообщения()
		//client.Message<-СообщениеКлиенту

		//ТекстСообщения2 := client.НайтиОтветНаСообщение(&mes)
		//
		//Ошибка("НУЖНО ПРОВЕРИТЬ ТекстСообщения2 %+v", ТекстСообщения2)

		//Инфо("client.Host  %+v", client.Host)
		//		if client.Host == "localhost" || client.Host == "10.26.6.25" || client.Host == "10.26.6.25:8081" || client.Host == "10.26.6.30"{
		//			Инфо(" перед СоздатьРабочийСтол %+v", mes)
		//			client.СоздатьРабочийСтол(mes)
		//		}
	}
}

// алгоритм Получить уровень доступа left join fssp_configs.уровень_доступа ON fssp_configs.уровень_доступа.уровень ?| (
//    select array_agg(level.уровень) from (
//           select jsonb_array_elements_text(уровень) уровень  from fssp_configs.уровень_доступа where логин = 'kukushkin@r26' OR отдел = 26911 OR должность = 44 OR (отдел IS NULL AND логин IS NULL)
//                                         )level
//    )
//SELECT distinct iobot.диалоги_ио.* FROM iobot.диалоги_ио
//join fssp_configs.уровень_доступа ON iobot.диалоги_ио.доступ ?| (
//select array_agg(level.уровень) from (
//select jsonb_array_elements_text(уровень) уровень  from fssp_configs.уровень_доступа where логин = 'kukushkin@r26' OR отдел = 26911 OR должность = 44 OR (отдел IS NULL AND логин IS NULL)
//)level
//)
func ПолучитьУровеньДоутспа(client *Client) {
	sql := `SELECT * FROM iobot.диалоги_ио
left join fssp_configs.уровень_доступа ON fssp_configs.уровень_доступа.уровень ?| (
    select array_agg(level.уровень) from (
           select jsonb_array_elements_text(уровень) уровень  from fssp_configs.уровень_доступа where логин = $1 OR отдел = $2 OR должность = $3 OR (отдел IS NULL AND логин IS NULL)
                                         )level
    )`
	_, err := sqlStruct{
		Name: "уровень_доступа",
		Sql:  sql,
		Values: [][]byte{
			[]byte(client.Login),
		},
	}.Выполнить(nil)
	if err != nil {
		log.Printf(">>>> Ошибка SQL запроса: %+v \n\n", err)
	}
}

func (client *Client) ReadMessage() {
	for {
		_, message, err := client.Ws.ReadMessage()
		if err != nil {
			//Ошибка("err %+v\n", err)
			//log.Printf("read:err %v, CloseGoingAway %v", err, websocket.CloseGoingAway)

			//log.Printf("websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure): %v",
			//	websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure))

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				//log.Printf("error: %v", err)
			}
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				//log.Printf("error: %v", err)
			}
			break
		}

		mes := Сообщение{}

		err = json.Unmarshal(message, &mes)
		//Инфо(" Unmarshal message%+v", string(message))
		if err != nil {
			Ошибка("err %+v message %+v\n", err, string(message))
		}
		if client.ВходящееСобщение == nil {
			client.ВходящееСобщение = map[int]Сообщение{}
		}
		client.ВходящееСобщение[mes.Id] = mes
		//Инфо("read message client.ВходящееСобщение %+v", client.ВходящееСобщение)
		if mes.От == "" {
			mes.От = client.Login
		}
		//mes.Id, mes.Время = СохранитьСообщение(mes)
		//СохранитьСообщение(mes)
		mes.СохранитьЛогСообщения()
		if mes.Кому != "io" {
			if server.Clients[mes.Кому] != nil {
				server.Clients[mes.Кому].Message <- &mes
			}
		} else {
			go client.ИО(mes)
		}

		//client.message<-message
	}
}

func (client *Client) Write(b []byte) (int, error) {
	go func() {
		responseMes := Сообщение{
			Id:          0,
			От:          "io",
			Кому:        client.Login,
			Текст:       string(b),
			MessageType: []string{"server_log"},
		}

		// проверим открыт ли канал, если открыт то отправим в него данные иначе удалим получателя из списка отправка ЛОГОВ
		//		if канал, ok := (<-client.message); ok {
		//			log.Printf("канал %+v\n", канал)
		client.Message <- &responseMes
		//		} else {
		//			log.Printf("Канал закрыт канал %+v %+v\n",канал, ok)
		//mw := io.MultiWriter(os.Stdout,IoLoger) //, server.Clients["maksimchuk@r26"]
		//log.SetOutput(mw)
		//}
	}()

	return len(b), nil
}

func (client *Client) SendMessage() {
	for {
		select {
		case q := <-client.Message:
			//log.Printf("q %+v,  %+v\n", q)
			response, err := json.Marshal(q)
			if err != nil {
				log.Printf("SendMessage response err Marshal: %+v\n", err)
			}
			//clinetR:=server.Clients[q["To"]]
			//log.Printf("Проверим активность канала %+v\n",client )
			//err =  client.ws.WriteMessage(1, response)
			//канал, ok  <-client.message
			//log.Printf("канал %+v, ok %+v\n", канал, ok)
			//if  ok {

			err = client.Ws.WriteMessage(1, response)

			if err != nil {
				//Ошибка("SendMessage client.ws (%+v)\n", client)
				//Ошибка("SendMessage client.message response (%+v)\n", string(response))
				//Ошибка("SendMessage WriteMessage err (%+v) response(%+v) client(%+v) \n", err, string(response), client)
				return
			}
			//} else {
			//	log.Printf("Канал пользователя  %+v канал %+v ok %+v закрыт %+v\n", client.Login, ok, канал,client.message)
			//	return
			//}

		}
	}
}

/*
ОбработатьОщибку, получает на вход sql строку, и данные для подстановки в запрос, Функция сама отправляет сообщение клиенту если обработчик умпешно отработал и вернул данные с полем ,message*/
func (client *Client) ОбработатьОшибкуБД(ОбработчикОшибки interface{}, data interface{}) {

	SQLStateСкрипт := ОбработчикОшибки.(map[string]interface{})["sql"]

	SQLStateString := textTpl.Must(textTpl.New("SQLStateСкриптЗапрос").Parse(SQLStateСкрипт.(string)))

	БайтБуферSQLState := new(bytes.Buffer)

	err := SQLStateString.Execute(БайтБуферSQLState, data)

	if err != nil {
		log.Println("executing template:", err)
	} else {
		//SQLСкрипт = БайтБуферSQLState.String()
	}

	Результат, err := sqlStruct{
		Name:   "",
		Sql:    БайтБуферSQLState.String(),
		Values: [][]byte{},
	}.Выполнить(nil)

	if err != nil {
		log.Printf(">>>> ERROR \n %+v \n\n", err)
	}

	// алгоритм т.к. это ошибка выполнения запроса, и нам нужно сообщить об этлм клиенту, то выведем сразу сообщение со стандартным шаблном ошибки.
	if len(Результат) == 1 {
		ИнформацияКлиенту := Результат[0]["message"]
		//type ДанныеОтвета struct {
		//	Контейнер string `json:"контейнер"`
		//	Данные interface{} `json:"данные"`
		//	HTML string `json:"html"`
		//	Обработчик string `json:"обработчик"` //JS функция или объект/класс/плагин для обработки данных (table..)
		//}
		ДанныеДляОтвета := &ДанныеОтвета{
			Обработчик: "FloatMessage",
			Данные:     ИнформацияКлиенту,
		}
		СообщениеКлиенту := &Сообщение{
			От:          "io",
			Кому:        client.Login,
			MessageType: []string{"error"},
			Контэнт:     ДанныеДляОтвета,
		}
		log.Printf("СообщениеКлиенту %+v\n", СообщениеКлиенту)
		СообщениеКлиенту.СохранитьИОтправить(client)
	}

}

/*
ОбработатьSQLСкрипты на вход получает объект со скриптами, выполняет каждый из них, сохраняет в карту, и возвращает карту с данными всех запросов
Если в резулттате какогто запроса возникла ошибка, и для этой ошибки есть обработчик, то обработчик ыполниться, и сообщит клиенту об ошибке
*/

func ВходящиеДанныеДействий(вопрос *Сообщение) map[string]interface{} {

	Данные := map[string]interface{}{}

	for НазваниеДействия, ДанныеДействия := range вопрос.Выполнить.Действие {
		Инфо("НазваниеДействия %+v ДанныеДействия %+v\n", НазваниеДействия, ДанныеДействия)
		//Инфо("reflect.TypeOf(ДанныеДействия) %+v", reflect.TypeOf(ДанныеДействия).Kind())
		for Имя, Значение := range ДанныеДействия {
			Данные[Имя] = Значение
		}
	}
	return Данные
}

//  ОтветитьКлиенту - обрабатывает входящий запрос от клиента в соответсвии со сценаием из диалоги_ио
func (client *Client) ОбработатьСообщениеИОтветитьКлиенту(ДейсвтияДляОтветаНаЗапрос map[string]interface{}, вопрос *Сообщение) {
	//_, ОтветЕсть := ДейсвтияДляОтветаНаЗапрос["ответ"]
	//Инфо("ответ %+v \n", ДейсвтияДляОтветаНаЗапрос)
	//Инфо("ОтветЕсть %+v\n", ОтветЕсть)
	//if ОтветЕсть {

	var ЗагруженныйФайл map[string]string
	for НазваниеДействия, ДанныеДействия := range вопрос.Выполнить.Действие {
		Инфо(" НазваниеДействия %+v", НазваниеДействия)
		if Файлы, ЕстьФайлы := ДанныеДействия["Файл"]; ЕстьФайлы {
			Инфо("ЕстьФайлы  %+v", ЕстьФайлы)
			/*
				Файлы = {
					"Папка":name,
					"Файл": ЧтениеФайла.result,
					"ИмяФайла": value.name
				}
			*/
			var ПутьДляЗагрузки interface{}
			var ЕстьПуть bool
			if ПутьДляЗагрузки, ЕстьПуть = ДанныеДействия["filePath"]; !ЕстьПуть {
				ПутьДляЗагрузки = "uploads/tmp"
			}
			var ОшибкаЗагрузки error

			НеМенятьИмя := ДанныеДействия["не_менять_имя"]
			if НеМенятьИмя == nil {
				НеМенятьИмя = "false"
			}
			ЗагруженныйФайл, ОшибкаЗагрузки = client.СохранитьФайл(Файлы.(map[string]interface{}), ПутьДляЗагрузки.(string), НеМенятьИмя.(string))
			//map[string]string{
			//	"ИмяФайла":ИмяФайла,
			//	"Папка":ПодПапка,
			//	"ПутьДляЗагрузки": ПутьДляЗагрузки,
			//}
			if ОшибкаЗагрузки != nil {
				Ошибка("  %+v", ОшибкаЗагрузки)
				ДанныеДляОтвета := &ДанныеОтвета{
					Обработчик: "FloatMessage",
					Данные:     "Не удалось сохранить файл " + ОшибкаЗагрузки.Error(),
				}

				СообщениеКлиенту := &Сообщение{
					От:          "io",
					Кому:        client.Login,
					MessageType: []string{"error"},
					Контэнт:     ДанныеДляОтвета,
				}
				client.Message <- СообщениеКлиенту
			} else {
				ДанныеДействия["Файл"] = ЗагруженныйФайл
			}

			//} else {
			//	if  !ЕстьПуть {
			//		Ошибка(" Нет filePath  %+v", ДанныеДействия)
			//	}
			//}

		}

	}

	if ДейсвтияДляОтветаНаЗапрос["ответ"] != nil {
		ОтветИО := ДейсвтияДляОтветаНаЗапрос["ответ"].(map[string]interface{})

		/*алгоритм  наличие утверждения в ответе исключает вопрос и  ожидание ответа от клиента
					Если есть утвверждение то сразу обращаемся к полю выполнить, и проверяем наличие "выполнить"
				   Если есть "выполнить" то выполняем его, и если "выполнить" вернуло данные с полем "ОтветКлиенту" то добавляем ответ в текст сообщения,  если "выполнить" вернуло данные с полями HTML и таргет то помещяем эти данные в соответсвтующие поля в сообщение клиенту
		  алгоритм
					цель: куда вставлять хтмл с результатами,
					всегда должен быть полный путь IDs блоков через точку (НЕ КЛАССОВ А ID)
				   например main_content.tickets_list   (если на странице есть ID то данные в этом блоке обновятьсяб если способ вставки HTML шаблона помечен как обновить)
				   обработка блоков для вставки идёт с права на лево, если есть tickets_list обновляем данные внутри,
					если нет tickets_list  , то ищем main_content и вставляем данные в него , если нет main_content то ничего не делаем
				   имя хтмл шаблона который будет возвращён клиенту в виде ответа с результатом sql запроса,
				   если нет sql запроса то возращаеться  просто хтмл
		*/

		ЕстьУтверждение := ОтветИО["утверждение"]
		//Инфо("ЕстьУтверждение %+v\n", ЕстьУтверждение)
		if ЕстьУтверждение != nil {
			утверждение := ЕстьУтверждение.(string)

			номер_диалога := strconv.Itoa(ДейсвтияДляОтветаНаЗапрос["номер_диалога"].(int))
			номер_сообщения := strconv.Itoa(ДейсвтияДляОтветаНаЗапрос["номер_сообщения"].(int))

			_, err := sqlStruct{
				Name: "история_диалогов_ио",
				Sql:  "INSERT INTO iobot.история_диалогов_ио (клиент, номер_диалога, номер_сообщения, время_сообщения, завершено, сообщение_клиента) VALUES ($1,$2,$3,NOW(),true, $4)",
				Values: [][]byte{
					[]byte(client.Login),
					[]byte(номер_диалога),
					[]byte(номер_сообщения),
					[]byte(вопрос.Текст),
				},
			}.Выполнить(nil)

			if err != nil {
				Ошибка(">>>> ERROR \n %+v \n\n", err)
			}
			//Инфо("%+v",утверждение)
			tpl := template.Must(template.New("Сообщение").Parse(утверждение))
			resHtml := new(bytes.Buffer)
			err = tpl.Execute(resHtml, client.UserInfo)
			if err != nil {
				Ошибка("executing template:", err)
			} else {
				утверждение = resHtml.String()
			}

			СообщениеКлиенту := &Сообщение{
				Текст: утверждение,
				От:    "io",
				Кому:  client.Login,
			}

			//Инфо("ОтветИО %+v\n", ОтветИО)

			if ДействиеДляОТвета := ОтветИО["выполнить"]; ДействиеДляОТвета != nil {
				Инфо("ОтветИО[выполнить] %+v ЕстьДействиеДляОТвета %+v\n", ОтветИО["выполнить"])
				//log.Printf("ОтветИО[выполнить] %+v\n", ОтветИО["выполнить"])
				client.ВыполнитьДействиеДляОтвета(ДействиеДляОТвета.(string), вопрос, СообщениеКлиенту)
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
		}

		вопросКлиенту := ОтветИО["вопрос"]
		if вопросКлиенту != nil {
			Инфо("вопросКлиенту %+v\n", вопросКлиенту)
			ожидает := ОтветИО["ожидает"]
			if ожидает != nil {

				номер_диалога := strconv.Itoa(ДейсвтияДляОтветаНаЗапрос["номер_диалога"].(int))
				номер_сообщения := strconv.Itoa(ДейсвтияДляОтветаНаЗапрос["номер_сообщения"].(int))

				_, _ = sqlStruct{
					Name: "история_диалогов_ио",
					Sql:  "INSERT INTO iobot.история_диалогов_ио (клиент, номер_диалога, номер_сообщения, время_сообщения, ожидает, завершено, сообщение_клиента) VALUES ($1,$2,$3,NOW(),$4,false, $5)",
					Values: [][]byte{
						[]byte(client.Login),
						[]byte(номер_диалога),
						[]byte(номер_сообщения),
						[]byte(ожидает.(string)),
						[]byte(вопрос.Текст),
					},
				}.Выполнить(nil)
			}
			ВариантыОжидаемыхОтветов := []string{}
			if ожидает.(string) == "вариант_ответа" {

				Ошибка("Нет обработчика %+v ", ожидает.(string))

				//Далее, ЕстьСледующийШаг := Диалог["далее"]
				//if ЕстьСледующийШаг && Далее != nil {
				//	for _, Вариант := range Далее.([]interface{}){
				//		if ОжидаемыйОтвет, ЕстьОжидаемыйВариант := Вариант.(map[string]interface{})["ответ"];ЕстьОжидаемыйВариант{
				//			ВариантыОжидаемыхОтветов= append(ВариантыОжидаемыхОтветов, ОжидаемыйОтвет.(string))
				//		}
				//	}
				//}
			}

			tpl := template.Must(template.New("Сообщение").Parse(вопросКлиенту.(string)))
			resHtml := new(bytes.Buffer)

			ДанныеДляГенерацииОтвета := map[string]interface{}{
				"client":                   client.UserInfo.Info,
				"ВариантыОжидаемыхОтветов": ВариантыОжидаемыхОтветов,
			}

			err := tpl.Execute(resHtml, ДанныеДляГенерацииОтвета)
			if err != nil {
				Ошибка("executing template:", err)
			} else {
				вопросКлиенту = resHtml.String()
			}

			if len(ВариантыОжидаемыхОтветов) > 0 {
				вопросКлиенту = вопросКлиенту.(string) + string(render("variantsArray", ДанныеДляГенерацииОтвета))
			}

			СообщениеКлиенту := &Сообщение{
				Текст: вопросКлиенту.(string),
				От:    "io",
				Кому:  client.Login,
			}

			if ДействиеДляОТвета, ЕстьДействиеДляОТвета := ОтветИО["выполнить"]; ЕстьДействиеДляОТвета {
				client.ВыполнитьДействиеДляОтвета(ДействиеДляОТвета.(string), вопрос, СообщениеКлиенту)
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
			//client.Message<-СообщениеКлиенту
		}

		if ЕстьУтверждение == nil && вопросКлиенту == nil {
			if ДействиеДляОТвета := ОтветИО["выполнить"]; ДействиеДляОТвета != nil {
				Инфо("ОтветИО[выполнить] %+v\n", ОтветИО["выполнить"])
				СообщениеКлиенту := &Сообщение{
					Текст: "",
					От:    "io",
					Кому:  client.Login,
				}
				client.ВыполнитьДействиеДляОтвета(ДействиеДляОТвета.(string), вопрос, СообщениеКлиенту)
				СообщениеКлиенту.СохранитьИОтправить(client)
			}
		}
	}

	/*Если во входящем Запросе есть загружаемый файл то сохраняем его в указаный filePath
	 */

	SQLЗапрос, ЕстьСкрипт := ДейсвтияДляОтветаНаЗапрос["sql_запрос"]

	Инфо("SQLЗапрос  %+v", SQLЗапрос)

	ДанныеHTMLШаблонов := ДейсвтияДляОтветаНаЗапрос["html_шаблон"]

	JSONДанные, ЕстьJSONДанные := ДейсвтияДляОтветаНаЗапрос["json_данные"]

	if ЕстьСкрипт && len(SQLЗапрос.([]interface{})) > 0 { //&& SQLЗапрос.([]interface{})[0] !=nil

		var ДанныеЗапросов map[string]interface{}

		ДанныеЗапросов = client.ОбработатьБезопасноSQLСкрипты(SQLЗапрос.([]interface{}), вопрос)

		//}

		//Инфо("ДанныеЗапросов %+v", ДанныеЗапросов)
		// алг: ВРЕМЕННЫЙ КОСТЫЛЬ, на новом сервере будет по другому

		if ДанныеЗапросов["действие"] != nil {
			//Инфо("ДанныеЗапросов[действие] %+v", ДанныеЗапросов)
			for _, навык := range ДанныеЗапросов["действие"].([]map[string]interface{}) {
				ЗапуститьСкрипт(client, навык)
			}
		}

		ВсеОшибкаЗапросов := map[string]string{}

		for ИмяДанных, Значения := range ДанныеЗапросов {
			//Инфо("ИмяДанных %+v", ИмяДанных)
			//Инфо("Значения %+v", Значения)
			if reflect.TypeOf(Значения).String() == "map[string]interface {}" {
				if ошибка, ok := Значения.(map[string]interface{})["ошибка"]; ok {
					ВсеОшибкаЗапросов[ИмяДанных] = ошибка.(string)
				}
			}
		}
		if len(ВсеОшибкаЗапросов) > 0 {
			for Имя, Ошибка := range ВсеОшибкаЗапросов {
				СообщениеКлиенту := &Сообщение{
					От:          "io",
					Кому:        client.Login,
					MessageType: []string{"error"},
					Контэнт: &ДанныеОтвета{
						Обработчик: "FloatMessage",
						Данные:     Имя + "| " + Ошибка,
					},
				}
				СообщениеКлиенту.СохранитьИОтправить(client)
			}
		}

		if client != nil {
			СообщениеКлиенту := Сообщение{
				Id:          0,
				От:          "io",
				Кому:        client.Login,
				Текст:       "Сбор данных завершён",
				MessageType: []string{"attention"},
				Content: struct {
					Target     string      `json:"target"`
					Data       interface{} `json:"data"`
					Html       string      `json:"html"`
					Обработчик string      `json:"обработчик"`
				}{
					Target: "progress",
				},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
		}

		//Инфо("ДанныеЗапросов  %+v", ДанныеЗапросов)
		/*
			Для постановки данных в sql передадим в data данные клиента, и данные из входящего вопроса: ВходящиеАргументы и вопрос.Выполнить.Действие
		*/

		ДанныеДляРендераСИнфоКлиента := map[string]interface{}{
			"client":  client.UserInfo.Info,
			"data":    ДанныеЗапросов,
			"OspList": client.ПолучитьСписокОтделов(),
			"Osp":     ПолучитьСписокОСП(),
			"вопрос":  вопрос,
			"данные":  ВходящиеДанныеДействий(вопрос),
		}
		//Инфо("ДанныеДляРендераСИнфоКлиента %+v", ДанныеДляРендераСИнфоКлиента["OspList"])

		if ДанныеHTMLШаблонов != nil {
			//Инфо("ДанныеДляРендераСИнфоКлиента %+v %+v", ДанныеДляРендераСИнфоКлиента, ДанныеHTMLШаблонов)
			for _, ДанныеHTMLШаблона := range ДанныеHTMLШаблонов.([]interface{}) {

				НазваниеШаблона := ДанныеHTMLШаблона.(map[string]interface{})["HTML"].(string)
				ЦельДляВставки := ДанныеHTMLШаблона.(map[string]interface{})["цель"]
				Контэнт := ДанныеHTMLШаблона.(map[string]interface{})["Контэнт"]

				if strings.Contains(ЦельДляВставки.(string), "{{") {
					tpl := textTpl.Must(textTpl.New("цель").Parse(ЦельДляВставки.(string)))
					буферШаблона := new(bytes.Buffer)

					err := tpl.Execute(буферШаблона, ДанныеДляРендераСИнфоКлиента)

					if err != nil {
						Ошибка("executing template:", err)
					} else {
						ЦельДляВставки = буферШаблона.String()
					}
				}

				ДанныеДляОтвета := &ДанныеОтвета{}

				if Контэнт != nil {
					данные := Контэнт.(map[string]interface{})["данные"]
					if данные != nil {
						//[{"HTML": "test_doc_upload", "цель": "AccessSCZI", "Контэнт": {"данные": {"{{.data.uid_number}}.{{.data.документ}}": "HTML"}}}]
						for ШаблонКоординатДанных, ЧтоВставлять := range данные.(map[string]interface{}) {
							for ИмяРезультата, МассивДанных := range ДанныеЗапросов {

								if len(МассивДанных.([]map[string]interface{})) > 0 {
									for _, СтрокаДанных := range МассивДанных.([]map[string]interface{}) {
										Инфо(" ИмяРезультата %+v СтрокаДанных %+v", ИмяРезультата, СтрокаДанных)
										ДанныеДляРендераЦели := map[string]interface{}{
											"data":   СтрокаДанных,
											"client": client.UserInfo.Info,
										}
										tpl, errParse := textTpl.New("КоординатыДанных").Parse(ШаблонКоординатДанных)
										if errParse != nil {
											Ошибка("  %+v", errParse)
										}
										буферШаблона := new(bytes.Buffer)

										err := tpl.Execute(буферШаблона, ДанныеДляРендераЦели)
										if err != nil {
											Ошибка("executing template: %+v, ДанныеДляРендераЦели: %+v", err, ДанныеДляРендераЦели)
										}

										Инфо("буферШаблона  %+v", буферШаблона)

										Данные := ЧтоВставлять
										if ЧтоВставлять == "HTML" {
											if strings.Contains(НазваниеШаблона, "{{") {
												tpl := textTpl.Must(textTpl.New("КоординатыДанных").Parse(НазваниеШаблона))
												НазваниеHTMLШаблона := new(bytes.Buffer)
												err := tpl.Execute(НазваниеHTMLШаблона, ДанныеДляРендераЦели)
												if err != nil {
													Ошибка("executing template: %+v, ДанныеДляРендера НазваниеHTMLШаблона: %+v", err, ДанныеДляРендераЦели)
												}
												НазваниеШаблона = НазваниеHTMLШаблона.String()
											}
											Инфо("НазваниеШаблона  %+v", НазваниеШаблона)
											Данные = string(render(НазваниеШаблона, ДанныеДляРендераЦели))
										}

										ДанныеДляВставки := map[string]interface{}{
											буферШаблона.String(): Данные,
										}

										ДанныеДляОтвета = &ДанныеОтвета{
											Контейнер: ЦельДляВставки.(string),
											Данные:    ДанныеДляВставки,
										}
										if СпособВставки := ДанныеHTMLШаблона.(map[string]interface{})["СпособВставки"]; СпособВставки != nil {
											ДанныеДляОтвета.СпособВставки = СпособВставки.(string)
										}
										Инфо("ДанныеДляОтвета  %+v", ДанныеДляОтвета)
										СообщениеКлиенту := &Сообщение{
											От:   "io",
											Кому: client.Login,
											//MessageType: []string{"irritation","io_action"},
											Контэнт: ДанныеДляОтвета,
										}

										СообщениеКлиенту.СохранитьИОтправить(client)
									}
								}
							}

						}

					}
				} else {
					if strings.Contains(НазваниеШаблона, "{{") {
						tpl := textTpl.Must(textTpl.New("КоординатыДанных").Parse(НазваниеШаблона))
						НазваниеHTMLШаблона := new(bytes.Buffer)
						err := tpl.Execute(НазваниеHTMLШаблона, ДанныеДляРендераСИнфоКлиента)
						if err != nil {
							Ошибка("executing template: %+v, ДанныеДляРендера НазваниеHTMLШаблона : %+v", err, ДанныеДляРендераСИнфоКлиента)
						}
						НазваниеШаблона = НазваниеHTMLШаблона.String()
					}

					ДанныеДляОтвета = &ДанныеОтвета{
						Контейнер: ЦельДляВставки.(string),
						HTML:      string(render(НазваниеШаблона, ДанныеДляРендераСИнфоКлиента)),
					}
					if СпособВставки := ДанныеHTMLШаблона.(map[string]interface{})["СпособВставки"]; СпособВставки != nil {
						ДанныеДляОтвета.СпособВставки = СпособВставки.(string)
					}
					СообщениеКлиенту := &Сообщение{
						От:   "io",
						Кому: client.Login,
						//MessageType: []string{"irritation","io_action"},
						Контэнт: ДанныеДляОтвета,
					}

					СообщениеКлиенту.СохранитьИОтправить(client)
				}
			}
		}

		if ЕстьJSONДанные && JSONДанные != nil {
			Инфо("ЕстьJSONДанные %+v JSONДанные %+v\n", ЕстьJSONДанные, JSONДанные)
			//пройдём по массиву , отрендерим и отправим клиенту
			for _, ЭлементДанных := range JSONДанные.([]interface{}) {

				Цель := ЭлементДанных.(map[string]interface{})["цель"].(string)

				tpl := template.Must(template.New("цель").Funcs(tplFunc()).Parse(Цель))

				БайтБуферSql := new(bytes.Buffer)

				err := tpl.Execute(БайтБуферSql, ДанныеДляРендераСИнфоКлиента)
				Инфо("ДанныеДляРендераСИнфоКлиента %+v Цель %+v", ДанныеДляРендераСИнфоКлиента, Цель)

				if err != nil {
					Ошибка(">>>> ERROR \n %+v \n\n", err)
				} else {
					Цель = БайтБуферSql.String()
					Инфо("ДанныеДляРендераСИнфоКлиента %+v Цель %+v", ДанныеДляРендераСИнфоКлиента, Цель)
				}
				ОбъектДанных := ЭлементДанных.(map[string]interface{})["DATA"].(string)

				Шаблон := ЭлементДанных.(map[string]interface{})["Шаблон"]

				Данные := ДанныеЗапросов[ОбъектДанных]
				var HTMLстрока string

				if HTML, ok := ЭлементДанных.(map[string]interface{})["HTML"]; ok {

					tpl, err := textTpl.New("xml").Funcs(tplFunc()).ParseGlob(РабочаяПапка + "/html/*.*")
					if err != nil {
						Ошибка(">>>> ERROR \n %+v  %+v \n\n", err, tpl)
					}
					//tplFiles, err := template.New("").Funcs(tplFunc()).ParseGlob(pattern["pattern"])
					БайтБуферSql := new(bytes.Buffer)

					err = tpl.ExecuteTemplate(БайтБуферSql, HTML.(string), ДанныеДляРендераСИнфоКлиента)
					//err = tpl.Execute(БайтБуферSql, ДанныеДляРендераСИнфоКлиента)
					//Инфо("ДанныеДляРендераСИнфоКлиента %+v Цель %+v", ДанныеДляРендераСИнфоКлиента, БайтБуферSql.String())

					if err != nil {
						Ошибка(">>>> ERROR \n %+v \n\n", err)
					}
					//Цель=БайтБуферSql.String()

					HTMLстрока = БайтБуферSql.String()

					//HTMLстрока = string(render(HTML.(string), ДанныеДляРендераСИнфоКлиента))
				}

				Инфо("ОбъектДанных %+v Шаблон %+v ДанныеЗапросов %+v HTML %+v", ОбъектДанных, Шаблон, ДанныеЗапросов, HTMLстрока)
				//,"{{(index .data.удалённый_документ 0).login}}.{{(index .data.удалённый_документ 0).Документ}}"
				Инфо("Данные %+v", Данные)
				СообщениеКлиенту := &Сообщение{
					От:   "io",
					Кому: client.Login,
				}

				if Данные != nil && Данные != "ошибка" {

					Инфо("ОбъектДанных %+v\n", ОбъектДанных)
					Инфо("получатель== 'сlient' %+v\n", ЭлементДанных.(map[string]interface{})["получатель"] == "сlient")

					if ЭлементДанных.(map[string]interface{})["получатель"] == "сlient" {
						ДанныеДляОтвета := &ДанныеОтвета{
							Контейнер: Цель,
							Данные:    Данные,
						}
						if HTMLстрока != "" {
							ДанныеДляОтвета.HTML = HTMLстрока
						}

						Обработчик := ЭлементДанных.(map[string]interface{})["обработчик"]

						Инфо("ЭлементДанных %+v, Цель %+v, Данные %+v Обработчик %+v ДанныеДляОтвета %+v", ЭлементДанных, Цель, Данные, Обработчик, ДанныеДляОтвета)

						if Обработчик != nil && Обработчик != "" {
							ДанныеДляОтвета.Обработчик = Обработчик.(string)
						}
						СообщениеКлиенту.Контэнт = ДанныеДляОтвета
						Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)

						СообщениеКлиенту.СохранитьИОтправить(client)

					} else if ЭлементДанных.(map[string]interface{})["получатель"] == "всем" { // всем кто онлайн
						for Логин, Получатель := range server.Clients {
							//log.Printf("sendLogin %+v clientOnline.Login %+v\n", )
							if Логин != client.Login {
								СообщениеКлиенту := &Сообщение{
									От:   "io",
									Кому: Логин,
									Контэнт: &ДанныеОтвета{
										Контейнер: Цель,
										Данные:    Данные,
									},
								}
								СообщениеКлиенту.СохранитьИОтправить(Получатель)
							}
						}
					} else if ЭлементДанных.(map[string]interface{})["получатель"] == "админу" { // всем админам онлайн
						// получить список админов, проверить кто из них онлайн и отправить ему данные

					} else if ЭлементДанных.(map[string]interface{})["получатель"] == "всем.на_странице" { // всем кто онлайн на той же странице

					} else if ЭлементДанных.(map[string]interface{})["получатель"] == "админу.на_странице" { // админам онлайн на той же странице

					} else {
						Получатель, ЕстьПолучатель := server.Clients[ЭлементДанных.(map[string]interface{})["получатель"].(string)]
						if ЕстьПолучатель {
							СообщениеКлиенту := &Сообщение{
								От:   "io",
								Кому: Получатель.Login,
								Контэнт: &ДанныеОтвета{
									Контейнер: Цель,
									Данные:    Данные,
								},
							}
							СообщениеКлиенту.СохранитьИОтправить(Получатель)
						}

					}

				}

				if Шаблон != nil {

					tpl, err := textTpl.New("Данные").Funcs(tplFunc()).Parse(Шаблон.(string))
					if err != nil {
						Ошибка("  %+v", err)
					}
					БайтБуферSql := new(bytes.Buffer)

					err = tpl.Execute(БайтБуферSql, ДанныеДляРендераСИнфоКлиента)

					Инфо("ДанныеДляРендераСИнфоКлиента %+s", БайтБуферSql)
					if err != nil {
						Ошибка(">>>> ERROR \n %+v \n\n", err)
					}

					ДанныеДляОтвета := &ДанныеОтвета{
						Контейнер: Цель,
						Данные:    БайтБуферSql.String(),
					}
					Обработчик := ЭлементДанных.(map[string]interface{})["обработчик"]
					if Обработчик != nil {
						ДанныеДляОтвета.Обработчик = Обработчик.(string)
					}
					СообщениеКлиенту := &Сообщение{
						От:      "io",
						Кому:    client.Login,
						Контэнт: ДанныеДляОтвета,
					}
					Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)

					СообщениеКлиенту.СохранитьИОтправить(client)
				}

			}
		}

	} else {
		if ДанныеHTMLШаблонов != nil {
			for _, ДанныеHTMLШаблона := range ДанныеHTMLШаблонов.([]interface{}) {
				Инфо("Рендерим шаблон %+v\n", ДанныеHTMLШаблона)

				НазваниеШаблона := ДанныеHTMLШаблона.(map[string]interface{})["HTML"].(string)
				ЦельДляВставки := ДанныеHTMLШаблона.(map[string]interface{})["цель"]

				Инфо("  %+v", reflect.TypeOf(ЦельДляВставки))
				ДанныеДляРендераСИнфоКлиента := map[string]interface{}{
					"client":  client.UserInfo.Info,
					"OspList": client.ПолучитьСписокОтделов(),
					"Osp":     ПолучитьСписокОСП(),
					"вопрос":  вопрос,
					"данные":  ВходящиеДанныеДействий(вопрос),
				}

				switch ЦельДляВставки.(type) {
				case []interface{}:
					for _, ШаблонЦели := range ЦельДляВставки.([]interface{}) {
						tpl := textTpl.Must(textTpl.New("цель").Parse(ШаблонЦели.(string)))
						буферШаблона := new(bytes.Buffer)

						err := tpl.Execute(буферШаблона, ДанныеДляРендераСИнфоКлиента)

						if err != nil {
							Ошибка("executing template:", err)
						}
						Инфо("буферШаблона  %+v", буферШаблона)
						ДанныеДляОтвета := &ДанныеОтвета{
							Контейнер: буферШаблона.String(),
							HTML:      string(render(НазваниеШаблона, ДанныеДляРендераСИнфоКлиента)),
						}
						if СпособВставки := ДанныеHTMLШаблона.(map[string]interface{})["СпособВставки"]; СпособВставки != nil {
							ДанныеДляОтвета.СпособВставки = СпособВставки.(string)
						}
						СообщениеКлиенту := &Сообщение{
							От:      "io",
							Кому:    client.Login,
							Контэнт: ДанныеДляОтвета,
						}
						СообщениеКлиенту.СохранитьИОтправить(client)
					}

				default:
					Инфо("ДанныеДляРендераСИнфоКлиента %+v\n", ДанныеДляРендераСИнфоКлиента)
					if strings.Contains(ЦельДляВставки.(string), "{{") {
						tpl := textTpl.Must(textTpl.New("цель").Parse(ЦельДляВставки.(string)))
						буферШаблона := new(bytes.Buffer)

						err := tpl.Execute(буферШаблона, ДанныеДляРендераСИнфоКлиента)

						if err != nil {
							Ошибка("executing template:", err)
						} else {
							ЦельДляВставки = буферШаблона.String()
						}
					}
					ДанныеДляОтвета := &ДанныеОтвета{
						Контейнер: ЦельДляВставки.(string),
						HTML:      string(render(НазваниеШаблона, ДанныеДляРендераСИнфоКлиента)),
					}
					if СпособВставки := ДанныеHTMLШаблона.(map[string]interface{})["СпособВставки"]; СпособВставки != nil {
						ДанныеДляОтвета.СпособВставки = СпособВставки.(string)
					}
					СообщениеКлиенту := &Сообщение{
						От:      "io",
						Кому:    client.Login,
						Контэнт: ДанныеДляОтвета,
					}
					СообщениеКлиенту.СохранитьИОтправить(client)
				}

			}

		}
	}
}
func (СообщениеКлиенту *Сообщение) СохранитьИОтправить(client *Client) {
	//Инфо("СохранитьИОтправить %+v", client.Ws.RemoteAddr())
	//СообщениеКлиенту.СохранитьЛогСообщения()

	//Инфо("client.ВходящееСобщение %+v СообщениеКлиенту.Id %+v", client.ВходящееСобщение, СообщениеКлиенту.Id)
	//Инфо("client.ВходящееСобщение[СообщениеКлиенту.Id] %+v", client.ВходящееСобщение[СообщениеКлиенту.Id])
	client.mu.Lock()
	if client.ВходящееСобщение[СообщениеКлиенту.Id].ОбратныйВызов != "" {
		СообщениеКлиенту.ОбратныйВызов = client.ВходящееСобщение[СообщениеКлиенту.Id].ОбратныйВызов
	}

	delete(client.ВходящееСобщение, СообщениеКлиенту.Id)
	client.mu.Unlock()
	client.Message <- СообщениеКлиенту
}
