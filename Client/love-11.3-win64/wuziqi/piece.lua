---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by wangboyuan.
--- DateTime: 2020/7/30 12:37
---
local class = require("middleclass")
piece = class("piece")
local pieces = {}
--local originalFont = love.graphics.getFont()

function piece:initialize(x, y,color)
    self.x = x
    self.y = y
    self.color = color or {150,150,150}
    self.originalColor = self.color

    self.beSetted = false

    self.id = { self.x , self.y }

    table.insert(pieces, self)
    return self
end

function piece:update()
    local x, y = love.mouse.getX(), love.mouse.getY()
    if x < self.x + 20 and x > self.x and y < self.y + 20 and y > self.y then
        if love.mouse.isDown(1) then
            self.beSetted = true
        end
        self.color = {self.color[1] + 20, self.color[2] + 20, self.color[3] + 20}
    else
        self.color = self.originalColor
    end
end

function piece:draw()

    if   self.beSetted then
        love.graphics.setColor({ 0,0,0 })
        love.graphics.circle("fill", self.x+10, self.y+10, 10) -- Draw white circle with 100 segments.
    end


end

function updatePieces()
    for i, v in pairs(pieces) do
        v:update()
    end
end

function drawPieces()
    for i, v in pairs(pieces) do
        v:draw()
    end
end

function clearPieces()
    pieces = {}
end