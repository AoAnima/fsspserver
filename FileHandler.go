package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	xlsx "github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	txtTpl "text/template"

	"os"
	"path/filepath"
	"strings"
	"time"
)

func (client *Client) СохранитьФайл(Файл map[string]interface{}, ПутьДляЗагрузки string, НеМенятьИмя string) (map[string]string, error) {
	/*
		Файл = {
			"Папка":name,
			"Файл": ЧтениеФайла.result,
			"ИмяФайла": value.name
		}
	*/

	РасширениеФайлы := Файл["ИмяФайла"].(string)[bytes.LastIndex([]byte(Файл["ИмяФайла"].(string)), []byte(".")):]
	ПодПапка := ""
	if Файл["Папка"] != nil {
		ПодПапка = Файл["Папка"].(string)
	}
	ИмяФайла := Файл["ИмяФайла"].(string)

	if НеМенятьИмя == "false" {
		dt := time.Now()
		ТекущееВремя := dt.Format("01-02-2006_15h04m05s.000000000")
		if ПодПапка != "" {
			ИмяФайла = ПодПапка + "_" + ТекущееВремя + РасширениеФайлы
		} else {
			ИмяФайла = ТекущееВремя + РасширениеФайлы
		}
	}

	ПутьЗагрузки := filepath.Dir(os.Args[0]) + "/" + ПутьДляЗагрузки + "/" + ПодПапка + "/"

	Инфо("ПутьЗагрузки  %+v;  ИмяФайла %+v", ПутьЗагрузки, ИмяФайла)

	if _, err := os.Stat(ПутьЗагрузки); os.IsNotExist(err) {
		// path/to/whatever does not exist
		err = os.MkdirAll(ПутьЗагрузки, 0777)
		if err != nil {
			Ошибка(" Не удалось создать директорию  %+v", err)
			return nil, err
		}
	}

	file, err := os.Create(ПутьЗагрузки + ИмяФайла)
	if err != nil {
		Инфо(">>>> ERROR \n %+v \n\n", err)
		return nil, err
	}

	dec, err := base64.StdEncoding.DecodeString(Файл["Файл"].(string)[bytes.Index([]byte(Файл["Файл"].(string)), []byte("base64,"))+7:])
	if err != nil {
		Ошибка(" %+v", err)
		return nil, err
	}

	defer file.Close()

	if _, err := file.Write(dec); err != nil {
		Ошибка(" %+v", err)
		return nil, err
	}
	if err := file.Sync(); err != nil {
		Ошибка(" %+v", err)
		return nil, err
	}
	Инфо(" Сохранённый файл  %+v   %+v   %+v ", file, ПодПапка, ИмяФайла)
	return map[string]string{
		"ИмяФайла":        ИмяФайла,
		"Папка":           ПодПапка,
		"ПутьДляЗагрузки": ПутьДляЗагрузки,
	}, nil
}

/*
Нужно продумать структуру БД чтобы можно было выполнять функцию/плагин в которыйбудет передаваться набор sql запросов которые должны выполняться

например Запрос на обработку реестра.
Выбираем файл который должен содержать столбцы в определённом порядке
Пишем sql запросы которые нужно выполнить с данными из файла в аис или БД_ИО
Определяем какие столбцы в каком порядке подставляються в запрос
Если нужно Пишем sql запросы которые будут выполняться с полученными данными , например сохранить в бд статистики
Если нужно во
*/

func ОбработатьВУЦ(client *Client, сообщение Сообщение) map[string]interface{} {

	for НазваниеДействия, ДанныеДействия := range сообщение.Выполнить.Действие {
		Инфо(" НазваниеДействия %+v", НазваниеДействия)
		if Файл, ЕстьФайл := ДанныеДействия["Файл"]; ЕстьФайл {

			Результат := ОбработатьДанныеУЦ(client, Файл.(map[string]string))

			//Инфо("  Результат %+v",Результат )

			return Результат
		}
	}
	return nil
}

