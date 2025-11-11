#include <iostream>

class Integer {
private: 
    int value = 0;
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
};

enum class CompareOperator
{
    LESS_THAN,
    LESS_THAN_EQUAL,
    GREATER_THAN,
    GREATER_THAN_EQUAL,
};

class CompareResult 
{
public:
    Integer lhs;
    Integer rhs;
    CompareOperator op;
public: 
    CompareResult(const Integer& lhs, const Integer& rhs, CompareOperator op): lhs(lhs), rhs(rhs), op(op)
    {}
    operator bool() const 
    {
        switch(op) 
        {
            case CompareOperator::LESS_THAN :
                return (int)lhs < (int)rhs;
            case CompareOperator::LESS_THAN_EQUAL :
                return (int)lhs <= (int)rhs;
            case CompareOperator::GREATER_THAN :
                return (int)lhs > (int)rhs;
            case CompareOperator::GREATER_THAN_EQUAL :
                return (int)lhs >= (int)rhs;
        }
        
        return false;
    }
};

CompareResult operator<(const Integer& lhs, const Integer& rhs) {
    return CompareResult(lhs, rhs, CompareOperator::LESS_THAN);
}

CompareResult operator>(const Integer& lhs, const Integer& rhs) {
    return CompareResult(lhs, rhs, CompareOperator::GREATER_THAN);
}

CompareResult operator<=(const Integer& lhs, const Integer& rhs) {
    return CompareResult(lhs, rhs, CompareOperator::LESS_THAN_EQUAL);
}

CompareResult operator>=(const Integer& lhs, const Integer& rhs) {
    return CompareResult(lhs, rhs, CompareOperator::GREATER_THAN_EQUAL);
}

bool operator<(const CompareResult& lhs, const Integer& rhs) {
    return (bool)lhs && (lhs.rhs < rhs);
}

bool operator>(const CompareResult& lhs, const Integer& rhs) {
    return (bool)lhs && (lhs.rhs > rhs);
}

bool operator<=(const CompareResult& lhs, const Integer& rhs) {
    return (bool)lhs && (lhs.rhs <= rhs);
}

bool operator>=(const CompareResult& lhs, const Integer& rhs) {
    return (bool)lhs && (lhs.rhs >= rhs);
}

using namespace std;

int main()
{
    Integer num1 = 10;
    Integer num2 = 20;
    Integer num3 = 30;
    
    if(num1 < num2 < num3) { // 10 < 20 < 30
        puts("붐");
    }
    
    if(num1 > num2 > num3) { // 10 > 20 > 30
        puts("안됨");
    }
    
    if(num1 >= num1 < num3) { // 10 >= 10 < 30
        puts("됨");
    }

    return 0;
}
