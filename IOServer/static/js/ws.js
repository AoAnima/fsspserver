function include() {
    let script = document.createElement('script');
     console.log(window.location.hostname, window.location.hostname !== "localhost");
    if (window.location.hostname !== "10.26.6.25" && window.location.hostname !== "localhost"){
        script.src = 'http://10.26.6.25/static/js/tpl.js';
    } else {
        script.src = '/static/js/tpl.js';
    }
    console.log(script);
    document.getElementsByTagName('head')[0].appendChild(script);
}

include();
// UserInfo struct {
//     Initials string `json:"Initials"`
//     FullName string `json:"FullName"`
//     FirstName string `json:"FirstName"`
//     LastName string `json:"LastName"`
//     MiddleName string `json:"MiddleName"`
//     OspName string `json:"OspName"`
//     OspNum int `json:"OspNum"`
//     PostName string `json:"PostName"`
//     Инициалы string `json:"Инициалы"`
//     ПолноеИмя string `json:"ПолноеИмя"`
//     Фамилия string `json:"Фамилия"`
//     Имя string `json:"Имя"`
//     Отчество string `json:"Отчество"`
//     ОСП string `json:"ОСП"`
//     КодОСП int `json:"КодОСП"`
//     Должность string `json:"Должность"`
// } `json:"UserInfo"`
let WS = {
    Наблюдатель:{},
    AdminMenu: {},
    module: "",
    ws: {},
    chatLog: null,
    messageBox: null,
    contactsList: null,
    uid: null,
    UserInfo: {
        uid: null,
        FullName: null,
        Initials: null,
        ПолноеИмя: null,
        Инициалы: null,
        Фамилия: null,
        Имя: null,
        Отчество: null,
        ОСП: null,
        КодОСП: null,
        Должность: null,
    },

    LastChats: null,
    // Сообщение_ : {
    //     Id:"",
    //     От:"",
    //     Кому:"",
    //     Текст:"",
    //     Время:"",
    //     ОтветНа:"",
    //     Файлы :"",
    //     send:function(){
    //         console.log(this);
    //         return JSON.stringify(this);
    //     }
    // },
    Сообщение: function (m) {
        this.f5=false;
        // this.f5 = !!m.Время && m.Время !== Date.now(); // то же самое что if (!!m.Время && m.Время !== Date.now()) {this.f5 =true}
        this.Id = m.Id ? m.Id : -0;
        this.От = m.От ? m.От : WS.uid;
        this.Кому = m.Кому ? m.Кому : "io";
        this.Текст = m.Текст ? m.Текст : null;
        this.Время = m.Время ? m.Время : Date.now();
        this.ОтветНа = m.ОтветНа ? m.ОтветНа : null;
        this.Файлы = m.Файлы ? m.Файлы : null;
        this.Выполнить = m.Выполнить ? m.Выполнить : null;
        this.ОбратныйВызов = m.ОбратныйВызов ? m.ОбратныйВызов : null;
        this.ПодготовитьКОтправке = function () {
            // console.log(this);
            return JSON.stringify(this);
        };
        return this;
    },

    init: function (uid, reconect = false, render = true) {

        // this.uid = uid;
        // this.UserInfo.uid = uid;
        let user="";
        if (!!uid || uid !== "") {
            this.uid = uid;
            this.UserInfo.uid = uid;
            user = "?login=" + this.uid
        }
        console.log("подключение...");
        // this.ws = new WebSocket("ws://10.26.4.20:8080/chat?user="+uid);
        if (reconect) {
            reconect = "&reconect=true"
        } else {
            reconect = ""
        }
        if (render) {
            render = "&render=true"
        } else {
            render = ""
        }
        let selfWs = this;


        // if (!!uid || uid === "") {
        //     // user= "?user=kanaev@r26"
        //     user = "?user=" + this.uid
        // } else {
        //     console.log("window.location.search", window.location.search === "");
        //     if (!!window.location.search && window.location.search !== "") {
        //         user = window.location.search
        //     }
        // }
        // console.log("user", user);
        // let port=8181;
        // if (window.location.port !== "" || !window.location.port || window.location.port!==8080){
        //     window.location.port
        // }

        // const  wsUrl = "ws://"+window.location.hostname+"/wsconnect"+user+reconect+render;//:8181 ":"+port+
        // const wsUrl = "ws://localhost/wsconnect" + user + reconect;//:8181 ":"+port+
        //

        console.log(WS.host, window.location.host)
        let wsUrl = "ws://10.26.6.25/wsconnect" + user + reconect;//:8181 ":"+port+
        if (!!WS.host) {
            if ( window.location.host==='10.26.6.25:8080') {
                WS.host='10.26.6.25:8080'
            }
            wsUrl =  "ws://"+WS.host+"/wsconnect" + user + reconect;
        } else {
            if ( window.location.host==='10.26.6.25:8080') {
                WS.host='10.26.6.25:8080'
                wsUrl =  "ws://"+WS.host+"/wsconnect" + user + reconect;
            }
            console.log(WS.host, window.location.host)

        }



        this.ws = new WebSocket(wsUrl);

        // дополним стандртный отправщик сообщений
        let originalSender = this.ws.send;
        this.ws.send = (...args) => {
            // console.log("args", args);
            let сообщение = args[0];
            let НеЛогировать = args[1]; // если передано любое значениеотличное от null то переданое сощение не будетлогироваться в url
            const хэш=сообщение.ПодготовитьКОтправке()
            if (!НеЛогировать){
                if (window.location.host !== "10.26.4.20") {
                    history.pushState("stateSendMessage", "title", "/?toio=" + хэш);

                }
            } else {
                 console.log("НеЛогировать", сообщение);
            }
            return originalSender.apply(this.ws, [хэш])

        };
        console.log("Устанаваливаем соединение ", this.ws);

        this.ws.onopen = function (evt) {
            console.log( "Соединение установленно!")
            let m = new WS.Сообщение({
                Текст: "Соединение установленно!",
                Время: Date.now(),
                От: "io",
                Кому: uid,
                Id: 0,
            });
            if (!!document.getElementById('ws_widget')) {
                WS.printMessage(m, "system info", document.getElementById("log_io"));
            }
            let urlParams = new URLSearchParams(window.location.search);
            let ToIO = urlParams.getAll('toio');

            if (ToIO.length>0 ){
                let СообщениеИО =  new WS.Сообщение(JSON.parse(ToIO[0]))
                console.log("Сообщениеm",СообщениеИО);
                СообщениеИО.f5 = true;
                СообщениеИО.От= uid;
                WS.ws.send(СообщениеИО,"нелогировать");
            } else {

                // if (window.location.hostname !== '10.26.4.20') {
                //     console.log(uid)
                //     // if (!!uid && uid != 'auth') {
                //     //     let СообщениеИО = new WS.Сообщение({
                //     //         Текст: "Создать рабочий стол",
                //     //         Время: Время(),
                //     //         От: uid,
                //     //         Кому: "io",
                //     //         Выполнить: {
                //     //             Действие: {"СоздатьРабочийСтол": {}}
                //     //         },
                //     //     });
                //     //     // WS.ws.send(СообщениеИО, "нелогировать");
                //     // }
                // }
            }
            console.log("Соединение установленно!", this.readyState);
            return this;
        };

        this.ws.onclose = function (evt) {

            console.log("Соединение потеряно. Пытаюсь подключиться", "this.uid", uid, evt);

            let m = new WS.Сообщение({
                Текст: "Соединение потеряно. ",
                Время: Время(),
                От: "io",
                Кому: uid,
                Id: "",
            });
            if (!!document.getElementById('ws_widget')) {
                WS.printMessage(m, "system error", document.getElementById("log_io"));
            }


            WS.ws = null;

            if (uid !== 'shhedrina@r26') { // постоянно переподключаеться, нужно проверить в чем проблема
                setTimeout(function () {
                    WS.init(uid, true);
                }, 3000)
            }
            // WS.init(uid, true);
        };
        this.ws.onmessage = function (evt) {

            let сообщение = JSON.parse(evt.data);
            if (!this.LastChats) {
                this.LastChats = document.getElementById("ws_last_chats");
            }
            let UserChatLog = document.getElementById("log_" + сообщение.От);
            let UserBlock = document.getElementById(сообщение.От);
            console.log("сообщение с сервера ", сообщение);
            if (!!сообщение.MessageType && сообщение.MessageType.includes("float_message") ) {
                WS.FloatMessage(сообщение);
            }

            if (сообщение.От === "io") {
                WS.IOMessageHandler(сообщение)
            }
            if (!!UserChatLog && сообщение.Текст !== "") {
                WS.printMessage(сообщение, "mes_to", UserChatLog);
            }
            if (!!UserBlock) {
                UserBlock.classList.add("new_message");
                let Clone = UserBlock.cloneNode(true);
                UserBlock.remove();
                this.LastChats.prepend(Clone);
            }

        };

        this.ws.onerror = function (evt) {
            console.log("onerror", evt);
            // console.log(this);
            let m = new WS.Сообщение({
                Текст: "Ошибка " + JSON.stringify(evt),
                Время: Date.now(),
                От: "io",
                Id: "",
                Кому: this.uid
            });
            console.log(document.getElementById('ws_widget'));
            if (!!document.getElementById('ws_widget')) {
                WS.printMessage(m, "system error", document.getElementById("log_io"));
            }


            if (this.readyState === 3) {
                console.log(this.readyState);
                // setTimeout(function () {
                //     this.ws.close();
                //     WS.init(uid, true);
                // }, 3000);
                // console.log("Не удаёться установить соединение с сервером? Проверим статус сервера");

                // let request = new Request("http://10.26.6.25/home/io/WSServer/start.php?cmd=checkWSServer");
                // request.headers.append("Request-method", "ajax/fetch");
                // fetch(request, {
                //     credentials: "include",
                //     mode: 'cors'
                // }).then(function (response) {
                //     console.log(response)
                // }).catch(function (response) {
                //     console.log(response)
                // });
            }

        };


        return this;
    },
    CloseRichEditor:function(сообщение){
       let ИДОкнаРедактора = сообщение.Контэнт.контейнер

        if (!!ИДОкнаРедактора){
            let ОкноРедактора = document.getElementById(ИДОкнаРедактора)
            if (!!ОкноРедактора) {
                ОкноРедактора.classList.add("hided"); // анимация сворачивания окна
                ОкноРедактора.parentNode.removeChild(ОкноРедактора)
            }
        }
    },
    IOMessageHandler: function (сообщение) {
        // console.log("IOMessageHandler сообщение", сообщение)

        if(!!сообщение.ОбратныйВызов) {
             console.log("СкрытьФормуВхода", сообщение.ОбратныйВызов);

            if (WS.hasOwnProperty(сообщение.ОбратныйВызов)){
                WS[сообщение.ОбратныйВызов](сообщение)
            }
            if (!!window[сообщение.ОбратныйВызов]){
                window[сообщение.ОбратныйВызов](сообщение)
            }
        }

        // console.log("сообщение.Контэнт: ", сообщение.Контэнт);

        if (!!сообщение.Content.обработчик || (!!сообщение.Контэнт && !!сообщение.Контэнт.обработчик)) {
             // console.log("обработчик ", сообщение.Контэнт.обработчик);
            if (сообщение.Content.обработчик == "table") {
                CreateTable(сообщение.Content.target, сообщение.Content);
            } else {
                if (WS.hasOwnProperty(сообщение.Контэнт.обработчик)){
                    WS[сообщение.Контэнт.обработчик](сообщение)
                }
            }
        }

        if (!!сообщение.Выполнить) {
            // console.log(window[сообщение.Выполнить.action]);
            if (!!window[сообщение.Выполнить.action]) {
                window[сообщение.Выполнить.action](сообщение.Выполнить.Arg)
            }
        }

        if (!!сообщение.AdminMenu) {
            this.AdminMenu = сообщение.AdminMenu;
        }


        if (!!сообщение.Online) {
            WS.ПользовательПодключён(сообщение)
        }
        if (!!сообщение.Offline) {
            WS.ПользовательОтключился(сообщение)
        }
        // console.log("сообщение.Контэнт", сообщение.Контэнт);

        if  (!!сообщение.Контэнт && !!сообщение.Контэнт.контейнер){
            if (сообщение.Контэнт.контейнер === "сообщение"){
                WS.printMessage(сообщение)
            } else {
                let контейнер;
                // console.log(сообщение.Контэнт.контейнер)
                let контейнеры = сообщение.Контэнт.контейнер.split(".");
                console.log("контейнеры", контейнеры)
                if (контейнеры.length>1){
                    // let НомерИДКонтейнера = контейнеры.length-1; // проверим самый последний контейнер в списке, если он есть, то обновляем его содержимое, иначе проходим по всем конетйнерам и обновляем первый попавшийся
                    let i=контейнеры.length;
                    // console.log(i);
                    for (i; i >= 0;i--){
                        // console.log(i);
                        // console.log("контейнеры[НомерИДКонтейнера]",i, контейнеры[i-1]);

                        const ТабличныеБлоки = ['table', 'tBodies', 'tFoot', 'tHead'];
                        let контейнер = контейнеры[i-1];
                        let ТабличныйБлокДляВставки;
                        let БлокДанных = document.getElementById(контейнер);

                         // console.log("контейнер",контейнер);
                         if (контейнер === 'modal') {
                             console.log(контейнеры)
                           let окно =  new МодальноеОкно({
                                ИдМодала: "modal_"+контейнеры[i],
                                ИдКонтэнта: "modal_content_"+контейнеры[i],
                                Контэнт:сообщение.Контэнт.html,
                             })
                              console.log(окно);
                         }

                         if (ТабличныеБлоки.includes(контейнер)) {
                            ТабличныйБлокДляВставки = контейнер;
                            i--;
                             console.log("контейнеры[i-1]", контейнеры[i-1]);
                            БлокДанных = document.getElementById(контейнеры[i-1]);
                            console.log("БлокДанныхконтейнер[i-1]",БлокДанных);
                        }

                        console.log("i =",i, "контейнеры.length =", контейнеры.length,"контейнер =", контейнер, "!!БлокДанных=", !!БлокДанных)

                        if(!!БлокДанных){ // Блок для обновления найден
                            if (i===контейнеры.length){
                                // если Контэнт.html
                                if (!!сообщение.Контэнт.html) {
                                    let fragment = document.createRange().createContextualFragment(сообщение.Контэнт.html)
                                    console.log("fragment", fragment)

                                    if (!!ТабличныйБлокДляВставки) {
                                        console.log("ТабличныйБлокДляВставки",ТабличныйБлокДляВставки)
                                        if (ТабличныйБлокДляВставки === "tBodies"){
                                            console.log( БлокДанных[ТабличныйБлокДляВставки])
                                            БлокДанных[ТабличныйБлокДляВставки][0].append(fragment)
                                        } else {
                                            console.log( БлокДанных[ТабличныйБлокДляВставки])
                                            БлокДанных[ТабличныйБлокДляВставки].append(fragment)
                                        }
                                    } else {
                                        console.log("Блок для обновления найден replaceWith",БлокДанных);
                                        console.log("fragment", fragment)

                                        if (!!сообщение.Контэнт.способ_вставки){
                                            console.log("способ_вставки", сообщение.Контэнт.способ_вставки);
                                            if (сообщение.Контэнт.способ_вставки === "обновить"){

                                                // контейнер.innerHTML= сообщение.Контэнт.html;
                                                БлокДанных.classList.add("update")
                                                setTimeout(function(){
                                                    БлокДанных.innerHTML= сообщение.Контэнт.html;
                                                    setTimeout(()=>БлокДанных.classList.remove("update"),200)

                                                },200)

                                            } else if (сообщение.Контэнт.способ_вставки === "добавить"){

                                                // контейнер.innerHTML= сообщение.Контэнт.html;
                                                БлокДанных.classList.add("update")
                                                 console.log("добавить");
                                                 console.log(БлокДанных);
                                                 console.log(сообщение.Контэнт.html);
                                                БлокДанных.append( сообщение.Контэнт.htm);
                                                // setTimeout(function(){
                                                //     БлокДанных.insertAdjacentHTML("beforeEnd", сообщение.Контэнт.htm);
                                                //     setTimeout(()=>БлокДанных.classList.remove("update"),200)
                                                //
                                                // },200)

                                            }else {
                                                 console.log("замена 1", БлокДанных);
                                                БлокДанных.replaceWith(fragment);
                                            }
                                        }else {
                                            console.log("замена 2", БлокДанных);
                                            БлокДанных.replaceWith(fragment);
                                        }



                                    }
                                }
                                if (!!сообщение.Контэнт.данные){
                                    // вероятней всего будет массив данных, с ключами равными id куда вставить значение
                                    // поэтому проходим по каждому элементу массива и вставляем данные
                                    for (let Данные of сообщение.Контэнт.данные) {
                                        // пока проверим есть контейнер и ключ равны то вставим данные в контенер, потом можно будет подумать как ещё обрабатывать
                                        // if (сообщение.Контэнт.контейнер=== Данные)
                                        console.log(Данные)
                                        if (!!Данные.option) {
                                            console.log(Данные.option)
                                            // данные для вставки в select
                                            Данные.option.forEach(function(el) {
                                                    console.log(el);
                                                    var opt = document.createElement("option");
                                                    opt.value = el.key
                                                    opt.text = el.value
                                                    БлокДанных.add(opt, null)
                                                }
                                            );
                                            // for (el in Данные.option){
                                            //      console.log(el);
                                            //     var opt = document.createElement("option");
                                            //     opt.value = el.key
                                            //     opt.text = el.value
                                            //     БлокДанных.add(opt,null)
                                            // }

                                        }

                                    }
                                }
                            } else {
                                if (i === 1){
                                    if (!!сообщение.Контэнт.html) {
                                        let fragment = document.createRange().createContextualFragment(сообщение.Контэнт.html);
                                        console.log("fragment", fragment)

                                        if (!!ТабличныйБлокДляВставки) {
                                            if (ТабличныйБлокДляВставки === "tBodies") {
                                                БлокДанных[ТабличныйБлокДляВставки][0].insertAdjacentHTML("beforeend", сообщение.Контэнт.html)
                                            } else {
                                                БлокДанных[ТабличныйБлокДляВставки].append(fragment)
                                            }
                                        } else {
                                            console.log("Блок для обновления найден innerHTML",БлокДанных);
                                            БлокДанных.innerHTML = сообщение.Контэнт.html;

                                        }
                                    }

                                } else {
                                    if (!!ТабличныйБлокДляВставки) {
                                        if (!!сообщение.Контэнт.html) {
                                            let fragment = document.createRange().createContextualFragment(сообщение.Контэнт.html);
                                            console.log("fragment", fragment)

                                            if (ТабличныйБлокДляВставки === "tBodies") {

                                                БлокДанных[ТабличныйБлокДляВставки][0].append(fragment)

                                            } else {
                                                console.log("Блок для обновления найден append",БлокДанных);
                                                if (ТабличныйБлокДляВставки == "table") {
                                                    var tBody = document.createElement ("tbody");
                                                         tBody.innerHTML=сообщение.Контэнт.html;
                                                         console.table(tBody);
                                                        БлокДанных.appendChild (tBody);
                                                    // БлокДанных.append(fragment)
                                                } else {
                                                    БлокДанных[ТабличныйБлокДляВставки].append(fragment)
                                                }


                                            }
                                        }

                                    } else {
                                        if (!!сообщение.Контэнт.html) {
                                            console.log("Блок для обновления найден beforeend",БлокДанных);
                                            БлокДанных.insertAdjacentHTML("beforeend", сообщение.Контэнт.html);
                                        }
                                    }
                                }
                            }
                            break
                        } else {
                             // console.log("контейнер",контейнер);
                             // console.log("БлокДанных", БлокДанных);
                             // console.log("i", i);
                            // НомерИДКонтейнера = НомерИДКонтейнера-1
                            if (i===1) {
                                console.log("контейнер не найден:",контейнер);
                            }
                        }
                    }

                    console.log(i)
                } else {
                    контейнер = document.getElementById(сообщение.Контэнт.контейнер);
                    console.log(контейнер)
                    /*
                    * Обновление таблицы в реальном времени по мере получения данных
                    * */
                    if (!!контейнер && контейнер.tagName === "TABLE"){
                        // if (контейнер.tagName === "TABLE") {
                            ВставитьДанныеВТаблицуИзРБД(контейнер, сообщение)
                        // }
                    }



                    if (!!контейнер && контейнер.tagName !== "TABLE" && !!сообщение.Контэнт.html && сообщение.Контэнт.html != ""){
                        let fragment = document.createRange().createContextualFragment(сообщение.Контэнт.html)
                        console.log(сообщение.Контэнт.способ_вставки)
                        console.log(сообщение.Контэнт.способ_вставки=== "обновить")

                        if (!!сообщение.Контэнт.способ_вставки){
                            if (сообщение.Контэнт.способ_вставки === "обновить"){
                                // контейнер.innerHTML= сообщение.Контэнт.html;
                                контейнер.classList.add("update")
                                setTimeout(function(){
                                    контейнер.innerHTML= сообщение.Контэнт.html;
                                    setTimeout(()=>контейнер.classList.remove("update"),200)
                                    let БлокиСРедкатором = контейнер.getElementsByClassName("suneditor");
                                    if (!!БлокиСРедкатором && БлокиСРедкатором.length >0){
                                        for (Блок of БлокиСРедкатором){
                                            if (!!Блок){
                                                console.log(Блок, Блок.id);
                                                if (!document.getElementById("suneditor_"+Блок.id)) {
                                                    console.log("создаём редактор в",Блок, Блок.id);
                                                    // initRichEditor(Блок);
                                                    Редактор(Блок);
                                                } else {
                                                    console.log('Редактор уже создан для блока ', Блок.id, document.getElementById("suneditor_"+Блок.id));
                                                }

                                            }
                                        }
                                    }
                                },200)
                            }  else if (сообщение.Контэнт.способ_вставки === "заменить"){
                                контейнер.replaceWith(document.createRange().createContextualFragment(сообщение.Контэнт.html))
                            } else {
                                console.log(сообщение.Контэнт.способ_вставки)
                                console.log("контейнер", контейнер)
                                контейнер.insertAdjacentHTML("beforeend", сообщение.Контэнт.html);
                            }
                        }else {
                            console.log("контейнер", контейнер)
                            контейнер.insertAdjacentHTML("beforeend", сообщение.Контэнт.html);
                        }
                    }

                    // if (!контейнер && !!сообщение.Контэнт.данные){
                    if (!!контейнер && !!сообщение.Контэнт.данные){
                        // вероятней всего будет массив данных, с ключами равными id куда вставить значение
                        // поэтому проходим по каждому элементу массива и вставляем данные
                        console.log("ЕСЛИ СЮДА ПОПАЛО ТО НАДО ПРОВЕРИТЬ ОТКУДА СЮДА ПРИХОДЯТ ЕЩЁ ДАННЫЕ, пока приходят данные в виде карты из диалога 'удалить файл из журанала скзс'", сообщение.Контэнт.данные)
                        console.log(typeof сообщение.Контэнт.данные)


                        for (let Цель in сообщение.Контэнт.данные){
                            console.log(Цель)
                            let контейнер = document.getElementById(Цель);
                                if (!!контейнер ) {
                                    if (контейнер.tagName === "INPUT") {
                                        контейнер.value = сообщение.Контэнт.данные[Цель]
                                    } else if (контейнер.tagName === "DIV") {
                                        контейнер.innerHTML = сообщение.Контэнт.данные[Цель]
                                    }
                                }
                        }

                        // for (let Данные of сообщение.Контэнт.данные) {
                            // пока проверим есть контейнер и ключ равны то вставим данные в контенер, потом можно будет подумать как ещё обрабатывать
                            // if (сообщение.Контэнт.контейнер=== Данные)

                            // console.log(Данные)
                            // console.log(Object.keys(Данные))
                            // for (let Цель in Данные) {
                            //     console.log(Цель)
                            //     let контейнер = document.getElementById(Цель);
                            //     if (!!контейнер ) {
                            //         if (контейнер.tagName === "INPUT") {
                            //             контейнер.value = Данные[Цель]
                            //         } else if (контейнер.tagName === "DIV") {
                            //             контейнер.innerHTML = Данные[Цель]
                            //         }
                            //     }
                            // }
                        // }
                    }
                }

                    let ЭлементыСНаблюдателем = document.querySelectorAll('[data-observ]')
                    // for (let Элемент of ЭлементыСНаблюдателем){
                    //     console.log("Элемент",Элемент)
                    //     console.log(this.Наблюдатель[Элемент])
                    //
                    //     // let funct = 'ObservElem(){'+Элемент.dataset.observ+'}';
                    //
                    //    // let F=new Function('a',Элемент.dataset.observ);
                    //
                    //     // if (!!this.Наблюдатель[Элемент]) {
                    //     //     console.log(this.Наблюдатель[Элемент])
                    //     // } else  {
                    //     //     let config = { attributes: true, childList: true, characterData: true, subtree: true }
                    //     //     let F=new Function('a',Элемент.dataset.observ);
                    //     //
                    //     //     this.Наблюдатель[Элемент] = new MutationObserver(изменения => {
                    //     //         console.log("изменения",изменения); // console.log(изменения)
                    //     //         console.log(Элемент.dataset.observ);
                    //     //         console.log(" F",  F);
                    //     //         F();
                    //     //     });
                    //     //     // this.Наблюдатель[Элемент] = new MutationObserver(F);
                    //     //     this.Наблюдатель[Элемент].observe(Элемент, config);
                    //     // }
                    // }


                // let контейнер = document.getElementById(сообщение.Контэнт.контейнер);
                // console.log(контейнер);
                // if (!!контейнер){
                //     контейнер.insertAdjacentHTML("beforeend", сообщение.Контэнт.html);
                // } else {
                //     console.log("контейнер ", сообщение.Контэнт.контейнер);
                // }

            }

        }

        if ((!!сообщение.Content && !!сообщение.Content["target"])){
            // console.log("io сообщение.Content.target", сообщение.Content["target"]);

            let target = null;
            if (сообщение.Content["target"] !== "body" ) {
                target = document.getElementById(сообщение.Content["target"]);
            } else {
                // если таргет body то предполагаем что авторизация пройдена успешно,
                target = document.body;
                //
                if (!!сообщение.Token && !!сообщение.Token.Hash){
                    console.log(сообщение)
                    console.log(window.location.hostname)
                    document.cookie = `Token=${сообщение.Token.Hash};expires=${сообщение.Token.Истекает};path=/;domain=${window.location.hostname}`;
                    console.log(document.cookie)

                    СкрытьФормуВхода()
                }
            }

            if (!!target) {
                // console.log("сообщение", сообщение);
                // console.log("сообщение.Content.data", сообщение.Content.data);

                if (сообщение.Content.html && сообщение.Content.html !== "") {

                    if (сообщение.Content["target"] === "log_"+сообщение.От){
                        WS.printMessage(сообщение)
                        // target.insertAdjacentHTML("beforeend", сообщение.Content["html"]);
                    } else {
                        target.insertAdjacentHTML("beforeend", сообщение.Content["html"]);

                    }

                }
                //TerminalLog
                if (!!сообщение.Content.data && сообщение.Content.data !== "") {
                    // console.log("сообщение TerminalLog", сообщение.Content.data);

                    if (!!сообщение.Content.data["TerminalLog"]) {
                        let LogStr = this.tpl("terminalLogString", сообщение.Content.data["TerminalLog"]);
                        target.insertAdjacentHTML("beforeend", LogStr);
                        target.scrollTop = target.scrollHeight;
                        if (!!сообщение.Content.data["SSHclient"]) {
                            let clientCmd = document.getElementById("cmd_" + сообщение.Content.data["SSHclient"].login);

                            // console.log("clientCmd",clientCmd);
                            // console.log(clientCmd.dataset.login);

                            if (!!clientCmd) {

                                if (!clientCmd.dataset.login || !!clientCmd.dataset.login && clientCmd.dataset.login === "") {
                                    console.log(сообщение.Content.data["SSHclient"].login);
                                    clientCmd.dataset.login = сообщение.Content.data["SSHclient"].login;

                                }
                                if (!clientCmd.dataset.ip || !!clientCmd.dataset.ip && clientCmd.dataset.ip === "") {
                                    console.log(сообщение.Content.data["SSHclient"].ip);
                                    clientCmd.dataset.ip = target.dataset.ip = сообщение.Content.data["SSHclient"].ip;
                                }
                            }
                        }
                        // }
                        // target.append(fragment);
                    }


                    // Если в поле data есть поле innerText то нужно данные вставить в target как текст и обновить
                    if(!!сообщение.Content.data["innerText"]){
                        target.innerText=сообщение.Content.data["innerText"];

                        // if (target.nodeName==="INPUT"){
                        //     target=target.parentNode;
                        // }
                        target.classList.add("saved");
                        console.log(target)
                        // setTimeout(function(target){
                        //     console.log(target);
                        //     target.classList.remove("saved")
                        // },5000, target)
                    }
                    console.log(сообщение.Content);
                    if(!!сообщение.Content.data["checked"]){
                        target.checked=сообщение.Content.data["checked"];

                        if (target.nodeName==="INPUT"){
                            target=target.parentNode;
                        }
                        target.classList.add("saved");
                        setTimeout(target=>{
                            target.classList.remove("saved")
                        },5000, target)
                    }
                }

                if (!this.chatLog) {
                    this.chatLog = document.getElementById("ws_chat_log");
                    this.messageBox = document.getElementById("ws_message_box");
                }
                if (!this.LastChats) {
                    this.LastChats = document.getElementById("ws_last_chats");
                }

                if (сообщение.Content["target"] === "body" ) { //&& сообщение.Content["target"] === "log_io"
                    console.log(сообщение);
                    // Так как вставка идёт в тело страницы, то 99% что это либо обновление либо открытие страницы. Поэтому пологаем что это инциилизация страницы.
                    let scroller = target.getElementsByClassName("scroll_wrapper");
                    if (scroller.length>0){
                        scroller[0].scrollTop = scroller[0].scrollHeight;
                    }
                    // let Время = PgTimeStamp();
                     // console.log("приветствие Время" , Время);
                    // Сообщим боту что Клиент мессенджера создан, чтобы он прислал лог сообщений и создал рабочий стол
                    // let mes = new this.Сообщение({
                    //     Id: Date.parse(Время) / 1000,
                    //     Текст: "приветствие",
                    //     Время: Время,
                    //     Кому: "io",
                    //     От: this.uid,
                    //     Выполнить: {"action": "ПолучитьЛогЧата", "arg": {"Login": "io"}},
                    // });
                    //
                    // this.ws.send(mes,"нелогировать");
                    //
                    //
                    // // Алгоритим: т.к. эта первая загрузка, проверим что наодиться в url и если там есть search toio то отправим соообщение с его содержимым

                    // let urlParams = new URLSearchParams(window.location.search);
                    // let ToIO = urlParams.getAll('toio');
                    //
                    //  if (ToIO.length>0 ){
                    //      let СообщениеИО =  new this.Сообщение(JSON.parse(ToIO[0]))
                    //      console.log("Сообщение",СообщениеИО);
                    //      СообщениеИО.f5 = true;
                    //      СообщениеИО.От= this.uid;
                    //      this.ws.send(СообщениеИО,"нелогировать");
                    // }



                    // let mesCp = new this.Сообщение({
                    //     Id: Date.parse(Время) / 1000,
                    //     Текст: "проверить данные рабочий станции",
                    //     Время: Время,
                    //     Кому: "io",
                    //     От: this.uid,
                    // });
                    //
                    // this.ws.send(mesCp);



                }

            }
        }

        if (!!сообщение.UserInfo && !!сообщение.UserInfo.FullName) {
            this.uid = сообщение.UserInfo.uid
            this.UserInfo = сообщение.UserInfo
            // this.UserInfo.FullName = сообщение.UserInfo.FullName;
            // this.UserInfo.Initials = сообщение.UserInfo.Initials;
            // UserInfo struct {
            //     Initials string `json:"Initials"`
            //     FullName string `json:"FullName"`
            //     FirstName string `json:"FirstName"`
            //     LastName string `json:"LastName"`
            //     MiddleName string `json:"MiddleName"`
            //     OspName string `json:"OspName"`
            //     OspNum int `json:"OspNum"`
            //     PostName string `json:"PostName"`
            //     Инициалы string `json:"Инициалы"`
            //     ПолноеИмя string `json:"ПолноеИмя"`
            //     Фамилия string `json:"Фамилия"`
            //     Имя string `json:"Имя"`
            //     Отчество string `json:"Отчество"`
            //     ОСП string `json:"ОСП"`
            //     КодОСП int `json:"КодОСП"`
            //     Должность string `json:"Должность"`
            // } `json:"UserInfo"`
            // uid: null,
            //     UserInfo: {
            //     uid: null,
            //         FullName: null,
            //         Initials: null,
            //         ПолноеИмя: null,
            //         Инициалы: null,
            //         Фамилия: null,
            //         Имя: null,
            //         Отчество: null,
            //         ОСП: null,
            //         КодОСП: null,
            //         Должность: null,
            // },

            // console.log("Данные пользователя ", this)

        }


        //алгоритм Проверим не появился ли на станице блок с классом ckeditor если появился то инциируем в этом блоке и во всех блоках с классом ckeditor функцию initRichEditor . Нужно проверить не инициирован ли уже редактор в блоке.
        // console.log("SUNEDITOR", SUNEDITOR);
 //       let БлокиСРедкатором = document.getElementsByClassName("suneditor");
 //
 // console.log("БлокиСРедкатором",БлокиСРедкатором);
        //  if (!!БлокиСРедкатором && БлокиСРедкатором.length >0){
        //     for (Блок of БлокиСРедкатором){
        //         if (!!Блок){
        //             console.log(Блок, Блок.id);
        //             if (!document.getElementById("suneditor_"+Блок.id)) {
        //                 console.log("создаём редактор в",Блок, Блок.id);
        //                 // initRichEditor(Блок);
        //                 initRichEditor(Блок);
        //             } else {
        //                 console.log('Редактор уже создан для блока ', Блок.id, document.getElementById("suneditor_"+Блок.id));
        //             }
        //
        //         }
        //     }
        // }

        return this

    },

     Zip: function (сообщение){

        if (сообщение.Контэнт.данные.length > 1) {
     
            XMLS = WS.XMLtemplate(сообщение)
            var zip = new JSZip();

            let cont = {}

            for (имя in XMLS ){

                // console.log(имя, XMLS);
                zip.file(имя, XMLS[имя]);

            }
            zip.generateAsync({type:"blob"}).then(function(content) {
                // see FileSaver.js
               // console.log(content);
                saveAs(content, сообщение.Контэнт.данные[0].osp+"_"+сообщение.Контэнт.данные[0].date+".zip");
            });


        } else {
            var zip = new JSZip();

            zip.file(сообщение.Контэнт.данные[0].osp+сообщение.Контэнт.данные[0].id+".xml", сообщение.Контэнт.html);
            // var img = zip.folder("images");
            // img.file("smile.gif", imgData, {base64: true});
            //  zip.file("hello.txt").async("string").then(function (сообщение.Контэнт.html){
            //      saveAs(content, "example.zip");
            //  });
            zip.generateAsync({type:"blob"}).then(function(content) {
                // see FileSaver.js
                // console.log(content);
                saveAs(content, сообщение.Контэнт.данные[0].osp+сообщение.Контэнт.данные[0].id+".zip");
            });
        }


    },

    XMLtemplate: function (сообщение){

        let XMLs = {};
        for (данные of сообщение.Контэнт.данные) {

            let xml= `<?xml version="1.0" encoding="UTF-8"?>
    <job doc="F360" nfile="${данные.osp}${данные.id}" datform="${данные.date}" datcome="" subject="* Адм. правонарушения - нов." to="a000001" to_desc="" maxrec="300" attachment="" realname="" oper="ins" sost="new" jr="">
            <request> 
            <DOC_F360 num="1">`
            for (ключ in данные) {

                if (данные[ключ] != null && (ключ.toUpperCase() !== "OSP" && ключ.toUpperCase() !== "FILE" && ключ.toUpperCase() !== "LOGIN" && ключ.toUpperCase() !== "DATE" && ключ.toUpperCase() !== "ID" )){
                    let значение = String(данные[ключ])
                    xml+=`<${ключ.toUpperCase()}>${значение.toUpperCase()}</${ключ.toUpperCase()}>`
                }
            }
            xml+=`</DOC_F360>
                    </request>
                    </job>`
            XMLs[данные.osp+данные.id+".xml"]=xml;
        }
return XMLs
    },

    DeleteFromDom:function(сообщение){
         console.log(сообщение.Контэнт.контейнер);
         let удаляемыйБлок =  document.getElementById(сообщение.Контэнт.контейнер)
         if (удаляемыйБлок) {
             удаляемыйБлок.parentNode.removeChild(удаляемыйБлок)
         }

    },
    ПользовательОтключился:function(сообщение){
         // console.log("ПользовательОтключился");
        let userOnline = document.getElementById(сообщение.Offline);
        if (!!userOnline) {
            userOnline.classList.remove("online");
            if (userOnline.parentNode.getAttribute("id") !== "ws_last_chats") {
                document.getElementById("ws_offline").appendChild(userOnline);
            }
        }
    },
    ПользовательПодключён:function(сообщение){
         // console.log("ПользовательПодключён");
        let userOnline = document.getElementById(сообщение.Online);

        if (!userOnline) {
            userOnline = this.tpl("UserWidget", сообщение);
            document.getElementById("ws_online").insertAdjacentHTML("beforeend", userOnline)
        } else if (!!userOnline) {
            userOnline.classList.add("online");
            if (!!userOnline.parentNode && userOnline.parentNode.getAttribute("id") !== "ws_last_chats") {
                document.getElementById("ws_online").appendChild(userOnline);
            }
        }
    },


    FloatMessage:function(сообщение){



        var message = document.createElement("div");
        mBlockId = "mesBlock_"+Date.now();
        message.id = mBlockId;

        message.classList.add("message_float","hide");
        if (!сообщение.MessageType ){
            сообщение.MessageType=["info"];

        }
        for (let type of сообщение.MessageType){
            console.log(type)
            message.classList.add(type)
        }
         console.log(сообщение);
        if (!!сообщение.Контэнт && !!сообщение.Контэнт.данные){
            message.innerHTML="<div>"+сообщение.Контэнт.данные+"</div>";
        } else if (сообщение.Текст){
            message.innerHTML="<div>"+сообщение.Текст+"</div>";
        }

        message.dataset.title="Кликните на сообщении чтобы закрыть";

        document.getElementById("noteContainer").appendChild(message);
        // message.classList.remove("hide");
        setTimeout(function(){
            message.classList.remove("hide");
        }, 100);

        message.onclick = function(){
            message.classList.add("hide");
            setTimeout(function(){
                message.parentNode.removeChild(message);
            },500);
        };

        console.log("сообщение.MessageType !error", !сообщение.MessageType.includes("error"))
        console.log("сообщение.MessageType !info", !сообщение.MessageType.includes("info"))

        if (!сообщение.MessageType.includes("error") && !сообщение.MessageType.includes("info")) {

            setTimeout(function(){
                message.classList.add("hide");
                setTimeout(function(){
                    message.parentNode.removeChild(message);
                },500);
                //document.getElementById(mBlockId).parentNode.removeChild(document.getElementById(mBlockId));
            }, 5000);
        }

        return mBlockId;
    },
    ShowChatBox: function () {
        document.getElementById("ws_widget").classList.add("show_ws")
    },
    CloseChatBox: function () {
        document.getElementById("ws_widget").classList.remove("show_ws")
    },
    ToggleChatBox: function () {
        let  ws_widget = document.getElementById("ws_widget")
        ws_widget.classList.toggle("show_ws");
        if (ws_widget.classList.contains("show_ws")){

        }
    },
    CreateSkill: function () {
        let Время =  PgTimeStamp();
        let mes = new this.Сообщение({
            Id: Date.parse(Время) / 1000,
            Текст: "Новый навык",
            Время: Время,
            Кому: "io",
            От: this.uid,
        });
        // let UserChatLog = document.getElementById("log_io"
        this.printMessage(mes,"mes_to", null );

        this.ws.send(mes);

    },
    ShowSkill: function () {
        let Время =  PgTimeStamp();
        let mes = new this.Сообщение({
            Id: Date.parse(Время) / 1000,
            Текст: "Показать навыки",
            Время: Время,
            Кому: "io",
            От: this.uid,
        });
        // let UserChatLog = document.getElementById("log_io"
        // this.printMessage(mes,"mes_to", null );
        // WS.send(mes);
        this.ws.send(mes);

    },
    //
    ОтправитьСообщениеИО:function (Сообщение, аргументы=null, неЛогировать=null){
        let Время =  PgTimeStamp();
        let ТекстСообщения = Сообщение
         console.log(аргументы);
        if (аргументы !== null){
            ТекстСообщения = null
            аргументы["время"] = Время
        } else {
            аргументы= {"время": Время}
        }

        let mes = new this.Сообщение({
            Id: Date.parse(Время) / 1000,
            Текст: ТекстСообщения,
            Время: Время,
            Кому: "io",
            От: this.uid,
            Выполнить:{
                Действие : {[Сообщение]: аргументы}
            }
        });
        this.printMessage(mes,"mes_to", null );

        this.ws.send(mes, неЛогировать);
    },
    SendFastMessageIO : function(Сообщение, аргументы, неЛогировать=null){

        console.log(аргументы)
        let Время =  PgTimeStamp();

        let mes = new this.Сообщение({
            Id: Date.parse(Время) / 1000,
            Текст: Сообщение,
            Время: Время,
            Кому: "io",
            От: this.uid,
        });
        if (!!аргументы){

            mes.ВходящиеАргументы=аргументы
        } else {
            mes.ВходящиеАргументы = {"время":Время};
        }
        console.log(mes)
        // let UserChatLog = document.getElementById("log_io"
        this.printMessage(mes,"mes_to", null );
        // WS.send(mes);

        this.ws.send(mes, неЛогировать);
    },

    SearchUid: function (event) {

    },

    ServerAction: function (actionObj) { // {"action":"collectData","arg":WS.module}

        let mes = new this.Сообщение({
            Время: Date.now(),
            От: this.UserInfo.uid,
            Кому: "io",
            Id: 0,
            Выполнить: actionObj
        });

        console.log(actionObj);
        console.log(mes);
        console.log(this.ws.readyState);

        if (this.ws.readyState < 1) {
            console.log(actionObj);
            this.ws.addEventListener('open', event => {
                this.ws.send(mes);
            });
        } else if (this.ws.readyState === 1) {
            this.ws.send(mes);
        }

        return this;
    },
    SkillEdit: function(event, id){
        console.log(event) ;
        console.log("id",id);

        let editTarget = event.originalTarget;
        let editFieldName = editTarget.dataset.skillfield;

        let mesWrapper=editTarget;

        while (true){
            mesWrapper = mesWrapper.parentNode;
           if(mesWrapper.classList.contains("message_wrapper")){
               break
           }
        }

        console.log(mesWrapper);
        let СтароеЗначение = editTarget.innerText;
        // if (mesWrapper.classList.contains("history_log")){
        //     let mes = new this.Сообщение({
        //         Время: PgTimeStamp(),
        //         Кому: this.uid,
        //         От: "io",
        //         Текст: `${this.UserInfo.Имя}, Нельзя редактировать старые сообщения`
        //     });
        //     this.printMessage(mes,"mes_from", null );
        //     return this;
        // } else {
        //     СтароеЗначение = editTarget.innerText
        // }
        СтароеЗначение = editTarget.innerText;
        console.log(editTarget,id,editFieldName, editTarget.dataset);
        if (!!id && !!editFieldName){
            if (editTarget.type === "checkbox" ){
                console.log(editTarget.checked);
                console.log(!editTarget.checked);
                let m = new WS.Сообщение({
                    Выполнить: {"Действие": {
                            "ИзменитьНавык": {
                                "НомерНавыка": id,
                                "ИзменяемоеПоле":editFieldName,
                                "НовоеЗначение":editTarget.checked,
                                "СтароеЗначение":!editTarget.checked}
                        }
                    },
                    Время: PgTimeStamp(),
                    Кому: "io",
                    От:this.uid,
                });
                WS.ws.send(m.send())
                return
            }



            editTarget.setAttribute("contenteditable","true");

            editTarget.classList.add("edit");
            editTarget.focus();

            editTarget.addEventListener("focusout",
                function(event){
                    console.log("focusout",event);
                    console.log(id);
                    console.log(this)
                    event.target.classList.remove("edit");
                    event.target.setAttribute("contenteditable","fasle")
                    let НовоеЗначение =event.target.innerText
                    if (НовоеЗначение === СтароеЗначение){
                        return
                    }
                    let m = new WS.Сообщение({
                        Выполнить: {"Действие": {
                                        "ИзменитьНавык": {
                                                "НомерНавыка": id,
                                                "ИзменяемоеПоле":editFieldName,
                                                "НовоеЗначение":НовоеЗначение,
                                                "СтароеЗначение":СтароеЗначение}
                                        }
                                    },
                        Время: PgTimeStamp(),
                        Кому: "io",
                        От:this.uid,
                    });
                    WS.ws.send(m.send())

                }, {once:true},true);


        }
        return this;
    },
    RunSkill: function (id) {
        console.log(id);
        let m = new WS.Сообщение({
            Выполнить: {"Навык": id},
            Время: PgTimeStamp(),
            Кому: "io",
        });



        this.ws.send(m);

        let mes = new WS.Сообщение({
            Текст: "Пытаюсь решить вашу проблему... ",
            Время: PgTimeStamp(),
            От : "io",
            Кому: this.uid,
        });


        WS.printMessage(mes, "mes_to" );

        return this
    },
    RunCmd: function (cmdBlock) {
        console.log(cmdBlock);
        let cmd = cmdBlock.innerText;
        console.log(cmd)
        let login=cmdBlock.dataset.login;
        let ip=cmdBlock.dataset.ip;
        if (!login){

        }
        let m = new WS.Сообщение({
            Выполнить: {"Комманду": cmd, "Arg": {"login": login}},
            Время: PgTimeStamp(),
            Кому: "io",
        });

        this.ws.send(m);
        let Log={
            prefix: `${WS.uid} >> ${login}(${ip})#`,
            text:cmd
        };

        // cmdBlock.innerHTML this.tpl("terminalLogString", TerminalLog);

        // WS.printMessage(m, "mes_from", "ws_terminal_log_"+login);
        // terminalLogString

        let LogStr = this.tpl("terminalLogString", Log);
        console.log(LogStr);
        console.log(cmdBlock);
        let TerminalLog = document.getElementById("ws_terminal_log_"+login);
            TerminalLog.insertAdjacentHTML("beforeend", LogStr);

        cmdBlock.innerHTML = "";
        cmdBlock.focus();


        return this;
    },
    ПолучитьЛогЧата: function (login) {
        // console.log("ПолучитьЛогЧата this", this);
        // console.log("ПолучитьЛогЧата WS", WS);
        let m = new WS.Сообщение({
            Текст: "приветствие",
            Выполнить: {"action": "ПолучитьЛогЧата", "arg": {"Login": login}},
            Время: PgTimeStamp(),
            Кому: "io",
        });

        // WS.printMessage(this.Сообщение, "outcoming" );
        this.ws.send(m, "нелогировать");
        return this
    },
    ActiveUserBlock: null,
    ShowChatLog: function (selectContact, login) {
        this.messageBox.innerHTML = "";
        if (!this.chatLog) {
            this.chatLog = document.getElementById("ws_chat_log");
        }

        for (let UChatLog of this.chatLog.children) {
            UChatLog.classList.add("hiden")
        }
        let UserChatLog = document.getElementById("log_wrapper_" + login);
        let TerminalLog = document.getElementById("ws_terminal_" + login);

        if (!UserChatLog) {
            //let ws_messages = `<div class="log_wrapper" id="log_wrapper_${login}"></div>`;
            let ws_messages = this.tpl("LogWrapper", {login: login, AdminMenu: this.AdminMenu});
            this.chatLog.insertAdjacentHTML("beforeend", ws_messages);
            this.ПолучитьЛогЧата(login);
        } else {
            UserChatLog.classList.remove("hiden");
            let scroller = UserChatLog.getElementsByClassName("scroll_wrapper");
            scroller[0].scrollTop = scroller[0].scrollHeight;
            if (!!TerminalLog) {
                TerminalLog.classList.remove("hiden");
            }

        }

        if (!!this.ActiveUserBlock) {
            this.ActiveUserBlock.classList.remove("active");
        }

        this.ActiveUserBlock = selectContact;
        selectContact.classList.add("active");
        selectContact.classList.remove("new_message");
        if (!this.messageBox) {
            this.messageBox = document.getElementById("ws_message_box");
        }
        this.messageBox.dataset.recipient = login;
        this.messageBox.focus();
        return this;
    },
    ToggleOUList: function (target) {
        let contacts_wrapper = document.getElementById("contacts_wrapper");
        for (let ouList of contacts_wrapper.children) {
            if (ouList != target.parentNode) {
                if (!ouList.classList.contains("closed")) {
                    ouList.classList.add("closed")
                }
            } else {
                ouList.classList.toggle("closed")
            }
        }
        return this;
    },
    stopServer: function () {
        let request = new Request("http://10.26.6.25/WSServer/start.php?cmd=stopWSServer");
        request.headers.append("Request-method", "ajax/fetch");
        fetch(request, {
            credentials: "include",
            mode: 'cors'
        }).then(function (response) {
            console.log(response)
        }).catch(function (response) {
            console.log(response)
        });
    },
    sendMessage: function (messageBox) {
        console.log(messageBox);
        console.log(this.messageBox);

        // if (this.messageBox.innerHTML === ""){
        //     return
        // }
        if (messageBox === "" || messageBox.innerHTML === "" || this.messageBox.innerHTML === "") {
            return
        }
        // console.log(this.messageBox.innerHTML)
        // console.log(this.messageBox.innerText)

        let Время = PgTimeStamp();
        let mes = new this.Сообщение({
            Id: Date.parse(Время) / 1000,
            Текст: this.messageBox.innerText,
            Время: Время,
            Кому: this.messageBox.dataset.recipient,
            От: this.uid,
        });

        let UserChatLog = document.getElementById("log_" + this.messageBox.dataset.recipient);
        let UserBlock = document.getElementById(this.messageBox.dataset.recipient);

        UserBlock.classList.remove("new_message");

        let Clone = UserBlock.cloneNode(true);
        UserBlock.remove();
        this.LastChats.prepend(Clone);

        console.log("UserChatLog", UserChatLog)
        WS.printMessage(mes, "mes_from", UserChatLog);
        this.messageBox.innerHTML = "";
        this.messageBox.focus();
        this.ws.send(mes);
        return this;
    },
    CloseForm: function (event) {
        event.preventDefault();
        event.stopPropagation();
        let form = event.target;
        form.parentNode.removeChild(form);
    },
    ЗагрузитьФайлы: function(event){
        event.preventDefault();
        event.stopPropagation();

        let form = event.target;

        console.log(form)

        let FormAction = event.target.getAttribute("action");
        let Действие={[FormAction]:{}};
        let ДанныеФормы = new FormData(form);

        let ФайлыДляЗагрузки=[]

        for(let [name, value] of ДанныеФормы) {

            if (name.includes("[]")){

                console.log(Действие);
                console.log(name, value);

                let ИмяПоля = name.replace("[]","")
                if (!Действие[FormAction][ИмяПоля]){
                    Действие[FormAction][ИмяПоля]=[value];
                } else {
                    Действие[FormAction][ИмяПоля].push(value);
                    // Действие[FormAction][ИмяПоля]=Действие[FormAction][ИмяПоля]+","+value;
                }
            } else  if (name.includes(".")){
                console.log(Действие);
                console.log(name, value);
                let ГруппаПолей =  name.split(".",-1);

                console.log("Действие[FormAction]", Действие[FormAction])

                if (!Действие[FormAction][ГруппаПолей[0]]){
                    Действие[FormAction][ГруппаПолей[0]]= {
                        [ГруппаПолей[1]]: value
                    };
                } else {
                    Действие[FormAction][ГруппаПолей[0]][ГруппаПолей[1]] =value
                }
            } else {
                /*алгоритм : если тип значение input равен object то поле является файлом, запускаем чтение файлв в формате DataUrl/base64 и добавляем в Дейсвтие поле файл с обхектом содержащим прочитаный файл*/

                if (typeof value === "object"){
                    // console.log(name, value, Object.keys(value))
                    ФайлыДляЗагрузки.push({
                        Папка:name,
                        Файл:value
                    })



                } else {
                    if (value!==""){
                        Действие[FormAction][name]=value;
                    }

                }

            }
        }

        // let ОжадитьЗагрузкуФайла= true;
        // while (ОжадитьЗагрузкуФайла){
        //     // console.log(ЧтениеФайла.readyState)
        //     if(ЧтениеФайла.readyState===2) {
        //         ОжадитьЗагрузкуФайла = false
        //         console.log('ЧтениеФайла.readyState', ЧтениеФайла.readyState);
        //     }
        //
        // }
        // ЧтениеФайла.onprogress=function(data){
        //     console.log(data)
        // }


        for (let Файл of ФайлыДляЗагрузки){
            let ЧтениеФайла = new FileReader();

            console.log(ЧтениеФайла.result , ЧтениеФайла.readyState)

            ЧтениеФайла.onload = function(e) {
                console.log("onload", Файл)

                Действие[FormAction]["Файл"] = {
                    "Папка":Файл.Папка,
                    "Файл": ЧтениеФайла.result,
                    "ИмяФайла": Файл.Файл.name,
                }
            }
            ЧтениеФайла.onloadend=function() {
                // По завершению загрузки отправляем файл на сервер
                // console.log(ЧтениеФайла.readyState, Действие)
                let mes = new WS.Сообщение({
                    Id: Date.parse(Время()) / 1000,
                    Время: Время(),
                    Кому: "io",
                    От: this.uid,
                    Выполнить: {
                        Действие:Действие,
                    }
                });

                WS.ws.send(mes, true);
            }
            ЧтениеФайла.readAsDataURL(Файл.Файл);
        }
        //     // очистим поле input file чтобы предоотвратить повторную загрузки при загрузке нового файла
            form[name].value=''

    },
    ЗагрузитьФайл: function(event){
        event.preventDefault();
        event.stopPropagation();

        let form = event.target;

        console.log(form)

        let FormAction = event.target.getAttribute("action");
        let Действие={[FormAction]:{}};
        let ДанныеФормы = new FormData(form);
        let ЧтениеФайла = new FileReader();
 console.log(ДанныеФормы);
        for(let [name, value] of ДанныеФормы) {
 console.log(name, value);
            if (name.includes("[]")){

                let ИмяПоля = name.replace("[]","")
                if (!Действие[FormAction][ИмяПоля]){
                    Действие[FormAction][ИмяПоля]=[value];
                } else {
                    Действие[FormAction][ИмяПоля].push(value);
                    // Действие[FormAction][ИмяПоля]=Действие[FormAction][ИмяПоля]+","+value;
                }
            } else  if (name.includes(".")){
                console.log(Действие);
                console.log(name, value);
                let ГруппаПолей =  name.split(".",-1);

                console.log("Действие[FormAction]", Действие[FormAction])

                if (!Действие[FormAction][ГруппаПолей[0]]){
                    Действие[FormAction][ГруппаПолей[0]]= {
                        [ГруппаПолей[1]]: value
                    };
                } else {
                    Действие[FormAction][ГруппаПолей[0]][ГруппаПолей[1]] =value
                }
            } else {
                    /*алгоритм : если тип значение input равен object то поле является файлом, запускаем чтение файлв в формате DataUrl/base64 и добавляем в Дейсвтие поле файл с обхектом содержащим прочитаный файл*/

                if (typeof value === "object"){
                    // console.log(name, value, Object.keys(value))
                    console.log(ЧтениеФайла, value)
                    console.log("onload", value.name)

                    if (!!value.name && value.name != "") {

                        ЧтениеФайла.onload = function (e) {
                            console.log("onload", value.name)
                            Действие[FormAction]["Файл"] = {
                                "Папка": name,
                                "Файл": ЧтениеФайла.result,
                                "ИмяФайла": value.name
                            }
                            // очистим поле input file чтобы предоотвратить повторную загрузки при загрузке нового файла
                            form[name].value = ''
                        }
                        console.log("Имя INPUT ", name, "значение INPUT", value)
                        ЧтениеФайла.readAsDataURL(value);
                    }
                } else {
                    if (value!==""){
                        Действие[FormAction][name]=value;
                    }

                }

            }
        }


        ЧтениеФайла.onloadend=function() {
             // По завершению загрузки отправляем файл на сервер
            console.log(ЧтениеФайла.readyState, Действие)
                let mes = new WS.Сообщение({
                    Id: Date.parse(Время()) / 1000,
                    Время: Время(),
                    Кому: "io",
                    От: this.uid,
                    Выполнить: {
                        Действие:Действие,
                    }
                });

                WS.ws.send(mes, true);
        }


    },
    SaveArticle: function (event){
        event.preventDefault();
        event.stopPropagation();
        let форма = event.target;
        let статья =  форма["статья"];
        console.log(форма);

        толькоТекст = статья.value.replace(/(<([^>]+)>)/ig,"");
        let ДанныеФормы = new FormData(форма);
            console.log("ДанныеФормы",ДанныеФормы);
            ДанныеФормы.append("текст",толькоТекст);

        const действиеФормы = форма.getAttribute("action");
        let АргументыДействия={};
        let ПреобразованиеФормыВДанные=function(v,k){
            if(!!АргументыДействия[k]){
                if (АргументыДействия[k]!== ""){
                    АргументыДействия[k]=АргументыДействия[k]+";"+v;
                }
                // ключ уже есть и пришёл второй ключ, значит нужно передать в качестве массива
                // if (typeof АргументыДействия[k] === "string"){
                //     let ТекущеееЗначение = АргументыДействия[k];
                //     // меняем значение ключа на массив и добавляем туда данные
                //     АргументыДействия[k]=[ТекущеееЗначение,v];
                //      console.log("АргументыДействия string",АргументыДействия[k]);
                //
                // } else if (typeof АргументыДействия[k] === "object"){
                //     АргументыДействия[k].push(v)
                //     console.log("АргументыДействия object",АргументыДействия[k]);
                //
                // }
            } else {
                АргументыДействия[k]=v;
            }
        };


        ДанныеФормы.forEach(ПреобразованиеФормыВДанные)
        let Действие={[действиеФормы]:АргументыДействия};

        let mes = new this.Сообщение({
            Id: Date.parse(Время()) / 1000,
            Время: Время(),
            Кому: "io",
            От: this.uid,
            Выполнить: {
                Действие: Действие,
            }
        });
        this.ws.send(mes,"нелогировать");

    },
    ПоискСтатей: function(event){
         console.log("количество символов",event.target.value.length);
        if (event.target.value.length > 3){
            event.target.form.dispatchEvent(new Event('submit', {cancelable: true}))

        }

    },
/*
* ОтправитьФорму(event, параметры)
* параметры = {
* форма: event.target
* скрыть: true/false - убирает форму со старницы
* сброс: true/false - очищае форму после отправки
* деактивировать: true/false - запрещает ввод в поля после отправки
* нелогировать: любая_строка_отличная_от_null/null или true - не логирует в URL данные отправки формы
* ОбратныйВызов: Функция обратного выхова при получении результата с сревре

* }
*/
    ОтправитьФорму(event, параметры) {
        event.preventDefault();
        event.stopPropagation();

         // console.log(event.submitter.getAttribute('formaction'));

        let ДанныеФормы, FormAction;
        let ОбратныйВызов, НеЛогировать, ДейсвтиеНаКнопке;

         console.log(event);
         if (!!event.submitter) {
              console.log(!!event.submitter);
              console.log(event.submitter);
             ДейсвтиеНаКнопке = event.submitter.getAttribute('formaction')
         }

        console.log("ДейсвтиеНаКнопке ",ДейсвтиеНаКнопке);

        if (!!ДейсвтиеНаКнопке && ДейсвтиеНаКнопке !== ""){
            FormAction = ДейсвтиеНаКнопке
        } else {
            if (!!параметры) {
                if (!!параметры.форма) {
                    FormAction = параметры.форма.getAttribute("action");
                }
            } else {
                FormAction = event.target.getAttribute("action");
            }
        }

        if (!!параметры){
            if (!!параметры.форма) {
                ДанныеФормы = new FormData(параметры.форма);
            }
            if (!!параметры.ОбратныйВызов){
                ОбратныйВызов=параметры.ОбратныйВызов
            }
            if (!!параметры.нелогировать){
                НеЛогировать=параметры.нелогировать
            }
        } else {
            ДанныеФормы = new FormData(event.target);
            НеЛогировать = true;
        }

 console.log("ДанныеФормы", ДанныеФормы);


        let Действие={[FormAction]:{}};
        /*алгоритм: Елси нужно передать массив значений из input с одинаковыми именами, то имя должно содерать [] ? например  <input name="таблицы[]" value="имя_таблицы_1"><input name="таблицы[]" value="имя_таблицы_2">
        * table[] всегда должны быть массивом
        *    <input type="text" class="" name="table[]" value="arests_from_period">
        * */
        let ТребуетСортировки = []

        for(let [name, value] of ДанныеФормы) {
 console.log("name",name);
  console.log(value);
              if (value==="") {
                  continue
              }
             if ((value.includes("{") && value.includes("}")) || (value.includes("[") && value.includes("]"))  ) {
                  console.log(value);
                 value = JSON.parse(value)
             }
            if (name.includes("[") && !name.includes(".")){
                 console.log(`name.includes("[]") && !name.includes(".")`);
                let ИмяПоля = name.replace("[]","")
                 // console.log("ИмяПоля", ИмяПоля);
                let ИдФильтра = "";

                if (ИмяПоля.includes("[")){
                    let МассивПоля_ИдФильтра = ИмяПоля.split(/[\[\]]/);
                    // console.log("МассивПоля_Индекс", МассивПоля_ИдФильтра);
                    ИмяПоля= МассивПоля_ИдФильтра[0]
                    ИдФильтра = МассивПоля_ИдФильтра[1]

                    // console.log("ИдФильтра", ИдФильтра,"value", value);
                }


                if (!Действие[FormAction][ИмяПоля]){

                    if (ИдФильтра === "" || !ИдФильтра) {
                        // console.log("ИдФильтра", ИдФильтра);
                        Действие[FormAction][ИмяПоля] = [value];
                    } else {
                        // console.log(`Действие[${FormAction}][${ИмяПоля}]`, Действие);
                        Действие[FormAction][ИмяПоля]= {[ИдФильтра]:[value]};
                        // console.log(`Действие[${FormAction}][${ИмяПоля}]`, Действие);
                    }

                } else {
                    if (ИдФильтра === "" || !ИдФильтра) {
                        // console.log("Действие ",  Действие[FormAction]);
                        Действие[FormAction][ИмяПоля].push(value);
                    } else {

                        if (!Действие[FormAction][ИмяПоля][ИдФильтра]) {
                            Действие[FormAction][ИмяПоля][ИдФильтра]=[];
                        }
                        // console.log(`Действие[${FormAction}][${ИмяПоля}][${ИдФильтра}] `,  Действие[FormAction][ИмяПоля]);
                        Действие[FormAction][ИмяПоля][ИдФильтра].push(value);
                    }
                    // Действие[FormAction][ИмяПоля].push(value);
                }
                 // console.log("Действие ", Действие);

            } else  if (name.includes(".") && !name.includes("[")){
                console.log("name.includes(\".\") && !name.includes(\"[]\")");
                // console.log(name, value);
                let ГруппаПолей =  name.split(".",-1);
                // console.log("Действие[FormAction]", Действие[FormAction])
                 console.log(value);
                if (!Действие[FormAction][ГруппаПолей[0]]){
                    Действие[FormAction][ГруппаПолей[0]]= {
                        [ГруппаПолей[1]]: value
                    };
                } else {
                    Действие[FormAction][ГруппаПолей[0]][ГруппаПолей[1]] =value
                }

            } else if (name.includes(".") && name.includes("[")){
                 // console.log("name 2", name);
                // имямассива[0].ИмяКлюча - в [] указываеться ключ группы для полей котрые будут входжить в один пакет массива
                // ключ удаляется. и на сервер не передаётся
                let КлючГруппыПолей;
                let ГруппаПолей =  name.split(".",-1);
                    // console.log("ГруппаПолей 2", ГруппаПолей);

                 if (ГруппаПолей[0].includes("[")){
                     let ИмяМассива = ГруппаПолей[0].split(/[\[\]]/);
                      console.log("ИмяМассива", ИмяМассива);
                      console.log("ИмяМассива[1]", ИмяМассива[1]);
                      console.log("ГруппаПолей", ГруппаПолей);
                     КлючГруппыПолей=ИмяМассива[1]

                     if (ГруппаПолей[1].includes("[")){
                         var ИмяВложенногоОбъетка = ГруппаПолей[1].split(/[\[\]]/)[0];
                     }


                     if (!Действие[FormAction][ИмяМассива[0]]){

                         if (ИмяВложенногоОбъетка !== "" || !!ИмяВложенногоОбъетка) {
                             Действие[FormAction][ИмяМассива[0]]= {}
                         } else {
                             Действие[FormAction][ИмяМассива[0]]=[]
                             ТребуетСортировки.push(ИмяМассива[0])
                         }
                         // ТребуетСортировки.push(Действие[FormAction][ИмяМассива[0]])

                     }


                     console.log(ИмяВложенногоОбъетка);

                     if (!Действие[FormAction][ИмяМассива[0]][ИмяМассива[1]]){

                       if ( !!ИмяВложенногоОбъетка && ИмяВложенногоОбъетка !== "") {

                             console.log("ИмяМассива[0]", ИмяМассива[0], "ИмяМассива[1]", ИмяМассива[1]);

                             Действие[FormAction][ИмяМассива[0]][ИмяМассива[1].toString()][ИмяВложенногоОбъетка]=[]

                              console.log(Действие[FormAction][ИмяМассива[0]]);
                         } else {
                             Действие[FormAction][ИмяМассива[0]][ИмяМассива[1]]={[ГруппаПолей[1]]: value}
                         }

                     } else {
                         if ( !!ИмяВложенногоОбъетка && ИмяВложенногоОбъетка !== "") {
                             if (!Действие[FormAction][ИмяМассива[0]][ИмяМассива[1]][ИмяВложенногоОбъетка]) {
                                 Действие[FormAction][ИмяМассива[0]][ИмяМассива[1].toString()][ИмяВложенногоОбъетка]=[]
                             }
                         }
                     }
                      console.log(Действие);
                     if ( !!ИмяВложенногоОбъетка && ИмяВложенногоОбъетка !== "")  {
                         console.log("ИмяВложенногоОбъетка", ИмяВложенногоОбъетка);
                         console.log("Действие[FormAction][ИмяМассива[0]]", Действие[FormAction][ИмяМассива[0]]);
                         Действие[FormAction][ИмяМассива[0]][ИмяМассива[1].toString()][ИмяВложенногоОбъетка].push(value)
                         ИмяВложенногоОбъетка=null
                     } else {
                         Действие[FormAction][ИмяМассива[0]][ИмяМассива[1]][ГруппаПолей[1]] = value
                     }
                      // console.log("name 3", name);
                      // console.log("Действие[FormAction] ", Действие[FormAction]);
                     // Действие[FormAction][ИмяМассива].push(ГруппаПолей[0])
                 }

            }

            else {
                // console.log(`Действие[${FormAction}][${name}] `,  Действие[FormAction][name]);
                Действие[FormAction][name]=value;

            }
        }
        // Действие[FormAction][ИмяМассива[0]][ИмяМассива[1]][ГруппаПолей[1]].flat()
        Действие[FormAction]["время"]=Время();
        // for (let СортируемыйОбъек in Действие[FormAction]){
        //      console.log("СортируемыйОбъек", СортируемыйОбъек);
            for (let Обект of ТребуетСортировки) {
                Действие[FormAction][Обект]= Действие[FormAction][Обект].flat()
            }
        // }
        console.log("ОтправитьФорму Действие", Действие);
        let ТекущееВремя = new Date
        let ms = ТекущееВремя.getMilliseconds()
        let mes = new this.Сообщение({
            Id: Date.parse(Время()) / 1000 + ms,
            Время: Время(),
            Кому: "io",
            От: this.uid,
            Выполнить: {
                Действие: Действие,
            }
        });

        if(!!ОбратныйВызов){
            mes.ОбратныйВызов = ОбратныйВызов;
        }

        this.ws.send(mes, НеЛогировать);



        if (!!параметры) {
            if (параметры.деактивировать) {
                for (let inputField of параметры.форма.elements) {
                    inputField.disabled = true;
                    inputField.classList.add("disabled")
                }
            }
            if (параметры.скрыть) {
                параметры.форма.parentNode.removeChild(параметры.форма);
            }
            if (параметры.сброс) {
                параметры.форма.reset()
            }
        }
    },

    SendForm: function (event, убратьФорму, reset=false, деактивировать = false ) {
        event.preventDefault();
        event.stopPropagation();
        // console.log(event)




        let form = event.target;

        let formData = new FormData(form);

        const action = form.getAttribute("action");
        let Действие={[action]:{}};


        for(let [name, value] of formData) {

            Действие[action][name]=value;
        }

        let mes = new this.Сообщение({
            Id: Date.parse(Время()) / 1000,
            Время: Время(),
            Кому: "io",
            От: this.uid,
            Выполнить: {
                Действие: Действие,
            }
        });
        console.log("SendForm на сервер", mes)

        this.ws.send(mes);
        if (деактивировать){
            for(let inputField of form.elements){

                // if (inputField.attributes.id !== ""){
                //     inputField.setAttribute("id","");
                inputField.disabled=true;
                inputField.classList.add("disabled")
                // }
                // if (inputField.nodeName === "BUTTON"){
                //     inputField.classList.add("disabled")
                // }
            }
        }
        if (убратьФорму){
            form.parentNode.removeChild(form);
            for(let inputField of form.elements){

                    inputField.disabled=true;
                    inputField.classList.add("disabled")

            }
        }
        if (reset){
            form.reset()
        }
    },
    ЛогинНеНайден: function(сообщение){
        console.log("Логин не найден, но вероятно его нет в локалной базе данных, поэттому ждём пароль и проверяяем в тма лдап")
        сообщение.Контэнт.данные = "Логин не найден";
        this.FloatMessage(сообщение)
        //
        // let user_fullname_wrapper = document.getElementById("user_fullname_wrapper")
        //     user_fullname_wrapper.classList.remove("hidden")
        //     user_fullname_wrapper.classList.add("error")
        // let user_fullname = document.getElementById("user_fullname");
        //     user_fullname.innerText="Логин не найден!!";
        //     user_fullname.classList.add("error")
    },
    ПоказатьФИОВФормеАтворизации: function(сообщение){
        console.log(сообщение)
        if (!!сообщение.Контэнт.данные) {
            let user_fullname_wrapper = document.getElementById("user_fullname_wrapper");
                user_fullname_wrapper.classList.remove("hidden", "error");
                // user_fullname_wrapper.classList.remove("error");
            let user_fullname = document.getElementById("user_fullname")
                user_fullname.innerText=сообщение.Контэнт.данные;
                user_fullname.classList.remove("error")
        }
    },

    printMessage: function (сообщение, type = "mes_from", blockLog = null) {
        // console.log("368 printMessage ",сообщение);
        // console.log("blockLog",blockLog);
        // console.log(" this.UserInfo", this.UserInfo);
        // console.log(" WS.UserInfo", WS.UserInfo);
        let Reaply = "";

        if (!!сообщение.MessageType && сообщение.MessageType.includes("server_log")){
            blockLog = document.getElementById("server_log");
        } else if (!!сообщение.MessageType && сообщение.MessageType.includes("attention")){
            blockLog = document.getElementById("float_doc");
        } else if(!!сообщение.MessageType && (сообщение.MessageType.includes("error") || сообщение.MessageType.includes("info"))){
            this.FloatMessage(сообщение);
            return
        }else {
            if (!blockLog) {
                // blockLog = this.chatLog;
                // if (!blockLog) {
                blockLog = document.getElementById("log_io");
                // }
            }
        }
        if (!blockLog){
            return
        }


        if (сообщение.ip !== "") {
            console.log(сообщение);
        }
        if (сообщение.Текст === "") {
            if (!!сообщение.Content["target"] && сообщение.Content["target"] === "log_"+сообщение.От){

            } else if(!!сообщение.Контэнт.html){

            } else {
                console.log("Нечего писать в сообщения")
                return this;
            }

        }



        if (!this.chatLog) {
            this.chatLog = document.getElementById("ws_chat_log");
        }
        if (!this.LastChats) {
            this.LastChats = document.getElementById("ws_last_chats");
        }
        if (!!сообщение.ОтветНа && сообщение.ОтветНа !== "") {
            Reaply = `<div class= "reaply_to">${сообщение.ОтветНа}</div>`
        }

        let messageHTML = this.tpl("Сообщение", сообщение);

        if (!!сообщение.MessageType && сообщение.MessageType.includes("attention")) {
            blockLog.innerHTML = messageHTML;
            // blockLog.insertAdjacentHTML("beforeend", messageHTML);
            // blockLog.scrollTop = blockLog.scrollHeight;

        } else {
            blockLog.insertAdjacentHTML("beforeend", messageHTML);
        }
 // console.log(blockLog.scrollTop);
 // console.log( blockLog.parentNode.scrollHeight);
 // console.log( blockLog.parentNode.scrollHeight);
        blockLog.parentNode.scrollTop = blockLog.parentNode.scrollHeight;
        return this;
    },
    terminalOpen: function (login) {
        let terminalWrapper = document.getElementById(`ws_terminal_${login}`);
        terminalWrapper.classList.remove("offline");
        terminalWrapper.classList.remove("mini");
        let m = new WS.Сообщение({
            Выполнить: {"action": "SSHConnect", "arg": {"Login": login}},
            Время: PgTimeStamp(),
            Кому: "io",
        });

        // WS.printMessage(this.Сообщение, "outcoming" );
        this.ws.send(m);
        return this

    },
    terminalClose:function(login){
        let terminalWrapper = document.getElementById(`ws_terminal_${login}`);
        terminalWrapper.classList.add("offline");
        terminalWrapper.classList.add("mini");
        this.CloseSSHConnection(login);
        return this
    },
    CloseSSHConnection: function (login) {
        let m = new WS.Сообщение({
            Выполнить: {"action": "CloseSSH", "arg": {"Login": login}},
            Время: PgTimeStamp(),
            Кому: "io",
        });
        // WS.printMessage(this.Сообщение, "outcoming" );
        this.ws.send(m);
        return this
    },
    tpl: function (tplName, values) {
        return tpl(tplName, values)
    },

    Auth: function () {

    },
    ОткрытьЗакрыть: function(targetId){
        let target = document.getElementById(targetId);
        target.classList.toggle("hidden")
    },
    ПоказатьДетализацию: function (event) {

        // В таблице в заголовке должна быть строка с нумерацией столбцов, в ней и должны быть атрибуты sql
        console.log(event);
        let ТекущаяЯчейка = event.originalTarget;
        if (ТекущаяЯчейка.tagName!=="TD"){
            console.log("не ячейка")
            return
        }

        if (ТекущаяЯчейка.classList.contains('nodetails')){

            return
        }

        let СтрокаЯчейки = ТекущаяЯчейка.parentElement

        // Посчитаем количество строк в заголовке, найдём объединённые по вертикали, если такие есть и они захватывают строку details_row то прибавим

        let Аргументы ={};
        // ищем ИД строки в которой находится ячейка, для дальнейшего получения ИД sql для сбора детализации
            if (СтрокаЯчейки.tagName==="TR") {

               // определение аргументов по id нужно удалить. использовалось в витрине  arests_from_period
               //  let ИД = СтрокаЯчейки.id
               //  let МассивИД = ИД.split(".");
               //  let ИмяАргумента = МассивИД[0]
               //  let КодОСП = null
               //  if(ИмяАргумента==="OSP"){
               //      КодОСП = МассивИД[1];
               //  }
               //  Аргументы[ИмяАргумента]=МассивИД[1];

                if (!!СтрокаЯчейки.dataset ) {
                    for (let ИмяАргумента in СтрокаЯчейки.dataset){
                        console.log(ИмяАргумента, СтрокаЯчейки.dataset[ИмяАргумента])
                        Аргументы[ИмяАргумента]=СтрокаЯчейки.dataset[ИмяАргумента];
                    }
                }

                let ТелоТаблицы = СтрокаЯчейки.parentElement
                let Таблица = ТелоТаблицы.parentElement
                // проверим если ID таблицы имеет префикс OSP то возьмум код осп отдута
                // if (КодОСП === null){
                //     let МассивИДТаблицы = Таблица.id.split(".")
                //     if (МассивИДТаблицы[0]==="OSP") {
                //         КодОСП = МассивИДТаблицы[1];
                //     }
                //     Аргументы[КодОСП]=МассивИДТаблицы[1];
                // }

                if (!!Таблица.dataset) {
                    for (let ИмяАргумента in Таблица.dataset){
                        console.log(ИмяАргумента, Таблица.dataset[ИмяАргумента])
                        Аргументы[ИмяАргумента]=Таблица.dataset[ИмяАргумента];
                    }
                }

                let ЗаголовокТаблица = Таблица.tHead;
                let СтолбецЯчейки=null;
                // let ЗначениеСтолбца=null;

                for (let строкаЗаголовка of ЗаголовокТаблица.rows){
                    if (!!строкаЗаголовка.dataset && !!строкаЗаголовка.dataset.details_arguments && строкаЗаголовка.dataset.details_arguments !== ""){


                        СтолбецЯчейки = строкаЗаголовка.cells[ТекущаяЯчейка.cellIndex]
                        console.log(СтолбецЯчейки)

                        if (!!СтолбецЯчейки.dataset ) {
                            for (let ИмяАргумента in СтолбецЯчейки.dataset){
                                console.log(ИмяАргумента, СтолбецЯчейки.dataset[ИмяАргумента])
                                Аргументы[ИмяАргумента]=СтолбецЯчейки.dataset[ИмяАргумента];
                            }
                        }
                    // else if (!!СтолбецЯчейки.id && СтолбецЯчейки.id !== "") {
                    //         ЗначениеСтолбца = СтолбецЯчейки.id
                    //         Аргументы["столбец"]=ЗначениеСтолбца;
                    //     }

                        console.log("СтолбецЯчейки", СтолбецЯчейки)

                        if (!!СтолбецЯчейки ){
                            break;
                        }
                    }

                }


                if (!!СтолбецЯчейки){

                  let ИДТаблицыЗапроса =  Таблица.id+"-details";

                  // Получаем ИД Столбца Конкатенируем к нему имя таблицы и это значение будет соответсвовать полю table в таблице fssp_data.sql
                  //  а входящиеАргументы будут номером ОСП

                   // проверим если данные собираються  за период. то найдём форму в которой задаёться период и получим из неё значения полуй date_from, date_to
                   let ПериодС=null;
                   let ПериодПо=null;

                    console.log("нужно написать обработчик для сбора периода")

                   if (!!ТекущаяЯчейка.querySelector(".date_from")){
                       ПериодС=ТекущаяЯчейка.querySelector(".date_from").innerText
                       ПериодПо=ТекущаяЯчейка.querySelector(".date_to").innerText
                   } else if (!!Таблица.dataset.period && Таблица.dataset.period !== ""){
                       console.log(Таблица.dataset.period)
                        let ФормаПериода = document.getElementById(Таблица.dataset.period);
                        ПериодС = ФормаПериода.querySelector("#date_from").value
                        ПериодПо= ФормаПериода.querySelector("#date_to").value
                    }
                    Аргументы["date_from"]=ПериодС
                    Аргументы["date_to"]=ПериодПо
                    Аргументы["table[]"]=[ИДТаблицыЗапроса]
                    Аргументы["id_таблицы"]=ИДТаблицыЗапроса // используйется для
                    Аргументы["КонтейнерРезультата"]="modal"
                    if (!!Таблица.dataset.action && Таблица.dataset.action !== "") {
                        WS.ОтправитьСообщениеИО(Таблица.dataset.action,Аргументы, true)
                    } else {
                        WS.ОтправитьСообщениеИО("Детализация",Аргументы, true)
                    }


                } else {
                     console.log("не удалось определить столбец");
                }

            } else {
                 console.log("Ищем строку");
            }

    },
};



