#define generic_equal(lhs, rhs) \
_Generic((lhs), \
 const char*: _Generic((rhs), const char*: (strcmp(lhs,rhs)==0), default: false), \
 default: (lhs==rhs))
 
#define generic_not_equal(lhs, rhs) \
_Generic((lhs), \
 const char*: _Generic((rhs), const char*: (strcmp(lhs,rhs)!=0), default: false), \
 default: (lhs!=rhs))
 
#define generic_less(lhs, rhs) \
_Generic((lhs), \
 const char*: _Generic((rhs), const char*: (strcmp(lhs,rhs)<0), default: false), \
 default: (lhs<rhs))
 
#define generic_less_equal(lhs, rhs) \
_Generic((lhs), \
 const char*: _Generic((rhs), const char*: (strcmp(lhs,rhs)<=0), default: false), \
 default: (lhs<=rhs))
 
#define generic_greater(lhs, rhs) \
_Generic((lhs), \
 const char*: _Generic((rhs), const char*: (strcmp(lhs,rhs)>0), default: false), \
 default: (lhs>rhs))
 
#define generic_greater_equal(lhs, rhs) \
_Generic((lhs), \
 const char*: _Generic((rhs), const char*: (strcmp(lhs,rhs)>=0), default: false), \
 default: (lhs>=rhs))
 
 
#define generic_compare(lhs,rhs) \
_Generic((lhs), \
 const char*: _Generic((rhs), const char*: (strcmp(lhs,rhs)), default: assert(false)), \
 default: (lhs-rhs)))