func СохранитьДанныеБухгалтерии(client *Client, сообщение Сообщение) map[string]interface{} {
	var ОткрытыйФайл *os.File
	var Ридер *csv.Reader
	ДанныеДействия := сообщение.Выполнить.Действие["загрузить данные бухгалтерии"]
	var Файл interface{}
	var ЕстьФайл bool
	Инфо("ДанныеДействия %+v %+v Файл %+v", сообщение.Выполнить.Действие, ДанныеДействия, ДанныеДействия["Файл"])

	if Файл, ЕстьФайл = ДанныеДействия["Файл"]; !ЕстьФайл {

		Инфо(" ЕстьФайл %+v", ЕстьФайл)
		//Результат := ОбработатьДанныеУЦ(client, Файл.(map[string]string))

		//Инфо("  Результат %+v",Результат )

		return nil
	}

	Путь := Файл.(map[string]string)["ПутьДляЗагрузки"]
	Папка := Файл.(map[string]string)["Папка"]
	ИмяФайла := Файл.(map[string]string)["ИмяФайла"]
	Инфо("ИмяФайла %+v", ИмяФайла)
	if strings.Contains(ИмяФайла, ".csv") {

		ОткрытыйФайл, Ридер = ОткрытьCSV(filepath.Dir(os.Args[0]) + "/" + Путь + "/" + Папка + "/" + ИмяФайла)
		Инфо(" %+v", Ридер)

		defer ОткрытыйФайл.Close()
		//for {
		//Строка, err := Ридер.Read()

		//Инфо("Строка  %+v", Строка)
		//if err != nil {
		//	Ошибка("  %+v", err)
		//	break
		//}
		//ФИО := Строка[0]
		//СерийныйНомер := Строка[1]
		//ОСПКод := Строка[2]
		//НачалоДействия := Строка[3]
		//ДатаОтзыва := Строка[4]
		//ИмяОСП := Строка[5]

		//}

	} else if strings.Contains(ИмяФайла, ".xlsx") {
		реестр, err := xlsx.OpenFile(filepath.Dir(os.Args[0]) + "/" + Путь + "/" + Папка + "/" + ИмяФайла)
		if err != nil {
			Ошибка(" %+v ", err)
		}
		//Листы := реестр.GetSheetList()
		ИмяРабочегоЛиста := реестр.GetSheetName(0)
		rows, err := реестр.Rows(ИмяРабочегоЛиста)
		Инфо(" %+v", rows)
		for rows.Next() {
			row, err := rows.Columns()
			if err != nil {
				Ошибка(" %+v ", err)
			}

			Материнка := row[3]
			Проц := row[4]
			Инвентарный := row[2]
			осп := ""
			if len(row) > 9 {
				осп = row[9]
			}

			//Материнка := row[1]
			//			Проц := row[2]
			//			memory := row[5]
			//			Инвентарный := row[3]
			//			осп := row[4]

			_, err = sqlStruct{
				Name: "Добавление новой записи рутокена",
				Sql:  `INSERT INTO oi.бухгалтерия (инвентарный,motherboard, cpu, осп, дата_загрузки) VALUES ($1,$2,$3,$4, current_date)`,
				Values: [][]byte{
					[]byte(Инвентарный),
					[]byte(Материнка),
					[]byte(Проц),
					[]byte(осп),
				},
			}.Выполнить(nil)

		}
		//Ошибка(" Необработчика для xlsx %+v", Файл["ИмяФайла"])
	} else {

		Ошибка(" Не потдерживаемый тип файла %+v", Файл.(map[string]string)["ИмяФайла"])
	}
	return map[string]interface{}{}
}