/*
* Данные = {
* ИдМодала: 'id',
* Заголовок: 'заголовок или ничего',
* Контэнт: html
* Кнопки: ['закрыть']
* }
* */
let МодальноеОкно = function (ВходящиеДанные){
    this.ИдМодала = ВходящиеДанные.ИдМодала ? ВходящиеДанные.ИдМодала : null;
    if (!this.ИдМодала){
         console.log("Нет ИД для модального окна");
         return;
    }
    this.Заголовок = ВходящиеДанные.Заголовок ? ВходящиеДанные.Заголовок : null;
    this.Контэнт = ВходящиеДанные.Контэнт ? ВходящиеДанные.Контэнт : null;
    if (!this.Контэнт){
        console.log("Нет Контэнта для модального окна");
        return;
    }
    this.Кнопки =  ВходящиеДанные.Кнопки ? ВходящиеДанные.Кнопки : ['закрыть'];
    this.Класс =  ВходящиеДанные.Класс ? ВходящиеДанные.Класс.push('modal').join(',') : 'modal';
     console.log(this.Класс);

    let Кнопки="";
    for (let ИмяКнопки of this.Кнопки){
      let ШаблоныКнопок = {
          закрыть: `<button type="button" class="close" onclick="Закрыть('${this.ИдМодала}')"></button>`
      }
        Кнопки += ШаблоныКнопок[ИмяКнопки]
    }
    this.ИдКонтэнта = ВходящиеДанные.ИдКонтэнта ? ВходящиеДанные.ИдКонтэнта : "modal_content";
    let ШаблонМодала = `<div class="${this.Класс}" id = "${this.ИдМодала}"><div id="" class="modal_header">
                        ${Кнопки}
                        </div>
                        <div class="modal_content" id="${this.ИдКонтэнта}">
                         ${this.Контэнт}
                        </div>
                    </div>`;
    let Модал = document.createRange().createContextualFragment(ШаблонМодала)
    console.log(Модал);
    let MainContent = document.getElementById("main_content")
    let modalContent
    if (!!MainContent){
        modalContent =MainContent.append(Модал);
    } else {
        modalContent = document.body.append(Модал);
    }
    document.body.classList.add('noscroll')

    return modalContent;
};


