package main
import (
    "fmt"
    "os"
    "net"
    "io"
)
func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:9876") 
    if err != nil {
        fmt.Println(err)
        return
    }

    defer conn.Close()

    io.Copy(os.Stdout, conn)
    fmt.Println("\nDone")
}
