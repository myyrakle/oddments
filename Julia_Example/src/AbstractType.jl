# 추상 타입 Runnable 선언
abstract type Runnable end

# Runnable의 하위 타입 BoomRun 정의
struct BoomRun <: Runnable
    run::Function

    function BoomRun()
        new(function() println("BOOM!!!") end)
    end
end

function main()
    runner::Runnable = BoomRun()
    runner.run()
end

main()