function CloseModal(modal_id){
    let modal = document.getElementById(modal_id);
        modal.classList.add("hided"); // анимация сворачивания окна
        modal.parentNode.removeChild(modal)
}
function Закрыть(id){
    let modal = document.getElementById(id);
        modal.classList.add("hided"); // анимация сворачивания окна
        modal.parentNode.removeChild(modal)
    document.body.classList.remove("noscroll")
}

/**
 * @return {boolean}
 */
function ScanKeys(event){
    // console.log(event)
    if (event.keyCode === 13 && event.ctrlKey) {
        // console.log(event.target);
        // br = document.createElement('br');
        // return event.target.appendChild(br);
        return false
    } else if (event.keyCode === 13){
        event.preventDefault();
        event.stopPropagation();
        WS.sendMessage(event.target)
        return true
    }
    // console.log(event.target.innerHTML);
    // console.log(event.target.innerText);
    return false
}
function sendLogError(errorObject){
    console.log("SendLogError"," ФУНКЦИЯ ОТПРАВКИ ОШИБКИ НА СЕРВЕР НЕ РЕАЛИЗОВАНА");
    console.log(errorObject);
}

function ReloadFiles(Args){
   for (let fileId in Args.other){
       // console.log(fileId);
       // console.log("window.location.host", window.location.host);
       ReloadedFile = document.getElementById(fileId);
       host="";
       if (window.location.host !== "http://10.26.6.25"){
           host = "http://10.26.6.25"
       }
       if (fileId === "wsstyle"){
           ReloadedFile.setAttribute("href", host+Args.other[fileId]+"?"+Date.now())
       }
       // else if (fileId === "wsjs"){
       //     ReloadedFile.setAttribute("src", host+Args.other[fileId]+"?"+Date.now())
       // }
   }

}

