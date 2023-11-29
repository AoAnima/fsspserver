в = {
    "jsonrpc": "2.0",
    "auth": "{{.}}",
    "id": 1,
    "method": "host.get",
    "params": {
    "output":"extend",
        "filter":{
        "host":["Discovered ARM OSP Linux"]
    }
}
}


let Master = {
    Назад: function (MasterId, e){
        let МастерКонтейнер = document.getElementById(MasterId);
        console.log(МастерКонтейнер)
        let ВсеШаги = МастерКонтейнер.querySelectorAll(".master_step")
        console.log(ВсеШаги)
        for (let Шаг of ВсеШаги){
            if (Шаг.classList.contains('active')){
                let НомерАктивногоШага = Шаг.dataset.master_step;
                let ПредидущийШаг = МастерКонтейнер.querySelector(`#step_${Number(НомерАктивногоШага)-1}`);
                if (!!ПредидущийШаг) {
                    ПредидущийШаг.classList.add('active')
                    ПредидущийШаг.classList.remove('hidden')
                    Шаг.classList.remove('active')
                    Шаг.classList.add('hidden')
                    // если ПредидущийШаг === 1 значит назад идти некуда выключим кнопку
                   if (ПредидущийШаг.dataset.master_step === "1"){
                       // e.target.classList.add('disabled')
                   }



                } else {
                    // e.target.classList.add('disabled')
                }

            }
        }
    },
    Далее: function(MasterId, e){
        let МастерКонтейнер = document.getElementById(MasterId);
        let ВсеШаги = МастерКонтейнер.querySelectorAll(".master_step")
        console.log(ВсеШаги)
        for (let Шаг of ВсеШаги){

            if (Шаг.classList.contains('active')){
                let НомерАктивногоШага = Шаг.dataset.master_step;
                console.log("НомерАктивногоШага", НомерАктивногоШага)
                let СледующийШаг = МастерКонтейнер.querySelector(`#step_${Number(НомерАктивногоШага)+1}`);
                console.log("Следующий шаг",`#step_${Number(НомерАктивногоШага)+1}`, СледующийШаг)
                if(!!СледующийШаг) {
                    СледующийШаг.classList.add('active')
                    СледующийШаг.classList.remove('hidden')
                    Шаг.classList.remove('active')
                    Шаг.classList.add('hidden')
                    // проверим есть ли ещё активные шаги после текущего, есл нет то выключим кнопку
                    if (!МастерКонтейнер.querySelector(`#step_${Number(НомерАктивногоШага)+2}`)){
                        // e.target.classList.add('disabled')
                    }
                    // останавливаем проход по вкладкам
                     break
                    // проверим и если нужэно включим кнопку назад
                } else {
                   let КонтейнерПоследнегоШага = document.getElementById(`step_${НомерАктивногоШага}`)
 console.log("КонтейнерПоследнегоШага", КонтейнерПоследнегоШага);
                    let форма = КонтейнерПоследнегоШага.querySelector('[data-finalform="true"]')
                     console.log(форма);
                    форма.dispatchEvent(new Event('submit', {cancelable: true}))

                }
            }
        }
    },
}



function SettingCustomSQL(event, ИДТаблицы){
    console.log(event.target)
    let Значение = event.target.value
    let Тип = event.target.dataset.sql_type
     console.log(Тип, "Значение", Значение);
    let НазваниеФильтра = event.target.dataset.caption
    // если спец фильтр был отмечен то добавим его в следующий шаг, и запросм его данные получим списко столбцов которые можно ыгрузить в детализацию
    // let КонтейнерНастроек = document.getElementById(ИДТаблицы)
    let Таблица = document.getElementById(ИДТаблицы)

    if (event.target.checked) {

        if (!Таблица.querySelector(`#${Тип}_wrapper_${Значение}`)){
            console.log(Таблица);

        let БлокСпецФильтра = `<th class="${Тип}_wrapper" id = "${Тип}_head_cell_${Значение}" class="">
            <div id ="${Тип}_name_${Значение}">${НазваниеФильтра}</div>
<!--            <input class="" readonly type="text" name="запрос[${Значение}].тип" value="${Тип}">-->
            <input class="" readonly type="text" name="запрос[${Значение}].имя_столбца" value="${НазваниеФильтра}">
            
            <label class="${Тип}_details" for="fields-${Тип}_id_${Значение}">
<!--   имя поля делаем как тип, т.к. инчае получается лишний инпут         -->
                     <input type="checkbox"
                             data-caption="${НазваниеФильтра}"
                             id="fields-${Тип}_id_${Значение}"
                             value="${Тип}"
                             name="запрос[${Значение}].тип"
                             form="new_data-table"
                             onchange="Детализация(event, '${Тип}_fields_', '${Тип}', ${Значение})">
                     <div class="input_name">
                        Нужна детализация по данным
                     </div>
               </label>  
               <div class="${Тип}_fields_wrapper" id="${Тип}_fields_${Значение}"></div>           
        </th>`
             console.log(Таблица.tHead);
            Таблица.tHead.rows[0].insertAdjacentHTML('beforeend',БлокСпецФильтра)
            // КонтейнерНастроек.insertAdjacentHTML('beforeend',БлокСпецФильтра)
            let ЯчейкаДанных =   `<td id="${Тип}_cell_${Значение}">0</td>`
            Таблица.tBodies[0].rows[0].insertAdjacentHTML('beforeend',ЯчейкаДанных)

        } else {
            WS.FloatMessage({MessageType:['error'],Контэнт:{данные:'Спецфильтр уже добавлен, в новую стат. таблицу'}})
        }

        // WS.ОтправитьСообщениеИО('Получить поля спей.фильтра', {"КонтейнерРезультата":ИДКонтейнераНастроек, "ИДСпецФильтра": Значение}, true)

    } else {
      let БлокЗаголовка  =  Таблица.querySelector(`#${Тип}_head_cell_${Значение}`)

         console.log(БлокЗаголовка);
         // console.log(КонтейнерНастроек);
        БлокЗаголовка.remove()
        let ЯчейкаДанных =  Таблица.querySelector(`#${Тип}_cell_${Значение}`)
        ЯчейкаДанных.remove()
    }

}

