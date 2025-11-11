defmodule Test do
    def test do
        IO.puts "붐"
    end
end

alias Test, as: Foo

Test.test
Foo.test
