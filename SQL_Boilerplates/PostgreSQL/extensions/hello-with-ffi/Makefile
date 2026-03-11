MODULE_big = simpleext
OBJS = simpleext.o
EXTENSION = simpleext
DATA = simpleext--1.0.sql

PG_CXXFLAGS = -std=c++17
SHLIB_LINK += -lstdc++

PG_CONFIG = pg_config
PGXS := $(shell $(PG_CONFIG) --pgxs)
include $(PGXS)
