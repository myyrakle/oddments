from dotenv import load_dotenv
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import StrOutputParser
from langchain_google_genai import ChatGoogleGenerativeAI
import os

# 환경변수 로드
load_dotenv()


def basic_prompt_example():
    print("=== Gemini LLM 체인 예제 ===")

    # API 키 확인
    if (
        not os.getenv("GOOGLE_API_KEY")
        or os.getenv("GOOGLE_API_KEY") == "your_gemini_api_key_here"
    ):
        print("⚠️  GOOGLE_API_KEY가 설정되지 않았습니다.")
        print("실제 Gemini를 사용하려면 .env 파일에 GOOGLE_API_KEY를 설정해주세요.")
        print("Google AI Studio에서 키 발급: https://aistudio.google.com/app/apikey")
        print()
        return

    try:
        # Gemini LLM 초기화
        llm = ChatGoogleGenerativeAI(model="gemini-2.0-flash", temperature=0.7)

        # 프롬프트 템플릿
        prompt = ChatPromptTemplate.from_messages(
            [("system", "당신은 친근한 AI 어시스턴트입니다."), ("human", "{question}")]
        )

        # 출력 파서
        output_parser = StrOutputParser()

        # 체인 구성
        chain = prompt | llm | output_parser

        # 체인 실행
        response = chain.invoke(
            {"question": "LangChain 1.0과 Gemini의 주요 특징을 간단히 알려주세요."}
        )

        print("✅ Gemini 응답:")
        print(response)
        print()

    except Exception as e:
        print(f"❌ Gemini 호출 중 오류 발생: {e}")
        print()


def main():
    """메인 함수"""
    basic_prompt_example()


if __name__ == "__main__":
    main()
