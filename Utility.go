package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	txt "text/template"
	"time"
)

func ОпределитьПользователяПоIp(IP string) map[string]interface{} {
	UID, err := sqlStruct{
		Name: "статичные_ip",
		Sql:  "SELECT uid FROM iobot.статичные_ip WHERE ip = $1",
		Values: [][]byte{
			[]byte(IP),
		},
	}.Выполнить(nil)
	if err != nil {
		Ошибка(">>>> Ошибка SQL запроса: %+v \n\n", err, IP)
		return nil
	}
	if len(UID) > 0 {
		FIO, errFio := sqlStruct{
			Name: "статичные_ip",
			Sql:  "SELECT second_name,givenname,initials, login FROM fssp_configs.users WHERE login = $1",
			Values: [][]byte{
				[]byte(UID[0]["uid"].(string)),
			},
		}.Выполнить(nil)

		if errFio != nil {
			Ошибка(">>>> ERROR \n %+v \n\n", err)
			return nil
		} else {
			if len(FIO) > 0 {
				return FIO[0]
			}
		}
		return nil
	} else {
		return nil
	}
}

func renderJS(tplName string, data interface{}) []byte {
	pattern := map[string]string{
		"name":    "tplFiles",
		"pattern": РабочаяПапка + "/html/*.*",
	}
	jsFiles, err := txt.New("").Funcs(tplFunc()).ParseGlob(pattern["pattern"])

	//for _, fn:=range jsFiles.Templates(){
	//	log.Printf("fn %+v\n", fn.Name())
	//}
	if err != nil {
		log.Print("render:", err)
		return nil
	}
	resJS := new(bytes.Buffer)

	err = jsFiles.ExecuteTemplate(resJS, tplName, data)

	if err != nil {
		log.Print("ExecuteTemplate:", err)
		return nil
	}

	return resJS.Bytes()
}

func render(tplName string, data interface{}) []byte {
	/*
		Данные клиента нужно добавлять в data вперед каждым вызовом render если в шаблоне предполагаеться использование данных пользователя

	*/
	Инфо(" render tplName %+v", tplName)
	pattern := map[string]string{
		"name":    "tplFiles",
		"pattern": РабочаяПапка + "/html/*.*",
	}

	tplFiles, err := template.New("").Funcs(tplFunc()).ParseGlob(pattern["pattern"])

	if err != nil {
		Ошибка(">>>> ERROR ОШИБКА ПАРСИНГА \n %+v \n\n", err)
		return nil
	}

	if data != nil {
		// если переданы данные для рендеринга content блока (форма авторизации, рабочий стол)
		//log.Printf("reflect.TypeOf(data) %+v\n", reflect.TypeOf(data).String() == "map[string]interface {}")
		//log.Printf("reflect.TypeOf(data).String() %+v\n", reflect.TypeOf(data).String())

		if reflect.TypeOf(data).String() == "map[string]interface {}" {
			if ContentData, ok := data.(map[string]interface{})["ContentData"]; ok {
				resHtml := new(bytes.Buffer)
				ContentTpls, err := tplFiles.Clone()
				if err != nil {
					Ошибка(">>>> ERROR tplFiles.Clone() \n %+v \n\n", err)
				}
				tplName := ContentData.(map[string]interface{})["tplName"].(string)
				data := ContentData.(map[string]interface{})["tplData"]
				err = ContentTpls.ExecuteTemplate(resHtml, tplName, data)
				if err != nil {
					Ошибка("\n\n>>>> ERROR ОШИБКА СОЗДАНИЯ HTML ШАБЛОНА <<<<: \n %+v", err)
					//Ошибка("\n\n>>>> ERROR ОШИБКА СОЗДАНИЯ HTML ШАБЛОНА <<<<: \n %+v \n\n >>>>> ПАРАМЕТРЫ ШАБЛОНА <<<<< \n %+v \n\n", err, data)
				}
				ContentHtml := `{{define "content"}}` + resHtml.String() + `{{end}}`
				tplFiles, err = tplFiles.Parse(ContentHtml)
				if err != nil {
					Ошибка(">>>> ERROR \n %+v \n\n", err)
				}
			}
		}
	}
	resHtml := new(bytes.Buffer)

	err = tplFiles.ExecuteTemplate(resHtml, tplName, data)

	if err != nil {
		Ошибка(">>>> ERROR ОШИБКА СОЗДАНИЯ HTML ШАБЛОНА <<<<: \n %+v", err)
		//Ошибка(">>>> ERROR ОШИБКА СОЗДАНИЯ HTML ШАБЛОНА <<<<: \n %+v \n\n >>>>> ПАРАМЕТРЫ ШАБЛОНА <<<<< \n %+v \n", err, data)
		return nil
	}
	//log.Printf("resHtml %+v\n", resHtml.String())
	return resHtml.Bytes()
}

///////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////tplFunc//////////////////////////////////////////////////////
//////////////////////////////tplFunc//////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

