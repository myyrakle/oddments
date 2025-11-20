"""
LangChain 기본 기능만으로 구현한 의도 분류 기반 멀티 에이전트 예제

이 예제는 다음과 같은 구조로 되어 있습니다:
1. Intro 에이전트: 사용자 입력의 의도를 분류 (HELP vs SMALLTALK)
2. Help 에이전트: 도움/질문/문제해결이 필요한 경우 전문적인 답변 제공
3. Smalltalk 에이전트: 일상대화/인사/잡담에 친근하게 응답

각 에이전트는 고유한 역할과 프롬프트를 가지고, 
Intro 에이전트의 분류 결과에 따라 적절한 전문 에이전트가 선택됩니다.
"""

import os
from dotenv import load_dotenv
from langchain_google_genai import ChatGoogleGenerativeAI
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import StrOutputParser

# 환경변수 로드
load_dotenv()

class Agent:
    """간단한 에이전트 클래스"""
    
    def __init__(self, name, role, prompt_template, llm):
        self.name = name
        self.role = role
        self.prompt = ChatPromptTemplate.from_template(prompt_template)
        self.llm = llm
        self.output_parser = StrOutputParser()
        self.chain = self.prompt | self.llm | self.output_parser
    
    def run(self, input_data):
        """에이전트 실행"""
        print(f"\n🤖 {self.name} 에이전트가 작업 중...")
        print(f"역할: {self.role}")
        print("-" * 50)
        
        try:
            result = self.chain.invoke(input_data)
            print(f"✅ {self.name} 완료!")
            return result
        except Exception as e:
            print(f"❌ {self.name} 오류: {str(e)}")
            return None

class MultiAgentSystem:
    """의도 분류 기반 멀티 에이전트 시스템"""
    
    def __init__(self, llm):
        self.llm = llm
        self.agents = {}
        self.setup_agents()
    
    def setup_agents(self):
        """에이전트들을 설정"""
        
        # 1. Intro 에이전트 - 의도 분류
        intro_prompt = """
당신은 사용자의 의도를 분석하는 전문가입니다.
사용자의 입력을 분석하고 다음 중 하나로 분류해주세요:

분류 옵션:
- HELP: 도움이나 문제 해결이 필요한 경우 (질문, 가이드 요청, 튜토리얼, 문제 해결 등)
- SMALLTALK: 일상 대화나 잡담 (인사, 날씨, 감정 표현, 개인적 이야기 등)

사용자 입력: {user_input}

응답 형식:
분류: [HELP 또는 SMALLTALK]
이유: [분류한 이유를 한 줄로 설명]

분류만 명확하게 해주세요.
"""
        
        # 2. Help 에이전트 - 도움 및 문제 해결
        help_prompt = """
당신은 친절하고 지식이 풍부한 도우미입니다.
사용자의 질문이나 문제를 해결하는 데 도움을 주세요.

사용자 질문: {user_input}

다음과 같이 도움을 제공해주세요:
- 명확하고 구체적인 답변 제공
- 필요시 단계별 가이드 제공
- 추가 참고사항이나 팁 포함
- 이해하기 쉬운 설명 사용

전문적이면서도 친근하게 답변해주세요.
"""
        
        # 3. Smalltalk 에이전트 - 일상 대화
        smalltalk_prompt = """
당신은 친근하고 공감능력이 뛰어난 대화 상대입니다.
사용자와 자연스러운 일상 대화를 나누세요.

사용자 말: {user_input}

다음과 같이 대화해주세요:
- 따뜻하고 친근한 톤 사용
- 적절한 감정 표현과 공감
- 자연스러운 대화 흐름 유지
- 필요시 관련된 질문이나 주제 확장

편안하고 즐거운 대화를 만들어주세요.
"""
        
        # 에이전트 생성
        self.agents['intro'] = Agent(
            "Intro", 
            "사용자 의도 분류",
            intro_prompt, 
            self.llm
        )
        
        self.agents['help'] = Agent(
            "Help",
            "도움 및 문제 해결", 
            help_prompt,
            self.llm
        )
        
        self.agents['smalltalk'] = Agent(
            "Smalltalk",
            "일상 대화 및 잡담",
            smalltalk_prompt,
            self.llm
        )
    
    def classify_intent(self, user_input):
        """사용자 의도 분류"""
        print("� 사용자 의도를 분석 중...")
        
        result = self.agents['intro'].run({"user_input": user_input})
        if not result:
            return "HELP"  # 기본값
        
        # 분류 결과 파싱
        if "HELP" in result.upper():
            return "HELP"
        elif "SMALLTALK" in result.upper():
            return "SMALLTALK"
        else:
            return "HELP"  # 기본값
    
    def run_conversation(self, user_input):
        """대화 시스템 실행"""
        print("🚀 멀티 에이전트 대화 시스템 시작!")
        print(f"사용자: {user_input}")
        print("=" * 60)
        
        # 1단계: 의도 분류
        intent = self.classify_intent(user_input)
        print(f"\n🎯 분류 결과: {intent}")
        print("=" * 60)
        
        # 2단계: 해당 에이전트 실행
        if intent == "HELP":
            response = self.agents['help'].run({"user_input": user_input})
        else:  # SMALLTALK
            response = self.agents['smalltalk'].run({"user_input": user_input})
        
        if response:
            print(f"\n� 최종 응답:")
            print(response)
        
        print("\n" + "=" * 60)
        return response

def main():
    print("🔗 LangChain 의도 분류 기반 멀티 에이전트 시스템")
    print("=" * 50)
    
    # API 키 확인
    api_key = os.getenv("GOOGLE_API_KEY")
    if not api_key:
        print("❌ GOOGLE_API_KEY가 설정되지 않았습니다.")
        print("💡 .env 파일에 API 키를 설정해주세요.")
        return
    
    try:
        # Gemini 모델 초기화
        llm = ChatGoogleGenerativeAI(
            model="gemini-2.0-flash",
            temperature=0.7,
            google_api_key=api_key
        )
        
        # 멀티 에이전트 시스템 생성
        multi_agent = MultiAgentSystem(llm)
        
        print("시스템이 준비되었습니다!")
        print("\n💡 사용법:")
        print("- 질문이나 도움이 필요하면 → Help 에이전트가 응답")
        print("- 인사나 일상 대화를 하면 → Smalltalk 에이전트가 응답")
        print("- 'quit' 또는 '종료'를 입력하면 종료")
        print("\n" + "=" * 50)
        
        # 대화형 루프
        while True:
            try:
                user_input = input("\n👤 당신: ").strip()
                
                if user_input.lower() in ['quit', 'exit', '종료', '끝']:
                    print("👋 대화를 종료합니다. 감사합니다!")
                    break
                
                if not user_input:
                    continue
                
                # 멀티 에이전트 대화 실행
                multi_agent.run_conversation(user_input)
                
            except KeyboardInterrupt:
                print("\n👋 대화를 종료합니다. 감사합니다!")
                break
            except Exception as e:
                print(f"❌ 오류 발생: {str(e)}")
        
    except Exception as e:
        print(f"❌ 시스템 초기화 중 오류 발생: {str(e)}")

if __name__ == "__main__":
    main()
