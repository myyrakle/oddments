-- 6자리 랜덤코드 목록 생성
with 
BASE58_LIST as (
	select 
		ARRAY[
			'1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i',
			'j', 'k', 'm', 'n', 'o', 'p', 'q', 'r', 's',
			't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 
			'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 
			'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V',
			'W', 'X', 'Y', 'Z'
		] as list
)
select 
	_1.c||_2.c||_3.c||_4.c||_5.c||_6.c as code
from 
(
	select 
		(select list from BASE58_LIST)[t.num] as c
	from 
	  	(select generate_series(1, 58) as num) t
) _1
cross join
(
	select 
		(select list from BASE58_LIST)[t.num] as c
	from 
		(select generate_series(1, 58) as num) t
) _2
cross join
(
	select 
		(select list from BASE58_LIST)[t.num] as c
	from 
		(select generate_series(1, 58) as num) t
) _3
cross join
(
	select 
		(select list from BASE58_LIST)[t.num] as c
	from 
		(select generate_series(1, 58) as num) t
) _4
cross join
(
	select 
		(select list from BASE58_LIST)[t.num] as c
	from 
		(select generate_series(1, 58) as num) t
) _5
cross join
(
	select 
		(select list from BASE58_LIST)[t.num] as c
	from 
  		(select generate_series(1, 58) as num) t
) _6
