defmodule Person do
    defstruct name: "", age: 0
end

defprotocol Print do
    def print(data)
end

defimpl Print, for: Person do
    def print(data) do
        IO.puts data.name
        IO.puts data.age
    end
end

defmodule Test do
    def test do
        john = %Person{name: "john", age: 20}
        Print.print john
    end
end

Test.test