/*
идея НУЖНО ДОБАВИЬ В таблицу sql скрипты источник откуда брать данные, rdb fssp local
	нужно написать несколько функций для подключения к бд для посгреса, для firebird для mysql и конфиги для подключения положить в local базу
*/

type СтруктураШаблонов struct {
	Основной   string
	ВсеШаблоны map[string]string
}

func ПолучитьДанныеШаблона(Название string) СтруктураШаблонов {

	HTMLфайлы, err := template.New("").Funcs(tplFunc()).ParseGlob(РабочаяПапка + "/html/*.*")
	if err != nil {
		Ошибка("err %+v ", err)
	}
	Шаблоны := СтруктураШаблонов{}
	ВсеШаблоны := map[string]string{}
	Шаблоны.Основной = Название
	ПолучитьHTMLШаблонов([]byte(Название), HTMLфайлы, &ВсеШаблоны)
	Шаблоны.ВсеШаблоны = ВсеШаблоны
	return Шаблоны
}
func ПолучитьТекстШаблона(Название []byte, HTMLфайлы *template.Template) string {
	Lookup := HTMLфайлы.Lookup(string(Название))
	Инфо("Lookup.Tree. %+v", Lookup.Tree)
	ФайлШаблона, err := os.Open(РабочаяПапка + "/html/" + string(Название) + ".html")
	stat, err := ФайлШаблона.Stat()
	if err != nil {
		Ошибка(" %+v ", err)
	}

	// чтение файла
	БуферЧтенияФайла := make([]byte, stat.Size())
	_, err = ФайлШаблона.Read(БуферЧтенияФайла)
	if err != nil {
		Ошибка(" %+v ", err)
	}

	return string(БуферЧтенияФайла)
}

func ПолучитьHTMLШаблонов(Название []byte, HTMLфайлы *template.Template, ВсеШаблоны *map[string]string) {

	HTMLШаблона := ПолучитьТекстШаблона(Название, HTMLфайлы)
	//if _, ok := Шаблоны.Основной[string(Название)];ok {
	//	Шаблоны.Основной[string(Название)]=HTMLШаблона
	//}

	(*ВсеШаблоны)[string(Название)] = HTMLШаблона
	ИменаВложенныхШаблонов := ПолучитьИменаВложенныхШаблонов(HTMLШаблона)
	for _, ИмяВложенногоШаблона := range ИменаВложенныхШаблонов {
		ПолучитьHTMLШаблонов(ИмяВложенногоШаблона[1], HTMLфайлы, ВсеШаблоны)
	}
}
func ПолучитьИменаВложенныхШаблонов(HTMLШаблона string) [][][]byte {
	// плучаем все назваия вложеных шаблонов
	re := regexp.MustCompile(`(?m){{.?template\s"(?P<key>[a-zA-Z]*)".+?}}`)
	ИменаВложенныхШаблонов := re.FindAllSubmatch([]byte(HTMLШаблона), -1)
	return ИменаВложенныхШаблонов
}

