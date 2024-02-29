local bit = require("bit")
local M = {}

M.VERSION = string.char(0)
M.Types = {
  Undefined = string.char(0),
  Config = string.char(1),
  RunCmd = string.char(2),
  Response = string.char(255)
}
M.DUMMY = string.char(0,0)

function M.format_payload_length(str)
  local length = #str

  local byte4 = bit.band(length, 0xFF)
  local byte3 = bit.band(bit.rshift(length, 8), 0xFF)
  local byte2 = bit.band(bit.rshift(length, 16), 0xFF)
  local byte1 = bit.band(bit.rshift(length, 24), 0xFF)

  return string.char(byte1,byte2,byte3,byte4)
end

function M.createMessage(type, payload)
  local payload_length = M.format_payload_length(payload)
  return M.VERSION..type..payload_length..M.DUMMY..payload
end

function M.getMessage(message)
  local version = string.sub(message, 1,1)
  local type = string.sub(message, 2,2)
  local payload = vim.json.decode(string.sub(message, 9))
  return {
    version = version,
    type = type,
    payload = payload
  }
end

return M
