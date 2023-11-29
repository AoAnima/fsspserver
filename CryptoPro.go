package main

import (
	"bytes"
	"context"
	"errors"
	"golang.org/x/crypto/ssh"
	"log"

	"strconv"
	"strings"
	"time"
)

func (client *Client) СинхронихироватьДатыКриптоПро (вопрос Сообщение){
	Результат ,err:= sqlStruct{
		DBSchema:"ecp",
		Name:   "pcs_lists",
		Sql:    "SELECT * FROM pcs_lists",
		Values: [][]byte{},
	}.Выполнить(nil)
	if err != nil{
		Ошибка(">>>> Ошибка SQL запроса: %+v \n\n",err)
	}
	for _, СтарыеДанныеПк := range Результат {
		Инфо("crypro_date %+v", СтарыеДанныеПк["crypro_date"].(string))
		if СтарыеДанныеПк["crypro_date"] != nil || СтарыеДанныеПк["crypro_date"].(string) != ""{
			_ ,err = sqlStruct{
				Name:   `записатьДатуУстановкиЛицензии`,
				Sql:    `UPDATE skzi.учёт_лицензий_криптопро SET дата_лицензии = cast ($1 as DATE) WHERE mac = $2`,
				Values: [][]byte{
					[]byte(СтарыеДанныеПк["crypro_date"].(string)),
					[]byte(СтарыеДанныеПк["mac"].(string)),
				},
				DBSchema:"skzi",
			}.Выполнить(nil)
			if err != nil {
				Ошибка(" %+v ", err)
			}
		}

	}
}

func (client *Client) СинхронихироватьКлючи (вопрос Сообщение){
	Результат ,err:= sqlStruct{
			DBSchema:"ecp",
			Name:   "pcs_lists",
			Sql:    "SELECT * FROM pcs_lists",
			Values: [][]byte{},
		}.Выполнить(nil)
	if err != nil{
	 Ошибка(">>>> Ошибка SQL запроса: %+v \n\n",err)
	}
	Инфо("%+v",Результат)

	for _, СтарыеДанныеПк := range Результат {
		КлючВНовойБазе ,err := sqlStruct{
			Name:   "сравнить_ключи",
			Sql:    `SELECT mac, лицензия FROM skzi.учёт_лицензий_криптопро WHERE mac = $1`,
			Values: [][]byte{
				[]byte(СтарыеДанныеПк["mac"].(string)),
			},
			DBSchema:"skzi",
		}.Выполнить(nil)
		if err != nil{
			Ошибка(">>>> Ошибка SQL запроса: %+v \n\n",err)
		}
		Инфо("КлючВНовойБазе %+v len(КлючВНовойБазе) %+v", КлючВНовойБазе, len(КлючВНовойБазе))
		if len(КлючВНовойБазе) < 1{
			_ ,err = sqlStruct{
				Name:   `добавитьПк`,
				Sql:    `INSERT INTO skzi.учёт_лицензий_криптопро (осп, ip, mac, uid,  время_обновления, ОС,   инвентарный_номер) VALUES ($1, $2, $3, $4, NOW(), $5,$6) ON CONFLICT (mac)  DO UPDATE SET время_обновления = NOW(), последний_вход = EXCLUDED.последний_вход, ip= EXCLUDED.ip`,
				Values: [][]byte{
					[]byte(strconv.Itoa(СтарыеДанныеПк["osp_code"].(int))),
					[]byte(СтарыеДанныеПк["ip"].(string)),
					[]byte(СтарыеДанныеПк["mac"].(string)),
					[]byte(СтарыеДанныеПк["login"].(string)),
					//[]byte(СтарыеДанныеПк["last_auth"].(string)),
					[]byte(СтарыеДанныеПк["os"].(string)),
					[]byte(СтарыеДанныеПк["pc_num"].(string)),
				},
				DBSchema:"skzi",
			}.Выполнить(nil)
			if err != nil {
				Ошибка(" %+v ", err)
			}
		} else {
			if КлючВНовойБазе[0]["лицензия"] == "" && СтарыеДанныеПк["cryptopro"].(string) != "" {
				_ ,err = sqlStruct{
					Name:   `добавитьПк`,
					Sql:    `UPDATE skzi.учёт_лицензий_криптопро SET лицензия = $1 WHERE mac = $2`,
					Values: [][]byte{
						[]byte(СтарыеДанныеПк["cryptopro"].(string)),
						[]byte(СтарыеДанныеПк["mac"].(string)),
					},
					DBSchema:"skzi",
				}.Выполнить(nil)
				if err != nil {
					Ошибка(" %+v ", err)
				}
			}
			Инфо("КлючВНовойБазе %+v", КлючВНовойБазе[0]["лицензия"])
			if КлючВНовойБазе[0]["лицензия"] != "" && КлючВНовойБазе[0]["лицензия"] != СтарыеДанныеПк["cryptopro"].(string){
				_ ,err = sqlStruct{
					Name:   `добавитьПк`,
					Sql:    `UPDATE skzi.учёт_лицензий_криптопро SET старый_ключ = $1 WHERE mac = $2`,
					Values: [][]byte{
						[]byte(СтарыеДанныеПк["cryptopro"].(string)),
						[]byte(СтарыеДанныеПк["mac"].(string)),
					},
					DBSchema:"skzi",
				}.Выполнить(nil)
				if err != nil {
					Ошибка(" %+v ", err)
				}
			}


		}

	}




}

