# try-catch
try 
    throw(DomainError("그냥 에러"))
catch e
    println("예외 던져짐: $e")
end


# try-finally
try 
    throw(DomainError("그냥 에러"))
finally
    println("종료됨")
end
