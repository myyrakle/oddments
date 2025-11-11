#include <iostream>

class Foo {
  template<class T> 
  void foo() {
    puts("fo3o");
  }
};

int main() {
  Foo obj;
  
  obj.template foo<int>();
  return 0;
}
