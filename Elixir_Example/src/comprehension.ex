list = for e <- 1..10, do: e*2
Enum.join(list, ", ") |> IO.puts

list = for e <- 1..10, e<5, do: e*2
Enum.join(list, ", ") |> IO.puts
