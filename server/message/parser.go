package message

import (
	"encoding/json"
)

type MessageType byte

const (
  Undefined = 0x00
  Config = 0x01
  RunCmd = 0x02
  Response = 0xFF
)

type MessageHeader struct {
  Version       byte
  Type          MessageType
  Length        uint32
}

type Configuration struct {
  Pattern       string                `json:"pattern"`
  Dir           string                `json:"dir"`
  Adapter       string                `json:"adapter"`
  Exclude       []string              `json:"exclude"`
  Properties    map[string]string     `json:"props"`
}

type RunCommand struct {
  Pattern       string                `json:"pattern"`
  Mode          string                `json:"mode"` // Mode can be SINGLE, FILE or ALL
}

func byteToMessageType(b byte) (MessageType) {
  switch b {
    case byte(Config):
      return Config
    case byte(RunCmd):
      return RunCmd
    default:
      return Undefined
  }
}

func toUint32(b []byte) (uint32) {
  return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

func fromUint32(u uint32) ([]byte) {
  buf := make([]byte, 4)
  buf[0] = byte(u >> 24)
  buf[1] = byte(u >> 16)
  buf[2] = byte(u >> 8)
  buf[3] = byte(u)
  return buf
}

func ParseMessage(m []byte) (MessageHeader, error) {
  var message MessageHeader
  message.Version = m[0]
  message.Type = byteToMessageType(m[1])
  message.Length = toUint32(m[2:6])
  return message, nil
}

func ParseConfig(c []byte) (Configuration, error) {
  var config Configuration
  err := json.Unmarshal(c, &config)
  return config, err
}

func ParseRunCmd(c []byte) (RunCommand, error) {
  var cmd RunCommand
  err := json.Unmarshal(c, &cmd)
  return cmd, err
}

func CreateMessageHeader(version byte, msgType MessageType, length uint32) []byte {
  buf := make([]byte, 2)
  buf[0] = version
  buf[1] = byte(msgType)
  buf = append(buf, fromUint32(length)...)
  dummy := []byte{0x00, 0x00}
  buf = append(buf, dummy...)
  return buf
}
