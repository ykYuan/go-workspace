package main

import (
    "fmt"
    "log"
    "net"
    "os"
    "bufio"
    "io"
)

var data = make(map[string]string)




func handle(conn net.Conn) {
    defer conn.Close()
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        ln := scanner.Text()
        fmt.Println(ln)
        }

}

func take(conn net.Conn) {
    scanner2 := bufio.NewScanner(os.Stdin)
    for scanner2.Scan() {
    ln2 := scanner2.Text()
    io.WriteString(conn, ln2)
}

}

func main() {
    li, err := net.Listen("tcp", ":9000")
    if err != nil {
        log.Fatalln(err)
    }
    defer li.Close()
    for {
        conn, err := li.Accept()
        if err != nil {
            log.Fatalln(err)
        }
        go handle(conn)
        go take(conn)
    }
}