function MinimizeTerminal(login){
    let terminal = document.getElementById("ws_terminal_"+login);
    terminal.classList.toggle('mini');

}


function CloseTerminal(login){
    let terminal = document.getElementById("ws_terminal_"+login);
    let target= document.getElementById("ws_chat_log");
    let user_log_wrapper= document.getElementById("log_wrapper_"+login);

        target.appendChild(terminal);
        terminal.classList.remove("float");
        terminal.querySelector(".ws_terminal_header").classList.remove("draggable");
        terminal.style=null;
        terminal.setAttribute('data-x', 0);
        terminal.setAttribute('data-y',  0);

        terminal.classList.add("mini","offline");

    if (user_log_wrapper.classList.contains("hiden")){
        terminal.classList.toggle("hiden")
    }
    let ticon = document.getElementById('ticon_'+login);
    if (!!ticon){
        ticon.parentNode.removeChild(ticon);
    }
    WS.CloseSSHConnection(login)
}

function AddTerminalToDoc(login){
    let doc_container = document.getElementById("doc_container");
    let ticon = document.getElementById(`ticon_${login}`);

    if (!doc_container.contains(ticon)){
        let terminalIcon = document.createElement("div");
        terminalIcon.classList.add('terminal_icon');
        terminalIcon.setAttribute('id', 'ticon_'+login);
        terminalIcon.innerHTML='<button class="tocon_name" onclick="activateTerminal(${login}})">' +login + '</button>';
        doc_container.append(terminalIcon)
    }

}

