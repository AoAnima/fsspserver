package main

import (
	"encoding/json"
	"log"
	"path/filepath"
	"plugin"
	"strconv"
	"strings"
)


/*
Actions
каждый Actions должен возвращать самостоятельно ответ пользователю
client.Message<-&Сообщение{}

Можно считать что каждое действие атомарно, и должно само позаботиться о доставке ответа клиенту

*/

var Actions = map[string]interface{}{
	"getChatLog": (*Client).ПолучитьЛогПереписки,
	"getIoMenu": (*Client).ПолучитьМенюБота,
	"collectData":(*Client).СобратьДанные,
	"GetData":(*Client).ЗагрузитьДанные,
	"SSHConnect":(*Client).SSHConnect,
	"CloseSSH":(*Client).СloseSSH,
	"CreateSkillName":(*Client).CreateSkillName,
	"CreateSkillDescription":(*Client).CreateSkillDescription,
	"CreateSkillCmd":(*Client).CreateSkillCmd,
	"AddSkillTags":(*Client).AddSkillTags,
	"EndSkillCreate":(*Client).EndSkillCreate,
	//"SkillTags":(*Client).SkillTags,
	//"GetData":DataCollector,
	//"CollectData": DataCollector,
}

func Plugin (client * Client , mes Сообщение) {
	//mes.Выполнить.Action
	DataCollector, err := plugin.Open("./plugins/DataCollector.so")
	if err != nil {
		log.Printf("Ошибка загрузки плагина DataCollector %+v\n", err)
	}

	ActionFunction, err := DataCollector.Lookup(mes.Выполнить.Action)
	if err != nil {
		log.Printf("err	 %+v\n", err)
	}
	ActionFunction.(func(*Client, Сообщение))(client, mes)
}

func (client * Client)Plugins (mes Сообщение){
	PluginsList, err := filepath.Glob("./plugins/*.so")
	if err != nil {
		log.Printf("ошибка открытия плагинов %+v\n", err)
		mes := Сообщение{
			Текст: err.Error(),
			От: "io",
			Кому:client.Login,
			Id:-2, // -2 ошибка
		}
		client.Message <- &mes
	}
	for _, filename := range PluginsList{

		pluginFile, err := plugin.Open(filename)
		if err != nil {
			log.Printf("ошибка открытия плагинов %+v\n", err)
			mes := Сообщение{
				Текст: "Не возможно загрузить плагины",
				От: "io",
				Кому:client.Login,
				Id:-2, // -2 ошибка
			}
			client.Message <- &mes
		}
		log.Printf("mes.Выполнить.Action %+v\n", mes.Выполнить.Action)
		Action, err := pluginFile.Lookup(mes.Выполнить.Action)
		log.Printf("Action %+v\n", Action)
		if err != nil {
			mes := Сообщение{
				Текст: "Действие не существует",
				От: "io",
				Кому:client.Login,
				Id:-2, // -2 ошибка
			}
			client.Message <- &mes
			return
		}




		argc := map[string]interface{}{
			"Clinet":client,
			"Message":&mes,
		}
		Action.(func(map[string]interface{}))(argc)
	}
}



func (client *Client)EndSkillCreate(mes Сообщение) (map[string]interface{}) {

	return nil
}
func (client *Client)AddSkillTags(mes Сообщение) (map[string]interface{}) {

	return nil
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
				log.Printf("err	 %+v\n", err)
			}
			log.Printf("ActionIdMap %+v\n", ActionIdMap)
			//ActionId = action_id[0]["result"].(map[string]string)["action_id"]
			if ActionIdMap["io_action"]!=""{
				ActionId=ActionIdMap["io_action"]
			}
		}

	} else {

		КомандаНавыка := map[string]interface{}{"error":"Не удалось добавить описание навыку :"+mes.Текст}
		return КомандаНавыка
	}

	cmd,err := json.Marshal(map[string]string{"cmd":mes.Текст})
	if err != nil {
		log.Printf("err	 %+v\n", err)
	}
	КомандаНавыка, _ := ВыполнитьPgSQL(sqlStruct{
		Name:   "io_actions",
		Sql:    "UPDATE iobot.io_actions SET actions = $1  WHERE id = $2 RETURNING *",
		Values: [][]byte{
			cmd,
			[]byte(ActionId),
		},
	})

	//log.Printf("КомандаНавыка %+v\n", КомандаНавыка)

	if len(КомандаНавыка)<1{
		КомандаНавыка[0]= map[string]interface{}{"error":"Не удалось добавить описание навыку :"+mes.Текст}
	} else {
		КомандаНавыка[0]["result"]=map[string]interface{}{
			"io_action":КомандаНавыка[0]["id"].(string),
		}
	}

	return КомандаНавыка[0]
}

