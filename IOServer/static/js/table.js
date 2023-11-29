
// type Table struct {
//     Setting map[string]interface{} `json:"setting"`
//     Id int `json:"id,omitempty"`
//     Name string `json:"name,omitempty"`
//     Description map[string]interface{} `json:"description,omitempty"`
//     Order int `json:"orders,omitempty"`
//     Alias string `json:"alias,omitempty"`
//     Columns		map[string]string `json:"columns"`
//     Header		[]map[string]interface{} `json:"header"`
//     Rows  		[]map[string]interface{} `json:"rows"`
//     Footer  		[]map[string]interface{} `json:"footer"`
//     CountRows uint32 `json:"countRows"`
//     Index 		map[string]map[string]int	`json:"index"` //interface{}
// }

// testLiteral();
// function testLiteral(){
//     start =  Date.now();
//     str = "dsfsdf asdf sdf sdf s ${value}";
//     result =[];
//     let value =1;
//     for (let i=0; i <=2;i++){
//         // pos[i] = str.indexOf("${");
//         // str.
//         let startVariable;
//         [str].map(ch=> console.log(ch))
//
//         // cellValue = new Function ("value","return `"+str+"`");
//         // cellValue.apply(null, [value]);
//         // result[i] = cellValue(value);
//     }
//     stop =  Date.now()
//      console.log(result);
//      console.log(stop-start);
// }
function tpl(tplName, tplValues){
    const tplMap = {

        dataSource: function (tplValues) {
            return `<div class="fieldset">
                                    <select name="source.${tplValues.num}" id="source${tplValues.num}" >
                                        <optgroup  label = "Источник данных">
                                        <option disabled hidden selected="selected" value="">Источник данных</option>
                                        <option value="manual">Ручной ввод</option>
                                        <option value="sql">SQL</option>
                                        <option value="extendApi">Внешний API</option>
                                         </optgroup>
                                         <optgroup  label = "Автоматическое заполнение"> 
                                        <option value="osp">ОСП</option>
                                        <option value="spi">СПИ</option>                                      
                                         </optgroup>
                                    </select>
                             </div>`
        },
        table_constructor: function (tplValues) {
             console.log("table_constructor tplValues", tplValues);
            let dataSource;
             if(!!tplValues.data_source){
                 dataSource = tplValues.data_source.source;

             } else if(!!tplValues.dataSource){
                 dataSource = tplValues.dataSource;
             }
            console.log(dataSource);

            // tplValues.orders = tplValues.orders?tplValues.orders:0;
            // bundIputs(['unique_${tplValues.num}', 'index_${tplValues.num}'], this,false)
            return {
                'lastOrder': `${tplValues.orders}`,
                'order': `<input readonly="true" size="2"
                                    data-collumn-orders="${tplValues.orders}"
                                    data-collumn-alias="collumn_${tplValues.orders}"
                                    data-collumn-caption="№ П.П."
                                    data-target-table="${tplValues.alias}"
                                    class=""
                                    type="text"
                                    placeholder=""
                                    name="collumn_order.${tplValues.orders}"
                                    id="collumn_order_${tplValues.orders}"
                                    value="${tplValues.orders}">`,
                'caption': `<input required 
                                    oninput="bundIputs(['collumn_alias_${tplValues.orders}'], this,true)" 
                                    onfocus="checkPrevValue(event);bundIputsFocus(['collumn_alias_${tplValues.orders}'], this,true)"
                                    data-collumn-orders="${tplValues.orders}"   
                                    data-collumn-alias="collumn_${tplValues.orders}"   
                                    data-collumn-caption="Название столбца (рус.)"   
                                    data-target-table="${tplValues.alias}"  
                                    class="" 
                                    type="text" 
                                    placeholder="" 
                                    name="collumn_name.${tplValues.orders}" 
                                    id="collumn_name_${tplValues.orders}" 
                                    value="${!!tplValues.caption ? tplValues.caption : ""}">`,
                'field': `<input required  
                                    oninput="translitSelf(this);checkEdit(event)"
                                    onfocus="checkPrevValue(event)"
                                    onchange="sql.setCollumnId(this.dataset.collumnAlias,this)"
                                    data-collumn-orders="${tplValues.orders}"
                                    data-collumn-caption="Алиас столбца (лат.)" 
                                    data-collumn-alias="collumn_${tplValues.orders}" 
                                    data-target-table="${tplValues.alias}"  
                                    title="Под этим именем создаёться столбец в БД, толькона латинском без пробелови др. спецсимволов.Латинские буквы, и цифры.\nПервый символ НЕ ДОЛЖЕН быть цифрой" 
                                    class="" 
                                    type="text" 
                                    pattern="\\b(\\D).*"
                                    placeholder="Латинские буквы, и цифры. Первый символ не должн быть цифрой" 
                                    name="collumn_alias.${tplValues.orders}" 
                                    id="collumn_alias_${tplValues.orders}" 
                                    value="${!!tplValues.field ? tplValues.field : ""}">`,
                'comment': `<input 
                                      onfocus="checkPrevValue(event)"
                                      oninput="checkEdit(event)"
                                      data-collumn-alias="collumn_${tplValues.orders}"
                                      data-collumn-caption="Описание столбца"   
                                      data-target-table="${tplValues.alias}"  
                                      data-collumn-orders="${tplValues.orders}"
                                      title="Описание столбца" 
                                      class="" 
                                      type="text" 
                                      placeholder="" 
                                      name="collumn_comment.${tplValues.orders}" 
                                      id="collumn_comment_${tplValues.orders}" 
                                      value="${!!tplValues.comment ? tplValues.comment : ""}">`,
                'sourceIde':`<div id ="${dataSource}_container_${tplValues.orders}">
                            <button type="button" 
                                    class="btn" 
                                    data-source-container="${dataSource}_container_${tplValues.orders}"
                                    onclick="sql.editDataSourceHandler({'data_type':'${dataSource}','eventTarget':this,'orders':'${tplValues.orders}','wrapper':'modal_handler_${tplValues.orders}'})">
                                    <i class="fas fa-edit"></i>
                                    </button></div>`,
                'source': `<div class="row nowrap">
                           <select required name="source.${tplValues.orders}" 
                                    onfocus="checkPrevValue(event)"
                                    onchange="checkEdit(event);checkSelected(event);"
                                    id="source_${tplValues.orders}"
                                     data-collumn-caption="Источник данных" 
                                     data-collumn-orders="${tplValues.orders}">     
                                        <optgroup selected label = "Источник данных">
                                            <option disabled hidden ${!dataSource?'selected="selected"':""} value="">Источник данных</option>
                                            <option  ${(!!dataSource && dataSource==="manual")?'selected="selected"':""} value="manual">Ручной ввод</option>
                                            <option ${(!!dataSource && dataSource==="sql")?'selected="selected"':""} data-handler = "sql.setDataSourceHandler" value="sql">SQL</option>
                                            <option ${(!!dataSource && dataSource==="api")?'selected="selected"':""}  data-handler = "sql.setDataSourceHandler" value="api">Внешний API</option>
                                        </optgroup>
                                        <optgroup  label = "Автоматическое заполнение">
                                            <option ${(!!dataSource && dataSource==="osp")?'selected="selected"':""}  value="osp">ОСП</option>
                                            <option ${(!!dataSource && dataSource==="spi")?'selected="selected"':""}  value="spi">СПИ</option>                                      
                                         </optgroup>
                                    </select>                                       
                                       ${(!!dataSource && (dataSource==="sql" || dataSource==="api"))?`<div id ="${dataSource}_container_${tplValues.orders}"><button type="button" 
                                            class="btn" 
                                             data-source-container="${dataSource}_container_${tplValues.orders}"
                                            onclick="sql.editDataSourceHandler({'data_type':'${dataSource}','eventTarget':this,'orders':'${tplValues.orders}','wrapper':'modal_handler_${tplValues.orders}'})">
                                            <i class="fas fa-edit"></i>
                                            </button></div>`:""} 
                                    </div>`,
                'action': `<button type="button" data-collumn-orders="${tplValues.orders}"
                                                data-target-alias="${tplValues.alias}" 
                                                data-collumn-alias="collumn_${tplValues.orders}"     
                                                class="btn" 
                                                onclick="sql.removeCollumnSetting(this.dataset.targetAlias, null, 0, ${tplValues.orders},  this)">                                
                                                    -
                              </button>`,
                'type': `<select name="collumn_type.${tplValues.orders}" 
                                onfocus="checkPrevValue(event)"  
                                onchange="checkEdit(event)"  
                                data-collumn-orders="${tplValues.orders}"
                                data-collumn-caption="Тип данных стобца"
                                id="dataType_${tplValues.orders}" class="">
                                                <option ${tplValues.type === "integer" ? 'selected="selected"' : null} value="integer">Целое число</option>
                                                <option ${tplValues.type === "numeric" ? 'selected="selected"' : null} value="numeric">Дробное число</option>
                                                <option ${tplValues.type === "varchar" ? 'selected="selected"' : null}  value="varchar">Строка/Текст</option>
                                                <option ${tplValues.type === "jsonb" ? 'selected="selected"' : null}  value="jsonb">Структура/Массив/JSON</option>
                                                <option ${tplValues.type === "timestamp" ? 'selected="selected"' : null}  value="timestamp">Дата-время (timestamp)</option>
                                            </select>`,
                'unique': `<input type="checkbox" 
                                    onchange="checkPrevValue(event);checkEdit(event)"  
                                    data-collumn-orders="${tplValues.orders}"
                                    data-collumn-caption="Уникальный стоблбец"  
                                                ${!!tplValues.setting ? tplValues.setting.unique ? 'checked' : null : null}
                                    value="true" 
                                    name="unique.${tplValues.orders}" 
                                    id="unique_${tplValues.orders}">`,
                'index': ` <input type="checkbox" ${!!tplValues.setting ? tplValues.setting.index ? 'checked' : null : null}
                                  onchange="checkPrevValue(event);checkEdit(event)"  
                                  data-collumn-orders="${tplValues.orders}"
                                  data-collumn-caption="Индексируемый столбец"
                                  value="true"  
                                   name="index.${tplValues.orders}" 
                                   id="index_${tplValues.orders}">`,
                'draftSaved': `<label for="draft_saved_${tplValues.orders}">
                                <i class="status" id="saved_status_${tplValues.orders}" title="Сохраненно или нет">
                                ${!!tplValues.last_change?"+":"-"}
                                </i>
                                <div id="last_change_${tplValues.orders}" title="Дата последнего сохранения">
                                 ${!!tplValues.last_change?tplValues.last_change:"Не сохраннено"}
                                </div>
                                </label>
                               <input data-collumn-orders="${tplValues.orders}"
                                    ${!!tplValues.last_change?'checked':null}
                                     readonly="true" 
                                     type="checkbox" 
                                     value="true"  
                                     name="draft_saved.${tplValues.orders}" 
                                     id="draft_saved_${tplValues.orders}">
                                <input data-collumn-orders="${tplValues.orders}"
                                        readonly="true" 
                                        type="checkbox"
                                        value=""
                                        name="editing.${tplValues.orders}"
                                        id="editing_${tplValues.orders}">`
            }
        },
    };
    result = tplMap[tplName](tplValues);
    return result;
};

