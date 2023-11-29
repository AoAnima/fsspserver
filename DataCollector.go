package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
	"time"

	"github.com/jackc/pgconn"
	"strconv"
	"strings"
)

type sqlStruct struct{
	ВходящиеАргументы map[string]interface{}
	SQLАргументы []interface{} // список кооординат для ВходящиеАргументы, обрабатываеться функцие получитьАргументы или получитьАргументыдляАИС, и помещаються в Values или АргументыДляАИС

	Name string `json:"name"`// Имя запроса sql
	Table string `json:"table"`
	Sql string `json:"sql"`

	ResultType	string `json:"result_type"`
	ColName string `json:"col_name"`
	Labels string `json:"labels"`
	Type	[]string `json:"type"`
	SubTable string `json:"sub_table"`

	ПеребратьАргументы []interface{} // Аргументы которые нужно перебрать поочередни в запросе

	Dbs             []interface{} `json:"dbs"`
	Values          [][]byte // Аргументы для подстановки в sql для postgresql, не требуют доп обработки
	АргументыДляАИС []interface{} // Аргументы для подстановки в sql для аис, не требуют доп обработки
	БазаДанных string // куда отправлять запрос io, fssp, rdb, osp
	Асинхронно bool // если это поле true то выполняем отправляем выполнение запроса в горутину
	Клиент *Client

	DBSchema        string        `json:"db_schema"`


	Conn *sql.DB
	ДанныеПодключения map[string]interface{} `json:"ДанныеПодключения"` // логин и пароль для подключения к осп

	DateSync string `json:"date_sync"`
	ResultReader *pgconn.ResultReader
	PgError *pgconn.PgError
	ОбновлениеВРеальномВремени bool
	ПараметрыСтолбца map[string]interface{}
}

/*
ЗагрузитьДанные - выбирает данные из БД и возвращает их пользовател, данные выбираються из таблицы по имени модуля/страницы или имени таблицы
*/
type Table struct {
	Setting map[string]interface{} `json:"setting"`
	Id int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Description map[string]interface{} `json:"description,omitempty"`
	Order int `json:"orders,omitempty"`
	Caption string `json:"caption"` // table caption <table><caption>Example Caption</caption>...</table>
	Alias string `json:"alias,omitempty"`
	Columns		map[string]string `json:"columns"`
	Header		[]map[string]interface{} `json:"header"`
	Rows  		[]map[string]interface{} `json:"rows"`
	Footer  		[]map[string]interface{} `json:"footer"`
	CountRows uint32 `json:"countRows"`
	Index 		map[string]map[string]int	`json:"index"` //interface{}
	SubTables bool `json:"sub_tables"`

}

func ПолучитьМетаДанныеТаблицы(tableName string ) ([]map[string]interface{}, string) {

	sqlSting := "SELECT col_name, label, col_parameters FROM sql WHERE sql.table = $1 OR  sql.table = 'any' ORDER BY col_parameters ->> 'order' ASC"


	sqlQuery := sqlStruct{
		Name:   tableName,
		Sql:    sqlSting,
		Values: [][]byte{
				[]byte(tableName),
		},
		DBSchema:"fssp_data",
	}
	Collumns, _ := ВыполнитьPgSQL(sqlQuery)

	sqlStingCaption := "SELECT caption FROM meta_table WHERE meta_table.table = $1"


	queryCaption := sqlStruct{
		Name:   tableName,
		Sql:    sqlStingCaption,
		Values: [][]byte{
			[]byte(tableName),
		},
		DBSchema:"fssp_data",
	}
	CaptionRow, _ := ВыполнитьPgSQL(queryCaption)
	if CaptionRow == nil {
		return nil, "Нет данных для таблицы: "+tableName
	}
		Caption := CaptionRow[0]["caption"].(string)



//Инфо("Header Rows %+v\n", Rows)
	return Collumns, Caption
}


func СохранитьКэшСтатистикиИзРБД(ДанныеДляСохранения map[string]РезультатSQLизАИС){
	//Инфо("ДанныеДляСохранения %+v", ДанныеДляСохранения)
	for ИмяСтолбца, Данные := range ДанныеДляСохранения{

		КартаСтрок, err := json.Marshal(Данные.КартаСтрок)
		if err != nil {
			Ошибка(" %+v ", err)
		}
		АргументыЗапроса, err := json.Marshal(Данные.ВходящиеАргументы)
		if err != nil {
			Ошибка(" %+v ", err)
		}

		if КодОсп, ЕстьКод := Данные.ВходящиеАргументы["osp_code"]; ЕстьКод {
			ИмяСтолбца=ИмяСтолбца+"."+КодОсп.(string)
		}



		Инфо("СохранитьКэшСтатистикиИзРБД  КартаСтрок %+s", КартаСтрок)

		if string(КартаСтрок) == "null"{
			КартаСтрок=nil
		}

		Значения := [][]byte{
			[]byte(Данные.Таблица),
			[]byte(ИмяСтолбца),
			КартаСтрок,
			АргументыЗапроса,
		}

		_ ,err = sqlStruct{
			Name:   "Сохранение кэша",
			Sql:    "INSERT INTO стат_кэш_рбд (таблица, столбец,  данные, аргументы_запроса, дата) VALUES ($1, $2, $3, $4, NOW()) ON CONFLICT ON CONSTRAINT стат_кэш_рбд_pk DO UPDATE SET дата = NOW(), данные = EXCLUDED.данные, аргументы_запроса=EXCLUDED.аргументы_запроса",
			Values: Значения,
			DBSchema:"fssp_data",
		}.Выполнить(nil)

		if err != nil{
			Ошибка(">>>> Ошибка SQL запроса: %+v \n",err)
		}
	}
}

func СохранитьКэшСтатистики(ДанныеДляСохранения map[string]РезультатSQLизАИС){
	//Инфо("ДанныеДляСохранения %+v", ДанныеДляСохранения)
	for ИмяСтолбца, Данные := range ДанныеДляСохранения{

		КартаСтрок, err := json.Marshal(Данные.КартаСтрок)
		if err != nil {
			Ошибка(" %+v ", err)
		}
		ОСП := strconv.Itoa(Данные.ОСП["osp_code"].(int))
		Количество := strconv.Itoa(Данные.Количество)
		Значения := [][]byte{
			[]byte(Данные.Таблица),
			[]byte(ИмяСтолбца),
			[]byte(ОСП),
			КартаСтрок,
			[]byte(Количество),
		}

		_ ,err = sqlStruct{
			Name:   "Сохранение кэша",
			Sql:    "INSERT INTO кэш_стат_данных (таблица, столбец, осп, данные, количество, дата) VALUES ($1, $2, $3,$4, $5, NOW()) ON CONFLICT ON CONSTRAINT кэш_стат_данных_pk DO UPDATE SET дата = NOW(), данные = EXCLUDED.данные, количество=EXCLUDED.количество",
			Values: Значения,
			DBSchema:"fssp_data",
		}.Выполнить(nil)

		if err != nil{
			Ошибка(">>>> Ошибка SQL запроса: %+v \n",err)
		}


	}

}


func (client *Client) ЗагрузитьДанные(mes Сообщение){


	tableNames := mes.Выполнить.Arg.Tables
	Tables:=map[string]Table{}
	//if len(tableNames) >1{
		for _, tableName := range tableNames {

			//tableName := mes.Выполнить.Arg.Tables
			sqlSting := "SELECT * FROM "+tableName+" WHERE date_sync = (select max(date_sync) from "+tableName+")"
			//sqlSting := "SELECT  json_build_object(complaints.sub_table, json_agg(complaints.*)) FROM "+tableName+" WHERE date_sync = (select max(date_sync) from "+tableName+")"



			Инфо("sqlSting %+v\n",sqlSting )
			sqlQuery := sqlStruct{
				Name:   tableName,
				Sql:    sqlSting,
				Values: nil,
				DBSchema:"fssp_data",

			}
			Rows, _ := ВыполнитьPgSQL(sqlQuery)

			Collumns, Caption := ПолучитьМетаДанныеТаблицы(tableName)

			table := Table{
				Alias:tableName,
				Rows:Rows,
				Header: Collumns,
				Caption: Caption,
			}
			Tables[tableName]=table
		}
	//}

	//sqlSting := "SELECT  json_build_object(complaints.sub_table, json_agg(complaints.*)) FROM "+tableName+" WHERE date_sync = (select max(date_sync) from "+tableName+")"
	//sqlSting := "SELECT  json_build_object(complaints.sub_table, json_agg(complaints.*)) FROM "+tableName+" WHERE date_sync = (select max(date_sync) from "+tableName+")"
	//
	//Инфо("sqlSting %+v\n",sqlSting )
	//sqlQuery := sqlStruct{
	//	Name:   tableName,
	//	Sql:    sqlSting,
	//	Values: nil,
	//}
	//Rows, ArrayRows := ВыполнитьPgSQL(sqlQuery)
	//Инфо("ArrayRows %+v\n", ArrayRows)
	//resultString, err  := json.Marshal(ArrayRows)
	////resultRowsString, err  := json.Marshal(Rows)
	//
	//Инфо("resultString %+v\n", resultString)
	//if err != nil {
	//	Инфо("err	 %+v\n", err)
	//}
	//tableData := map[string]interface{}{
	//
	//}
	//resultStringtable, err  := json.Marshal(table)
	//if err != nil {
	//	Инфо("err	 %+v\n", err)
	//}
	//data := map[string]interface{}{"table":table}

	/*
	Data: map[string]interface{}{"tables":Tables}, // map[string]interface{}{"tables":Tables},
	*/
	messToLient := Сообщение{
		Id:      0,
		От:      "io",
		Кому:    client.Login,
		Текст:   "",
		Content: struct {
			Target string `json:"target"`
			Data interface{} `json:"data"`
			Html string `json:"html"`
			Обработчик string `json:"обработчик"`
		}{
			Data: map[string]interface{}{"tables":Tables}, // map[string]interface{}{"tables":Tables},
			Html:"",
			Обработчик: "table",
		},
	}
	client.Message<-&messToLient
}

