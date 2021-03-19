local Scene = require("lib.Scene")
local Anim = require("lib.Animation")
local Sprite = require("lib.Sprite")

local T = Scene:derive("Test")

local hero_atlas

local spr
local idle = Anim(16, 16, 16, 16, 4, 4, 6 )
local walk = Anim(16, 32, 16, 16, 6, 6, 12)
local swim = Anim(16, 64, 16, 16, 6, 6, 12)
local punch = Anim(16, 80, 16, 16, 3, 3, 10, false)
local snd

function T:new(scene_mgr) 
    self.super:new(scene_mgr)

    hero_atlas = love.graphics.newImage("assets/gfx/hero.png")
    spr = Sprite(hero_atlas, 16,16, 100,100, 10, 10)
    spr:add_animations({idle = idle, walk = walk, swim = swim, punch = punch})
    spr:animate("swim")

    snd = love.audio.newSource("assets/sfx/hit01.wav", "static")
end

local entered = false
function T:enter()
    if not entered then
        entered = true
        print("enter test")
    end
end

function T:update(dt)
    -- if Key:key_down("space") or spr.current_anim ~= "punch" then
    --     spr:animate("punch")
    --     love.audio.stop(snd)
    --     love.audio.play(snd)
    -- elseif Key:key_down("escape") then
    --     love.event.quit()
    -- end   

    -- if spr.current_anim == "punch" and spr:animation_finished() then
    --     spr:animate("idle")
    -- end
    spr:update(dt)

end

function T:draw()
    love.graphics.clear(64,64,255)
    spr:draw()
end

return T
