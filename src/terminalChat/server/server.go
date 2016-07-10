package main

import (
    "fmt"
    "net"
    "bufio"
    "io"
)

type client struct {
    Name string
    Output chan message
}

type message struct {
    Username string
    Text string
}

type chatRoom struct {
    clients map[string]client
    join chan client
    leave chan client
    input chan message
}

func (cr *chatRoom) Run() {
    for {
        select {
        case msg := <-cr.input:
            for _, cl := range cr.clients {
                cl.Output <- msg
                }

        case cl := <-cr.join:
            cr.clients[cl.Name]=cl
            go func () {
                cr.input <- message {
                    Username: "System",
                    Text: fmt.Sprintf("%s joined us", cl.Name),
                    }
            }()

        case cl := <-cr.leave:
            delete(cr.clients, cl.Name)
            go func () {
                cr.input <- message {
                   Username: "System",
                   Text: fmt.Sprintf("%s left us", cl.Name),
                }
            }()

        }
    }
}

func handle (conn net.Conn, cr *chatRoom) {
    defer conn.Close()
    io.WriteString(conn, "Enter your Username:")
    scanner := bufio.NewScanner(conn)
    scanner.Scan()
    cl := client{
        Name: scanner.Text(),
        Output: make(chan message),
        }
    cr.join <- cl
    defer func () {
        cr.leave <- cl
    }()
    //constantly read information from the client
    go func () {
        for scanner.Scan() {
            txt := scanner.Text()
            cr.input <- message{
                Username: cl.Name,
                Text: txt,
            }
        }
    }()
    //print client messages
    for msg := range cl.Output {
        if msg.Username != cl.Name {
            _, err := io.WriteString(conn, msg.Text + "\n")
            if err != nil {
                break
            }
       }
    }

}


func main() {
    ln, err := net.Listen("tcp", ":9000")
    if err != nil {
      panic(err)
    }
    defer ln.Close()

    cr := &chatRoom {
        clients: make(map[string]client),
        join:  make(chan client),
        leave: make(chan client),
        input: make(chan message),
        }
    go cr.Run()

    for {
        conn, err := ln.Accept()
        if err != nil {
            panic(err)
        }
        go handle (conn, cr)

    }
}
