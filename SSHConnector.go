package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)


func(client *Client)СloseSSH(mes Сообщение){


	if server.Clients[mes.Выполнить.Arg.Login]!= nil && server.Clients[mes.Выполнить.Arg.Login].Ssh != nil{

		err := server.Clients[mes.Выполнить.Arg.Login].Ssh.Close()
		if err != nil {
			Ошибка("err	 %+v\n", err)
		}

	responseMes := Сообщение{
		Id:      0,
		От:      "io",
		Кому:    client.Login,
		Content:  struct {
			Target string `json:"target"`
			Data interface{} `json:"data"`
			Html string `json:"html"`
			Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
		}{
			Target:"ws_terminal_log_"+mes.Выполнить.Arg.Login,
			//Data: render("terminalLog", "Соединение установленно"),
			Data: map[string]map[string]string{
				"TerminalLog": {
					"prefix":"IO#",
					"text": "Соединение зыкрыто c ip:" +server.Clients[mes.Выполнить.Arg.Login].Ip,
				},
				"SSHclient":{"login":mes.Выполнить.Arg.Login, "ip":server.Clients[mes.Выполнить.Arg.Login].Ip},
			},
		},
	}
	Инфо("responseMes %+v\n", responseMes)
	//client.Message<-&responseMes
		responseMes.СохранитьИОтправить(client)
	} else {
		responseMes := Сообщение{
			Id:      0,
			От:      "io",
			Кому:    client.Login,
			Content:  struct {
				Target string `json:"target"`
				Data interface{} `json:"data"`
				Html string `json:"html"`
				Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
			}{
				Target:"ws_terminal_log_"+mes.Выполнить.Arg.Login,
				//Data: render("terminalLog", "Соединение установленно"),
				Data: map[string]map[string]string{
					"TerminalLog": {
						"prefix":"IO#",
						"text": "Пользователь не в сети",
					},
					"SSHclient":{"login":mes.Выполнить.Arg.Login, "ip":server.Clients[mes.Выполнить.Arg.Login].Ip},
				},
			},
		}
		Инфо("responseMes %+v\n", responseMes)
		responseMes.СохранитьИОтправить(client)
		//client.Message<-&responseMes
	}
}


