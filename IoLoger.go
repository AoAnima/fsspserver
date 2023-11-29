package main

import (
	"context"
	//"fmt"
	"github.com/jackc/pgconn"
	"time"
)

type IoLog struct{
	DB *pgconn.PgConn
	QueryName string
	ctx context.Context
	cancel context.CancelFunc
}

func (iolog *IoLog) Write (b []byte)(int, error){
	//fmt.Printf("IoLog Write%+v time.Now %+v\n", string(b), strconv.FormatInt(time.Now().Unix(), 10))
	//fmt.Printf("IoLog Write%+v time.Now %+v\n", string(b), strconv.FormatInt(time.Now().Unix(), 10))
	//2019-12-10 10:45:04.431334
	timestamp := time.Now().Format("2006-01-02 15:04:05.000000")
	ResultReader := iolog.DB.ExecPrepared(iolog.ctx, iolog.QueryName, [][]byte{[]byte(timestamp),b}, nil, nil)
	if ResultReader.Read().Err != nil {
		Ошибка("ResultReader %+v\n", ResultReader.Read().Err.Error())
	}
	return len(b), nil
}