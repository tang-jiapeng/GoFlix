local key=KEYS[1]
local limit=ARGV[1]
local offset=ARGV[2]


local exists=redis.call("EXISTS",key)
if exists==0 then
    return {3}
end

local status=0

local ex=redis.call("ZRange",key,0,0)
local ttl=redis.call("TTL",key)
if tonumber(ex[1])>=tonumber(ttl) then
    status=2
else
    status=1
end

local res=redis.call("ZRevRangeByScore",key,"+inf",0,"Limit",tonumber(offset),tonumber(limit))
table.insert(res,status)
return res