func tplFunc() map[string]interface{} {
	var funcMap = template.FuncMap{

		"ИменаПолей": func(Данные []interface{}) string {
			Таблицы := map[string]string{}
			for _, Столбец := range Данные {
				ТаблицаСтолбец := strings.Split(Столбец.(string), ".")
				if Таблицы[ТаблицаСтолбец[0]] == "" {
					Таблицы[ТаблицаСтолбец[0]] = "'" + ТаблицаСтолбец[1] + "'"
				} else {
					Таблицы[ТаблицаСтолбец[0]] = "'" + ТаблицаСтолбец[1] + "'" + "," + Таблицы[ТаблицаСтолбец[0]]
				}
			}
			Where := ""
			for Таблица, Столбцы := range Таблицы {
				if Where == "" {
					Where = "(sys_objects.SOBJ_NAME = '" + Таблица + "' and SYS_FIELDS.SFLD_NAME IN (" + Столбцы + "))"
				} else {
					Where = Where + " OR (sys_objects.SOBJ_NAME = '" + Таблица + "' and SYS_FIELDS.SFLD_NAME IN (" + Столбцы + "))"
				}

			}
			return Where

		},

		"РазбитьСтроку": func(строка string, разделитель string) []string {

			return strings.Split(строка, разделитель)

		},

		"Месяцы": func(data []map[string]interface{}) map[string]string {
			Инфо("Месяцы %+v", data)
			Месяцы := map[string]string{"1": "Январь", "2": "Февраль", "3": "Март", "4": "Апрель", "5": "Май", "6": "Июнь", "7": "Июль", "8": "Август", "9": "Сентябрь", "10": "Октябрь", "11": "Ноябрь", "12": "Декабрь"}
			результат := map[string]string{}
			for _, номер := range data {

				результат[номер["мес"].(string)] = Месяцы[номер["мес"].(string)]
			}
			return результат
		},
		//"ДобавитьВМассив": func(Массив []interface{}, данные interface{}){
		//	return
		//},
		"Пусто": func(data interface{}) interface{} {

			if data == nil {
				return "Нет данных"
			} else {
				switch data.(type) {
				case int64:
					return strconv.Itoa(int(data.(float64)))
				case map[string]interface{}:
					return data
				}

			}

			return data
		},

		"Строка": func(Строки ...string) string {
			Результат := ""
			for _, Стр := range Строки {
				Результат = Результат + Стр
			}
			return Результат
		},
		"вJSON": func(data interface{}) interface{} {
			//Инфо("  reflect.TypeOf %+v",  reflect.TypeOf(data).Kind())
			//if reflect.TypeOf(data).Kind().String() == "map" {
			if data == nil {
				return nil
			}
			r, e := json.Marshal(data)
			Инфо("вJSON data %+v r %+v", data, r)
			if e != nil {
				Ошибка("  %+v", e)
			}
			return string(r)
			//}
		},
		"СоздатьКонекст": func(ВесьКонтекст map[string]interface{}, ДанныеДляШаблона interface{}) map[string]interface{} {
			//Инфо("ВесьКонтекст %+v", ВесьКонтекст["data"])
			//ВесьКонтекст["data"] = ДанныеДляШаблона
			//var НовыйКонтекст interface{}CopyMap

			НовыйКонтекст := map[string]interface{}{}
			for Имя, Значение := range ВесьКонтекст {
				if Имя == "data" {
					НовыйКонтекст["data"] = ДанныеДляШаблона
				} else {
					НовыйКонтекст[Имя] = Значение
				}
			}

			return НовыйКонтекст
		},
		"ПолучитьДанныеШаблона": ПолучитьДанныеШаблона,
		"СтрокуВHTML": func(Строка string, Data map[string]interface{}) template.HTML {
			resHtml := new(bytes.Buffer)
			tpl, err := template.New("HTMLTpl").Parse(Строка)
			if err != nil {
				Ошибка(" %+v ", err)
			}
			err = tpl.Execute(resHtml, Data)
			if err != nil {
				Ошибка(" %+v ", err)
			}
			return template.HTML(resHtml.String())
		},
		"render": func(Имя string, Данные interface{}) template.HTML {
			return template.HTML(string(render(Имя, Данные)))
		},

		"СтрокаСодержит": func(Строка interface{}, Подстрока string) bool {
			Инфо("  %+v", Подстрока)
			if reflect.TypeOf(Строка).Kind() == reflect.String {
				return strings.Contains(Строка.(string), Подстрока)
			} else if reflect.TypeOf(Строка).Kind() == reflect.Float64 {
				return strings.Contains(strconv.Itoa(int(Строка.(float64))), Подстрока)
			} else if reflect.TypeOf(Строка).Kind() == reflect.Int64 {
				return strings.Contains(strconv.Itoa(int(Строка.(int64))), Подстрока)
			} else if reflect.TypeOf(Строка).Kind() == reflect.Int {

				return strings.Contains(strconv.Itoa(Строка.(int)), Подстрока)
			}
			return false
		},
		"TypeOf": func(n interface{}) string {
			return reflect.TypeOf(n).String()
		},
		"Тип": func(n interface{}) string {
			return reflect.TypeOf(n).String()
		},
		/*НайтиВ
		Цель Масив Карт в котором ищем
		Поле  название искомого поля
		Значение  значение в Поле
		*/
		"НайтиВ": func(Цель []map[string]interface{}, Поле string, Значение interface{}) map[string]interface{} {
			Инфо("Цель %+v Поле %+v Значение %+v\n", Цель, Поле, Значение)
			for _, elem := range Цель {
				ЗначениеПоля, ЕстьПоле := elem[Поле]
				if ЕстьПоле {
					//log.Printf("ЗначениеПоля == Значение %+v\n", ЗначениеПоля == Значение)
					if ЗначениеПоля == Значение {
						Инфо("elem %+v\n", elem)
						return elem
					}
				}
			}
			//a = append(a[:i], a[i+1:]...)
			return nil
		},
		"Итого": func(слагаемые ...interface{}) float64 {
			var Итог float64
			for _, слагаемое := range слагаемые {
				var СлагаемоеЧисло float64
				switch reflect.TypeOf(слагаемое).Kind() {
				case reflect.String:
					СлагаемоеЧисло, _ = strconv.ParseFloat(слагаемое.(string), 64)
					break
				case reflect.Int:
					СлагаемоеЧисло = float64(слагаемое.(int))
				case reflect.Float64:
					СлагаемоеЧисло = слагаемое.(float64)
				}

				//Инфо("слагаемое %+v", reflect.TypeOf(слагаемое).Kind())
				Итог += СлагаемоеЧисло
			}
			return Итог
		},
		"Сумма": func(f interface{}, s interface{}) int {
			//Инфо("f %+v s %+v", f,s)
			if f == nil {
				f = 0
			}
			if s == nil {
				s = 0
			}
			var fv int
			var sv int
			if reflect.TypeOf(f).Kind() == reflect.Float64 {
				fv = int(f.(float64))
			} else if reflect.TypeOf(f).Kind() == reflect.Int {
				fv = f.(int)
			} else if reflect.TypeOf(f).Kind() == reflect.String {
				fv, _ = strconv.Atoi(f.(string))
			}
			if reflect.TypeOf(s).Kind() == reflect.Float64 {
				sv = int(s.(float64))
			} else if reflect.TypeOf(s).Kind() == reflect.Int {
				sv = s.(int)
			} else if reflect.TypeOf(s).Kind() == reflect.String {
				sv, _ = strconv.Atoi(s.(string))
			}
			//Инфо(" sv %+v,fv %+v",sv,fv, sv+fv)
			return sv + fv
		},
		"Разница": func(f interface{}, s interface{}) int {
			//Инфо("f %+v s %+v", f,s)
			if f == nil {
				f = 0
			}
			if s == nil {
				s = 0
			}
			var fv int
			var sv int
			if reflect.TypeOf(f).Kind() == reflect.Float64 {
				fv = int(f.(float64))
			} else if reflect.TypeOf(f).Kind() == reflect.Int {
				fv = f.(int)
			} else if reflect.TypeOf(f).Kind() == reflect.String {
				fv, _ = strconv.Atoi(f.(string))
			}
			if reflect.TypeOf(s).Kind() == reflect.Float64 {
				sv = int(s.(float64))
			} else if reflect.TypeOf(s).Kind() == reflect.Int {
				sv = s.(int)
			} else if reflect.TypeOf(s).Kind() == reflect.String {
				sv, _ = strconv.Atoi(s.(string))
			}
			//Инфо(" sv %+v,fv %+v",sv,fv, sv+fv)
			return sv - fv
		},

		"Последовательность": func(КоличествоЦиклов interface{}) []int {
			var МассивИтераций []int
			if reflect.TypeOf(КоличествоЦиклов).Kind() == reflect.String {
				var err error
				КоличествоЦиклов, err = strconv.Atoi(КоличествоЦиклов.(string))
				if err != nil {
					Ошибка("  %+v", err)
				}
				МассивИтераций = make([]int, КоличествоЦиклов.(int))
				for i, _ := range МассивИтераций {
					МассивИтераций[i] = i + 1
				}
			}

			return МассивИтераций
		},
		"Плюс": func(fistSlag int, secondSlag int) int {
			return fistSlag + secondSlag
		},
		"StringsJoin": func(StringsInter interface{}, sep string) string {
			//log.Printf("\n\n StringsInter %+v \n\n", StringsInter)
			if StringsInter == nil {
				return ""
			}
			Strings := make([]string, len(StringsInter.([]interface{})))
			for i, str := range StringsInter.([]interface{}) {
				Strings[i] = str.(string)
			}
			return strings.Join(Strings, sep)
		},
		"ПолучитьМенюБота": ПолучитьМенюБота,
		"CreatInitails":    CreatInitails,
		"Целое": func(число interface{}) int {
			if reflect.TypeOf(число).Kind() == reflect.Float64 {
				//log.Printf("int(число.(float64)) %+v\n", int(число.(float64)))
				//str  := strconv.FormatFloat(число.(float64),  'E', -1, 64)
				return int(число.(float64))
			}
			if reflect.TypeOf(число).Kind() == reflect.String {
				//log.Printf("int(число.(float64)) %+v\n", int(число.(float64)))
				//str  := strconv.FormatFloat(число.(float64),  'E', -1, 64)
				целое, err := strconv.Atoi(число.(string))
				if err != nil {
					Ошибка(" %+v ", err)
				}
				return целое
			}
			return 0
		},
		"ЧислоВСтроку": func(число interface{}) string {
			// Только с целыми числами / float64 так же не должно содержать дробной части
			//log.Printf("число %+v\n", число)
			//log.Printf("reflect.TypeOf(число).Kind() %+v\n", reflect.TypeOf(число).Kind())
			if число == nil {
				return "Нет данных"
			}

			if reflect.TypeOf(число).Kind() == reflect.Float64 {
				//log.Printf("int(число.(float64)) %+v\n", int(число.(float64)))
				//str  := strconv.FormatFloat(число.(float64),  'E', -1, 64)
				return strconv.Itoa(int(число.(float64)))
			} else {
				return strconv.Itoa(число.(int))
			}

		},
		"CurrentTime": func() string { return time.Now().Format("02.01.2006 15:04:05") },
		"Сегодня":     Сегодня,
		"СегодняДата": func() string { return time.Now().Format("02.01.2006") },
		"ParseTime": func(dateTime string) string {
			if dateTime == "" {
				return ""
			}
			//log.Printf("dateTime %+v\n", dateTime)

			t, err := time.Parse("2006-01-02T15:04:05.999999", dateTime)

			if err != nil {
				//log.Printf("\n !! ERR %+v\n", err)
				t, err = time.Parse("2006-01-02 15:04:05.999999", dateTime)
				if err != nil {
					t, err = time.Parse("2006-01-02T00:00:00+03:00", dateTime)
					if err != nil {
						t, err = time.Parse(time.RFC3339, dateTime)
						if err != nil {
							return err.Error()
						} else {
							return t.Format("02.01.2006 15:04:05")
						}
					} else {
						return t.Format("02.01.2006 15:04:05")
					}

				} else {
					return t.Format("02.01.2006 15:04:05")
				}

			} else {
				//log.Printf("t %+v\n", t)
				return t.Format("02.01.2006 15:04:05")
			}
		},

		"ParseDate": func(dateTime string) string {
			if dateTime == "" {
				return ""
			}
			//log.Printf("dateTime %+v\n", dateTime)

			t, err := time.Parse("2006-01-02T15:04:05.999999", dateTime)

			if err != nil {
				//log.Printf("\n !! ERR %+v\n", err)
				t, err = time.Parse("2006-01-02 15:04:05.999999", dateTime)
				if err != nil {
					t, err = time.Parse("2006-01-02T00:00:00+03:00", dateTime)
					if err != nil {
						return err.Error()
					} else {
						return t.Format("02.01.2006")
					}

				} else {
					return t.Format("02.01.2006")
				}

			} else {
				//log.Printf("t %+v\n", t)
				return t.Format("02.01.2006")
			}
		},
		"СоздатьДеревоКаталогов": func(Данные []map[string]interface{}, ИмяШаблона string) map[int]template.HTML {
			if Данные[0]["меню"] == nil {
				return nil
			}

			КартаМеню := map[int]template.HTML{}

			Категории := Данные[0]["меню"].(map[string]interface{})

			for ид_категории, ПунктКатегории := range Категории {
				if ПунктКатегории.(map[string]interface{})["размещение"].(float64) == 0 {
					ид, _ := strconv.Atoi(ид_категории)
					//Инфо("ид %+v", ид)
					HTMLКатегории := СоздатьПунктМеню(Категории, ид_категории, ИмяШаблона)
					КартаМеню[ид] = HTMLКатегории
				}
			}

			return КартаМеню

		},
		"УчетПК": УчетПК,
		"Заглавные": func(строка string) string {
			return strings.ToUpper(строка)
		},
		"ПолучитьДату": func(dateTime string) string {
			if dateTime == "" {
				return ""
			}
			//log.Printf("dateTime %+v\n", dateTime)

			t, err := time.Parse("2006-01-02T15:04:05.999999", dateTime)

			//parsing time "1987-02-20T00:00:00+03:00": extra text: "T00:00:00+03:00"    0001-01-01 00:00:00 +0000 UTC

			if err != nil {
				//log.Printf("\n !! ERR %+v\n", err)
				t, err := time.Parse("2006-01-02 15:04:05.999999", dateTime)
				if err != nil {
					t, err = time.Parse("2006-01-02", dateTime)
					if err != nil {
						Ошибка("err %+v t %+v", err, t)
					}
					return t.Format("02.01.2006")
				} else {
					return t.Format("02.01.2006")
				}

			} else {
				//log.Printf("t %+v\n", t)
				return t.Format("01.02.2006")
			}
		},

		"ДатаФормыАП": func(dateTime string) string {
			if dateTime == "" {
				return ""
			}
			//log.Printf("dateTime %+v\n", dateTime)

			t, err := time.Parse("2006-01-02T15:04:05.999999", dateTime)

			//parsing time "1987-02-20T00:00:00+03:00": extra text: "T00:00:00+03:00"    0001-01-01 00:00:00 +0000 UTC

			if err != nil {
				//log.Printf("\n !! ERR %+v\n", err)
				t, err := time.Parse("2006-01-02 15:04:05.999999", dateTime)
				if err != nil {
					t, err = time.Parse("2006-01-02", dateTime)
					if err != nil {
						Ошибка("err %+v t %+v", err, t)
					}
					return t.Format("02012006150405")
				} else {
					return t.Format("02012006150405")
				}

			} else {
				//log.Printf("t %+v\n", t)
				return t.Format("02012006150405")
			}
		},
	}
	return funcMap
}
func УчетПК(Данные map[string]interface{}) interface{} {
	//Инфо("Данные %+v", Данные["пк_из_тма"])
	//insert into iobot.запросы (имя)
	//values  ('пк_из_тма'),
	//	('пк_из_заббикс'),
	//	('пк_из_бух');

	ВсеПк := []map[string]interface{}{}

	заббикс := map[string]map[string]interface{}{}
	for _, пк_из_заббикс := range Данные["пк_из_заббикс"].([]map[string]interface{}) {
		заббикс[пк_из_заббикс["hw.addr"].(string)] = пк_из_заббикс
	}

	тма := map[string]map[string]interface{}{}
	for _, пк_из_тма := range Данные["пк_из_тма"].([]map[string]interface{}) {

		тма[пк_из_тма["inventory_num"].(string)] = пк_из_тма
		if заббикс[пк_из_тма["main_mac"].(string)] != nil {
			ВсеПк = append(ВсеПк, map[string]interface{}{
				"заббикс": заббикс[пк_из_тма["main_mac"].(string)],
				"тма":     тма[пк_из_тма["inventory_num"].(string)],
			})
			//заббикс[пк_из_тма["main_mac"].(string)]["тма"] = тма[пк_из_тма["inventory_num"].(string)]
		} else {
			ВсеПк = append(ВсеПк, map[string]interface{}{
				"заббикс": nil,
				"тма":     тма[пк_из_тма["inventory_num"].(string)],
			})
		}
	}

	Инфо("заббикс %+v", ВсеПк)
	//бухгалтерия := map[string][]map[string]interface{}{}
	//for _, пк_из_б := range Данные["пк_из_б"].([][]map[string]interface{}) {
	//	бухгалтерия["инвентарный"]=пк_из_б
	//}
	//
	//
	//
	//
	//for _, пк_из_тма := range Данные["пк_из_тма"].([]map[string]interface{}) {
	//	for  _, пк := range ВсеПк {
	//		if пк_из_тма["main_mac"] == пк["hw.addr"] {
	//			пк["main_mac"] = пк_из_тма["main_mac"]
	//			пк["inventory_num"] = пк_из_тма["inventory_num"]
	//			пк["motherboard_model"] = пк_из_тма["motherboard_model"]
	//			пк["ram_total"] = пк_из_тма["ram_total"]
	//			пк["cpu_name"] = пк_из_тма["cpu_name"]
	//			пк["осп"] = пк_из_тма["осп"]
	//		}
	//	}
	//}
	//
	//for _, пк_из_б := range Данные["пк_из_б"].([]map[string]interface{}) {
	//	for  _, пк := range ВсеПк {
	//		if пк_из_б["инвентарный"] == пк["inventory_num"] {
	//			пк["инвентарный"] = пк_из_б["инвентарный"]
	//			пк["motherboard"] = пк_из_б["motherboard"]
	//			пк["memory"] = пк_из_б["memory"]
	//			пк["cpu"] = пк_из_б["cpu"]
	//			пк["осп"] = пк_из_б["осп"]
	//		}
	//	}
	//}

	//Инфо("ВсеПк %+v", ВсеПк)

	return nil
}
func Сегодня() string {
	return time.Now().Format("02.01.2006 15:04:05")
}

