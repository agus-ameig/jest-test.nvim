local M = {}

function M.split(str, delimiter)
  local res = {}

  for item in str:gmatch("([^"..delimiter.."]+)") do
    table.insert(res, item)
  end

  return res
end

function M.find(arr, callback)
  for _,item in ipairs(arr) do
    if callback(item) then
      return item
    end
  end
end

function M.tab_string(amount)
  local res = ""
  for _=1,amount do
    res = res.."\t"
  end
  return res
end

return M
