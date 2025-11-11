# 별칭 부여
IntOrBool = Union{Int, Bool}

function main()
    num::IntOrBool = 10 # Int값 할당
    println(num)

    num = true # Bool도 할당 가능
    println(num)
end

main()
