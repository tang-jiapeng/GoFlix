local key=KEYS[1]
local value=ARGV[1]

local res=redis.call("Get",key)

if res==nil then
    return "lock has been delete"
end

if res~=value then
    return "input value is not the current lock value"
end

redis.call("Del",key)

return nil
