package main

import (
	"bytes"
	"golang.org/x/crypto/ssh"

	"strconv"
	"strings"
	"time"
)

/*
Выбираем данные подсетей в которых есть рабочии станции сотрудников ( в базе данных есть поле start_ip в котором записан начальный ип адрес с которого начинают раздаваться адреса рабчим станциям отдела)
*/
//func ПолучитьСписокСетей() []map[string]interface{}{
//	Результат ,err:= sqlStruct{
//			Name:   "osp_address",
//			Sql:    `SELECT * FROM fssp_configs.osp_address WHERE start_ip IS NOT NUL`,
//			Values: [][]byte{},
//		}.Выполнить(nil)
//	if err != nil{
//	log.Printf(">>>> Ошибка SQL запроса: %+v \n\n",err)
//	}
//	return Результат
//}



func НеУстановился (Ключ string) {
	Результат ,err:= sqlStruct{
		Name:   "лицензии_криптопро",
		Sql:    `UPDATE skzi.лицензии_криптопро SET установлен = false WHERE лицензия = $1'`,
		Values: [][]byte{
			[]byte(Ключ),
			//[]byte(strconv.Itoa(ОспКод.(int))),
		},
		DBSchema:"skzi",
	}.Выполнить(nil)
	if err != nil {
		Ошибка(">>>> ERROR \n %+v \n\n", err)
	}
	Инфо("Результат %+v\n", Результат)
}


func СоединитьсяСПКиУстановитьКлючКриптоПро(ip string, Сеть map[string]interface{},  Канал chan<- string){
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("123QWer"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 5 * time.Second,
	}
	//log.Printf("client %+v\n", config)
	client, err := ssh.Dial("tcp", ip+":22", config)

	if err != nil || client == nil {
		//log.Printf("client не удача  %+v ip %+v\n",client, ip )
		Канал <- "end"
	} else {
		session, err := client.NewSession()

		if err != nil {
			Ошибка("err %+v %+v\n", err )

		}
		var Data bytes.Buffer
		session.Stdout = &Data
		session.Stderr = &Data
		//log.Printf("Запрашиваем данные криптопро у клиент  %+v\n",ip )
		if err := session.Run("/opt/cprocsp/sbin/amd64/cpconfig -license -view"); err != nil {
			//log.Printf("err %+v\n", err.Error()=="Process exited with status 3")

			//if err.Error()=="Process exited with status 3"{
			//	log.Printf("Data.String() %+v\n", Data.String())
			//	errClose := session.Close()
			//	if err != nil {
			//		log.Printf(">>>> ERROR \n %+v \n\n", errClose)
			//	}
			//	Newsession, Newerr := client.NewSession()
			//	if Newerr != nil {
			//		log.Printf(">>>> ERROR \n %+v \n\n", err)
			//	}
			//	if err := Newsession.Run("dmidecode");err != nil{
			//		log.Printf("err %+v\n", )
			//	}
			//	log.Printf("Newsession  Stdout%+v\n", Newsession.Stdout)
			//	log.Printf("Newsession Stderr %+v\n", Newsession.Stderr)
			//	//if strings.Contains(Data.String(), "ThinkCentre M920x") {
			//	//	log.Printf("LENOVO ThinkCentre M920x %+v\n", ip )
			//	//}
			//}
			//Manufacturer: LENOVO
			//	Product Name: 10S0S0G600
			//Version: ThinkCentre M920x
			//	Serial Number: PC17TMHL
			//UUID: 10589a80-c828-11e9-91d2-ea4bb06e1f00
			//	Wake-up Type: Power Switch
			//	SKU Number: LENOVO_MT_10S0_BU_Think_FM_ThinkCentre M920x
			//Family: ThinkCentre M920x

			//log.Printf("ip %+v\n", ip)
			//log.Printf("не удаёться Получить Установленный Ключ: " + err.Error())
			//log.Printf("Data.String(): %+v" , Data.String())
			//log.Printf("License is expired or not yet valid: %+v ip %+v" , strings.Contains(Data.String(), "License is expired or not yet valid") , ip)

			//log.Printf("Не валлидная %+v\n", strings.Contains(Data.String(), "License is expired or not yet valid"))
			if strings.Contains(Data.String(), "License is expired or not yet valid")	{
				errClose := session.Close()
				if errClose != nil {

					Ошибка(">>>> ERROR \n %+v \n\n", errClose)
				}
				sessionMacData, err := client.NewSession()

				if err != nil {
					Ошибка("err %+v %+v\n", err )

				}
				var MacData bytes.Buffer
				sessionMacData.Stdout = &MacData
				sessionMacData.Stderr = &MacData

				if err := sessionMacData.Run("cat /sys/class/net/*/address"); err != nil {
					Инфо("Failed to run: " + err.Error())

				}
				macs := strings.Split(MacData.String(), "\n")

				Инфо(" macs %+v ip %+v \n",macs , ip)
				if macs[0]!= ""{
					НовыйКлюч := ВзятьСвободныйКлюч(macs[0], Сеть["osp_code"])
					//log.Printf("НовыйКлюч %+v\n", НовыйКлюч["лицензия"])
					if Лицензия, ЕстьЛицензия := НовыйКлюч["лицензия"];ЕстьЛицензия{
						sessionЛицензия, err := client.NewSession()
						if err != nil {
							Ошибка("err %+v %+v\n", err )

						}
						var ЛицензияData bytes.Buffer
						sessionЛицензия.Stdout = &ЛицензияData
						sessionЛицензия.Stderr = &ЛицензияData

						if err := sessionЛицензия.Run("/opt/cprocsp/sbin/amd64/cpconfig -license -set "+Лицензия.(string)); err != nil {
							//log.Fatal("Failed to run: " + err.Error())
							Ошибка("err /cpconfig -license -set %+v %+v macs[0] %+v\n", err , Лицензия.(string), macs[0])
							НеУстановился (Лицензия.(string))
						} else {
							Инфо("Установлена %+v ip %+v macs[0] %+v\n",Лицензия.(string), ip, macs[0] )
							//12:51:40 OSPLoop.go:167: Установлена 4040R-L0000-01YYR-FGM1U-5QPZ3 ip 10.26.109.171 macs[0] 98:fa:9b:36:bd:f6  98:fa:9b:58:ad:62 98:fa:9b:36:a1:aa

						}

					}
				}

			}


			Канал <- "end"
		} else {

			//log.Printf("Data.String() %+v\n", Data.String())
			license := strings.Split(Data.String(), "\n")
			//log.Printf("license %+v len(license) %+v ip %+v Data.String(), %+v\n", license, len(license), ip, Data.String(),)

			if len(license)>1 {
				//license[1] = strings.TrimSpace(license[1])
				ВидЛицензии := strings.TrimSpace(license[2])

				if ВидЛицензии == "License is expired or not yet valid"{
					Инфо("ИСТЁКШАЯ ЛИЦЕНЗИЯ  %+v ip %+v\n", license, ip)
				}
				//log.Printf("Data %+v ВидЛицензии %+v\n", Data.String() , ВидЛицензии)
				session.Close()
			}
			Канал <- "end"
		}


	}
}

