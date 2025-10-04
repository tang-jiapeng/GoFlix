local key=KEYS[1]
local ttl=KEYS[2]
local data=ARGV

local exists=redis.call("EXISTS",key)
if exists==1 then
    local ex=redis.call("ZRange",key,1,1)
    local last=redis.call("TTL",key)
    if tonumber(ex[1])>=tonumber(last) then
        redis.call("DEL",key)
    else
        return
    end
end

for i=1,#data,2
do redis.call("ZAdd",key,data[i],data[i+1])
end

redis.call("EXPIRE",key,ttl)

return