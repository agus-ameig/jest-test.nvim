local M = {}

M.default_config = {
  pattern = ".*[.]spec[.]ts",
  adapter = "jest",
  exclude = { "node_modules", ".git" },
}


function M.create_configuration(opts)
  local config = opts or M.default_config
  config.pattern = config.pattern or M.default_config.pattern
  config.adapter = config.adapter or M.default_config.adapter
  config.exclude = config.exclude or M.default_config.exclude
  config.props = config.props

  config.dir = config.dir or vim.fn.getcwd()

  return vim.json.encode(config)
end

return M
