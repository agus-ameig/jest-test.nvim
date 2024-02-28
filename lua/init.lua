package.loaded["lua.network.network_bus"] = nil
package.loaded["lua.network.network_utils"] = nil
local Network = require("lua.network.network_bus")
local bus = Network.new()


local function on_data()
  local result = bus:getLatestMessage()
  vim.print(vim.json.encode(result[3]))
end

bus:enable()
bus:send_and_rcv(
     string.char(1),
     "{\"pattern\":\".*[.]spec[.]js\",\"dir\":\"/Users/agustinameigenda/Documents/personal/test\",\"adapter\":\"jest\",\"exclude\":[\"node_modules\"],\"props\":{}}",
     on_data
)