func СобратьДанныеПКизТМА() map[string]РезультатSQL {
	//бд :=FSSPconnect()
	Результат ,err := sqlStruct{
		Dbs: []interface{}{"fssp"},
		Name:   "список_компов_из_тма",
		Sql:    `SELECT  SCAN_USER uid, DNS_UNIT_HARDWARE.MAIN_MAC mac, DNS_UNIT_HARDWARE.MAIN_IP ip, OS_NAME,DNS_UNIT_HARDWARE.CREATION_TIMESTAMP last_auth, MB_NAME, CPU_NAME, COMMENT, OSP_DEP_NAME, DEPARTMENT, INVNUM_SYSTEM_UNIT from DNS_UNIT_HARDWARE join (select distinct DNS_UNIT_HARDWARE.MAIN_MAC, max(DNS_UNIT_HARDWARE.CREATION_TIMESTAMP) TIMES from DNS_UNIT_HARDWARE group by MAIN_MAC) t2 ON t2.MAIN_MAC = DNS_UNIT_HARDWARE.MAIN_MAC AND CREATION_TIMESTAMP = TIMES left join osp ON osp.DIV_NAME = OSP_DEP_NAME WHERE  IS_SERVER = 0`,
		Values: [][]byte{},
		//DBSchema:"skzi",
	}.ВыполнитьSQLвАИС()

	if err != nil {
		Ошибка(" %+v ", err)
	} else {
		//Инфо("Результат %+v", Результат)
	}

	//Таблица := map[string]interface{}{
	//	"Контейнер":"journal_skzi.result",
	//	"HTML": render("TableGenerator", Результат),
	//}
	//log.Printf("Таблица Строки %+v",  Строки)
	//Данные := &ДанныеОтвета{
	//	Контейнер:"journal_skzi.result",
	//	HTML: string(render("TableGenerator", Результат["список_компов_из_тма"])),
	//}

	//СообщениеКлиенту:= &Сообщение{
	//			Текст:   "Данные из ТМА получены, обновляю таблицу с рабочими станциями",
	//			От: "io",
	//			Кому:client.Login,
	//			MessageType: []string{"io_action"},
	//			//Контэнт: Данные,
	//
	//		}
	//СообщениеКлиенту.СохранитьИОтправить(client)
	//return  Таблица
	return Результат
}