var АктивныеДиалоги []map[string]interface{}

func (client *Client)ПолучитьАктивныйДиалог (mes Сообщение){
	// Sql:    "SELECT * FROM iobot.io_dialogs_log WHERE uid = $1 AND status->>'завершено' IS NULL",
	АктивныеДиалоги,_ = ВыполнитьPgSQL(sqlStruct{
		Name:   "dialogs_log",
		Sql:    "SELECT * FROM iobot.io_dialogs_log LEFT JOIN iobot.io_dialogs ON iobot.io_dialogs_log.dialog_id = iobot.io_dialogs.dialog_id AND iobot.io_dialogs_log.message_id = iobot.io_dialogs.message_id WHERE uid = $1 AND status->>'завершено' IS NULL",
		Values: [][]byte{
			[]byte(client.Login),
		},
	})
}

func (client *Client)CreateSkillDescription(mes Сообщение) (map[string]interface{}) {


	action_id, _ := ВыполнитьPgSQL(sqlStruct{
		Name:   "io_actions",
		Sql:    "SELECT status->>'result' as result FROM iobot.io_dialogs_log WHERE uid = $1 AND status->>'завершено' IS NULL",
		Values: [][]byte{
			[]byte(client.Login),
		},
	})
	log.Printf("io_action %+v\n", action_id)
	var ActionId string

	if len(action_id[0])>0{
		if action_id[0]["result"] != nil{
			var ActionIdMap map[string]string
			err := json.Unmarshal([]byte(action_id[0]["result"].(string)), &ActionIdMap)
			if err != nil {
				log.Printf("err	 %+v\n", err)
			}
			log.Printf("ActionIdMap %+v\n", ActionIdMap)
			//ActionId = action_id[0]["result"].(map[string]string)["action_id"]
			if ActionIdMap["io_action"]!=""{
				ActionId=ActionIdMap["io_action"]
			}
		}

	} else {

		ОписаниеНавыка := map[string]interface{}{"error":"Не удалось добавить описание навыку :"+mes.Текст}
		return ОписаниеНавыка
	}

log.Printf("ActionId %+v\n", ActionId)
	ОписаниеНавыка, _ := ВыполнитьPgSQL(sqlStruct{
		Name:   "io_actions",
		Sql:    "UPDATE iobot.io_actions SET description = $1  WHERE id = $2 RETURNING *",
		Values: [][]byte{
			[]byte(mes.Текст),
			[]byte(ActionId),
		},
	})

	log.Printf("ОписаниеНавыка %+v\n", ОписаниеНавыка)

	if len(ОписаниеНавыка)<1{
		ОписаниеНавыка[0]= map[string]interface{}{"error":"Не удалось добавить описание навыку :"+mes.Текст}
	} else {
		ОписаниеНавыка[0]["result"]=map[string]interface{}{
			"io_action":ОписаниеНавыка[0]["id"].(string),
		}
	}

	return ОписаниеНавыка[0]
}



func (client *Client)CreateSkillName(mes Сообщение) (map[string]interface{}) {

	НовыйНавык, _ := ВыполнитьPgSQL(sqlStruct{
				Name:   "io_actions",
				Sql:    "INSERT INTO iobot.io_actions (name, creator) VALUES ($1,$2) ON CONFLICT (name) DO NOTHING RETURNING *",
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
			"io_action":НовыйНавык[0]["id"].(string),
		}
	}

	return НовыйНавык[0]
}

