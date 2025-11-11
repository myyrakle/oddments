fn compute_gcd(a: i32, b: i32) -> i32 {
    if b == 0 {
        return a;
    }
    
    compute_gcd(b, a%b)
}

fn main() {
    println!("최대공약수 {:?}", compute_gcd(1230, 1040));
    
    return ();
}
