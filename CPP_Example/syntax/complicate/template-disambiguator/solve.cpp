#include <iostream>
  2
  3 template<class Boom>
  4 struct Foo {
  5   template<class T>
  6   void foo() {
  7     puts("fo3o");
  8   }
  9 };
 10
 11 template <class T>
 12 void run() {
 13   Foo<T> obj;
 14   obj.template foo<T>();
 15 }
 16
 17 int main() {
 18   run<int>();
 19
 20   return 0;
 21 }
