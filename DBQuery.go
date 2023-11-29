package main

import (

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgconn"
	"reflect"
	"strconv"
	"strings"
	txtTpl "text/template"
)

/*
Выполнить вызывает PrepAndExecSql, получет указатель на читатель данных из базы.
Разбирает результат выборки в массив карт
return []map[string]interface{}{
	map[string]interface{}{
		"ИмяПоля": значение поля любого тип
	}
}
Значения полей с типом json (Postgres IOD = 3802) парсяться в interface{}
*/

func (sqlData sqlStruct) Выполнить(ctx context.Context) ([]map[string]interface{}, error) {
	//Инфо("Выполнить  %+v", sqlData.Name)

	var cancel context.CancelFunc
	if ctx == nil{
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
	}


	РезультатЗапроса, err := sqlData.PgSqlResultReader(ctx)

	if err != nil{
			Ошибка("Выполнить %+s \n  Sql %+v", err,  sqlData.Sql)
		for n, v := range  sqlData.Values {
			Инфо(" %+v = %+v", n+1, string(v), sqlData.Values)
		}
		// Если возникла ошибка в исполнении или подготовке запроса то вернём ошибку и nil
		return nil, err
	}

	ОписаниеПолей := РезультатЗапроса.FieldDescriptions
	Строки := РезультатЗапроса.Rows

	if len(Строки) == 0 {
		// Если в резулдьтате выполнения запроса не возникло ошибок проверим сколько было затронуто строк, изменено, вставлено, удалено и т.д
		if РезультатЗапроса.CommandTag.RowsAffected()==0{
			// Запрос не затронул ни одной строки....

			//return []map[string]interface{}{} , errors.New("Запрос не затронул ни одной строки "+sqlData.Sql)
			return []map[string]interface{}{} , nil
		} else {
			// запрос выполнился, было затронуто более нуля строк, но запрос не вернулрезультат видимо не должен
			return []map[string]interface{}{{"Затронуто строк":РезультатЗапроса.CommandTag.RowsAffected()}}, nil
		}
	}

	МассивСтрок := make([]map[string]interface{}, РезультатЗапроса.CommandTag.RowsAffected())

	for i,строка := range Строки {
		КартаСтрок := map[string]interface{}{}
		for номер ,ячейка :=range строка {
			ИмяПоля := string(ОписаниеПолей[номер].Name)

			if ОписаниеПолей[номер].DataTypeOID == 3802{ // JSONB
				if len(ячейка)>0{
					var jsonCell interface{}
					err:=json.Unmarshal(ячейка, &jsonCell)

					if err != nil {
						Ошибка("\n !! ERR %+v Поле %+v string(ячейка) %+v\n", err, string(ОписаниеПолей[номер].Name), ячейка, len(ячейка))
						КартаСтрок[ИмяПоля] = string(ячейка)
					}

					//Инфо(" JSONB КартаСтрок %+v", КартаСтрок)
					КартаСтрок[ИмяПоля]=jsonCell
				} else {
					КартаСтрок[ИмяПоля]=nil
				}
			} else if ОписаниеПолей[номер].DataTypeOID == 16 { // строка

				КартаСтрок[ИмяПоля] = string(ячейка)

			}else if ОписаниеПолей[номер].DataTypeOID == 23 { //целое

				if ячейка == nil {
					КартаСтрок[ИмяПоля] = nil
				}  else {
					КартаСтрок[ИмяПоля], err = strconv.Atoi(string(ячейка))
					if err != nil {
						Ошибка(">>>> ERROR \n %+v ИмяПоля %+v string(ячейка) %+v \n\n", err, ИмяПоля, string(ячейка))
					}
				}
			}else {
				КартаСтрок[ИмяПоля] = string(ячейка)
			}

		}

		МассивСтрок[i]= КартаСтрок

	}

	return МассивСтрок, nil
}


/*
ВыполнитьПолучитьРидер подговтавливает и выполняет запрос, возвращает указатель на ResultReader
*/

