function just_print(value::T) where {T}
    println(value)
end

function main()
    just_print(100::Int)
    just_print("foo"::String)
end

main()