/*
Получает информацию о sql которые нужно выполнить в аис, по имени витрины/таблицы
выполняет все запросы, запускает ВыполнитьSQL()
*/

func (client *Client) СобратьДанные(mes Сообщение){
	//tableName := mes.Выполнить.Arg.(string)
	//Инфо("tableName %+v\n", tableName)

	sqls, _:=mes.ПолучитьSQLДляСбораДанныхИзОСП()

	result := make(chan string)
	counter := 0
	//Инфо("sqls %+v\n", sqls)
	for _, sqlData := range sqls{
		counter++
		go sqlData.ВыполнитьSQLвАИСиСохранитьДанные( result)
	}
	Инфо("counter %+v\n", counter)


	for {
		q := <-result
		if q=="end"{
			Инфо("q %+v\n", q)
			Инфо("result канал %+v\n", result)
			counter--
		}
		Инфо("counter %+v\n", counter)
		if counter==0{
			Инфо("result канал по идее закрыт %+v\n", result)
			//close(result)
			//Инфо("counter = %+v Выключаем сервер\n", counter)
			//if err := Srv.Shutdown(context.Background()); err != nil {
			//	// Error from closing listeners, or context timeout:
			//	Инфо("HTTP server Shutdown: %v", err)
			//}
		}
	}
}


/*
Получает информацию о sql которые нужно выполнить в аис, по имени витрины/таблицы
выполняет все запросы, запускает ВыполнитьSQL() и возвращает html таблицу с данными
*/
func (client *Client) СобратьДанныеИзОСПВТаблицу(mes Сообщение){
	//tableName := mes.Выполнить.Arg.(string)
	//Инфо("tableName %+v\n", tableName)

	sqls, _ :=mes.ПолучитьSQLДляСбораДанныхИзОСП()
	ДанныеПокдлюченийКОсп, err:=ПолучитьДанныеПодключенийАИСОСП() // получает данные подключений к ОСП в том числе к управлению  243,fssp где хранятся данные только управления
	result := make(chan map[string]РезультатSQLизАИС, len(ДанныеПокдлюченийКОсп))
	counter := 0
	Инфо("sqls %+v\n", sqls)

	if err != nil {
		Ошибка(" %+v ", err)
	}



	РезульатСбораДанных:= map[string]map[string]РезультатSQLизАИС{}

	type СтруктураСтолбцов struct {
		ИмяПеременной string
		ИмяСтолбца string
		НомерСтолбца int
	}
	ПорядокСтолбцов := map[string]СтруктураСтолбцов{}

	for _, ДанныеОСП := range ДанныеПокдлюченийКОсп{
		РезульатСбораДанных[strconv.Itoa(ДанныеОСП["osp_code"].(int))]=map[string]РезультатSQLизАИС{}
		for Порядок, sqlData := range sqls{
			ПорядокСтолбцов[sqlData.ColName]=СтруктураСтолбцов{
				ИмяПеременной: sqlData.ColName,
				ИмяСтолбца: sqlData.Labels,
				НомерСтолбца: Порядок,
			}
			counter++

			go sqlData.ВыполнитьSQLвОСП(ДанныеОСП, result)
		}
	}

	Инфо("counter %+v\n", counter)

	сообщениеКлиенту := Сообщение{
		Id:      0,
		От:      "io",
		Кому:    client.Login,
		Текст:   "Запущенно "+strconv.Itoa(counter)+" запросов",
		MessageType: []string{"attention"},
		Content: struct {
			Target string `json:"target"`
			Data interface{} `json:"data"`
			Html string `json:"html"`
			Обработчик string `json:"обработчик"`
		}{
			Target:"progress",


		},
	}
	client.Message<-&сообщениеКлиенту

	ОСП := map[string]interface{}{}
	for counter >0{
		ДанныеЗапросаИзОСП := <-result
		if ДанныеЗапросаИзОСП != nil{
			for ИмяСтолбца, Данные:= range ДанныеЗапросаИзОСП{
				РезульатСбораДанных[strconv.Itoa(Данные.ОСП["osp_code"].(int))][ИмяСтолбца]=Данные

				_, ok:=ОСП[strconv.Itoa(Данные.ОСП["osp_code"].(int))]
				Инфо("ОСП %+v", ok)
				if !ok {
					ОСП[strconv.Itoa(Данные.ОСП["osp_code"].(int))]=Данные.ОСП
				}

			}
			Инфо("ДанныеИзОСП %+v\n", ДанныеЗапросаИзОСП)
			Инфо("result канал %+v\n", result)

			counter--
		}
		Инфо("counter %+v\n", counter)
		сообщениеКлиенту := Сообщение{
			Id:      0,
			От:      "io",
			Кому:    client.Login,
			Текст:   "Осталось получить данные из"+strconv.Itoa(counter)+" запросов",
			MessageType: []string{"attention"},
			Content: struct {
				Target string `json:"target"`
				Data interface{} `json:"data"`
				Html string `json:"html"`
				Обработчик string `json:"обработчик"`
			}{
				Target:"progress",


			},
		}
		client.Message<-&сообщениеКлиенту
		if counter==0{
			Инфо("result канал по идее закрыт %+v\n", result)
			//close(result)
			//Инфо("counter = %+v Выключаем сервер\n", counter)
			//if err := Srv.Shutdown(context.Background()); err != nil {
			//	// Error from closing listeners, or context timeout:
			//	Инфо("HTTP server Shutdown: %v", err)
			//}
			Обработчик := &ДанныеОтвета{
				Обработчик: "УдалитьЭлемент",
				Контейнер: "progress",
			}
			сообщениеКлиенту = Сообщение{
				Id:      0,
				От:      "io",
				Кому:    client.Login,
				Текст:   "",
				MessageType: []string{},
				Контэнт: Обработчик,
			}
			continue
		}
	}

	ДанныеДляШаблона := map[string]interface{}{
		"ПорядокСтолбцов":ПорядокСтолбцов,
		"Данные":РезульатСбораДанных,
		"ОСП":ОСП,
	}

	html:=render("StatisticTableGenerator", ДанныеДляШаблона)

	сообщениеКлиенту = Сообщение{
		Id:      0,
		От:      "io",
		Кому:    client.Login,
		Текст:   "Данные собраны",
		MessageType: []string{"attention"},
		Content: struct {
			Target string `json:"target"`
			Data interface{} `json:"data"`
			Html string `json:"html"`
			Обработчик string `json:"обработчик"`
		}{
			Target:"main_content.table_wrapper_"+mes.Выполнить.Arg.Module,
			Html:string(html),
		},
	}
	client.Message<-&сообщениеКлиенту
}


func (client *Client)  РендерТаблиц (Таблицы map[string]СтруктураТаблицыДанных, Данные interface{}) {
	//ИдТаблицы := Таблица
	Инфо("РендерТаблиц %+v", Таблицы)

	for ИдТаблицы, Структура := range Таблицы {
		if Структура.ОбновлятьВРеальномВремени {
			ДанныеДляРендера := map[string]interface{}{
				"СтруктураТаблицы": Структура,
				"Данные":           Данные,
			}

			html := render(ИдТаблицы, ДанныеДляРендера)

			if html != nil {

				client.Message <- &Сообщение{
					Id:   0,
					От:   "io",
					Кому: client.Login,
					Контэнт: &ДанныеОтвета{
						Контейнер:     "main_content.table_wrapper_" + ИдТаблицы,
						Данные:        nil,
						HTML:          string(html),
						Обработчик:    "",
						СпособВставки: "обновить",
					},
				}
			}
		}
	}
}

type СтруктураСтолбца struct {
	ИмяПеременной string
	ИмяСтолбца string
	НомерСтолбца int
	ОбновлениеВРеальномВремени bool
	ПараметрыСтолбца map[string]interface{}
	Запрос string
	Аргументы []interface{}
}
type СтруктураТаблицыДанных struct {
	ВходящиеАргументы map[string]interface{}
	ИмяТаблицы string
	СтруктураСтолбцов map[string]СтруктураСтолбца
	ОбновлятьВРеальномВремени bool
}

