package main

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func (client *Client)ПоказатьСписокЗаявок (mes Сообщение) map[string]interface{} {

	ЗаявкиВОчереди ,err:= sqlStruct{
		Name:   "tickets",
		Sql:    "SELECT * FROM iobot.tickets WHERE status = 'в очереди'",
		Values: [][]byte{},
	}.Выполнить(nil)
	if err != nil{
		Инфо(">>>> Ошибка SQL запроса: %+v \n\n",err)
	}

	СписокЗаявок := map[string]interface{}{
		"Контейнер":"main_content.tickets_list",
		"HTML": render("TicketsList", ЗаявкиВОчереди),
	}

	return СписокЗаявок
}

func (client *Client)СоздатьНавык (вопрос Сообщение) {
	Инфо(">>>>> СоздатьНавык mes %+v\n", вопрос)
	Аргументы := вопрос.Выполнить.Действие["СоздатьНавык"]

	Значения := [][]byte{}
	МеткиПодстановки:=""
	Столбцы := ""

	if Название, ЕстьНазвание :=  Аргументы["skill_name"];ЕстьНазвание {
		Столбцы = Столбцы+ "name"
		Значения=append(Значения, []byte(Название.(string)))
		МеткиПодстановки = МеткиПодстановки + "$"+strconv.Itoa(len(Значения))
	} else {
		// нет имени не возможно создать навык без имени
		СообщениеКлиенту:= &Сообщение{
			Текст:   "Название навыка обязательно!",
			От: "io",
			Кому:client.Login,
			MessageType:[]string{"irritation","io_action"},
		}
		//СообщениеКлиенту.СохранитьЛогСообщения()
		//client.Message<-СообщениеКлиенту
		СообщениеКлиенту.СохранитьИОтправить(client)
		return
	}
	if Комманды, ЕстьКомманды :=  Аргументы["skill_action"];ЕстьКомманды {
		Столбцы = Столбцы+ ", actions"
		Значения=append(Значения, []byte(`{"cmd":"`+Комманды.(string)+`"}`))
		//Инфо("Комманды %+v\n", `{"cmd":"`+Комманды.(string)+`"}`)
		МеткиПодстановки = МеткиПодстановки + ",$"+strconv.Itoa(len(Значения))
	}else {
		// нет имени не возможно создать навык без имени
		СообщениеКлиенту:= &Сообщение{
			Текст:   "Комманды навыка обязательны!",
			От: "io",
			Кому:client.Login,
			MessageType:[]string{"irritation","io_action"},
		}
		//СообщениеКлиенту.СохранитьЛогСообщения()
		//client.Message<-СообщениеКлиенту
		СообщениеКлиенту.СохранитьИОтправить(client)
		return
	}
	if Разрешение, ЕстьРазрешение :=  Аргументы["skill_self_use"];ЕстьРазрешение {
		Столбцы = Столбцы + ",self_use"
		Значения=append(Значения, []byte(Разрешение.(string)))
		МеткиПодстановки = МеткиПодстановки + ",$"+strconv.Itoa(len(Значения))

		Столбцы = Столбцы + ",access"
		Значения=append(Значения, []byte(`["user"]`))
		МеткиПодстановки = МеткиПодстановки + ",$"+strconv.Itoa(len(Значения))
	}
	if Описание, ЕстьОписание :=  Аргументы["skill_description"];ЕстьОписание {
		if len(Описание.(string))>0 {
			Столбцы = Столбцы + ",description"
			Значения=append(Значения, []byte(Описание.(string)))
			МеткиПодстановки = МеткиПодстановки + ",$"+strconv.Itoa(len(Значения))
		}
	}

	if Проблемы, ЕстьПроблемы :=  Аргументы["skill_problems"];ЕстьПроблемы {
		if len(Проблемы.(string)) > 0 {
			МассивПроблем := strings.Split(Проблемы.(string),",")
			Столбцы = Столбцы + ",problem"
			СтрокаПроблем := ""
			//Инфо("МассивПроблем %+v\n", МассивПроблем)
			for _, проблема := range МассивПроблем{
				//МассивПроблем[н]=strings.TrimSpace(проблема)
				if СтрокаПроблем == ""{
					СтрокаПроблем = `"`+strings.TrimSpace(проблема)+`"`
				} else {
					СтрокаПроблем = СтрокаПроблем +`,"`+strings.TrimSpace(проблема)+`"`
				}

			}

			СтрокаПроблем="["+СтрокаПроблем+"]"
			Значения=append(Значения, []byte(СтрокаПроблем))
			//Инфо("СтрокаПроблем %+v\n", СтрокаПроблем)
			МеткиПодстановки = МеткиПодстановки + ",$"+strconv.Itoa(len(Значения))
		}
	}
	if Метки, ЕстьМетки :=  Аргументы["skill_keywords"];ЕстьМетки {
		if len(Метки.(string)) > 0 {
			МассивМеток := strings.Split(Метки.(string),",")
			Столбцы = Столбцы + ",keywords"
			//Инфо("МассивМеток %+v\n", МассивМеток)
			СтрокаМеток := ""
			for _, метка := range МассивМеток {
				//МассивПроблем[н]=strings.TrimSpace(проблема)
				if СтрокаМеток == "" {
					СтрокаМеток = `"` + strings.TrimSpace(метка) + `"`
				} else {
					СтрокаМеток = СтрокаМеток + `,"` + strings.TrimSpace(метка) + `"`
				}

			}
			СтрокаМеток = "[" + СтрокаМеток + "]"
			//Инфо("СтрокаПроблем %+v\n", СтрокаМеток)
			Значения = append(Значения, []byte(СтрокаМеток))
			МеткиПодстановки = МеткиПодстановки + ",$" + strconv.Itoa(len(Значения))
		}
	}

	Столбцы = Столбцы + ",creator"
	Значения=append(Значения, []byte(client.Login))
	МеткиПодстановки = МеткиПодстановки + ",$"+strconv.Itoa(len(Значения))

	Результат ,err:= sqlStruct{
		Name:   "io_actions",
		Sql:    "INSERT INTO iobot.io_actions ("+Столбцы+",date_create) VALUES ("+МеткиПодстановки+",NOW()) ON CONFLICT (name) DO NOTHING RETURNING *",
		Values: Значения,
	}.Выполнить(nil)
	if err != nil{
		Инфо(">>>> Ошибка SQL запроса: %+v \n sql %+v \n",err , "INSERT INTO iobot.io_actions ("+Столбцы+",date_create) VALUES ("+МеткиПодстановки+",NOW()) ON CONFLICT (name) DO NOTHING RETURNING *")

		for _, значение := range Значения{
			Инфо("Значение %+v\n", string(значение))
		}


		СообщениеКлиенту:= &Сообщение{
			Текст:   "Ошибка сохранения навыка! " + err.Error(),
			От: "io",
			Кому:client.Login,
			MessageType:[]string{"irritation"},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
		return
		//client.СохранитьИОтправитьСообщениеКлиенту(СообщениеКлиенту)
	} else {
		if len(Результат) > 0{

			СообщениеКлиенту:= &Сообщение{
				Текст:   "Навык создан. " + string(render("SkillTpl", Результат[0])),
				От: "io",
				Кому:client.Login,
				//Контэнт: &ДанныеОтвета{
				//	Контейнер:"сообщение",
				//	HTML:string(render("SkillTpl", Результат[0])),
				//},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
			return
		} else {
			СообщениеКлиенту:= &Сообщение{
				Текст:   "Ошибка сохранения навыка! " + err.Error(),
				От: "io",
				Кому:client.Login,
				MessageType:[]string{"irritation"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
		}
	}
	//Инфо("Результат СоздатьНавык %+v\n", Результат)

}


func (client *Client)РедакторНавыка (mes Сообщение) map[string]interface{} {
	Инфо(">>>>> РедакторНавыка mes %+v\n", mes)
	return nil

}

func (client *Client)СоздатьЗаявку(вопрос Сообщение) map[string]interface{}{
	СоздатьНовыйТикет ,err:= sqlStruct{
		Name:   "tickets",
		Sql:    "INSERT INTO iobot.tickets (login, ip, problem) VALUES ($1,$2,$3) RETURNING *",
		Values: [][]byte{
			[]byte(client.Login),
			[]byte(client.Ip),
			[]byte(вопрос.Текст),
		},
	}.Выполнить(nil)
	if err != nil{
		Инфо(">>>> Ошибка SQL запроса: %+v \n\n",err)
	}
	if len(СоздатьНовыйТикет)>0 {
		НовыйТикет := СоздатьНовыйТикет[0]
		ЗапросОчереди ,err:= sqlStruct{
			Name:   "tickets",
			Sql:    "SELECT count(*) FROM iobot.tickets WHERE (status = 'в очереди' OR status = 'в работе') AND ticket_id < $1",
			Values: [][]byte{
				[]byte(НовыйТикет["ticket_id"].(string)),
			},
		}.Выполнить(nil)
		if err != nil{
			Инфо(">>>> Ошибка SQL запроса: %+v \n\n",err)
		}
		if len(ЗапросОчереди)>0{
			Очередь := ЗапросОчереди[0]
			НовыйТикет["очередь"] = Очередь["count"]
		}
		return map[string]interface{}{
			"ОтветКлиенту":render("NewTicket", НовыйТикет),
		}

	}

	return nil
}



func (client *Client)НовыйНавык (вопрос Сообщение) map[string]interface{} {
	Инфо(">>>>> НовыйНавык mes %+v\n", вопрос)

	return map[string]interface{}{"Контейнер":"сообщение","HTML": render("NewSkillForm", nil)}
}



func (client *Client) УдалитьНавык (вопрос Сообщение) map[string]interface{} {
	//Аргументы:= вопрос.Выполнить.Действие["УдалитьНавык"]
	Аргументы:= вопрос.ВходящиеАргументы

	if ИД, ЕстьИД :=  Аргументы["action_id"];ЕстьИД {

		Результат ,err:= sqlStruct{
			Name:   "io_actions",
			Sql:    "UPDATE iobot.io_actions SET delete = true WHERE action_id = $1 RETURNING name",
			Values: [][]byte{
				[]byte(ИД.(string)),
			},
		}.Выполнить(nil)
		if err != nil{
			Инфо(">>>> Ошибка SQL запроса: %+v \n\n",err)
			СообщениеКлиенту:= &Сообщение{
				Текст:   "Ошибка удалёния навыка "+err.Error(),
				От: "io",
				Кому:client.Login,
				MessageType:[]string{"irritation","io_action"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
		}
		ТекстСообщения :=""
		if len(Результат)>0{
			Инфо("Результат УдалитьНавык %+v\n", Результат)
			ТекстСообщения = "Навык "+Результат[0]["name"].(string)+", помечен как удалённый, и не будет отображаться в списке активных навыков."
		} else if len(Результат)==0{
			ТекстСообщения = "Ни один навык не был удалён"
		}
		СообщениеКлиенту:= &Сообщение{
			Текст:  ТекстСообщения,
			От: "io",
			Кому:client.Login,
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
	} else {
		СообщениеКлиенту:= &Сообщение{
			Текст:   "ИД Навыка не определён",
			От: "io",
			Кому:client.Login,

		}
		СообщениеКлиенту.СохранитьИОтправить(client)
		//СообщениеКлиенту.СохранитьЛогСообщения()
		//client.Message<-СообщениеКлиенту
	}
	return nil
}
func (client *Client) СохранитьИзмененияНавыка (вопрос Сообщение) {
	Инфо(">>>>> СохранитьИзмененияНавыка mes %+v\n", вопрос)

	Аргументы:= вопрос.Выполнить.Действие["СохранитьИзмененияНавыка"]

	Значения := [][]byte{}
	МеткиПодстановки:=""
	Столбцы := ""
	if ИД, ЕстьИД :=  Аргументы["skill_id"];ЕстьИД {
		//Столбцы = Столбцы+ "action_id"
		Значения=append(Значения, []byte(ИД.(string)))
		//МеткиПодстановки = МеткиПодстановки + "$"+strconv.Itoa(len(Значения))
	}
	if Название, ЕстьНазвание :=  Аргументы["skill_name"];ЕстьНазвание {
		Столбцы = Столбцы+ "name"
		Значения=append(Значения, []byte(Название.(string)))
		МеткиПодстановки = МеткиПодстановки + "$"+strconv.Itoa(len(Значения))
	} else {
		// нет имени не возможно создать навык без имени
		СообщениеКлиенту:= &Сообщение{
			Текст:   "Название навыка обязательно!",
			От: "io",
			Кому:client.Login,
			MessageType:[]string{"irritation","io_action"},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
		//СообщениеКлиенту.СохранитьЛогСообщения()
		//client.Message<-СообщениеКлиенту
		return
	}
	if Комманды, ЕстьКомманды :=  Аргументы["skill_action"];ЕстьКомманды {
		Столбцы = Столбцы+ ", actions"
		Значения=append(Значения, []byte(`{"cmd":"`+Комманды.(string)+`"}`))
		//Инфо("Комманды %+v\n", `{"cmd":"`+Комманды.(string)+`"}`)
		МеткиПодстановки = МеткиПодстановки + ",$"+strconv.Itoa(len(Значения))
	}else {
		// нет имени не возможно создать навык без имени
		СообщениеКлиенту:= &Сообщение{
			Текст:   "Комманды навыка обязательны!",
			От: "io",
			Кому:client.Login,
			MessageType:[]string{"irritation","io_action"},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
		//СообщениеКлиенту.СохранитьЛогСообщения()
		//client.Message<-СообщениеКлиенту
		return
	}
	if Разрешение, ЕстьРазрешение :=  Аргументы["skill_self_use"];ЕстьРазрешение {
		Столбцы = Столбцы + ",self_use"
		Значения=append(Значения, []byte(Разрешение.(string)))
		МеткиПодстановки = МеткиПодстановки + ",$"+strconv.Itoa(len(Значения))

		Столбцы = Столбцы + ",access"
		Значения=append(Значения, []byte(`["user"]`))
		МеткиПодстановки = МеткиПодстановки + ",$"+strconv.Itoa(len(Значения))
	} else {
		Столбцы = Столбцы + ",access"
		Значения=append(Значения, []byte(`["admin","curators","architect"]`))
		МеткиПодстановки = МеткиПодстановки + ",$"+strconv.Itoa(len(Значения))
	}

	if Описание, ЕстьОписание :=  Аргументы["skill_description"];ЕстьОписание {
		if len(Описание.(string))>0 {
			Столбцы = Столбцы + ",description"
			Значения=append(Значения, []byte(Описание.(string)))
			МеткиПодстановки = МеткиПодстановки + ",$"+strconv.Itoa(len(Значения))
		}
	}

	if Проблемы, ЕстьПроблемы :=  Аргументы["skill_problems"];ЕстьПроблемы {
		if len(Проблемы.(string)) > 0 {
			МассивПроблем := strings.Split(Проблемы.(string),",")
			Столбцы = Столбцы + ",problem"
			СтрокаПроблем := ""
			Инфо("Проблемы.(string) %+v\n", Проблемы.(string))
			Инфо("МассивПроблем %+v\n", МассивПроблем)

			for _, проблема := range МассивПроблем{
				//МассивПроблем[н]=strings.TrimSpace(проблема)
				if СтрокаПроблем == ""{
					СтрокаПроблем = `"`+strings.TrimSpace(проблема)+`"`
				} else {
					СтрокаПроблем = СтрокаПроблем +`,"`+strings.TrimSpace(проблема)+`"`
				}

			}

			СтрокаПроблем="["+СтрокаПроблем+"]"
			Инфо("СтрокаПроблем %+v\n", СтрокаПроблем)
			Значения=append(Значения, []byte(СтрокаПроблем))
			//Инфо("СтрокаПроблем %+v\n", СтрокаПроблем)
			МеткиПодстановки = МеткиПодстановки + ",$"+strconv.Itoa(len(Значения))
		}
	}
	if Метки, ЕстьМетки :=  Аргументы["skill_keywords"];ЕстьМетки {
		if len(Метки.(string)) > 0 {
			МассивМеток := strings.Split(Метки.(string),",")
			Столбцы = Столбцы + ",keywords"
			//Инфо("МассивМеток %+v\n", МассивМеток)
			СтрокаМеток := ""
			for _, метка := range МассивМеток {
				//МассивПроблем[н]=strings.TrimSpace(проблема)
				if СтрокаМеток == "" {
					СтрокаМеток = `"` + strings.TrimSpace(метка) + `"`
				} else{
					СтрокаМеток = СтрокаМеток + `,"` + strings.TrimSpace(метка) + `"`
				}

			}
			СтрокаМеток = "[" + СтрокаМеток + "]"
			//Инфо("СтрокаПроблем %+v\n", СтрокаМеток)
			Значения = append(Значения, []byte(СтрокаМеток))
			МеткиПодстановки = МеткиПодстановки + ",$" + strconv.Itoa(len(Значения))
		}
	}

	СтарыйНавык ,err:= sqlStruct{
		Name:   "io_actions",
		Sql:    "SELECT * FROM iobot.io_actions WHERE action_id = $1",
		Values: [][]byte{
			[]byte(Аргументы["skill_id"].(string)),
		},
	}.Выполнить(nil)
	if err != nil{
		Инфо(">>>> Ошибка SQL запроса: %+v \n\n",err)
	}

	//СтарыйНавыкЛог, ошибка := json.Marshal(СтарыйНавык[0])
	//if ошибка != nil {
	//	Инфо(">>>> ERROR \n %+v \n\n", ошибка)
	//}
	//Инфо("\n\n СтарыйНавыкЛог %+v\n\n\n", `[{"Дата":"`+time.Now().Format("2006-01-02 15:04:05.000000")+`", "Пользователь":`+client.Login+`,"Предыдущее":`+string(СтарыйНавыкЛог)+`}]`)

	СтарыйНавыкЛог, ошибка := json.Marshal([]map[string]interface{}{
		{"Дата":time.Now().Format("2006-01-02 15:04:05.000000"),
			"Пользователь":client.Login,
			"Предыдущее":СтарыйНавык[0],
		},
	})
	if ошибка != nil {
		Инфо(">>>> ERROR \n %+v \n\n", ошибка)
	}
	Инфо("\n\n >>>>>> СтарыйНавыкЛог %+v\n", string(СтарыйНавыкЛог))

	Столбцы = Столбцы + ",change_log"
	Значения=append(Значения, СтарыйНавыкЛог)
	МеткиПодстановки = МеткиПодстановки + ",change_log||$"+strconv.Itoa(len(Значения))

	Инфо(" Значение 8 %+v\n", string(Значения[7]))


	Инфо("\n\n SQL UPDATE  %+v\n\n", "UPDATE iobot.io_actions SET ("+Столбцы+",date_create) = ("+МеткиПодстановки+",NOW()) WHERE action_id=$1 RETURNING *")


	for i, s:=range Значения{
		Инфо("i %+v\n", i)
		Инфо("s %+v\n", string(s))
	}


	Результат ,err:= sqlStruct{
		Name:   "io_actions",
		Sql:    "UPDATE iobot.io_actions SET ("+Столбцы+",date_create) = ("+МеткиПодстановки+",NOW()) WHERE action_id=$1 RETURNING *",
		Values: Значения,
	}.Выполнить(nil)
	if err != nil{
		Инфо(">>>> Ошибка SQL запроса: %+v \n\n",err)
		СообщениеКлиенту:= &Сообщение{
			Текст:   "Ошибка сохранения навыка! " + err.Error(),
			От: "io",
			Кому:client.Login,
			MessageType:[]string{"irritation"},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
		return
		//client.СохранитьИОтправитьСообщениеКлиенту(СообщениеКлиенту)
	} else {
		if len(Результат) > 0{
			Результат[0]["readonly"]="t"
			Навык := map[string]interface{}{"Skills":Результат}
			СообщениеКлиенту:= &Сообщение{
				Текст:   "Навык изменён. " + string(render("SkillList",Навык)),
				От: "io",
				Кому:client.Login,
				//Контэнт: &ДанныеОтвета{
				//	Контейнер:"сообщение",
				//	HTML:string(render("SkillTpl", Результат[0])),
				//},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
			return
		} else {
			СообщениеКлиенту:= &Сообщение{
				Текст:   "Ошибка сохранения навыка! " + err.Error(),
				От: "io",
				Кому:client.Login,
				MessageType:[]string{"irritation"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
		}
	}
	return
}

func (client *Client)ПоказатьНавыки (mes Сообщение) map[string]interface{} {


	SQL :=`SELECT distinct iobot.io_actions.* FROM iobot.io_actions
       join fssp_configs.уровень_доступа ON iobot.io_actions.access ?| (
    select array_agg(level.уровень) from (
       select jsonb_array_elements_text(уровень) уровень  from fssp_configs.уровень_доступа where логин = $1 OR отдел = $2 OR должность = $3 OR (отдел IS NULL AND логин IS NULL) )level
)WHERE delete <> true`
	OSPCode := strconv.Itoa(client.UserInfo.Info.OspCode)
	Post := strconv.Itoa(client.UserInfo.Info.Post)
	Навыки, _ := sqlStruct{
		Name:   "io_actions",
		Sql:    SQL,
		Values: [][]byte{
			[]byte(client.Login),
			[]byte(OSPCode),
			[]byte(Post),

		},
	}.Выполнить(nil)

	Инфо(">>> ПоказатьНавыки %+v\n", Навыки)
	Навык := map[string]interface{}{"Skills":Навыки,"client":client.UserInfo.Info}
	//Навык["client"] = client.UserInfo.Info

	html:=render("SkillList", Навык)
	//Инфо(" \n \n !!!! >>>> HTML \n", html)
	return map[string]interface{}{"Контейнер":"сообщение","HTML": html}

	//return map[string]interface{}{
	//	"html":html,
	//	"target":"log_"+mes.Кому,
	//}
	//mesOut:= &Сообщение{
	//					Текст:  string(render("NewSkill", Навык)),
	//					От: "io",
	//					Кому:client.Login,
	//					MessageType:[]string{"io_action"},
	//				}
	//Инфо("НавыкНавыкНавыкНавык %+v\n", Навык)
	//mesOut.Id, mesOut.Время =  СохранитьСообщение(*mesOut)
	////удалим из памяти бота информацию об актиновм диалоге. закончили общаться
	//client.АктивныеДиалоги = nil
	//client.Message<-mesOut
	//
	//return nil
}
func (client *Client)ПоказатьНавык (mes Сообщение) map[string]interface{} {

	Навыки, _ := sqlStruct{
		Name:   "io_actions",
		Sql:    "SELECT * FROM iobot.io_actions WHERE action_id = (select  MAX(action_id) from iobot.io_actions)",
		Values: [][]byte{},
	}.Выполнить(nil)
	Инфо("НавыкНавыкНавыкНавык %+v\n", Навыки)
	Навык := Навыки[0]
	Навык["client"] = client.UserInfo.Info

	html:=string(render("NewSkill", Навык))

	return map[string]interface{}{
		"html":html,
		"target":"log_"+mes.Кому,
	}
	//mesOut:= &Сообщение{
	//					Текст:  string(render("NewSkill", Навык)),
	//					От: "io",
	//					Кому:client.Login,
	//					MessageType:[]string{"io_action"},
	//				}
	//Инфо("НавыкНавыкНавыкНавык %+v\n", Навык)
	//mesOut.Id, mesOut.Время =  СохранитьСообщение(*mesOut)
	////удалим из памяти бота информацию об актиновм диалоге. закончили общаться
	//client.АктивныеДиалоги = nil
	//client.Message<-mesOut

	//return nil
}

type СтруктураНавыка struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Actions map[string]interface{} `json:"actions"`
	Message map[string]string `json:"message"`
	ВариантыОшибок map[string]interface{} `json:"варианты_ошибок"`
}
/*
получитьНавыкДляВыполненияПоИмени Пока используется для запуска скилов из диалога

*/
func (client *Client)получитьНавыкДляВыполненияПоИмени(имя string) СтруктураНавыка {
	SQLstr := `SELECT action_id, name, actions, сообщение_при_удаче, сообщение_при_ошибке, варианты_ошибок  FROM iobot.io_actions WHERE name= $1`
	//if client.UserInfo.Info.OspCode != 26911 {
	//	SQLstr= `SELECT action_id, name, actions, сообщение_при_удаче, сообщение_при_ошибке FROM iobot.io_actions WHERE name= $1 AND self_use=true`
	//}

	Навыки, err := sqlStruct{
		Name:   `io_actions`,
		Sql:    SQLstr,
		Values: [][]byte{
			[]byte(имя),
		},
	}.Выполнить(nil)
	if err != nil {
		Инфо(">>>> ERROR \n %+v \n\n", err)
	}
	Навык := СтруктураНавыка{}
	if len(Навыки) >0 {

		Инфо("Навык.Действия[cmd] %+v\n", Навыки)

		actions := Навыки[0]["actions"].(map[string]interface{})

		cmd, ok := actions["cmd"]
		if !ok {
			Инфо(" нет cmd%+v\n", actions)
		}

		tpl := template.Must(template.New("action").Parse(cmd.(string)))
		Команда := new(bytes.Buffer)

		errExecute := tpl.Execute(Команда, client.UserInfo.Info)
		if errExecute != nil {
			Ошибка("executing template:", err)
		}
		Инфо("Команда %+v\n", Команда)
		//Навык.Действия["cmd"] = resHtml.String()

		Навык = СтруктураНавыка{
			Id:Навыки[0]["action_id"].(string),
			Actions:map[string]interface{}{"cmd":Команда.String()},
			Name:Навыки[0]["name"].(string),
			Message: map[string]string{},
			ВариантыОшибок: map[string]interface{}{},
		}
		if СообщениеПриОшибке, ЕстьСообщениеПриОшибке := Навыки[0]["сообщение_при_ошибке"];ЕстьСообщениеПриОшибке && СообщениеПриОшибке != ""{
			Навык.Message["сообщение_при_ошибке"] = СообщениеПриОшибке.(string)
		}
		if СообщениеПриУдаче, ЕстьСообщениеПриУдаче := Навыки[0]["сообщение_при_удаче"];ЕстьСообщениеПриУдаче && СообщениеПриУдаче != ""{
			Навык.Message["сообщение_при_удаче"] = СообщениеПриУдаче.(string)
		}
		if ВариантыОшибок, ЕстьВариантыОшибок := Навыки[0]["варианты_ошибок"];ЕстьВариантыОшибок && ВариантыОшибок != nil{
			Навык.ВариантыОшибок  = ВариантыОшибок.(map[string]interface{})
		}
	}


	//BotMenu := []*BotMenuStruct{}
	//ActionMap, _ := ВыполнитьPgSQL(sqlQuery)
	//Action := СтруктураНавыка{}
	//
	//for _, row := range ActionMap{
	//	for cellName, cell := range row{
	//		if cellName == "actions"{
	//			actions := map[string]interface{}{}
	//			json.Unmarshal([]byte(cell.(string)), &actions)
	//			Action.Действия=actions
	//		}
	//	}
	//}
	return Навык
}

func ЗапуститьСкрипт (client *Client, навык map[string]interface{})  {

	//skillId := mes.Выполнить.Навык.String()

	Инфо("навык %+v", навык)

	tpl := template.Must(template.New("action").Parse(навык["actions"].(map[string]interface{})["cmd"].(string)))
	Команда := new(bytes.Buffer)

	errExecute := tpl.Execute(Команда, client.UserInfo.Info)
	if errExecute != nil {
		Ошибка("executing template:", errExecute)
	}


	logCmd, errors := client.RunSkillCmd(Команда.String())
	Инфо("logCmd  %+v, errors  %+v\n", logCmd, errors )

	if errors != nil{
		//Инфо("навык.ВариантыОшибок %+v\n", навык.ВариантыОшибок)
		//for _, ошибка := range errors {
		//	Инфо("ошибка '%+v'\n", ошибка)
		//
		//	Инфо("навык.ВариантыОшибок[ошибка] %+v\n", навык.ВариантыОшибок[ошибка])
		//
		//	if ДействиеПиОшибке, ЕстьВариантОшибки := навык.ВариантыОшибок[ошибка];ЕстьВариантОшибки{
		//		Инфо("ДействиеПиОшибке %+v\n", ДействиеПиОшибке)
		//		if ДействиеПиОшибке == "далее" {
		//			return logCmd, nil
		//		}
		//	}
		//}

		if СообщениеПриОшибке, ЕстьСообщениеПриОшибке := навык["сообщение_при_ошибке"];ЕстьСообщениеПриОшибке{

			tpl := template.Must(template.New("action").Parse(СообщениеПриОшибке.(string)))
			ТекстСообщения := new(bytes.Buffer)
			errExecute := tpl.Execute(ТекстСообщения, client.UserInfo.Info)

			if errExecute != nil {
				Ошибка("executing template: %+v", errExecute)
			}
			СообщениеКлиенту:= &Сообщение{
				Текст:   ТекстСообщения.String(),
				От: "io",
				Кому:client.Login,
				MessageType: []string{"io_action"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)

		} else {
			СообщениеКлиенту:= &Сообщение{
				Текст:   "Возможно я не смог исправить Вашу проблему: \n "+ strings.Join(errors, "; "),
				От: "io",
				Кому:client.Login,
				MessageType: []string{"irritation","io_action"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
			СообщениеКлиенту= &Сообщение{
				Текст:   "Но, на всякий случай проверьте решена ли Ваша проблема",
				От: "io",
				Кому:client.Login,
				MessageType: []string{"io_action"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
		}


	} else {

		if СообщениеПриУдаче, ЕстьСообщениеПриУдаче :=  навык["сообщение_при_удаче"];ЕстьСообщениеПриУдаче{
			tpl := template.Must(template.New("action").Parse(СообщениеПриУдаче.(string)))
			ТекстСообщения := new(bytes.Buffer)
			errExecute := tpl.Execute(ТекстСообщения, client.UserInfo.Info)

			if errExecute != nil {
				Ошибка("executing template: %+v", errExecute)
			}
			СообщениеКлиенту:= &Сообщение{
				Текст:   ТекстСообщения.String(),
				От: "io",
				Кому:client.Login,
				MessageType: []string{"io_action"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)

		} else {
			СообщениеКлиенту:= &Сообщение{
				Текст:   "Проверьте решена ли Ваша проблема",
				От: "io",
				Кому:client.Login,
				MessageType: []string{"io_action"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
		}

	}

	client.СохранитьЛог(навык["name"].(string), logCmd, errors)

	if errors != nil {
		Инфо("errors %+v\n", errors)
	}
	Инфо("навык %+v logCmd: %+v\n", навык, logCmd)
	//return logCmd, errors
}

func (client *Client)ЗапуститьНавыкПоИмени(имя string) (string, []string) {

	//skillId := mes.Выполнить.Навык.String()
	навык := client.получитьНавыкДляВыполненияПоИмени(имя)

	logCmd, errors := client.RunSkillCmd(навык.Actions["cmd"].(string))
	Инфо("logCmd  %+v, errors  %+v\n", logCmd, errors )

	if errors != nil{
		Инфо("навык.ВариантыОшибок %+v\n", навык.ВариантыОшибок)
		for _, ошибка := range errors {
			Инфо("ошибка '%+v'\n", ошибка)

			Инфо("навык.ВариантыОшибок[ошибка] %+v\n", навык.ВариантыОшибок[ошибка])

			if ДействиеПиОшибке, ЕстьВариантОшибки := навык.ВариантыОшибок[ошибка];ЕстьВариантОшибки{
				Инфо("ДействиеПиОшибке %+v\n", ДействиеПиОшибке)
				if ДействиеПиОшибке == "далее" {
					return logCmd, nil
				}
			}
		}

		if СообщениеПриОшибке, ЕстьСообщениеПриОшибке := навык.Message["сообщение_при_ошибке"];ЕстьСообщениеПриОшибке{

			tpl := template.Must(template.New("action").Parse(СообщениеПриОшибке))
			ТекстСообщения := new(bytes.Buffer)
			errExecute := tpl.Execute(ТекстСообщения, client.UserInfo.Info)

			if errExecute != nil {
				Ошибка("executing template: %+v", errExecute)
			}
			СообщениеКлиенту:= &Сообщение{
				Текст:   ТекстСообщения.String(),
				От: "io",
				Кому:client.Login,
				MessageType: []string{"io_action"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)

		} else {
			СообщениеКлиенту:= &Сообщение{
				Текст:   "Возможно я не смог исправить Вашу проблему: \n "+ strings.Join(errors, "; "),
				От: "io",
				Кому:client.Login,
				MessageType: []string{"irritation","io_action"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
			СообщениеКлиенту= &Сообщение{
				Текст:   "Но, на всякий случай проверьте решена ли Ваша проблема",
				От: "io",
				Кому:client.Login,
				MessageType: []string{"io_action"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
		}


	} else {

		if СообщениеПриУдаче, ЕстьСообщениеПриУдаче :=  навык.Message["сообщение_при_удаче"];ЕстьСообщениеПриУдаче{
			tpl := template.Must(template.New("action").Parse(СообщениеПриУдаче))
			ТекстСообщения := new(bytes.Buffer)
			errExecute := tpl.Execute(ТекстСообщения, client.UserInfo.Info)

			if errExecute != nil {
				Ошибка("executing template: %+v", errExecute)
			}
			СообщениеКлиенту:= &Сообщение{
				Текст:   ТекстСообщения.String(),
				От: "io",
				Кому:client.Login,
				MessageType: []string{"io_action"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)

		} else {
			СообщениеКлиенту:= &Сообщение{
				Текст:   "Проверьте решена ли Ваша проблема",
				От: "io",
				Кому:client.Login,
				MessageType: []string{"io_action"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
		}

	}

	client.СохранитьЛог(навык.Id, logCmd, errors)

	if errors != nil {
		Инфо("errors %+v\n", errors)
	}
	Инфо("навык %+v logCmd: %+v\n", навык, logCmd)
	return logCmd, errors
}

func (client *Client)получитьНавыкДляВыполнения(skillId string) СтруктураНавыка {
	SQLstr := `SELECT action_id, name, actions , сообщение_при_удаче, сообщение_при_ошибке FROM iobot.io_actions WHERE action_id= $1`
	if client.UserInfo.Info.OspCode != 26911 {
		SQLstr= `SELECT action_id, name, actions , сообщение_при_удаче, сообщение_при_ошибке FROM iobot.io_actions WHERE action_id= $1 AND self_use=true`
	}

	Навыки, err := sqlStruct{
		Name:   `io_actions`,
		Sql:    SQLstr,
		Values: [][]byte{
			[]byte(skillId),
		},
	}.Выполнить(nil)
	if err != nil {
		Инфо(">>>> ERROR \n %+v \n\n", err)
	}
	Навык := СтруктураНавыка{}
	if len(Навыки) >0 {

		Инфо("Навык.Действия[cmd] %+v\n", Навыки)

		actions := Навыки[0]["actions"].(map[string]interface{})

		cmd, ok := actions["cmd"]
		if !ok {
			Инфо(" нет cmd%+v\n", actions)
		}

		tpl := template.Must(template.New("action").Parse(cmd.(string)))
		Команда := new(bytes.Buffer)

		errExecute := tpl.Execute(Команда, client.UserInfo.Info)
		if errExecute != nil {
			Ошибка("executing template: %+v", err)
		}
		Инфо("Команда %+v\n", Команда)
		//Навык.Действия["cmd"] = resHtml.String()

		Навык = СтруктураНавыка{
			Id:Навыки[0]["action_id"].(string),
			Actions:map[string]interface{}{"cmd":Команда.String()},
			Name:Навыки[0]["name"].(string),
			Message: map[string]string{},
		}
		if СообщениеПриОшибке, ЕстьСообщениеПриОшибке := Навыки[0]["action_id"];ЕстьСообщениеПриОшибке{
			Навык.Message["сообщение_при_ошибке"] = СообщениеПриОшибке.(string)
		}
		if СообщениеПриУдаче, ЕстьСообщениеПриУдаче := Навыки[0]["action_id"];ЕстьСообщениеПриУдаче{
			Навык.Message["сообщение_при_удаче"] = СообщениеПриУдаче.(string)

		}
	}

	//tions.go:1149: >>>>> СоздатьНавык mes {Token:{Hash: Истекает:} Id:10303 Ip: От:maksimchuk@r26 Кому:io Текст: MessageType:[] Время:2020-02-11T16:14:57.477065 ОтветНа: Файлы:[] Offline: Online: AdminMenu:[] ВходящиеАргументы:map[] Выполнить:{Action: Действие:map[СоздатьНавык:map[skill_action:killall soffice.bin;rm -f /home/${uid}/.config/libreoffice/4/.lock; skill_description:убивает все процессы soffice.bin, удаляет /home/${uid}/.config/libreoffice/4/.lock  skill_keywords: OpenOffice. LibreOffice skill_name:Сброс процессов OpenOffice skill_problems:Не открываются документы, зависает OpenOffice/LibreOffice, не запускается  OpenOffice/LibreOffice skill_self_use:on]] Skill:0 Навык: Cmd: Комманду: Arg:{Module: Tables:[] Login: Other:map[]}} Контэнт:<nil> Content:{Target: Data:<nil> Html: Обработчик:} UserInfo:{Uid: Initials: FullName: FirstName: LastName: MiddleName: OspName: OspNum:0 PostName: Инициалы: ПолноеИмя: Фамилия: Имя: Отчество: ОСП: КодОСП:0 Должность:}}
	//16:14:57 Действия.go:1256: >>>> Ошибка SQL запроса: ОШИБКА: неверный синтаксис для типа json (SQLSTATE 22P02)



	//BotMenu := []*BotMenuStruct{}
	//ActionMap, _ := ВыполнитьPgSQL(sqlQuery)
	//Action := СтруктураНавыка{}
	//
	//for _, row := range ActionMap{
	//	for cellName, cell := range row{
	//		if cellName == "actions"{
	//			actions := map[string]interface{}{}
	//			json.Unmarshal([]byte(cell.(string)), &actions)
	//			Action.Действия=actions
	//		}
	//	}
	//}
	return Навык
}

func (client *Client)ЗапуститьНавык(mes Сообщение){

	skillId := mes.Выполнить.Навык.String()
	навык := client.получитьНавыкДляВыполнения(skillId)

	logCmd, errors := client.RunSkillCmd(навык.Actions["cmd"].(string))


	if errors != nil{
		СообщениеКлиенту:= &Сообщение{
			Текст:   "Возможно я не смог исправить Вашу проблему: \n "+ strings.Join(errors, "; "),
			От: "io",
			Кому:client.Login,
			MessageType: []string{"irritation","io_action"},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
		СообщениеКлиенту= &Сообщение{
			Текст:   "Но, на всякий случай проверьте решена ли Ваша проблема",
			От: "io",
			Кому:client.Login,
			MessageType: []string{"io_action"},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
	} else {
		СообщениеКлиенту:= &Сообщение{
			Текст:   "Проверьте решена ли Ваша проблема",
			От: "io",
			Кому:client.Login,
			MessageType: []string{"io_action"},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
	}

	client.СохранитьЛог(skillId, logCmd, errors)

	if errors != nil {
		Инфо("errors %+v\n", errors)
	}
	Инфо("навык %+v logCmd: %+v\n", навык, logCmd)
}

//func (client *Client)ИзменитьНавык (mes Сообщение) map[string]interface{} {
//
//	Аргументы := mes.ВходящиеАргументы
//	if Аргументы !=nil  {
//		if actionId , ok := Аргументы["action_id"];ok {
//
//			ИзменяемыйНавык ,err:= sqlStruct{
//					Name:   "io_actions",
//					Sql:    "SELECT * FROM iobot.io_actions WHERE action_id = $1",
//					Values: [][]byte{
//						[]byte(actionId.(string)),
//					},
//				}.Выполнить(nil)
//			if err != nil{
//			Инфо(">>>> Ошибка SQL запроса: %+v \n\n",err)
//			}
//
//			Навык := map[string]interface{}{"Skill":ИзменяемыйНавык,"client":client.UserInfo.Info}
//			//Навык["client"] = client.UserInfo.Info
//
//			html:=render("EditSkillForm" ,Навык)
//
//			return map[string]interface{}{"Контейнер":"сообщение","HTML": html}
//
//		}
//	} else {
//		СообщениеКлиенту:= &Сообщение{
//					Текст:   "Не пойму какой навык Вы хотите изменить",
//					От: "io",
//					Кому:client.Login,
//				}
//				СообщениеКлиенту.СохранитьЛогСообщения()
//				client.Message<-СообщениеКлиенту
//				return nil
//	}
//	return nil
//}
//func (client *Client)ИзменитьНавык_Диалог (mes Сообщение){
//	/*
//	 Выполнить: {"Действие": {
//		"ИзменитьНавык": {
//		"НомерНавыка": id,
//		"ИзменяемоеПоле":editFieldName,
//		"НовоеЗначение":event.target.innerText,
//		"СтароеЗначение":СтароеЗначение}
//	}
//	},
//	*/
//	 Аргументы, ok :=mes.Выполнить.Действие["ИзменитьНавык"]
//	if !ok{
//		Инфо("Для изменения навыка не хватает аргументов в ok=%+v, аргументе = %+v\n", ok, Аргументы )
//		mesOut:= &Сообщение{
//			Текст:  "Для изменения навыка не хватает аргументов",
//			От: "io",
//			Кому:client.Login,
//			MessageType:[]string{"io_action"},
//		}
//		mesOut.Id, mesOut.Время = СохранитьСообщение(*mesOut)
//		client.Message<-mesOut
//		return
//	}
//
//	// для обновления навыка отправим запрос в котором оновлять будем тот навыык в котором совпадает номер навыка и старое значение изменяемого поля
//	var СтароеЗначение string
//	var НовоеЗначение string
//	var НомерНавыка string
//	var ИзменяемоеПоле string
//	if _, ок := Аргументы["НомерНавыка"];ок{
//		НомерНавыка = Аргументы["НомерНавыка"].(string)
//	} else {
//		Инфо("НомерНавыка пусто: %+v\n", НомерНавыка, ок)
//	}
//	if _, ок := Аргументы["ИзменяемоеПоле"];ок{
//		ИзменяемоеПоле = Аргументы["ИзменяемоеПоле"].(string)
//	} else {
//		Инфо("ИзменяемоеПоле пусто: %+v\n", ИзменяемоеПоле, ок)
//	}
//
//	if _, ок := Аргументы["СтароеЗначение"];ок{
//		if ИзменяемоеПоле == "self_use"{
//			СтароеЗначение = strings.TrimSpace(strconv.FormatBool(Аргументы["СтароеЗначение"].(bool)))
//		} else {
//			СтароеЗначение = strings.TrimSpace(Аргументы["СтароеЗначение"].(string))
//		}
//
//	} else {
//		Инфо("СтароеЗначение пусто: %+v\n", СтароеЗначение, ок)
//	}
//	if _, ок := Аргументы["НовоеЗначение"];ок{
//
//		if ИзменяемоеПоле == "self_use"{
//			НовоеЗначение = strings.TrimSpace(strconv.FormatBool(Аргументы["НовоеЗначение"].(bool)))
//		} else {
//			НовоеЗначение = strings.TrimSpace(Аргументы["НовоеЗначение"].(string))
//		}
//	} else {
//		Инфо("НовоеЗначение пусто: %+v\n", НовоеЗначение, ок)
//	}
//
//
//	switch ИзменяемоеПоле{
//		case "actions_cmd":
//			ИзменяемоеПоле="actions"
//			НовоеЗначение = `{"cmd":"`+НовоеЗначение+`"}`
//			СтароеЗначение = `{"cmd":"`+СтароеЗначение+`"}`
//			break
//		case "actions_file":
//			ИзменяемоеПоле="actions"
//			НовоеЗначение = `{"file":"`+НовоеЗначение+`"}`
//			СтароеЗначение = `{"file":"`+СтароеЗначение+`"}`
//			break
//		case "keywords":
//			Тэги := strings.Split(НовоеЗначение, ",")
//			for н, тэг := range Тэги{
//				Тэги[н]=strings.TrimSpace(тэг)
//			}
//			НовоеЗначениеСтрокаБайт, err:= json.Marshal(Тэги)
//			if err != nil {
//				Ошибка("\n !! ERR %+v\n", err)
//			}
//			НовоеЗначение = string(НовоеЗначениеСтрокаБайт)
//
//			ТэгиСтароеЗначение := strings.Split(СтароеЗначение, ",")
//			for н, тэг := range ТэгиСтароеЗначение{
//				ТэгиСтароеЗначение[н]=strings.TrimSpace(тэг)
//			}
//			СтароеЗначениеСтрокаБайт, err:= json.Marshal(ТэгиСтароеЗначение)
//			if err != nil {
//				Ошибка("\n !! ERR %+v\n", err)
//			}
//			СтароеЗначение = string(СтароеЗначениеСтрокаБайт)
//			break
//		case"problem":
//			Проблемы := strings.Split(НовоеЗначение, ";")
//			for н, проблема := range Проблемы{
//				Проблемы[н]=strings.TrimSpace(проблема)
//			}
//			НовоеЗначениеСтрокаБайт, err:= json.Marshal(Проблемы)
//			if err != nil {
//				Ошибка("\n !! ERR %+v\n", err)
//			}
//			НовоеЗначение = string(НовоеЗначениеСтрокаБайт)
//
//			//if СтароеЗначение.String != ""{
//				ПроблемыСтароеЗначение := strings.Split(СтароеЗначение, ";")
//				for н, проблема := range ПроблемыСтароеЗначение{
//					ПроблемыСтароеЗначение[н]=strings.TrimSpace(проблема)
//				}
//				СтароеЗначениеСтрокаБайт, err:= json.Marshal(ПроблемыСтароеЗначение)
//
//				if err != nil {
//					Инфо("\n !! ERR %+v\n", err)
//				}
//				СтароеЗначение = string(СтароеЗначениеСтрокаБайт)
//			//}
//			break
//		case "dir":
//		default:
//
//	}
//
//	Инфо("\n НомерНавыка %+v\n\n НовоеЗначение %+v\n\n СтароеЗначение %+v\n", НомерНавыка, НовоеЗначение, СтароеЗначение)
//	Инфо("ИзменяемоеПоле %+v\n", ИзменяемоеПоле)
//
//	//SQL := "UPDATE iobot.io_actions SET "+ИзменяемоеПоле+"=$2 WHERE action_id=$1 AND "+ИзменяемоеПоле+"= $3 RETURNING *"
//	SQLstr :=  "UPDATE iobot.io_actions SET "+ИзменяемоеПоле+"=$2 WHERE action_id=$1 AND "+ИзменяемоеПоле+"= $3 RETURNING "+ИзменяемоеПоле
//	Инфо("SQLstr %+v\n", SQLstr)
//	ИзменённыйНавык,_:=sqlStruct{
//						Name:   "io_action",
//						Sql:    SQLstr,
//						Values: [][]byte{
//							[]byte(НомерНавыка),
//							[]byte(НовоеЗначение),
//							[]byte(СтароеЗначение),
//						},
//						}.Выполнить(nil)
//
//	Инфо(" len(ИзменённыйНавык)>0 %+v\n", len(ИзменённыйНавык))
//	if len(ИзменённыйНавык)>0{
//		ИзменённыйНавык[0]["client"] = client.UserInfo.Info
//		var TargetData map[string]interface{}
//
//		if ИзменяемоеПоле == "self_use"{
//			TargetData = map[string]interface{}{"checked":ИзменённыйНавык[0][ИзменяемоеПоле]}
//		} else {
//			TargetData = map[string]interface{}{"innerText":ИзменённыйНавык[0][ИзменяемоеПоле]}
//		}
//
//		mesOut := &Сообщение{
//			//Текст: "Навык изменён: "+ string(render("NewSkill", ИзменённыйНавык[0])),
//			Текст: "Навык изменён",
//			Кому:client.Login,
//			От:"io",
//			Content: struct {
//				Target string `json:"target"`
//				Data interface{} `json:"data"`
//				Html string `json:"html"`
//				Обработчик string `json:"обработчик"`
//			}{
//				Target: "skill_"+ИзменяемоеПоле+"_"+НомерНавыка,
//				Data: TargetData,
//			},
//		}
//
//		mesOut.Id, mesOut.Время =  СохранитьСообщение(*mesOut)
//		client.Message<-mesOut
//	}
//
//
//}

func (client *Client)EndSkillCreate(mes Сообщение) (map[string]interface{}) {

	return nil
}
func (client *Client)AddSkillTags(mes Сообщение) (map[string]interface{}) {


	if mes.Текст == "" {
		return nil
	}

	Тэги := strings.Split(mes.Текст, ",")
	for i, тэг := range Тэги{
		Тэги[i]=strings.TrimSpace(тэг)
	}
	jsonTags, err := json.Marshal(Тэги)
	if err != nil {
		Инфо("err	 %+v\n", err)
	}
	ОбновлённыйНавык , _ := ВыполнитьPgSQL(sqlStruct{
		Name:   "io_actions",
		Sql:    "UPDATE iobot.io_actions SET keywords = $1  WHERE action_id = $2 RETURNING *",
		Values: [][]byte{
			jsonTags,
			[]byte(client.АктивныеДиалоги[0]["status"].(map[string]interface{})["result"].(map[string]interface{})["io_action"].(string)),
		},
	})

	if len(ОбновлённыйНавык)<1{
		//err = errors.New("Не удалось создать новый навык, ")
		ОбновлённыйНавык= []map[string]interface{}{0:{"error":"Не удалось создать новый навык, веорятно имя уже существует :"+mes.Текст}}
	} else {
		ОбновлённыйНавык[0]["result"]=map[string]interface{}{
			"io_action":ОбновлённыйНавык[0]["action_id"].(string),
		}
	}

	return ОбновлённыйНавык[0]
}


/*
Сообщение.ПолучитьДейсвтие() получает данные из сообщения  mes.Выполнить.Действие ["ИмяДействия"]
*/
func (mes Сообщение) ПолучитьДействие(Название string) (map[string]interface{} , bool){
	if Название == ""{
		return nil, false
	}
	Аргументы, ok :=mes.Выполнить.Действие[Название]
	return Аргументы, ok
}

/*
Сообщение.ДобавитьДейсвтие() добавялет данные в сообщение  mes.Выполнить.Действие ["ИмяДействия"]
//Выполнить: {"Действие": {
	//	"ИзменитьНавык": {
	//			"НомерНавыка": id,
	//			"ИзменяемоеПоле":editFieldName,
	//			"НовоеЗначение":event.target.innerText,
	//			"СтароеЗначение":СтароеЗначение}
	//}
	//},
*/

func (mes Сообщение) ДобавитьДейсвтие(Название string, НовоеДействие map[string]interface{}) bool {

	if Название == ""{
		return false
	}
	if НовоеДействие != nil{
		if mes.Выполнить.Действие == nil {
			mes.Выполнить.Действие = map[string]map[string]interface{}{Название: НовоеДействие}
		} else {
			mes.Выполнить.Действие[Название]=НовоеДействие
		}
	}
	return true
}

/*
ПолучитьНавык получает данные навыков, не рендерит, возвращает массив навыков
Выполнить: {"Действие": {
			"ПолучитьНавыки": {
				"НомераНавыков": [id...],
				"Описание":"описание содержит",
				"Название":"содержит слова"
				"Метки":[],
				"Проблемы":[],
				"Создатель":[],
				"РазрешеноВсем":true/false/nil
			]
*/

func (client *Client)ShowNewSkill (mes Сообщение) (map[string]interface{}) {
	//Инфо("\n ShowNewSkill %+v\n", client.АктивныеДиалоги)

	IdНовогоНавыка:=client.АктивныеДиалоги[0]["status"].(map[string]interface{})["result"].(map[string]interface{})["io_action"].(string)

	НовыйНавык, err :=sqlStruct{
		Name:   "io_actions",
		Sql:    "SELECT * FROM iobot.io_actions WHERE action_id = $1",
		Values: [][]byte{
			[]byte(IdНовогоНавыка),
		},
	}.Выполнить(nil)
	if err != nil {
		Инфо("\n !! ERR НовыйНавык %+v\n", err)
		return nil
	}
	//удалим из памяти бота информацию об актиновм диалоге. закончили общаться
	client.АктивныеДиалоги = nil

	if len(НовыйНавык)>0{
		НовыйНавык[0]["client"] = client.UserInfo.Info
		//Навык := НовыйНавык[0]
		//
		//КомандыНавыка := map[string]interface{}{}
		//
		//err := json.Unmarshal([]byte(Навык["actions"].(string)), &КомандыНавыка)
		//
		//if err != nil {
		//	Инфо("\n !! ERR >> %+v\n\n ", err)
		//} else {
		//	Навык["actions"] = КомандыНавыка
		//}
		//
		//var РешаемыеПроблемы []string
		//
		//
		//err = json.Unmarshal([]byte(Навык["problem"].(string)), &РешаемыеПроблемы)
		//if err != nil {
		//	Инфо("\n !! ERR >> %+v\n\n ", err)
		//} else {
		//	Навык["problem"] = strings.Join(РешаемыеПроблемы, "; ")
		//}
		//
		//var КлючевыеСлова []string
		//err = json.Unmarshal([]byte(Навык["keywords"].(string)), &КлючевыеСлова)
		//if err != nil {
		//	Инфо("\n !! ERR >> %+v\n\n ", err)
		//} else {
		//	Навык["keywords"] = strings.Join(КлючевыеСлова, ", ")
		//}

		//render("NewSkill",НовыйНавык[0] )
		return map[string]interface{}{"html":render("NewSkill",НовыйНавык[0])}

		//return map[string]interface{}{"data":НовыйНавык[0],"render":"NewSkill"}

	} else {
		return nil
	}
}


func (client *Client)AddProblemDescription(mes Сообщение) (map[string]interface{}) {


	if mes.Текст == "" {
		return nil
	}

	Проблемы := strings.Split(mes.Текст, ";")

	for i, проблема := range Проблемы{
		Проблемы[i]=strings.TrimSpace(проблема)
	}
	jsonПроблемы, err := json.Marshal(Проблемы)
	if err != nil {
		Инфо("err	 %+v\n", err)
	}
	ОбновлённыйНавык , _ := ВыполнитьPgSQL(sqlStruct{
		Name:   "io_actions",
		Sql:    "UPDATE iobot.io_actions SET problem = $1  WHERE action_id = $2 RETURNING *",
		Values: [][]byte{
			jsonПроблемы,
			[]byte(client.АктивныеДиалоги[0]["status"].(map[string]interface{})["result"].(map[string]interface{})["io_action"].(string)),
		},
	})

	if len(ОбновлённыйНавык)<1{
		//err = errors.New("Не удалось создать новый навык, ")
		ОбновлённыйНавык= []map[string]interface{}{0:{"error":"Не удалось добавить описание проблем :"+mes.Текст}}
	} else {
		ОбновлённыйНавык[0]["result"]=map[string]interface{}{
			"io_action":ОбновлённыйНавык[0]["action_id"].(string),
		}
	}

	return ОбновлённыйНавык[0]
}
func (client *Client)AllowSelfUse(mes Сообщение) (map[string]interface{}) {


	if mes.Текст == "" {
		return map[string]interface{}{"error":"Не удалось разрешить самсотоятельное выполнение скрипта: "+mes.Текст}
	}

	ОбновлённыйНавык , _ := ВыполнитьPgSQL(sqlStruct{
		Name:   "io_actions",
		Sql:    "UPDATE iobot.io_actions SET self_use = true  WHERE action_id = $1 RETURNING *",
		Values: [][]byte{
			[]byte(client.АктивныеДиалоги[0]["status"].(map[string]interface{})["result"].(map[string]interface{})["io_action"].(string)),
		},
	})

	if len(ОбновлённыйНавык)<1{
		//err = errors.New("Не удалось создать новый навык, ")
		ОбновлённыйНавык= []map[string]interface{}{
			0:{
				"error":"Не удалось разрешить самсотоятельное выполнение скрипта:"+mes.Текст,
			},
		}
	} else {
		ОбновлённыйНавык[0]["result"]=map[string]interface{}{
			"io_action":ОбновлённыйНавык[0]["action_id"].(string),
		}
	}

	return ОбновлённыйНавык[0]
}
func (client *Client)CreateSkillCmd(mes Сообщение) (map[string]interface{}) {
	action_id, _ := ВыполнитьPgSQL(sqlStruct{
		Name:   "io_actions",
		Sql:    "SELECT status->>'result' as result FROM iobot.io_dialogs_log WHERE uid = $1 AND status->>'завершено' IS NULL",
		Values: [][]byte{
			[]byte(client.Login),
		},
	})
	var ActionId string

	if len(action_id[0])>0{
		if action_id[0]["result"] != nil{
			var ActionIdMap map[string]string
			err := json.Unmarshal([]byte(action_id[0]["result"].(string)), &ActionIdMap)
			if err != nil {
				Ошибка("err	 %+v\n", err)
			}
			Инфо("ActionIdMap %+v\n", ActionIdMap)
			//ActionId = action_id[0]["result"].(map[string]string)["action_id"]
			if ActionIdMap["io_action"]!=""{
				ActionId=ActionIdMap["io_action"]
			}
		}

	} else {

		КомандаНавыка := map[string]interface{}{"error":"Не удалось добавить описание навыку :"+mes.Текст}
		return КомандаНавыка
	}

	cmd,err := json.Marshal(map[string]string{"cmd":strings.TrimSpace(mes.Текст)})
	if err != nil {
		Инфо("err	 %+v\n", err)
	}
	КомандаНавыка, _ := ВыполнитьPgSQL(sqlStruct{
		Name:   "io_actions",
		Sql:    "UPDATE iobot.io_actions SET actions = $1  WHERE action_id = $2 RETURNING *",
		Values: [][]byte{
			cmd,
			[]byte(ActionId),
		},
	})

	//Инфо("КомандаНавыка %+v\n", КомандаНавыка)

	if len(КомандаНавыка)<1{
		КомандаНавыка[0]= map[string]interface{}{"error":"Не удалось добавить описание навыку :"+mes.Текст}
	} else {
		КомандаНавыка[0]["result"]=map[string]interface{}{
			"io_action":КомандаНавыка[0]["action_id"].(string),
		}
	}

	return КомандаНавыка[0]
}
func (client *Client)CreateSkillDescription(mes Сообщение) (map[string]interface{}) {
	action_id, _ := ВыполнитьPgSQL(sqlStruct{
		Name:   "io_actions",
		Sql:    "SELECT status->>'result' as result FROM iobot.io_dialogs_log WHERE uid = $1 AND status->>'завершено' IS NULL",
		Values: [][]byte{
			[]byte(client.Login),
		},
	})
	Инфо("io_action %+v\n", action_id)
	var ActionId string

	if len(action_id[0])>0{
		if action_id[0]["result"] != nil{
			var ActionIdMap map[string]string
			err := json.Unmarshal([]byte(action_id[0]["result"].(string)), &ActionIdMap)
			if err != nil {
				Ошибка("err	 %+v\n", err)
			}
			Инфо("ActionIdMap %+v\n", ActionIdMap)
			//ActionId = action_id[0]["result"].(map[string]string)["action_id"]
			if ActionIdMap["io_action"]!=""{
				ActionId=ActionIdMap["io_action"]
			}
		}

	} else {

		ОписаниеНавыка := map[string]interface{}{"error":"Не удалось добавить описание навыку :"+mes.Текст}
		return ОписаниеНавыка
	}

	Инфо("ActionId %+v\n", ActionId)
	ОписаниеНавыка, _ := ВыполнитьPgSQL(sqlStruct{
		Name:   "io_actions",
		Sql:    "UPDATE iobot.io_actions SET description = $1  WHERE action_id = $2 RETURNING *",
		Values: [][]byte{
			[]byte(mes.Текст),
			[]byte(ActionId),
		},
	})

	Инфо("ОписаниеНавыка %+v\n", ОписаниеНавыка)

	if len(ОписаниеНавыка)<1{
		ОписаниеНавыка= []map[string]interface{}{0:{"error":"Не удалось добавить описание навыку :"+mes.Текст}}
	} else {
		ОписаниеНавыка[0]["result"]=map[string]interface{}{
			"io_action":ОписаниеНавыка[0]["action_id"].(string),
		}
	}

	return ОписаниеНавыка[0]
}
func (client *Client)CreateSkillName(mes Сообщение) (map[string]interface{}) {

	НовыйНавык, _ := ВыполнитьPgSQL(sqlStruct{
		Name:   "io_actions",
		Sql:    "INSERT INTO iobot.io_actions (name, creator,date_create) VALUES ($1,$2, NOW()) ON CONFLICT (name) DO NOTHING RETURNING *",
		Values: [][]byte{
			[]byte(mes.Текст),
			[]byte(client.Login),
		},
	})
	//var err error

	if len(НовыйНавык)<1{
		//err = errors.New("Не удалось создать новый навык, ")
		НовыйНавык= []map[string]interface{}{0:{"error":"Не удалось создать новый навык, веорятно имя уже существует :"+mes.Текст}}
	} else {
		НовыйНавык[0]["result"]=map[string]interface{}{
			"io_action":НовыйНавык[0]["action_id"].(string),
		}
	}

	return НовыйНавык[0]
}