func (sqlData sqlStruct) PgSqlResultReader (ctx context.Context) (*pgconn.Result, error) {

	if sqlData.DBSchema == ""{
		sqlData.DBSchema = "iobot"
	}

	var PgConn *pgconn.PgConn
	var ОшибкаПодключения error

	if sqlData.DBSchema == "ecp" {
		Инфо("Подключаемся  %+v", sqlData.DBSchema)
		PgConn, ОшибкаПодключения = MVVPGConnect(sqlData.DBSchema, ctx)
	} else {
		PgConn, ОшибкаПодключения = PGConnect(sqlData.DBSchema, ctx)
	}

	if ОшибкаПодключения != nil {
		Ошибка("err	 %+v\n", ОшибкаПодключения)
		return nil, fmt.Errorf("Нет соединения с БД; %+v \n", ОшибкаПодключения.Error())
	}
	_, err := PgConn.Prepare(ctx, sqlData.Name, sqlData.Sql, nil)
//Инфо("SD %+v", SD)
	if err != nil{
		Ошибка("\n Ошибка подготовки запроса %+v sqlData.Name %+v , sqlData %+v", err, sqlData.Name, sqlData)

		err:=PgConn.Close(ctx)
		if err != nil {
			Ошибка(">>>> ERROR \n %+v \n\n", err)
		}
		return nil, fmt.Errorf("Ошибка подготовки запроса %+v sqlData.Name %+v , sqlData.Sql %+v", err, sqlData.Name, sqlData.Sql)
	} else {
		ResultReader := PgConn.ExecPrepared(ctx, sqlData.Name, sqlData.Values, nil, nil)
		Result := ResultReader.Read()
		var ОшибкаЗапроса *pgconn.PgError
		if Result.Err != nil {
			ОшибкаЗапроса = Result.Err.(*pgconn.PgError)
			Инфо("Code %+v Message %+v Detail %+v sqlData.Values %+v", ОшибкаЗапроса.Code, ОшибкаЗапроса.Message, ОшибкаЗапроса.Detail, sqlData.Values)
		}


		err:=PgConn.Close(ctx)
		if err != nil {
			Ошибка(">>>> ERROR \n %+v \n\n", err)
		}

		return Result, Result.Err
	}
}