func (client *Client) ЗапросБоту(mes Сообщение){

	/*
		 Если боту пришло сообщение
		1. Проверить есть ли с этим пользователем не оконченный диалог
			если нет то:
				1. Проверим есть ли в диолгах команда пришедшая в тексте
					а. если есть то получим данные ответа
				2. если в диалогах нет похожей команды то перенаправим вопрос куратору который в сети
						-если ни одного куратора в сети нет, или куратор не ответил в течении нескольких секунд то ИО ответит что не понял вопроса
			если есть, то получим
	*/
	client.ПолучитьАктивныйДиалог(mes)

	//ДиалогиВОжиданииОтвета, _ := ВыполнитьPgSQL(sqlStruct{
	//	Name:   "dialogs_log",
	//	Sql:    "SELECT * FROM iobot.io_dialogs_log WHERE uid = $1 AND status->>'завершено' IS NULL",
	//	Values: [][]byte{
	//		[]byte(client.Login),
	//	},
	//})

	//log.Printf("ДиалогиВОжиданииОтвета %+v\n", ДиалогиВОжиданииОтвета)


	if len(АктивныеДиалоги) >0{
		//
		var АктивныйДиалог map[string]interface{}
		log.Printf("len(АктивныеДиалоги) %+v\n", len(АктивныеДиалоги))
		if len(АктивныеДиалоги) == 1 {
			АктивныйДиалог = АктивныеДиалоги[0]
			// получим id диалога и сообщения, затем получи action io_dialogs = id диалога+сообщения
			Действие, _ := ВыполнитьPgSQL(sqlStruct{
				Name:   "io_dialogs",
				Sql:    "SELECT * FROM iobot.io_dialogs WHERE dialog_id = $1 AND message_id = $2 ",
				Values: [][]byte{
					[]byte(АктивныйДиалог["dialog_id"].(string)),
					[]byte(АктивныйДиалог["message_id"].(string)),
				},
			})

			СтатусАктивногоДиалога := map[string]interface{}{}
			err := json.Unmarshal([]byte(АктивныйДиалог["status"].(string)), &СтатусАктивногоДиалога)
			if err != nil {
				log.Printf("err	 %+v\n", err)
			}

			if СтатусАктивногоДиалога["await"] != nil && СтатусАктивногоДиалога["await"] == "user_message" {
				// проверим есть ли ответ пользователя в вариантах ответа которые возможны для последнего активного впороса
				СледующееСообщение,_:=ВыполнитьPgSQL(sqlStruct{
					Name:   "io_dialogs",
					Sql:    "SELECT * FROM iobot.io_dialogs left JOIN (SELECT jsonb_array_elements(next)::integer nextId FROM iobot.io_dialogs WHERE dialog_id = $2 AND  message_id = $3) as nextSteps ON iobot.io_dialogs.id IN (nextSteps.nextId) WHERE  user_message ? lower($1)",
					Values: [][]byte{
						[]byte(strings.TrimSpace(mes.Текст)),
						[]byte(АктивныйДиалог["dialog_id"].(string)),
						[]byte(АктивныйДиалог["message_id"].(string)),
					},
				})
				// Если вариант найдет то нужно проверить Action активного диалога и выполнить если надо.
				// Если Actions нету то никаких действийне требуеться, пометим активный вопрос как завершённый и отправим новый вопрос
				if Actions[Действие[0]["action"].(string)] == nil{

					status := map[string]interface{}{
						"завершено":"ok (Action nil)",
						"result":СтатусАктивногоДиалога["result"],
					}
					statusByte,err := json.Marshal(status)
					if err != nil {
						log.Printf("err	 %+v\n", err)
					}
					_,_ = ВыполнитьPgSQL(sqlStruct{
						Name:   "io_dialogs",
						Sql:    "UPDATE iobot.io_dialogs_log SET status = status || $2  WHERE id = $1",
						Values: [][]byte{
							[]byte(АктивныйДиалог["id"].(string)),
							statusByte,
						},
					})

					СообщениеИО := map[string]interface{}{}
					err = json.Unmarshal([]byte(СледующееСообщение[0]["io"].(string)), &СообщениеИО)
					if err != nil {
						log.Printf("err	 %+v\n", err)
					}

					// Добавим АктивныйВопрос

					status = map[string]interface{}{
							"await":СообщениеИО["await"].(string),
							"result":СтатусАктивногоДиалога["result"],
						}

					statusByte,err = json.Marshal(status)
					if err != nil {
						log.Printf("err	 %+v\n", err)
					}
					_,_=ВыполнитьPgSQL(sqlStruct{
						Name:   "io_dialogs_log",
						Sql:    "INSERT INTO iobot.io_dialogs_log (uid, message, status, time_log, dialog_id, message_id) VALUES ($1,$2,$3,NOW(), $4, $5) RETURNING *",
						Values: [][]byte{
							[]byte(client.Login),
							[]byte(СообщениеИО["question"].(string)),
							statusByte,
							[]byte(СледующееСообщение[0]["dialog_id"].(string)),
							[]byte(СледующееСообщение[0]["message_id"].(string)),
						},
					})

					// Отправим пользователю сообщение
					if СообщениеИО["question"] != nil{
						ТекстСообщения := СообщениеИО["question"].(string)
						if СообщениеИО["variants"] !=nil{
							ТекстСообщения = ТекстСообщения + string(render("variants", СообщениеИО))
						}
						mes:= &Сообщение{
							Текст:  ТекстСообщения,
							От: "io",
							Кому:client.Login,
							Content:struct {
								Target string `json:"target"`
								Data interface{} `json:"data"`
								Html string `json:"html"`
								Обработчик string `json:"обработчик"`
							}{
								//Data: map[string]string{"SkillId":НовыйНавык["id"].(string)},
								//Обработчик: ОтветИО[0]["action"].(string),
							},
						}
						idMes := СохранитьСообщение(*mes)
						mes.Id=idMes
						client.Message<-mes
					}
				}
			}

			if  Actions[Действие[0]["action"].(string)] != nil{

				РезультатДейсвтия:=Actions[Действие[0]["action"].(string)].(func(*Client, Сообщение)map[string]interface{})(client, mes)
				// Сделать проверку, что действие выполнилось успешно


				//Если не удалось выполнить действие вернём ошибку
				if РезультатДейсвтия["error"] != nil{

					mes:= &Сообщение{
						Текст:  РезультатДейсвтия["error"].(string),
						От: "io",
						Кому:client.Login,
					}
					idMes := СохранитьСообщение(*mes)
					mes.Id=idMes
					client.Message<-mes

					return
				} else {
					var status map[string]interface{}

					if РезультатДейсвтия["result"] !=nil{
						status = map[string]interface{}{
							"завершено":"ok",
							"result":РезультатДейсвтия["result"],
						}
					}
					statusByte,err := json.Marshal(status)
					if err != nil {
						log.Printf("err	 %+v\n", err)
					}
					ОбновимСтатусСообщения,_ := ВыполнитьPgSQL(sqlStruct{
						Name:   "io_dialogs",
						Sql:    "UPDATE iobot.io_dialogs_log SET status = status || $2  WHERE id = $1",
						Values: [][]byte{
							[]byte(АктивныйДиалог["id"].(string)),
							statusByte,
						},
					})
					log.Printf("ОбновимСтатусСообщения %+v\n", ОбновимСтатусСообщения)

				}

				if Действие[0]["next"] != ""{

					log.Printf("Действие %+v\n", Действие)

					NextMessages := Действие[0]["next"]
					var NextIds []json.Number
					err := json.Unmarshal([]byte(NextMessages.(string)), &NextIds)

					log.Printf("err %+v\n", err)

					ОтветИО,_ := ВыполнитьPgSQL(sqlStruct{
						Name:   "io_dialogs",
						Sql:    "SELECT * FROM iobot.io_dialogs WHERE dialog_id = $1 AND message_id = $2",
						Values: [][]byte{
							[]byte(Действие[0]["dialog_id"].(string)),
							[]byte(NextIds[0]),
						},
					})

					if len(ОтветИО)>0 && ОтветИО[0]["io"] != nil {

						СообщениеИО := map[string]interface{}{}

						err := json.Unmarshal([]byte(ОтветИО[0]["io"].(string)), &СообщениеИО)
						if err != nil {
							log.Printf("err	 %+v\n", err)
						}

						var status map[string]interface{}

						if РезультатДейсвтия["result"] !=nil{
							status = map[string]interface{}{
								"await":СообщениеИО["await"].(string),
								"result":РезультатДейсвтия["result"],
							}
						}
						statusByte,err := json.Marshal(status)

						_,_=ВыполнитьPgSQL(sqlStruct{
							Name:   "io_dialogs_log",
							Sql:    "INSERT INTO iobot.io_dialogs_log (uid, message, status, time_log, dialog_id, message_id) VALUES ($1,$2,$3,NOW(), $4, $5) RETURNING *",
							Values: [][]byte{
								[]byte(client.Login),
								[]byte(СообщениеИО["question"].(string)),
								statusByte,
								[]byte(ОтветИО[0]["dialog_id"].(string)),
								[]byte(ОтветИО[0]["message_id"].(string)),
							},
						}) //[]byte(`{"await":"`+СообщениеИО["await"].(string)+`"}`),

						// ОБРАБОТКА ОТВЕТА БОТА
						// НЕОБХОДИМО ДОБАВИТЬ ОБРАБОТКУ КНОПОК



						if СообщениеИО["question"] != nil{
							ТекстСообщения := СообщениеИО["question"].(string)
							if СообщениеИО["variants"] !=nil{
								ТекстСообщения = ТекстСообщения + string(render("variants", СообщениеИО))
							}
							mes:= &Сообщение{
								Текст:  ТекстСообщения,
								От: "io",
								Кому:client.Login,
								Content:struct {
									Target string `json:"target"`
									Data interface{} `json:"data"`
									Html string `json:"html"`
									Обработчик string `json:"обработчик"`
								}{
									//Data: map[string]string{"SkillId":НовыйНавык["id"].(string)},
									//Обработчик: ОтветИО[0]["action"].(string),
								},
						}
							idMes := СохранитьСообщение(*mes)
							mes.Id=idMes
							client.Message<-mes
						}
					}

				}

			} else {

				//Если обработчик Actions не создан, значит он не нужен


				// Нужно будет убрать это сообщение или выводить только мне
				if client.Login == "maksimchuk@r26" {
					mes:= &Сообщение{
						Текст: "Нет обрабочтика для действия "+Действие[0]["action"].(string),
						От: "io",
						Кому:client.Login,
					}
					idMes := СохранитьСообщение(*mes)
					mes.Id=idMes
					client.Message<-mes
				}
			}
		}
	} else {
		//log.Printf("нет активных диалогов с ботом \n")
		//log.Printf("strings.TrimSpace(mes.Текст) %+v\n", strings.TrimSpace(mes.Текст))

		ОтветИО,_:=ВыполнитьPgSQL(sqlStruct{
			Name:   "io_dialogs",
			Sql:    "SELECT * FROM iobot.io_dialogs WHERE user_message ? lower($1) ",
			Values: [][]byte{
				[]byte(strings.TrimSpace(mes.Текст)),
			},
		})

		//log.Printf("ОтветИО %+v\n", ОтветИО)

		if len(ОтветИО) >0{

			СообщениеИО := map[string]string{}

			if ОтветИО[0]["io"].(string) != ""{
				err := json.Unmarshal([]byte(ОтветИО[0]["io"].(string)), &СообщениеИО)
				if err != nil {
					log.Printf("err	 %+v\n", err)
				}
			}
			status := `{"await":"`+СообщениеИО["await"]+`"}`
			if  status == ""{
				status = `{"завершено":"ok"}`
			}
			message := СообщениеИО["question"]


			_,_=ВыполнитьPgSQL(sqlStruct{
				Name:   "io_dialogs_log",
				Sql:    "INSERT INTO iobot.io_dialogs_log (uid, message, status, time_log, dialog_id, message_id) VALUES ($1,$2,$3,NOW(), $4, $5) RETURNING *",
				Values: [][]byte{
					[]byte(client.Login),
					[]byte(message),
					[]byte(status),
					[]byte(ОтветИО[0]["dialog_id"].(string)),
					[]byte(ОтветИО[0]["message_id"].(string)),
				},
			})

			mes := &Сообщение{
				Текст: СообщениеИО["question"],
				От: "io",
				Кому:client.Login,
			}
			//idMes :=
			mes.Id=СохранитьСообщение(*mes)
			client.Message<-mes


		}



	}

	// Сохраним лог диалога, добавим актинвый диалог

}

