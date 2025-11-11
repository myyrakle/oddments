fn compute_gcd(a: i32, b: i32) -> i32 {
    // a의 약수 목록 구하기
    let a_list: Vec<i32> = (1..=a).into_iter()
        .filter(|e| a%e == 0)
        .collect();
    println!("a의 약수 {:?}", a_list);
    
    // b의 약수 목록 구하기
    let b_list: Vec<i32> = (1..=b).into_iter()
        .filter(|e| b%e == 0)
        .collect();
    println!("b의 약수 {:?}", b_list);
    
    // 교집합으로 공약수만 추출하기
    let cd_list: Vec<i32> = a_list.into_iter()
        .filter(|e| b_list.contains(e))
        .collect();
    println!("공약수 {:?}", cd_list);
    
    // 가장 큰 수를 반환
    let gcd = *cd_list.last().unwrap();
    
    return gcd;
}

fn main() {
    println!("최대공약수 {:?}", compute_gcd(123000, 104000));
    
    return ();
}