function Детализация (event, контейнер, тип , SqlId){

    if (event.target.checked){
        if (тип === "sql"){
            WS.ОтправитьСообщениеИО('получить поля sql запроса',{[тип+'_id']:`${SqlId}`,'контэйнер_результата':'modal.'+контейнер+SqlId, 'тип':тип})
        } else {
            WS.ОтправитьСообщениеИО('получить поля спец.фильтра',{[тип+'_id']:`${SqlId}`,'контэйнер_результата':'modal.'+контейнер+SqlId, 'тип':тип})
        }


    } else {
        console.log(контейнер+SqlId)
        document.getElementById(контейнер+SqlId).innerHTML=""
    }
}


let ШаблоныПолей = {
    argument: function (метка) {
        return `<div class="flex-row argument" id="arg_wrapper-${метка}">
            <label for="arg_name-${метка}">
                <div class="input_name">Название</div>
                <input type="text" name="аргументы[${метка}].имя_аргумента" id="arg_name-${метка}" placeholder="">
            </label>
            <label for="arg_id-${метка}">
                <div class="input_name">идентификатор</div>
                <input type="text" name="аргументы[${метка}].идентификатор" id="arg_id-${метка}" placeholder="">
            </label>
            <label for="arg_type-${метка}">
                <div class="input_name">тип</div>
                <select name="аргументы[${метка}].тип" id="arg_type-${метка}">
                     <option selected value="" disabled>Выбрать Тип</option>
                    <option value="date">Дата</option>
                    <option value="number">Число</option>
                    <option value="text" selected>Текст</option>
                </select>
            </label>
<!--            <button data-tpl="argument_tpl" data-target="arguments_wrapper" data-number="${метка}" onclick="ДобавитьАргументы(event)" class="plus" type="button">+</button>-->
            <button data-target="arg_wrapper-${метка}" data-class_name="argument" data-number="${метка}" class="minus" onclick="УдалитьАргумент(event)" type="button">-</button>
         </div>`
    },
    detail_field: function(метка){
       return  `<div class="flex-row detail_field" id="detail_field-${метка}" >
                                                <label for="field_name-${метка}">
                                                    <div class="input_name auto">Название</div>
                                                    <input value="" type="text" name="детализация[${метка}].имя_поля" id="field_name-${метка}" placeholder="">
                                                </label>
                                                <label for="field_id-${метка}">
                                                    <div class="input_name auto">SQL идентификатор</div>
                                                    <input type="text" name="детализация[${метка}].идентификатор" id="field_id-${метка}" placeholder="" value="">
                                                </label>
                                                <button data-target="detail_field-${метка}" data-class_name="detail_field"  data-number="${метка}" class="minus" onclick="УдалитьАргумент(event)" type="button">-</button>
                                            </div>`
    }

}


