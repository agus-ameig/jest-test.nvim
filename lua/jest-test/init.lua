package.loaded["lua.jest-test.network.network_bus"] = nil
package.loaded["lua.jest-test.network.network_utils"] = nil
package.loaded["lua.jest-test.internals.configuration"] = nil
package.loaded["lua.jest-test.internals.test_tree"] = nil
package.loaded["lua.jest-test.internals.utils"] = nil
local Network = require("lua.jest-test.network.network_bus")
local TestTree = require("lua.jest-test.internals.test_tree")
local configuration = require("lua.jest-test.internals.configuration")
local bus = Network.new()

local tree = nil
local config = configuration.create_configuration({
  pattern = ".*[.]spec[.]js",
  dir = "/Users/agustinameigenda/Documents/personal/test",
  adapter = "jest",
  exclude = { "node_modules", ".git" },
})

-- vim.print(config)

local function on_data()
  local result = bus:getLatestMessage()
  tree = TestTree.new(result.payload, vim.json.decode(config))
  vim.print(vim.json.encode(tree:get_test_file_by_path("/Users/agustinameigenda/Documents/personal/test/index.spec.js")))
  vim.print(tree:pretty_print())
end

bus:enable()
bus:send_and_rcv(
     string.char(1),
     config,
     on_data
)

--[[ local path = vim.api.nvim_buf_get_name(vim.api.nvim_get_current_buf())
vim.print(Utils)
local arr =  Utils.split(path, "/")

for i, word in ipairs(arr) do
    print(i, word)
end]]--

-- "{\"pattern\":\".*[.]spec[.]js\",\"dir\":\"/Users/agustinameigenda/Documents/personal/test\",\"adapter\":\"jest\",\"exclude\":[\"node_modules\"],\"props\":{}}",
