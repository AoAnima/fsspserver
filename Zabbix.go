package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	txtTpl "text/template"
)
type ЗабиксДанныеАвторизации struct {
	Логин string `json:"login"`
	Пароль string `json:"password"`
}
type ЗабиксРезультат struct {
	Jsonrpc string `json:"jsonrpc"`
	Результат interface{} `json:"result"`
	Ид int `json:"id"`
}
//type ЗабиксЗапрос type
//{
	//Jsonrpc string `json:"jsonrpc"`
	//Ид int `json:"id"`
	//Метод string `json:"method"`
	//Параметры map[string]interface{} `json:"params"`
	//Авторизация *string `json:"auth, omitempty,_"`

//}




func ЗапросВЗабикс(запрос string) (ЗабиксРезультат, error) {

	/*
	{
			Auth:{{.}},
			Id:1,
			Jsonrpc:"2.0",
			Method:"item.get",
			Params:map[string]interface{}{
				"filter":map[string]interface{}{
					"key_":[]string{"hw.addr","create.date","create.date.cprocsp","cryptopro.license","system.uname","	cryptopro_license"},
					//"key_":"2c:4d:54:53:10:61",
				},
				"output": []string{"hostid","key_","lastvalue",},
			},
		}
	*/
	Инфо("zabix  запрос %+v", запрос)


	tpl, err := txtTpl.New("ЗаббиксЗапрос").Parse(запрос)
	if err != nil {
		Ошибка("err %+v tpl %+v ", err, tpl)
	}

	БайтБуферДанных := new(bytes.Buffer)
	Инфо("Авторизация %+v", Авторизация)
	err = tpl.Execute(БайтБуферДанных, Авторизация)
	if err != nil {
		Ошибка(" %+v ", err)
	} else {
		запрос = БайтБуферДанных.String()
		запрос = strings.TrimSpace(запрос)
	}

	Инфо("zabix  запрос %+v", запрос)
	//jsonQuery, err := json.Marshal(запрос)
	//if err != nil{
	//	log.Printf("err	 %+v\n", err)
	//}

	отправитьЗапрос, err := http.NewRequest("GET", "http://10.26.6.28/zabbix/api_jsonrpc.php", strings.NewReader(запрос))

	отправитьЗапрос.Header.Set("content-type", "application/json")
	//proxyUrl, _ := url.Parse("http://root:Servroot88#@10.26.12.134:8080")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			//Proxy:           http.ProxyURL(proxyUrl),
		},
	}
	ответ , err := client.Do(отправитьЗапрос)
	Инфо(" %+v", ответ)

	//log.Printf("jsonQuery %+v\n", string(jsonQuery))
	//resp, err:=http.Post( "http://10.26.6.28/zabbix/api_jsonrpc.php", "application/json",bytes.NewBuffer(jsonQuery))

	if err != nil {
		Ошибка("err	 %+v\n", err)
	}
	//log.Printf("resp %+v\n", resp, err)
	defer ответ.Body.Close()

	body, err := ioutil.ReadAll(ответ.Body)

	//Инфо("ответ из забикса %+v\n", string(body), err)
	var zabixResult ЗабиксРезультат
	err = json.Unmarshal(body, &zabixResult)
	Инфо("zabixResult %+v", zabixResult)
	if nil != err {
		Ошибка("Не удаёться прочитать тело ответа", err)
	}
	//log.Printf("resp %+v\n", string(body))
	return zabixResult,err
}

var Авторизация string

func ЗабиксАвторизация(){

	if Авторизация != ""  {
		return
	}


	query := `{
		"id":1,
		"jsonrpc":"2.0",
		"method":"user.login",
		"params":{
			"user":"Maksimchuk",
			"password":"Aa111111"
		},
 		"auth": null
	}`
	Инфо("sendZabbix query%+v\n", query)
	ОтветИзЗабикса,err := ЗапросВЗабикс(query)
	Инфо("ответ из забикса при авторизации %+v\n", ОтветИзЗабикса)
	//response = {"jsonrpc":"2.0","result":"4da0def83b1bf30d63c70e993bd244b7","id":1}

	//var ОтветИзЗабикса ЗабиксРезультат
	if err!=nil{
		Ошибка("err %+v\n", err)
	}
	//else {
	//	_= json.Unmarshal(ответ, &ОтветИзЗабикса)
	//}
	Авторизация = ОтветИзЗабикса.Результат.(string)
	//return ОтветИзЗабикса.Результат.(string)
}
