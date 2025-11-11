//에러 타입 FooError 정의.
//인자는 문자열 하나 받음
exception FooError of string

//에러 타입 BarError 정이
exception BarError of string

try //여기서부턴 예외 발생가능
  //예외 투척
  raise (BarError("그냥 던져봄"))
with //예외 발생하면 해당 예외 타입의 절로 이동
  | FooError(text) -> printfn "Foo:%s" text
  | BarError(text) -> printfn "Bar:%s" text
