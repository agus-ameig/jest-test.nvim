local utils = require("lua.jest-test.internals.utils")
local pretty_print_table = {
  ["DIR"] = function (dir, level)
    print("abc")
    return utils.tab_string(level).."D "..dir.name.."\n"
  end,
  ["FILE"] = function (file, level)
    local tests_pp = ""
    for _, suite in ipairs(file.test_cases)  do
      tests_pp = tests_pp..utils.tab_string(level+1).."TS "..suite.name.."\n"
      print("t", tests_pp)
      for _, test in ipairs(suite.tests) do
        tests_pp = tests_pp..utils.tab_string(level+2).."T "..test.name.."\n"
      end
    end
    return utils.tab_string(level).."F "..file.name.."\n"..tests_pp
  end
}
local TestTree = {}
TestTree.__index = TestTree


function TestTree.new(tree, config)
  local self =  {
    tree = tree,
    config = config
  }
  return setmetatable(self, TestTree)
end

function TestTree:get_test_file_by_path(path)
 local relative_path = path:gsub(self.config.dir, "")
 local path_components = utils.split(relative_path, "/")
 local current_item = self.tree
 for _,name in ipairs(path_components) do
    current_item = utils.find(current_item.children, function (item)
      return item.name == name
    end)
 end
 return current_item
end


local function pretty_print_children(tree, level)
  print("Level", level)
  local func = pretty_print_table[tree.type]
  print(vim.json.encode(tree))
  print(tree.type, func)
  local result = ""
  if func then
    result = func(tree, level)
    print("res", result)
  end

  if #tree.children == 0 then
    return result
  end

  for _,item in ipairs(tree.children) do
    result = result..pretty_print_children(item, level + 1)
  end
  return result
end

function TestTree:pretty_print()
  return pretty_print_children(self.tree, 0)
end

return TestTree
