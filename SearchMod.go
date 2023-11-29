package main

import (
	"github.com/opesun/goquery"
	"strings"
	"time"
)

var ЗапросыКлиентов = map[string]СтруктураЗапроса{}
type СтруктураЗапроса struct {
	Запрос string
	ТокеныЗапроса []string
	СловарьТекущегоВвода map[string]interface{}
}
/*
алгоритм отправлять запрос от клиента на сервер либо по нажатию на enter либо по пробелу или знаку припенания
		сохранять реузльтат предидущего слова с структуразапроса,
	-можно отправлять и каждое нажатие если количество символов после пробела больше 3 трёх,
	но предлагать варианты вводимого слова из словаря для избежания опечаток,
	по проеблу пытаться найти словосечтение и при этом продолжат обработку воода

алгоритм и неужно ещё на клиенте реализовать словарь, в который будет передаваться вводимое слово,
		тоесть человек начинает вводить слово,
		после нажатия третий буквы обращаемся в словарь и ищем совпадения по триграммам начала слов и передаём все слова клиенту показывая варианты в ниспадающем меню, и до пробела обрабатываем ввод из этого меню
		при этом нужно сделат что если нажат влево то вводим первое слово,
		если нажато вниз то переходим в ниспадающий список для выбора слова
		после ввода трёх букв нового слова показывать варианты слов и попытаться найти в бд сочетание первого слова и введёных букв второго

идея нужно почитать про пул соединений для postgres

*/
func (client *Client) НайтиИРАспарсить (вопрос Сообщение) {

	/*алгоритм обработаем вхдящие аргумннты ДЛЯ ДЕЙСТВИ НАЙТИ ИЛИ ПОСК ИЛИ НАЙДИ
	*/

	Инфо("Найти в %+v", вопрос.Выполнить.Действие["Найти"])
	АргументыПоиска := map[string]interface{}{}
	var ЧтоИскать interface{}
	var ГдеИскать interface{}
	if АргументыПоиска =вопрос.Выполнить.Действие["Найти"];АргументыПоиска != nil{
		ГдеИскать = АргументыПоиска["в"]
		ЧтоИскать = АргументыПоиска["что"]
	} else if АргументыПоиска =вопрос.Выполнить.Действие["Найди"];АргументыПоиска != nil{
		ГдеИскать = АргументыПоиска["в"]
		ЧтоИскать = АргументыПоиска["что"]
	} else if АргументыПоиска =вопрос.Выполнить.Действие["Поиск"];АргументыПоиска != nil{
		ГдеИскать = АргументыПоиска["в"]
		ЧтоИскать = АргументыПоиска["что"]
	}

	if ГдеИскать != nil {
		/*алгоритм вызываем функцию поиска в указаном месте*/
	}
	if ЧтоИскать != nil {
		if ЧтоИскать == "распарсить" {
//Инфо("498008 %+v", 498008)
			//for i := 498008; i > 497908; i-- {
			//	Инфо("i %+v", i)

					x, err := goquery.ParseUrl("https://habr.com/ru/all/page5/")

					if err != nil  {
						Ошибка("goquery err %+v ",err)
					}
				//<a href="https://habr.com/ru/post/497980/" class="post__title_link">Сколько воды утекло? Решаем задачу лунной походкой на Haskell</a>
					Node := x.Find(".post__title_link")
					Инфо("Node %+v", Node)
					for _ , нода := range Node {
						Инфо(" %+v", нода)
						Инфо("Node %+v", нода.Attr[0].Val)
						страница := ПолучитьСтраницуПоСсылке (нода.Attr[0].Val)
						сек:=0
						for {

							Инфо(" %+v сек", сек)
							time.Sleep(1 * time.Second)
							сек++
							if сек ==10 {
								break
							}
						}

						СохранитьСтраницу(страница)
					}


					//s := strings.TrimSpace()
					//ПолучитьСтраницуПоСсылке (Node.Attr("href"))

//					if  s != "" {
//
//					Результат ,ОшибкаЗапроса := sqlStruct{
//							Name:   `НазваниеЗапроса`,
//							Sql:    `INSERT INTO  knowledge_base.статья
//									(название, статья, автор, доступ,опубликован,осп,отдел,  текст, индекс) VALUES
//									($$`+s[0:200]+`$$ , $$`+s[0:5000]+`$$,  'maksimchuk@r26','["user"]',true,26911,
//26911,$$`+s[0:5000]+`$$,to_tsvector('russian', $$`+s[0:5000]+`$$)) RETURNING ид_статьи, название,автор,дата_создания, substring(текст for 200) текст`,
//									Values: [][]byte{
//								//[]byte(s[0:50]), // название первые 5 симовлов
//								//[]byte(s[0:5000]), // вся статья 5 0000 символов
//								//[]byte(client.Login), // вся статья 5 0000 символов
//							},
//							DBSchema:``,
//						}.Выполнить(nil)
//
//						if ОшибкаЗапроса != nil{
//						 Ошибка(">>>> Ошибка SQL запроса: %+v \n\n", ОшибкаЗапроса)
//						} else {
//							Инфо("Сохранили статью %+v", Результат)
//						}
//					}

				//time.Sleep(20 * time.Second)
			//}
		}
/*алгоритм парсим строку поиска и пытаемся найти то что передано,
					при этом создадим для текущего клиента объект в котором будем хранить в течении некоторого времени предидущий запрос поиска, да завершения ввода, или до нажатия на один из прелагаемых результатов поиска
*/

		//ТокеныЗапроса:= strings.Split(ЧтоИскать.(string), " ")
		// берём последнее слово для того чтобы найти в словаре подсказу и покать клиенту....
		//но пока что не буду это делать, сейчас сделаю только по пробелу или enter
		//ПоследнееСлово := ТокеныЗапроса[len(ТокеныЗапроса)-1]
		// Поищем в словаре в бд
		//РезультатЗапроса ,ОшибкаЗапроса := sqlStruct{
		//		Name:   `Поиск в словаре`,
		//		Sql:    `SELECT * FROM `,
		//		Values: [][]byte{
		//			[]byte(),
		//		},
		//		DBSchema:``,
		//	}.Выполнить(nil)
		//if ОшибкаЗапроса != nil{
		// Ошибка(">>>> Ошибка SQL запроса: %+v \n\n", ОшибкаЗапроса)
		//}
	}

}
func СохранитьСтраницу (s string){


						if  s != "" {

						Результат ,ОшибкаЗапроса := sqlStruct{
								Name:   `НазваниеЗапроса`,
								Sql:    `INSERT INTO  knowledge_base.статья
										(название, статья, автор, доступ,опубликован,осп,отдел,  текст, индекс) VALUES
										($$`+s[0:200]+`$$ , $$` +s+ `$$,  'maksimchuk@r26','["user"]',true,26911,
	26911,$$`+s+`$$,to_tsvector('russian', $$`+s+`$$)) RETURNING ид_статьи, название,автор,дата_создания, substring(текст for 200) текст`,
										Values: [][]byte{
									//[]byte(s[0:50]), // название первые 5 симовлов
									//[]byte(s[0:5000]), // вся статья 5 0000 символов
									//[]byte(client.Login), // вся статья 5 0000 символов
								},
								DBSchema:``,
							}.Выполнить(nil)

							if ОшибкаЗапроса != nil{
							 Ошибка(">>>> Ошибка SQL запроса: %+v \n\n", ОшибкаЗапроса)
							} else {
								Инфо("Сохранили статью %+v", Результат)
							}
						}



}
func ПолучитьСтраницуПоСсылке (адрес string ) string {

страница, err := goquery.ParseUrl(адрес)

if err != nil  {
Ошибка("goquery err %+v ",err)
}

	s := strings.Join(strings.Fields(страница.Text()), " ")
//Инфо("s %+v", s)
return s
}