func (client *Client)AiHandler(mes Сообщение){
	log.Printf("mes.Выполнить %+v\n", mes.Выполнить)

log.Printf("Actions[mes.Выполнить.Action] %+v\n", mes.Выполнить.Action)
log.Printf("Actions %+v\n",Actions)


	if mes.Текст != ""{
		client.ЗапросБоту(mes)

	} else if mes.Выполнить.Action != "" && Actions[mes.Выполнить.Action] != nil {
		Actions[mes.Выполнить.Action].(func(*Client, Сообщение))(client, mes)
	} else {
		/*
		Пробеум загрузить плагин и выполнить действие из него
		*/
		client.Plugins (mes)
	}
	if mes.Выполнить.Cmd != ""{

	}
	if mes.Выполнить.Skill != 0{
		client.ЗапуститьНавык(mes)
	}



}

type СтруктураНавыка struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Actions map[string]interface{} `json:"actions"`
}

func (client *Client)получитьНавыкOLD(skillId int) СтруктураНавыка {
	conn := DBconn()
	defer conn.Close()
	sqlStr := `SELECT id , name , actions FROM io_actions WHERE id= $1`
	skillQuery,err := conn.Query(sqlStr, skillId)
	if err != nil{
		log.Printf("err %+v\n", err)
	}
	var id int
	var name string
	var actions []byte
	var actionMap map[string]interface{}
	for skillQuery.Next(){
		skillQuery.Scan(&id, &name, &actions)
		err := json.Unmarshal(actions, &actionMap)
		if err != nil {
			log.Printf("err  %+v Unmarshal( %+v , &actionMap) \n", err, actions)
		}
	}
	return СтруктураНавыка{
		Id:      id,
		Name:    name,
		Actions: actionMap,
	}
}

