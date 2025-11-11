defmodule MyOp do
    def l ** r, do: :math.pow(l, r)
end

import MyOp

IO.puts(10 ** 3)
