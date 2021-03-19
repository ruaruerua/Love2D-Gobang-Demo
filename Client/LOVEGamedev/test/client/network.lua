local socket = require "socket"

UDP_MSG = {
    Connect = "connect",
    Connected = "connected",
    Disconnect = "disconnect",
    Message = "msg"
}
local network = {}
function network.connet(address,port,connectedFunc,disconnectFunc,recvMsgFunc)
    udp = socket.udp()
    network.connectedFunc = connectedFunc
    network.disconnectFunc = disconnectFunc
    network.recvMsgFunc = recvMsgFunc
    network.isConnect = false
    -- normally socket reads block until they have data, or a
    -- certain amout of time passes.
    -- that doesn't suit us, so we tell it not to do that by setting the
    -- 'timeout' to zero
    udp:settimeout(0)

    -- unlike the server, we'll just be talking to the one machine,
    -- so we'll "connect" this socket to the server's address and port
    -- using setpeername.
    --
    -- [NOTE: UDP is actually connectionless, this is purely a convenience
    -- provided by the socket library, it doesn't actually change the
    --'bits on the wire', and in-fact we can change / remove this at any time.]
    udp:setpeername(address, port)

    -- seed the random number generator, so we don't just get the
    -- same numbers each time.
    math.randomseed(os.time())

    -- thats...it, really. the rest of this is just putting this context and practical use.
    udp:send(UDP_MSG.Connect) -- the magic line in question.
end

function network.update()

    -- there could well be more than one message waiting for us, so we'll
    -- loop until we run out!
    repeat
        -- and here is something new, the much anticipated other end of udp:send!
        -- receive return a waiting packet (or nil, and an error message).
        -- data is a string, the payload of the far-end's send. we can deal with it
        -- the same ways we could deal with any other string in lua (needless to
        -- say, getting familiar with lua's string handling functions is a must.
        data, msg = udp:receive()

        if data then -- you remember, right? that all values in lua evaluate as true, save nil and false?
            print("recv data:",msg,data)
            if data==UDP_MSG.Connected then
                network.connected()
            elseif data==UDP_MSG.Disconnect then
                network.disconnect()
            end
            if network.recvMsgFunc then
                network.recvMsgFunc(data)
            end
        elseif msg ~= 'timeout' then
            error("Network error: "..tostring(msg))
        end
    until not data

end

function network.connected()
    print("connected")
    network.isConnect = true
    if network.connectedFunc then
        network.connectedFunc()
    end
end

function network.disconnect()
    network.isConnect = false
    if network.disconnectFunc then
        network.disconnectFunc()
    end
end


function network.send(msg)
    if not network.isConnect then
        print("network is not connected")
        return
    end
    udp:send(UDP_MSG.Message.."|"..msg)
end

return network