func (client *Client)получитьНавык(skillId int) СтруктураНавыка {
	sqlQuery := sqlStruct{
		Name:   `io_actions`,
		Sql:    `SELECT id , name , actions FROM io_actions WHERE id= $1`,
		Values: [][]byte{
			[]byte(strconv.Itoa(skillId)),
		},
	}
	//BotMenu := []*BotMenuStruct{}
	ActionMap, _ := ВыполнитьPgSQL(sqlQuery)
	Action := СтруктураНавыка{}

	for _, row := range ActionMap{
		for cellName, cell := range row{
			if cellName == "actions"{
				actions := map[string]interface{}{}
				json.Unmarshal([]byte(cell.(string)), &actions)
				Action.Actions=actions
			}
		}
	}
return Action
}


func (client *Client)ЗапуститьНавык(mes Сообщение){
	skillId := mes.Выполнить.Skill
	навык := client.получитьНавык(skillId)

	log.Printf("навык %+v logCmd: %+v\n", навык)

	logCmd, errors := client.RunCmd(навык.Actions["cmd"].(string))

	client.СохранитьЛог(skillId, logCmd, errors)

	if errors != nil {
		log.Printf("errors %+v\n", errors)
	}
	log.Printf("навык %+v logCmd: %+v\n", навык, logCmd)
}

