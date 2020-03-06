package main
import (
    "fmt"
    "os"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"bufio"
)

type Message struct {
	gorm.Model
	Content string
	Status byte
}

func main() {
	//Connect to data storage
	db, err := gorm.Open("sqlite3", "client-msg.db")
	if err != nil {
		panic("Ошибка подключения к базе данных")
	}
	defer db.Close()
	// Migrate the schema
	db.AutoMigrate(&Message{})
	//Run tasks
	srvAddr := "127.0.0.1:9876"
	go TasksRun(db, srvAddr)
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
		fmt.Print("\n")
	}
	fmt.Print("Работа завершена")

}
