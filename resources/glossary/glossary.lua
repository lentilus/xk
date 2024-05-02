local lfs = require("lfs")
local M = {}


local function debugging(msg)
    texio.write_nl("DEBUG: " .. msg)
end

debugging("sourced lua")

local xettelkasten_dir = os.getenv("ZETTEL_DATA")
local healthy_list = os.getenv("healthy_string") or ""

debugging("healthy list is: "..healthy_list)
debugging("xettelkasten dir is: "..xettelkasten_dir)

local sep = "%s"
local healthy={}
for str in string.gmatch(healthy_list, "([^"..sep.."]+)") do
    debugging("healthy: " .. str)
    table.insert(healthy, str)
end

M.preamble = xettelkasten_dir .. "/standalone_preamble.tex"

debugging("using preamble: "..M.preamble)

local function no_references()
    debugging("no references")
    tex.print("\\vspace{3pt}", "\\textit{no references}")
end

M.refs = function(zettel)
    local file = io.open(xettelkasten_dir..zettel.."/references", "r")
    if file == nil or file:lines() == nil then
        no_references()
        return
    end

    local lines = file:lines()
    local found_refs = false
    for line in lines do
        debugging("here is a line")
        if line ~= "" then
            found_refs = true
            line = line:gsub("_", " ")
            tex.print("$\\rightarrow$ \\hyperlink{"..line.."}{\\textit{"..line.."}} \\newline")
        end
    end

    if not found_refs then
        no_references()
    end

    collectgarbage("collect")
end

local function sort_alphabetical(a, b)
    return a:lower() < b:lower()
end

-- this function sucks, it should not be recursive
local function insert(file)
    local path = file .. "/zettel.tex"
    debugging("adding "..path)
    tex.print("\\input{"..path.."}")
end

M.insert_zettel = function()
    debugging("insert_zettel was called")
    table.sort(healthy, sort_alphabetical)
    for _,file in pairs(healthy) do
        insert(file)
    end
end

M.tags = function()
    debugging("Adding tags (not implemented)")
end

return M