func ПолучитьSQLАргументы (Аргументы interface{}, Данные map[string]interface{}) [][]byte{
	//Инфо("  %+v", "ПолучитьSQLАргументы")
	АргументыSQL := make([][]byte, len(Аргументы.([]interface{})))
	var ДанныеПоАдресу interface{}
	for индекс, Адрес := range Аргументы.([]interface{}){
		ШагиАдреса := strings.Split(Адрес.(string), ".")
		ДанныеПоАдресу = Данные
		for номерШагаАдреса, ШагАдреса := range ШагиАдреса {

			//Инфо(" %+v - ШагАдреса  %+v ШагиАдреса %+v",номерШагаАдреса,  ШагАдреса, ШагиАдреса)

			if ШагАдреса != "" {

				//Инфо("  %+v == %+v", номерШагаАдреса, len(ШагиАдреса)-1)
				//Инфо("reflect.TypeOf(%+v).Kind()  %+v индекс %+v ШагАдреса %+v ",ДанныеПоАдресу, reflect.TypeOf(ДанныеПоАдресу).Kind(), индекс, ШагАдреса)


				switch reflect.TypeOf(ДанныеПоАдресу).Kind() {
				case reflect.Int:
					АргументыSQL[индекс] = []byte(strconv.Itoa(ДанныеПоАдресу.(int)))
				case reflect.Ptr:
					//reflect.New(reflect.TypeOf(ДанныеПоАдресу))
					//Инфо(" reflect.Ptr  %+v", reflect.TypeOf(ДанныеПоАдресу).Elem().Kind())

					switch reflect.TypeOf(ДанныеПоАдресу).Elem().Kind(){
					case reflect.Struct :
						if номерШагаАдреса == len(ШагиАдреса)-1 {
							ЗначениеПоля := reflect.ValueOf(ДанныеПоАдресу).Elem().FieldByName(ШагАдреса)
							//Инфо("  ЗначениеПоля %+v, ЗначениеПоля.Kind() %+v", ЗначениеПоля, ЗначениеПоля.Kind())
							switch ЗначениеПоля.Kind(){
							case  reflect.String:
								АргументыSQL[индекс] = []byte(ЗначениеПоля.String())
							case  reflect.Int:
								АргументыSQL[индекс] = []byte(strconv.Itoa(int(ЗначениеПоля.Int())))
							case reflect.Map:
								var err error
								АргументыSQL[индекс], err = json.Marshal(ДанныеПоАдресу)
								if err != nil {
									Ошибка("  %+v", err)
								}
							}

						} else {
							ДанныеПоАдресу = reflect.ValueOf(ДанныеПоАдресу).Elem().FieldByName(ШагАдреса)
						}
					}
				case reflect.Interface :
					if номерШагаАдреса == len(ШагиАдреса)-1 {

						АргументыSQL[индекс] = reflect.ValueOf(ДанныеПоАдресу).Bytes()
					} else {
						//Инфо(" reflect.Interface %+v",ДанныеПоАдресу)

						ДанныеПоАдресу = reflect.ValueOf(ДанныеПоАдресу).Elem()
						//Инфо(" reflect.Interface ДанныеПоАдресу Elem() %+v", ДанныеПоАдресу)
					}

					//Инфо("ДанныеПоАдресу reflect.Interface  %+v", ДанныеПоАдресу)

				case reflect.Map :

					//Инфо("reflect.Map  %+v",ДанныеПоАдресу)

					ВложенныеЭлемент := reflect.TypeOf(ДанныеПоАдресу).Elem()
					//Инфо(" reflect.Map; ВложенныеЭлемент %+v ; ВложенныеЭлемент.Kind() %+v",ВложенныеЭлемент,  ВложенныеЭлемент.Kind())

					//Инфо("  %+v = %+v; %+v, ДанныеПоАдресу %+v", номерШагаАдреса , len(ШагиАдреса)-1, ШагАдреса, ДанныеПоАдресу)

					switch ВложенныеЭлемент.Kind() {
					case reflect.String:
						АргументыSQL[индекс] = []byte( ДанныеПоАдресу.(map[string]string)[ШагАдреса])
					case reflect.Interface:
						//Инфо("reflect.Interface: %+v", )
						if номерШагаАдреса == len(ШагиАдреса)-1 {
							var err error
							var ЕстьДанные bool
							ДанныеПоАдресу , ЕстьДанные =  ДанныеПоАдресу.(map[string]interface{})[ШагАдреса]
							//Инфо("ДанныеПоАдресу  %+v %+v ЕстьДанные %+s", ДанныеПоАдресу, reflect.TypeOf(ДанныеПоАдресу).Kind(), ЕстьДанные)

							if ДанныеПоАдресу == nil && !ЕстьДанные {
								АргументыSQL[индекс] = nil
							} else {

								//Инфо("ДанныеПоАдресу  %+s; ШагАдреса %+s ; reflect.TypeOf(ДанныеПоАдресу).Kind() %+v",ДанныеПоАдресу, ШагАдреса, reflect.TypeOf(ДанныеПоАдресу).Kind())
								switch reflect.TypeOf(ДанныеПоАдресу).Kind() {

								case reflect.Slice:
									//Инфо("reflect.Slice ШагАдреса=%+s ДанныеПоАдресу  %+v",ШагАдреса,  ДанныеПоАдресу)
									АргументыSQL[индекс], err = json.Marshal(ДанныеПоАдресу)
									if err != nil {
										Ошибка("  %+v", err)
									}
								case reflect.Map:
									//Инфо(" reflect.Map %+s", ДанныеПоАдресу)
									if ДанныеПоАдресу == nil {
										АргументыSQL[индекс]=nil
									}else{
										АргументыSQL[индекс], err = json.Marshal(ДанныеПоАдресу)
										//Инфо("АргументыSQL[%+v]  %+s", ШагАдреса, АргументыSQL[индекс])
										if err != nil {
											Ошибка("  %+v", err)
										}
									}
								case reflect.String:
									//Инфо(" reflect.String %+s", ДанныеПоАдресу)
									if ДанныеПоАдресу.(string) == ""{
										АргументыSQL[индекс]= nil
									}else {
										АргументыSQL[индекс] = []byte(ДанныеПоАдресу.(string))
									}
								case reflect.Float64:
									if ДанныеПоАдресу.(float64) == 0{
										АргументыSQL[индекс]= nil
									}else {
										АргументыSQL[индекс] = []byte(strconv.Itoa(int(ДанныеПоАдресу.(float64))))
									}
								case reflect.Int:
									Инфо(" reflect.Int %+s", ДанныеПоАдресу)
									if ДанныеПоАдресу.(int) == 0{
										АргументыSQL[индекс]= nil
									}else {
										АргументыSQL[индекс] = []byte(strconv.Itoa(ДанныеПоАдресу.(int)))
									}
								case reflect.Ptr:
									//Инфо(" reflect.Ptr %+s", ДанныеПоАдресу)
									АргументыSQL[индекс], err = json.Marshal(ДанныеПоАдресу)
									//Инфо("АргументыSQL[%+v]  %+s", ШагАдреса, АргументыSQL[индекс])
									if err != nil {
										Ошибка("  %+v", err)
									}
								}
							}
						} else {
							ДанныеПоАдресу = ДанныеПоАдресу.(map[string]interface{})[ШагАдреса]
						}

					}

				case reflect.String:

					//ДанныеПоАдресу = ДанныеПоАдресу.(string)
					АргументыSQL[индекс] = []byte(ДанныеПоАдресу.(string))
				}

			}
		}

		//Инфо(" АргументыSQL %+v", АргументыSQL)


		//tpl, err := txtTpl.New("АргументыЗапроса").Funcs(tplFunc()).Parse("{{"+Адрес.(string)+"}}")
		//if err != nil {
		//	Ошибка(" %+v ", err)
		//}
		//
		//БайтБуферДанных := new(bytes.Buffer)
		//
		//err = tpl.Execute(БайтБуферДанных, Данные)
		//if err != nil {
		//	Ошибка(" %+v ", err)
		//}
		////Инфо("БайтБуферДанных.String()  %+s", БайтБуферДанных.String())
		//
		//if БайтБуферДанных.String() == "<no value>" || БайтБуферДанных.String() == `""`{
		//	SQLАргументы[индекс] = nil
		//} else {
		//
		//	SQLАргументы[индекс] = БайтБуферДанных.Bytes()
		//}

		//Инфо("SQLАргументы[индекс] %+s", SQLАргументы[индекс])
	}
	return АргументыSQL
}

