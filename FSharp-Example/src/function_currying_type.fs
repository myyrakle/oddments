let mul (lhs:int) (rhs:int):int = 
  lhs * rhs
  
//커링 함수 생성(lhs만 2로 고정)
let twice:(int->int) = mul 2

printfn "%d" (twice 10)
