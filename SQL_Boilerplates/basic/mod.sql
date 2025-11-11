select 
	mod(10, 3) as std, /* 표준함수 */
	10%3 as psql /* Postgres 비표준 */
