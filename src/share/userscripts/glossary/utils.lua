M = {}

local split = function(inputstr, sep)
    if sep == nil then
        sep = "%s"
    end
    local t = {}
    for str in string.gmatch(inputstr, "([^" .. sep .. "]+)") do
        table.insert(t, str)
    end
    return t
end

M.sort_alphabetical = function(a, b)
    return a:lower() < b:lower()
end

M.xk = function(command)
    local handle = io.popen("xk " .. command)
    if handle == nil then
        return
    end
    local res = handle:read("a")
    handle:close()

    return split(res, "\n")
end

M.health = function(zettel)
    return os.execute("xk script checkhealth " .. zettel)
end

return M
