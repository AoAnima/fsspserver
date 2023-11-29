package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"strconv"
	txtTpl "text/template"
)

type РезультатSQLизАИС struct {
	ВходящиеАргументы         map[string]interface{}
	Таблица                   string
	Столбцы                   []string
	Строки                    [][]interface{}
	КартаСтрок                []map[string]string
	МассивСтрок               []map[string]interface{}
	Количество                int
	ОСП                       map[string]interface{}
	Ошибка                    string
	ОбвнолятьВРеальномВремени bool
}

/*
Возвращает результат в карте с ключём равным имени столбца
*/
func (sqlData sqlStruct) ВыполнитьSQLвРДБ(result chan<- map[string]РезультатSQLизАИС) {

	Данные := РезультатSQLизАИС{
		ВходящиеАргументы:         sqlData.ВходящиеАргументы,
		Таблица:                   sqlData.Table,
		ОбвнолятьВРеальномВремени: sqlData.ОбновлениеВРеальномВремени,
	}

	ТаблицаДанных := map[string]РезультатSQLизАИС{
		sqlData.ColName: Данные,
	}

	var conn *sql.DB
	conn = RDBconnect()
	defer conn.Close()

	if conn == nil {
		Ошибка("Данные.Ошибка %+v ТаблицаДанных %+v", Данные.Ошибка, ТаблицаДанных)

		Данные.Ошибка = "Нет связи с РДБ"
		ТаблицаДанных[sqlData.ColName] = Данные
		result <- ТаблицаДанных
		return
	}

	Аргументы := sqlData.SQLАргументы // аргументы которые нужно получить из данных
	var SQLАргументы []interface{}

	if Аргументы != nil {
		SQLАргументы = make([]interface{}, len(Аргументы))
		for индекс, Адрес := range Аргументы {
			tpl, err := txtTpl.New("SqlЗапрос").Funcs(tplFunc()).Parse("{{" + Адрес.(string) + "}}")
			if err != nil {
				Ошибка(" %+v ", err)
			}
			БайтБуферДанных := new(bytes.Buffer)

			err = tpl.Execute(БайтБуферДанных, sqlData.ВходящиеАргументы)
			if err != nil {
				Ошибка(" %+v ", err)
			}
			if БайтБуферДанных.String() == "<no value>" {
				SQLАргументы[индекс] = nil
			} else {
				SQLАргументы[индекс] = БайтБуферДанных.String()
			}
		}
	} else if sqlData.АргументыДляАИС != nil {
		SQLАргументы = sqlData.АргументыДляАИС
	}

	ПодготовленныйЗапрос, err := conn.Prepare(sqlData.Sql)

	if err != nil {

		Данные.Ошибка = "Ошибка подготовки запроса " + sqlData.Sql + err.Error()
		Ошибка("Ошибка подготовки запроса %+v %+v", err.Error(), sqlData.Sql)
		ТаблицаДанных[sqlData.ColName] = Данные
		result <- ТаблицаДанных
		return
	}
	РезультатЗапроса, err := ПодготовленныйЗапрос.Query(SQLАргументы...)

	if err != nil {
		Данные.Ошибка = "Ошибка выполнения запроса " + sqlData.Sql
		Ошибка("Ошибка выполнения запроса %+v %+v sql.Sql %+v err\t %+v \n", sqlData.Name, sqlData.ColName, sqlData.Sql, err)
		ТаблицаДанных[sqlData.ColName] = Данные
		result <- ТаблицаДанных
		return
	}

	if РезультатЗапроса == nil {
		Ошибка("Запрос '%+v' не вернул данных %+v\n", sqlData.Name, РезультатЗапроса)
		result <- ТаблицаДанных
		return
	}

	Данные.Столбцы, err = РезультатЗапроса.Columns()

	if err != nil {
		Ошибка(" не удалось получить список имён колонок	 %+v\n", err)
	}
	//Инфо("Данные.Столбцы %+v\n", Данные.Столбцы)

	for РезультатЗапроса.Next() {

		ЗначенияСтрокиNullString := make([]interface{}, len(Данные.Столбцы))

		for i, _ := range Данные.Столбцы {
			ЗначенияСтрокиNullString[i] = new(sql.NullString)
		}
		errScan := РезультатЗапроса.Scan(ЗначенияСтрокиNullString...)
		if errScan != nil {
			Ошибка("errScan	 %+v a %+v\n", errScan, ЗначенияСтрокиNullString)
			Данные.Ошибка = "Ошибка ЗначенияСтрокиNullString " + errScan.Error()
		}
		ЗначенияСтроки := make([]interface{}, len(ЗначенияСтрокиNullString))
		КартаСтроки := map[string]string{}

		for i, Значение := range ЗначенияСтрокиNullString {
			ЗначенияСтроки[i] = Значение.(*sql.NullString).String
			КартаСтроки[Данные.Столбцы[i]] = Значение.(*sql.NullString).String
		}

		Данные.Строки = append(Данные.Строки, ЗначенияСтроки)
		Данные.КартаСтрок = append(Данные.КартаСтрок, КартаСтроки)
	}

	ТаблицаДанных[sqlData.ColName] = Данные

	result <- ТаблицаДанных

}

