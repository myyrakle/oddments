#![feature(portable_simd)]
use std::simd::{u8x64, Simd};

use wasm_bindgen::prelude::wasm_bindgen;

#[wasm_bindgen]
pub fn add(lhs: Vec<u8>, rhs: Vec<u8>) -> Vec<u8> {
    let mut result = Vec::with_capacity(lhs.len());

    for i in 0..lhs.len() {
        result.push(lhs[i] + rhs[i]);
    }

    return result;
}

#[cfg(target_arch = "wasm32")]
#[target_feature(enable = "simd128")]
#[wasm_bindgen]
pub fn add_with_simd(lhs: Vec<u8>, rhs: Vec<u8>) -> Vec<u8> {
    let mut result = vec![0; lhs.len()];

    let mut i = 0;
    while i < lhs.len() {
        // 길이가 64보다 작은 경우에는 그냥 계산
        if i + u8x64::splat(0).lanes() > lhs.len() {
            result[i] = lhs[i] + rhs[i];
            continue;
        }

        let lhs_simd: u8x64 = Simd::from_slice(&lhs[i..]);
        let rhs_simd: u8x64 = Simd::from_slice(&rhs[i..]);
        let result_simd = lhs_simd + rhs_simd;
        result_simd.copy_to_slice(&mut result[i..]);
        i += u8x64::splat(0).lanes();
    }

    return result;
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_add() {
        let lhs = vec![1, 2, 3, 4];
        let rhs = vec![5, 6, 7, 8];
        let result = add(lhs, rhs);
        assert_eq!(result, vec![6, 8, 10, 12]);
    }

    #[test]
    fn test_add_with_simd() {
        let lhs = vec![1, 2, 3, 4, 1, 2, 3, 4, 9];
        let rhs = vec![5, 6, 7, 8, 5, 6, 7, 8, 8];
        let result = add_with_simd(lhs, rhs);
        assert_eq!(result, vec![6, 8, 10, 12, 6, 8, 10, 12, 17]);
    }
}