func (client *Client) СобратьДанныеИзРБД(mes Сообщение)  {

	SQLЗапросы, Таблицы :=mes.ПолучитьSQLДляСбораДанныхИзРДБ()

	//Инфо("ПолучитьSQLДляСбораДанныхИзОСП Таблицы  %+v SQLЗапросы %+v", Таблицы, SQLЗапросы)

	ОСП, err:=ПолучитьДанныеОСП() // получает данные подключений к ОСП в том числе к управлению  243,fssp где хранятся данные только управления
	ДанныеИзАИС := make(chan map[string]РезультатSQLизАИС, len(SQLЗапросы))

	КоличествоЗапросов := 0

	if err != nil {
		Ошибка(" %+v ", err)
	}

	// Рендерим HTMl код для отображения таблицы на клиенте, и отдаём её клиенту если таблица обновляется в реальном времени
	//ДанныеДляРендера := map[string]interface{}{
	//	"ОСП":ОСП,
	//}


	РезультатСбораДанных:= map[string]РезультатSQLизАИС{}
	for _, sqlData := range SQLЗапросы{
		go sqlData.ВыполнитьSQLвРДБ(ДанныеИзАИС)
		КоличествоЗапросов++
	}

	client.Message<-&Сообщение{
		Id:      0,
		От:      "io",
		Кому:    client.Login,
		Текст:   "Отправленно запросов: "+strconv.Itoa(КоличествоЗапросов)+ ".",
		MessageType: []string{"attention"},
	}

	Инфо("КоличествоЗапросов %+v", КоличествоЗапросов)

	for КоличествоЗапросов > 0{
		ДанныеЗапросаИзОСП := <-ДанныеИзАИС
		if ДанныеЗапросаИзОСП != nil{
			СохранитьКэшСтатистикиИзРБД(ДанныеЗапросаИзОСП)
			for ИдСтолбца, Данные:= range ДанныеЗапросаИзОСП{
				РезультатСбораДанных[ИдСтолбца]=Данные

				//РезульатСбораДанныхДляВставкиВтаблицу[Данные.Таблица][ИмяСтолбца]=Данные
				/*
					ИмяСтолбца в базе данных задаёться через точку ид_строки_в_таблице.ид_столбца_в_таблице
					col_name: osp.in_progress - значит на клиенте нужно osp заменить на номер осп из данных, и столбец с id = "in_progress"
					col_name: aj.in_progress - значит на клиенте нужно найти строку с id = "aj" и столбец с id = "in_progress"

				*/
				//КонтейнерДляВставки := Данные.Таблица+"."+ИдСтолбца
				//РезульатСбораДанныхДляВставкиВтаблицу:= map[string]map[string]РезультатSQLизАИС{
				//	Данные.Таблица : map[string]РезультатSQLизАИС{
				//		ИдСтолбца:Данные,
				//	},
				//}

				if Данные.ОбвнолятьВРеальномВремени {
					РезульатСбораДанныхДляВставкиВтаблицу:= map[string]РезультатSQLизАИС{
						ИдСтолбца:Данные,
					}
					client.Message<-&Сообщение{
						Id:      0,
						От:      "io",
						Кому:    client.Login,
						Текст:   "",
						MessageType: []string{},
						Контэнт: &ДанныеОтвета{
							СпособВставки: "",
							Контейнер:     Данные.Таблица,
							Данные:        РезульатСбораДанныхДляВставкиВтаблицу,
							Обработчик:    "",
						},
					}
				}

			}
			КоличествоЗапросов--
			//Инфо("Осталось КоличествоЗапросов %+v", КоличествоЗапросов)
		}

		client.Message<-&Сообщение{
			Id:      0,
			От:      "io",
			Кому:    client.Login,
			Текст:   "Осталось получить данные от "+strconv.Itoa(КоличествоЗапросов)+" источников",
			MessageType: []string{"attention"},
		}


		if КоличествоЗапросов == 0{
			Инфо("КоличествоЗапросов %+v", КоличествоЗапросов)
			//continue
			client.Message<-&Сообщение{
				Id:      0,
				От:      "io",
				Кому:    client.Login,
				Текст:   "Данные собраны",
				MessageType: []string{"attention"},
			}
			break
		}
	}



	/* Отсюда идёт рендер всей таблицы если не было обновления в релаьном времени*/
	ДанныеДляШаблона := map[string]interface{}{
		"Таблицы":Таблицы,
		"Данные":РезультатСбораДанных,
		"ОСП":ОСП,
	}
	client.ОтобразитьСтатТаблицы (ДанныеДляШаблона, mes)

}

//func(client *Client) ПолучитьСтатДанные (mes Сообщение){
//	//ДанныеДляШаблона := client.СобратьДанныеИзРБД(mes)
//	client.ОтобразитьСтатТаблицы (ДанныеДляШаблона, mes)
//}
func(client *Client) ОтобразитьСтатТаблицы (ДанныеДляШаблона map[string]interface{}, mes Сообщение){

	for ИдТаблицы, СтруктураТаблицы  := range ДанныеДляШаблона["Таблицы"].(map[string]СтруктураТаблицыДанных) {
		Инфо("ОбновлятьВРеальномВремени %+v", СтруктураТаблицы.ОбновлятьВРеальномВремени)

		if !СтруктураТаблицы.ОбновлятьВРеальномВремени {

			//Данные := &ДанныеОтвета{
			//	Контейнер:  "main_content.table_wrapper_" + ИдТаблицы,
			//	Данные:     nil,
			//	HTML:       string(render(ИдТаблицы, ДанныеДляШаблона)),
			//	Обработчик: "",
			//}

			Инфо("отправляем данные клиенту Таблица/Шаблон %+v", ИдТаблицы)

			/*алгоритм: если во входящих аргументах есть поле КОнтейнерРЕзультат то добавляем знаечение к полюч контейнер*/


			СообщениеКлиенту := &Сообщение{
				Id:          0,
				От:          "io",
				Кому:        client.Login,
				Текст:       "Данные собраны",
				MessageType: []string{"attention_end"},
				Контэнт:     &ДанныеОтвета{
					Контейнер:  "main_content.table_wrapper_" + ИдТаблицы,
					Данные:     nil,
					HTML:       string(render(ИдТаблицы, ДанныеДляШаблона)),
					Обработчик: "",
				},
			}

			if mes.Выполнить.Действие["СобратьДанныеИзРБД"]["КонтейнерРезультата"] != nil {
				СообщениеКлиенту.Контэнт.Контейнер = mes.Выполнить.Действие["СобратьДанныеИзРБД"]["КонтейнерРезультата"].(string)+"."+ИдТаблицы
			}

			client.Message <- СообщениеКлиенту
		}
	}


}

/* собирает данные по структуркам*/
func (client *Client) СобратьДанныеИзАИС(mes Сообщение){
	//tableName := mes.Выполнить.Arg.(string)
	//Инфо("tableName %+v\n", tableName)

	SQLЗапросы, Таблицы :=mes.ПолучитьSQLДляСбораДанныхИзОСП()

	Инфо("ПолучитьSQLДляСбораДанныхИзОСП Таблицы  %+v", Таблицы)

	ДанныеПокдлюченийКОсп, err:=ПолучитьДанныеПодключенийАИСОСП() // получает данные подключений к ОСП в том числе к управлению  243,fssp где хранятся данные только управления
	ДанныеИзАИС := make(chan map[string]РезультатSQLизАИС, len(ДанныеПокдлюченийКОсп)*len(SQLЗапросы))

	КоличествоЗапросов := 0

	if err != nil {
		Ошибка(" %+v ", err)
	}

	// Рендерим HTMl код для отображения таблицы на клиенте, и отдаём её клиенту если таблица обновляется в реальном времени

	client.РендерТаблиц(Таблицы, ДанныеПокдлюченийКОсп)

	РезульатСбораДанных:= map[string]map[string]РезультатSQLизАИС{}
	//РезульатСбораДанныхДляВставкиВтаблицу:= map[string]map[string]РезультатSQLизАИС{}
	//Таблицы :=  map[string]map[string]СтруктураСтолбцов{}
	//Таблицы :=  map[string]СтруктураТаблицыДанных{}

	ОСП := map[string]interface{}{}

	Инфо("ДанныеПокдлюченийКОсп %+v", len(ДанныеПокдлюченийКОсп))

	for _, ДанныеОСП := range ДанныеПокдлюченийКОсп{

		ОСП[strconv.Itoa(ДанныеОСП["osp_code"].(int))]=ДанныеОСП

		РезульатСбораДанных[strconv.Itoa(ДанныеОСП["osp_code"].(int))]=map[string]РезультатSQLизАИС{}

		for _, sqlData := range SQLЗапросы{
			go sqlData.ВыполнитьSQLвОСП( ДанныеОСП, ДанныеИзАИС)
			КоличествоЗапросов++

		}
	}

	client.Message<-&Сообщение{
		Id:      0,
		От:      "io",
		Кому:    client.Login,
		Текст:   "Отправленно запросов: "+strconv.Itoa(КоличествоЗапросов)+ ". В " + strconv.Itoa(len(ДанныеПокдлюченийКОсп)) + " ОСП",
		MessageType: []string{"attention"},
	}

	Инфо("КоличествоЗапросов start %+v", КоличествоЗапросов)

	for КоличествоЗапросов > 0{
		ДанныеЗапросаИзОСП := <-ДанныеИзАИС
		if ДанныеЗапросаИзОСП != nil{

			СохранитьКэшСтатистики(ДанныеЗапросаИзОСП)
			for ИдСтолбца, Данные:= range ДанныеЗапросаИзОСП{
				РезульатСбораДанных[strconv.Itoa(Данные.ОСП["osp_code"].(int))][ИдСтолбца]=Данные


				//РезульатСбораДанныхДляВставкиВтаблицу[Данные.Таблица][ИмяСтолбца]=Данные
				/*
				ИмяСтолбца в базе данных задаёться через точку ид_строки_в_таблице.ид_столбца_в_таблице
				col_name: osp.in_progress - значит на клиенте нужно osp заменить на номер осп из данных, и столбец с id = "in_progress"
				col_name: aj.in_progress - значит на клиенте нужно найти строку с id = "aj" и столбец с id = "in_progress"


				Данные могут содержать как одно значение , так и массив строк, в этом случае нужно вывести количество строк в ячейке, а массив строк вывести в качестве детализации...

				Нужно Создать метод для рендера таблиц детализации...
				Нужно подтягивать из АИС имена столбцов

				*/
				//КонтейнерДляВставки := Данные.Таблица+"."+ИдСтолбца
				//РезульатСбораДанныхДляВставкиВтаблицу:= map[string]map[string]РезультатSQLизАИС{
				//	Данные.Таблица : map[string]РезультатSQLизАИС{
				//		ИдСтолбца:Данные,
				//	},
				//}

				if Данные.ОбвнолятьВРеальномВремени {
					РезульатСбораДанныхДляВставкиВтаблицу:= map[string]РезультатSQLизАИС{
						ИдСтолбца:Данные,
					}
					client.Message<-&Сообщение{
						Id:      0,
						От:      "io",
						Кому:    client.Login,
						Текст:   "",
						MessageType: []string{},
						Контэнт: &ДанныеОтвета{
							СпособВставки: "",
							Контейнер:     Данные.Таблица,
							Данные:        РезульатСбораДанныхДляВставкиВтаблицу,
							Обработчик:    "",
						},
					}
				}

			}
			КоличествоЗапросов--
			//Инфо("Осталось КоличествоЗапросов %+v", КоличествоЗапросов)
		}

		client.Message<-&Сообщение{
			Id:      0,
			От:      "io",
			Кому:    client.Login,
			Текст:   "Осталось получить данные от "+strconv.Itoa(КоличествоЗапросов)+" источников",
			MessageType: []string{"attention"},
		}


		if КоличествоЗапросов == 0{
			Инфо("КоличествоЗапросов %+v", КоличествоЗапросов)
			//continue
			break
		}
	}



	/* Отсюда идёт рендер всей таблицы если не было обновления в релаьном времени*/
	ДанныеДляШаблона := map[string]interface{}{
		"Таблицы":Таблицы,
		"Данные":РезульатСбораДанных,
		"ОСП":ОСП,
	}

	for ИдТаблицы, СтруктураТаблицы  := range Таблицы {

		if !СтруктураТаблицы.ОбновлятьВРеальномВремени {

			//Данные := &ДанныеОтвета{
			//	Контейнер:  "main_content.table_wrapper_" + ИдТаблицы,
			//	Данные:     nil,
			//	HTML:       string(render(ИдТаблицы, ДанныеДляШаблона)),
			//	Обработчик: "",
			//}

			Инфо("отправляем данные клиенту Таблица/Шаблон %+v", ИдТаблицы)

			client.Message <- &Сообщение{
				Id:          0,
				От:          "io",
				Кому:        client.Login,
				Текст:       "Данные собраны",
				MessageType: []string{"attention_end"},
				Контэнт:     &ДанныеОтвета{
					Контейнер:  "main_content.table_wrapper_" + ИдТаблицы,
					Данные:     nil,
					HTML:       string(render(ИдТаблицы, ДанныеДляШаблона)),
					Обработчик: "",
				},
			}
		}
	}

}

