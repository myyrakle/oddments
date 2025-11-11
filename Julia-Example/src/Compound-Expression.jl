result = begin
    a = 10
    b = 20
    a + b 
end

println(result)

result = (a = 10; b = 20; c = 30; a+b+c)

println(result)