/*
SSHConnect подключается по ssh  к клиенту, если в сообщение в mes.Выполнить.Arg.Login передан логин то выполняется попытка найти ip адрес по логину в памяти сервера, если ip найден то выполняется попытка подклчения по ip
Иначе если в сообщении нет логина то подключение производится к текущему клиенту от которого пришёл запрос
*/
func (client *Client) SSHConnect (mes Сообщение) {
	IP :=  client.Ip

	ConnectedLogin := mes.Выполнить.Arg.Login

	if ConnectedLogin != "" && server.Clients[ConnectedLogin].Ssh != nil{
			responseMes := Сообщение{
				Id:      0,
				От:      "io",
				Кому:    client.Login,
				Content:  struct {
					Target string `json:"target"`
					Data interface{} `json:"data"`
					Html string `json:"html"`
					Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
				}{
					Target:"ws_terminal_log_"+ConnectedLogin,
					//Data: render("terminalLog", "Соединение установленно"),
					Data: map[string]map[string]string{
						"TerminalLog": {
							"prefix":"root@"+server.Clients[ConnectedLogin].Ip+"#",
							"text": "Соединение Активно",
						},
						"SSHclient":{"login":ConnectedLogin, "ip":server.Clients[ConnectedLogin].Ip},
					},
				},
			}
			Инфо("responseMes %+v\n", responseMes)
			//client.Message<-&responseMes
			responseMes.СохранитьИОтправить(client)
		return
		}

	var config *ssh.ClientConfig
	config = &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("123QWer"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 5 * time.Second,
	}

	if ConnectedLogin != ""{
		IP = server.Clients[ConnectedLogin].Ip
	}
	//СообщениеКлиенту := Сообщение{
	//	Id:      0,
	//	От:      "io",
	//	Кому:    client.Login,
	//	Content:  struct {
	//		Target string `json:"target"`
	//		Data interface{} `json:"data"`
	//		Html string `json:"html"`
	//		Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
	//	}{
	//		Target:"ws_terminal_log_"+ConnectedLogin,
	//		//Data: render("terminalLog", "Соединение установленно"),
	//		Data: map[string]map[string]string{
	//			"TerminalLog": {
	//				"prefix":"IO@"+client.Login+"#",
	//				"text": "Устанавливаю соединени с "+ConnectedLogin+"("+IP+")",
	//			},
	//			"SSHclient":{"login":ConnectedLogin, "ip":IP},
	//		},
	//	},
	//}
	//	client.Message<-СообщениеКлиенту
	СообщениеКлиенту:= Сообщение{
		Текст:   "Устанавливаю соединени с "+ConnectedLogin+"("+IP+")",
		От: "io",
		Кому:client.Login,
		MessageType: []string{"info"},
	}
	СообщениеКлиенту.СохранитьИОтправить(client)


	clientSsh, err := ssh.Dial("tcp",IP+":22", config)
	if err != nil || client == nil {

		Ошибка("err ssh.Dial %+v IP %+v \n", err, IP)

		if Ping(IP) {
			//print("Комп пингуеться, вероятно не подходит пароль")

			СообщениеКлиенту = Сообщение{
				Id:      0,
				От:      "io",
				Кому:    client.Login,
				Content:  struct {
					Target string `json:"target"`
					Data interface{} `json:"data"`
					Html string `json:"html"`
					Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
				}{
					Target:"ws_terminal_log_"+ConnectedLogin,
					//Data: render("terminalLog", "Соединение установленно"),
					Data: map[string]map[string]string{
						"TerminalLog": {
							"prefix":"IO@"+IP+"#",
							"text": "Не удаётся установить соединение с "+ConnectedLogin + " ("+IP+")",
						},
						"SSHclient":{"login":ConnectedLogin, "ip":IP},
					},
				},
			}

			//client.Message<-&responseMes
			СообщениеКлиенту.СохранитьИОтправить(client)
			СообщениеКлиенту = Сообщение{
				Текст:   "Не удаётся установить соединение с "+ConnectedLogin + " ("+IP+")",
				От: "io",
				Кому:client.Login,
				MessageType: []string{"errroк"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
			return
		}
		СообщениеКлиенту = Сообщение{
			Текст:   "Не удаётся установить соединение с "+ConnectedLogin + " ("+IP+")",
			От: "io",
			Кому:client.Login,
			MessageType: []string{"error"},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
		return
	}
	//log.Printf("\n \n clientSsh.User() %+v\n \n",clientSsh )
	Инфо("mes.Выполнить.Arg.Login %+v ConnectedLogin %+v\n",mes.Выполнить.Arg.Login, ConnectedLogin )

	if ConnectedLogin != "" {
		if server.Clients[ConnectedLogin]== nil{
			СообщениеКлиенту = Сообщение{
				Id:      0,
				От:      "io",
				Кому:    client.Login,
				Content:  struct {
					Target string `json:"target"`
					Data interface{} `json:"data"`
					Html string `json:"html"`
					Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
				}{
					Target:"ws_terminal_log_"+ConnectedLogin,
					//Data: render("terminalLog", "Соединение установленно"),
					Data: map[string]map[string]string{
						"TerminalLog": {
							"prefix":"root@"+IP+"#",
							"text": "Соединение не установленно, клиент разорвал соединение или выключил компьютер",
						},
						"SSHclient":{"login":ConnectedLogin, "ip":IP},
					},
				},
			}

			//client.Message<-&responseMes
			СообщениеКлиенту.СохранитьИОтправить(client)
			СообщениеКлиенту = Сообщение{
				Текст:   "Соединение не установленно, клиент разорвал соединение или выключил компьютер",
				От: "io",
				Кому:client.Login,
				MessageType: []string{"info"},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
			return
		}

		server.mx.Lock()
		server.Clients[ConnectedLogin].Ssh = clientSsh
		server.mx.Unlock()

		СообщениеКлиенту = Сообщение{
			От:      "io",
			Кому:    client.Login,
			Content:  struct {
				Target string `json:"target"`
				Data interface{} `json:"data"`
				Html string `json:"html"`
				Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
			}{
				Target:"ws_terminal_log_"+ConnectedLogin,
				//Data: render("terminalLog", "Соединение установленно"),
				Data: map[string]map[string]string{
					"TerminalLog": {
						"prefix":"root@"+IP+"#",
						"text": "Соединение установленно",
					},
					"SSHclient":{"login":ConnectedLogin, "ip":IP},
				},
			},
		}
		//log.Printf("responseMes %+v\n", responseMes)
		СообщениеКлиенту.СохранитьИОтправить(client)
		СообщениеКлиенту = Сообщение{
			Текст:   "Соединение установленно",
			От: "io",
			Кому:client.Login,
			MessageType: []string{"info"},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
		return
		//client.Message<-&responseMes
	} else {
		client.Ssh = clientSsh
		СообщениеКлиенту= Сообщение{
					Текст:   "Соединение установленно c IP: " + IP,
					От: "io",
					Кому:client.Login,
					MessageType: []string{"info"},
				}
		СообщениеКлиенту.СохранитьИОтправить(client)
				//СообщениеКлиенту.СохранитьЛогСообщения()
				//client.Message<-СообщениеКлиенту
	}
}

func  (client *Client) ВыполнитьНесколькоКоммандПоSSH (сообщение Сообщение){
	ConnectedLogin := сообщение.Выполнить.Arg.Login
	Инфо("\n  >>>>> ВыполнитьНесколькоКммандПоSSH ConnectedLogin %+v\n", ConnectedLogin)
	if ConnectedLogin == "" {
		СообщениеКлиенту:= &Сообщение{
			Текст:   "Не могу понять к кому необходимо установить соединение, не передан логин пользователя "+ConnectedLogin,
			От: "io",
			Кому:client.Login,
			MessageType: []string{"irritation","io_action"},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
		return
	} else {
		client.SSHConnect(Сообщение{})
	}

	SSH := server.Clients[ConnectedLogin].Ssh
	var session *ssh.Session
	var errSession error
	var stdin io.WriteCloser
	var commandLog bytes.Buffer
	var errorLog bytes.Buffer

	if server.Clients[ConnectedLogin].SSHSession == nil {
		server.mx.Lock()
		session, errSession = SSH.NewSession()
		Инфо("session %+v, errSession %+v\n", session, errSession)
		server.Clients[ConnectedLogin].SSHSession = session


		stdin, errSession = session.StdinPipe()
		Инфо("stdin %+v\n", stdin)
		if errSession != nil {
			Ошибка(">>>> ERROR errSession \n %+v \n\n", errSession)
		}
		session.Stdout = &commandLog
		session.Stderr = &errorLog
		server.Clients[ConnectedLogin].SshLog = commandLog
		server.Clients[ConnectedLogin].SshErr = errorLog
		server.Clients[ConnectedLogin].SshIn = stdin
		err := session.Shell()
		if err != nil {
			Ошибка(">>>> ERROR \n %+v \n\n", err)
		}
		server.mx.Unlock()

	} else {
		server.mx.RLock()
		session = server.Clients[ConnectedLogin].SSHSession
		commandLog = server.Clients[ConnectedLogin].SshLog
		errorLog = server.Clients[ConnectedLogin].SshErr
		stdin =server.Clients[ConnectedLogin].SshIn
		server.mx.RUnlock()
	}

	if errSession != nil {
		Ошибка(">>>> ERROR \n %+v \n\n", errSession)
	}

	if errSession != nil {
		СообщениеКлиенту := Сообщение{
			От:      "io",
			Кому:    client.Login,
			Content:  struct {
				Target string `json:"target"`
				Data interface{} `json:"data"`
				Html string `json:"html"`
				Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
			}{
				Target:"ws_terminal_log_"+ConnectedLogin,
				//Data: render("terminalLog", "Соединение установленно"),
				Data: map[string]map[string]string{
					"TerminalLog": {
						"prefix":"root@"+server.Clients[ConnectedLogin].Ip+"#",
						"text": "Не удаётся открыть сеанс для удёлнного выполнения комманд",
					},
					"SSHclient":{"login":ConnectedLogin, "ip":server.Clients[ConnectedLogin].Ip},
				},
			},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
		return
	}

	//commands := []string{
	//	"pwd",
	//	"whoami",
	//	"echo 'bye'",
	//	"exit",
	//}
	//log.Printf("stdin %+v\n", stdin)
	//for _, cmd := range commands {
	//	_, err := fmt.Fprintf(stdin, "%s\n", "pwd")
	//	_, err := fmt.Fprintf(stdin, "%s\n", "whoami")
	//	_, err := fmt.Fprintf(stdin, "%s\n", "echo 'bye'")
	//	_, err := fmt.Fprintf(stdin, "%s\n", "exit")
	//	if err != nil {
	//		log.Printf("err %+v \n",err)
	//	}
	//}


	cmd, err := fmt.Fprintf(stdin, "%s\n", сообщение.Выполнить.Комманду)
	Инфо("cmd	 %+v\n", cmd)
	if err != nil {
		Ошибка("err %+v \n",err)
	}

	countByte, errWrite := stdin.Write([]byte(сообщение.Выполнить.Комманду))
	Инфо("сообщение.Выполнить.Комманду %+v, countByte %+v errWrite %+v\n",сообщение.Выполнить.Комманду,countByte, errWrite,)
	countByte, errWrite = stdin.Write([]byte("exit"))
	Инфо("countByte %+v errWrite %+v\n",countByte, errWrite)

	//_, _= fmt.Fprintf(stdin, "%s\n", "exit")

	errClose := stdin.Close()
	if errClose != nil {
		log.Printf(">>>> ERROR \n %+v \n\n", errClose)
	}

	//log.Printf("commandLog %+v errorLog %+v\n",commandLog, errorLog)
	errWait := session.Wait()
	if errWait != nil {
		log.Printf("errWait %+v \n",commandLog, errWait)
	}

	СообщениеКлиенту := Сообщение{
		Id:      0,
		От:      "io",
		Кому:    client.Login,
		Content:  struct {
			Target string `json:"target"`
			Data interface{} `json:"data"`
			Html string `json:"html"`
			Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
		}{
			Target:"ws_terminal_log_"+ConnectedLogin,
			Data: map[string]map[string]string{
				"TerminalLog": {
					"prefix":"root@"+server.Clients[ConnectedLogin].Ip+"#",
					"text": "<pre>"+commandLog.String()+"</pre>",
				},
				"SSHclient":{"login":ConnectedLogin, "ip":server.Clients[ConnectedLogin].Ip},
			},
		},
	}
	СообщениеКлиенту.СохранитьИОтправить(client)
}



/*
Выполняет подключение к удалённоу компьютеру, выполняет команду, отправляет ответ клиенту
*/
func  (client *Client) ВыполнитьКоммандуПоSSH (сообщение Сообщение){

	ConnectedLogin := сообщение.Выполнить.Arg.Login
	log.Printf("\n  >>>>> ВыполнитьКоммандуПоSSH ConnectedLogin %+v\n", ConnectedLogin)
	if ConnectedLogin == "" {
		СообщениеКлиенту:= &Сообщение{
					Текст:   "Не могу понять к кому необходимо установить соединение, не передан логин пользователя "+ConnectedLogin,
					От: "io",
					Кому:client.Login,
					MessageType: []string{"irritation","io_action"},
				}
				СообщениеКлиенту.СохранитьИОтправить(client)
		return
	} else {
		client.SSHConnect(Сообщение{})
	}



	if server.Clients[ConnectedLogin].Ssh != nil {
		SSH := server.Clients[ConnectedLogin].Ssh
		session, err := SSH.NewSession()
		if err != nil {
			СообщениеКлиенту := Сообщение{
				От:      "io",
				Кому:    client.Login,
				Content:  struct {
					Target string `json:"target"`
					Data interface{} `json:"data"`
					Html string `json:"html"`
					Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
				}{
					Target:"ws_terminal_log_"+ConnectedLogin,
					//Data: render("terminalLog", "Соединение установленно"),
					Data: map[string]map[string]string{
						"TerminalLog": {
							"prefix":"root@"+server.Clients[ConnectedLogin].Ip+"#",
							"text": "Не удаётся открыть сеанс для удёлнного выполнения комманд",
						},
						"SSHclient":{"login":ConnectedLogin, "ip":server.Clients[ConnectedLogin].Ip},
					},
				},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
			return
		}
		var commandLog bytes.Buffer
		var errorLog bytes.Buffer

			session.Stdout = &commandLog
			session.Stderr = &errorLog
			//stdBuf, err := session.StdoutPipe()

		 errCmd := session.Run(сообщение.Выполнить.Комманду)

		//		Shell := session.Shell()
		//		Shell.
		////log.Printf("errCmd %+v\n", errCmd)

		log.Printf(" commandLog.String() %+v\n",  commandLog.String())

		if errCmd != nil {
			//err := session.Run(сообщение.Выполнить.Комманду);
			//log.Printf("errorLog %+v\n",errorLog )
			//log.Printf("out %+v\n",out )
			//log.Printf("commandLog.String() %+v\n", errorLog.String())
			//text:= strings.Split(errorLog.String(), "\n")
			//log.Printf("Split text %+v\n", text)
			//textNew:= strings.ReplaceAll(errorLog.String(), "\n", "<br>")
			//textNew= strings.ReplaceAll(textNew, "\t", "&emsp;")
			//textNew= strings.ReplaceAll(textNew, " ", "&emsp;")

			log.Printf("errCmd %+v сообщение.Выполнить.Комманду %+v\n", errCmd, сообщение.Выполнить.Комманду)

			СообщениеКлиенту := Сообщение{
				Id:      0,
				От:      "io",
				Кому:    client.Login,
				Content:  struct {
					Target string `json:"target"`
					Data interface{} `json:"data"`
					Html string `json:"html"`
					Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
				}{
					Target:"ws_terminal_log_"+ConnectedLogin,
					Data: map[string]map[string]string{
						"TerminalLog": {
							"prefix":"root@"+server.Clients[ConnectedLogin].Ip+"#",
							"text": "<pre>"+errorLog.String()+"</pre>",//errorLog.String() + commandLog.String(),
						},
						"SSHclient":{"login":ConnectedLogin, "ip":server.Clients[ConnectedLogin].Ip},
					},
				},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
			return
		} else {

			//
			//log.Printf("commandLog.String() %+v\n", commandLog.String())
			//text:= strings.Split(commandLog.String(), "\n")
			//log.Printf("Split text %+v\n", text)
			//textNew:= strings.ReplaceAll(commandLog.String(), "\n", "<br>")
			//textNew= strings.ReplaceAll(textNew, "\t", "&nbsp;")
			//
			//log.Printf("ReplaceAll text %+v\n", textNew)
			СообщениеКлиенту := Сообщение{
				Id:      0,
				От:      "io",
				Кому:    client.Login,
				Content:  struct {
					Target string `json:"target"`
					Data interface{} `json:"data"`
					Html string `json:"html"`
					Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
				}{
					Target:"ws_terminal_log_"+ConnectedLogin,
					Data: map[string]map[string]string{
						"TerminalLog": {
							"prefix":"root@"+server.Clients[ConnectedLogin].Ip+"#",
							"text": "<pre>"+commandLog.String()+"</pre>",
						},
						"SSHclient":{"login":ConnectedLogin, "ip":server.Clients[ConnectedLogin].Ip},
					},
				},
			}
			СообщениеКлиенту.СохранитьИОтправить(client)
		}
		defer session.Close()
		return
	}
}
func isMatch(bytes []byte, t int, matchingBytes []byte) bool {
	if t >= len(matchingBytes) {
		for i := 0; i < len(matchingBytes); i++ {
			if bytes[t - len(matchingBytes) + i] != matchingBytes[i] {
				return false
			}
		}
		return true
	}
	return false
}
var escapePrompt = []byte{'$', ' '}

func readUntil(r io.Reader, matchingByte []byte) (*string, error) {
	var buf [64 * 1024]byte
	var t int
	for {
		n, err := r.Read(buf[t:])
		Инфо("n %+v\n", )
		if err != nil {
			Ошибка("err %+v\n", err)
			return nil, err
		}
		t += n
		if isMatch(buf[:t], t, matchingByte) {
			stringResult := string(buf[:t])
			Ошибка("stringResult %+v\n",stringResult )
			return &stringResult, nil
		}
	}
}
func write(w io.WriteCloser, command string) error {
	_, err := w.Write([]byte(command + "\n"))
	return err
}

func (client *Client) RunSkillCmd (cmd string) (string, []string) {

	client.SSHConnect(Сообщение{})
	var commandLog bytes.Buffer
	var errorLog bytes.Buffer
	//var in bytes.Buffer


	if client.Ssh == nil{
		СообщениеКлиенту:= &Сообщение{

			Текст:   "Не удаётся установить соединение c IP "+ client.Ip,
			От: "io",
			Кому:client.Login,
			MessageType: []string{"irritation","io_action","error"},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
		return "", []string{"Не удаётся установить соединение c IP "+ client.Ip}
	}
	session, err := client.Ssh.NewSession()

	if err != nil {
		СообщениеКлиенту := Сообщение{
			От:      "io",
			Кому:    client.Login,
			Content:  struct {
				Target string `json:"target"`
				Data interface{} `json:"data"`
				Html string `json:"html"`
				Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
			}{
				Target:"ws_terminal_log_"+ client.Login,
				//Data: render("terminalLog", "Соединение установленно"),
				Data: map[string]map[string]string{
					"TerminalLog": {
						"prefix":"root@"+client.Ip+"#",
						"text": "Не удаётся открыть сеанс для удёлнного выполнения комманд",
					},
					"SSHclient":{"login": client.Login, "ip":client.Ip},
				},
			},
		}
		СообщениеКлиенту.СохранитьИОтправить(client)
		return "", []string{err.Error()}
	}


	session.Stdout = &commandLog
	session.Stderr = &errorLog
	//in , errStderrPipe := session.StdinPipe()
	//if errStderrPipe != nil {
	//	log.Printf(">>>> ERROR \n %+v \n\n", errStderrPipe)
	//}
	//out, errStdoutPipe := session.StdoutPipe()
	//if errStdoutPipe != nil {
	//	log.Printf(">>>> ERROR \n %+v \n\n", errStdoutPipe)
	//}
	//cmdstring := flag.String("cmd", "display arp statistics all", "cmdstring")
	Инфо("\n\n!!!!!! >>>>>>>>  ВЫПОЛНЕНИЕ НАВЫКА коммана %+v на машине %+v \n\n", cmd, client.Ip)
	СообщениеКлиенту:= &Сообщение{
				Текст:   "Пытаюсь исправить Ваше проблему. Ожидайте...",
				От: "io",
				Кому:client.Login,
				MessageType: []string{"info"},
				//MessageType: []string{"irritation","io_action"},
			}
	СообщениеКлиенту.СохранитьИОтправить(client)
	//modes := ssh.TerminalModes{
	//	ssh.ECHO:          1,     // disable echoing
	//	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	//	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	//}
	//if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
	//	log.Fatal("request for pseudo terminal failed: ", err)
	//}
	//w, err := session.StdinPipe()
	//if err != nil {
	//	panic(err)
	//}
	//r, err := session.StdoutPipe()
	//if err != nil {
	//	panic(err)
	//}
	//e, err := session.StderrPipe()
	//if err != nil {
	//	panic(err)
	//}
	//
	//in, out := MuxShell(w, r, e)
	//if err := session.Shell(); err != nil {
	//	log.Fatal(err)
	//}
	//i:=<-out //ignore the shell output
	//log.Printf("i %+v\n", i)
	//in <- *cmdstring
	//fmt.Printf("%s\n", <-out)
	//
	//
	//
	//if _, err := w.Write([]byte("whoim\r")); err != nil {
	//	fmt.Printf("Failed to run: " + err.Error())
	//}
	//if _, err := w.Write([]byte(cmd)); err != nil {
	//	fmt.Printf("Failed to run: " + err.Error())
	//}
	//l := <-out
	//log.Printf("l %+v\n", l)
	//session.Wait()
	//err = session.Shell();
	//
	//if err != nil {
	//	log.Printf(
	//		">>>> ERROR \n %+v \n\n", err)
	//}
	//
	//
	//_, errWrite := in.Write([]byte(cmd+"/\n"))
	//err = session.Run(cmd)
//	if err := session.Start("/bin/sh"); err != nil {
//		log.Fatal(err)
//	}
//
//	var p []byte
//
//	outW, errRead := readUntil(out, escapePrompt)
//	log.Printf("outW %+v\n", outW)
//	fmt.Printf("out: %s\n", out)
//	if errRead != nil {
//		log.Printf(">>>> ERROR \n %+v \n\n", errRead)
//	}
//	_, errWrite := in.Write([]byte(cmd+"/\n"))
//	if errWrite != nil {
//		log.Printf(">>>> ERROR \n %+v \n\n", errWrite)
//	}
//	outW, errRead = readUntil(out, escapePrompt)
//	log.Printf("outW %+v\n", outW)
//	fmt.Printf("out: %s\n", out)
//	if errRead != nil {
//		log.Printf(">>>> ERROR \n %+v \n\n", errRead)
//	}
//

	if err := session.Run(cmd); err != nil || errorLog.String() != ""  {
		session.Wait()

		log.Printf("errorLog %+v\n", errorLog)
		log.Printf("commandLog %+v\n", commandLog.String())
		log.Printf("err %+v errorLog %+v \n", err, errorLog.String())
		ЛогВыполнения := ""
		if commandLog.String() != "" {
			ЛогВыполнения=commandLog.String()
		}
		if err != nil{
			return ЛогВыполнения, []string{strings.TrimSpace(err.Error()),strings.TrimSpace(errorLog.String())}
		} else {
			return ЛогВыполнения, []string{strings.TrimSpace(errorLog.String())}
		}

	}

	ЛогВыполнения := commandLog.String()

log.Printf("cmdLog %+v\n", ЛогВыполнения)
log.Printf("errorLog %+v\n", errorLog.String())
//log.Printf("commandLog %+v\n", commandLog.String())

	defer session.Close()
	return  "nil", nil
}

func checkError(err error, info string) {
	if err != nil {
		fmt.Printf("%s. error: %s\n", info, err)
		os.Exit(1)
	}
}

func MuxShell(w io.Writer, r, e io.Reader) (chan<- string, <-chan string) {
	in := make(chan string, 5)
	out := make(chan string, 5)
	var wg sync.WaitGroup
	wg.Add(1) //for the shell itself
	go func() {
		for cmd := range in {
			wg.Add(1)
			w.Write([]byte(cmd + "\n"))
			wg.Wait()
		}
	}()

	go func() {
		var (
			buf [1024 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:])
			if err != nil {
				fmt.Println(err.Error())
				close(in)
				close(out)
				return
			}
			t += n
			result := string(buf[:t])
			if strings.Contains(string(buf[t-n:t]), "More") {
				w.Write([]byte("\n"))
			}
			if strings.Contains(result, "username:") ||
				strings.Contains(result, "password:") ||
				strings.Contains(result, ">") {
				out <- string(buf[:t])
				t = 0
				wg.Done()
			}
		}
	}()
	return in, out
}

func Ping(ip string) bool{

	_, err := exec.Command("ping", ip, "-c 2", "-i 2").Output()
	//Инфо("Ping %+v",string(s), err )
	//out, _ := exec.Command("ssh", "root@"+ip).Output()
	if err != nil{
		return false
	} else {
		//log.Printf("out %+v\n", string(out))
		return true
	}
}