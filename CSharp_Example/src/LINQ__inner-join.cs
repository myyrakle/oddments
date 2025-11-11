using System;
using System.Collections.Generic;
using System.Linq;

//사람 클래스
//이름과 나이만 저장
class Person
{
  public string name;
  public int age;
  public Person(String name, int age)
  {
    this.name = name;
    this.age = age;
  }
}

//학생 클래스
//이름과 학교명만 저장
class Student
{
  public string name;
  public string school;
  public Student(string name, string school)
  {
    this.name = name;
    this.school = school;
  }
}
 
namespace Dcoder
{
  public class Program
  {
    public static void Main(string[] args)
    {
      Person[] persons = {
        new Person("john", 15),
        new Person("tom", 20),
        new Person("anna", 13),
        new Person("foo", 22)
      };
      
      Student[] students = {
        new Student("john", "조선대"),
        new Student("tom", "엄석대"),
        new Student("anna", "예일대"),
        new Student("bar", "첨성대")
      };
      
      var result = 
        from person in persons
        join student in students 
        on person.name equals student.name
        select new {
          name = person.name,
          age = person.age,
          school = student.school
        };
       
      foreach(var e in result)
      {
        Console.Write($"이름:{e.name}, ");
        Console.Write($"나이:{e.age}, ");
        Console.Write($"학교:{e.school}\n");
      }
    }
  }
}
    
    
    
    
    
    
