---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by wangboyuan.
--- DateTime: 2020/8/2 16:56
---
local socket = require("socket")
local json = require("json")
--local t = 0

local tcpnetwork = {}

while true do
    input = io.read()
    if #input > 0 then
        assert(sock:send(input .. "\n"))
    end
    local recvt, sendt, status = socket.select({sock}, nil, 1)
    while #recvt > 0 do
        local response, receive_status = sock:receive()
        if receive_status ~= "closed" then
            if response then
                print("-------------")
                print(response)
                recvt, sendt, status = socket.select({sock}, nil, 1)
            end
        else
            break
        end
    end
end

TCP_MSG = {
    Connect = "connect",
    Connected = "connected",
    Disconnect = "disconnect",
    Message = "msg"
}

function tcpnetwork.connet(address,port,connectedFunc,disconnectFunc,recvMsgFunc)
    --udp = socket.udp()
    --tcp = socket.tcp()
    tcp = socket.connect(address, port)
    tcpnetwork.connectedFunc = connectedFunc
    tcpnetwork.disconnectFunc = disconnectFunc
    tcpnetwork.recvMsgFunc = recvMsgFunc
    tcpnetwork.isConnect = false

    tcp:settimeout(0)

    --udp:setpeername(address, port)

    math.randomseed(os.time())

    --udp:send(UDP_MSG.Connect) -- the magic line in question.
    tcp:send(TCP_MSG.connect .. "\n")
end

function tcpnetwork.update()

    -- there could well be more than one message waiting for us, so we'll
    -- loop until we run out!
    repeat
        -- and here is something new, the much anticipated other end of udp:send!
        -- receive return a waiting packet (or nil, and an error message).
        -- data is a string, the payload of the far-end's send. we can deal with it
        -- the same ways we could deal with any other string in lua (needless to
        -- say, getting familiar with lua's string handling functions is a must.
        data, msg = udp:receive()

        local recvt, sendt, status = socket.select({tcp}, nil, 1)
        while #recvt > 0 do
            local response, receive_status = tcp:receive()
            if receive_status ~= "closed" then
                if response then
                    print("-------------")
                    print(response)
                    recvt, sendt, status = socket.select({tcp}, nil, 1)
                end
            else
                break
            end
        end

        if data then -- you remember, right? that all values in lua evaluate as true, save nil and false?
            print("recv data:",msg,data)
            if data==TCP_MSG.Connected then
                tcpnetwork.connected()
            elseif data==TCP_MSG.Disconnect then
                tcpnetwork.disconnect()
            end
            if tcpnetwork.recvMsgFunc then
                tcpnetwork.recvMsgFunc(data)
            end
        elseif msg ~= 'timeout' then
            error("Network error: "..tostring(msg))
        end
    until not data

end

function tcpnetwork.connected()
    print("connected")
    tcpnetwork.isConnect = true
    if tcpnetwork.connectedFunc then
        tcpnetwork.connectedFunc()
    end
end

function tcpnetwork.disconnect()
    tcpnetwork.isConnect = false
    if tcpnetwork.disconnectFunc then
        tcpnetwork.disconnectFunc()
    end
end


function tcpnetwork.send(msg)
    if not tcpnetwork.isConnect then
        print("network is not connected")
        return
    end
    tcp:send(TCP_MSG.Message.."|"..msg)
end

return tcpnetwork