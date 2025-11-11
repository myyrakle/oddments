function lcm(n1, n2) {
    function _gcd(a, b){
        if (!b) {
            return a;
        }

        return _gcd(b, a % b);
    }

    const gcdValue = _gcd(n1, n2);

    return (n1 * n2) / gcdValue;
}