func Детализация (client *Client, вопрос Сообщение) {
	//ЕстьКэш := client.ЗагрузитьКэшДетализации(mes)
	//if !ЕстьКэш {
	//	ТекущееДействие := mes.Выполнить.Действие["ЗагрузитьДетализацию"]
	//	mes.Выполнить.Действие = map[string]map[string]interface{}{
	//		"СобратьДанныеИзРБД":ТекущееДействие,
	//	}
	//	client.СобратьДанныеИзРБД(mes)
	//}


	// алгоритм:
	//	1.Проверить Кэш
	// 		1.1. Если Кэш есть то отобразить данные из кэша
	//	2.Если Кэша нет получить Запрос для сбора данных
	//  3.	Выполнить запрос
	// 		3.1 сохранить данные в Кэш
	// 		3.2 отобразить данные на клиенте



	// Получаем данные из АИС результат приходит в виде map[string]РезультатSQLизАИС
	for ИмяДействия, Аргументы  := range вопрос.Выполнить.Действие{
		Инфо("Детализация ИмяДействия  %+v", ИмяДействия)
		if Аргументы["sql_id"] != nil{


			SQLЗапрос ,ОшибкаЗапроса := sqlStruct{
				Name:   "Получить запрос для сбора статистики	",
				Sql:    "SELECT jsonb_object_agg(coalesce(запросы.имя, 'нет'), jsonb_build_object('скрипт',запросы.скрипт,'аргументы',запросы.аргументы, 'динамический', запросы.динамический, 'аргументы_динамического_запроса', запросы.аргументы_динамического_запроса, 'динамический_шаблон', запросы.динамический_шаблон,'обработчик', запросы.обработчик,'данные_обработчика',запросы.данные_обработчика, 'база_данных',запросы.база_данных)) sql_запрос     FROM iobot.запросы WHERE ид_запроса = $1",

				Values: [][]byte{
					[]byte(Аргументы["sql_id"].(string)),
				},
			}.Выполнить(nil)
			if ОшибкаЗапроса != nil {
				Ошибка("  %+v", ОшибкаЗапроса)
			}

			// переводим []map[string]interface{} в []interface{} если прийдёт в голову более элегантное решениенужно переписать

			SQLЗапросы := make([]interface{}, len(SQLЗапрос))
			for i, v := range SQLЗапрос{
				SQLЗапросы[i]=v
			}

				//SQLЗапросы = SQLЗапрос
				ДанныеЗапросов := client.ОбработатьБезопасноSQLСкрипты(SQLЗапросы, &вопрос)



			Инфо("ДанныеЗапросов  %+v", ДанныеЗапросов)


			for _, РезультатЗапросов := range ДанныеЗапросов {

				for Имя, Данные := range РезультатЗапросов.(map[string]РезультатSQLизАИС) {
					Инфо("Имя данных %+v", Имя)
					if МассивСтрок := Данные.МассивСтрок; МассивСтрок != nil {

						//for _,  МассивСтрок := Данные.(map[string]РезультатSQLизАИС);МассивСтрок != nil{
						//
						ДанныеДляГенратораТаблиц := map[string]interface{}{}

						for _, СтрокаДанных := range МассивСтрок {
							if СтрокаДанных["COLLUMNS"] != nil {
								var Столбцы []string

								err := json.Unmarshal([]byte(СтрокаДанных["COLLUMNS"].(string)), &Столбцы)
								if err != nil {
									Ошибка("  %+v", err)
								}
								ДанныеДляГенратораТаблиц["Столбцы"] = Столбцы

							}
							if СтрокаДанных["LIST"] != nil {
								var Строки interface{}
								err := json.Unmarshal([]byte(СтрокаДанных["LIST"].(string)), &Строки)

								if err != nil {
									Ошибка("  %+v", err, СтрокаДанных["LIST"].(string))
								}
								ДанныеДляГенратораТаблиц["Строки"] = Строки

							}
						}

						ИдТаблицы := Аргументы["id_таблицы"].(string)

						СообщениеКлиенту := &Сообщение{
							Id:          0,
							От:          "io",
							Кому:        client.Login,
							Текст:       "Данные собраны",
							MessageType: []string{"attention_end"},
							Контэнт: &ДанныеОтвета{
								Контейнер:  "main_content.table_wrapper_" + ИдТаблицы,
								Данные:     nil,
								HTML:       string(render("TableGenerator", ДанныеДляГенратораТаблиц)),
								Обработчик: "",
							},
						}

						if Аргументы["КонтейнерРезультата"] != nil {
							СообщениеКлиенту.Контэнт.Контейнер = Аргументы["КонтейнерРезультата"].(string) + "." + ИдТаблицы
						}

						client.Message <- СообщениеКлиенту
					}
				}
			}
		}

	}
}



func (client *Client) ЗагрузитьДетализацию(mes Сообщение){

	ЕстьКэш := client.ЗагрузитьКэшДетализации(mes)
	if !ЕстьКэш {
		ТекущееДействие := mes.Выполнить.Действие["ЗагрузитьДетализацию"]
		mes.Выполнить.Действие = map[string]map[string]interface{}{
			"СобратьДанныеИзРБД":ТекущееДействие,
		}
		client.СобратьДанныеИзРБД(mes)
	}
}

