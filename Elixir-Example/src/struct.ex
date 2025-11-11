defmodule Person do
    defstruct name: "", age: 0
end

defmodule Test do
    def test do
        john = %Person{name: "john", age: 20}                john2 = %{john | age: 99, name: "john2"} # update
        IO.puts john2.name
        IO.puts john2.age
    end
end

Test.test
