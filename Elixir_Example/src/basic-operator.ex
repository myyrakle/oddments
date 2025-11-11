IO.puts 10+3 #덧셈
IO.puts 10-3 #뺄셈
IO.puts 10*3 #곱셈
IO.puts 10/3 #나눗셈
IO.puts div(10, 3) #정수 나눗셈
IO.puts rem(10, 3) #나머지셈

IO.inspect [3, 4] ++ [5, 1] #리스트 연결
IO.inspect [1, 2, 3, 4] -- [1] # 리스트 뺄셈

IO.puts 10 == 10.0
IO.puts 10 != "10"
IO.puts 10 === 10.0
IO.puts 10 !== 10.0
IO.puts 10 > 5
IO.puts 10 < 5
IO.puts 13 >= 5
IO.puts 13 <= 5

IO.puts true and true
IO.puts true && true
IO.puts false or true
IO.puts false || true
IO.puts not true
IO.puts !true

IO.puts 5 in [1, 3, 5]
IO.puts 6 not in [1, 3]
