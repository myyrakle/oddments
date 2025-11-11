#include <string>

std::string operator*(const std::string& lhs, size_t rhs);
std::string operator*(size_t lhs, const std::string& rhs);

std::string operator*(const std::string& lhs, size_t rhs)
{
  std::string sum;
  sum.reserve(lhs.size()*rhs);
  
  for(int i = 0; i<rhs; i++)
  {
    sum.append(lhs);
  }
  
  return sum;
}

std::string operator*(size_t lhs, const std::string& rhs)
{
  return rhs*lhs;
}
