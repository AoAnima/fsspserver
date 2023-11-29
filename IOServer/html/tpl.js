{{define "tplsJs"}}
let tpl = function(tplName, values){
        // console.log("tpl values", values)
        const tpls = {

            метка:function(метка){
                const TagId=`tag_${new Date().getTime()}`
                const newTag = `<label id="${TagId}" class="hashtag"><input readonly type="text" name="метки" value="${метка}"><button class="remove_tag" onclick="Закрыть('${TagId}')">x</button></label>`
                return newTag
            },
            terminalLogString:function(TerminalLog){
                // console.log(TerminalLog);
                let logString = `<div class="terminal_log_wrapper" id="">  
                                     <div class="log_prefix">${TerminalLog.prefix}</div>
                                     <div class="log_text">${TerminalLog.text}</div>
                                  </div>`;
                return logString
            },
            LogWrapper:function (value){
                let ws_messages = `<div class="log_wrapper" id="log_wrapper_${value.login}">                                     
                                   </div>${WS.tpl("Terminal", value)}`;
                return ws_messages
            },
            /**
             * @return {string}
             */
            SkillWrapper:function (AdminMenu){
                let SkillsHtml =`<div class="skillWrapper">`;
                console.log(AdminMenu);
                for(let skill of AdminMenu){
                    console.log(skill.problem);
                    SkillsHtml = SkillsHtml + `<button class="btn_problem" title="${skill.description}" onclick="WS.RunSkill(${skill.action_id})">${skill.problem}</button>`;
                }
                SkillsHtml = SkillsHtml + `</div>`;
                return SkillsHtml
            },
            UserWidget:function (сообщение) {

                const user = сообщение.Content.data;
                // console.log(user);
                let widgetHtml = `<div onclick="WS.ShowChatLog(this, '${user.info.login}')" id="${user.info.login}" data-ip="${user.info.ip}" data-initials="${user.Initials}" data-ip="${сообщение.ip}" data-fullname="${user.FullName}" class="ws_contact online">                        
                         <div class="avatar">
                             ${user.Initials}
                        </div>
                        <div class="user_info">
                             <div class="name">
                              ${user.FullName}
                            </div>  
                       
                    <div class="curator_post">
                            
                    </div>
                       
                    <div class="div_name">
                           ${user.info.osp_name}
                    </div>
                        <div class="post">
                             ${user.info.post_name}
                    </div>
                    </div>
                    </div>`;
                return widgetHtml;
            },
            /**
             * @return {string}
             */
            Terminal:function(values){
                console.log("Terminal",values);
                const login = values.login;
                const AdminMenu = values.AdminMenu;
                // console.log("AdminMenu", AdminMenu)
                if (!!AdminMenu && Object.keys(AdminMenu).length>0) {
                    return `<div id="ws_terminal_${login}" data-ip="" data-x=0 data-y=0  data-login="${login}" class="log_terminal_wrapper mini offline">
                                <div class="ws_terminal_header" >   
                                    <button class="terminal_connect_btn" onclick="WS.terminalOpen('${login}')">Открыть терминал</button>   
                                    <button class="terminal_close_btn" onclick="WS.terminalClose('${login}')">Закрыть терминал</button>   
                                   
                                </div>
                                <div class="ws_terminal_body">                                   
                                    <div id="ws_terminal_log_${login}" class="ws_terminal_log"></div>       
                                                         
                                    <div class="ws_terminal_input" onkeypress="ScanKeys(event)?WS.RunCmd(event.target):null" contenteditable="true" name="cmd" id="cmd_${login}" data-login="${login}" data-ip=""></div>   
                                    <div class="fast_cmd">                                    
                                        <button class="toggler" onclick="this.parentNode.classList.toggle('expand')">Навыки ИО</button>    
                                        ${WS.tpl("SkillWrapper", AdminMenu)}            
                                    </div>     
                                </div>                                        
                            </div>`;
                } else {
                    return "";
                }

            },

            /**
             * @return {string}
             */
            UserLog: function(){
                return `<div class="ws_user_log" id="log_${login}"></div>`;
            },
            /**
             * @return {string}
             */
            Сообщение : function(m){
                // console.log("m",m.Текст);
                // console.log("m",`${m.Текст}`);
                let text =  m.Текст;

                // console.log("сожержит перенос строки",text.includes('\n'))
                let mesClass = "mes_from";
                let UserInfo = {};
                UserInfo.FullName=WS.UserInfo.FullName;
                UserInfo.Initials=WS.UserInfo.Initials;
                // console.log(WS.UserInfo);
                // console.log("m.От",m.От, "WS.UserInfo.uid", WS.UserInfo.uid)

                if (m.От !== WS.UserInfo.uid){
                    let contactBlock;
                    // console.log("m.От", m.От)
                    // console.log("m.От", document.getElementById(m.От))
                    // if (m.От === "io"){
                    //     contactBlock = document.getElementById("io");
                    // } else {
                        contactBlock = document.getElementById(m.От);
                    // }
                    mesClass = "mes_to";
                    UserInfo.FullName = contactBlock.dataset.fullname;
                    UserInfo.Initials = contactBlock.dataset.initials;
                }

                if (!!m.messageType){
                    mesClass = mesClass+" "+m.messageType.join(' ');
                }
//${m.От+"-"+m.Id+"-"+m.Кому}

                if (!!m.Content && !!m.Content["html"] && !!m.Content["target"] && m.Content["target"] === "log_"+m.От){
                    text =text+m.Content["html"]
                }
                if (!!m.Контэнт && !!m.Контэнт.контейнер && m.Контэнт.контейнер === "сообщение" && m.Контэнт.html){
                    text= text + "<br>" +m.Контэнт.html
                }

                return `<div onclick="this.parentNode.removeChild(this);" id="mes_${m.Id}" class="${mesClass} message_wrapper">
                    <div class="avatar">
                        ${UserInfo.Initials}
                    </div>
                    <div class="ws_message_block">
                        <div class="autor">
                            <div class="name">
                             ${UserInfo.FullName}
                            </div>
                            <div class="time">
                                ${m.Время}
                            </div>
                        </div>
                        <div class="ws_text_message">${text}</div>
                    </div>
                </div>`;
            },
        };
        let result = tpls[tplName](values);
        return result;

};

{{end}}