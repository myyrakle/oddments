#include <random>

template <class Engine = std::mt19937, class Integer = int>
Integer random(Integer begin, Integer end)
{
  static Engine engine = []{
    Engine engine;
    std::random_device device;
    engine.seed(device());
    return engine;
  }();
  
  static std::uniform_int_distribution<Integer> dist{};
  using param_t = typename decltype(dist)::param_type;
  
  return dist(engine, param_t(begin, end-1));
}
