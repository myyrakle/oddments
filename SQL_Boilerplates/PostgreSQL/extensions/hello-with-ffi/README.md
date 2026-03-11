# simpleext (PostgreSQL Extension)

현재 프로젝트에서 바로 빌드 가능한 가장 간단한 C++ PostgreSQL 확장 템플릿입니다.

## 파일
- `simpleext.control`
- `simpleext--1.0.sql`
- `simpleext.cpp`
- `Makefile`

## 설치 방법
1. `pg_simpleext` 폴더에서 빌드 후 설치
   ```sql
   make
   make install
   ```
2. PostgreSQL에 접속해서 설치
   ```sql
   CREATE EXTENSION simpleext;
   SELECT simpleext.hello_world(1); -- 2 (입력값 + 1)
   SELECT simpleext.add_ints(2, 3); -- 5
   ```

## 삭제
```sql
DROP EXTENSION simpleext;
```

## 참고
- `LANGUAGE C` 함수는 `simpleext.cpp`에서 `extern "C"` + `PG_FUNCTION_INFO_V1`로 선언된 엔트리 함수와 매핑됩니다.