func ПоСети (Сеть map[string]interface{}, Канал chan<- string) {

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

		ip := "10.26." + Подсеть + "." + strconv.Itoa(НачальныйАдрес)
		//if Подсеть == "6"{
		//	Канал <- "end"
		//}
		КоличествоПотоковПК++
		go СоединитьсяСПКиУстановитьКлючКриптоПро (ip, Сеть ,КаналПК)
	}
	for {
		СообщениеОтПк := <- КаналПК
		if СообщениеОтПк == "end"{
			КоличествоПотоковПК--
		}
		if КоличествоПотоковПК ==0 {
			Инфо("КоличествоПотоковПК %+v Сеть %+v\n", КоличествоПотоковПК, Сеть)
			close(КаналПК)
			break
		}
	}
	Канал <-"end"
}


func ПоВсемРабочимСтанциям(){
	СписокСетей := ПолучитьСписокСетей()
	Канал := make(chan string)
	var КоличествоПотоков int
	for i, Сеть := range СписокСетей {
		КоличествоПотоков=i
		go ПоСети(Сеть, Канал)
	}
	for {
		СообщениеИЗПотока := <-Канал
		Инфо("СообщениеИЗПотока %+v КоличествоПотоков %+v\n", СообщениеИЗПотока, КоличествоПотоков)
		if СообщениеИЗПотока == "end"{
			КоличествоПотоков--
		}
		if КоличествоПотоков ==0 {
			Инфо("КоличествоПотоков %+v\n", КоличествоПотоков)
			close(Канал)
			break
		}

	}

}


func ПодключитьсяКПк (){

}

