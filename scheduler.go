package main
import (
	"github.com/carlescere/scheduler"
	"github.com/jinzhu/gorm"
	"fmt"
	"net"
	"io"
	"strings"
)

func TasksRun(db *gorm.DB, srvAddr string) {
	scheduler.Every(10).Seconds().Run(MsgSender)
	scheduler.Every(20).Seconds().Run(MsgSender(1, db, srvAddr))
}

func MsgSender(status byte, db *gorm.DB, srvAddr string){
	msgs := GetMsgs(db, status)
	//Change status SAVED to SCHEDUL
	for _, msg := range msgs {
		SetMsgStatus(db, 1, msg)
	}
	//Connect to server
	conn, cerr := Connect(srvAddr)
	defer conn
	if cerr != nil {
		fmt.Println(cerr)
		return
	}
	//Send maessage
	for _, msg := range msgs {
		resp, werr := SendMsg(conn, msg.Content, 5)
		if resp == 0 || werr != nil {
			fmt.Println("Не удалось отправить сообщение серверу...")
			return
		} else {
			//Change status SCHEDUL to SENT
			SetMsgStatus(db, 2, msg)
			fmt.Println("Cообщение отправлен")
		}
	}
}

func Connect(srvAddr string) ( net.Conn, error) {
	return net.Dial("tcp", srvAddr)
}

func SendMsg(conn *net.Conn, msg string, timeout int) ( int64, error){
	buf := make([]byte, 1024)
	//Send maessage to data server
	conn.SetDeadline(time.Now().Add(time.Second*timeout))
	return io.CopyBuffer(conn, strings.NewReader(msg), buf)

//	//Get response from server.
//	conn.SetDeadline(time.Now().Add(time.Second*timeout))
//	var respBuf strings.Builder
//	resp, werr = io.CopyBuffer(&respBuf, conn, buf)
//	if resp <= 0 || werr != nil {
//		fmt.Println("Не удалось получить ответ от сервера")
//		return
//	} else {
//		fmt.Printf("Ответ от сервера: %s", respBuf.String())
//	}
}

func GetMsgs(db *gorm.DB, status byte) (msgs []Message){
	db.Where("status = ?", status).Find(&msgs)
}

func SetMsgStatus(db *gorm.DB, status byte, msg *Message){
	db.Model(msg).Update("status", status)
}
