package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"jest/scanner/message"
	"jest/scanner/scanner"
)

var configuration message.Configuration

func main() {
  ln, err := net.Listen("tcp", ":8080")
  if err != nil {
    log.Println("Error starting: ", err)
    return
  }

  for {
    conn, err := ln.Accept()
    if err != nil {
      log.Println("Could not accept", err)
      continue
    }
    go handleConnection(conn)
  }

}

func handleConnection(conn net.Conn) {
  defer conn.Close()

  buf := make([]byte, 8)
  _, err := conn.Read(buf)

  if err != nil {
    fmt.Println("Error reading: ", err)
    return
  }

  msgHeader, err := message.ParseMessage(buf)
  fmt.Printf("Received: version: %d, type: %d, payload length: %d\n", msgHeader.Version, msgHeader.Type, msgHeader.Length)
  payload := make([]byte, msgHeader.Length)
  _, err = conn.Read(payload)
  if err != nil {
    log.Println("Error reading data: ", err)
  }

  fmt.Printf("Payload: %s", payload)
  msg,err := recieveAndDecodeMsg(msgHeader, payload)
  responseHeader := message.CreateMessageHeader(0x01, message.Response, uint32(len(msg)))
  response := append(responseHeader, msg...)
  log.Println("%s", msg)
  conn.Write(response)
}


func recieveAndDecodeMsg(header message.MessageHeader, payload []byte) ([]byte, error) {
  switch header.Type {
    case message.Config:
      configuration, err := message.ParseConfig(payload)
      if err != nil {
        return []byte{}, err
      }
      log.Println("Calling set up")
      testTree := scanner.SetUp(configuration)
      testTreeJson, err := json.Marshal(testTree)
      if err != nil {
        log.Println("Error", err)
      }
      return testTreeJson, err
    case message.RunCmd:
      _, err := message.ParseRunCmd(payload)
      if err != nil {
        return []byte{}, err
      }
      return []byte{}, nil
    default:
      log.Println("Type not recognized")
  }

  return []byte{}, nil
}