func  (client *Client)СохранитьЛог(skillId int,logCmd string, errors []string){
	log.Printf("skillId %+v logCmd %+v errors %+v\n",skillId , logCmd , errors)
	conn := DBconn()
	defer conn.Close()
	errorsBytes, err := json.Marshal(errors)
	if err != nil {
		log.Printf("err	 %+v\n", err)
	}
	logBytes, err := json.Marshal(logCmd)
	if err != nil {
		log.Printf("err	 %+v\n", err)
	}
	log.Printf("skillId %+v logCmd %+v errors %+v\n",skillId , logCmd , errors)
	print(client.Login,client.Ip, skillId, errorsBytes,logBytes)
	sqlStr := `INSERT INTO io_action_log (login, ip, action_id, action_error, action_out) VALUES($1,$2,$3,$4,$5)`
	i , err := conn.Exec(sqlStr, client.Login,client.Ip, skillId, errorsBytes,logBytes)
	//log.Printf("INSERT i %+v\n", i)
	//log.Printf("INSERT err %+v\n", err)
	if err != nil{
		log.Printf("INSERT err: %+v ; i: %+v; \n", err, i)
	}
}
type BotMenuStruct struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Problem []string `json:"problem"`
	Actions map[string]interface{} `json:"actions"`
	Tags []string `json:"groups"`
}
func (client *Client)ПолучитьМенюКуратора()  []BotMenuStruct{
	if client.UserInfo.Info.OspCode== 26911{
		return client.ПолучитьМенюБота()
	}
	return nil
}
func (client *Client)ПолучитьМенюБота()  []BotMenuStruct {//
	sqlQuery := sqlStruct{
		Name:   `io_actions`,
		Sql:    `SELECT * FROM io_actions`,
	}
	//BotMenu := []*BotMenuStruct{}
	BotMenuMap, _ := ВыполнитьPgSQL(sqlQuery)
	BotMenu :=  make([]BotMenuStruct, len(BotMenuMap))

	for i, row := range BotMenuMap{
		for cellName, cell := range row{
			if cellName == "actions"{
				actions := map[string]interface{}{}
				json.Unmarshal([]byte(cell.(string)), &actions)
				BotMenu[i].Actions=actions
			}
			if cellName == "problem"{
				problem := []string{}
				json.Unmarshal([]byte(cell.(string)), &problem)
				BotMenu[i].Problem=problem
			}
		}
	}


//
//log.Printf("BotMenu %+v\n", BotMenu[0])
//log.Printf("BotMenuMap %+v\n", BotMenuMap)
return BotMenu
}