func ДобавитьПК (ПК map[string]string){

	ОСП :="260"+ПК["DEPARTMENT"]
	ИнвентарныйНомер, ЕстьИнвентарныйНомер:=ПК["INVNUM_SYSTEM_UNIT"]
	if !ЕстьИнвентарныйНомер {
		ИнвентарныйНомер = ""
	}

	Мак, ЕстьМак := ПК["MAC"]
	if ЕстьМак {
		Мак = strings.ReplaceAll(Мак, "-", ":")
	}

	//MB_NAME, CPU_NAME, COMMENT
	//МатПлата := strings.ReplaceAll(ПК["MB_NAME"], "\n", "")
	МатПлата := strings.Join(strings.Fields(ПК["MB_NAME"]), " ")
	Процессор := strings.Join(strings.Fields(ПК["CPU_NAME"]), " ")
	//Процессор := strings.ReplaceAll(ПК["CPU_NAME"], "\n", "")

	ПКИнфо:= `{"Материнская плата":"`+МатПлата+`","Процессор":"`+Процессор+`", "Комментарий":"`+ПК["COMMENT"]+`"}`
	//Инфо("ПКИнфо %+v", ПКИнфо)
	Списан := strings.Contains(ПК["COMMENT"], "Списан")

	_ ,err:= sqlStruct{
			Name:   `добавитьПк`,
			Sql:    `INSERT INTO skzi.учёт_лицензий_криптопро (осп, ip, mac, uid, последний_вход, время_обновления, ОС,  пк, инвентарный_номер, списан) VALUES ($1, $2, $3, $4, $5, NOW(), $6,$7,$8,$9) ON CONFLICT (mac)  DO UPDATE SET время_обновления = NOW(), последний_вход = EXCLUDED.последний_вход, ip= EXCLUDED.ip`,
			Values: [][]byte{
				[]byte(ОСП),
				[]byte(ПК["IP"]),
				[]byte(Мак),
				[]byte(ПК["UID"]),
				[]byte(ПК["LAST_AUTH"]),
				[]byte(ПК["OS_NAME"]),
				[]byte(ПКИнфо),
				[]byte(ИнвентарныйНомер),
				[]byte(strconv.FormatBool(Списан)),
			},
			DBSchema:"skzi",
		}.Выполнить(nil)

	if err != nil{
		Ошибка(">>>> Ошибка SQL запроса: %+v \n\n",err)
	}


}

func (client *Client) ОбновитьДанныеПоПК (вопрос Сообщение) {

	//список_компов_из_тма := СобратьДанныеПКизТМА()
	//КомпыИзТма := список_компов_из_тма["список_компов_из_тма"].КартаСтрок
	//
	//for _, ПК  := range КомпыИзТма {
	//	//Инфо("НомерСтроки %+v ПК %+v\n", НомерСтроки, ПК)
	//	ДобавитьПК(ПК)
	//}
	СканерСети()
}
/*
	СканерСети
		Сканирует Сети ОСП
		При вохможности подключиться к ПК, подключается и выполняет указаную команду

*/
func СканерСети() {
	СписокСетей:=ПолучитьСписокСетей()
	Канал := make(chan string)
	var КоличествоПотоков int
	for i, Сеть := range СписокСетей {
		КоличествоПотоков=i
		ПросканироватьСеть(Сеть, Канал)
		for {
			СообщениеИЗПотока := <-Канал
			//log.Printf("СообщениеИЗПотока %+v КоличествоПотоков %+v\n", СообщениеИЗПотока, КоличествоПотоков)
			if СообщениеИЗПотока == "end"{
				КоличествоПотоков--
			}
			if КоличествоПотоков ==0 {
				log.Printf("КоличествоПотоков %+v\n", КоличествоПотоков)
				close(Канал)
				break
			}
		}
	}
	//for {
	//	СообщениеИЗПотока := <-Канал
	//	//log.Printf("СообщениеИЗПотока %+v КоличествоПотоков %+v\n", СообщениеИЗПотока, КоличествоПотоков)
	//	if СообщениеИЗПотока == "end"{
	//		КоличествоПотоков--
	//	}
	//	if КоличествоПотоков ==0 {
	//		log.Printf("КоличествоПотоков %+v\n", КоличествоПотоков)
	//		close(Канал)
	//		break
	//	}
	//}
}


