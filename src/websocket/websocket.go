package main

import (
    "fmt"
    "golang.org/x/net/websocket"
    "os"
    "net/http"
)

func getAndRespond (ws *websocket.Conn) {
    var err error
    for {
        var reply string
        if err = websocket.Message.Receive (ws, &reply); err != nil {
            break
        }
        fmt.Println("Received back from client: " + reply)
    }
}


func main() {

    http.Handle("/", websocket.Handler(getAndRespond))
    err := http.ListenAndServe(":9000", nil)
    checkError(err)
}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }

}
