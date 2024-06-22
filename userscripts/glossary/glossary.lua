local M = {}

local u = require("utils")
local xk = u.xk

M.preamble = function()
    local path = xk("path")[1]
    if path == nil then
        return
    end
    tex.print('\\input{' .. path .. '/standalone_preamble.tex}')
end

-- Print references from zettel in pretty format
local references = function(z)
    local refs = xk("ref ls -z " .. z)
    if refs == nil then
        return
    end

    if next(refs) == nil then
        tex.print('no references\\', "")
        return
    end

    tex.print('references:\\', "")
    for _, r in pairs(refs) do
        tex.print(r .. '\\', "")
    end
end

-- input zettel source file
local zettel = function(z)
    local path = xk("path -z " .. z)[1] .. "/zettel.tex"
    tex.print('\\input{' .. path .. '}', "")
end

-- input all zettels from zettelkasten
M.all = function()
    local zettels = xk("ls")
    if zettels == nil then
        return
    end

    table.sort(zettels, u.sort_alphabetical)

    for _, z in pairs(zettels) do
        if u.health(z) == 0 then
            zettel(z)
            references(z)
        end
    end
end

return M
