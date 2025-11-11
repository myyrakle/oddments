//에러 타입 FooError 정의.
//인자는 문자열 하나 받음
exception FooError of string

//예외 투척
raise (FooError("흠"))
