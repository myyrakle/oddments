select
	1 = 1 as eq, /* 같으면 true */
	1 != 2 as neq, /* 다르면 true */
	1 <> 2 as neq2, /* 다르면 true */
	1 < 5 as lt, /* 왼쪽이 작으면 true */
	20 > 5 as gt, /* 왼쪽이 크면 true */
	1 <= 5 as lte, /* 왼쪽이 작거나 같으면 true */
	20 >= 5 as gte /* 왼쪽이 크거나 같으면 true */
