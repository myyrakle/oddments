class MyError implements Exception
{
    errMsg() => "으악";
}

main()
{
    //예외 발생구역
    try
    {
        throw new MyError();
    }
    on MyError catch(e) //예외 발생시 이동
    {
        print("내 예외");
        print(e.errMsg());

    }
    finally
    {
        print("종료");
    }

    print('foo');
}