func ОбработатьДанныеУЦ(client *Client, Файл map[string]string) map[string]interface{} {
	var ОткрытыйФайл *os.File
	var Ридер *csv.Reader

	Путь := Файл["ПутьДляЗагрузки"]
	Папка := Файл["Папка"]
	ИмяФайла := Файл["ИмяФайла"]

	if strings.Contains(ИмяФайла, ".csv") {

		ОткрытыйФайл, Ридер = ОткрытьCSV(filepath.Dir(os.Args[0]) + "/" + Путь + "/" + Папка + "/" + ИмяФайла)

		defer ОткрытыйФайл.Close()

	} else if strings.Contains(ИмяФайла, ".xlsx") {
		Ошибка(" Необработчика для xlsx %+v", Файл["ИмяФайла"])
	} else {
		Ошибка(" Не потдерживаемый тип файла %+v", Файл["ИмяФайла"])
	}

	//_, err := sqlStruct{
	//	Name: "Очистим таблицу с данными УЦ",
	//	Sql:  `TRUNCATE skzi.new_ecp`,
	//	Values: [][]byte{},
	//}.Выполнить(nil)
	//if err != nil {
	//	Ошибка("  %+v", err)
	//}

	Результат := map[string]interface{}{}

	//ДанныеИзАИС := make(chan map[string]РезультатSQLизАИС)
	НомерСтроки := 0
	for {
		Строка, err := Ридер.Read()

		Инфо("%+v Строка  %+v", НомерСтроки, Строка)
		//Инфо("Строка  %+v", Строка)
		if err != nil {
			Ошибка("  %+v", err)
			break
		}

		ФИО := Строка[1]
		СерийныйНомер := Строка[2]
		ОСП := Строка[0]
		НачалоДействия := Строка[3]
		ДатаОтзыва := Строка[4]

		if НачалоДействия == "" || ДатаОтзыва == "" {
			continue
		}
		//ИмяОСП := Строка[4]

		//РезультатЗапросаИзАИС ,ОшибкаЗапроса := sqlStruct{
		//	Name:            `Получить ОСП пользовтяеля`,
		//	Sql:             `SELECT distinct  SPI_DEPARTMENT_CAPTION  FROM SPI
		//						join sys_users on SPI.SUSER_ID = sys_users.SUSER_ID
		//						where suser_fio = ?`,
		//	АргументыДляАИС: []interface{}{ФИО},
		//	БазаДанных:"fssp",
		//	Клиент: client,
		//}.ВыполнитьЗапросВАИС()
		//if ОшибкаЗапроса != nil {
		//	Ошибка(" %+v ", ОшибкаЗапроса)
		//}
		//
		//Инфо("РезультатЗапросаИзАИС %+v", РезультатЗапросаИзАИС)
		//var ОСП string
		//if len(РезультатЗапросаИзАИС["Получить ОСП пользовтяеля"]["26000"].Строки)>0{
		//	ОСП = РезультатЗапросаИзАИС["Получить ОСП пользовтяеля"]["26000"].Строки[0][0].(string)
		//}

		НомерСтроки++
		NewEcpsql := `INSERT INTO skzi.new_ecp (osp_name,  serial_number, date_start, date_end, common_name) VALUES ($1,$2, $3, $4, $5) ON CONFLICT (serial_number) DO UPDATE SET date_end = EXCLUDED.date_end  returning *;`
		//ФИО := Строка[0]
		//СерийныйНомер := Строка[1]
		//ОСПКод := Строка[2]
		//НачалоДействия := Строка[3]
		//ДатаОтзыва := Строка[4]
		Аргументы := [][]byte{
			[]byte(ОСП),
			[]byte(СерийныйНомер),
			[]byte(НачалоДействия),
			[]byte(ДатаОтзыва),
			[]byte(ФИО),
		}
		if ДатаОтзыва == "" {
			Аргументы[4] = nil
		}

		Данные, err := sqlStruct{
			Name:   "Добавление новой записи рутокена",
			Sql:    NewEcpsql,
			Values: Аргументы,
		}.Выполнить(nil)

		if err != nil {
			Ошибка("  %+v", err)
		}

		СообщениеКлиенту := &Сообщение{
			Текст: "",
			От:    "io",
			Кому:  client.Login,
			Контэнт: &ДанныеОтвета{
				HTML:      string(render("JournalSCZIRow", Данные[0])),
				Контейнер: "journal_form_1", СпособВставки: "добавить",
			},
		}
		client.Message <- СообщениеКлиенту

		//	if !strings.Contains(ФИО, "УФССП") {
		//
		//		sqlStr := `SELECT I_USER_CERT_REQUEST.DEPARTMENT_CODE_VKSP AS "osp_code",
		//		I_USER_CERT_REQUEST.ORGANIZATION_UNIT AS "osp_name",
		//		I_USER_CERT_REQUEST.CONTAINER_NAME AS "container_name",
		//		I_USER_CERT_REQUEST.TOKEN_NUMBER AS "token_number",
		//		CERTIFICATE.SERIAL_NUMBER AS "serial_number",
		//		CERTIFICATE.ISSUE_DATE AS "issue_date",
		//		I_USER_CERT_REQUEST.EXPIRATION_DATE AS "EXPIRATION_DATE",
		//		CERTIFICATE.EXPIRE_DATE AS "EXPIRE_DATE",
		//		CERTIFICATE.KEY_EXPIRE_DATE AS "KEY_EXPIRE_DATE",
		//		CERTIFICATE.CERT_STATUS AS "CERT_STATUS",
		//		I_USER_CERT_REQUEST.COMMON_NAME AS "common_name",
		//		document.DOC_NUMBER AS "doc_number",
		//		document.DOC_DATE AS "doc_date",
		//		CERTIFICATE.ISSUE_DATE AS "entry_date"
		//	FROM document
		//	JOIN document AS certdoc ON document.id= certdoc.PARENT_ID
		//	JOIN CERTIFICATE ON CERTDOC.id = CERTIFICATE.id
		//	JOIN I_USER_CERT_REQUEST ON document.id = I_USER_CERT_REQUEST.id
		//	WHERE CERTIFICATE.SERIAL_NUMBER = ?`
		//
		//		if strings.Contains(СерийныйНомер, " ") {
		//			continue
		//		}
		//
		//		ПараметрыЗапроса := sqlStruct{
		//			Name:            "рутокен",
		//			Sql:             sqlStr,
		//			АргументыДляАИС: []interface{}{СерийныйНомер},
		//		}
		//
		//		go ПараметрыЗапроса.ВыполнитьSQLвФССП(ДанныеИзАИС)
		//
		//
		//	ДанныеРутокена := map[string]РезультатSQLизАИС{}
		//
		//	for {
		//		ДанныеРутокена = <-ДанныеИзАИС
		//		Инфо(" ДанныеРутокена %+v", ДанныеРутокена)
		//		if ДанныеРутокена != nil {
		//			if len(ДанныеРутокена["рутокен"].КартаСтрок) > 0 {
		//				Рутокен := ДанныеРутокена["рутокен"].КартаСтрок[0]
		//
		//				NewEcpsql := `INSERT INTO skzi.new_ecp (osp_code,osp_name, container_name, token_number, serial_number, date_start, date_end, common_name, doc_number, doc_date) VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (serial_number) DO UPDATE SET date_end = EXCLUDED.date_end, token_number = EXCLUDED.token_number, container_name = EXCLUDED.container_name returning *;`
		//
		//				Аргументы := [][]byte{
		//					[]byte(Рутокен["osp_code"]),
		//					[]byte(Рутокен["osp_name"]),
		//					[]byte(Рутокен["container_name"]),
		//					[]byte(Рутокен["token_number"]),
		//					[]byte(Рутокен["serial_number"]),
		//					[]byte(НачалоДействия),
		//					[]byte(ДатаОтзыва),
		//					[]byte(Рутокен["common_name"]),
		//					[]byte(Рутокен["doc_number"]),
		//					[]byte(Рутокен["doc_date"]),
		//				}
		//				if ДатаОтзыва == "" {
		//					Аргументы[6] = nil
		//				}
		//
		//				Данные, err := sqlStruct{
		//					Name:   "Добавление новой записи рутокена",
		//					Sql:    NewEcpsql,
		//					Values: Аргументы,
		//				}.Выполнить(nil)
		//				if err != nil {
		//					Ошибка("  %+v", err)
		//				}
		//
		//				СообщениеКлиенту := &Сообщение{
		//					Текст: "",
		//					От:    "io",
		//					Кому:  client.Login,
		//					Контэнт: &ДанныеОтвета{
		//						HTML:      string(render("JournalSCZIRow", Данные[0])),
		//						Контейнер: "journal_form_1", СпособВставки: "добавить",
		//					},
		//				}
		//				client.Message <- СообщениеКлиенту
		//
		//				//Результат = append(Результат, Данные...)
		//				// прерываем цикл чтения из канала и продолжаем обрабатывать строки
		//			}
		//			break
		//		}
		//	}
		//} else {
		//		NewEcpsql := `INSERT INTO skzi.new_ecp (osp_code,osp_name,  serial_number, date_start, date_end, common_name) VALUES ($1,$2, $3, $4, $5, $6) ON CONFLICT (serial_number) DO UPDATE SET date_end = EXCLUDED.date_end  returning *;`
		//		//ФИО := Строка[0]
		//		//СерийныйНомер := Строка[1]
		//		//ОСПКод := Строка[2]
		//		//НачалоДействия := Строка[3]
		//		//ДатаОтзыва := Строка[4]
		//		Аргументы := [][]byte{
		//			[]byte(ОСПКод),
		//			[]byte(ИмяОСП),
		//			[]byte(СерийныйНомер),
		//			[]byte(НачалоДействия),
		//			[]byte(ДатаОтзыва),
		//			[]byte(ФИО),
		//		}
		//		if ДатаОтзыва == "" {
		//			Аргументы[4] = nil
		//		}
		//
		//		Данные, err := sqlStruct{
		//			Name:   "Добавление новой записи рутокена",
		//			Sql:    NewEcpsql,
		//			Values: Аргументы,
		//		}.Выполнить(nil)
		//		if err != nil {
		//			Ошибка("  %+v", err)
		//		}
		//
		//		СообщениеКлиенту := &Сообщение{
		//			Текст: "",
		//			От:    "io",
		//			Кому:  client.Login,
		//			Контэнт: &ДанныеОтвета{
		//				HTML:      string(render("JournalSCZIRow", Данные[0])),
		//				Контейнер: "journal_form_1", СпособВставки: "добавить",
		//			},
		//		}
		//		client.Message <- СообщениеКлиенту
		//	}
	}
	Инфо(" Результат %+v", Результат)
	return map[string]interface{}{
		"ОтветКлиенту": "Обработка данных УЦ завершена",
	}
}