//
//func ПолучитьСписокОсп(conn *sql.DB) (map[string]Key, error) {
//	log.Printf("ПолучитьСписокОсп %+v\n", )
//	cryptoKeyQuery, err := conn.Query("select osp_code, ip_ais, osp_name from osp_address WHERE osp_name is not NULL")
//	if err != nil {
//		log.Printf("err %+v\n", err)
//
//		return nil, err
//	}
//
//
//	resultData := map[string]Key{}
//
//	for cryptoKeyQuery.Next() {
//
//		var OspCodeS string
//		var ip []byte
//		var OspName string
//		err = cryptoKeyQuery.Scan(&OspCodeS, &ip, &OspName)
//		if err != nil {
//
//			return nil, err
//		} else {
//
//			subNet := strings.Split(string(ip),".")
//
//			resultData[OspCodeS]=Key{
//				ServerIp: string(ip),
//				OspCode:    OspCodeS,
//				SubNet:subNet[2],
//				OspName:OspName,
//			}
//		}
//	}
//
//	return resultData, nil
//}
//
//func ПолучитьУстановленныйКлюч(client *ssh.Client) ([]string, error){
//	var b bytes.Buffer
//
//	session, err := client.NewSession()
//	session.Stdout = &b
//	if err != nil {
//		log.Printf("err %+v %+v\n", err )
//		return nil, err
//	}
//	if err := session.Run("/opt/cprocsp/sbin/amd64/cpconfig -license -view"); err != nil {
//		log.Printf("не удаёться Получить Установленный Ключ: " + err.Error())
//		return nil, err
//	}
//	license := strings.Split(b.String(), "\n")
//	license[1] = strings.TrimSpace(license[1])
//	//log.Printf("license %+v\n", license)
//	session.Close()
//	return license, nil
//}
//func УстановитьКлюч(client *ssh.Client, cpKeys string) ([]string, error){
//	//log.Printf("Установим ключ %+v\n", cpKeys.CryptoKeys)
//
//	//if len(cpKeys.CryptoKeys)==0{
//	//	//checkLic, err := ПроверитьКлюч(client)
//	//	err := errors.New("Не хватает ключей")
//	//	return nil, err
//	//}
//
//	var setCpKey bytes.Buffer
//	session, err := client.NewSession()
//	if err != nil {
//		return nil, err
//	}
//	session.Stdout = &setCpKey
//	//var CPKey string
//	//log.Printf("Количество ключей cpKeys.CryptoKeys %+v\n", len(cpKeys.CryptoKeys))
//
//	//CPKey, cpKeys.CryptoKeys = cpKeys.CryptoKeys[len(cpKeys.CryptoKeys)-1], cpKeys.CryptoKeys[:len(cpKeys.CryptoKeys)-1]
//	//
//	log.Printf("Устанваливаем ключ cpKeys %+v\n", cpKeys)
//	if err := session.Run("/opt/cprocsp/sbin/amd64/cpconfig -license -set "+strings.TrimSpace(cpKeys)); err != nil {
//		//log.Fatal("Failed to run: " + err.Error())
//		log.Printf("err /cpconfig -license -set %+v %+v %+v\n", err , cpKeys, client)
//		return nil, err
//	}
//
//	checkLic, err := ПолучитьУстановленныйКлюч(client)
//	log.Printf("Получили установленный ключ   %+v, устанавливали ключ %+v Ключи равны ?: %+v\n", checkLic[1], cpKeys, strings.TrimSpace(checkLic[1])==strings.TrimSpace(cpKeys))
//
//	if strings.TrimSpace(checkLic[1]) == strings.TrimSpace(cpKeys) {
//		return checkLic, err
//	} else {
//		return nil, err
//	}
//}
//func ПолучитьМакАдрес(client *ssh.Client)(string, error) {
//
//	var mac bytes.Buffer
//	session, err := client.NewSession()
//
//	if err != nil {
//		log.Printf("err %+v %+v\n", err )
//		return "", err
//	}
//	session.Stdout = &mac
//
//	if err := session.Run("cat /sys/class/net/*/address"); err != nil {
//		log.Printf("Failed to run: " + err.Error())
//		return "", err
//	}
//	macs := strings.Split(mac.String(), "\n")
//	session.Close()
//
//	//log.Printf("macs %+v\n", macs)
//	if macs[0] == "Auth User/Pass with PS...fail...Please reconnect!."{
//		log.Printf("macs %+v\n", macs)
//		err := errors.New("Auth User/Pass with PS...fail...Please reconnect!.")
//		return "", err
//	}
//
//
//	return macs[0], nil
//}

//func ПроверитьВерсиюОС(client *ssh.Client)(string, error){
//	var osV bytes.Buffer
//	session, err := client.NewSession()
//	if err != nil {
//		log.Printf("err %+v %+v\n", err )
//		return "", err
//	}
//	session.Stdout = &osV
//
//	if err := session.Run("lsb_release -d"); err != nil {
//		log.Printf("Не удаёться полуичть версию ОС %+v %+v\n", err )
//		return "", err
//	}
//	osVersion := strings.Split(osV.String(), "\t")
//	//log.Printf("osV %+v\n", osV.String() == "Auth User/Pass with PS...fail...Please reconnect!.")
//	//log.Printf("(IC5) %+v\n", strings.Contains(osV.String(), "(IC5)"))
//	if len(osVersion)<=1 || osV.String() == "Auth User/Pass with PS...fail...Please reconnect!."{
//		//log.Printf("osVersion %+v\n", osVersion)
//		err := errors.New("Не удаёться получить версию ОС")
//		return "", err
//	}
//	//log.Printf("osVersion %+v, osVersion 1 %+v\n", osVersion, osVersion[1])
//	//osVersion := []string{osV.String()}
//	session.Close()
//	return osVersion[1], nil
//
//}