function activateTerminal(login){
    console.log("Подсветить терминал ", login)
}

function FloatTerminal(login){

   const position = { x: 0, y: 0 };

   let termanals_container = document.getElementById("termanals_container");
   let terminal = document.getElementById("ws_terminal_"+login);

   if (termanals_container.contains(terminal)){
       terminal.classList.toggle("float");
       let target= document.getElementById("ws_chat_log");
           target.appendChild(terminal);
       return
   }
        terminal.classList.toggle("float");
        terminal.querySelector(".ws_terminal_header").classList.toggle("draggable");

    termanals_container.appendChild(terminal);

    AddTerminalToDoc(login);

    let terminals_draggable =  interact('.draggable');
    let terminals_float =  interact('.float');
    let  target;

    terminals_draggable.draggable({
        modifiers: [
            interact.modifiers.restrict({
                restriction: document.body,
                // restriction: terminal.parentNode,
                // elementRect: { top: 0, left: 0, bottom: 1, right: 1 },
                endOnly: true
            })
        ],
        listeners: {
            start (event) {
                target = event.target.parentNode;
                position.x = (parseFloat(target.getAttribute('data-x')) || 0);
                position.y = (parseFloat(target.getAttribute('data-y')) || 0);
                console.log(target, position)
                // positon = JSON.parse(event.target.dataset.position)
            },
            move (event) {
                position.x += event.dx;
                position.y += event.dy;
                // event.target.dataset.position=JSON.stringify({x: position.x,y: position.y});
                target.style.transform = `translate(${position.x}px, ${position.y}px)`;
                target.setAttribute('data-x', position.x);
                target.setAttribute('data-y',  position.y);
            },


        },
        inertia: true,
    });
    terminals_float.resizable({
        // resize from all edges and corners
        edges: { left: false, right: true, bottom: true, top: false },

        modifiers: [
            // keep the edges inside the parent
            // interact.modifiers.restrictEdges({
            //     // outer: 'parent',
            //     endOnly: true
            // }),

            // minimum size
            interact.modifiers.restrictSize({
                min: { width: 300, height: 180 }
            })
        ],
        // inertia: true
    }).on('resizemove', function (event) {

       target = event.target;
       position.x = (parseFloat(target.getAttribute('data-x')) || 0);
       position.y = (parseFloat(target.getAttribute('data-y')) || 0);

            // update the element's style
            target.style.width = event.rect.width + 'px';
            target.style.height = event.rect.height + 'px';

            // translate when resizing from top or left edges
               position.x += event.deltaRect.left;
               position.y  += event.deltaRect.bottom;

            target.style.webkitTransform = target.style.transform ='translate(' +  position.x + 'px,' + position.y + 'px)';

            target.setAttribute('data-x', position.x);
            target.setAttribute('data-y', position.y);

        })

}