function ДобавитьАргументы (event, ИмяШаблона){
    let Кнопка = event.target
    let метка = Кнопка.dataset.count
    // let МассивУдалённыхАргументов = Кнопка.dataset.deleted
    // if (МассивУдалённыхАргументов !== "[]") {
    //      console.log(МассивУдалённых);
    //     // метка = МассивУдалённыхАргументов
    //     МассивУдалённых=JSON.parse(МассивУдалённыхАргументов)
    //     метка = МассивУдалённых.shift()
    //     // var myIndex = myArray.indexOf('two');
    //     // if (myIndex !== -1) {
    //     //     myArray.splice(myIndex, 1);
    //     // }
    //
    //     Кнопка.dataset.deleted = JSON.stringify(МассивУдалённых)
    // } else{
         // console.log(метка);
        метка = Number(метка)
        метка =метка+1
        Кнопка.dataset.count = метка
    // }
     // console.log(метка);
    let шаблонАргумента = ШаблоныПолей[ИмяШаблона](метка)

    // let шаблонАргумента = `<div class="flex-row argument" id="arg_wrapper-${метка}">
    //         <label for="arg_name-${метка}">
    //             <div class="input_name">Название</div>
    //             <input type="text" name="аргументы[${метка}].имя_аргумента" id="arg_name-${метка}" placeholder="">
    //         </label>
    //         <label for="arg_id-${метка}">
    //             <div class="input_name">идентификатор</div>
    //             <input type="text" name="аргументы[${метка}].идентификатор" id="arg_id-${метка}" placeholder="">
    //         </label>
    //         <label for="arg_type-${метка}">
    //             <div class="input_name">тип</div>
    //             <select name="аргументы[${метка}].тип" id="arg_type-${метка}">
    //                  <option selected >Выбрать Тип</option>
    //                 <option value="date">Дата</option>
    //                 <option value="number">Число</option>
    //                 <option value="text">Текст</option>
    //             </select>
    //         </label>
// <!--            <button data-tpl="argument_tpl" data-target="arguments_wrapper" data-number="${метка}" onclick="ДобавитьАргументы(event)" class="plus" type="button">+</button>-->
//             <button data-target="arg_wrapper-${метка}" data-number="${метка}" class="minus" onclick="УдалитьАргумент(event)" type="button">-</button>
//          </div>`

    if (!!Кнопка.dataset.target && Кнопка.dataset.target != "") {
        let Цель = document.getElementById(Кнопка.dataset.target)
        Цель.insertAdjacentHTML('beforeend',шаблонАргумента )
    }

}

function УдалитьАргумент (event) {
    // ВсеАргументы = document.querySelectorAll( ".argument")
    НомерУдаляемого = Number(event.target.dataset.number)
     // console.log(НомерУдаляемого);

    let Цель = document.getElementById(event.target.dataset.target)
     // console.log(Цель);
    Цель.remove()
    // ВсеАргументы = document.querySelectorAll("."+event.target.dataset.class_name)

     // console.log(ВсеАргументы);
    // МассивУдалённых=JSON.parse(document.getElementById("btn_add_argument").dataset.deleted)
    // МассивУдалённых.push()
    // МассивУдалённых.sort()
    // document.getElementById("btn_add_argument").dataset.deleted = JSON.stringify(МассивУдалённых)
    let СледующийНомер = -1
   //
   // let  ИдБлока = event.target.dataset.target
   //   console.log(ИдБлока);
   //  let ИмяБлока = ИдБлока.split("-",-1)[0]
   //   console.log(ИмяБлока);
   //  // for (арг of ВсеАргументы){
   //  //      // console.log("арг", арг);
   //  //     НомерПеребираемогоАргумента = Number(арг.id.split("-",-1)[1])
   //  //
   //  //      console.log("НомерПеребираемогоАргумента",НомерПеребираемогоАргумента, "НомерУдаляемого", НомерУдаляемого) ;
   //  //      // console.log("СледующийНомер", СледующийНомер);
   //  //     if (НомерПеребираемогоАргумента > НомерУдаляемого){
   //  //          // console.log("поменять");
   //  //          if (СледующийНомер == -1) {
   //  //              // арг.id = "arg_wrapper-"+НомерУдаляемого
   //  //              арг.id = ИмяБлока+"-"+НомерУдаляемого
   //  //              for (child of арг.children){
   //  //                  if (child.tagName === "LABEL") {
   //  //                      child.for = child.for.split("-",-1)[0]
   //  //                      let input = child.control
   //  //                      input.id = input.id.split("-",-1)[0]+НомерУдаляемого
   //  //                      name = input.name.split(/[\[\]]/)
   //  //                       console.log(name);
   //  //                  }
   //  //              }
   //  //              СледующийНомер = НомерУдаляемого+1
   //  //          } else {
   //  //              арг.id = ИмяБлока+"-"+НомерУдаляемого
   //  //              for (child of арг.children){
   //  //                  if (child.tagName === "LABEL") {
   //  //                      child.for = child.for.split("-",-1)[0]
   //  //                      let input = child.control
   //  //                      input.id = input.id.split("-",-1)[0]+НомерУдаляемого
   //  //                      name = input.name.split(/[\[\]]/)
   //  //                      console.log(name);
   //  //                  }
   //  //              }
   //  //              // арг.id = "arg_wrapper-"+СледующийНомер
   //  //              СледующийНомер++
   //  //          }
   //  //     }
   //  //     // console.log("ВсеАргументы.length", ВсеАргументы.length, "НомерУдаляемого", НомерУдаляемого, "СледующийНомер",СледующийНомер);
   //  //
   //  //     if (ВсеАргументы.length == НомерУдаляемого) {
   //  //         // console.log("ВсеАргументы.length", ВсеАргументы.length, "НомерУдаляемого", НомерУдаляемого, "НомерУдаляемого-1", (НомерУдаляемого-1));
   //  //         СледующийНомер=НомерУдаляемого-1
   //  //     }
   //  //     // console.log("арг", арг);
   //  //
   //  // }
   //   console.log(ИмяБлока, "btn_add_"+ИмяБлока);
   //  document.getElementById("btn_add_"+ИмяБлока).dataset.count=(ВсеАргументы.length-1)
}