func ПолучитьSQLАргументыДляАис (Аргументы interface{}, Данные map[string]interface{}) []interface{}{
	//Инфо("  %+v", "ПолучитьSQLАргументыДляАис", Аргументы)

	АргументыSQL := make([]interface{}, len(Аргументы.([]interface{})))

	var ДанныеПоАдресу interface{}
	for индекс, Адрес := range Аргументы.([]interface{}){
		//Инфо("Адрес  %+v", Адрес)
		ШагиАдреса := strings.Split(Адрес.(string), ".")
		ДанныеПоАдресу = Данные
		for номерШагаАдреса, ШагАдреса := range ШагиАдреса {

			//Инфо("номерШагаАдреса %+v - ШагАдреса  %+v",номерШагаАдреса,  ШагАдреса)

			if ШагАдреса != "" {

				//Инфо("  %+v == %+v", номерШагаАдреса, len(ШагиАдреса)-1)
				//Инфо("ДанныеПоАдресу %+v reflect.TypeOf(%+v).Kind()  %+v ",ДанныеПоАдресу)

				if ДанныеПоАдресу == nil{
					АргументыSQL[индекс] = nil
					continue
				}

				switch reflect.TypeOf(ДанныеПоАдресу).Kind() {
				case reflect.Int:
					АргументыSQL[индекс] = ДанныеПоАдресу.(int)
				case reflect.Ptr:
					//reflect.New(reflect.TypeOf(ДанныеПоАдресу))
					//Инфо(" reflect.Ptr  %+v", reflect.TypeOf(ДанныеПоАдресу).Elem().Kind())

					switch reflect.TypeOf(ДанныеПоАдресу).Elem().Kind(){
					case reflect.Struct :
						if номерШагаАдреса == len(ШагиАдреса)-1 {
							ЗначениеПоля := reflect.ValueOf(ДанныеПоАдресу).Elem().FieldByName(ШагАдреса)
							//Инфо("  ЗначениеПоля %+v, ЗначениеПоля.Kind() %+v", ЗначениеПоля, ЗначениеПоля.Kind())
							switch ЗначениеПоля.Kind(){
							case  reflect.String:
								АргументыSQL[индекс] = ЗначениеПоля.String()
							case  reflect.Int:
								АргументыSQL[индекс] = ЗначениеПоля.Int()
							case reflect.Map:
								var err error
								АргументыSQL[индекс], err = json.Marshal(ДанныеПоАдресу)
								if err != nil {
									Ошибка("  %+v", err)
								}
							}

						} else {
							ДанныеПоАдресу = reflect.ValueOf(ДанныеПоАдресу).Elem().FieldByName(ШагАдреса)
						}
					}
				case reflect.Interface :
					if номерШагаАдреса == len(ШагиАдреса)-1 {

						АргументыSQL[индекс] = reflect.ValueOf(ДанныеПоАдресу).Bytes()
					} else {
						//Инфо(" reflect.Interface %+v",ДанныеПоАдресу)

						ДанныеПоАдресу = reflect.ValueOf(ДанныеПоАдресу).Elem()
						//Инфо(" reflect.Interface ДанныеПоАдресу Elem() %+v", ДанныеПоАдресу)
					}

					//Инфо("ДанныеПоАдресу reflect.Interface  %+v", ДанныеПоАдресу)

				case reflect.Map :

					//Инфо("reflect.Map  %+v",ДанныеПоАдресу)

					ВложенныеЭлемент := reflect.TypeOf(ДанныеПоАдресу).Elem()
					//Инфо(" reflect.Map; ВложенныеЭлемент %+v ; ВложенныеЭлемент.Kind() %+v",ВложенныеЭлемент,  ВложенныеЭлемент.Kind())

					//Инфо("  %+v = %+v; %+v, ДанныеПоАдресу %+v", номерШагаАдреса , len(ШагиАдреса)-1, ШагАдреса, ДанныеПоАдресу)

					switch ВложенныеЭлемент.Kind() {
					case reflect.String:
						АргументыSQL[индекс] = ДанныеПоАдресу.(map[string]string)[ШагАдреса]
					case reflect.Interface:
						if номерШагаАдреса == len(ШагиАдреса)-1 {
							var err error
							var ЕстьДанные bool
							ДанныеПоАдресу , ЕстьДанные =  ДанныеПоАдресу.(map[string]interface{})[ШагАдреса]
							//Инфо("ДанныеПоАдресу  %+v %+v ЕстьДанные %+s", ДанныеПоАдресу, reflect.TypeOf(ДанныеПоАдресу).Kind(), ЕстьДанные)

							if ДанныеПоАдресу == nil && !ЕстьДанные {
								АргументыSQL[индекс] = nil
							} else {

								//Инфо("ДанныеПоАдресу  %+s; ШагАдреса %+s ; reflect.TypeOf(ДанныеПоАдресу).Kind() %+v",ДанныеПоАдресу, ШагАдреса, reflect.TypeOf(ДанныеПоАдресу).Kind())
								switch reflect.TypeOf(ДанныеПоАдресу).Kind() {

								case reflect.Slice:
									//Инфо("reflect.Slice ШагАдреса=%+s ДанныеПоАдресу  %+v",ШагАдреса,  ДанныеПоАдресу)
									АргументыSQL[индекс], err = json.Marshal(ДанныеПоАдресу)
									if err != nil {
										Ошибка("  %+v", err)
									}
								case reflect.Map:
									//Инфо(" reflect.Map %+s", ДанныеПоАдресу)
									if ДанныеПоАдресу == nil {
										АргументыSQL[индекс]=nil
									}else{
										АргументыSQL[индекс], err = json.Marshal(ДанныеПоАдресу)
										//Инфо("АргументыSQL[%+v]  %+s", ШагАдреса, АргументыSQL[индекс])
										if err != nil {
											Ошибка("  %+v", err)
										}
									}
								case reflect.String:
									//Инфо(" reflect.String %+s", ДанныеПоАдресу)
									if ДанныеПоАдресу.(string) == ""{
										АргументыSQL[индекс]= nil
									}else {
										АргументыSQL[индекс] = ДанныеПоАдресу.(string)
									}
								case reflect.Float64:
									if ДанныеПоАдресу.(float64) == 0{
										АргументыSQL[индекс]= nil
									}else {
										АргументыSQL[индекс] = ДанныеПоАдресу.(float64)
									}
								case reflect.Int:
									Инфо(" reflect.Int %+s", ДанныеПоАдресу)
									if ДанныеПоАдресу.(int) == 0{
										АргументыSQL[индекс]= nil
									}else {
										АргументыSQL[индекс] = ДанныеПоАдресу.(int)
									}
								case reflect.Ptr:
									//Инфо(" reflect.Ptr %+s", ДанныеПоАдресу)
									АргументыSQL[индекс], err = json.Marshal(ДанныеПоАдресу)
									//Инфо("АргументыSQL[%+v]  %+s", ШагАдреса, АргументыSQL[индекс])
									if err != nil {
										Ошибка("  %+v", err)
									}
								}
							}
						} else {
							ДанныеПоАдресу = ДанныеПоАдресу.(map[string]interface{})[ШагАдреса]
						}

					}

				case reflect.String:

					//ДанныеПоАдресу = ДанныеПоАдресу.(string)
					АргументыSQL[индекс] = ДанныеПоАдресу.(string)
				}

			}
		}

		//Инфо(" АргументыSQL %+s", АргументыSQL)

	}
	return АргументыSQL
}
func ПолучитьДанныеИзАИСвПотоке(SqlData sqlStruct, канал chan interface{}) {
	РезультатЗапроса ,ОшибкаЗапроса := SqlData.ВыполнитьЗапросВАИС()
	if ОшибкаЗапроса != nil {
		Ошибка("  %+v", ОшибкаЗапроса)
		канал <- ОшибкаЗапроса
	}

	канал <- РезультатЗапроса

	Инфо(" выходи из Функции ПолучитьДанныеИзАИСвПотоке : %+v", SqlData.Name)
}