//func СоздатьДерево (Данные []map[int]interface{})  {
//	    КартаМеню := map[int]template.HTML{}
//		Категории := Данные[0]
//
//		for ид, ПунктКатегории	:= range Категории {
//			if ПунктКатегории.(map[string]interface{})["размещение"] ==0 {
//				HTMLКатегории := СоздатьПунктМеню(Категории, ид)
//				КартаМеню[ид] = HTMLКатегории
//			}
//		}
//		Инфо("КартаМеню %+v", КартаМеню)
//}
//func renderMenuItem(menuElement ctx.MenuStruct) template.HTML {
//	tplFiles := ctx.MainCtx.Value("tplFiles").(*template.Template)
//	TplBlock, err := tplFiles.Clone()
//
//
//	var menuItem template.HTML
//	var htmlMenuItem = new(bytes.Buffer)
//	//log.Printf("menuElement 540 %+v\n", menuElement)
//	TplBlock.ExecuteTemplate(htmlMenuItem, "menu_item", menuElement)
//	menuItem = template.HTML(htmlMenuItem.String())
//	//log.Printf("%+v\n",menuItem)
//	return menuItem
//}
func СоздатьПунктМеню(Категории map[string]interface{}, ид string, ИмяШаблона string) template.HTML {

	Инфо("рендерим пункт меню с ид %+v  %+v", ид, Категории[ид])

	if ВложенныеКатегории := Категории[ид].(map[string]interface{})["вложенные_категории"].([]interface{}); len(ВложенныеКатегории) > 0 && ВложенныеКатегории[0] != nil {

		Категории[ид].(map[string]interface{})["Вложения"] = map[int]template.HTML{}
		ВложенныйHTML := make(map[int]template.HTML, len(ВложенныеКатегории))
		for Номер, идКатегории := range ВложенныеКатегории {

			ВложенныйHTML[Номер] = СоздатьПунктМеню(Категории, strconv.Itoa(int(идКатегории.(float64))), ИмяШаблона)

			Категории[ид].(map[string]interface{})["Вложения"].(map[int]template.HTML)[Номер] = ВложенныйHTML[Номер]
		}
	}
	// если вложенного пункта нету то рендерим HTML
	//Инфо("Категории %+v", Категории)
	HTMLПункта := template.HTML(string(render(ИмяШаблона, Категории[ид])))
	Инфо("HTMLПункта %+v", HTMLПункта)
	return HTMLПункта

}