function Время(){
    return PgTimeStamp()
}
function PgTimeStamp(){

    const monthArray = [
        "01",
        "02",
        "03",
        "04",
        "05",
        "06",
        "07",
        "08",
        "09",
        "10",
        "11",
        "12"];

    let Data = new Date();
    let Year = Data.getFullYear();
    let Month = monthArray[Data.getMonth()];

    let Day = Data.getDate();
    let Hour = Data.getHours();
    let Minutes = Data.getMinutes();
    let Seconds = Data.getSeconds();
    let Milliseconds = Data.getMilliseconds();
    // 2019-10-02 12:36:28.421287

    timestamp = Year+"-"+Month+"-"+Day+" "+Hour+":"+Minutes+":"+Seconds//+":"+Milliseconds;
    // Date.parse(timestamp)/1000;
    // console.log(Data.getMonth());
    console.log("timestamp", timestamp);
    return timestamp
}

function divBlur(event){
    console.log(event)
}

function DeleteForm(event) {
    let form=event.target;
    form.parentNode.classList.add("hidden");
    setTimeout(function(){
        // form.parentNode.removeChild(form)
    },3000)

}
function UserChange() {
let user_fullname_wrapper = document.getElementById("user_fullname_wrapper");
    user_fullname_wrapper.classList.toggle("hidden")
}