func КэшироватьДанные(ИмяЗапроса string, Пользователь Client, Данные interface{},база string){

	ДанныеПользователя := map[string]interface{}{
		"Логин": Пользователь.Login,
		"IP": Пользователь.Ip,
		"информация": Пользователь.UserInfo,
	}

	пользователь, err := json.Marshal(ДанныеПользователя)
	if err != nil {
		Ошибка("  %+v", err)
	}

	данные, err := json.Marshal(Данные)
	if err != nil {
		Ошибка("  %+v", err)
	}
		_ ,err = sqlStruct{
			Name:   "Кэширование данных из АИС",
			Sql:    `INSERT INTO logs.кэш_аис_запросов (имя_запроса, дата, пользователь, данные, база) VALUES ($1,NOW(),$2,$3,$4)`,
			Values: [][]byte{
				[]byte(ИмяЗапроса),
				пользователь,
				данные,
				[]byte(база),
			},
		}.Выполнить(nil)
	if err != nil {
		Ошибка("  %+v", err)
	}
}

func (client *Client) ОбработатьБезопасноSQLСкрипты(SQLЗапросыInterface []interface{}, вопрос *Сообщение) map[string]interface{} {


	//Инфо("ОбработатьБезопасноSQLСкрипты SQLСкрипт %+v\n", SQLЗапросыInterface)

	SQLЗапросы := SQLЗапросыInterface//.([]interface{})//["sql"].(map[string]interface{})
	//SQLСкрипт, ЕстьSqlСкрипт := SQLЗапрос.(map[string]interface{})//["sql"].(map[string]interface{})

	Данные:= map[string]interface{}{
		"client" : client.UserInfo.Info,
		"data": map[string]interface{}{},
		"РезультатВсехЗапросов": map[string]interface{}{},
	}

	// Засунем все данные пришедшие от клиента в одну карту,
	if вопрос.Выполнить.Действие !=nil{
		for НазваниеДействия, ДанныеДействия := range вопрос.Выполнить.Действие{
			Инфо("НазваниеДействия %+v ДанныеДействия %+v\n", НазваниеДействия, ДанныеДействия)
			//Инфо("reflect.TypeOf(ДанныеДействия) %+v", reflect.TypeOf(ДанныеДействия).Kind())
			Данные["data"]=ДанныеДействия
		}
	}

	if len(SQLЗапросы)==0{
		Инфо(" нет SQL срикптов. SQLЗапросы: %+v\n", SQLЗапросы)
		return nil
	} else {
		РезультатВсехЗапросов := map[string]interface{}{}

		for _ , SQLСкрипт := range SQLЗапросы{
			if SQLСкрипт == nil {
				continue
			}
			for _, КартаСкрипта := range SQLСкрипт.(map[string]interface{}) {
				//Инфо("Очерёдность %+v ", Очерёдность, )
				for ИмяСкрипта, ДанныеЗапроса := range КартаСкрипта.(map[string]interface{}){
					Инфо(" ИмяСкрипта %+v ",ИмяСкрипта)
					//Инфо(" Выполним Sql %+v Данные %+v", ДанныеЗапроса.(map[string]interface{}), Данные)

					Скрипт :=ДанныеЗапроса.(map[string]interface{})["скрипт"]
					БазаДанных := ДанныеЗапроса.(map[string]interface{})["база_данных"]
					ОбработчикДанных := ДанныеЗапроса.(map[string]interface{})["обработчик"]
					ДанныеОбработчка := ДанныеЗапроса.(map[string]interface{})["данные_обработчика"]
					комментарий := ДанныеЗапроса.(map[string]interface{})["комментарий"]

					if комментарий != nil {
						if client != nil{
							СообщениеКлиенту := Сообщение{
								Id:      0,
								От:      "io",
								Кому:    client.Login,
								Текст:   комментарий.(string),
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
							//Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)
							СообщениеКлиенту.СохранитьИОтправить(client)
						}
					}


					//Инфо("ДанныеОбработчка  %+v", ДанныеОбработчка)
					//Инфо("БазаДанных  %+v", БазаДанных)
					//Обработчик := ДанныеЗапроса.(map[string]interface{})["обработчик"]
					Аргументы := ДанныеЗапроса.(map[string]interface{})["аргументы"] // аргументы которые нужно получить из данных
						//var SQLАргументы [][]byte
						var АргументыSQL  [][]byte
						var АргументыSQLдляАИС  []interface{}



						if Аргументы != nil{
							//SQLАргументы = make([][]byte, len(Аргументы.([]interface{})))
							if БазаДанных == "fssp" ||  БазаДанных =="osp" ||  БазаДанных =="rdb"||  БазаДанных =="osp_rbd" ||  БазаДанных =="osp_code" {
								АргументыSQLдляАИС = ПолучитьSQLАргументыДляАис (Аргументы, Данные)
							} else if БазаДанных == "io" {
								АргументыSQL = ПолучитьSQLАргументы (Аргументы, Данные)
							} else if БазаДанных == "zabbix" {

							}

						}
						//Инфо("SQLАргументы %+s", SQLАргументы)

						if Динамический := ДанныеЗапроса.(map[string]interface{})["динамический"];  Динамический!=nil && Динамический.(bool)  {
							/* Если запрос динамический то соберём его*/
							if ДинамическийШаблон := ДанныеЗапроса.(map[string]interface{})["динамический_шаблон"]; ДинамическийШаблон!=nil {
								Скрипт = ДинамическийШаблон.(string) + Скрипт.(string)
							}
								Инфо("Скрипт.(string)  %+v", Скрипт.(string), Данные)

								tpl, err := txtTpl.New("SqlЗапрос").Funcs(tplFunc()).Parse(Скрипт.(string))
								if err != nil {
									Ошибка("err %+v tpl %+v ", err, tpl)
								}

								БайтБуферДанных := new(bytes.Buffer)
								Инфо("Данные передаваемые в динамический sql %+v", Данные)
								err = tpl.Execute(БайтБуферДанных, Данные)
								if err != nil {
									Ошибка(" %+v ", err)
								} else {
									Скрипт = БайтБуферДанных.String()
									Скрипт = strings.TrimSpace(Скрипт.(string))
								}
								Инфо("Скрипт  %+v \n БайтБуферДанных %+s \n АргументыSQL %+s ", Скрипт, БайтБуферДанных, АргументыSQL)
						}

						var ОтветЗабикса ЗабиксРезультат
						var РезультатЗапроса []map[string]interface{}
						var РезультатЗапросаИзАИС  map[string]map[string]РезультатSQLизАИС //[ИмяЗапроса][ОСП]РезультатSQLизАИС
						var ОшибкаЗапроса interface{}
						if БазаДанных != "io" && БазаДанных != "zabbix" {
							/* Подключаемся к аис */

							var Асинхронно bool // если true то не ожидаем завершения а продолжаем обработку данных

							if ДанныеЗапроса.(map[string]interface{})["асинхронно"] != nil {

									Асинхронно = ДанныеЗапроса.(map[string]interface{})["асинхронно"].(bool)

									канал := make(chan interface{})
									go ПолучитьДанныеИзАИСвПотоке(sqlStruct{
										Name:            ИмяСкрипта,
										Sql:             Скрипт.(string),
										АргументыДляАИС: АргументыSQLдляАИС,
										БазаДанных:БазаДанных.(string),
										Клиент: client,
										Асинхронно: Асинхронно,
									}, канал)

									go ОжиданиеДанныхИзАИС (канал)

							} else {
								Инфо("АргументыSQLдляАИС %+v", АргументыSQLдляАИС)

								РезультатЗапросаИзАИС ,ОшибкаЗапроса = sqlStruct{
									Name:            ИмяСкрипта,
									Sql:             Скрипт.(string),
									АргументыДляАИС: АргументыSQLдляАИС,
									БазаДанных:БазаДанных.(string),
									Клиент: client,
								}.ВыполнитьЗапросВАИС() // Идея: Может нужно обработчикДанных засунуть внутрь, чтобы данные из осп обрабатывались по мере получения... А если данные нужно обработать не по осп а всей кучей....?
							}

						} else if  БазаДанных == "io" {
							РезультатЗапроса ,ОшибкаЗапроса = sqlStruct{
								Name:   ИмяСкрипта,
								Sql:    Скрипт.(string),
								Values: АргументыSQL,
							}.Выполнить(nil)
						} else if БазаДанных == "zabbix" {
							ЗабиксАвторизация()
							ОтветЗабикса ,ОшибкаЗапроса = ЗапросВЗабикс(Скрипт.(string))
						}


						if ОшибкаЗапроса != nil {

							Ошибка("%+v", ОшибкаЗапроса)

							Ошибка("%+v ОшибкаЗапроса %#v", reflect.TypeOf(ОшибкаЗапроса).Kind(), ОшибкаЗапроса)

							 var ДеталиОшибкиЗапроса *pgconn.PgError

								//ДеталиОшибкиЗапроса = ОшибкаЗапроса.(*pgconn.PgError)

							Инфо("ДеталиОшибкиЗапросаPG %+v  %+v", ДеталиОшибкиЗапроса == nil, reflect.TypeOf(ДеталиОшибкиЗапроса).Kind())
							//Инфо(" %+v", ДеталиОшибкиЗапроса.Code, ДеталиОшибкиЗапроса.ColumnName, ДеталиОшибкиЗапроса.Detail)
							//if strings.Contains(ОшибкаЗапроса.Error(), "SQLSTATE 23505"){
							if ДеталиОшибкиЗапроса != nil {
								if ДеталиОшибкиЗапроса.Code == "23505" {
										Инфо(" %+v %+v %+v", ДеталиОшибкиЗапроса.Code, ДеталиОшибкиЗапроса.ColumnName, ДеталиОшибкиЗапроса.Detail)
										Ошибка("НЕТ ОБРАБОТЧИКА SQLSTATE 23505 %+v ", ОшибкаЗапроса)
										//алгоритм. Если запрос вернул ошибку с ограничением целостности уникального значения , тогда проверим есть ли обработчик SQLSTATE. и Наверное пока sql для обрабтки ошибок будет возвращать сразу текст который будет выводиться клиенту,
										//SQLStateОбъект, ЕстьSQLStateОбъект := SQLСкрипт.(map[string]interface{})["SQLState"]
										//
										//if !ЕстьSQLStateОбъект{
										//	log.Printf("ЕстьSQLStateОбъект %+v\n", ЕстьSQLStateОбъект)
										//}
										//
										//ОбработчикОшибки := SQLStateОбъект.(map[string]interface{})["23505"]
										//
										//if ОбработчикОшибки != nil{
										////	/*алгоритм Функция ОбработатьОщибку, получает на вход sql строку, и данные для подстановки в запрос, Функция сама отправляет сообщение клиенту если обработчик умпешно отработал и вернул данные с полем ,message*/
										//	client.ОбработатьОшибкуБД (ОбработчикОшибки , Данные)
										//}
										РезультатВсехЗапросов[ИмяСкрипта] = map[string]interface{}{
											"ошибка": "Данные НЕ были сохраннены! Причина: Повторяющееся значения уникального поля " + ДеталиОшибкиЗапроса.Detail,
									}

									} else {
										РезультатВсехЗапросов[ИмяСкрипта] = map[string]interface{}{
											"ошибка" : ДеталиОшибкиЗапроса.Code +" " + ДеталиОшибкиЗапроса.Detail,
										}
									}
							}
							//Ошибка("НЕТ ОБРАБОТЧИКА SQLSTATE; ОшибкаЗАпроса %+v",  ОшибкаЗапроса)
							//ДанныеДляОтвета := &ДанныеОтвета{
							//	Обработчик:"FloatMessage",
							//	Данные: "Не удалось сохранить файл "+ ОшибкаЗапроса.Error(),
							//}
							//
							//СообщениеКлиенту:= &Сообщение{
							//	От: "io",
							//	Кому:client.Login,
							//	MessageType: []string{"error"},
							//	Контэнт:ДанныеДляОтвета,
							//}
							//client.Message<-СообщениеКлиенту
						} else{
							if len(РезультатЗапроса) == 1{
								//Инфо(" РезультатЗапроса %+v", РезультатЗапроса)
								if сообщениеКлиенту, ЕстьСообщение :=  РезультатЗапроса[0]["message"];ЕстьСообщение {
									//ДанныеДляОтвета := &ДанныеОтвета{
									//	Обработчик:"FloatMessage",
									//	Данные:сообщениеКлиенту,
									//}
									СообщениеКлиенту:= &Сообщение{
										От: "io",
										Кому:client.Login,
										MessageType: []string{"note"},
										Контэнт:&ДанныеОтвета{
											Обработчик:"FloatMessage",
											Данные:сообщениеКлиенту,
										},
									}
									Инфо("СообщениеКлиенту %+v\n", СообщениеКлиенту)
									СообщениеКлиенту.СохранитьИОтправить(client)
								}
							}


							if ОбработчикДанных != nil{

								Инфо(" ИмяСкрипта %+v", ИмяСкрипта)
								Инфо(" ОбработчикДанных %+v", ОбработчикДанных)
								//Инфо(" Данные %+v", Данные)
								//Инфо(" РезультатЗапросаИзАИС %+v", РезультатЗапросаИзАИС)
								//Инфо(" РезультатЗапросаИзАИС[ИмяСкрипта].КартаСтрок %+v", РезультатЗапросаИзАИС[ИмяСкрипта])
								Инфо(" ОтветЗабикса != nil %+v", ОтветЗабикса.Результат != nil)

								if РезультатЗапросаИзАИС != nil {
									РезультатВсехЗапросов[ИмяСкрипта] = ОбработатьДанные(ОбработчикДанных.(string), ДанныеОбработчка, client, РезультатЗапросаИзАИС, Данные) // [ИмяЗапроса][ОСП]РезультатSQLизАИС
								}
								if РезультатЗапроса != nil  {
									РезультатВсехЗапросов[ИмяСкрипта] = ОбработатьДанные(ОбработчикДанных.(string), ДанныеОбработчка, client, РезультатЗапроса, Данные)
								}
								if ОтветЗабикса.Результат != nil  {
									РезультатВсехЗапросов[ИмяСкрипта] = ОбработатьДанные(ОбработчикДанных.(string), ДанныеОбработчка, client, ОтветЗабикса.Результат, Данные)
								}
							} else if РезультатЗапросаИзАИС != nil {
								for имяСкрипта, ДанныеИзЗапроса := range РезультатЗапросаИзАИС {
									РезультатВсехЗапросов[имяСкрипта] = ДанныеИзЗапроса
								}

							} else if РезультатЗапроса != nil {
								РезультатВсехЗапросов[ИмяСкрипта] = РезультатЗапроса
							} else if ОтветЗабикса.Результат != nil {
								РезультатВсехЗапросов[ИмяСкрипта] = ОтветЗабикса.Результат
							}
						}
						Данные["РезультатВсехЗапросов"]=РезультатВсехЗапросов

					}
			}
		}

		return РезультатВсехЗапросов
	}
}


func ОжиданиеДанныхИзАИС (канал chan interface{}){
	for {
		данные := <- канал
		if данные != nil {
			Инфо("ОжиданиеДанныхИзАИС: Результат: %+v", данные)
		}
	}
}


func ПолучитьСписокСетей() []map[string]interface{}{

	Результат ,err:= sqlStruct{
		Name:   "osp_address",
		Sql:    `SELECT * FROM fssp_configs.osp_address WHERE start_ip IS NOT NULL`,
		Values: [][]byte{},
		DBSchema:"fssp_configs",
	}.Выполнить(nil)
	if err != nil{
		Ошибка(">>>> Ошибка SQL запроса: %+v \n\n",err)
	}
	return Результат
}

func ПолучитьДанныеОСП() (map[string]string, error) {
	sqlSting := "SELECT osp_code,osp_name FROM fssp_configs.osp_address WHERE pwd is not null and osp_code > 26"

	Результат, err := sqlStruct{
		Name:   "Данные ОСП",
		Sql:    sqlSting,
		//Values:	[][]byte{[]byte(stringByte)},
		DBSchema:"fssp_configs",
	}.Выполнить(nil)

	if err != nil {
		Ошибка(" %+v ", err)
		return nil, err
	}

	ОСП := map[string]string{}

	for _, ДанныеОСП := range Результат{
		ОСП[strconv.Itoa(ДанныеОСП["osp_code"].(int))]=ДанныеОСП["osp_name"].(string)
	}

	return ОСП, err
}

func ПолучитьДанныеПодключенийАИСОСП()([]map[string]interface{}, error){
	sqlSting := "SELECT osp_code,osp_name,ip_ais,pwd  FROM fssp_configs.osp_address WHERE pwd is not null and osp_code > 26 and osp_code <> 26000 and osp_code <> 26911"

	sqlQuery, err := sqlStruct{
		Name:   "OSPConnections",
		Sql:    sqlSting,
		//Values:	[][]byte{[]byte(stringByte)},
		DBSchema:"fssp_configs",
	}.Выполнить(nil)
	if err != nil {
		Ошибка(" %+v ", err)
	}
	return sqlQuery, err
}

