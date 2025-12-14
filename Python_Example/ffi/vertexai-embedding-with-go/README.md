# Vertexai Go FFI

- vertexai의 python 모듈의 비정상적인 메모리 사용량을 회피하기 위한 FFI 구성 예제
- 그나마 정상적인 구조를 가진 Go SDK를 사용해서 필요한 부분을 shared library로 컴파일하고, 그걸 Python에서 사용하는 예제
- 개선 효과: 240mb (python sdk) => 30mb (.so with go)

## Setup

- go 필요
- `sh go-ffi/build.sh` 실행해서 컴파일 => shared library 생성
- `uv run test.py` 사용
