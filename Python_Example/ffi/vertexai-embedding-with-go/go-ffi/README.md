# Vertex AI Multimodal Embedding - Go Shared Library

이 프로젝트는 Google Cloud Vertex AI의 Multimodal Embedding API를 Go로 구현하고, Python에서 사용할 수 있도록 C shared library로 컴파일한 것입니다.

## 특징

- **고성능**: Go로 구현되어 빠른 성능
- **Text Embedding**: 텍스트를 벡터로 변환
- **Image Embedding**: 이미지를 벡터로 변환
- **Python 통합**: Python에서 ctypes를 통해 쉽게 사용

## 빌드 방법

### 필요 사항

- Go 1.21 이상
- CGO 활성화 (기본적으로 활성화되어 있음)
- GCC 또는 다른 C 컴파일러

### 컴파일

```bash
cd vertexai
CGO_ENABLED=1 go build -buildmode=c-shared -o libvertexai.so lib.go
```

이 명령은 다음 파일들을 생성합니다:
- `libvertexai.so`: shared library (Linux)
- `libvertexai.h`: C 헤더 파일

### 다른 플랫폼용 빌드

**macOS:**
```bash
CGO_ENABLED=1 go build -buildmode=c-shared -o libvertexai.dylib lib.go
```

**Windows:**
```bash
CGO_ENABLED=1 go build -buildmode=c-shared -o libvertexai.dll lib.go
```

## 사용 방법

### Python에서 사용

```python
from vertexai_client import VertexAIEmbedding

# 클라이언트 초기화
client = VertexAIEmbedding(
    project_id="your-gcp-project-id",
    location="us-central1",
    credentials_json="/path/to/credentials.json"  # 선택사항
)

# 텍스트 임베딩 가져오기
text = "Hello, world!"
text_embedding = client.get_text_embedding(text)
print(f"Embedding dimension: {len(text_embedding)}")

# 이미지 임베딩 가져오기
image_embedding = client.get_image_embedding(image_path="image.jpg")
print(f"Embedding dimension: {len(image_embedding)}")
```

### 인증 설정

다음 방법 중 하나로 GCP 인증을 설정할 수 있습니다:

**1. 환경 변수 사용 (권장)**
```bash
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/credentials.json"
```

**2. 코드에서 직접 지정**
```python
client = VertexAIEmbedding(
    project_id="your-project-id",
    credentials_json="/path/to/credentials.json"
)
```

**3. JSON 문자열로 직접 전달**
```python
import json

credentials = {
    "type": "service_account",
    "project_id": "...",
    # ... 나머지 필드
}

client = VertexAIEmbedding(
    project_id="your-project-id",
    credentials_json=json.dumps(credentials)
)
```

## API 레퍼런스

### Go 함수

#### GetTextEmbedding
```go
//export GetTextEmbedding
func GetTextEmbedding(projectID *C.char, location *C.char, text *C.char, credentialsJSON *C.char) *C.char
```

텍스트의 임베딩 벡터를 반환합니다.

**Parameters:**
- `projectID`: GCP 프로젝트 ID
- `location`: GCP 리전 (예: "us-central1")
- `text`: 임베딩할 텍스트
- `credentialsJSON`: 인증 정보 JSON 문자열 (빈 문자열이면 기본 인증 사용)

**Returns:**
JSON 문자열:
```json
{
  "text_embedding": [0.1, 0.2, ...],
  "error": ""
}
```

#### GetImageEmbedding
```go
//export GetImageEmbedding
func GetImageEmbedding(projectID *C.char, location *C.char, imageBase64 *C.char, credentialsJSON *C.char) *C.char
```

이미지의 임베딩 벡터를 반환합니다.

**Parameters:**
- `projectID`: GCP 프로젝트 ID
- `location`: GCP 리전
- `imageBase64`: Base64로 인코딩된 이미지 데이터
- `credentialsJSON`: 인증 정보 JSON 문자열

**Returns:**
JSON 문자열:
```json
{
  "image_embedding": [0.1, 0.2, ...],
  "error": ""
}
```

### Python 클래스

#### VertexAIEmbedding

**`__init__(project_id, location="us-central1", credentials_json=None)`**

클라이언트를 초기화합니다.

**`get_text_embedding(text: str) -> List[float]`**

텍스트 임베딩을 가져옵니다.

**`get_image_embedding(image_path: str = None, image_bytes: bytes = None) -> List[float]`**

이미지 임베딩을 가져옵니다.

## 에러 처리

Go 함수는 에러가 발생하면 JSON 응답의 `error` 필드에 에러 메시지를 담아 반환합니다.

Python 클라이언트는 에러가 있을 경우 `RuntimeError`를 발생시킵니다.

```python
try:
    embedding = client.get_text_embedding("test")
except RuntimeError as e:
    print(f"Error occurred: {e}")
```

## 성능 고려사항

- Go shared library는 약 32MB입니다 (Vertex AI SDK 포함)
- 각 API 호출은 새로운 클라이언트 연결을 생성합니다
- 대량 처리 시 적절한 배치 처리를 고려하세요

## 의존성

### Go 의존성
```
cloud.google.com/go/aiplatform/apiv1
google.golang.org/api/option
google.golang.org/protobuf/types/known/structpb
```

### Python 의존성
- Python 3.7+
- ctypes (표준 라이브러리)

## 라이선스

MIT License

## 참고자료

- [Vertex AI Multimodal Embedding API](https://cloud.google.com/vertex-ai/docs/generative-ai/embeddings/get-multimodal-embeddings)
- [Go CGO Documentation](https://golang.org/cmd/cgo/)
