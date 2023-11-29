package main

import(
	"github.com/jlaffaye/ftp"
	"time"
)


func FTP () *ftp.ServerConn{
	c, err := ftp.Dial("10.26.4.14:22", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		Ошибка("%+v", err)
	}
	err = c.Login("root", "Servroot88#")
	if err != nil {
		Ошибка("%+v", err)
	}
	return c
}

func СписокФайловНаFTP(c *ftp.ServerConn, папка string){

}