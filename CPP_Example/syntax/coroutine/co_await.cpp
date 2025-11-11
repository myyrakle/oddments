#include <iostream>
#include <coroutine>
#include <thread>
#include <chrono>

struct coroutine {
    struct promise_type {
        coroutine get_return_object() { return {}; }
        std::suspend_never initial_suspend() { return {}; }
        std::suspend_never final_suspend() noexcept { return {}; }
        void return_void() {}
        void unhandled_exception() {}
    };
};

template <class F>
struct awaitable {
    awaitable(F f) {
        this->task = std::thread([&] {
            f();
            });
    }

    std::thread task;

    bool await_ready() {
        if (this->task.joinable()) {
            this->task.join();
            return true;
        }
        else {
            return false;
        }
    }

    void await_suspend(std::coroutine_handle<> h) {
    }

    void await_resume() {
    }
};

template <class F>
auto async_thread(F f) {
    return awaitable<F>(f);
}

coroutine handle() {
    co_await async_thread([] {
        for (int i = 0; i < 5; i++) {
            puts("#####");
            std::this_thread::sleep_for(std::chrono::seconds(1));
        }
        
    });
     co_await async_thread([] {
        for (int i = 0; i < 5; i++) {
            puts("$$$$$");
            std::this_thread::sleep_for(std::chrono::seconds(1));
        }
    });
}


int main()
{
    handle();
}
