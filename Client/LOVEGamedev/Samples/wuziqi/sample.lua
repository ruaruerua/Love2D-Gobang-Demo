
local network = require("network")
-- the address and port of the server
local address, port = "localhost", 12345

-- love.load, hopefully you are familiar with it from the callbacks tutorial
function love.load()
    print("test")
    network.connet(address, port)
    network.connectedFunc = function()
        print("connect ok!")
    end
    network.disconnectFunc = function()
        print("disconnect!")
    end
    network.recvMsgFunc = function(msg)
        print("recv=>",msg," size=",string.len(msg))
    end

    t = 0 -- (re)set t to 0
end

-- love.update, hopefully you are familiar with it from the callbacks tutorial
function love.update(deltatime)
    network.update()
    t = t + 1
    if t%10==0 then
        local dg = string.format("%d", t)
        network.send(dg)
    end

end

-- love.draw, hopefully you are familiar with it from the callbacks tutorial
function love.draw()

end

-- And thats the end of the udp client example.


