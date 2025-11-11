#include <iostream>
#include <coroutine>

// 밑에서 다시 정의할 예정
struct promise;

template <class T>
struct coroutine;

// 코루틴 객체
template <class T>
struct coroutine: std::coroutine_handle<promise>{
	struct promise {
		coroutine get_return_object()
		{
			return coroutine(std::coroutine_handle<promise>::from_promise(*this));
		}

		std::suspend_always initial_suspend() noexcept
		{
			return {};
		}

		std::suspend_always final_suspend() noexcept
		{
			return {};
		}

		void return_void() {
		}

		void unhandled_exception() { std::puts("예외 발생"); }

		T value;

		std::suspend_always yield_value(T from) {
			this->value = from; // caching the result in promise
			return {};
		}
	};

	using promise_type = struct promise;
	std::coroutine_handle<promise> handle;


	coroutine(std::coroutine_handle<promise> handle) {
		this->handle = handle;
	}

	T operator()() {
		handle(); // resume
		auto result = handle.promise();
		return result.value;
	}
};

coroutine<int> new_counter() {
	int count = 0;

	while (true) {
		count++;
		co_yield count;
	}
}


int main()
{
	coroutine counter = new_counter();

	std::cout << counter() << std::endl;
	std::cout << counter() << std::endl;
	std::cout << counter() << std::endl;
	std::cout << counter() << std::endl;
}
