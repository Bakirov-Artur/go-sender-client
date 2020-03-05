package main
import (
    "fmt"
    "os"
    "net"
    "io"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"bufio"
	"strings"
	"time"
)

type Message struct {
	gorm.Model
	Content string
	Status byte
}

func main() {
	//Create connect to data server
    conn, err := net.Dial("tcp", "127.0.0.1:9876")
	//Set deadline connect
    if err != nil {
        fmt.Println(err)
        return
    }

    defer conn.Close()
	//Connect to data storage
	db, err := gorm.Open("sqlite3", "client-msg.db")
	if err != nil {
		panic("Ошибка подключения к базе данных")
	}
	defer db.Close()
	// Migrate the schema
	db.AutoMigrate(&Message{})
	//User input message
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Message shell")
	fmt.Println("---------------------")
	buf := make([]byte, 1024)
	for {
		fmt.Print("Введите сообщение \n-> ")
		text, _ := reader.ReadString('\n')
		db.Create(&Message{Content: text, Status: 0})
		fmt.Println("Cообщение сохранен\n")
		//Send maessage to data server
		conn.SetDeadline(time.Now().Add(time.Second*5))
		resp, werr := io.CopyBuffer(conn, strings.NewReader(text), buf)
		if resp == 0 || werr != nil {
			fmt.Println("Не удалось отправить сообщение серверу...")
			continue
		} else {
			fmt.Println("Cообщение отправлен")
		}

		//Get response from server.
		conn.SetDeadline(time.Now().Add(time.Second*5))
		var respBuf strings.Builder
		resp, werr = io.CopyBuffer(&respBuf, conn, buf)
		if resp <= 0 || werr != nil {
			fmt.Println("Не удалось получить ответ от сервера")
			continue
		} else {
			fmt.Printf("Ответ от сервера: %s", respBuf.String())
		}
		fmt.Print("\n")
	}
	fmt.Print("Работа завершена")

}