func (client *Client) ЗагрузитьКэшДетализации(mes Сообщение) bool {

	ЕстьКэш := false

	for _, АргументыДействия := range mes.Выполнить.Действие {
		if Таблицы, есть := АргументыДействия["table[]"]; есть {

			Инфо("Таблицы %+v", Таблицы)

			for _, таблица := range Таблицы.([]interface{}) {

				Инфо(" %+v", таблица)

				if КодОсп, ЕстьКод := АргументыДействия["osp_code"]; ЕстьКод {

					Таблица:=таблица.(string)
					Инфо(" %+v", АргументыДействия)

					//if Есть := АргументыДействия["столбец"]; {
					//
					//}

					Столбец := АргументыДействия["столбец"].(string)+"."+КодОсп.(string)

					СтрокаПериода := ""

					АргументыЗапроса := [][]byte{
						[]byte(Таблица),
						[]byte(Столбец),
					}

					if НачалоПериода, естьПериод :=  АргументыДействия["date_from"].(string); естьПериод {
						СтрокаПериода= СтрокаПериода + ` and аргументы_запроса->>'date_from' = $3 AND аргументы_запроса->>'date_to' = $4`
						АргументыЗапроса = append(АргументыЗапроса, []byte(НачалоПериода))
						АргументыЗапроса = append(АргументыЗапроса, []byte(АргументыДействия["date_to"].(string)))
					}



					Данные ,err := sqlStruct{
						Name:   "Получение кэша детализации для "+Таблица,
						Sql:    `select таблица,столбец, данные as "КартаСтрок", дата, аргументы_запроса from fssp_data.стат_кэш_рбд where таблица = $1 AND столбец = $2`+ СтрокаПериода,
						DBSchema:"fssp_data",
						Values: АргументыЗапроса,
					}.Выполнить(nil)

					if err != nil{
						Ошибка(">>>> Ошибка SQL запроса: %+v \n",err , Таблица, Столбец)
					}
					if len(Данные)>0{
						ЕстьКэш = true
					} else {
						continue
					}

					ОСП, err := ПолучитьДанныеОСП()

					if err != nil {
						Ошибка("Данные осп %+v ", err)
					}

					ДанныеДляШаблона := map[string]interface{}{
						"Данные":Данные,
						"ОСП": ОСП,
					}


					HTML := render(Таблица, ДанныеДляШаблона)

					СообщениеКлиенту := &Сообщение{
						Id:          0,
						От:          "io",
						Кому:        client.Login,
						Текст:       "Данные собраны",
						MessageType: []string{"attention_end"},
						Контэнт:     &ДанныеОтвета{
							Контейнер:  "main_content.table_wrapper_" + Таблица,
							Данные:     nil,
							HTML:       string(HTML),
							Обработчик: "",
						},
					}

					if АргументыДействия["КонтейнерРезультата"] != nil {
						СообщениеКлиенту.Контэнт.Контейнер = АргументыДействия["КонтейнерРезультата"].(string)+"."+Таблица
					}
					client.Message <- СообщениеКлиенту
				}
			}
		}
	}
return ЕстьКэш
}

func (client *Client) ЗагрузитьКэшСтатистикиРБД(mes Сообщение)  {

	СтрокаТаблиц := ""

	for _, АргументыДействия := range mes.Выполнить.Действие {
		if Таблицы, есть := АргументыДействия["table[]"]; есть {
			Инфо("Таблицы %+v", Таблицы)
			for _, таблица := range Таблицы.([]interface{}) {
				Инфо(" %+v", таблица)
				if СтрокаТаблиц != "" {
					СтрокаТаблиц = СтрокаТаблиц + ", "
				}
				СтрокаТаблиц = СтрокаТаблиц + "'" + таблица.(string) + "'"
			}
		}

	}



Инфо("СтрокаТаблиц %+v", СтрокаТаблиц)


	Кэш ,err := sqlStruct{
		Name:   "Получение кэша",
		Sql:    `select  jsonb_build_object(таблица, jsonb_object_agg(t1.OSP ,t1.данные)) данные from
    (select таблица, t.OSP, jsonb_object_agg(t.столбец, jsonb_build_object('количество', t.количество, 'время',дата, 'аргументы_запроса',аргументы_запроса)) данные from
        (select таблица, столбец, jsonb_array_elements(данные)->> 'OSP_CODE' as OSP, jsonb_array_elements(данные)->> 'количество' as "количество" , дата, аргументы_запроса from fssp_data.стат_кэш_рбд WHERE таблица in (`+СтрокаТаблиц+`)
        )t group by t.OSP, t.таблица
    ) t1 group by t1.таблица;`,
		DBSchema:"fssp_data",
	}.Выполнить(nil)
	if err != nil{
		Ошибка(">>>> Ошибка SQL запроса: %+v \n",err)
	}

	ОСП, err := ПолучитьДанныеОСП()
	if err != nil {
		Ошибка("Данные осп %+v ", err)
	}


for _ , Таблицы := range Кэш {
	for ИмяТаблицы, Данные := range Таблицы["данные"].(map[string]interface{}) {
		ДанныеДляШаблона := map[string]interface{}{
			"Данные":Данные,
			"ОСП": ОСП,
		}

		HTML := render(ИмяТаблицы, ДанныеДляШаблона)

		client.Message <- &Сообщение{
			Id:          0,
			От:          "io",
			Кому:        client.Login,
			Текст:       "Данные собраны",
			MessageType: []string{"attention_end"},
			Контэнт:     &ДанныеОтвета{
				Контейнер:  "main_content.table_wrapper_" + ИмяТаблицы,
				Данные:     nil,
				HTML:       string(HTML),
				Обработчик: "",
			},
		}
	}
}

}

func (client *Client) ЗагрузитьКэшСтатистики(mes Сообщение){

	СтрокаТаблиц := ""
	for _, АргументыДействия := range mes.Выполнить.Действие {

		if Таблицы, есть := АргументыДействия["table[]"]; есть {
			Инфо("Таблицы %+v", Таблицы)
			for _, таблица := range Таблицы.([]interface{}) {

				Инфо(" %+v", таблица)
				if СтрокаТаблиц != "" {
					СтрокаТаблиц = СтрокаТаблиц + ", "
				}
				СтрокаТаблиц = СтрокаТаблиц + "'" + таблица.(string) + "'"

			}
		}
	}
	_ ,err := sqlStruct{
		Name:   "Получение кэша",
		Sql:    "SELECT * FROM кэш_стат_данных WHERE таблица in ("+СтрокаТаблиц+")",
		DBSchema:"fssp_data",
	}.Выполнить(nil)

	if err != nil{
		Ошибка(">>>> Ошибка SQL запроса: %+v \n",err)
	}

}

func (client *Client) СобратьДанныеИзОСПВEXCEL(mes Сообщение){
	//tableName := mes.Выполнить.Arg.(string)
	Инфо("СобратьДанныеИзОСПВEXCEL %+v\n", mes)

	sqls, _:=mes.ПолучитьSQLДляСбораДанныхИзОСП()

	ДанныеПокдлюченийКОсп, err:=ПолучитьДанныеПодключенийАИСОСП()
	result := make(chan map[string]РезультатSQLизАИС, len(ДанныеПокдлюченийКОсп))
	counter := 0
	//Инфо("sqls %+v\n", sqls)

	if err != nil {
		Ошибка(" %+v ", err)
	}

	РезульатСбораДанных:= map[string]map[string]РезультатSQLизАИС{}

	type СтруктураСтолбцов struct {
		ИмяПеременной string
		ИмяСтолбца string
		НомерСтолбца int
	}
	ПорядокСтолбцов := map[string]СтруктураСтолбцов{}

	for _, ДанныеОСП := range ДанныеПокдлюченийКОсп{
		РезульатСбораДанных[strconv.Itoa(ДанныеОСП["osp_code"].(int))]=map[string]РезультатSQLизАИС{}
		for Порядок, sqlData := range sqls{
			ПорядокСтолбцов[sqlData.ColName]=СтруктураСтолбцов{
				ИмяПеременной: sqlData.ColName,
				ИмяСтолбца: sqlData.Labels,
				НомерСтолбца: Порядок,
			}
			counter++

			go sqlData.ВыполнитьSQLвОСП(ДанныеОСП, result)
		}
	}

	Инфо("counter %+v\n", counter)
	сообщениеКлиенту := Сообщение{
		Id:      0,
		От:      "io",
		Кому:    client.Login,
		Текст:   "Запущенно "+strconv.Itoa(counter)+" запросов",
		MessageType: []string{"attention"},
		Content: struct {
			Target string `json:"target"`
			Data interface{} `json:"data"`
			Html string `json:"html"`
			Обработчик string `json:"обработчик"`
		}{
			Target:"progress",


		},
	}
	client.Message<-&сообщениеКлиенту

	ОСП := map[string]interface{}{}
	for counter >0{
		ДанныеЗапросаИзОСП := <-result
		if ДанныеЗапросаИзОСП != nil{
			for ИмяСтолбца, Данные:= range ДанныеЗапросаИзОСП{
				РезульатСбораДанных[strconv.Itoa(Данные.ОСП["osp_code"].(int))][ИмяСтолбца]=Данные

				if _, ok:=ОСП[strconv.Itoa(Данные.ОСП["osp_code"].(int))]; !ok { // информация об осп код имя адрес
					ОСП[strconv.Itoa(Данные.ОСП["osp_code"].(int))]=Данные.ОСП
				}

			}
			//Инфо("ДанныеИзОСП %+v\n", ДанныеЗапросаИзОСП)
			//Инфо("result канал %+v\n", result)

			counter--
		}
		Инфо("counter %+v\n", counter)
		сообщениеКлиенту := Сообщение{
			Id:      0,
			От:      "io",
			Кому:    client.Login,
			Текст:   "Осталось получить данные из"+strconv.Itoa(counter)+" запросов",
			MessageType: []string{"attention"},
			Content: struct {
				Target string `json:"target"`
				Data interface{} `json:"data"`
				Html string `json:"html"`
				Обработчик string `json:"обработчик"`
			}{
				Target:"progress",
			},
		}
		client.Message<-&сообщениеКлиенту
		if counter==0{
			Инфо("result канал по идее закрыт %+v\n", result)
			//close(result)
			//Инфо("counter = %+v Выключаем сервер\n", counter)
			//if err := Srv.Shutdown(context.Background()); err != nil {
			//	// Error from closing listeners, or context timeout:
			//	Инфо("HTTP server Shutdown: %v", err)
			//}
			break
		}
	}
	Инфо(" РезульатСбораДанных %+v", РезульатСбораДанных)
	f := excelize.NewFile()
	НомерСтроки := 2
	f.SetCellValue("Sheet1", "A1", "Пользователь")
	f.SetCellValue("Sheet1", "B1", "Вход")
	f.SetCellValue("Sheet1", "C1", "Выход")
	f.SetCellValue("Sheet1", "D1", "IP")

	f.SetColWidth("Sheet1", "A", "D", 30)

	for _, data:= range РезульатСбораДанных{
		cellName,err := excelize.CoordinatesToCellName(1,НомерСтроки)
		НомерСтроки++
		f.SetCellValue("Sheet1", cellName, data["user_session"].ОСП["osp_name"])

		for _, rowdata := range data["user_session"].Строки	{

			for cellIdx, cellData := range rowdata{
				cellName,err = excelize.CoordinatesToCellName(cellIdx+1,НомерСтроки)
				if err != nil {
					Ошибка(" %+v ", err)
				}
				f.SetCellValue("Sheet1", cellName, cellData)
			}
			НомерСтроки++
		}
	}
	dt := time.Now()
	dt.Format("01-02-2006 15:04:05")
	if err := f.SaveAs("сессии "+dt.String()+".xlsx"); err != nil {
		Ошибка (err.Error())
	}


	//f := excelize.NewFile()
	//// Создать новый лист
	////index := f.NewSheet("Sheet1")
	//
	//f.SetCellValue("Sheet1", "B2", 100)
	//
	//for ospCode,


	//сообщениеКлиенту = Сообщение{
	//	Id:      0,
	//	От:      "io",
	//	Кому:    client.Login,
	//	Текст:   "Данные собраны",
	//	MessageType: []string{"attention"},
	//	Content: struct {
	//		Target string `json:"target"`
	//		Data interface{} `json:"data"`
	//		Html string `json:"html"`
	//		Обработчик string `json:"обработчик"`
	//	}{
	//		Target:"table_wrapper_"+mes.Выполнить.Arg.Module,
	//		//Html:string(html),
	//	},
	//}
	//client.Message<-&сообщениеКлиенту
}

