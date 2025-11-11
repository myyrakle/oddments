defmodule Test do
    def double(n), do: n * 2
    def add3(n), do: n + 3
    def half(n), do: n / 2
end

num = 4
num = Test.half(Test.add3(Test.double num))
IO.puts num

num = 4
num = Test.double(num) |> Test.add3() |> Test.half()
IO.puts num
