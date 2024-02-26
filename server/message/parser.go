package message

import (
	"encoding/json"
)

type Message struct {
  Version       int                   `json:"version"`
  Type          string                `json:"type"`
  Data          json.RawMessage       `json:"data"`
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

func ParseMessage(m []byte) (Message, error) {
  var message Message
  err := json.Unmarshal(m, &message)
  return message, err
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
