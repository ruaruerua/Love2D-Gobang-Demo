---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by wangboyuan.
--- DateTime: 2020/8/3 11:38
---
local network = require("network")
require("Component/Button")
require("Component/Player")
require("Component/piece")
require("Component/DrawnWindow")
require("Function")
local json = require("json")
local ranklist = {}
local debugtext = "Hall"
local isMarching = false

local refreshT = 0

function love.load() --资源加载回调函数，仅初始化时调用一次
    clearButton()
    clearPieces()
    clearplayer()
    clearDrawnWindow()
    love.graphics.setBackgroundColor(0.3,0.3,0.3)
    isMarching = false
    button:new(function()
        if not isMarching then
            OnMarch()
        end
    end,"Match", 600, 500, 5, 5, {1,0,0}, love.graphics.newFont(30))

    DrawnWindow:new("test")

    local _ip = string.split(ServerHost,":")
    network.connet(_ip[1],_ip[2])
    print(_ip[1]..":".._ip[2])
    network.recvMsgFunc = function(data)
        HALLCallBack(data)
    end

end

function love.update(dt) --更新回调函数，每周期调用
    updateButtons()
    updateplayer()
    updateDrawnWindow()
    network.update()

    --if refreshT % 1000 == 0 then
    --    --SetPlayer()
    --    local _json = {}
    --    _json.CMD =
    --    network.send(_json,8)
    --end

    --if refreshT > 10000 then
    --    refreshT = 1
    --end
    --refreshT = refreshT + 1
end


function love.draw() --绘图回调函数，每周期调用
    love.graphics.print("ServerHost"..ServerHost.."MyID"..MyID.."Myscore"..Myscore, 10, 10)
    love.graphics.print(debugtext,  love.graphics.newFont(30),600, 10)
    drawButtons()
    love.graphics.print("Players",  love.graphics.newFont(20),400, 50)
    drawplayer()
    love.graphics.print("Rank",  love.graphics.newFont(20),600, 50)

    drawDrawnWindow()
end


function love.keypressed(key) --键盘检测回调函数，当键盘事件触发是调用
end



function love.mousepressed(x,y,key) --回调函数释放鼠标按钮时触发。
end

--region 界面逻辑
function DrawRank()
    local x = 600
    local y = 10
    for i, v in pairs(ranklist) do
        love.graphics.print(y)
        y = y + 20
    end
end

function SetPlayer()
    clearplayer()
    local x = 400
    local y = 100
    --for i, v in pairs() do
    --    player:new()
    --end
    for i = 1, 10 do
        player:new(i,x,y)
        y = y + 40
    end
end

--endregion

--region 按钮逻辑
function OnMarch()
    debugtext = "isMarching!!!"
    network.send("",4)
    isMarching = true
end
--endregion

--region callback
function HALLCallBack(data)
    local _json = data
    if _json.CMD == 7 then
        local body = json.decode(_json.Body)
        fP = body.FirstPlayer
        sP = body.LastPlayer
        RoomID = body.RoomID
        --print(fP.."   "..sP.."  "..RoomID)
        if fP.ID == MyID then
            MyPlayer = fP
            PieceBeSetOn()
            MyP = "B"
        else
            MyPlayer = sP
            PieceBeSetOff()
            MyP = "W"
        end
        isMarching = false
        SwitchScence("Game")
    end

    --if _json.CMD == 8 then
    --    SetPlayer()
    --end

end


--endregion
