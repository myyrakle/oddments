type Person(_name, _age) =
  //필드 선언
  let name: string = _name
  let age: int = _age
    
  //메서드
  member this.print() = 
    printfn "이름:%s,나이:%d" name age
    
let john = new Person("john", 14)
john.print()
    