function СкрытьФормуВхода(){
    let auth_block = document.getElementById("auth_block");
    if (!!auth_block){
        auth_block.classList.add("auth");
        setTimeout(function(){
            auth_block.parentNode.removeChild(auth_block)
        }, 2000)
    }
}

function СохранитьСтатью (Editor){
    console.log("СохранитьСтатью Editor", Editor);
    let ctx = Editor.getContext();
     console.log("ctx",ctx);
    let РодитеслькаяФорма =  ctx.element.originElement.form;
    const ОригинальныйAction = РодитеслькаяФорма.getAttribute("action");;
         console.log("ОригинальныйAction", ОригинальныйAction);
         РодитеслькаяФорма.action = "сохранить черновик статьи в базу знаний";

    let времяИзменения =  РодитеслькаяФорма.elements["время_изменения"];
        времяИзменения.value=Время();
        // времяИзменения.setAttribute("value",Время());
         console.log("времяИзменения",времяИзменения);

    // let ид_статьи = РодитеслькаяФорма.elements["ид_статьи"];
    Editor.save();
    РодитеслькаяФорма.dispatchEvent(new Event('submit', {cancelable: true}));
    РодитеслькаяФорма.action=ОригинальныйAction
}


let TreeMenu = {
    Переключить: function (цель) {
         // console.log(event.target);
        document.getElementById(цель).classList.toggle("open")
    },
};


