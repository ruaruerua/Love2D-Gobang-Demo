---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by wangboyuan.
--- DateTime: 2020/7/29 20:45
---
require("Function")
require('Component/Button')
require("Component/piece")
--local json = require("json")
--local Init = Class("Scences/Init")
--local test
function love.load() --资源加载回调函数，仅初始化时调用一次
    clearButton()
    clearPieces()
    love.graphics.setBackgroundColor(0.3,0.3,0.3)
    button:new(function()
        OnStart()
    end,"Start", 360, 300, 5, 5, {1,0,0}, love.graphics.newFont(30))

    button:new(function()
        OnQuit()
    end,"QUIT", 360, 400, 5, 5, {1,0,0}, love.graphics.newFont(30))

    --local a = {}
    --local key = "aa"
    --a[key] = "bb"
    --test = json.encode(a)
end

function OnStart()
    --SwitchScence('Hall')
    loginServer()
end

function OnQuit()
    love.event.quit ()
end

function love.update(dt) --更新回调函数，每周期调用
    updateButtons()
    updatePieces()
    if ServerHost ~= 0 and ServerHost ~= nil then
        print(ServerHost)
        SwitchScence('Hall')
        --love.graphics.print(ServerHost..MyID..Myscore, 50, 50)
    end
end




function love.draw() --绘图回调函数，每周期调用
    drawButtons()
    drawPieces()

    love.graphics.print('Gobang', 360, 100)
    --love.graphics.print(test, 360, 100)
end



function love.keypressed(key) --键盘检测回调函数，当键盘事件触发是调用





end



function love.mousepressed(x,y,key) --回调函数释放鼠标按钮时触发。





end

