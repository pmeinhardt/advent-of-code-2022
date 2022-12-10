local cpu = {cycle=1, x=1, fn=nil, t=0}

function cpu:execute (input)
  return function ()
    if self:isidle() then
      local instruction = input()
      if not instruction then return nil end
      self:read(instruction)
    end

    self:exec()
    self:tick()

    return self
  end
end

function cpu:read (instruction)
  local tokens = string.gmatch(instruction, "([^%s]+)")
  local op = tokens()

  local fn, t

  if op == "addx" then
    local value = tonumber(tokens())
    t, fn = 1, function () self.x = self.x + value end
  elseif op == "noop" then
    t, fn = 0, function () end
  else
    error(string.format("error in instruction: %s", instruction))
  end

  self.fn, self.t = fn, t
end

function cpu:tick ()
  self.cycle = self.cycle + 1
  self.t = math.max(0, self.t - 1)
end

function cpu:exec ()
  if self.fn and self.t == 0 then
    self.fn()
    self.fn = nil
  end
end

function cpu:isidle ()
  return self.fn == nil
end

return cpu