func ОткрытьCSV(Путь string) (*os.File, *csv.Reader) {

	Инфо("Путь %+v", Путь)
	файл, ошибка := os.Open(Путь)
	if ошибка != nil {
		Ошибка("Не удаётся открыть файл %+v. Ошибка %+v", Путь, ошибка)
	}
	Инфо("файл  %+v", файл)
	ФайлРидер := csv.NewReader(bufio.NewReader(файл))
	ФайлРидер.Comma = ';'
	ФайлРидер.LazyQuotes = true

	return файл, ФайлРидер
}

/*
1 вариант
EXCEL шаблон - шаблон с готовой разметкой шапки
Нужно указать с какой строки начинаеться заполнение данными.
И В какой столбец какие данные вносить.

Если шаблона нет, то выгружать данные с именами из sql запроса, или если есть


2 вариант это нужно найти или придумать язык разметки заголовка, какие колонки объединять и т.д.. - НО ЛУЧШЕ НЕ ЗАМАРАЧИВАТЬСЯ И Либо заставить загружать xlsx шаблон, либо сделать визуальный редактор на странице, из которого генерироватьи  сохранять шаблон



*/

func СохранитьДанныеВEXCEL(client *Client, ДанныеОбработчка interface{}, ДанныеИзБд interface{}, ВходящиеДанные map[string]interface{}) map[string]interface{} {
	//(func(*Client, interface{}, interface{}) map[string]interface{})(client,ДанныеОбработчка, ДанныеИзБд)
	//var ИмяШаблона string
	//var КартаЗаполнения map[string]interface{}
	//var ПерваяСтрока int
	if ДанныеОбработчка == nil {
		//err := fmt.Errorf("Нет данных XLSX шаблона %+v", ДанныеОбработчка)

		return map[string]interface{}{
			"Ошибка": map[string]interface{}{
				"текст": fmt.Sprintf("Нет данных XLSX шаблона %+v", ДанныеОбработчка),
			},
		}
	}

	Инфо("ДанныеИзБд  %+v", ДанныеИзБд)
	//ДанныеОбработчка.(map[string]interface{})
	//ИмяШаблона = ДанныеОбработчка.(map[string]interface{})["шаблон"].(string)
	//КартаЗаполнения = ДанныеОбработчка.(map[string]interface{})["столбцы"].(map[string]interface{})
	//var err error
	//ПерваяСтрока, err = strconv.Atoi(ДанныеОбработчка.(map[string]interface{})["первая_строка"].(string))
	//if err != nil {
	//	Ошибка("  %+v", err)
	//}

	var Строки []map[string]interface{}
	Файлы := map[string]interface{}{} // map[string]map[string]interface{}{"файлы": {принадлежностьФайла: ссылка на файл}} принадлежностьФайла например осп,
	//reflect.TypeOf(ДанныеИзБд)
	switch ДанныеИзБд.(type) {
	case []map[string]interface{}:
		Строки = ДанныеИзБд.([]map[string]interface{})
		файл := ЗаписатьВEXCEL(Строки, ДанныеОбработчка.(map[string]interface{}), "")
		Файлы["файл"] = файл["файл"]
	case map[string]map[string]РезультатSQLизАИС: // [ИмяЗапроса][ОСП]РезультатSQLизАИС
		for ИмяДанных, ДанныеПоОСП := range ДанныеИзБд.(map[string]map[string]РезультатSQLизАИС) {
			for ОСП, Данные := range ДанныеПоОСП {
				файл := ЗаписатьВEXCEL(Данные.МассивСтрок, ДанныеОбработчка.(map[string]interface{}), ИмяДанных+"_"+ОСП+".xlsx")
				Файлы["файлы"] = map[string]interface{}{
					ОСП: файл["файл"],
				}
			}

		}

	}
	//Строки := ДанныеИзБд.([]map[string]interface{})
	//
	//
	//
	//
	//for N, строка := range Строки{
	//	НомерСтроки := ПерваяСтрока + N
	//	var Значение string
	//	for Столбец, ШаблонДанных := range КартаЗаполнения {
	//		if ШаблонДанных == "номер по порядку" {
	//			Значение =strconv.Itoa(N+1)
	//		} else {
	//			tpl, err := txtTpl.New("ДанныеЯчейки").Funcs(tplFunc()).Parse(ШаблонДанных.(string))
	//			if err != nil {
	//				Ошибка(" %+v ", err)
	//			}
	//			БайтБуферДанных := new(bytes.Buffer)
	//			err = tpl.Execute(БайтБуферДанных, map[string]map[string]interface{}{"строка": строка})
	//			if err != nil {
	//				Ошибка(" %+v ", err)
	//			}
	//			if БайтБуферДанных.String() == "<no value>" {
	//				Значение = ""
	//			} else {
	//				Значение = БайтБуферДанных.String()
	//			}
	//
	//		}
	//
	//		err = шаблон.SetCellValue("Лист1", Столбец+strconv.Itoa(НомерСтроки), Значение)
	//		if err != nil {
	//			Ошибка("  %+v", err)
	//		}
	//	}
	//}
	//НовыйФайл := "/uploads/tmp/"+ ИмяШаблона
	//err =	шаблон.SaveAs(filepath.Dir(os.Args[0]) + НовыйФайл)
	//if err != nil {
	//	//err := fmt.Errorf("Ну удалось сохранить временный файл XLSX %+v, ошибка %+v", НовыйФайл, err)
	//	return map[string]interface{}{
	//		"Ошибка": fmt.Sprintf("Ну удалось сохранить временный файл XLSX %+v, ошибка %+v", НовыйФайл, err),
	//	}
	//}
	Результат := map[string]interface{}{
		"файлы": Файлы,
	}
	Результат = Файлы

	СообщениеКлиенту := Сообщение{
		Id:          0,
		От:          "io",
		Кому:        client.Login,
		Текст:       "Данные сохраненные в EXCEL",
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
	СообщениеКлиенту.СохранитьИОтправить(client)

	return Результат
}

func ЗаписатьВEXCEL(Строки []map[string]interface{}, ДанныеОбработчка map[string]interface{}, ИмяФайла string) map[string]string {

	ИмяШаблона := ДанныеОбработчка["шаблон"].(string)
	КартаЗаполнения := ДанныеОбработчка["столбцы"].(map[string]interface{})
	var err error
	ПерваяСтрока, err := strconv.Atoi(ДанныеОбработчка["первая_строка"].(string))
	if err != nil {
		Ошибка("  %+v", err)
	}

	шаблон, err := xlsx.OpenFile(filepath.Dir(os.Args[0]) + "/uploads/document_tpl/" + ИмяШаблона)
	if err != nil {
		Ошибка("  %+v", err)
	}

	for N, строка := range Строки {
		НомерСтроки := ПерваяСтрока + N
		var Значение string
		for Столбец, ШаблонДанных := range КартаЗаполнения {
			if ШаблонДанных == "номер по порядку" {
				Значение = strconv.Itoa(N + 1)
			} else {
				tpl, err := txtTpl.New("ДанныеЯчейки").Funcs(tplFunc()).Parse(ШаблонДанных.(string))
				if err != nil {
					Ошибка(" %+v ", err)
				}
				БайтБуферДанных := new(bytes.Buffer)
				err = tpl.Execute(БайтБуферДанных, map[string]map[string]interface{}{"строка": строка})
				if err != nil {
					Ошибка(" %+v ", err)
				}
				if БайтБуферДанных.String() == "<no value>" {
					Значение = ""
				} else {
					Значение = БайтБуферДанных.String()
				}

			}
			//Инфо("Значение %+v", Значение)
			err = шаблон.SetCellValue("Лист1", Столбец+strconv.Itoa(НомерСтроки), Значение)
			if err != nil {
				Ошибка("  %+v", err)
			}
		}
	}
	if ИмяФайла == "" {
		ИмяФайла = "(" + Сегодня() + ")" + ИмяШаблона
	}
	НовыйФайл := "/uploads/tmp/" + ИмяФайла
	err = шаблон.SaveAs(filepath.Dir(os.Args[0]) + НовыйФайл)
	if err != nil {
		//err := fmt.Errorf("Ну удалось сохранить временный файл XLSX %+v, ошибка %+v", НовыйФайл, err)
		return map[string]string{
			"Ошибка": fmt.Sprintf("Ну удалось сохранить временный файл XLSX %+v, ошибка %+v", НовыйФайл, err),
		}
	}
	return map[string]string{"файл": НовыйФайл}
}
