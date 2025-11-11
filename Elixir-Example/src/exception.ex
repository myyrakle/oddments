defmodule MyError do
    defexception message: "오류"
end

try do
    raise MyError
rescue
    e in MyError -> IO.puts e.message
    e in RuntimeError -> IO.puts e.message
end

IO.puts "테스트"
