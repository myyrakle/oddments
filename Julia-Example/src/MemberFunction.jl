mutable struct Person
    name::String
    age::Int
    print::Function # 멤버함수 필드

    function Person(name::String, age::Int)
        this = new() # 객체 생성
        this.name = name
        this.age = age
        # 멤버함수 삽입
        this.print=function() Person__print(this) end
        this # 반환
    end
end

# 멤버함수 구현
function Person__print(p::Person)
    println("이름:$(p.name), 나이:$(p.age)")
end

p = Person("홍길동", 22) # 객체 생성
p.print()
