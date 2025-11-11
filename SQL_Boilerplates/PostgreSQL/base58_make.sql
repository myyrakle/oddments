WITH
BASE58_LIST AS (
	SELECT 
		ARRAY[
			'1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i',
			'j', 'k', 'm', 'n', 'o', 'p', 'q', 'r', 's',
			't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 
			'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 
			'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V',
			'W', 'X', 'Y', 'Z'
		] AS list
)
SELECT 
	-- 6자리 코드 생성
	(SELECT list FROM BASE58_LIST)[(RANDOM()*57)::INTEGER + 1]
	||
	(SELECT list FROM BASE58_LIST)[(RANDOM()*57)::INTEGER + 1]
	||
	(SELECT list FROM BASE58_LIST)[(RANDOM()*57)::INTEGER + 1]
	||
	(SELECT list FROM BASE58_LIST)[(RANDOM()*57)::INTEGER + 1]
	||
	(SELECT list FROM BASE58_LIST)[(RANDOM()*57)::INTEGER + 1]
	||
	(SELECT list FROM BASE58_LIST)[(RANDOM()*57)::INTEGER + 1]
	AS code
FROM 
(SELECT GENERATE_SERIES(1, 10)) T -- 코드 10개 생성. 숫자만 늘리면 그만큼 늘어남
