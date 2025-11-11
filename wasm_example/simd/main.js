import wasmUrl from "./rust/pkg/rust_bg.wasm?url";

export let wasm_instance = null;
export let wasm_module = null;

const main = async () => {
  const responsePromise = fetch(wasmUrl);
  const { module, instance } = await WebAssembly.instantiateStreaming(
    responsePromise
  );

  wasm_instance = instance;
  wasm_module = module;
};

main();