func СоздатьEXCEL(ДанныеТаблиц map[string]interface{}){
	//Данные := ДанныеТаблиц["Данные"].(map[int]map[string]РезультатSQLизАИС)
	//Данные[КодОСП][ИмяСтолбца]=РезультатЗапроса
	Таблицы := ДанныеТаблиц["Таблицы"].(map[string]map[string]СтруктураСтолбца)
	//Таблицы[ИмяТаблицы][ИмяСТолбца]=СтруктураСтолбцов{
	//	ИмяПеременной: sqlData.ColName,
	//	ИмяСтолбца: sqlData.Labels,
	//	НомерСтолбца: Порядок,
	//	ТипРезультата: sqlData.ResultType,
	//}
	//if len(Таблицы)>0{
	//
	//}
	for ИмяТаблицы, Структура := range Таблицы  {
		Инфо("ИмяТаблицы %+v Структура %+v ", ИмяТаблицы, Структура)

	}
}

func (mes *Сообщение) ПолучитьSQLДляСбораДанныхИзРДБ() ([]sqlStruct, map[string]СтруктураТаблицыДанных){
	/*
	   NOTE: необходимо найти спсоб как передавать массив значений в качестве одного параметра - В КАЧЕСТВЕ СПИСКА
	*/
	//sqlSting := "SELECT sql.table, sql.sql, col_name, label, sql.type, sub_table, dbs, current_timestamp as date_sync  FROM fssp_data.sql WHERE sql.table in ($1) AND sql.sql IS NOT NULL"

	//tablesNames:=[][]byte{}
	СтрокаТаблиц :=""
	СтрокаСтолбца :=""
	//tablesArray:=[]string{}
	//Инфо("Получаем скрипты для таблиц %+v\n", mes.Выполнить.Arg.Tables)
	if mes.Выполнить.Arg.Tables!= nil{
		for _, v := range mes.Выполнить.Arg.Tables{
			//tablesNames=append(tablesNames, []byte(v))
			if СтрокаТаблиц != ""{
				СтрокаТаблиц=СтрокаТаблиц+", "
			}
			СтрокаТаблиц=СтрокаТаблиц+"'"+v+"'"
			//tablesArray=append(tablesArray, "'"+v+"'")
		}
	}

	ВходящиеАргументы := map[string]interface{}{}
	for _, АргументыДействия := range mes.Выполнить.Действие {
		ВходящиеАргументы = АргументыДействия
		if Таблицы, есть := АргументыДействия["table[]"]; есть {
			Инфо("Таблицы %+v", Таблицы)
			for _, таблица := range Таблицы.([]interface{}) {
				Инфо(" %+v", таблица)
				if СтрокаТаблиц != "" {
					СтрокаТаблиц = СтрокаТаблиц + ", "
				}
				СтрокаТаблиц = СтрокаТаблиц + "'" + таблица.(string) + "'"
			}
		}
		if Столбец, есть := АргументыДействия["столбец"];есть{
			СтрокаСтолбца =" AND col_name = '"+Столбец.(string)+"'"
			//for _, столбец := range Столбцы.(interface{}) {
			//	Инфо(" %+v", столбец)
			//	if СтрокаСтолбцов != "" {
			//		СтрокаСтолбцов = СтрокаСтолбцов + ", "
			//	}
			//	СтрокаСтолбцов = СтрокаСтолбцов + "'" + столбец.(string) + "'"
			//}
		}
	}
	//if СтрокаСтолбцов != ""{
	//	СтрокаСтолбцов = " AND col_name in ("+СтрокаСтолбцов+")"
	//}



	Инфо(" %+v", СтрокаТаблиц)
	//stringByte := strings.Join(tablesArray, ",\x20") // x20 = space and x00 = null
	//
	//log.Println([]byte(stringByte))

	//Инфо("stringByte %+v \n", stringByte, []byte(stringByte))

	sqlSting := `SELECT sql.table, sql.sql, col_name, label, sql.type, sub_table, dbs, current_timestamp as date_sync, аргументы, real_time, col_parameters FROM fssp_data.sql WHERE sql.table in (`+СтрокаТаблиц+`) AND sql.sql IS NOT NULL AND dbs <@ '["rbd"]' `+СтрокаСтолбца

	sqlQuery := sqlStruct{
		Name:   "Запрос SQL для сбора данных из АИС",
		Sql:    sqlSting,
		//Values:	[][]byte{[]byte(stringByte)},
		DBSchema:"fssp_data",
	}

	resultRowsArrayMap, _ := ВыполнитьPgSQL(sqlQuery)

	SQLresult := make([]sqlStruct, len(resultRowsArrayMap))
	Таблицы :=  map[string]СтруктураТаблицыДанных{}



	for i, row:= range resultRowsArrayMap{

		var	АргументыДляПодстановки []interface{}
		var	ПараметрыСтолбца map[string]interface{}

		if  row["аргументы"] != nil{
			АргументыДляПодстановки = row["аргументы"].([]interface{})
		}



		if row["col_parameters"] != nil{
			ПараметрыСтолбца = row["col_parameters"].(map[string]interface{})
		}

		Инфо("row[real_time] %+s", row["real_time"])


		ОбновлениеВРеальномВремени := false
		if  row["real_time"]!= nil && row["real_time"].(string) == "t"{
			ОбновлениеВРеальномВремени = true
		}

		SQLresult[i] = sqlStruct{
			ВходящиеАргументы: ВходящиеАргументы,
			SQLАргументы: АргументыДляПодстановки,
			Name:     row["table"].(string),
			Table:    row["table"].(string),
			Sql:      row["sql"].(string),
			ColName:  row["col_name"].(string),
			Labels:   row["label"].(string),
			Type:     []string{},
			SubTable: row["sub_table"].(string),
			Dbs:      row["dbs"].([]interface{}),
			Values:   nil,
			DateSync: row["date_sync"].(string),
			ОбновлениеВРеальномВремени: ОбновлениеВРеальномВремени,
			ПараметрыСтолбца: ПараметрыСтолбца,
		}

		if _, ok := Таблицы[row["table"].(string)];!ok{
			//Таблицы[sqlData.Table]=map[string]СтруктураСтолбцов{}
			Таблицы[row["table"].(string)]= СтруктураТаблицыДанных{
				ИмяТаблицы: row["table"].(string),
				СтруктураСтолбцов: map[string]СтруктураСтолбца{},
				ОбновлятьВРеальномВремени: ОбновлениеВРеальномВремени,
				ВходящиеАргументы: ВходящиеАргументы,
			}
		}

		var НомерСтолбца int
		if номерСтолбца, Есть := ПараметрыСтолбца["order"]; Есть {
			НомерСтолбца = номерСтолбца.(int)
		}

		Таблицы[row["table"].(string)].СтруктураСтолбцов[row["col_name"].(string)]=СтруктураСтолбца{
			ИмяПеременной: row["col_name"].(string),
			ИмяСтолбца: row["label"].(string),
			ОбновлениеВРеальномВремени: ОбновлениеВРеальномВремени,
			НомерСтолбца:НомерСтолбца,
			ПараметрыСтолбца: ПараметрыСтолбца,
			Запрос: row["sql"].(string),
			Аргументы: АргументыДляПодстановки,
		}
	}


	return SQLresult, Таблицы

}