//func (client *Client)ПолучитьМенюБотаOLD() []BotMenuStruct {
//return nil
//	conn := DBconn()
//	defer conn.Close()
//	//sqlStr := `SELECT * FROM messages WHERE (autor = $1 AND recipient = $2) OR (recipient = $1 AND autor = $2) ORDER BY date ASC LIMIT 10`  // AND date >= CURRENT_DATE - INTERVAL '1 day'
//	sqlStr := `SELECT * FROM io_actions`
//
//	botMenu,err := conn.Query(sqlStr)
//	if err != nil{
//		log.Printf("err %+v\n", err)
//	}
//
//	// ДОБАВИТЬ СООБЩЕНИЕ ПРИВЕТСВИЯ
//	var id int
//	var name string
//	var description string
//	var problem []byte
//	var actions []byte
//
//	var BotMenu []BotMenuStruct
//	for botMenu.Next(){
//		err := botMenu.Scan(&id, &name, &description,&problem, &actions)
//		if err != nil{
//			log.Printf("err %+v\n", err)
//		}
//		var problemString []string
//		err = json.Unmarshal(problem, &problemString)
//		if err != nil {
//			log.Printf("err	 %+v\n", err)
//		}
//		var actionsString map[string]interface{}
//		err = json.Unmarshal(actions, &actionsString)
//		if err != nil {
//			log.Printf("err	 %+v\n", err)
//		}
//
//
//		BotMenu = append(BotMenu, BotMenuStruct{id, name, description,problemString, actionsString})
//	}
//
//	return BotMenu
//}

//
//cd /opt/ && mkdir cert && cd cert && wget ftp://10.0.29.228/cryptopro/fssp_sertif/c*.zip && unzip crl*.zip && ./ser_inst.sh && unzip cert.zip && cd cert && ./ser_inst.sh && rm -Rf /opt/cert