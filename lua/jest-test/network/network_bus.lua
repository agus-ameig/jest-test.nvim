local network_utils = require("lua.jest-test.network.network_utils")

local NetworkAdapter = {}
NetworkAdapter.__index = NetworkAdapter

function NetworkAdapter.new(opts)
  opts = opts or {}
  opts.uri = opts.uri or '127.0.0.1'
  opts.port = opts.port or 8080

  local self = {
    state = "stopped",
    messages = {},
    opts = opts,

    client = nil
  }

  return setmetatable(self, NetworkAdapter)
end

function NetworkAdapter:enable()
   local uv = vim.loop
   local client = uv.new_tcp()

   self.client = client
   self.state = "connecting"

   client:connect(self.opts.uri, self.opts.port, function (err)
      if self.state == "stopped" and err == nil then
          pcall(self.client.close, self.client)
          return
      end
      if err ~= nil then
        self.state = "error"
        error("Failed to connect to server"..err)
      else
        self.state = "connected"
      end
   end)
end

function NetworkAdapter:send(type, payload)
  if self.state ~= "connected" then
    return
  end
  self.client:write(network_utils.createMessage(type, payload))
end

function NetworkAdapter:rcv(callback)
  if self.state ~= "connected" then
    return
  end
  self.client:read_start(function (err, data)
    if err then
      self.state = "error"
      error("Failed to read data"..err)
      return
    end
    if data then
      self.messages[#self.messages + 1] = data
      pcall(callback)
      self.client:read_stop()
      return
    end

  end)

end

function NetworkAdapter:send_and_rcv(type, payload, callback)
  self.client:write(network_utils.createMessage(type, payload), function (err)
    if err then
      self.state = "error"
      error("Failed to write data"..err)
      return
    end
    self:rcv(callback)
  end)
end

function NetworkAdapter:getLatestMessage()
  if #self.messages == 0 then
    return
  end
  local msg = self.messages[#self.messages]
  table.remove(self.messages)
  return network_utils.getMessage(msg)
end

return NetworkAdapter