func ПросканироватьСеть (Сеть map[string]interface{}, Канал chan<- string) {

	//var hostKey ssh.PublicKey
	//log.Printf("cpKeys %+v\n", cpKeys)


	IpСрез := strings.Split(Сеть["start_ip"].(string),".")
	Подсеть := IpСрез[2]
	Инфо("Подсеть %+v\n", Подсеть)
	НачальныйАдрес,err := strconv.Atoi(IpСрез[3])
	if err != nil {
		Ошибка(">>>> ERROR \n %+v \n\n", err)
	}
	КаналПК := make(chan string)
	КоличествоПотоковПК :=0
	for  ; НачальныйАдрес<254; НачальныйАдрес++ {
		//log.Printf("НачальныйАдрес %+v Подсеть %+v\n", НачальныйАдрес, Подсеть)
		ip := "10.26." + Подсеть + "." + strconv.Itoa(НачальныйАдрес)
		//if Подсеть == "6"{
		//	Канал <- "end"
		//}
		КоличествоПотоковПК++
		Инфо("КоличествоПотоковПК %+v Сеть %+v\n", КоличествоПотоковПК, Сеть)
		СоединитьсяСПК (ip, Сеть ,КаналПК)

		for {
			СообщениеОтПк := <- КаналПК
			if СообщениеОтПк != ""{
				КоличествоПотоковПК--
			}
			if КоличествоПотоковПК ==0 {
				Инфо("КоличествоПотоковПК %+v Сеть %+v\n", КоличествоПотоковПК, Сеть)
				//close(КаналПК)
				Канал <-"end"
				break
			}
		}
	}

		//for {
		//	СообщениеОтПк := <- КаналПК
		//	if СообщениеОтПк != ""{
		//		КоличествоПотоковПК--
		//	}
		//	if КоличествоПотоковПК ==0 {
		//		Инфо("КоличествоПотоковПК %+v Сеть %+v\n", КоличествоПотоковПК, Сеть)
		//		//close(КаналПК)
		//		Канал <-"end"
		//		break
		//	}
		//}

}

func ПолучитьТаблицуМакАдресов (client *ssh.Client){
	session, err := client.NewSession()

	if err != nil {
		Ошибка("err %+v %+v\n", err )
	}
	var Data bytes.Buffer
	var ErrData bytes.Buffer
	session.Stdout = &Data
	session.Stderr = &ErrData

	if err := session.Run(`ip  neigh`); err != nil {
		Ошибка("Failed to run: " + err.Error())

	}
	if ErrData.Len() >0 {
		Ошибка("ErrData %+v ", ErrData.String())

	}
	Инфо("Мак адреса сетевых плат на пк %+v", Data.String())

	//macs := strings.Split(Data.String(), "\n")
	session.Close()

}


func ПолучитьМакАдрес (client *ssh.Client) (string, error){
	session, err := client.NewSession()

	if err != nil {
		Ошибка("err %+v %+v\n", err )
	}
	var Data bytes.Buffer
	var ErrData bytes.Buffer
	session.Stdout = &Data
	session.Stderr = &ErrData
	defer session.Close()
	if err := session.Run("cat /sys/class/net/*/address"); err != nil {
		Ошибка("Failed to run: " + err.Error())
		return "", err
	}
	if ErrData.Len() >0 {
		Ошибка("ErrData %+v ", ErrData.String())
		return "", errors.New(ErrData.String())
	}
	Инфо("Мак адреса сетевых плат на пк %+v", Data.String())

	macs := strings.Split(Data.String(), "\n")

	return macs[0], nil
}


