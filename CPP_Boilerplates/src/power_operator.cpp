#include <iostream>
#include <cmath>

using std::cout;

class Integer {
private: 
    int value = 0;
public:
    struct Power{
        int num;
        Power(int n): num(n) {}
    };
public:
    using Self = Integer;
    Integer(int init): value(init) 
    {}
    ~Integer() = default;
    Integer(const Self&) = default;
    Integer(Self&&) = default;
    Integer& operator=(const Self&) = default;
    Integer& operator=(Self&&) = default;
public:
    operator const int&() const
    {
        return value;
    }
    operator int&() 
    {
        return value;
    }
    Power operator*() const {
        return Power(this->value);
    }
    Integer operator*(const Power& rhs) const {
        return pow(this->value, rhs.num);
    }
};

int main()
{
    Integer a = 5;
    Integer b = 3;
    Integer c = a**b;
    cout<<c;
}
