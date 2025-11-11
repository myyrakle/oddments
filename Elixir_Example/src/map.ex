map = %{ "john"=>3, "tom"=>44 }
IO.puts map["john"]

new_map = Map.put_new map, "홍길동", 999
IO.puts new_map["홍길동"]

new_map = %{ new_map | "홍길동"=> 8888 }
IO.puts new_map["홍길동"]
