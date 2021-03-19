

function love.load()
	love.graphics.setColor(255,255,255,255)
	
	idle = {}
	idle.width,idle.height=137,143
	idle.img = love.graphics.newImage( "spirits/idle.png" )
	idle.frame = {}
	idle.ox = 0
	idle.oy = -60
	for i=1,4,1 do
		table.insert(idle.frame,love.graphics.newQuad(0,idle.height*(i-1),idle.width,idle.height,idle.img:getDimensions()))
	end
	
	run = {}
	run.width,run.height=210,210
	run.img = love.graphics.newImage( "spirits/run.png" )
	run.frame = {}
	-- 105,117,168,207,247,326,361,387,411,444,467
	-- run.frameOffset = {23,12,51,39,40,79,55,26,24,33}
	run.frameTime = {0.04,0.04,0.07,0.12,0.12,0.1,0.08,0.08,0.06,0.04}
	run.frameSpeed = {100,180,240,300,260,260,240,140,220,180}
	run.ox = 0
	run.oy = 0
	for i=1,10,1 do
		table.insert(run.frame,love.graphics.newQuad(0,run.height*(i-1),run.width,run.height,run.img:getDimensions()))
	end
	
	cat = {x=110,y=110,statu="idle", direction=1, frame=idle, frameN=1}
end

t = 0
function love.update(dt)
	t = t+dt
	if love.keyboard.isDown('a') then
		if cat.statu ~= "run" or cat.direction ~= -1 then
			cat.frame = run
			cat.frameN = 1
		end
		cat.direction = -1
		cat.statu = "run"
	elseif love.keyboard.isDown('d') then
		if cat.statu ~= "run" or cat.direction ~= 1 then
			cat.frame = run
			cat.frameN = 1
		end
		cat.direction = 1
		cat.statu = "run"
	else
		if cat.statu ~= "idle" then
			cat.frame = idle
			cat.frameN = 1
		end
		cat.statu = "idle"
	end
	
	if cat.statu == "idle" then
		if t>0.2 then
			t = 0
			cat.frameN = cat.frameN+1
			if cat.frame.frame[cat.frameN]==nil then
				cat.frameN = 1
			end
		end
	elseif cat.statu == "run" then
		if t>cat.frame.frameTime[cat.frameN] then
			t = 0
			cat.frameN = cat.frameN+1
			if cat.frame.frame[cat.frameN]==nil then
				cat.frameN = 1
			end
		end
		cat.x = cat.x + 2*cat.direction*cat.frame.frameSpeed[cat.frameN]*dt
	end
end
 
function love.draw()
	love.graphics.draw(cat.frame.img,cat.frame.frame[cat.frameN],cat.x-cat.frame.width*cat.direction/2,cat.y,0,cat.direction,1,cat.frame.ox,cat.frame.oy)
end

