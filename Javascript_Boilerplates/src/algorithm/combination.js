function get_combination(n, m) {
    const table = [];
    
    function combination_recursion(n, m) {
	    // nCn = 1, nC0 = 1
	    if (n == m || m == BigInt(0)) {
            return BigInt(1);
        }

        if(table[n] === undefined) {
            table[n] = [];
        }

	    if (table[n][m] !== undefined && table[n][m] !== null) { 
            return table[n][m];
        }

	    // nCm = n-1Cm-1 + n-1Cm
	    table[n][m] = combination_recursion(n - BigInt(1), m - BigInt(1)) + combination_recursion(n - BigInt(1), m);
            
        return table[n][m];
    }
    
    return combination_recursion(n, m);
}
