package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"jest/scanner/message"
	"jest/scanner/scanner"
)

var upgrader = websocket.Upgrader{
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,
}
var configuration message.Configuration

func main() {
  http.HandleFunc("/", ping)
  http.HandleFunc("/test", testWebsocket)
  http.HandleFunc("/scanner", startWebsocket)
  log.Fatal(http.ListenAndServe(":8080", nil))
}


func ping(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w,"I'm alive")
}

func testWebsocket(w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return
  }

  defer conn.Close()

  for {
    messageType, p, err := conn.ReadMessage()
    if err != nil {
      log.Println(err)
      return
    }

    if err := conn.WriteMessage(messageType, p); err != nil {
      log.Println(err)
      return
    }
  }
}

func startWebsocket(w http.ResponseWriter, r *http.Request) {
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return
  }

  defer ws.Close()

  for {
    err = recieveAndDecodeMsg(ws)
    if err != nil {
      log.Println(err)
      return
    }
  }
}

func recieveAndDecodeMsg(ws *websocket.Conn) error {
  messageType, p, err := ws.ReadMessage()
  if err != nil {
    return err
  }

  msg, err := message.ParseMessage(p)
  if err != nil {
    return err
  }
  log.Println(msg)

  switch msg.Type {
    case "config":
      configuration, err = message.ParseConfig(msg.Data)
      if err != nil {
        return err
      }
      log.Println("Calling set up")
      testTree := scanner.SetUp(configuration)
      testTreeJson, err := json.Marshal(testTree)
      if err != nil {
        log.Println("Error", err)
      }
      err = ws.WriteMessage(messageType, testTreeJson)
      if err != nil {
        log.Println("Error sending message", err)
      }
    case "runcmd":
      cmd, err := message.ParseRunCmd(msg.Data)
      if err != nil {
        return err
      }
      log.Println(cmd)
    default:
      log.Println("Type not recognized")
  }

  return nil
}