func (sqlData sqlStruct) ВыполнитьSQLАИС(result chan<- map[string]РезультатSQLизАИС) {

	Данные := РезультатSQLизАИС{
		ВходящиеАргументы:         sqlData.ВходящиеАргументы,
		Таблица:                   sqlData.Table,
		ОбвнолятьВРеальномВремени: sqlData.ОбновлениеВРеальномВремени,
		Количество:                0,
	}
	if sqlData.ДанныеПодключения != nil {
		Данные.ОСП = sqlData.ДанныеПодключения
	}

	var ИмяКонтейнераДанных string
	if sqlData.ColName != "" {
		ИмяКонтейнераДанных = sqlData.ColName
	} else if sqlData.Name != "" {
		ИмяКонтейнераДанных = sqlData.Name
	}

	ТаблицаДанных := map[string]РезультатSQLизАИС{
		ИмяКонтейнераДанных: Данные,
	}

	//var conn *sql.DB
	conn := sqlData.Conn

	//Инфо(" conn %+v", conn)

	//sqlData.Conn
	if conn == nil {
		Инфо("Данные.ОСП %+v Данные.Ошибка %+v ТаблицаДанных %+v", Данные.ОСП, Данные.Ошибка, ТаблицаДанных)
		ДанныеПокдлюченияОСП := sqlData.ДанныеПодключения
		if sqlData.Клиент != nil {
			СообщениеКлиенту := Сообщение{
				Id:          0,
				От:          "io",
				Кому:        sqlData.Клиент.Login,
				Текст:       "Соединяюсь с  " + ДанныеПокдлюченияОСП["osp_name"].(string),
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
			//Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)
			go СообщениеКлиенту.СохранитьИОтправить(sqlData.Клиент)
		}

		Connection := OSPconnect(ДанныеПокдлюченияОСП["pwd"].(string), ДанныеПокдлюченияОСП["ip_ais"].(string))
		if Connection == nil {
			if sqlData.Клиент != nil {
				СообщениеКлиенту := Сообщение{
					Id:          0,
					От:          "io",
					Кому:        sqlData.Клиент.Login,
					Текст:       "Нет связи с  " + ДанныеПокдлюченияОСП["osp_name"].(string),
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
				//Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)
				go СообщениеКлиенту.СохранитьИОтправить(sqlData.Клиент)
			}

			Данные.Ошибка = "Нет связи с ФСС"
			ТаблицаДанных[ИмяКонтейнераДанных] = Данные
			result <- ТаблицаДанных

			return
		} else {
			conn = Connection
		}
	}

	//Инфо("АргументыДляАИС %+v", sqlData.АргументыДляАИС)
	//Инфо("SQLАргументы %+v", sqlData.SQLАргументы)

	Аргументы := sqlData.SQLАргументы // аргументы которые нужно получить из данных
	var SQLАргументы []interface{}

	if Аргументы != nil {
		SQLАргументы = make([]interface{}, len(Аргументы))
		for индекс, Адрес := range Аргументы {
			tpl, err := txtTpl.New("SqlЗапрос").Funcs(tplFunc()).Parse("{{" + Адрес.(string) + "}}")
			if err != nil {
				Ошибка(" %+v ", err)
			}
			БайтБуферДанных := new(bytes.Buffer)

			err = tpl.Execute(БайтБуферДанных, sqlData.ВходящиеАргументы)
			if err != nil {
				Ошибка(" %+v ", err)
			}
			if БайтБуферДанных.String() == "<no value>" {
				SQLАргументы[индекс] = nil
			} else {
				SQLАргументы[индекс] = БайтБуферДанных.String()
			}
		}
	} else if sqlData.АргументыДляАИС != nil {
		SQLАргументы = sqlData.АргументыДляАИС
	}
	//Инфо("SQLАргументы  %+v", SQLАргументы)
	ПодготовленныйЗапрос, err := conn.Prepare(sqlData.Sql)
	//Инфо(" sqlData.Sql %+v", sqlData.Sql)

	if err != nil {

		Данные.Ошибка = "Ошибка подготовки запроса " + sqlData.Sql + err.Error()
		Ошибка("Ошибка подготовки запроса %+v %+v", err.Error(), sqlData.Sql)
		ТаблицаДанных[ИмяКонтейнераДанных] = Данные
		result <- ТаблицаДанных
		return
	}
	//if sqlData.Клиент != nil{
	//	СообщениеКлиенту := Сообщение{
	//		Id:      0,
	//		От:      "io",
	//		Кому:    sqlData.Клиент.Login,
	//		Текст:   "Нет связи с  "+sqlData.ДанныеПодключения["osp_name"].(string),
	//		MessageType: []string{"attention"},
	//		Content: struct {
	//			Target string `json:"target"`
	//			Data interface{} `json:"data"`
	//			Html string `json:"html"`
	//			Обработчик string `json:"обработчик"`
	//		}{
	//			Target:"progress",
	//		},
	//	}
	//	//Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)
	//	go СообщениеКлиенту.СохранитьИОтправить(sqlData.Клиент)
	//}
	//Инфо(" SQLАргументы для подстановки в запрос %+v", SQLАргументы)
	//если sqlData.ПеребратьАргументы то в цикле вызывать функцию и подставлять аргументы из мссива, при этом
	РезультатЗапроса, err := ПодготовленныйЗапрос.Query(SQLАргументы...)

	//Инфо("sqlData.Sql  %#v РезультатЗапроса %#v \n %+v",sqlData.Sql, РезультатЗапроса, err)

	if err != nil {
		Данные.Ошибка = "Ошибка выполнения запроса " + sqlData.Sql
		Ошибка("Ошибка выполнения запроса %+v ColName %+v sql.Sql %+v err\t %+v \n SQLАргументы %+v", sqlData.Name, sqlData.ColName, sqlData.Sql, err, SQLАргументы)
		ТаблицаДанных[sqlData.ColName] = Данные
		result <- ТаблицаДанных
		return
	}

	if РезультатЗапроса == nil {
		Ошибка("Запрос '%+v' не вернул данных %+v\n", ИмяКонтейнераДанных, РезультатЗапроса)
		result <- ТаблицаДанных
		return
	}

	Данные.Столбцы, err = РезультатЗапроса.Columns()

	if err != nil {
		Ошибка(" не удалось получить список имён колонок	 %+v\n", err)
	}
	//Инфо("Данные.Столбцы %+v\n", Данные.Столбцы)

	//ЗначенияСтрокиNullString := make([]interface{}, len(Данные.Столбцы))

	for РезультатЗапроса.Next() {

		ЗначенияСтрокиNullString := make([]interface{}, len(Данные.Столбцы))

		typs, err := РезультатЗапроса.ColumnTypes()
		//кол, _ :=РезультатЗапроса.Columns()
		//Инфо("кол %+v", кол)
		//Инфо("ColumnTypes %+v", typs)
		for i, s := range typs {
			//Инфо("%#v %#v РезультатЗапроса.ScanType() %#v", i,s, s.ScanType())
			//Инфо("i %+v,s %+v, s.ScanType().String() %+v", i,s,  s.ScanType().String())
			switch s.ScanType().String() {
			case "int64":
				//Инфо("int64 %+v",i, s.ScanType().String())
				ЗначенияСтрокиNullString[i] = new(sql.NullInt64)
			case "decimal.Decimal":
				//Инфо("i Decimal %+v",i, s.ScanType().String())
				ЗначенияСтрокиNullString[i] = new(sql.NullFloat64)
			case "string":
				//Инфо("string %+v",i, s.ScanType().String())
				ЗначенияСтрокиNullString[i] = new(sql.NullString)
			default:
				//Инфо("default %+v",i, s.ScanType().String())
				ЗначенияСтрокиNullString[i] = new(sql.NullString)
			}
		}
		//Инфо("ЗначенияСтрокиNullString %+v", ЗначенияСтрокиNullString)

		//for i, _ := range Данные.Столбцы {
		//	ЗначенияСтрокиNullString[i] = new(sql.NullString)
		//}
		//Инфо(" %+v \n %#v", РезультатЗапроса, РезультатЗапроса)
		errScan := РезультатЗапроса.Scan(ЗначенияСтрокиNullString...)
		if errScan != nil {

			Ошибка("errScan	 %+v ЗначенияСтрокиNullString %#v\n", errScan, ЗначенияСтрокиNullString)
			Данные.Ошибка = "Ошибка ЗначенияСтрокиNullString " + errScan.Error()
			log.Fatal(errScan)
		}
		//Инфо("ЗначенияСтрокиNullString  %+v", ЗначенияСтрокиNullString)
		ЗначенияСтроки := make([]interface{}, len(ЗначенияСтрокиNullString))
		//КартаСтроки :=map[string]string{}
		МассивСтрок := map[string]interface{}{}

		for i, Значение := range ЗначенияСтрокиNullString {

			var ЗначениеВСтроке interface{}
			//Инфо("Столбец  %+v", Данные.Столбцы[i])
			if strings.Contains(strings.ToLower(Данные.Столбцы[i]), "_json") {
				err := json.Unmarshal([]byte(Значение.(*sql.NullString).String), &ЗначениеВСтроке)
				if err != nil {
					Ошибка("\n !! ERR %+v Поле %+v Значение %+v\n", err, Данные.Столбцы[i], Значение.(*sql.NullString).String)
				}
			} else {
				//Инфо("Значение %+v reflect %+v", Значение, reflect.TypeOf(Значение))
				switch Значение.(type) {
				case *sql.NullInt64:
					//Инфо("ЗначениеВСтроке int64 %+v", Значение)
					ЗначениеВСтроке = strconv.Itoa(int(Значение.(*sql.NullInt64).Int64))
				case *sql.NullFloat64:
					//Инфо("ЗначениеВСтроке float64 %+v", Значение)
					ЗначениеВСтроке = strconv.Itoa(int(Значение.(*sql.NullFloat64).Float64))
				case *sql.NullString:
					//Инфо("ЗначениеВСтроке float64 %+v", Значение)
					ЗначениеВСтроке = Значение.(*sql.NullString).String
				default:
					//Инфо("ЗначениеВСтроке default %+v", Значение)
					ЗначениеВСтроке = Значение.(*sql.NullString).String
				}

				//switch s.ScanType().String() {
				//case "int64":
				//	ЗначенияСтрокиNullString[i] = new(sql.NullInt64)
				//case "decimal.Decimal":
				//	ЗначенияСтрокиNullString[i] = new(sql.NullFloat64)
				//case "string":
				//	ЗначенияСтрокиNullString[i] = new(sql.NullString)
				//default:
				//	ЗначенияСтрокиNullString[i] = new(sql.NullString)
				//}
				//Инфо("Значение %+v", Значение)
				//				ЗначениеВСтроке = Значение.(*sql.NullString).String
			}

			ЗначенияСтроки[i] = ЗначениеВСтроке
			МассивСтрок[Данные.Столбцы[i]] = ЗначениеВСтроке
		}
		if COUNT, ЕстьCOUNT := МассивСтрок["COUNT"]; ЕстьCOUNT {
			Данные.Количество, err = strconv.Atoi(COUNT.(string))
			if err != nil {
				Ошибка("Данные.Количество %+v COUNT %+v", err, COUNT)
			}
		} else {
			Данные.Количество++
		}

		Данные.Строки = append(Данные.Строки, ЗначенияСтроки)
		Данные.МассивСтрок = append(Данные.МассивСтрок, МассивСтрок)
	}
	//Инфо("Данные %+v", Данные)
	ТаблицаДанных[ИмяКонтейнераДанных] = Данные

	//Инфо("ТаблицаДанных  %+v", ТаблицаДанных)

	result <- ТаблицаДанных
	defer conn.Close()
}

/**
ВыполнитьSQLвОСП
1 запрос возвращает значение  ячейки/столбца в карте с ключём равным имени столбца запроса sql.col_name
Стат данные в разрезе ОСП лучше и быстрее собирать в ОСП.
Только при отсутсвии связ с осп можно отправлять запрос в РБД - региональную базу данных
*/

func (sqlData sqlStruct) ВыполнитьSQLвОСП(ДанныеПокдлюченияОСП map[string]interface{}, result chan<- map[string]РезультатSQLизАИС) {

	pwd := ДанныеПокдлюченияОСП["pwd"].(string)
	ip := ДанныеПокдлюченияОСП["ip_ais"].(string)

	Данные := РезультатSQLизАИС{
		ОСП:                       ДанныеПокдлюченияОСП,
		Таблица:                   sqlData.Table,
		Количество:                0,
		ОбвнолятьВРеальномВремени: sqlData.ОбновлениеВРеальномВремени,
	}

	ТаблицаДанных := map[string]РезультатSQLизАИС{
		sqlData.ColName: Данные,
	}

	//Инфо("Собираем данные из  %+v", ДанныеПокдлюченияОСП["osp_name"])

	var conn *sql.DB

	conn = OSPconnect(pwd, ip)

	defer conn.Close()

	if conn == nil {
		Ошибка("нет подключения 	 %+v к %+v\n", conn, ДанныеПокдлюченияОСП["osp_name"])

		//алгоритм  если к ОСП нет подключения / нет связи то подключаемся к РДБ куда стекаються все данные с края, но тогда нужно в sql ещё добавить строку чтобы получать данные касательно этого ОСП

		Данные.Ошибка = "Нет связи с " + ДанныеПокдлюченияОСП["osp_name"].(string)
		Ошибка("Данные.Ошибка %+v ТаблицаДанных %+v", Данные.Ошибка, ТаблицаДанных)
		ТаблицаДанных[sqlData.ColName] = Данные
		result <- ТаблицаДанных
		return
	}

	Аргументы := sqlData.SQLАргументы // аргументы которые нужно получить из данных
	var SQLАргументы []interface{}

	if Аргументы != nil {
		SQLАргументы = make([]interface{}, len(Аргументы))
		for индекс, Адрес := range Аргументы {
			tpl, err := txtTpl.New("SqlЗапрос").Funcs(tplFunc()).Parse("{{" + Адрес.(string) + "}}")
			if err != nil {
				Ошибка(" %+v ", err)
			}
			БайтБуферДанных := new(bytes.Buffer)

			err = tpl.Execute(БайтБуферДанных, sqlData.ВходящиеАргументы)
			if err != nil {
				Ошибка(" %+v ", err)
			}
			//Инфо("БайтБуферДанных %s %+v", БайтБуферДанных,БайтБуферДанных.String()== "<no value>")
			if БайтБуферДанных.String() == "<no value>" {
				SQLАргументы[индекс] = nil
			} else {
				SQLАргументы[индекс] = БайтБуферДанных.String()
			}

		}
	} else if sqlData.АргументыДляАИС != nil {
		SQLАргументы = sqlData.АргументыДляАИС
	}

	//Инфо("SQLАргументы %+s", SQLАргументы)

	//РезультатЗапроса, err := conn.Query(sqlData.Sql)
	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()
	ПодготовленныйЗапрос, err := conn.Prepare(sqlData.Sql)

	if err != nil {

		Данные.Ошибка = "Ошибка подготовки запроса: " + err.Error() + "; " + sqlData.Sql
		Ошибка("Ошибка подготовки запроса %+v %+v", err.Error(), sqlData.Sql)
		ТаблицаДанных[sqlData.ColName] = Данные
		result <- ТаблицаДанных
		return
	}
	РезультатЗапроса, err := ПодготовленныйЗапрос.Query(SQLАргументы...)

	if err != nil {

		Данные.Ошибка = "Ошибка выполнения запроса " + sqlData.Sql

		Ошибка("%+v %+v sql.Sql %+v err\t %+v \n", sqlData.Name, sqlData.ColName, sqlData.Sql, err)
		ТаблицаДанных[sqlData.ColName] = Данные
		result <- ТаблицаДанных
		return
	}

	if РезультатЗапроса == nil {
		Инфо("Запрос '%+v' не вернул данных %+v\n", sqlData.Name, РезультатЗапроса)
		result <- ТаблицаДанных

		return
	}

	Данные.Столбцы, err = РезультатЗапроса.Columns()

	if err != nil {
		Ошибка(" не удалось получить список имён колонок	 %+v\n", err)
	}
	//Инфо("Данные.Столбцы %+v\n", Данные.Столбцы)

	for РезультатЗапроса.Next() {

		ЗначенияСтрокиNullString := make([]interface{}, len(Данные.Столбцы))
		for i, _ := range Данные.Столбцы {
			ЗначенияСтрокиNullString[i] = new(sql.NullString)
		}

		errScan := РезультатЗапроса.Scan(ЗначенияСтрокиNullString...)
		if errScan != nil {
			Ошибка("err	 %+v a %+v\n", errScan, ЗначенияСтрокиNullString)
			Данные.Ошибка = "Ошибка ЗначенияСтрокиNullString " + errScan.Error()
		}

		ЗначенияСтроки := make([]interface{}, len(ЗначенияСтрокиNullString))
		КартаСтроки := map[string]string{}

		for i, Значение := range ЗначенияСтрокиNullString {
			ЗначенияСтроки[i] = Значение.(*sql.NullString).String
			КартаСтроки[Данные.Столбцы[i]] = Значение.(*sql.NullString).String
		}

		if COUNT, ЕстьCOUNT := КартаСтроки["COUNT"]; ЕстьCOUNT {
			Данные.Количество, err = strconv.Atoi(COUNT)
			if err != nil {
				Ошибка("Данные.Количество %+v COUNT %+v", err, COUNT)
			}
		} else {
			Данные.Количество++
		}

		Данные.Строки = append(Данные.Строки, ЗначенияСтроки)
		Данные.КартаСтрок = append(Данные.КартаСтрок, КартаСтроки)
	}

	ТаблицаДанных[sqlData.ColName] = Данные

	result <- ТаблицаДанных

}

func (client *Client) ДобавитьСпецФильтр_2(mes Сообщение) { //map[string]interface{}
	ДанныеПокдлюченийКОсп, err := ПолучитьДанныеПодключенийАИСОСП()
	if err != nil {
		Ошибка(" %+v ", err)
	}

	//СпецФильтрДляДобавления ,ОшибкаЗапроса := sqlStruct{
	//	Name:   "СпецФильтр",
	//	Sql:    Скрипт.(string),
	//	Values: АргументыSQL,
	//}.Выполнить(nil)

	//495-721-35-05 -- настройка киско
	//	АргументыЗапроса := []interface{}{
	//		"DOC_DEPOSIT",
	//		"ЕПГУ, перечисление не превышает 7 дней где сумма долга ИП = сумма ПП РАБОЧИЙ 1",
	//		" join (\n    select distinct dd.id from doc_deposit dd join document d on d.id = dd.id  and d.docstatusid = 70\n                                   JOIN DOC_DEPOSIT_IP ON DOC_DEPOSIT_IP.ID = dd.id\n                                   join doc_ip on DOC_DEPOSIT_IP.IP_ID = doc_ip.id\n                                   join o_ip_rm_regfk zkr on zkr.paydoc_id = dd.id -- Связь ПП с ЗКР\n                                   join o_ip_rm_bankstlist kredit on kredit.paydoc_id = zkr.id -- связь ЗКР с кредитовой строкой\n    where kredit.payment_date <= (\n        select first 1 skip 7 dates_date from dates\n    where\n        dates.dates_date >= d.doc_date\n      and dates.dates_isbusiness = 1\n    order by dates_date\n)  and amount = ID_DEBTSUM )  t on t.id = doc_deposit.id",
	//		"left(doc_deposit.doc_deposit_contracc, 4)= '3023' and upper(doc_deposit.doc_deposit_memo) not similar to '%ПРИЛОЖЕНИЕ%|%МЕРЧАНТ%'",
	//		nil,
	//	}
	АргументыЗапроса := [][]interface{}{
		{
			"DOC_DEPOSIT",
			"Оплата через ЕПГУ, ИП окончено фактом более чем через 7 дней.",
			"join(\nselect \n doc_deposit.id, \n doc_ip.point,\n  doc_ip.article,\n   doc_ip.subpoint,\n   document.amount,\n   doc_ip.id_debtsum,\n  DOC_DEPOSIT.CARRY_DATE,\n   doc_ip.ip_date_finish\n FROM O_IP_RM_PLPDBT O_IP_RM_PLPDBT\nJOIN DOC_DEPOSIT_IP ON O_IP_RM_PLPDBT.ID=DOC_DEPOSIT_IP.ID\nJOIN DOC_DEPOSIT ON DOC_DEPOSIT_IP.ID=DOC_DEPOSIT.ID\nJOIN DOCUMENT ON DOC_DEPOSIT.ID=DOCUMENT.ID\njoin doc_ip_doc on doc_deposit_ip.ip_id = doc_ip_doc.id\nJOIN DOC_IP ON DOC_IP_DOC.ID=DOC_IP.ID JOIN DOCUMENT d ON DOC_IP.ID=d.ID\n) t on t.id = doc_deposit.id\n",
			"left(DOC_DEPOSIT.DOC_DEPOSIT_CONTRACC,4)= '3023' AND upper(DOC_DEPOSIT.DOC_DEPOSIT_MEMO) not containing 'ПРИЛОЖЕНИЕ'\nAND upper(DOC_DEPOSIT.DOC_DEPOSIT_MEMO) not containing 'МЕРЧАНТ'\n and t.point  = '1'\n       and t.article = '47'\n       and t.subpoint = '1'\n       and t.amount = t.id_debtsum\n       and t.CARRY_DATE < dateadd(day, -7,  t.ip_date_finish)",
			nil,
		},
		{
			"DOC_DEPOSIT",
			"Оплата через ЕПГУ, ИП неокончено более чем через 7 дней.",
			"join(\nselect \n doc_deposit.id, \n doc_ip.point,\n  doc_ip.article,\n   doc_ip.subpoint,\n   document.amount,\n   doc_ip.id_debtsum,\n  DOC_DEPOSIT.CARRY_DATE,\n   doc_ip.ip_date_finish\n FROM O_IP_RM_PLPDBT O_IP_RM_PLPDBT\nJOIN DOC_DEPOSIT_IP ON O_IP_RM_PLPDBT.ID=DOC_DEPOSIT_IP.ID\nJOIN DOC_DEPOSIT ON DOC_DEPOSIT_IP.ID=DOC_DEPOSIT.ID\nJOIN DOCUMENT ON DOC_DEPOSIT.ID=DOCUMENT.ID\njoin doc_ip_doc on doc_deposit_ip.ip_id = doc_ip_doc.id\nJOIN DOC_IP ON DOC_IP_DOC.ID=DOC_IP.ID JOIN DOCUMENT d ON DOC_IP.ID=d.ID\n) t on t.id = doc_deposit.id",
			"left(DOC_DEPOSIT.DOC_DEPOSIT_CONTRACC,4)= '3023' AND upper(DOC_DEPOSIT.DOC_DEPOSIT_MEMO) not containing 'ПРИЛОЖЕНИЕ'\nAND upper(DOC_DEPOSIT.DOC_DEPOSIT_MEMO) not containing 'МЕРЧАНТ'\nand t.ip_date_finish is null\nand t.amount = t.id_debtsum\nand t.CARRY_DATE < dateadd(day, -7,  current_date)",
			nil,
		},
		{
			"DOC_DEPOSIT",
			"Оплата через ЕПГУ,ИП окончено фактом менее чем через 7 дней.",
			"JOIN (\n  SELECT\n    DOC_DEPOSIT.ID\n  FROM DOC_DEPOSIT DOC_DEPOSIT JOIN DOCUMENT ON DOC_DEPOSIT.ID=DOCUMENT.ID\n  join(\n  select \n   doc_deposit.id, \n   doc_ip.point,\n    doc_ip.article,\n     doc_ip.subpoint,\n     document.amount,\n     doc_ip.id_debtsum,\n    DOC_DEPOSIT.CARRY_DATE,\n     doc_ip.ip_date_finish\n   FROM O_IP_RM_PLPDBT O_IP_RM_PLPDBT\n  JOIN DOC_DEPOSIT_IP ON O_IP_RM_PLPDBT.ID=DOC_DEPOSIT_IP.ID\n  JOIN DOC_DEPOSIT ON DOC_DEPOSIT_IP.ID=DOC_DEPOSIT.ID\n  JOIN DOCUMENT ON DOC_DEPOSIT.ID=DOCUMENT.ID\n  join doc_ip_doc on doc_deposit_ip.ip_id = doc_ip_doc.id\n  JOIN DOC_IP ON DOC_IP_DOC.ID=DOC_IP.ID JOIN DOCUMENT d ON DOC_IP.ID=d.ID\n  ) t on t.id = doc_deposit.id",
			"left(DOC_DEPOSIT.DOC_DEPOSIT_CONTRACC,4)= '3023' AND upper(DOC_DEPOSIT.DOC_DEPOSIT_MEMO) not containing 'ПРИЛОЖЕНИЕ'\n  AND upper(DOC_DEPOSIT.DOC_DEPOSIT_MEMO) not containing 'МЕРЧАНТ'\n   and t.point  = '1'\n         and t.article = '47'\n         and t.subpoint = '1'\n         and t.amount = t.id_debtsum\n         and t.ip_date_finish  <= dateadd(day, 7,  t.CARRY_DATE)",
			nil,
		},
	}

	запущено := 0

	for i, ДанныеОСП := range ДанныеПокдлюченийКОсп {
		go func(ДанныеОСП map[string]interface{}, запущено *int) {

			pwd := ДанныеОСП["pwd"].(string)
			ip := ДанныеОСП["ip_ais"].(string)

			Conn := OSPconnect(pwd, ip)
			if Conn == nil {

				Ошибка("нет связи с Conn %+v  ДанныеОСП %+v ", Conn, ДанныеОСП)

			}
			Инфо(" %+v %+v sql %+v", ДанныеОСП["osp_name"], i, "проверка количества обращений")

			ПодготовленныйЗапросПроверкаИмениФильтра, err := Conn.Prepare(`select ID from CUSTOM_SQL_FILTER where SQL_FILTER_NAME= ?`)
			if err != nil {
				Ошибка("Ошибка подготовки запроса %+v %+v", err.Error())
				return
			}
			РезультатЗапросаSelect, err := ПодготовленныйЗапросПроверкаИмениФильтра.Query("15 дней")
			var IDзапроса interface{}
			if err != nil {
				Ошибка("  %+v", err)
			}
			for РезультатЗапросаSelect.Next() {

				errScan := РезультатЗапросаSelect.Scan(&IDзапроса)
				if errScan != nil {
					Ошибка("errScan	 %+v a %+v\n", errScan, IDзапроса)
					continue
				}

			}
			Инфо("%+v i %+v IDзапроса %+v", ДанныеОСП["osp_name"], i, IDзапроса)

			switch IDзапроса.(type) {
			case int32:
				if IDзапроса.(int32) < 1 {

					return
				}
			case int64:
				if IDзапроса.(int64) < 1 {

					return
				}
			}

			*запущено++

			Инфо(" АргументыЗапроса %+v", АргументыЗапроса)

			if err != nil {
				Ошибка("Ошибка подготовки запроса %+v %+v", err.Error())
				return
			}

			//РезультатЗапроса, err := ПодготовленныйЗапрос.Query()
			//Инфо("РезультатЗапроса %+v err %+v", РезультатЗапроса, err)
			//for РезультатЗапроса.Next() {
			//	var SJOB_ID interface{}
			//	var SJOB_ACTIVE interface{}
			//
			//	errScan := РезультатЗапроса.Scan(&SJOB_ID, &SJOB_ACTIVE)
			//	Инфо("errScan %+v", errScan, SJOB_ID, SJOB_ACTIVE)
			//}
			for _, Аргументы := range АргументыЗапроса {
				ПодготовленныйЗапрос, err := Conn.Prepare(`INSERT INTO CUSTOM_SQL_FILTER (ID, OBJECT_NAME, SQL_FILTER_NAME,SQL_JOIN_PART, SQL_WHERE_PART, PARENT_ID) VALUES ((SELECT NEXT VALUE FOR SEQ_CUSTOM_SQL_FILTER FROM RDB$DATABASE),?,?,?,?,?) RETURNING ID`)
				if err != nil {
					Ошибка("Ошибка подготовки запроса %+v %+v", err.Error())
					return
				}
				РезультатЗапроса, err := ПодготовленныйЗапрос.Query(Аргументы...)
				if err != nil {
					Ошибка(" %+v ", err)
				}
				for РезультатЗапроса.Next() {
					var ID sql.NullInt64
					errScan := РезультатЗапроса.Scan(&ID)
					if errScan != nil {
						Ошибка(" %+v ", errScan)
					}
					*запущено--
					Инфо("ДанныеОСП  %+v ID %+v запущено %+v", ДанныеОСП, ID, &запущено)

				}
			}

		}(ДанныеОСП, &запущено)
	}
	Инфо("запущено %+v", запущено)

	for {
		if запущено <= 0 {
			Инфо("запущено %+v", запущено)
			return
		}
	}

	//return nil
}

func (client *Client) УдалитьСпецФильтр(mes Сообщение) {
	Инфо(" %+v", "УдалитьСпецФильтр")

	ДанныеПокдлюченийКОсп, err := ПолучитьДанныеПодключенийАИСОСП()
	if err != nil {
		Ошибка(" %+v ", err)
	}
	//АргументыЗапроса := [][]interface{}{
	//	{"Оплата через ЕПГУ, ИП неокончено более чем через 7 дней не входит в сводное"},
	//}

	ВходящиеАргументы := mes.Выполнить.Действие["УдалитьСпецФильтр"]["имя_спецфильтра"].([]interface{})

	АргументыЗапроса := [][]interface{}{
		ВходящиеАргументы,
	}

	запущено := 0
	for _, ДанныеОСП := range ДанныеПокдлюченийКОсп {
		go func(ДанныеОСП map[string]interface{}, запущено *int) {

			pwd := ДанныеОСП["pwd"].(string)
			ip := ДанныеОСП["ip_ais"].(string)

			Conn := OSPconnect(pwd, ip)
			if Conn == nil {
				Ошибка("нет связи с Conn %+v  ДанныеОСП %+v ", Conn, ДанныеОСП)
			}

			for _, Аргументы := range АргументыЗапроса {
				ПодготовленныйЗапрос, err := Conn.Prepare(`DELETE FROM CUSTOM_SQL_FILTER WHERE SQL_FILTER_NAME =  ?`)
				if err != nil {
					Ошибка("Ошибка подготовки запроса %+v %+v", err.Error())
					return
				}
				РезультатЗапроса, err := ПодготовленныйЗапрос.Query(Аргументы...)
				if err != nil {
					Ошибка(" %+v ", err)
				}
				for РезультатЗапроса.Next() {
					var ID sql.NullInt64
					errScan := РезультатЗапроса.Scan(&ID)
					if errScan != nil {
						Ошибка(" %+v ", errScan)
					}
					*запущено--
					Инфо("ДанныеОСП  %+v ID %+v запущено %+v", ДанныеОСП, ID, &запущено)
				}
				Инфо("заверщшено %+v ID %+v запущено %+v", ДанныеОСП, &запущено)
			}

		}(ДанныеОСП, &запущено)
	}

}

func (client *Client) ИзменитьСпецФильтр(mes Сообщение) { //map[string]interface{}

	ДанныеПокдлюченийКОсп, err := ПолучитьДанныеПодключенийАИСОСП()
	if err != nil {
		Ошибка(" %+v ", err)
	}

	//СпецФильтрДляДобавления ,ОшибкаЗапроса := sqlStruct{
	//	Name:   "СпецФильтр",
	//	Sql:    Скрипт.(string),
	//	Values: АргументыSQL,
	//}.Выполнить(nil)

	//495-721-35-05 -- настройка киско
	//АргументыЗапроса := []interface{}{
	//	nil,
	//	"document.id in (select id from document where (docstatusid = 105  OR docstatusid = 800) AND (METAOBJECT_CAPTION LIKE 'Жалоба%' OR  METAOBJECT_CAPTION LIKE '%обращен%' OR  METAOBJECT_CAPTION LIKE '%ходат%' OR  METAOBJECT_CAPTION LIKE '%59%' ) )",
	//	"Обращения, жалобы",
	//}

	запущено := 0

	for i, ДанныеОСП := range ДанныеПокдлюченийКОсп {
		//if ДанныеОСП["osp_code"].(string) == "26000" || ДанныеОСП["osp_code"].(string) == "26911" {
		//	continue
		//}
		go func(ДанныеОСП map[string]interface{}, запущено *int) {

			pwd := ДанныеОСП["pwd"].(string)
			ip := ДанныеОСП["ip_ais"].(string)

			Conn := OSPconnect(pwd, ip)
			if Conn == nil {

				Ошибка("нет связи с Conn %+v  ДанныеОСП %+v ", Conn, ДанныеОСП)

				//ОСП, _ := json.Marshal(ДанныеОСП)
				//запрос, _ := json.Marshal() //АргументыЗапроса
				//_ ,ОшибкаЗапроса := sqlStruct{
				//	Name:   "Результат добавления спецфильтра",
				//	Sql:   `INSERT INTO public.статус_добавления_спецфильтров (осп, запрос, статус) values ($1,$2,$3)`,
				//	Values: [][]byte{
				//		ОСП,
				//		запрос,
				//		[]byte("нет связи"),
				//	},
				//	DBSchema:"public",
				//}.Выполнить(nil)
				//
				//if ОшибкаЗапроса != nil {
				//	Ошибка("  %+v", ОшибкаЗапроса)
				//}
			}
			//Инфо(" %+v %+v sql %+v", ДанныеОСП["osp_name"], i, "проверка количества обращений")
			//			ПодготовленныйЗапросПроверкаИмениФильтра, err := Conn.Prepare(`SELECT
			//    count(I.ID)
			//FROM I JOIN DOCUMENT ON I.ID=DOCUMENT.ID
			//WHERE document.docstatusid in (2, 3, 185)
			//  and i.uptodate < current_date
			//  and DOCUMENT.DOC_NUMBER LIKE '%Х'`)
			ПодготовленныйЗапросПроверкаИмениФильтра, err := Conn.Prepare(`select ID from CUSTOM_SQL_FILTER where SQL_FILTER_NAME= ?`)
			if err != nil {
				Ошибка("Ошибка подготовки запроса %+v %+v", err.Error())
				return
			}
			РезультатЗапросаSelect, err := ПодготовленныйЗапросПроверкаИмениФильтра.Query("Оплата через ЕПГУ, ИП неокончено более чем через 7 дней не входит в сводное")
			var IDзапроса interface{}
			if err != nil {
				Ошибка("  %+v", err)
			}
			for РезультатЗапросаSelect.Next() {

				errScan := РезультатЗапросаSelect.Scan(&IDзапроса)
				if errScan != nil {
					Ошибка("errScan	 %+v a %+v\n", errScan, IDзапроса)
					continue
				}

			}
			Инфо("%+v %+v IDзапроса %+v", ДанныеОСП["osp_name"], i, IDзапроса)

			switch IDзапроса.(type) {
			case int32:
				if IDзапроса.(int32) < 1 {

					return
				}
			case int64:
				if IDзапроса.(int64) < 1 {

					return
				}
			}

			//if IDзапроса != nil{
			//	Инфо("ПодготовленныйЗапросПроверкаИмениФильтра %+v РезультатЗапросаSelect ID %+v",ДанныеОСП["osp_name"] , IDзапроса)
			//
			//	continue
			//}
			//ПодготовленныйЗапрос, err := Conn.Prepare(`execute block as
			//declare variable ID D_ID;
			//begin
			//  for select D.ID
			//      from DOCUMENT D
			//      where D.DOC_NUMBER in ('MVV-CUSTOM-0000518')
			//      into :ID
			//  do
			//  begin
			//    delete from DOCUMENT D where D.ID = :ID;
			//    delete from SYS_PATCH SP where SP.ID = :ID;
			//  end
			//end`)
			*запущено++

			ПодготовленныйЗапрос, err := Conn.Prepare(`UPDATE CUSTOM_SQL_FILTER SET SQL_JOIN_PART = ?, SQL_WHERE_PART = ? WHERE ID = ? RETURNING ID`)
			if err != nil {
				Ошибка("Ошибка подготовки запроса %+v %+v", err.Error())
				return
			}
			АргументыЗапроса := []interface{}{
				`join(
					select
					 doc_deposit.id,
					 doc_ip.point,
					  doc_ip.article,
					   doc_ip.subpoint,
					   document.amount,
					   doc_ip.id_debtsum,
					  DOC_DEPOSIT.CARRY_DATE,
					   doc_ip.ip_date_finish,
					   DBTRGRP
					 FROM O_IP_RM_PLPDBT O_IP_RM_PLPDBT
					JOIN DOC_DEPOSIT_IP ON O_IP_RM_PLPDBT.ID=DOC_DEPOSIT_IP.ID
					JOIN DOC_DEPOSIT ON DOC_DEPOSIT_IP.ID=DOC_DEPOSIT.ID
					JOIN DOCUMENT ON DOC_DEPOSIT.ID=DOCUMENT.ID
					join doc_ip_doc on doc_deposit_ip.ip_id = doc_ip_doc.id
					JOIN DOC_IP ON DOC_IP_DOC.ID=DOC_IP.ID JOIN DOCUMENT d ON DOC_IP.ID=d.ID
					
					) t on t.id = doc_deposit.id`,
				`left(DOC_DEPOSIT.DOC_DEPOSIT_CONTRACC,4)= '3023' AND upper(DOC_DEPOSIT.DOC_DEPOSIT_MEMO) not containing 'ПРИЛОЖЕНИЕ'
					AND upper(DOC_DEPOSIT.DOC_DEPOSIT_MEMO) not containing 'МЕРЧАНТ'
					and t.ip_date_finish is null
					and t.amount = t.id_debtsum
					and t.CARRY_DATE < dateadd(day, -7,  current_date) 
					 and DBTRGRP is null`,
				IDзапроса,
			}
			Инфо(" АргументыЗапроса %+v", АргументыЗапроса)
			//ПодготовленныйЗапрос, err := Conn.Prepare(`-- update sys_job set SJOB_ACTIVE =0 where SJOB_ID =-15345 RETURNING SJOB_ID, SJOB_ACTIVE`)
			//	ПодготовленныйЗапрос, err := Conn.Prepare(`-- UPDATE CUSTOM_SQL_FILTER SET SQL_JOIN_PART = 'join (
			//select document.id from
			//  (
			//  select  document.doc_number, count(document.doc_number) count_doc  from sendlist
			//  join document on document.id = sendlist.id where SENDLIST_OUT_DATE BETWEEN cast(:DATE as date) and  cast(:TO as date)
			//  AND lower(SENDLIST_CONTR_TYPE) =lower( :TYPE_RECIPIENT)
			//  group by document.doc_number
			//  ) t
			//  join document on document.DOC_NUMBER = t.DOC_NUMBER
			//WHERE  t.count_doc >1
			//)  s on s.id = sendlist.id', SQL_WHERE_PART = 'lower(SENDLIST_CONTR_TYPE) =lower( :TYPE_RECIPIENT)' WHERE SQL_FILTER_NAME = ? RETURNING ID`)

			//ПодготовленныйЗапрос, err := Conn.Prepare(` SELECT ID FROM CUSTOM_SQL_FILTER WHERE SQL_FILTER_NAME = ?`)
			//ПодготовленныйЗапрос, err := Conn.Prepare(` SELECT CF_ID FROM CUSTOM_SQL_FILTER_EXTRA_FILTERS WHERE CSFEF_PARAM_NAME = ?`)
			//ПодготовленныйЗапрос, err := Conn.Prepare(`DELETE FROM CUSTOM_SQL_FILTER_EXTRA_FILTERS WHERE CSFEF_PARAM_NAME = ?`)
			if err != nil {
				Ошибка("Ошибка подготовки запроса %+v %+v", err.Error())
				return
			}

			//РезультатЗапроса, err := ПодготовленныйЗапрос.Query()
			//Инфо("РезультатЗапроса %+v err %+v", РезультатЗапроса, err)
			//for РезультатЗапроса.Next() {
			//	var SJOB_ID interface{}
			//	var SJOB_ACTIVE interface{}
			//
			//	errScan := РезультатЗапроса.Scan(&SJOB_ID, &SJOB_ACTIVE)
			//	Инфо("errScan %+v", errScan, SJOB_ID, SJOB_ACTIVE)
			//}

			РезультатЗапроса, err := ПодготовленныйЗапрос.Exec(АргументыЗапроса...)
			if err != nil {
				Ошибка(" %+v АргументыЗапроса %+v  ДанныеОСП %+v ", err, АргументыЗапроса, ДанныеОСП)
			}

			Инфо("ДанныеОСП %+v РезультатЗапроса %+v", ДанныеОСП, РезультатЗапроса)
			//		РезультатЗапроса, err := ПодготовленныйЗапрос.Query(
			//				`O`,`Контроль`,
			//				nil,
			//`o.doc_electron=1 and docstatusid in (1,5,180,181,182,183,186,861)`,
			//       nil,
			//		)
			//		РезультатЗапроса, err := ПодготовленныйЗапрос.Query(`Дубликаты по номеру`)
			//		РезультатЗапроса, _ := Conn.Exec(`DELETE FROM CUSTOM_SQL_FILTER_EXTRA_FILTERS WHERE CSFEF_PARAM_NAME = 'DATE'`)
			//		r, e:=РезультатЗапроса.RowsAffected()
			//
			//		Инфо("РезультатЗапроса  %+v %+v %+v", r, e, ДанныеОСП["osp_name"])
			//		РезультатЗапроса, err := ПодготовленныйЗапрос.Query(
			//			`SENDLIST`,`Дубликаты по номеру`,
			//			`join (
			//select document.id from
			//  (select  document.doc_number, count(document.doc_number) count_doc  from sendlist
			//  join document on document.id = sendlist.id where SENDLIST_OUT_DATE >= cast(:DATE as date)
			//  group by document.doc_number) t
			//  join document on document.DOC_NUMBER = t.DOC_NUMBER
			//  where t.count_doc >1
			//) s on s.id = sendlist.id`,nil, nil,)
			//		if err != nil {
			//			Ошибка("  %+v", err)
			//			continue
			//		}

			//for РезультатЗапроса.Next() {
			//	var ID sql.NullInt64
			//	errScan := РезультатЗапроса.Scan(&ID)
			//	if errScan != nil {
			//		Ошибка(" %+v ", errScan)
			//	}
			//	*запущено--
			//	Инфо("ДанныеОСП  %+v ID %+v запущено %+v", ДанныеОСП, ID, *запущено)
			//
			//}

			//		if err != nil {
			//			Ошибка("  %+v", err)
			//
			//			ОСП, _ := json.Marshal(ДанныеОСП)
			//			запрос, _ := json.Marshal(АргументыЗапроса)
			//
			//
			//			_ ,ОшибкаЗапроса := sqlStruct{
			//											Name:   "Результат добавления спецфильтра",
			//											Sql:   `INSERT INTO public.статус_добавления_спецфильтров (осп, запрос, статус) values ($1,$2,$3)`,
			//											Values: [][]byte{
			//												ОСП,
			//												запрос,
			//												[]byte(err.Error()),
			//
			//											},
			//											DBSchema:"public",
			//										}.Выполнить(nil)
			//
			//			if ОшибкаЗапроса != nil {
			//				Ошибка("  %+v", ОшибкаЗапроса)
			//			}
			//
			//		}
			//
			//		for РезультатЗапроса.Next() {
			//			var ID sql.NullInt64
			//			errScan := РезультатЗапроса.Scan(&ID)
			//			Инфо("ID  %+v %+v",strconv.Itoa(int(ID.Int64)) , ДанныеОСП["osp_name"])
			//			if errScan != nil {
			//				Ошибка("errScan	 %+v a %+v\n", errScan, ID)
			//				continue
			//			} else {
			//				ids  := strconv.Itoa(int(ID.Int64))
			//				ОСП, _ := json.Marshal(ДанныеОСП)
			//				запрос, _ := json.Marshal(АргументыЗапроса)
			//				_ ,ОшибкаЗапроса := sqlStruct{
			//					Name:   "Результат добавления спецфильтра",
			//					Sql:   `INSERT INTO public.статус_добавления_спецфильтров (осп, запрос, статус) values ($1,$2,$3)`,
			//					Values: [][]byte{
			//						ОСП,
			//						запрос,
			//						[]byte(ids),
			//					},
			//					DBSchema:"public",
			//				}.Выполнить(nil)
			//
			//				if ОшибкаЗапроса != nil {
			//					Ошибка("  %+v", ОшибкаЗапроса)
			//				}
			//
			//
			//				Поля := [][]interface{}{
			//					{ids, -5, "COUNTDOC","Кличество документов >",1,"20"},
			//
			//					{ids, 3, "DATEFROM","Дата регистрации с",2,"current_date"},
			//					{ids, 4, "DATETO", "по",3,"current_date"},
			//				}
			//
			//		for _, Данные := range Поля {
			//			ПодготовленныйЗапросПоля, err := Conn.Prepare(`INSERT INTO CUSTOM_SQL_FILTER_EXTRA_FILTERS (CF_ID, EXTRA_FILTER_ID, CSFEF_PARAM_NAME, CSFEF_PARAM_CAPTION, CSFEF_PARAM_ORDER,CSFEF_FILTER_DEFAULT) VALUES (?,?,?,?,?,?) RETURNING CF_ID`)
			//			if err != nil {
			//				Ошибка("  %+v", err)
			//			}
			//
			//			_, err = ПодготовленныйЗапросПоля.Query(Данные...)
			//			if err != nil {
			//				Ошибка("  %+v", err)
			//			}
			//		}
			//
			//
			//
			//				//sqlStr := `INSERT INTO CUSTOM_SQL_FILTER_EXTRA_FILTERS (CF_ID, EXTRA_FILTER_ID, CSFEF_PARAM_NAME, CSFEF_PARAM_CAPTION, CSFEF_PARAM_ORDER) VALUES (`+ids+`,4,'TO','Дата по',2) RETURNING CF_ID`
			//
			//				//sqlStr := `INSERT INTO CUSTOM_SQL_FILTER_EXTRA_FILTERS (CF_ID, EXTRA_FILTER_ID, CSFEF_PARAM_NAME, CSFEF_PARAM_CAPTION, CSFEF_PARAM_ORDER) VALUES (`+ids+`,3,'DATE','Дата исходящего с',1) RETURNING CF_ID`
			//				//_, err := Conn.Exec(sqlStr)
			//				//if err != nil {
			//				//	Ошибка("  %+v", err)
			//				//}
			//				//t, _:=r.RowsAffected()
			//				//Инфо(" CF_ID %+v", t)
			////for r.Next(){
			////	var re interface{}
			////	err :=r.Scan(&re)
			////	if err != nil {
			////		Ошибка("  %+v", err)
			////	}
			////	Инфо("  %+v", re)
			////}
			////				sqlStr = `INSERT INTO CUSTOM_SQL_FILTER_EXTRA_FILTERS (CF_ID, EXTRA_FILTER_ID, CSFEF_PARAM_NAME, CSFEF_PARAM_CAPTION, CSFEF_PARAM_ORDER) VALUES (`+ids+`,-5,'TYPE_RECIPIENT','Тип получателя ',3) RETURNING CF_ID`
			////				_, err = Conn.Exec(sqlStr)
			////				if err != nil {
			////					Ошибка("  %+v", err)
			////				}
			//				//t, _=r.RowsAffected()
			//				//Инфо(" CF_ID %+v", t)
			//				//for r.Next(){
			//				//	var re interface{}
			//				//	err :=r.Scan(&re)
			//				//	if err != nil {
			//				//		Ошибка("  %+v", err)
			//				//	}
			//				//	Инфо("  %+v", re)
			//				//}
			//				//for РезультатЗапроса_1.Next() {
			//				//	var ID interface{}
			//				//	errScan := РезультатЗапроса.Scan(&ID)
			//				//	if errScan != nil {
			//				//		Ошибка("errScan	 %+v a %+v\n", errScan, ID)
			//				//
			//				//	}
			//				//	Инфо("CF_ID  %+v", ID)
			//				//}
			//				//Инфо(" РезультатЗапроса_1 %+v", РезультатЗапроса_1)
			//
			//			}
			//
			//		}
		}(ДанныеОСП, &запущено)
	}
	Инфо("запущено %+v", запущено)

	for {
		if запущено <= 0 {
			Инфо("запущено %+v", запущено)
			return
		}
	}

	//return nil
}

func ПроверитьИмяФильтра(ИмяФильтра string, Conn *sql.DB) interface{} {

	ПодготовленныйЗапросПроверкаИмениФильтра, err := Conn.Prepare(`select ID from CUSTOM_SQL_FILTER where SQL_FILTER_NAME= ?`)
	if err != nil {
		Ошибка("Ошибка подготовки запроса %+v %+v", err.Error())
		return nil
	}
	РезультатЗапросаSelect, err := ПодготовленныйЗапросПроверкаИмениФильтра.Query("Обращения, жалобы")
	var IDФильтра interface{}
	if err != nil {
		Ошибка("  %+v", err)
	}
	for РезультатЗапросаSelect.Next() {

		errScan := РезультатЗапросаSelect.Scan(&IDФильтра)
		if errScan != nil {
			Ошибка("errScan	 %+v a %+v\n", errScan, IDФильтра)
			continue
		}

	}
	Инфо("ИмяФильтра %+v IDзапроса %+v", ИмяФильтра, IDФильтра)
	return IDФильтра
}

func (client *Client) ДобавитьСпецФильтр(mes Сообщение) {
	Инфо("  %+v", "ДобавитьСпецФильтр")
	ДанныеПокдлюченийКОсп, err := ПолучитьДанныеПодключенийАИСОСП()
	if err != nil {
		Ошибка(" %+v ", err)
	}
	ДанныеИзАИС := make(chan map[string]interface{}, 40)
	количество := 0
	Инфо(" mes.Выполнить.Действие  %+v", mes.Выполнить.Действие)
	ВходящиеАргументы := mes.Выполнить.Действие["ДобавитьСпецФильтр"]
	ОтветКлиенту := map[string]interface{}{}

	if ВходящиеАргументы["sql_filter_name"] == nil {
		ОтветКлиенту["ОтветКлиенту"] = "Не указан sql_filter_name"
		Инфо(" %+v", ОтветКлиенту)
		return
	}
	if ВходящиеАргументы["object_name"] == nil {
		ОтветКлиенту["ОтветКлиенту"] = "Не указан object_name"
		Инфо(" %+v", ОтветКлиенту)
		return
	}
	if ВходящиеАргументы["sql_join_part"] == nil && ВходящиеАргументы["sql_where_part"] == nil {
		ОтветКлиенту["ОтветКлиенту"] = "Не указан sql_join_part и/или sql_where_part"
		Инфо(" %+v", ОтветКлиенту)
		return
	}

	sqlData := sqlStruct{
		ВходящиеАргументы: ВходящиеАргументы,
		SQLАргументы: []interface{}{
			".object_name",
			".sql_filter_name",
			".sql_join_part",
			".sql_where_part",
			".parent_id",
		},
		Name:  "Добавление CUSTOM_SQL_FILTER",
		Table: "CUSTOM_SQL_FILTER",
		Sql:   `INSERT INTO CUSTOM_SQL_FILTER (ID, OBJECT_NAME, SQL_FILTER_NAME,SQL_JOIN_PART, SQL_WHERE_PART, PARENT_ID) VALUES ((SELECT NEXT VALUE FOR SEQ_CUSTOM_SQL_FILTER FROM RDB$DATABASE),?,?,?,?,?) RETURNING ID`,
		//Sql:   `DELETE FROM CUSTOM_SQL_FILTER  WHERE SQL_FILTER_NAME= 'Сводное ИП с лицами, достигшими пенсионного возраста'`,
		//Sql:   `UPDATE CUSTOM_SQL_FILTER SQL_JOIN_PART = ?, SQL_WHERE_PART = ?  WHERE SQL_FILTER_NAME= 'Сводное ИП с лицами, достигшими пенсионного возраста'`,
	}

	if ВходящиеАргументы["аис"] == nil {
		ОтветКлиенту["ОтветКлиенту"] = "Необходимо указать в какой АИС добавлять фильтр"
		Инфо(" %+v", ОтветКлиенту)
		return
	}

	for _, аис := range ВходящиеАргументы["аис"].([]interface{}) {
		//if аис == "фссп" {
		//
		//	количество++
		//
		//	sqlData.Conn  = FSSPconnect()
		//	sqlData.ДанныеПодключения  = map[string]interface{}{
		//		"ОСП": "Управление ФССП",
		//		"КодОСП":26000,
		//	}
		//	go sqlData.ИзменитьДанныеВАИС(ДанныеИзАИС)
		//
		//}

		if аис == "осп" {
			for _, ДанныеОСП := range ДанныеПокдлюченийКОсп {
				количество++

				pwd := ДанныеОСП["pwd"].(string)
				ip := ДанныеОСП["ip_ais"].(string)

				sqlData.Conn = OSPconnect(pwd, ip)

				sqlData.ДанныеПодключения = map[string]interface{}{
					"ОСП":    ДанныеОСП["osp_name"].(string),
					"КодОСП": ДанныеОСП["osp_code"],
				}

				Инфо(" %+v", sqlData)

				sqlData.ИзменитьДанныеВАИС(ДанныеИзАИС)
			}
		}
	}

	for количество > 0 {
		ДанныеЗапросаИзОСП := <-ДанныеИзАИС
		if ДанныеЗапросаИзОСП != nil {
			Инфо("ДанныеЗапросаИзОСП  %+v", ДанныеЗапросаИзОСП)
			количество--
		}
		if количество == 0 {
			close(ДанныеИзАИС)
		}
	}
	Инфо(" %+v", ОтветКлиенту)
	return
}

/*
в зависимости от значения  выполняет подключение к соответсвующей базе данных и  вызывает функцию sqlData.ВыполнитьSQLАИС()
*/
func (sqlData sqlStruct) ВыполнитьЗапросВАИС() (map[string]map[string]РезультатSQLизАИС, error) {
	//Инфо("Выполняем:  %+v", "ВыполнитьЗапросВАИС")
	var ДанныеИзАис map[string]РезультатSQLизАИС
	var РезультатЗапроса chan map[string]РезультатSQLизАИС
	КоличествоГорутин := 0
	if sqlData.Клиент != nil {
		СообщениеКлиенту := Сообщение{
			Id:          0,
			От:          "io",
			Кому:        sqlData.Клиент.Login,
			Текст:       "Начинаю сбор данных в АИС",
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
		//Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)
		СообщениеКлиенту.СохранитьИОтправить(sqlData.Клиент)
	}

	//Инфо("sqlData 1014 %+v", sqlData)
	switch sqlData.БазаДанных {
	case "osp_code":
		Инфо("Выполняем скрипт на  %+v", "osp_code")
		КоличествоГорутин = 1
		РезультатЗапроса = make(chan map[string]РезультатSQLизАИС, 1)

		IPS := strings.Split(sqlData.Клиент.UserInfo.Info.Ip, ".")
		Подсеть := IPS[2]
		if IPS[2] != "6" {
			ПодсетьЧисло, _ := strconv.Atoi(IPS[2])
			ПодсетьЧисло = ПодсетьЧисло - 1
			Подсеть = strconv.Itoa(ПодсетьЧисло)
		}
		СообщениеКлиенту := Сообщение{
			Id:          0,
			От:          "io",
			Кому:        sqlData.Клиент.Login,
			Текст:       "Получение данных осп " + strconv.Itoa(sqlData.Клиент.UserInfo.Info.OspCode) + " подсеть " + Подсеть,
			MessageType: []string{"float_message", "error"},
			Content: struct {
				Target     string      `json:"target"`
				Data       interface{} `json:"data"`
				Html       string      `json:"html"`
				Обработчик string      `json:"обработчик"`
			}{
				//Target:"progress",
			},
		}
		//Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)
		СообщениеКлиенту.СохранитьИОтправить(sqlData.Клиент)

		ДанныеПокдлюченияОСП, ОшибкаЗапроса := sqlStruct{
			Name:   "Данные подключений к ОСП",
			Sql:    "SELECT * FROM fssp_configs.osp_address WHERE ip_ais = $1",
			Клиент: sqlData.Клиент,
			Values: [][]byte{
				[]byte("10.26." + Подсеть + ".34"),
			},
		}.Выполнить(nil)
		if ОшибкаЗапроса != nil {
			Ошибка(" %+v ", ОшибкаЗапроса)
		}
		Инфо("ДанныеПокдлюченияОСП %+v", ДанныеПокдлюченияОСП)
		sqlData.Conn = OSPconnect(ДанныеПокдлюченияОСП[0]["pwd"].(string), ДанныеПокдлюченияОСП[0]["ip_ais"].(string))
		sqlData.ДанныеПодключения = ДанныеПокдлюченияОСП[0]
		Инфо("sqlData %+v", sqlData)

		//sqlData.Sql =  `SELECT  * FROM V_SYS_CONNECTIONS WHERE LOGIN = UPPER(?)`
		//sqlData.АргументыДляАИС = []interface{}{sqlData.Клиент.UserInfo.Info.Login}

		go sqlData.ВыполнитьSQLАИС(РезультатЗапроса)

	case "fssp":
		Инфо("Выполняем скрипт на  %+v", "fssp")
		КоличествоГорутин = 1
		РезультатЗапроса = make(chan map[string]РезультатSQLизАИС, 1)
		sqlData.Conn = FSSPconnect()
		//go sqlData.ВыполнитьSQLвФССП(РезультатЗапроса)
		sqlData.ДанныеПодключения = map[string]interface{}{
			"osp_code": 26000,
			"osp_name": "Управление ФССП",
		}
		go sqlData.ВыполнитьSQLАИС(РезультатЗапроса)
	case "rdb":
		Инфо("Выполняем скрипт на  %+v", "rdb")
		КоличествоГорутин = 1
		РезультатЗапроса = make(chan map[string]РезультатSQLизАИС, 1)
		sqlData.Conn = RDBconnect()
		sqlData.ДанныеПодключения = map[string]interface{}{
			"osp_code": 260,
			"osp_name": "РДБ Управление ФССП",
		}
		go sqlData.ВыполнитьSQLАИС(РезультатЗапроса)
		//go sqlData.ВыполнитьSQLвРДБ(РезультатЗапроса)
	case "osp":
		Инфо("Выполняем скрипт на  %+v", "osp")
		ДанныеПодключенийАИСОСП, err := ПолучитьДанныеПодключенийАИСОСП()
		if err != nil {
			Ошибка("  %+v", err)
		}
		//КоличествоГорутин=2
		КоличествоГорутин = len(ДанныеПодключенийАИСОСП)
		РезультатЗапроса = make(chan map[string]РезультатSQLизАИС, КоличествоГорутин)
		for _, ДанныеПокдлюченияОСП := range ДанныеПодключенийАИСОСП {
			//sqlData.Conn =  OSPconnect(ДанныеПокдлюченияОСП["pwd"].(string) , ДанныеПокдлюченияОСП["ip_ais"].(string) )
			sqlData.ДанныеПодключения = ДанныеПокдлюченияОСП

			//if sqlData.Клиент != nil{
			//	СообщениеКлиенту := Сообщение{
			//		Id:      0,
			//		От:      "io",
			//		Кому:    sqlData.Клиент.Login,
			//		Текст:   "Запрашиваю данные в "+ДанныеПокдлюченияОСП["osp_name"].(string),
			//		MessageType: []string{"attention"},
			//		Content: struct {
			//			Target string `json:"target"`
			//			Data interface{} `json:"data"`
			//			Html string `json:"html"`
			//			Обработчик string `json:"обработчик"`
			//		}{
			//			Target:"progress",
			//		},
			//	}
			//	Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)
			//	СообщениеКлиенту.СохранитьИОтправить(sqlData.Клиент)
			//}
			//if sqlData.ПеребратьАргументы != nil {
			//
			//}

			go sqlData.ВыполнитьSQLАИС(РезультатЗапроса)
			//go sqlData.ВыполнитьSQLвОСП(ДанныеПокдлюченияОСП, РезультатЗапроса)
		}
	case "osp_rbd":

		ДанныеПодключенийАИСОСП, err := ПолучитьДанныеПодключенийАИСОСП()
		if err != nil {
			Ошибка("  %+v", err)
		}

		КоличествоГорутин = len(ДанныеПодключенийАИСОСП) + 1
		РезультатЗапроса = make(chan map[string]РезультатSQLизАИС, КоличествоГорутин+1)
		for _, ДанныеПокдлюченияОСП := range ДанныеПодключенийАИСОСП {

			sqlData.Conn = OSPconnect(ДанныеПокдлюченияОСП["pwd"].(string), ДанныеПокдлюченияОСП["ip_ais"].(string))
			sqlData.ДанныеПодключения = ДанныеПокдлюченияОСП

			if sqlData.Клиент != nil {
				СообщениеКлиенту := Сообщение{
					Id:          0,
					От:          "io",
					Кому:        sqlData.Клиент.Login,
					Текст:       "Запрашиваю данные в " + ДанныеПокдлюченияОСП["osp_name"].(string),
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
				//Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)
				СообщениеКлиенту.СохранитьИОтправить(sqlData.Клиент)
			}
			go sqlData.ВыполнитьSQLАИС(РезультатЗапроса)

		}

		sqlData.Conn = RDBconnect()
		sqlData.ДанныеПодключения = map[string]interface{}{
			"osp_code": 26000,
			"osp_name": "РДБ Управление ФССП",
		}
		go sqlData.ВыполнитьSQLАИС(РезультатЗапроса)

	default:
		if sqlData.Клиент != nil {
			СообщениеКлиенту := Сообщение{
				Id:          0,
				От:          "io",
				Кому:        sqlData.Клиент.Login,
				Текст:       "Не понятно куда отправлять запрос sqlData.БазаДанных" + sqlData.БазаДанных,
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
			Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)
			СообщениеКлиенту.СохранитьИОтправить(sqlData.Клиент)
		}
		Инфо(" Не понятно куда отправлять запрос sqlData.БазаДанных %+v", sqlData.БазаДанных)
		return nil, errors.New("Не понятно куда отправлять запрос sqlData.БазаДанных" + sqlData.БазаДанных)
	}

	ВсеДанные := map[string]map[string]РезультатSQLизАИС{
		sqlData.Name: {},
	} // [ИмяЗапроса][ОСП]РезультатSQLизАИС

	for КоличествоГорутин > 0 {
		ДанныеИзАис = <-РезультатЗапроса
		if ДанныеИзАис != nil {
			var ОСП string
			//Инфо(" ДанныеИзАис %+v", ДанныеИзАис)
			if ОСПкод := ДанныеИзАис[sqlData.Name].ОСП["osp_code"]; ОСПкод != nil {
				ОСП = strconv.Itoa(ОСПкод.(int))
			} else {
				ОСП = "ФССП"
			}

			ВсеДанные[sqlData.Name][ОСП] = ДанныеИзАис[sqlData.Name]

			//КэшироватьДанные(sqlData.Name, *sqlData.Клиент, ДанныеИзАис, ДанныеИзАис[sqlData.Name].ОСП["osp_name"].(string))

			КоличествоГорутин--

			if sqlData.Клиент != nil {
				СообщениеКлиенту := Сообщение{
					Id:          0,
					От:          "io",
					Кому:        sqlData.Клиент.Login,
					Текст:       "Данные получены из " + ДанныеИзАис[sqlData.Name].ОСП["osp_name"].(string) + ". Осталось получить ещё из " + strconv.Itoa(КоличествоГорутин) + " ОСП",
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
				//Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)
				go СообщениеКлиенту.СохранитьИОтправить(sqlData.Клиент)
			}
		}
		if КоличествоГорутин == 0 {
			break
		}
	}

	if sqlData.Клиент != nil {
		СообщениеКлиенту := Сообщение{
			Id:          0,
			От:          "io",
			Кому:        sqlData.Клиент.Login,
			Текст:       "Данные получены, обрабатываю результат",
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
		//Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)
		go СообщениеКлиенту.СохранитьИОтправить(sqlData.Клиент)
	}
	return ВсеДанные, nil
}

/*
Выполняет запрос в аис
*/

func (sqlData sqlStruct) ИзменитьДанныеВАИС(result chan<- map[string]interface{}) {

	//Данные := РезультатSQLизАИС{
	//	ВходящиеАргументы: sqlData.ВходящиеАргументы,
	//	Таблица: sqlData.Table,
	//	ОбвнолятьВРеальномВремени: sqlData.ОбновлениеВРеальномВремени,
	//}

	РезультатВыполнения := map[string]interface{}{
		"ДанныеПодключения": sqlData.ДанныеПодключения,
		"ошибка":            nil,
		"ЗатронутоСтрок":    0,
		"ИзменённыеДанные":  nil,
	}

	if sqlData.Conn == nil {
		Ошибка("Данные.Ошибка %+v sqlData.ВходящиеАргументы %+v", sqlData.ВходящиеАргументы)

		РезультатВыполнения["ошибка"] = "Нет связи с ФСС"
		result <- РезультатВыполнения
		return
	}

	var SQLАргументы []interface{}

	if sqlData.SQLАргументы != nil {
		SQLАргументы = ПолучитьSQLАргументыДляАис(sqlData.SQLАргументы, sqlData.ВходящиеАргументы)
	}

	ПодготовленныйЗапрос, err := sqlData.Conn.Prepare(sqlData.Sql)

	if err != nil {
		Ошибка("Ошибка подготовки запроса %+v %+v", err.Error(), sqlData.Sql)
		РезультатВыполнения["ошибка"] = "Ошибка подготовки запроса " + sqlData.Sql + err.Error()
		result <- РезультатВыполнения
		return
	}

	Инфо("ДанныеПодключения %+v \n SQLАргументы %+v", sqlData.ДанныеПодключения, SQLАргументы)

	РезультатЗапроса, err := ПодготовленныйЗапрос.Query(SQLАргументы...)

	if err != nil {

		Ошибка("Ошибка выполнения запроса %+v %+v sql.Sql %+v err\t %+v \n", sqlData.Name, sqlData.ColName, sqlData.Sql, err)

		РезультатВыполнения["ошибка"] = "Ошибка выполнения запроса " + sqlData.Sql + ";" + err.Error()
		result <- РезультатВыполнения
		return
	}

	if РезультатЗапроса == nil {
		Ошибка("Запрос '%+v' не вернул данных %+v\n", sqlData.Name, РезультатЗапроса)
		РезультатВыполнения["ЗатронутоСтрок"] = 0
		result <- РезультатВыполнения

		return
	}

	//ЗатронутоСтрок, err := РезультатЗапроса.RowsAffected()

	if err != nil {
		Ошибка(" не удалось получить количество затронутых строк	 %+v\n", err)
		РезультатВыполнения["ошибка"] = err.Error()

		result <- РезультатВыполнения
	}

	Столбцы, _ := РезультатЗапроса.Columns()

	ИзменённыеСтроки := []map[string]string{}
	for РезультатЗапроса.Next() {

		ЗначенияСтрокиNullString := make([]interface{}, len(Столбцы))

		for i, _ := range Столбцы {
			ЗначенияСтрокиNullString[i] = new(sql.NullString)
		}
		errScan := РезультатЗапроса.Scan(ЗначенияСтрокиNullString...)
		if errScan != nil {
			Ошибка("errScan	 %+v a %+v\n", errScan, ЗначенияСтрокиNullString)

		}
		ЗначенияСтроки := make([]interface{}, len(ЗначенияСтрокиNullString))
		КартаСтроки := map[string]string{}

		for i, Значение := range ЗначенияСтрокиNullString {
			ЗначенияСтроки[i] = Значение.(*sql.NullString).String
			КартаСтроки[Столбцы[i]] = Значение.(*sql.NullString).String
		}
		Инфо("КартаСтроки  %+v", КартаСтроки)
		ИзменённыеСтроки = append(ИзменённыеСтроки, КартаСтроки)

	}
	РезультатВыполнения["ЗатронутоСтрок"] = len(ИзменённыеСтроки)
	РезультатВыполнения["ИзменённыеДанные"] = ИзменённыеСтроки

	result <- РезультатВыполнения
}
