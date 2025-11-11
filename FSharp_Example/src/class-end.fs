type Person = class //클래스 시작
  //필드 선언
  val name: string
  val age: int
  
  //생성자
  new (_name:string, _age:int) =
    {name=_name; age=_age}
    
  //메서드
  member this.print() = 
    printfn "이름:%s,나이:%d" this.name this.age
    
  end //클래스 끝
  
let john = new Person("john", 14)
john.print()
