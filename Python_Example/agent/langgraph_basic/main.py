"""
LangGraph 1.0을 사용한 멀티 에이전트 시스템 예제

이 예제는 다음과 같은 구조로 되어 있습니다:
1. Intro 에이전트: 사용자 입력의 의도를 분류 (HELP vs SMALLTALK)
2. Help 에이전트: 도움/질문/문제해결이 필요한 경우 전문적인 답변 제공
3. Smalltalk 에이전트: 일상대화/인사/잡담에 친근하게 응답

LangGraph 1.0을 사용하여 에이전트 간의 흐름을 그래프로 정의하고,
조건부 라우팅을 통해 의도에 따라 적절한 에이전트가 선택됩니다.
"""

import os
from typing import TypedDict, Literal, Annotated
from dotenv import load_dotenv
from langchain_google_genai import ChatGoogleGenerativeAI
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import StrOutputParser
from langgraph.graph import StateGraph, START, END
from langgraph.types import Command

# 환경변수 로드
load_dotenv()


class State(TypedDict):
    """에이전트 시스템의 상태"""
    user_input: str
    intent: Literal["HELP", "SMALLTALK", "UNKNOWN"]
    response: str
    classification_result: str


class MultiAgentSystem:
    """LangGraph 기반 멀티 에이전트 시스템"""

    def __init__(self, llm):
        self.llm = llm
        self.graph = self._build_graph()
        self.app = self.graph.compile()

    def _build_graph(self) -> StateGraph:
        """LangGraph 그래프 구축"""
        graph = StateGraph(State)

        # 에이전트 노드 정의
        graph.add_node("intro_node", self._intro_agent)
        graph.add_node("help_node", self._help_agent)
        graph.add_node("smalltalk_node", self._smalltalk_agent)

        # 엣지 정의 (그래프 흐름)
        graph.add_edge(START, "intro_node")
        graph.add_conditional_edges(
            "intro_node",
            self._route_to_agent,
            {
                "help_node": "help_node",
                "smalltalk_node": "smalltalk_node"
            }
        )
        graph.add_edge("help_node", END)
        graph.add_edge("smalltalk_node", END)

        return graph

    def _intro_agent(self, state: State) -> State:
        """의도 분류 에이전트"""
        print("\n🎯 Intro 에이전트: 사용자 의도를 분석 중...")

        intro_prompt = ChatPromptTemplate.from_template("""
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
""")

        chain = intro_prompt | self.llm | StrOutputParser()
        result = chain.invoke({"user_input": state["user_input"]})

        # 분류 결과 파싱
        intent = "UNKNOWN"
        if "HELP" in result.upper():
            intent = "HELP"
        elif "SMALLTALK" in result.upper():
            intent = "SMALLTALK"

        print(f"✅ 분류 결과: {intent}")

        state["intent"] = intent
        state["classification_result"] = result
        return state

    def _help_agent(self, state: State) -> State:
        """도움 및 문제 해결 에이전트"""
        print("\n💡 Help 에이전트: 문제를 해결 중...")

        help_prompt = ChatPromptTemplate.from_template("""
당신은 친절하고 지식이 풍부한 도우미입니다.
사용자의 질문이나 문제를 해결하는 데 도움을 주세요.

사용자 질문: {user_input}

다음과 같이 도움을 제공해주세요:
- 명확하고 구체적인 답변 제공
- 필요시 단계별 가이드 제공
- 추가 참고사항이나 팁 포함
- 이해하기 쉬운 설명 사용

전문적이면서도 친근하게 답변해주세요.
""")

        chain = help_prompt | self.llm | StrOutputParser()
        response = chain.invoke({"user_input": state["user_input"]})

        print(f"✅ Help 에이전트 완료!")

        state["response"] = response
        return state

    def _smalltalk_agent(self, state: State) -> State:
        """일상 대화 에이전트"""
        print("\n😊 Smalltalk 에이전트: 친근하게 대화 중...")

        smalltalk_prompt = ChatPromptTemplate.from_template("""
당신은 친근하고 공감능력이 뛰어난 대화 상대입니다.
사용자와 자연스러운 일상 대화를 나누세요.

사용자 말: {user_input}

다음과 같이 대화해주세요:
- 따뜻하고 친근한 톤 사용
- 적절한 감정 표현과 공감
- 자연스러운 대화 흐름 유지
- 필요시 관련된 질문이나 주제 확장

편안하고 즐거운 대화를 만들어주세요.
""")

        chain = smalltalk_prompt | self.llm | StrOutputParser()
        response = chain.invoke({"user_input": state["user_input"]})

        print(f"✅ Smalltalk 에이전트 완료!")

        state["response"] = response
        return state

    def _route_to_agent(self, state: State) -> Literal["help_node", "smalltalk_node"]:
        """Intro 에이전트의 결과에 따라 다음 에이전트를 선택"""
        if state["intent"] == "HELP":
            return "help_node"
        elif state["intent"] == "SMALLTALK":
            return "smalltalk_node"
        else:
            return "help_node"  # 기본값

    def run(self, user_input: str) -> str:
        """멀티 에이전트 시스템 실행"""
        initial_state = {
            "user_input": user_input,
            "intent": "UNKNOWN",
            "response": "",
            "classification_result": ""
        }

        final_state = self.app.invoke(initial_state)
        return final_state["response"]


def main():
    """메인 함수"""
    print("🔗 LangGraph 1.0 멀티 에이전트 시스템")
    print("=" * 60)

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
        system = MultiAgentSystem(llm)

        print("✅ 시스템이 준비되었습니다!\n")
        print("💡 사용법:")
        print("- 질문이나 도움이 필요하면 → Help 에이전트가 응답")
        print("- 인사나 일상 대화를 하면 → Smalltalk 에이전트가 응답")
        print("- 'quit' 또는 '종료'를 입력하면 종료")
        print("\n" + "=" * 60)

        # 대화형 루프
        while True:
            try:
                user_input = input("\n👤 당신: ").strip()

                if user_input.lower() in ['quit', 'exit', '종료', '끝']:
                    print("👋 대화를 종료합니다. 감사합니다!")
                    break

                if not user_input:
                    continue

                # 멀티 에이전트 시스템 실행
                print("\n" + "=" * 60)
                response = system.run(user_input)

                print("\n" + "=" * 60)
                print(f"\n🤖 응답:")
                print(response)
                print("\n" + "=" * 60)

            except KeyboardInterrupt:
                print("\n👋 대화를 종료합니다. 감사합니다!")
                break
            except Exception as e:
                print(f"❌ 오류 발생: {str(e)}")

    except Exception as e:
        print(f"❌ 시스템 초기화 중 오류 발생: {str(e)}")


if __name__ == "__main__":
    main()