func (mes *Сообщение) ПолучитьSQLДляСбораДанныхИзОСП() ([]sqlStruct, map[string]СтруктураТаблицыДанных){
/*
NOTE: необходимо найти спсоб как передавать массив значений в качестве одного параметра
*/
	//sqlSting := "SELECT sql.table, sql.sql, col_name, label, sql.type, sub_table, dbs, current_timestamp as date_sync  FROM fssp_data.sql WHERE sql.table in ($1) AND sql.sql IS NOT NULL"

	//tablesNames:=[][]byte{}
	СтрокаТаблиц :=""
	//tablesArray:=[]string{}
	//Инфо("Получаем скрипты для таблиц %+v\n", mes.Выполнить.Arg.Tables)
	if mes.Выполнить.Arg.Tables!= nil{
		for _, v := range mes.Выполнить.Arg.Tables{
			//tablesNames=append(tablesNames, []byte(v))
			if СтрокаТаблиц != ""{
				СтрокаТаблиц=СтрокаТаблиц+", "
			}
			СтрокаТаблиц=СтрокаТаблиц+"'"+v+"'"
			//tablesArray=append(tablesArray, "'"+v+"'")
		}
	}

	ВходящиеАргументы := map[string]interface{}{}
	for _, АргументыДействия := range mes.Выполнить.Действие {
		ВходящиеАргументы = АргументыДействия
		if Таблицы, есть := АргументыДействия["table[]"]; есть {
			Инфо("Таблицы %+v", Таблицы)
			for _, таблица := range Таблицы.([]interface{}) {
				Инфо(" %+v", таблица)
				if СтрокаТаблиц != "" {
					СтрокаТаблиц = СтрокаТаблиц + ", "
				}
				СтрокаТаблиц = СтрокаТаблиц + "'" + таблица.(string) + "'"

			}
		}
	}
	Инфо(" %+v", СтрокаТаблиц)
	//stringByte := strings.Join(tablesArray, ",\x20") // x20 = space and x00 = null
	//
	//log.Println([]byte(stringByte))

//Инфо("stringByte %+v \n", stringByte, []byte(stringByte))

	sqlSting := `SELECT sql.table, sql.sql, col_name, label, sql.type, sub_table, dbs, current_timestamp as date_sync, аргументы, real_time, col_parameters FROM fssp_data.sql WHERE sql.table in (`+СтрокаТаблиц+`) AND sql.sql IS NOT NULL AND dbs <@ '["osp"]'`

	sqlQuery := sqlStruct{
		Name:   "Запрос SQL для сбора данных из АИС",
		Sql:    sqlSting,
		//Values:	[][]byte{[]byte(stringByte)},
		DBSchema:"fssp_data",
	}

	resultRowsArrayMap, _ := ВыполнитьPgSQL(sqlQuery)

	SQLresult := make([]sqlStruct, len(resultRowsArrayMap))
	Таблицы :=  map[string]СтруктураТаблицыДанных{}



	for i, row:= range resultRowsArrayMap{

	var	АргументыДляПодстановки []interface{}
	var	ПараметрыСтолбца map[string]interface{}

		if  row["аргументы"] != nil{
			АргументыДляПодстановки = row["аргументы"].([]interface{})
		}



		if row["col_parameters"] != nil{
			ПараметрыСтолбца = row["col_parameters"].(map[string]interface{})
		}

		Инфо("row[real_time] %+s", row["real_time"])


		ОбновлениеВРеальномВремени := false
		if  row["real_time"]!= nil && row["real_time"].(string) == "t"{
			ОбновлениеВРеальномВремени = true
		}

		SQLresult[i] = sqlStruct{
			ВходящиеАргументы: ВходящиеАргументы,
			SQLАргументы: АргументыДляПодстановки,
			Name:     row["table"].(string),
			Table:    row["table"].(string),
			Sql:      row["sql"].(string),
			ColName:  row["col_name"].(string),
			Labels:   row["label"].(string),
			Type:     []string{},
			SubTable: row["sub_table"].(string),
			Dbs:      row["dbs"].([]interface{}),
			Values:   nil,
			DateSync: row["date_sync"].(string),
			ОбновлениеВРеальномВремени: ОбновлениеВРеальномВремени,
			ПараметрыСтолбца: ПараметрыСтолбца,
		}

		if _, ok := Таблицы[row["table"].(string)];!ok{
			//Таблицы[sqlData.Table]=map[string]СтруктураСтолбцов{}
			Таблицы[row["table"].(string)]= СтруктураТаблицыДанных{
				ИмяТаблицы: row["table"].(string),
				СтруктураСтолбцов: map[string]СтруктураСтолбца{},
				ОбновлятьВРеальномВремени: ОбновлениеВРеальномВремени,
				ВходящиеАргументы: ВходящиеАргументы,
			}
		}

		var НомерСтолбца int
		if номерСтолбца, Есть := ПараметрыСтолбца["order"]; Есть {
			НомерСтолбца = номерСтолбца.(int)
		}

		Таблицы[row["table"].(string)].СтруктураСтолбцов[row["col_name"].(string)]=СтруктураСтолбца{
			ИмяПеременной: row["col_name"].(string),
			ИмяСтолбца: row["label"].(string),
			ОбновлениеВРеальномВремени: ОбновлениеВРеальномВремени,
			НомерСтолбца:НомерСтолбца,
			ПараметрыСтолбца: ПараметрыСтолбца,
			Запрос: row["sql"].(string),
			Аргументы: АргументыДляПодстановки,
		}
	}


	return SQLresult, Таблицы

}

/*
ВыполнитьSQL: Выполняет запрос к БД.
Возвращает мап с ключами равными названию столбцов, и значениями соответсвующими столбцам. Карта создаёться на лету.
sqlStingCaption := "SELECT caption FROM meta_table WHERE meta_table.table = $1"
	queryCaption := sqlStruct{
		Name:   tableName,
		Sql:    sqlStingCaption,
		Values: [][]byte{
			[]byte(tableName),
		},
	}
	CaptionRow, _ := ВыполнитьPgSQL(queryCaption)
*/

//func ПолучитьЗапрос (ctx context.Context, название string)( map[string]interface{} , error){
//	var cancel context.CancelFunc
//	if ctx == nil{
//		ctx, cancel = context.WithCancel(context.Background())
//		defer cancel()
//	}
//
//	ПолучитьЗапрос ,err:= sqlStruct{
//			Name:   "sql_запросы",
//			Sql:    "SELECT * FROM iobot.sql_запросы WHERE название = $1",
//			Values: [][]byte{
//				[]byte(название),
//			},
//		}.Выполнить(nil)
//	if err != nil{
//		Инфо(">>>> Ошибка SQL запроса: %+v \n\n",err)
//		return nil, err
//	}
//	if len(ПолучитьЗапрос)>0 {
//		return ПолучитьЗапрос[0], nil
//	} else {
//		return nil, err
//	}
//
//}

func ВыполнитьPgSQL(sqlData sqlStruct) ([]map[string]interface{} , [][]string) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	Read, err := sqlData.PgSqlResultReader(ctx)
	if err != nil {
		Ошибка("\n !! ERR %+v\n", err)
		return nil, nil
	}
	//Rreaded := ResultReader.Read()
	ОписаниеПолей := Read.FieldDescriptions
	//if Rreaded.Err != nil {
	//	Инфо("\n Rreaded.Err >> %+v\n", Rreaded.Err.Error())
	//}

	rows := Read.Rows
	if len(rows) <1 {
		return nil,nil
	}
	МассивКартСтрок := []map[string]interface{}{}
	МассивСтрок := [][]string{}
	Строка := []string{}

	for _,row := range rows {
		КартаСтрок := map[string]interface{}{}
		for collIdx ,ячейка :=range row {

			var ЗначениеЯчейки interface{}
			ИмяПоля := string(ОписаниеПолей[collIdx].Name)
			ЗначениеЯчейки = string(ячейка)

			if ОписаниеПолей[collIdx].DataTypeOID == 3802{ // JSONB
				if len(ячейка)>0{
					var jsonCell interface{}
					err:=json.Unmarshal(ячейка, &jsonCell)

					if err != nil {
						Ошибка("\n !! ERR %+v Поле %+v string(ячейка) %+v\n", err, string(ОписаниеПолей[collIdx].Name), ячейка, len(ячейка))
						КартаСтрок[ИмяПоля] = string(ячейка)
					}
					КартаСтрок[ИмяПоля]=jsonCell
				} else {
					КартаСтрок[ИмяПоля]=nil
				}

			} else {
				КартаСтрок[ИмяПоля]=ЗначениеЯчейки
			}
			Строка =append(Строка, ЗначениеЯчейки.(string))

		}
		МассивСтрок=append(МассивСтрок, Строка)
		МассивКартСтрок= append(МассивКартСтрок, КартаСтрок)
	}
	return МассивКартСтрок, МассивСтрок
}

type РезультатSQL struct {
	Столбцы []string
	Строки [][]interface{}
	КартаСтрок []map[string]string
}