func ПолучитьМенюБота() string {
	return "Добрый день! Если Вам нужна помощь куратора, выберите его справа и напишите сообщение."
}

func МассивСодержит(slice []interface{}, item interface{}) bool {
	set := make(map[interface{}]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

//http://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html
//var StdLog = log.New(os.Stderr, "", log.Lshortfile|log.Ltime)
var StdLog = log.New(os.Stderr, "", log.Llongfile)

//var StdinFile  = os.NewFile(1, РабочаяПапка+"/stdoutIo")
//
//var StdLogIO = log.New(StdinFile, "", log.Llongfile)

func Инфо(формат string, данные ...interface{}) {

	формат = strings.ReplaceAll(формат, "%+v", "\u001b[38;5;48m %+v  \u001b[0m\u001b[38;5;75m ")
	формат = strings.ReplaceAll(формат, "%#v", "\u001b[38;5;48m %#v  \u001b[0m\u001b[38;5;75m ")

	НоваяСтрокаФорматирования := " \u001b[0m\u001b[36m ИНФО: \u001b[38;5;75m " + формат + " \u001b[0m \n"

	str := fmt.Sprintf(НоваяСтрокаФорматирования, данные...)

	err := StdLog.Output(2, str)

	if err != nil {
		//log.Printf(" Ошибка вывода в консоль %+v", err)
		fmt.Fprintf(os.Stderr, " Ошибка вывода в консоль %+v данные %+v", err, данные)
	}
}

func Ошибка(формат string, данные ...interface{}) {

	формат = strings.ReplaceAll(формат, "%+v", "\u001b[38;5;204m %+v  \u001b[0m\u001b[38;5;1m")
	//fmt.Fprintf(os.Stderr," формат %+v",формат)
	формат = strings.ReplaceAll(формат, "%#v", "\u001b[38;5;204m %#v  \u001b[0m\u001b[38;5;1m")
	//fmt.Fprintf(os.Stderr," формат %+v",формат)
	str := fmt.Sprintf(" \u001b[48;5;124m ОШИБКА >> \u001b[0m  \u001b[38;5;1m  "+формат+" \u001b[0m \n", данные...)

	err := StdLog.Output(2, str)

	if err != nil {
		fmt.Fprintf(os.Stderr, " Ошибка вывода в консоль %+v данные %+v", err, данные)
	}

}

//func FrontReload() {
//	watcher, err := fsnotify.NewWatcher()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer watcher.Close()
//	err = watcher.Add(РабочаяПапка + "/static/css")
//	//err = watcher.Add(РабочаяПапка+"/static/js")
//
//	if err != nil {
//		log.Printf("РабочаяПапка %+v\n", РабочаяПапка)
//		log.Fatal(err)
//	}
//
//	//done := make(chan bool)
//
//	for {
//		select {
//		case event, ok := <-watcher.Events:
//
//			if !ok {
//				return
//			}
//
//			if event.Op&fsnotify.Write == fsnotify.Write {
//
//				fileName := event.Name[17:]
//				fileType := "wsstyle"
//				//if strings.Contains(fileName, "js"){
//				//	fileType = "wsjs"
//				//}
//
//				if server.Clients["maksimchuk@r26"] != nil {
//
//					responseMes := Сообщение{
//						Id:   0,
//						От:   "io",
//						Кому: "maksimchuk@r26",
//						Выполнить: struct {
//							Action   string                            `json:"action"`
//							Действие map[string]map[string]interface{} `json:"Действие"` // "НазваниеДействия" :{"имяАргумента_1":"значение аргумента_1"... }
//							Skill    int                               `json:"skill"`
//							Навык    json.Number                       `json:"Навык"`
//							Cmd      string                            `json:"cmd"`
//							Комманду string                            `json:"Комманду"`
//							Arg      struct {
//								Module string                 `json:"module"`
//								Tables []string               `json:"tables"`
//								Login  string                 `json:"login"`
//								Other  map[string]interface{} `json:"other"`
//							} `json:"Arg"`
//						}{
//							Action: "ReloadFiles",
//							Arg: struct {
//								Module string                 `json:"module"`
//								Tables []string               `json:"tables"`
//								Login  string                 `json:"login"`
//								Other  map[string]interface{} `json:"other"`
//							}{
//								Other: map[string]interface{}{
//									fileType: fileName,
//								},
//							},
//						},
//					}
//
//					server.Clients["maksimchuk@r26"].Message <- &responseMes
//				}
//			}
//		case err, ok := <-watcher.Errors:
//			log.Printf("watcher.Events: %+v\n", ok)
//			if !ok {
//				return
//			}
//			log.Println("error:", err)
//		}
//	}
//}

//
//func renderMenuFromTable (menuData ctx.Table) map[int]template.HTML {
//	start := time.Now()
//	menuHtml := map[int]template.HTML{}
//	log.Printf("menuData %+v\n", menuData)
//
//	menu := toMenuStruct(menuData.Rows)
//	//log.Printf("menuData %+v\n", menu)
//
//	//jsonMenu, _ := json.Marshal(menuData)
//	//json.Unmarshal(jsonMenu, &menu)
//
//	for idx := range menu {
//		//log.Printf("id %+v\n", idx,  menu[idx])
//		if  menu[idx].Parrent ==0{//
//
//			menuString := renderMenuString(menu, menuData.Index["id"], idx)
//
//
//			//log.Printf("menuString %+v\n",menuString)
//
//
//			menuHtml[idx] = menuString[idx]
//
//			//log.Printf("menuHtml[id] %+v\n",menuHtml[idx])
//		}
//	}
//	//log.Printf("menuHtml %+v\n", menuHtml)
//	fmt.Printf("время рендера меню %v\n",  time.Since(start))
//	return menuHtml
//}
//

//
//func renderMenuString (menuData []ctx.MenuStruct, menuIndex map[string]int , idx int) map[int]template.HTML {
//	itemHtml := map[int]template.HTML{}
//
//	if len(menuData[idx].Child) >0 {
//
//		//menuData[menuData[idx].Parrent].Children = make(map[int]template.HTML, len(menuData[idx].Child))
//		menuData[idx].Children = make(map[int]template.HTML, len(menuData[idx].Child))
//
//		for num, ChildId := range menuData[idx].Child{
//			//log.Printf("ChildId514 %+v idx %+v\n",ChildId , idx, menuData[idx])
//
//			childItemRowNum := menuIndex[strconv.Itoa(ChildId)]
//
//			itemHtmls :=renderMenuString(menuData, menuIndex, childItemRowNum)
//
//			itemHtml[idx] = itemHtmls[childItemRowNum]
//
//			menuData[idx].Children[num] = itemHtml[idx]
//			//menuData[menuData[idx].Parrent].Children[num] = itemHtml[idx]
//
//		}
//		//log.Printf("menuData Children%+v\n", menuData[idx])
//	}
//	//log.Printf("menuData Children%+v\n", menuData[idx])
//	itemHtml[idx] =renderMenuItem(menuData[idx])
//	//log.Printf("187 itemHtml%+v\n", itemHtml)
//	return itemHtml
//}
//
//func renderMenuItem(menuElement ctx.MenuStruct) template.HTML {
//	tplFiles := ctx.MainCtx.Value("tplFiles").(*template.Template)
//	TplBlock, err := tplFiles.Clone()
//
//
//	var menuItem template.HTML
//	var htmlMenuItem = new(bytes.Buffer)
//	//log.Printf("menuElement 540 %+v\n", menuElement)
//	TplBlock.ExecuteTemplate(htmlMenuItem, "menu_item", menuElement)
//	menuItem = template.HTML(htmlMenuItem.String())
//	//log.Printf("%+v\n",menuItem)
//	return menuItem
//}

//func toMenuStruct(menuData []map[string]interface{}) []ctx.MenuStruct {
//	menuStruct := make([]ctx.MenuStruct, len(menuData))
//
//	for i, menuItem := range menuData{
//		log.Printf("393 menuData %+v\n", menuItem)
//		log.Printf("menuItem[icon] %+v\n", menuItem["icon"].(interface{}))
//		if menuItem["parrent"] != nil {
//			menuStruct[i].Parrent = int(menuItem["parrent"].(float64))
//		}
//		if menuItem["description"] != nil {
//			menuStruct[i].Description = menuItem["description"].(string)
//		}
//		if menuItem["id"] != nil {
//			menuStruct[i].Id = int(menuItem["id"].(float64))
//		}
//		if menuItem["url"] != nil {
//			menuStruct[i].Url = menuItem["url"].(string)
//		}
//		if menuItem["title"] != nil {
//			menuStruct[i].Title = menuItem["title"].(string)
//		}
//		if menuItem["child"] != nil  {
//			menuStruct[i].Child = toIntSlice(menuItem["child"].([]interface{}))
//		}
//		if menuItem["access"] != nil  {
//			menuStruct[i].Access = toStringSlice(menuItem["access"].([]interface{}))
//		}
//
//		if menuItem["icon"].(map[string]interface{})["font"] != nil {
//
//			menuStruct[i].Icon.Font = menuItem["icon"].(map[string]interface{})["font"].(string)
//		}
//		//log.Printf("icon %+v\n",  menuItem["icon"].(map[string]interface{}))
//		if menuItem["icon"].(map[string]interface{})["image"] !=nil && menuItem["icon"].(map[string]interface{})["image"] !="" {
//			menuStruct[i].Icon.Image = menuItem["icon"].(map[string]interface{})["image"].(string)
//
//			//if menuItem["icon"].(map[string]interface{})["image"].(map[string]interface{})["name"] != nil {
//			//	menuStruct[i].Icon.Image.Name = menuItem["icon"].(map[string]interface{})["image"].(map[string]interface{})["name"].(string)
//			//}
//			//if menuItem["icon"].(map[string]interface{})["image"].(map[string]interface{})["pack_name"] != nil {
//			//	menuStruct[i].Icon.Image.PackName = menuItem["icon"].(map[string]interface{})["image"].(map[string]interface{})["pack_name"].(string)
//			//}
//		}
//	}
//	return menuStruct
//}
