macro_rules! is_trait{
    ($name:ty, $trait_name:path)=>{{
        trait __InnerMarkerTrait{
            fn __is_trait_inner_method()->bool{
                false
            }
        }
        struct __TraitTest<T>(T);
        impl<T:$trait_name> __TraitTest<T> {
            fn __is_trait_inner_method()->bool{
                true
            }
        }
        impl<T> __InnerMarkerTrait for __TraitTest<T>{}
        __TraitTest::<$name>::__is_trait_inner_method()
    }}
}
fn main() {
    println!("{}", is_trait!(i32, Copy)); //true
    println!("{}", is_trait!(String, Copy)); //false
}