func (sqlData sqlStruct) ВыполнитьSQLвАИС () ( map[string]РезультатSQL, error) {

	Данные := РезультатSQL{}
	ТаблицаДанных := map[string]РезультатSQL{
		sqlData.Name: Данные,
	}
	for _, dbConnector := range sqlData.Dbs{
		//Инфо("dbConnector %+v\n", dbConnector)
		var conn *sql.DB
		switch dbConnector {
		case "rdb":
			conn = RDBconnect() //244
			break
		case "fssp":
			conn = FSSPconnect() //243
			break
		case "osp":
			//conn = FSSPconnect()
			continue
			break
		}
		defer conn.Close()

		РезультатЗапроса, err := conn.Query(sqlData.Sql)
		if err != nil {
			Ошибка("err	 %+v sql.Sql %+v \n", err, sqlData.Sql)
			return nil, err
		}
		if РезультатЗапроса == nil{
			Инфо("Запрос '%+v' не вернул данных %+v\n",sqlData.Name, РезультатЗапроса)
			return ТаблицаДанных, nil
		}
		Данные.Столбцы, err = РезультатЗапроса.Columns()
		if err != nil {
			Ошибка(" не удалось получить список имён колонок	 %+v\n", err)
		}
		Инфо("Данные.Столбцы %+v\n", Данные.Столбцы)

		for РезультатЗапроса.Next(){

			ЗначенияСтрокиNullString := make([]interface{}, len(Данные.Столбцы))
			for i, _ := range Данные.Столбцы {
				ЗначенияСтрокиNullString[i] = new(sql.NullString)
			}


			errScan := РезультатЗапроса.Scan(ЗначенияСтрокиNullString...)
			if errScan != nil {
				Ошибка("err	 %+v a %+v\n", errScan, ЗначенияСтрокиNullString)
			}
			ЗначенияСтроки := make([]interface{}, len(ЗначенияСтрокиNullString))
			КартаСтроки :=map[string]string{}
			for i, Значение := range ЗначенияСтрокиNullString{
				ЗначенияСтроки[i] = Значение.(*sql.NullString).String
				КартаСтроки[Данные.Столбцы[i]] = Значение.(*sql.NullString).String
			}
			Данные.Строки=append(Данные.Строки, ЗначенияСтроки)
			Данные.КартаСтрок=append(Данные.КартаСтрок, КартаСтроки)
		}


		ТаблицаДанных[sqlData.Name] = Данные
	}
return ТаблицаДанных, nil
}

func (sqlData sqlStruct) ВыполнитьSQLвАИСиСохранитьДанные (result chan<- string)  {

	sqlResults := []interface{}{}
	sqlResultMap := map[string][]map[string]interface{}{
		sqlData.Name: {},
	}
	//ch1 := make(chan string)

	for _, dbConnector := range sqlData.Dbs{
		//Инфо("dbConnector %+v\n", dbConnector)
		var conn *sql.DB
		switch dbConnector {
		case "rdb":
			conn = RDBconnect()
			break
		case "fssp":
			conn = FSSPconnect()
			break
		case "osp":
			//conn = FSSPconnect()
			continue
			break
		}
		defer conn.Close()
		//go func() {

		/**/
		//Инфо("sqlData.Sql %+v\n", sqlData.Sql)


		resultSql, err := conn.Query(sqlData.Sql)
		if err != nil {
			Ошибка("err	 %+v sql.Sql %+v \n", err, sqlData.Sql)
		}
		if resultSql == nil{
			Инфо("resultSql %+v\n", resultSql)
			result <- "end"
			return
		}


		resultMaps := map[string]interface{}{}
		savedValuesStrings := []string{}

		savedRows := [][][]byte{}
		savedRowsStrings:= [][]string{}

		for resultSql.Next(){
			cols, err:=resultSql.Columns()
			if err != nil {
				Ошибка("err	 %+v\n", err)
			}
			Инфо("cols %+v\n", cols)
			vals := make([]interface{}, len(cols))
			for i, _ := range cols {
				vals[i] = new(sql.NullString)
			}

			savedValues := [][]byte{}

			//vals := make([]interface{}, len(cols))
			Инфо("len vals %+v len(cols) %+v cols %+v\n", len(vals), len(cols), cols)
			for i, _ := range cols {
				vals[i] = new(sql.NullString)
			}
			err = resultSql.Scan(vals...)

			if err != nil {
				Ошибка("err	 %+v a %+v\n", err, vals)
			}


			for i, colName := range cols {
				//resultMaps[colName] = vals[i].(*sql.NullString).String
				Инфо("vals[i] %+v\n", vals[i].(*sql.NullString).String)

				resultMaps[colName] =vals[i].(*sql.NullString).String
				savedValues=append(savedValues, []byte(vals[i].(*sql.NullString).String))
				savedValuesStrings=append(savedValuesStrings, vals[i].(*sql.NullString).String)

			}

			savedRowsStrings=append(savedRowsStrings, savedValuesStrings)
			savedRows = append(savedRows, savedValues)
			//Инфо("vals %+v\n", vals)

			sqlResults= append(sqlResults, vals)
			sqlResultMap[sqlData.Name] = append(sqlResultMap[sqlData.Name], resultMaps)
			//Инфо("vals %+v cols %+v\n",vals, cols )

			//if sqlData.SubTable != ""{
			//	sqlResultMap[sqlData.Name] = append(sqlResultMap[sqlData.Name], map[string]interface{}{"sub_table":[]byte(sqlData.SubTable)})
			//	cols=append(cols, "sub_table")
			//	savedValues = append(savedValues, []byte(sqlData.SubTable))
			//}

			if sqlData.DateSync != ""{
				sqlResultMap[sqlData.Name]=append(sqlResultMap[sqlData.Name], map[string]interface{}{"date_sync":[]byte(sqlData.DateSync)})
				cols=append(cols, "date_sync")
				savedValues = append(savedValues, []byte(sqlData.DateSync))
			}

			//Инфо("cols %+v savedValues %+v   sqlData Name %+v\n", cols, savedValues, sqlData.Name)
			//Инфо("savedRowsStrings %+v\n", savedRowsStrings)

			СохранитьСтатДанные(cols, savedValues, sqlData.Name)

		}
		//ch1 <- "end"
		/***/
		//}()
	}
	result <- "end"
}
//func СохранитьСтатДанные(sqlResultMap map[string][]map[string]interface{}){
func получитьИндексТаблицы (tableName string)  map[string]string {
	indexFieldsQuers := "SELECT column_name FROM information_schema.key_column_usage WHERE table_name = $1" // +tableName
	sqlQueryIndexF := sqlStruct{
		Name:   tableName,
		Sql:    indexFieldsQuers,
		Values: [][]byte{[]byte(tableName)},
	}

	indexFieldsArrayMap,_ := ВыполнитьPgSQL(sqlQueryIndexF)
	indexMapVales := map[string]string{}
	for _, IndexMap := range indexFieldsArrayMap{
		indexMapVales[IndexMap["column_name"].(string)]=IndexMap["column_name"].(string)
	}
	return indexMapVales
}

func получитьИмяИключения (tableName string)  string {
	indexFieldsQuers := "SELECT indexname FROM pg_indexes WHERE tablename  = $1" // +tableName
	sqlQueryIndexF := sqlStruct{
		Name:   tableName,
		Sql:    indexFieldsQuers,
		Values: [][]byte{[]byte(tableName)},
	}

	indexFieldsArrayMap, _ := ВыполнитьPgSQL(sqlQueryIndexF)
	indexMapVales := ""

	for _, IndexMap := range indexFieldsArrayMap{
		indexMapVales =IndexMap["indexname"].(string)
	}
	return indexMapVales
}

func СохранитьСтатДанные(cols []string, values [][]byte, tableName string){

	//for tableName, dataArrays:=range sqlResultMap{
	colsString:=strings.Join(cols,", ")
	valuesString :=""
	DoUpdate:=""

	Инфо("Сохраняем данные в tableName %+v\n", tableName)

	indexMapVales := получитьИндексТаблицы(tableName) // Ключевые  поля для исключения
	имяИсключения := получитьИмяИключения(tableName) //  ИМЯ исключения для update  в случае совпдания полей

	for i, colName := range cols{
		if valuesString !=""{
			valuesString=valuesString+", "
		}
		valuesString = valuesString+"$"+strconv.Itoa(i+1)
		//Инфо("colName %+v\n", colName)

		if indexMapVales[colName] == ""{
			if DoUpdate != ""{
				DoUpdate =DoUpdate+", "
			} else {
				DoUpdate = "DO UPDATE SET "
			}
			DoUpdate=DoUpdate + colName+" = $"+strconv.Itoa(i+1)
		}
	}
	//Инфо("colsString %+v\n", colsString)
	//Инфо("colsString %+v\n", DoUpdate)
	//Инфо("values %+v\n", values)

	sqlSqtring := "INSERT INTO "+tableName+" ("+colsString+") VALUES ("+valuesString+") ON CONFLICT ON CONSTRAINT "+имяИсключения+" "+DoUpdate

	//Инфо("sqlSqtring %+v\n", sqlSqtring)
	//
	//
	sqlQuery := sqlStruct{
		Name:   tableName,
		Sql:    sqlSqtring,
		Values: values,
		DBSchema: "fssp_data",
	}
	//}
	//
	_, _= ВыполнитьPgSQL(sqlQuery)
	//Инфо("resultRowsArray %+v\n", resultRowsArray)
}