func ПолучитьЛицензиюПОМак(client  *ssh.Client) map[string]interface{} {
	МакАдрес, ОшибкаМакАдрес := ПолучитьМакАдрес(client)
	if ОшибкаМакАдрес != nil {
		Ошибка("ОшибкаМакАдрес %+v ", ОшибкаМакАдрес)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Результат ,err:= sqlStruct{
		Name:   "проверить_ключ_в_бд",
		Sql:    `SELECT ключ as лицензия, мак FROM установленно WHERE мак = $1`,
		Values: [][]byte{
			[]byte(МакАдрес),
		},
		DBSchema:"skzi",
	}.Выполнить(ctx)
	if err != nil && !strings.Contains(err.Error(), "Запрос не затронул ни одной строки"){
		Ошибка("%+v МакАдрес %+v \n\n",err, МакАдрес)
	}
	if len(Результат)>0{
		return Результат[0]
	}
return nil
}

func ПолучитьМакПоЛицензии(ключ string ) map[string]interface{} {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Результат ,err:= sqlStruct{
		Name:   "проверить_ключ_в_бд",
		Sql:    `SELECT ключ as лицензия, мак FROM установленно WHERE ключ = $1`,
		Values: [][]byte{
			[]byte(ключ),
		},
		DBSchema:"skzi",
	}.Выполнить(ctx)
	if err != nil && !strings.Contains(err.Error(), "Запрос не затронул ни одной строки"){
		Ошибка("%+v ключ %+v \n\n",err, ключ)
	}
	if len(Результат)>0{
		return Результат[0]
	}
	return nil
}


func  СоединитьсяСПК(ip string, Сеть map[string]interface{},  Канал chan<- string){
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("123QWer"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 5 * time.Second,
	}
	Инфо("Соединяемся с ip %+v", ip)

	client, err := ssh.Dial("tcp", ip+":22", config)

	if err != nil || client == nil {
		Инфо("нет соединения с ip %+v", ip)
		Канал <- "нет соединения: "+err.Error()+" ip: "+ip
		return
	} else {
		defer client.Close()
		//ПолучитьТаблицуМакАдресов(client)
		session, err := client.NewSession()

		if err != nil {
			Ошибка("err %+v ip %+v\n", err, ip)
			Инфо("ОТключаемся от ip %+v", ip)
			Канал <- "нет соединения: "+err.Error()+" ip: "+ip
			return
		}
		var Data bytes.Buffer
		session.Stdout = &Data
		session.Stderr = &Data





		ЛицензияИМакИзБД := ПолучитьЛицензиюПОМак(client)

		errcp := session.Run("/opt/cprocsp/sbin/amd64/cpconfig -license -view");

		Инфо("ЛицензияИМакИзБД %+v;\n errcp %+v;\n Data.String() %+v", ЛицензияИМакИзБД, errcp, Data.String())

		if  ЛицензияИМакИзБД == nil || errcp != nil {
				Ошибка("errcp  %+v", errcp)

			if (strings.Contains(Data.String(), "License is expired or not yet valid")) || strings.Contains(Data.String(),"Error code:3"){

				Ошибка("ИСТЁКШАЯ ЛИЦЕНЗИЯ  %+v ip %+v\n", Data.String(), ip)

				МакАдрес, ОшибкаМакАдрес := ПолучитьМакАдрес(client)
				if ОшибкаМакАдрес != nil {
					Ошибка("ОшибкаМакАдрес %+v ", ОшибкаМакАдрес)
				} else {
					Инфо("МакАдрес %+v ip %+v",МакАдрес, ip)

				// алгоритм . Ключ истёк необходимо переустановить ключ
				НовыйКлюч := ВзятьСвободныйКлюч(МакАдрес, Сеть["osp_code"])

				//log.Printf("НовыйКлюч %+v\n", НовыйКлюч["лицензия"])
				if Лицензия, ЕстьЛицензия := НовыйКлюч["ключ"];ЕстьЛицензия{
					sessionЛицензия, err := client.NewSession()
					if err != nil {
						Ошибка("err %+v %+v\n", err )
					}

					var ЛицензияData bytes.Buffer
					sessionЛицензия.Stdout = &ЛицензияData
					sessionЛицензия.Stderr = &ЛицензияData

					if err := sessionЛицензия.Run("/opt/cprocsp/sbin/amd64/cpconfig -license -set "+Лицензия.(string)); err != nil {
						//log.Fatal("Failed to run: " + err.Error())
						Ошибка("err /cpconfig -license -set %+v %+v macs[0] %+v\n", err , Лицензия.(string),МакАдрес)
						НеУстановился (Лицензия.(string))
					} else {
						Инфо("Установлена %+v ip %+v macs[0] %+v\n",Лицензия.(string), ip,МакАдрес )
						//12:51:40 OSPLoop.go:167: Установлена 4040R-L0000-01YYR-FGM1U-5QPZ3 ip 10.26.109.171 macs[0] 98:fa:9b:36:bd:f6  98:fa:9b:58:ad:62 98:fa:9b:36:a1:aa

					}
				}
			}
			} else {

				Инфо("Data.String() %+v", Data.String())

				КриптоПроКлюч := strings.Split(Data.String(), "\n")

				МакИКлючИзБазы := ПолучитьМакПоЛицензии(КриптоПроКлюч[1])
				МакАдрес, ОшибкаМакАдрес := ПолучитьМакАдрес(client)
				if ОшибкаМакАдрес != nil {
					Ошибка("ОшибкаМакАдрес %+v ", ОшибкаМакАдрес)
				} else {
					Инфо("МакАдрес %+v ip %+v",МакАдрес, ip)
				}

				if МакИКлючИзБазы == nil || МакИКлючИзБазы["мак"] != МакАдрес {
					sessionЛицензия, err := client.NewSession()
					if err != nil {
						Ошибка("err %+v %+v\n", err)

					}

					НовыйКлюч := ВзятьСвободныйКлюч(МакАдрес, Сеть["osp_code"])
					Инфо("ВзятьСвободныйКлюч %+v", НовыйКлюч)
					Лицензия:= НовыйКлюч["ключ"]
					var ЛицензияData bytes.Buffer
					sessionЛицензия.Stdout = &ЛицензияData
					sessionЛицензия.Stderr = &ЛицензияData

					if err := sessionЛицензия.Run("/opt/cprocsp/sbin/amd64/cpconfig -license -set " + Лицензия.(string)); err != nil {
						//log.Fatal("Failed to run: " + err.Error())
						Ошибка("err /cpconfig -license -set %+v %+v macs[0] %+v\n", err, Лицензия.(string), ЛицензияИМакИзБД["мак"].(string))
						НеУстановился(Лицензия.(string))
					} else {
						Инфо("Установлена %+v ip %+v macs[0] %+v\n", Лицензия.(string), ip, ЛицензияИМакИзБД["мак"].(string))
						//12:51:40 OSPLoop.go:167: Установлена 4040R-L0000-01YYR-FGM1U-5QPZ3 ip 10.26.109.171 macs[0] 98:fa:9b:36:bd:f6  98:fa:9b:58:ad:62 98:fa:9b:36:a1:aa

					}
				}

				Инфо("ОТключаемся от ip %+v", ip)
				Канал <- "end"
				return

			}
			Инфо("ОТключаемся от ip %+v ЛицензияИМакИзБД %+v", ip, ЛицензияИМакИзБД)

			session.Close()
		} else {
			session.Close()
			//log.Printf("Data.String() %+v\n", Data.String())
			//log.Printf("Data.String() %+v\n", Data)
			КриптоПроКлюч := strings.Split(Data.String(), "\n")
			//log.Printf("license %+v len(license) %+v ip %+v Data.String(), %+v\n", license, len(license), ip, Data.String(),)
				Инфо("КриптоПроКлюч %+v", КриптоПроКлюч[1])
				Инфо("ЛицензияИМакИзБД %+v", ЛицензияИМакИзБД["лицензия"])
				if ЛицензияИМакИзБД["лицензия"] != КриптоПроКлюч[1] {
								sessionЛицензия, err := client.NewSession()
								if err != nil {
									Ошибка("err %+v %+v\n", err)

								}
								НовыйКлюч := ВзятьСвободныйКлюч(ЛицензияИМакИзБД["мак"].(string), Сеть["osp_code"])
								Лицензия:= НовыйКлюч["ключ"]
								var ЛицензияData bytes.Buffer
								sessionЛицензия.Stdout = &ЛицензияData
								sessionЛицензия.Stderr = &ЛицензияData

								if err := sessionЛицензия.Run("/opt/cprocsp/sbin/amd64/cpconfig -license -set " + Лицензия.(string)); err != nil {
									//log.Fatal("Failed to run: " + err.Error())
									Ошибка("err /cpconfig -license -set %+v %+v macs[0] %+v\n", err, Лицензия.(string), ЛицензияИМакИзБД["мак"].(string))
									НеУстановился(Лицензия.(string))
								} else {
									Инфо("Установлена %+v ip %+v macs[0] %+v\n", Лицензия.(string), ip, ЛицензияИМакИзБД["мак"].(string))
									//12:51:40 OSPLoop.go:167: Установлена 4040R-L0000-01YYR-FGM1U-5QPZ3 ip 10.26.109.171 macs[0] 98:fa:9b:36:bd:f6  98:fa:9b:58:ad:62 98:fa:9b:36:a1:aa

								}
				}



			//if len(КриптоПроКлюч)>1 {
			//	//license[1] = strings.TrimSpace(license[1])
			//	ВидЛицензии := strings.TrimSpace(КриптоПроКлюч[2])
			//	if ВидЛицензии == "License is expired or not yet valid"{
			//		Ошибка("ИСТЁКШАЯ ЛИЦЕНЗИЯ  %+v ip %+v \n", КриптоПроКлюч, ip)
			//		МакАдрес, ОшибкаМакАдрес := ПолучитьМакАдрес(client)
			//
			//
			//
			//		if ОшибкаМакАдрес != nil {
			//			Ошибка("ОшибкаМакАдрес %+v ", ОшибкаМакАдрес)
			//		} else {
			//
			//		НовыйКлюч := ВзятьСвободныйКлюч(МакАдрес, Сеть["osp_code"])
			//		Инфо("МакАдрес  %+v НовыйКлюч %+v\n",МакАдрес, НовыйКлюч["лицензия"])
			//		if Лицензия, ЕстьЛицензия := НовыйКлюч["лицензия"];ЕстьЛицензия {
			//			sessionЛицензия, err := client.NewSession()
			//			if err != nil {
			//				Ошибка("err %+v %+v\n", err)
			//
			//			}
			//			var ЛицензияData bytes.Buffer
			//			sessionЛицензия.Stdout = &ЛицензияData
			//			sessionЛицензия.Stderr = &ЛицензияData
			//
			//			if err := sessionЛицензия.Run("/opt/cprocsp/sbin/amd64/cpconfig -license -set " + Лицензия.(string)); err != nil {
			//				//log.Fatal("Failed to run: " + err.Error())
			//				Ошибка("err /cpconfig -license -set %+v %+v macs[0] %+v\n", err, Лицензия.(string), МакАдрес)
			//				НеУстановился(Лицензия.(string))
			//			} else {
			//				Инфо("Установлена %+v ip %+v macs[0] %+v\n", Лицензия.(string), ip, МакАдрес)
			//				//12:51:40 OSPLoop.go:167: Установлена 4040R-L0000-01YYR-FGM1U-5QPZ3 ip 10.26.109.171 macs[0] 98:fa:9b:36:bd:f6  98:fa:9b:58:ad:62 98:fa:9b:36:a1:aa
			//
			//			}
			//		}}
			//	}  else {
			//		Инфо("ВидЛицензии %+v", ВидЛицензии)
			//
			//		МакАдрес, ОшибкаМакАдрес := ПолучитьМакАдрес(client)
			//		if ОшибкаМакАдрес != nil {
			//			Ошибка("ОшибкаМакАдрес %+v ", ОшибкаМакАдрес)
			//		} else {
			//			Инфо("МакАдрес %+v ip %+v",МакАдрес, ip)
			//		}
			//		Инфо("МакАдрес %+v license %+v ip %+v\n", МакАдрес , КриптоПроКлюч, ip)
			//
			//		//ПроверитьИОбновитьЗаписьКриптоПроКлючаВБД (МакАдрес,  КриптоПроКлюч)
			//	}
			//}
			Инфо("ОТключаемся от ip %+v %+v %+v", ip, КриптоПроКлюч[0], ЛицензияИМакИзБД)
			Канал <- "end"
			return
		}
	}
}



func ПроверитьИОбновитьЗаписьКриптоПроКлючаВБД (МакАдрес string, КриптоПроКлюч []string){
	Инфо("КриптоПроКлюч 0 %+v МакАдрес %+v КриптоПроКлюч 1 %+v КриптоПроКлюч 2 %+v", КриптоПроКлюч[0], МакАдрес , КриптоПроКлюч[1], КриптоПроКлюч[2])


	Ключ := strings.TrimSpace(КриптоПроКлюч[1])

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Результат ,err:= sqlStruct{
			Name:   "проверить_ключ_в_бд",
			Sql:    `SELECT лицензия, пк FROM учёт_лицензий_криптопро WHERE mac = $1`,
			Values: [][]byte{
				[]byte(МакАдрес),
			},
			DBSchema:"skzi",
		}.Выполнить(ctx)
	if err != nil && !strings.Contains(err.Error(), "Запрос не затронул ни одной строки"){
		Ошибка("%+v МакАдрес %+v \n\n",err, МакАдрес)
	} else {

	Инфо("лицензия в бд %+v len(Результат) %+v на компе %+v", Результат , len(Результат),  Ключ)

	 if len(Результат)<1 || Результат[0]["лицензия"] != Ключ {

		 ВидЛицензии := strings.TrimSpace(КриптоПроКлюч[2])
		 ТипЛицензии := strings.TrimSpace(КриптоПроКлюч[3])

		 if strings.Contains(ВидЛицензии, "permanent"){
			 ВидЛицензии = "Постоянная"
		 }

		 if strings.Contains(ТипЛицензии, "Client"){
			 ТипЛицензии = "Клиентская"
		 }
		 _ ,err := sqlStruct{
			 Name:   "обновить_ключ_в_бд",
			 Sql:    `UPDATE учёт_лицензий_криптопро SET лицензия = $1, тип_лицензии = $3, время_обновления = NOW() WHERE mac = $2`,
			 Values: [][]byte{
				 []byte(Ключ),
				 []byte(МакАдрес),
				 []byte(ВидЛицензии+", "+ТипЛицензии),
			 },
			 DBSchema:"skzi",
		 }.Выполнить(ctx)
		 if err != nil{
			 Ошибка(">>>> Ошибка SQL запроса: %+v \n\n",err)
		 }
	 } else {
	 	Инфо("Результат %+v КриптоПроКлюч %+v ", Результат, КриптоПроКлюч )
	 }

	 if len(Результат) > 0 && Результат[0]["пк"] !=nil &&  strings.Contains(Результат[0]["пк"].(map[string]interface{})["Материнская плата"].(string), "LENOVO , 3135") {
		 _ ,err := sqlStruct{
			 Name:   "обновить_ключ_в_бд",
			 Sql:    `UPDATE учёт_лицензий_криптопро SET дата_лицензии = $1 WHERE mac = $2`,
			 Values: [][]byte{
				 []byte("06-03-2020"),
				 []byte(МакАдрес),
			 },
			 DBSchema:"skzi",
		 }.Выполнить(ctx)
		 if err != nil{
			 Ошибка(">>>> Ошибка SQL запроса: %+v \n\n",err)
		 }
	 }


	}

}

func ВзятьСвободныйКлюч (Мак string, ОспКод interface{}) map[string]interface{}{
	//UPDATE skzi.лицензии_криптопро SET mac_установки = $1, осп = $2 WHERE лицензия = (SELECT лицензия FROM key) RETURNING *
	Инфо("Мак %+v ОспКод %+v\n", Мак, ОспКод)
	Результат ,err:= sqlStruct{
		Name:   "лицензии_криптопро",
		Sql:    `WITH key AS (
						SELECT * FROM skzi.лицензии_криптопро
										  left join  skzi.установленно on ключ = лицензия
						WHERE  mac_установки IS NULL  AND мак IS NULL LIMIT 1
					)
					INSERT INTO skzi.установленно (мак, ключ) VALUES ($1, (SELECT лицензия FROM key)) RETURNING *`,
		Values: [][]byte{
			[]byte(Мак),
			//[]byte(strconv.Itoa(ОспКод.(int))),
		},
		DBSchema:"skzi",
	}.Выполнить(nil)
	if err != nil{
		//log.Printf("Мак %+v ОспКод %+v\n", Мак, ОспКод)
		Ошибка(">>>> Ошибка SQL запроса: %+v Мак %+v ОспКод %+v \n\n",err, Мак, ОспКод)
		return nil
	} else {
		if len(Результат)>0{
			return Результат[0]
		}
		return nil
	}
}