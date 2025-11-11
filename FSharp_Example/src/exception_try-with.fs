//에러 타입 FooError 정의.
//인자는 문자열 하나 받음
exception FooError of string

try //여기서부턴 예외 발생가능
  //예외 투척
  raise (FooError("그냥 던져봄"))
with //예외 발생하면 받음
  | FooError(text) -> printfn "%s" text