let Tabs = {
     Активная : null,
     Контейнер: null,
     Заглушка: null, // вкладка заглушка если нет данных или вкладки, если вкладки нет то данные не пришли с сервера
     ShowTab: function(контейнер, таб) {
         // вначале проверим есть ли контейнер в памяти, если нет то найдём его и запомним, так же найдём и запомним активную вкладку

        // if (!this.Контейнер){
        this.Контейнер = document.getElementById(контейнер);
        // }
        //  console.log("Контейнер", this.Контейнер.innerHTML)
        //  console.log("Контейнер", this.Контейнер)
         if (!this.Заглушка){
             if (!!this.Контейнер ) {
                 this.Заглушка =  this.Контейнер.querySelector('.empty_tab')
                 console.log("Заглушка",this.Заглушка)
             } else {
                 this.Заглушка = document.getElementById('empty_tab')
                 console.log("Заглушка",this.Заглушка)
             }
         }
         // console.log("this.Заглушка", this.Заглушка)
         console.log(" Активная вкладка ", this.Активная)


        if (!this.Активная){
            // Если в памяти нет активной вкладки, пройдём по всем вкладкам конетйнера и найдём ту у которой нет класска hidden она и будет активной
            let AllTabs = this.Контейнер.querySelectorAll('.tab');
            let BufferActiveTabs=[];

            console.log("Прячем все вкладки", AllTabs);
            for (let Tab of AllTabs) {
                if (!Tab.classList.contains('hidden')){
                   BufferActiveTabs.push(Tab)
                }
            }
            if (BufferActiveTabs.length >1) {
                for (let ActiveTab of BufferActiveTabs){
                    console.log("ActiveTab", ActiveTab)
                }
            } else if (BufferActiveTabs.length === 1){
                this.Активная = BufferActiveTabs[0]
            } else {
                console.log("нет активных влкдаок")

            }
        }
         let АктивируемаяВкладка
        // теперь переключим активную вкладку
         if (таб === "all"){
             let AllTabs = this.Контейнер.querySelectorAll('.tab');
             if (this.Активная==='all') {
                 for (let Tab of AllTabs) {
                     console.log(Tab.id)
                     if (Tab.id !== "empty_tab") {
                         Tab.classList.add('hidden')
                     }
                 }
                 this.Активная=null
             } else {


                 for (let Tab of AllTabs) {
                     console.log(Tab.id)
                     if (Tab.id !== "empty_tab") {
                         Tab.classList.remove('hidden')
                     }
                 }
                 this.Активная='all'
             }

         } else {
             АктивируемаяВкладка = this.Контейнер.querySelector('#'+таб);

             console.log("АктивируемаяВкладка", АктивируемаяВкладка)


                if (!!АктивируемаяВкладка) {
                    if (!!this.Активная) {
                        if (this.Активная === 'all') {
                            let AllTabs = this.Контейнер.querySelectorAll('.tab');

                            for (let Tab of AllTabs) {
                               Tab.classList.add('hidden')
                            }
                            if (this.Активная.id !== таб) {
                                АктивируемаяВкладка.classList.remove('hidden')
                                this.Активная = АктивируемаяВкладка
                            } else {
                                this.Активная = this.Заглушка;
                                if (!!this.Заглушка){
                                    this.Заглушка.classList.remove('hidden')
                                }
                            }

                        } else {

                            if (this.Контейнер.id !== this.Активная.id){
                                this.Активная.classList.add('hidden')
                                if (this.Активная.id !== таб) {
                                    АктивируемаяВкладка.classList.remove('hidden')
                                    this.Активная = АктивируемаяВкладка
                                } else {
                                    this.Активная = this.Заглушка;
                                    if (!!this.Заглушка){
                                        this.Заглушка.classList.remove('hidden')
                                    }
                                }
                            } else {
                                АктивируемаяВкладка.classList.remove('hidden')
                                this.Активная = АктивируемаяВкладка
                            }


                            // this.Заглушка.classList.add('hidden')
                        }
                    } else {
                        АктивируемаяВкладка.classList.remove('hidden')
                        this.Активная = АктивируемаяВкладка
                    }
                    // if (this.Активная.id !== таб) {

                    // }

                } else {
                    console.log("АктивируемаяВкладка", АктивируемаяВкладка)
                    console.log("Запросить данные с сервера и добавить вкладку", АктивируемаяВкладка)
                    if (!!this.Активная) {
                        if (this.Активная === 'all') {
                            let AllTabs = this.Контейнер.querySelectorAll('.tab');
                            for (let Tab of AllTabs) {
                                Tab.classList.add('hidden')
                            }
                        } else {
                            this.Активная.classList.add('hidden')
                        }
                    }

                    // if (!!this.Заглушка){
                    //     this.Заглушка.classList.remove('hidden');
                    //     this.Активная = this.Заглушка
                    // }
                }
         }
         // this.Контейнер.querySelector('#'+таб).classList.remove('hide')
    },
};






function  openFile(event){
    // let files = event.target;
    //
    // let reader = new FileReader();
    // reader.onload = function(){
    //     var dataURL = reader.result;
    //     var output = document.getElementById('output');
    //     console.log('dataURL',dataURL);
    //     output.src = dataURL;
    // };
    //
    // fetch
    //
    // reader.readAsDataURL(files.files[0]);
    var file = event.target.files[0];
    var reader = new FileReader();
    var rawData = new ArrayBuffer();

    reader.loadend = function() {

    }

    reader.onload = function(e) {

        rawData = e.target.result;
        console.log(rawData)
        // ws.send(rawData);

    }

    reader.readAsArrayBuffer(file);
}

function ПарсингМеток(event, ОберткаМеток){ //, inputМассивВсехДанных - более не ниспользуется
     if (event.keyCode === 59 || event.keyCode === 13){
         event.preventDefault();
         event.stopPropagation();
         console.log(event);
         var метка = event.target.value.replace(/(\D+)/g, '$1');
         Обертка = document.getElementById(ОберткаМеток)


         if (!!Обертка.children[0] && Обертка.children[0].classList.contains("tags_placeholder")){
             Обертка.innerHTML = `<!--<div class="hashtag">${метка}<button class="remove_tag">x</button></div>-->`

             Обертка.innerHTML = tpl("метка", метка)
         } else {
             Обертка.insertAdjacentHTML("beforeend",  tpl("метка", метка))
         }

         event.target.value="";
         // input в котором накапливаем все метки через запятую, более не нужен
         // let tags_input = document.getElementById(inputМассивВсехДанных)
         //
         // if (tags_input.value==="") {
         //     tags_input.value=метка
         // } else{
         //     tags_input.value=tags_input.value+","+метка
         // }
         // let tags = tags_input.value
         // let tagsArray = [];
         // if (tags!=="") {
         //     tagsArray  = JSON.parse(tags)
         // }
         // tagsArray.push(метка)
         // let tagsString = JSON.stringify(tagsArray)
         // tags_input.value= tagsString
     }

}

function monacoValue(id){
   monaco = document.getElementById(id)
     console.log(monaco);
     console.log(monaco.getValue());

}
function СвернутьПанель(id, event ){
     console.log(id);
   let панель = document.getElementById(id)
    панель.classList.toggle("mini")

    if (!!event.target.dataset.button_name && event.target.dataset.button_name !==""){
        let НовоеИмя = event.target.dataset.button_name
            event.target.dataset.button_name =  event.target.innerText
            event.target.innerText=НовоеИмя
    }
}

/*
* ВставитьДанныеВТаблицу если данные пришлли из ОСП
* */
function ВставитьДанныеВТаблицу(контейнер, сообщение){

    if (!!сообщение.Контэнт.данные) {

        for (let ИдЯчейки in сообщение.Контэнт.данные) {

            let ДанныеЯчейки = сообщение.Контэнт.данные[ИдЯчейки]

            let КоординатыЯчейки = ИдЯчейки.split(".");
            let ИДСтроки = КоординатыЯчейки[0];
            let ИДСтолбца = КоординатыЯчейки[1];
            let Строка;

            if (ИДСтроки === "osp") { // если ИД строки равно osp тогда ИДСтроки равен коду осп
                ИДСтроки = ДанныеЯчейки.ОСП.osp_code;
                Строка = контейнер.rows[`OSP.${ИДСтроки}`]
            } else {
                Строка = контейнер.rows[`${ИДСтроки}`]
            }
             console.log("контейнер",контейнер);


            let НомерЯчейки;
            for (let строкаЗаголовка of контейнер.tHead.rows){
                console.log("ИДСтолбца", ИДСтолбца);
                НомерЯчейки = строкаЗаголовка.cells[ИДСтолбца].cellIndex;
                console.log("НомерЯчейки",НомерЯчейки);
                if (!!НомерЯчейки ){
                    break;
                }
            }
            console.log("Строка 1",Строка,"ИДСтроки ",ИДСтроки);
            if (!!Строка) {
                console.log("Строка 2",Строка,"ИДСтроки ",ИДСтроки);
                console.log("ДанныеЯчейки.Количество", ДанныеЯчейки.Количество);

                console.log(" Строка.cells[НомерЯчейки]", Строка.cells[НомерЯчейки]);
                if(ДанныеЯчейки.Количество !== null && ДанныеЯчейки.Количество !== undefined){
                     console.log(" Строка.cells[НомерЯчейки]", Строка.cells[НомерЯчейки]);
                    Строка.cells[НомерЯчейки].innerHTML = ДанныеЯчейки.Количество;
                // if (!!ДанныеЯчейки.КартаСтрок && ДанныеЯчейки.КартаСтрок.length >0 ) {
                //     Строка.cells[НомерЯчейки].insertAdjacentElement("beforeend", СоздатьДетализацию(ДанныеЯчейки))}
                }
                if (!!ДанныеЯчейки.Ошибка && ДанныеЯчейки.Ошибка !==""){
                    Строка.cells[НомерЯчейки].insertAdjacentText("beforeend", ДанныеЯчейки.Ошибка);
                }

            }

        }
    }
}

function ВставитьДанныеВТаблицуИзРБД(контейнер, сообщение){
    console.log(контейнер)
    if (!!сообщение.Контэнт.данные) {
        console.log("данные", сообщение.Контэнт.данные)
        for (let ИдЯчейки in сообщение.Контэнт.данные) {

            let ДанныеЯчейки = сообщение.Контэнт.данные[ИдЯчейки]

            let КоординатыЯчейки = ИдЯчейки.split(".");
            let ИДСтроки = КоординатыЯчейки[0];
            let ИДСтолбца = КоординатыЯчейки[1];
            let СтрокаТаблицы;
            let НомерЯчейки;
            for (let НомерСтрокиЗаголовка in контейнер.tHead.rows){

                let строкаЗаголовка = контейнер.tHead.rows[НомерСтрокиЗаголовка]


                if (!!строкаЗаголовка.cells[ИДСтолбца]){
                    НомерЯчейки = строкаЗаголовка.cells[ИДСтолбца].cellIndex;

                    if (!!НомерЯчейки ){
                        break;
                    }
                }

            }


           if(!!ДанныеЯчейки.КартаСтрок) {
                for (let СтрокаДанных of ДанныеЯчейки.КартаСтрок) {
                    if (ИДСтроки === "osp") { // если ИД строки равно osp тогда ИДСтроки равен коду осп
                        if (!!СтрокаДанных.OSP_CODE) {
                            let ИДСтрокиОСП = СтрокаДанных.OSP_CODE;
                            СтрокаТаблицы = контейнер.rows[`OSP.${ИДСтрокиОСП}`]
                        } else {
                            console.log("нет кода осп в данных стоки", СтрокаДанных);
                        }
                    } else {
                        СтрокаТаблицы = контейнер.rows[`${ИДСтроки}`]
                    }
                    if (!!СтрокаДанных.количество) {
                        СтрокаТаблицы.cells[НомерЯчейки].innerHTML = СтрокаДанных.количество;
                    } else {

                        console.log("нет СтрокаДанных.COUNT. СтрокаДанных: ", СтрокаДанных);
                    }
                    // if (!!СтрокаДанных.Ошибка && СтрокаДанных.Ошибка !==""){
                    //     Строка.cells[НомерЯчейки].insertAdjacentText("beforeend", ДанныеЯчейки.Ошибка);
                    // }
                }
            } else {
               console.log(ДанныеЯчейки)
               if (ДанныеЯчейки.Ошибка){
                    console.log(ДанныеЯчейки.Ошибка);
               } else {
                   console.log(ДанныеЯчейки)
                   console.log("ИДСтроки",ИДСтроки)
                   console.log("контейнер.rows",контейнер.rows)
                   СтрокаТаблицы = контейнер.rows[`${ИДСтроки}`]
                   console.log("СтрокаТаблицы",СтрокаТаблицы)
                   СтрокаТаблицы.cells[НомерЯчейки].innerHTML = ДанныеЯчейки;



               }
           }
        }
    } else if (!!сообщение.Контэнт.html){
        console.log(" вставляем строки ", контейнер, сообщение.Контэнт.способ_вставки)

        if (!сообщение.Контэнт.способ_вставки ){
            контейнер.tBodies[0].innerHTML = сообщение.Контэнт.html
        } else  if (сообщение.Контэнт.способ_вставки ==='добавить') {
            console.log(сообщение.Контэнт.html, контейнер.tBodies[0])
            контейнер.tBodies[0].insertAdjacentHTML("beforeend", сообщение.Контэнт.html);
        }
    }
}

function  СоздатьДетализацию(КартаСтрок){

}

function СортироватьТаблицу(event){
    console.log(event)

}

function CвернутьРазвернуть(targetId){
    let target = document.getElementById(targetId);
    target.classList.toggle("turn")
}

function ДобавитьВремя(ИдКонтейнер, ДеньНедели, event){

   let Строка = `<label class="flex_justify-center"><input type="time"  name="время[${ДеньНедели}][]"  min="00:00" max="23:59" required value=""> <button type="button" class="minus " onclick="УдалитьВремя(event)">-</button> </label>`

    event.target.insertAdjacentHTML('beforebegin',Строка )
    // document.getElementById(ИдКонтейнер).ap
}

function УдалитьВремя(event){
    event.target.parentNode.remove()
}

function ПоказатьСкрытьСелектор(event){

    article_wrapper.querySelectorAll('select:not(.hiden)').forEach(elem => {
         console.log(elem);
        elem.classList.add('hiden');
        elem.setAttribute('hidden', 'true');
        elem.removeAttribute('required')}
    );
 console.log(event.target);
    elId=event.target.options[event.target.selectedIndex ].value;
    console.log('elId :' ,elId);
    selector=document.getElementById('PR_VID2_'+elId);
    selector.classList.remove('hiden');
    selector.removeAttribute('hidden');
    selector.setAttribute('required', 'required');
    console.log(selector)
}
//
// function ВставитШаблон(event) {
//     let Кнопка = event.target
//      console.log(Кнопка);
//     if (!!Кнопка.dataset.target && Кнопка.dataset.target != "") {
//        let Цель = document.getElementById(Кнопка.dataset.target)
//        let шаблон = document.getElementById(Кнопка.dataset.tpl)
//          console.log(шаблон);
//         let блок
//        if (!!шаблон) {
//             блок = шаблон.content.cloneNode(true);
//        }
//
//        console.log(блок);
//
//        let labels =  блок.querySelectorAll("label")
//         console.log(labels);
//         if (labels.length > 0){
//
//             for (label of labels){
//                  console.log(label);
//                 данныеFor = label.getAttribute("for")
//                  console.log(данныеFor);
//             }
//         }
//
//
//         Цель.querySelectorAll("label")
//         Цель.appendChild(блок)
//     }
// }