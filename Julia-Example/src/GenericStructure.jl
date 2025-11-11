# 제너릭 구조체
struct Wrap{T}
    value::T
end

function main()
    # Int로 구체화
    num :: Wrap{Int} = Wrap{Int}(100)
    println(num.value)

    # Float64로 구체화
    fnum :: Wrap{Float64} = Wrap{Float64}(0.133)
    println(fnum.value)
end

main()
