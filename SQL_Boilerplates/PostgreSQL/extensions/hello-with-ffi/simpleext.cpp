extern "C"
{
#include "postgres.h"
#include "fmgr.h"

PG_MODULE_MAGIC;

PG_FUNCTION_INFO_V1(simpleext_add_one);
PG_FUNCTION_INFO_V1(simpleext_add_ints);

Datum
simpleext_add_one(PG_FUNCTION_ARGS)
{
    int32 x = PG_GETARG_INT32(0);
    PG_RETURN_INT32(x + 1);
}

Datum
simpleext_add_ints(PG_FUNCTION_ARGS)
{
    int32 a = PG_GETARG_INT32(0);
    int32 b = PG_GETARG_INT32(1);
    PG_RETURN_INT32(a + b);
}
}