var ioTable = {
    Tables:{},  // maxOrder - последнее значение номера строки     index - индексы таблицы,

    Container:"", // контейнер для множества таблиц - родитеский контейнер куда при необходимости будут вставлены несколько таблиц если не найен на странице tContaner<Allias>
    tContainer:{}, // контенеры для конкретных таблиц, ключ = Allias
    HeadMap:{}, // tableAliias :{collName:order,index:true/false}
    dataFromServer:{},
    // tableAlias - ID куда будет вставлена таблица
    // tableAlias должен соответствовать action_code - действия данные которого будут выбираться
    // или  вфефЫуе может содержать action_code(string), тогда tableAlias может быть любой другой
    // header - Для имён колонок нужно сделать в БД таблицу в которой будут имена таблиц
    // Либо в таблицу action
    uopdateCell:function(cellData){
        console.log(cellData);
    },
    /*
    * this.index - устанавливает или извлекает индексы таблицы
    * indexData -  объект {tableAlias:"альяс таблицы","fieldName":{"cellValue":"rowIdx"},}  где fieldName - имя столбюца по которому индексируеться, cellValue значение ячейки, rowIdx номер строки которому соотвествует значение
    * indexData.get={"field":[value, value...]}
    * indexData.set={"field":[{rowNum:0, cellValue:value}, {rowNum:0, cellValue:value}...]}
    * */
    index:function(indexDate){
         console.log(indexDate);

        if(!!indexDate.tableAlias &&  !!this.Tables[tableAlias]){
            let rowsNumber=[];
            if (!!indexDate.get){
                for (let field in indexDate.get){
                    if (!!field && !!this.Tables[tableAlias].index[field]){
                        // если индекс существует то выберем все строки с указаным значение
                        for (let idx in indexDate.get[field]){
                            rowsNumber[idx] = this.Tables[indexDate.tableAlias].index[field][field[idx]]
                        }
                    }
                }
            }
            if (!!indexDate.set){
                for (let field in indexDate.set){
                    if(!!field){
                        if (!this.Tables[tableAlias].index[field]){
                            this.Tables[tableAlias].index[field] = {};
                        }
                        for (let idx in indexDate.set[field]) {
                            this.Tables[tableAlias].index[field][field[idx].value]=field[idx].rowNum
                            rowsNumber[idx] = this.Tables[tableAlias].index[field][field[idx].value];
                        }

                    }

                }
            }
        return rowsNumber
        }


    },


    updateRow: function(tableAlias, data){
         console.log(this.Tables[tableAlias]);
         console.log(this.Tables[tableAlias].rowTpl);
           if (!this.Tables[tableAlias]){
               this.Tables[tableAlias] = document.getElementById(tableAlias);
           }
           const tplName = this.Tables[tableAlias].rowTpl;
           for (let row in data.rows){
               let rowData = data.rows[row];
                if (!!tplName){
                   let newRow = tpl(tplName, rowData);
                   if (!!this.Tables[tableAlias].rows["R"+rowData.order]){
                       const tHead = this.Tables[tableAlias].tHead;
                       let updatedRow = this.Tables[tableAlias].rows["R"+rowData.order];
                        console.log(updatedRow);
                       for(const headRow of tHead.rows){

                           for(let headCell of headRow.cells){

                               let colAlias = headCell.id;

                               if(!!newRow[colAlias]) {
                                   const cellIndex = headCell.cellIndex;
                                   updatedRow.cells[cellIndex].innerHTML = newRow[colAlias];

                               }
                           }
                       }
                   }

                }
           }
    },
    updateBody: function (tableAlias){
        let data = this.dataFromServer[tableAlias];
        let tableRows;
        if (!data.ActionData.table || data.ActionData.table[tableAlias].rows.lenght<1){
            messages({message:"Данных для отображения таблицы "+JSON.stringify(tableAlias) + " нет",type:"info"});
            tableRows = {data:"Нет данных для отображения"};

        } else {
            tableRows = data.ActionData.table[tableAlias].rows;
        }

        console.log(tableRows);
        console.log(tableAlias);
        console.log(this.Tables[tableAlias]);
        console.log(this.Tables[tableAlias].tBodies);
        this.Tables[tableAlias].tBodies[0].replaceWith(this.Body({"tableAlias":tableAlias,
            "tableRows": tableRows,
        }));

    },
    load:async function (tableAlias){
        messages({message:"Загружаем данных для отображения таблицы "+JSON.stringify(tableAlias),type:"info"});
        console.log(tableAlias);
        // var request = new Request('/?jsonql={"get":{"data":["'+tableAlias+'"]}}',
        let block ={};
            block["table"]= {// id/имя блока
            "block_name":"table",
            "format": "html",
            "params":{
                "alias":tableAlias
            },
        };
        let jsonql = {
            "get":{
                "blocks": block,
                "pages": null,
                "layout": null,
                "format": "html"
            }
        };
        console.log(window.location);
        var request = new Request("?jsonql="+JSON.stringify(jsonql),
            {
                credentials: "include",
                mode: 'cors',
                method: 'GET'});

         console.log(request);
        request.headers.append("Request-method","ajax/fetch");
        // let tableScoop = this;
        // Синхронный запрос
        const response = await fetch(request);

        this.dataFromServer[tableAlias] = await response.json();

        messages({message:"Данные получены "+JSON.stringify(tableAlias),type:"info"});

        this.updateBody(tableAlias)
    },
    /*
    * tablesContainer - div  в который помещаються таблицы, не обязательный параметр
    * tableAlias -  alias и id таблицы, к div контейнеру добавляеться приставка tContainer<+Allias>
    * */
    // init: function (tablesContainer, tableAlias, headers, footers, tableRows = null){
    // index = {aliasField : {значение поля: номер_строки}}

    init: function (tablesData){
        console.log(tablesData.data.tables)

        for (let TableName in tablesData.data.tables) {

            let TableStruct=tablesData.data.tables[TableName];
                console.log("init tableData",TableStruct);



            let caption, tableAlias, footers, headers, tablesContainer, tableRows, tableIndex; //TableStruct,

            // if (!!TableStruct && tableData.hasOwnProperty("data")){
            if (!!TableStruct){
                // TableStruct = tableData.data;

                // console.log("TableStruct", TableStruct.table.alias);

                    if (!!TableStruct.alias && TableStruct.alias !== ""){
                        tableAlias = TableStruct.alias;
                    } else if (!!TableStruct.alias && TableStruct.alias !== ""){
                        tableAlias = TableStruct.alias;
                    }

                // if (!!TableStruct.table){
                // if (!!TableStruct.table && Object.keys(TableStruct.table).length>0){
                    if(!!TableStruct.footer && Object.keys(TableStruct.footer).length>0){
                        footers = TableStruct.footer
                    }
                    if(!!TableStruct.header && Object.keys(TableStruct.header).length>0){
                        headers = TableStruct.header
                    }
                    if(!!TableStruct.rows && Object.keys(TableStruct.rows).length>0){
                        tableRows = TableStruct.rows
                    }
                    if(!!TableStruct.index && Object.keys(TableStruct.index).length>0){
                        tableIndex = TableStruct.index
                    }
                    if(!!TableStruct.caption && Object.keys(TableStruct.caption).length>0){
                        caption = TableStruct.caption
                    }
                // }
            } else if (!TableStruct || !TableStruct.rows) {
                 console.log("нет данных для рендера таблицы выходим");
                 return;
            }
            // if(!!tableData.page_name){
            //     tablesContainer = "tContainer_"+tableData.page_name;
            // } else {
            //     tablesContainer =  tableAlias;
            // }

            // если в tableRows пришёл null или строка то получаем данные с сервера по tableAlias
            // console.log(tablesContainer, tableAlias, headers, footers, tableRows );

            this.checkTableDOM(tablesContainer, tableAlias); // проверим есть ли на страницу разметка для таблицы

            let Head;
            let rowsTpl;

            if (!!caption){
                let captionElem = document.createElement("caption");
                    captionElem.innerText=caption;
                this.Tables[tableAlias].appendChild(captionElem);
            }
            if (!!headers){
                Head = this.Header(tableAlias,headers,tableIndex);
                this.Tables[tableAlias].appendChild(Head);
            }
            let Footer;
            if(!!footers){
                 console.log("Создаём футер таблицы "+tableAlias);
                Footer = this.Footer(tableAlias,footers);
                // this.Tables[tableAlias].appendChild(Footer);
            }
            // console.log(this.tContainer[tableAlias]);
            // console.log(this.Container);



                if (!tableRows){
                    if (!!TableStruct.setting && !!TableStruct.setting.published){
                        tableRows={"data":"Загружаем данные"};
                        this.load(tableAlias);
                    } else if (!TableStruct.setting){
                        tableRows={}
                    }

                } else  if(!!tableRows && (tableRows.length<1 ||  (typeof tableRows === "string"))) {
                    // this.load(tableAlias);
                }

                let BodyData = {"tableAlias":tableAlias,
                                "tableRows":tableRows,
                                "tableIndex":tableIndex };

            let Body = this.Body(BodyData);

            // this.Tables[tableAlias].appendChild(Body);

            // if (!!this.tContainer[tableAlias]) {
            //     this.tContainer[tableAlias].appendChild(this.Tables[tableAlias]);
            // } else if(!!this.Container){
            //     this.Container.appendChild(this.Tables[tableAlias]);
            // }

            // dataFromServer.then(function (tableRows){
            //     console.log(tableRows);
            // });



        }
        return this;
    },
    getCountRows:function(tableAlias){
         console.log(tableAlias);
      return  document.getElementById(tableAlias).tBodies[0].rows.length
    },
    /**
     * setMaxOrder - устанавливает для таблицы из tableAlias значение максимального числа для поля order колонки
     * @param tableAlias - id таблицы
     * @param order - текущий order
     */
    setMaxOrder:function(tableAlias, order=1){

        if (!this.Tables[tableAlias].maxOrder){
            this.Tables[tableAlias].maxOrder =order
        } else {
            if (this.Tables[tableAlias].maxOrder<order){
                this.Tables[tableAlias].maxOrder = order
            } else if (this.Tables[tableAlias].maxOrder=order){

                console.log("maxOrder=tableRows[Ridx].order");
            }
        }
    },
    /**
     * getNextOrder - возвращает следующее значние order для новой колонки для таблицы из tableAlias
     * @param tableAlias  id таблицы
     */
    getNextOrder:function(tableAlias){
         console.log("getNextOrder", tableAlias);

        if (!!this.Tables[tableAlias] && !!this.Tables[tableAlias].maxOrder){

            console.log("560 getNextOrder", this.Tables[tableAlias].maxOrder);

            this.Tables[tableAlias].maxOrder = this.Tables[tableAlias].maxOrder+1;

            console.log("564 getNextOrder", this.Tables[tableAlias].maxOrder);

            return this.Tables[tableAlias].maxOrder

        } else if (!!this.Tables[tableAlias] && !this.Tables[tableAlias].maxOrder){

            console.log("568 getNextOrder", this.Tables[tableAlias].maxOrder);

            return this.Tables[tableAlias].maxOrder = 1
        } else if (!this.Tables[tableAlias]){
             console.log("Таблица не найдена: ",tableAlias );
        }
    },


    Body:function(BodyData){
     //BodyData = {"tableAlias","tableRows", "tBodyId"=null, "tableIndex"=null }
            console.log("Создаём body таблицы ", BodyData.tableAlias);
            console.log("BodyData", BodyData);
            console.log("this.Tables", this.Tables);
            let tBody =document.createElement("tbody"); //Создадим элемет tbody
            if(!BodyData.tBodyId){
                tBody.id = "tbody_"+BodyData.tableAlias;
                BodyData.tBodyId = tBody.id
            } else {
                tBody.id = BodyData.tBodyId;
            }
            console.log(BodyData.tBodyId);

        this.Tables[BodyData.tableAlias].appendChild(tBody);

        this.Tables[BodyData.tableAlias].index=BodyData.tableIndex;

        // console.log("Body tableRows",BodyData.tableRows);
        // console.log("Body tObject.keys(BodyData.tableRows).length",Object.keys(BodyData.tableRows).length);
        // console.log("Body tObject.keys(BodyData.tableRows)",Object.keys(BodyData.tableRows));

        if (!BodyData.tableRows || BodyData.tableRows.length<1 || (!!BodyData.tableRows && !Array.isArray(BodyData.tableRows) && Object.keys(BodyData.tableRows).length<1)){
            // messages({message:"Нет данных для отображения таблицы "+JSON.stringify(BodyData.tableRows),type:"warning"})
            // sendLogError({"type":"warning", "message":"Нет данных для отображения таблицы"});
            // let colspan = Object.keys(this.HeadMap[tableAlias]).length;
            // let tr =document.createElement("tr");
            //     tr.setAttribute("id","0");
            // let td = document.createElement("td");
            //     td.align ="center";
            //     td.colSpan=colspan;
            //     td.innerHTML = "<div id='loader'"+tableAlias+">Загружаем данные</div>";
            //     tr.appendChild(td);
            //     tBody.appendChild(tr);
             console.log("недобавляем строки");
        }else if (!!BodyData.tableRows && !Array.isArray(BodyData.tableRows)){
            //
            // console.log("Body BodyData.tableRows",BodyData.tableRows);
            // console.log("Body tObject.keys(BodyData.tableRows).length",Object.keys(BodyData.tableRows).length);
            // console.log("Body tObject.keys(BodyData.tableRows)",Object.keys(BodyData.tableRows));

/*Строка с обединёнными ячейками для вывода сообщения в таблице*/
            if(Object.keys(BodyData.tableRows).length > 0){

                 // console.log(this.HeadMap[BodyData.tableAlias]);
                let colspan = Object.keys(this.HeadMap[BodyData.tableAlias]).length;
                let tr =document.createElement("tr");
                    tr.setAttribute("id","0");
                let td = document.createElement("td");
                    td.align ="center";
                    td.colSpan=colspan;
                    td.innerHTML = "<div id='loader'"+BodyData.tableAlias+">"+BodyData.tableRows.data+"</div>";
                    tr.appendChild(td);
                    tBody.appendChild(tr);
           }
        } else {
            console.log(BodyData.tableRows)
            // for (let Ridx in BodyData.tableRows) {

                    // console.log(this.Tables[BodyData.tableAlias].rowTpl);
                    if(!!this.Tables[BodyData.tableAlias].rowTpl && this.Tables[BodyData.tableAlias].rowTpl !== ""){
                        // ТУТ должен быть код для конструктора табиц
                        // this.setMaxOrder(BodyData.tableAlias, BodyData.tableRows[Ridx].order);
                        // const rowsData = tpl("table_constructor", BodyData.tableRows[Ridx]); // this.Tables[tableAlias].rowTpl
                        // this.addRows(BodyData.tableAlias, [rowsData], BodyData.tBodyId, 0, BodyData.tableIndex)
                    } else {
                        this.addRows(BodyData.tableAlias, BodyData.tableRows, BodyData.tBodyId, 0, BodyData.tableIndex)
                    }
            // }
        }

        return tBody;
    },
    PrepareData:function(value, data){
        /*Подготавливает данные для функции addRows*/
        let result ={};
        start =  Date.now();


        for (let cell of data){
            if (cell.caption.includes('${')){
                cellValue = new Function ("value","return `"+cell.caption+"`");
                cellValue.apply(null, [value]);
                result[cell.field] = cellValue(value);
            } else {
                result[cell.field] = cellValue(cell.caption);
            }
        }
        stop =  Date.now();
         console.log(stop-start);
         console.log(result);
        return [result];
    },
    Header:function(tableAlias, headers,tableIndex=null){

        console.log("Создаём head таблицы ",tableAlias);
         console.log("tableIndex ",tableIndex);

        let tHead = document.createElement("thead");
        this.HeadMap[tableAlias]={};

        console.log(headers);

        // if (headers[0].alias === "table_constructor"){
        //     rowsTpl = "table_constructor";
        // }
        if(!!headers && headers.length >0){

            let td="";
            let tr = document.createElement("tr");
            for (let orderStr in headers){

               let order = Number.parseInt(orderStr);
                // console.log(Number.isInteger(order));

                if (!Number.isInteger(order)){
                    continue;
                }
                let header = headers[order];

                // console.log(header);

                if ((!this.Tables[tableAlias].rowTpl || this.Tables[tableAlias].rowTpl ==="") && header.alias !== tableAlias){
                    this.Tables[tableAlias].rowTpl = header.alias;
                }

                let width;
                if (!header.width){
                    width = '';
                } else {
                    width = 'width="'+header.width+'"';
                }


                if (!!header.field ||  header.caption) {
                    console.log("Переименовал field в col_name, caption в label")
                }

                console.log(header);

                let td = document.createElement("th");
                    td.setAttribute("id", header.col_name);
                    td.setAttribute("width", width);
                    td.dataset.field = header.col_name;
                    td.innerHTML = header.label;

                 // console.log(header);
                if (!!tableIndex && !!tableIndex[header.col_name]){
                    td.dataset.index="true"
                }
                tr.appendChild(td);
                // let col_parameters=JSON.parse(header.col_parameters)
                this.HeadMap[tableAlias][header.col_name]=order;

            }
            tHead.appendChild(tr);
            console.log(this.HeadMap);
            // tHead = `<thead><tr>${td}</tr></thead>`;
        }
        return tHead

    },
    Footer:function(tableAlias,footers){

        let tFoot=document.createElement("tfoot");
            tFoot.id = "tfoot_"+tableAlias;
        this.Tables[tableAlias].appendChild(tFoot);
         // console.log("Footer", this.Tables[tableAlias].tFoot);

        // if(!!footers && footers.length >0){
        //
        //
        //     let tr = document.createElement("tr");
        //     for (let footer of footers){
        //         let td = document.createElement("td");
        //         td.setAttribute("id", header.field);
        //         td.dataset.field = header.field;
        //         td.innerHTML = header.caption;
        //         tr.appendChild(td);
        //     }
        //     tFoot.appendChild(tr);
        //
        // }
 console.log( this.Tables[tableAlias]);

        for (let Ridx in footers){
            // if(!!this.Tables[tableAlias].rowTpl && this.Tables[tableAlias].rowTpl !== ""){
            //     this.addRows(tableAlias, [this.tpl(this.Tables[tableAlias].rowTpl, footers[Ridx])], tFoot.id, 0)
            // } else {
                this.addRows(tableAlias, this.PrepareData({"tableAlias":tableAlias},footers), tFoot.id, 0)
            // }
        }

        return tFoot;
    },
    checkTableDOM: function (tablesWrapper, tableAlias){

        // console.log("tableAlias", tableAlias);
        // console.log("this.Tables[tableAlias] ",this.Tables[tableAlias]);

        // if (!!this.Tables[tableAlias]){ // проверим есть ли инициированная таблица в хранилище таблиц
        //     // если есть то просто возвращаемся и продолжаем обработку данных
        //     return;
        // }
        // если таблица новая то проверим существует ли для неё контенер или сама таблица на страницу

        // if (!document.getElementById("tContainer_"+tableAlias)){ // контенера для текущей таблицы нету
        if (!document.getElementById("tContainer_"+tableAlias)){ // Разметки для таблицы нету

            // this.Container = document.getElementById(tablesWrapper); // Проверим есть ли разметка для таблиц на текущей странице
            // messages({message:"Не найдена разметка для таблицы "+tableAlias+" на странице  "+tablesWrapper,type:"warning"});
            sendLogError({"type":"warning", "message":"Не найден контейнер таблицы"});

            // if (!this.Container){ // проверим есть ли контенер для множества таблиц
            //     messages({message:"Не найден контейнер для таблиц на странице  "+tablesWrapper,type:"warning"});
            //     sendLogError({"type":"warning", "message":"Не найден контейнер таблицы"});
            // } else {
            //     let tCont = document.createElement("div");
            //     tCont.id = "table_wrapper_"+tableAlias;
            //     this.Container.appendChild(tCont);
            //     this.tContainer[tableAlias] = tCont;
            // }
        } else {// контенер для текущей таблицы найден

            this.Tables[tableAlias] = document.getElementById("tContainer_"+tableAlias);

            // let tablesWrapperOnPage = document.getElementById(tablesWrapper);
            // this.Container =  tablesWrapperOnPage ; // Если на странице есть разметка для вставки таблиц то  тут этот элемнт не обязателен
            // this.tContainer[tableAlias]=document.getElementById("table_wrapper_"+tableAlias);
        }
        // console.log("this.Tables[tableAlias] ",this.Tables[tableAlias]);
        // this.tablesWrapper = document.getElementById(tablesWrapper);
        //
        //
        // if (!this.tablesWrapper){
        //     messages({message:"Не найден контейнер для таблиц на странице  "+tablesWrapper,type:"warning"});
        //     sendLogError({"type":"warning", "message":"Не найден контейнер таблицы"});
        // }
    },
    deleteCollumn: function(tableAlias, columnId){
        let table = this.Tables[tableAlias];
        if (!table){
            table = document.getElementById(tableAlias);
            this.Tables[tableAlias] = table;
        }
        table.deleteColl(null, columnId)

    },
    /*
    * cells - массив объектов, с данными для чейки добавляемого столбца, индекс ячейки соответсвует строке
    * */
    addCollumn: function (tableAlias, target, header, body, footer=null){
        let table = this.Tables[tableAlias];
        if (!table){
            table = document.getElementById(tableAlias);
            this.Tables[tableAlias] = table;
        }


        let newCell = table.insertColls(-1, header, body, footer);

    },

    /*
    * addRow добавляем строку с данными в tElement = элемент таблицы в который нужно вставить строку, table.tbodies[0]
    * table.tFoot
    * */
 //    addRow:function(tableAlias, rowData, tElement){
 //        /**/
 //         console.log(tElement);
 //        const table = document.getElementById(tableAlias);
 //
 //        if (!!tElement){
 //            messages({"message":"Не задан элемент таблицы "+tableAlias+" в который вставлять строку" , "type":"error"});
 //        }
 //        let newRow = tElement.insertRow(-1);
 //
 //        // const rowId = rowData.order?rowData.order:rowData.id?rowData.id:tElement.rows.length;
 //
 // console.log("addRow rowData",rowData);
 //        if(tElement.tagName === "tfoot"){
 //            newRow.id = rowData.lastOrder?"F"+rowData.lastOrder+1:"F"+tElement.rows.length
 //        } else if (tElement.tagName === "tbody"){
 //            newRow.id = rowData.lastOrder?"R"+rowData.lastOrder+1:"R"+tElement.rows.length
 //        }
 //        console.log("addRow newRow",newRow);
 //
 //
 //        const tHead = table.tHead;
 //        for(const headRow of tHead.rows) {
 //            for (let headCell of headRow.cells) {
 //                let colAlias = headCell.id;
 //                if(!!rowData[colAlias]){
 //                    const cellIndex = headCell.cellIndex;
 //                        newCell = newRow.insertCell(cellIndex);
 //                        newCell.innerHTML = rowData[colAlias];
 //                        continue;
 //                }else {
 //                    messages({"message":"Таблица "+tableAlias+" Не найдены данные для столбца с ID "+colAlias , "type":"info"});
 //                    // newCell = newRow.insertCell(-1);
 //                }
 //            }
 //        }
 //
 //    },

    /*
    * addRows добавляет 1 строку в таблицу
    * если tElementId не задан то считаем что добавляем строки в tBodies[0]
    * Иначе находим элемент с указаным id=tElementId  и него вставляем строку
    * TODO: Нужно реализовать : 1. Парсинг параметров каждой ячейки coolspan, key, hidden,
    *  TODO: Присваивать строке или data атрибут или id в в соответсвии с переданными параметрами
    * */
    // Table.addRows(this.dataset.targetAlias, [Table.tpl('table_constructor', {order: Table.getCountRows(this.dataset.targetAlias), alias:this.dataset.targetAlias})], null, null);sql.saveDraft(this.form);
    addRows:function(tableAlias, rowsData, tElementId=null, tBodyIdx=0,tableIndex=null){

         // console.log(tableAlias, rowsData, tElementId, tBodyIdx, tableIndex);

         if (!tableIndex){
             tableIndex = this.Tables[tableAlias].index
         }
        let table = this.Tables[tableAlias];
        // console.log(table)
        if(!table){
            table = document.getElementById(tableAlias);
        }

           const newRows = {};
           let tElement; // tElement - Элемент таблицы куда будут производиться встаавки новых строк tbody/tfoot/thead
           if(!tElementId){
                if(!!table.tBodies && table.tBodies.length >0){
                    tElement = table.tBodies[!!tBodyIdx?tBodyIdx:0]
                } else {
                    if(!table.tBodies){
                       return this.Body({"tableAlias":tableAlias,
                                                  "tableRows": rowsData,
                                                  "tBodyId":tElementId
                                                    })
                    }
                }
           } else {
               tElement = document.getElementById(tElementId);
               if (!tElement){
                   tElement = table.children[tElementId];
               }
           }


           let iota=0;
           for(let rowData of rowsData){
                let newRow = tElement.insertRow(-1);
               if(tElement.tagName.toLowerCase() === "tfoot"){
                   newRow.id = rowData.lastOrder?"F"+rowData.lastOrder:"F"+tElement.rows.length
               } else if (tElement.tagName.toLowerCase() === "tbody"){
                   newRow.id = rowData.lastOrder?"R"+rowData.lastOrder:"R"+tElement.rows.length
               }

               // console.log("addRow newRow",newRow.id);
               newRows[newRow.id]= newRow;

                const tHead = table.tHead;
                    // console.log("tHead.rows", tHead.rows);
                    for(const headRow of tHead.rows){
                        for(let headCell of headRow.cells){
                            let colAlias = headCell.id;
                            if (rowData[colAlias]===""){
                                rowData[colAlias]="0"
                            }
                            if (colAlias=="iota"){
                                iota++;
                                rowData[colAlias]=iota;
                            }

                            if(!!rowData[colAlias]){
                                const cellIndex = headCell.cellIndex;

                                let newCell;

                                    if (newRow.cells.length>=cellIndex){
                                        if (!!newRow.cells[cellIndex] && newRow.cells[cellIndex]===null) {
                                            newCell = newRow.cells[cellIndex];
                                        } else if(!newRow.cells[cellIndex]) {
                                            newCell = newRow.insertCell(cellIndex);
                                            //
                                            // if (colAlias === "order"){
                                            //     this.index({"tableAlias":tableAlias, "set":{"value":rowData[colAlias], rowNum:newRow.rowIndex}})
                                            // }

                                            if (!!tableIndex && !!tableIndex[colAlias]){
                                                newCell.dataset.index="true";
                                            }
                                            newCell.dataset.field=colAlias;

                                        }
                                    } else  if (newRow.cells.length<cellIndex){
                                        if (newRow.cells.length<=cellIndex-1){
                                            // /*diffCountCells сколько ячеек не хватает междудобавляемой и имеющимся кол.*/
                                            let diffCountCells = cellIndex - newRow.cells.length;
                                            for (let i=0;i<=diffCountCells;i++){
                                                newCell = newRow.insertCell(-1);
                                                 // console.log(newCell);
                                                // if (colAlias === "order"){
                                                //     this.index({"tableAlias":tableAlias, "set":{"value":rowData[colAlias], rowNum:newRow.rowIndex}})
                                                // }
                                                if (!!tableIndex && !!tableIndex[colAlias]){
                                                    newCell.dataset.index="true";
                                                }
                                                newCell.dataset.field=colAlias;
                                            }
                                        }
                                    }
                                     if (!!rowData["date_sync"]){
                                        newCell.setAttribute("title", "Дата расчёта показателя: "+rowData["date_sync"])
                                     }
                                        newCell.innerHTML = rowData[colAlias];
                                        continue;
                            } else {
                                    // messages({"message":"Таблица "+tableAlias+"Не найдены данные для столбца с ID "+colAlias , "type":"info"});
                                    // newCell = newRow.insertCell(-1);
                             }
                        }
                    }
           }

 // console.log("tElement", tElement);
 // console.log("newRows", newRows);
        return newRows;
    },
    removeRows:function(rowData){
        // rowData = {
        //     tableAlias,
        //     tBodyId=null,
        //     tBodyIdx=0,
        //     rowId
        // };
        const tableAlias = rowData.tableAlias;
        const tBodyId = rowData.tBodyId;
        const tBodyIdx = rowData.tBodyIdx;
        const rowId= rowData.rowId;

        const table = document.getElementById(tableAlias);
         console.log("tableAlias", tableAlias, "tBodyId",tBodyId, "tBodyIdx", tBodyIdx, "rowId",rowId);

        if (!!tBodyId){
             console.log("tBodyId",tBodyId);
             console.log("tBodyId rows",table.deleteRow(table.tBodies[tBodyId].rows));
             // console.log(table.tBodies[tBodyId].rows[rowId].rowIndex);
             if(!!table.deleteRow(table.tBodies[tBodyId].rows[rowId])){
                table.deleteRow(table.tBodies[tBodyId].rows[rowId].rowIndex);

                /*Проверим сколько строк в таблице соталось после удаления, и какой макимальный индекс записан,
                    если есть оставшиеся строки то проверим какое значение установленно максимальное и установим его в качестве максимального индекса строк
                */

                if (table.tBodies[tBodyId].rows.length < 1 ) {
                    Table[tableAlias].maxOrder = 0
                } else {
                    if (table.tBodies[tBodyId].rows.length > Table[tableAlias].maxOrder){
                         console.log(table.tBodies[tBodyId].rows);
                        for (let row of table.tBodies[tBodyId].rows){
                             console.log(row);
                        }
                    }
                }


             } else {
                 messages({"message":`СТрока с индексом ${tBodyId} не найдена`,"type":"error"}, true);
             }
        }

        if (!!tBodyIdx || tBodyIdx !== null){
            if(!!table.tBodies[tBodyIdx].rows[rowId]){

                table.deleteRow(table.tBodies[tBodyIdx].rows[rowId].rowIndex);

                if (table.tBodies[tBodyIdx].rows.length < 1 ) {
                    this.Tables[tableAlias].maxOrder = 0
                } else {
                    if (table.tBodies[tBodyIdx].rows.length > this.Tables[tableAlias].maxOrder){
                        console.log(table.tBodies[tBodyIdx].rows);
 console.log("Нужно написать функцию для пересчёта номера строки");
                        for (let row of table.tBodies[tBodyIdx].rows){
                            console.log(row);

                        }
                    } else {
                        for (let row of table.tBodies[tBodyIdx].rows){
                            console.log(row);
                        }
                    }
                }
            }else {
                messages({"message":`СТрока с ID ${rowId} не найдена`,"type":"error"}, true);

            }

        }
    },

};