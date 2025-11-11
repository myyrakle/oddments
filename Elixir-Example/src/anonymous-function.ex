callback = fn(e) -> e+1 end
nums = [1, 2, 3, 4, 5]
nums = Enum.map nums, callback
str = Enum.join(nums, ", ")
IO.puts str

callback = &(&1 + 1)
nums = [1, 2, 3, 4, 5]
nums = Enum.map nums, callback
str = Enum.join(nums, ", ")
IO.puts str
