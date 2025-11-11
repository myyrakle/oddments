// 백준 node.js에서 입력 읽어오기
function get_input_from_baekjun() {
    const fs = require('fs');
    const input = fs.readFileSync('/dev/stdin').toString();
    return input;
}
