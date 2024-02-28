package main

import (
	"net"
  "log"
)
// ===========================================
// ============= Protocol Message ============
// ===========================================
// | V | TYPE |  PAYLOAD LENGTH| DUMMY  | PAYLOAD |
// | 1 |  1   |     4          |   2    | ...    |
// ===========================================
func main() {
  var v byte = 1
  var t byte = 1
  var d byte = 0
  data := "{\"pattern\":\".*[.]spec[.]js\",\"dir\":\"/Users/agustinameigenda/Documents/personal/test\",\"adapter\":\"jest\",\"exclude\":[\"node_modules\"],\"props\":{}}"
  datalen := uint32(len(data))
  conn, err := net.Dial("tcp", "localhost:8080")
  if err != nil {
    log.Println(err)
    return
  }

  datalenBytes := make([]byte, 4)
  datalenBytes[0] = byte(datalen >> 24)
  datalenBytes[1] = byte(datalen >> 16)
  datalenBytes[2] = byte(datalen >> 8)
  datalenBytes[3] = byte(datalen)
  var concat []byte
  concat = append(concat,v,t,datalenBytes[0],datalenBytes[1],datalenBytes[2],datalenBytes[3],d,d)
  concat = append(concat, []byte(data)...)

  _, err = conn.Write(concat)
  if err != nil {
    log.Println(err)
    return
  }

    conn.Read()


    conn.Close()
}
