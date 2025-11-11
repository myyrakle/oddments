struct Person
    name::String
    age::Int
end

p = Person("john", 10) # create object
println("이름:$(p.name), 나이:$(p.age)")
