---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by wangboyuan.
--- DateTime: 2020/7/29 20:43
---
local http = require "socket.http"


local Server = "192.168.92.204:9999"
--local Server = "127.0.0.1:9999"
ServerHost = 0
MyID = 0
Myscore = 0
MyP = "B"
canBeSet = true
StepNum = 0
fP = {}
sP = {}
MyPlayer = {}
RoomID = 0



--切换场景
function SwitchScence(scence)
    -- 将重要的函数赋予空值，以免冲突
    love.update = nil
    love.draw = nil
    love.keypressed = nil

    -- 将需要的场景加载进来，并执行load函数
    love.filesystem.load ('Scences/'..scence..'.lua') ()
    love.load ()
end

function loginServer()
    local response_body = {}
    local res, code, response_headers = http.request{
        url = "http://"..Server,
        method = "POST",
        sink = ltn12.sink.table(response_body),
    }
    local txt = table.concat(response_body)
    local v = string.split(txt,"|")
    ServerHost = v[1]
    MyID = v[2]
    Myscore = v[3]
end

function ReSetData()
    --ServerHost = 0
    --MyID = 0
    --Myscore = 0

    canBeSet = true
    MyP = "B"
    StepNum = 0
    fP = {}
    sP = {}
    RoomID = 0
end

function string.split(input, delimiter)
    input = tostring(input)
    delimiter = tostring(delimiter)
    if (delimiter=='') then return false end
    local pos,arr = 0, {}
    local str
    -- for each divider found
    for st,sp in function() return string.find(input, delimiter, pos, true) end do
        str = string.sub(input, pos, st - 1)
        if (str~='') then
            table.insert(arr, str)
        end
        pos = sp + 1
    end
    str = string.sub(input, pos)
    if str~='' then
        table.insert(arr, str)
    end
    return arr
end