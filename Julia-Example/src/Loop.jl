# while statement
i = 0

while i<5
    println("BOOM!!")
    global i+=1
end


# for statement
for i = 1:5
    println("BOOM!! $i")
end


# for-in statement
arr = ["foo", "bar", "BOOM!!!", "HO"]

for e in arr
    println(e)
end


# nested for (9x9 multiply)
for i = 1:9, j = 1:9 
    println("$(i)X$(j)=$(i*j)")
end


# nested for-in
for b ∈ [true, false], e ∈ [1,2,3,4]
    println("$b:$e")
end
