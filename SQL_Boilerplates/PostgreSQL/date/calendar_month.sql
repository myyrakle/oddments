/*달력형태의 한 달 범위(한 주에 끼인 이전달과 다음달 날까지)*/
SELECT 
	DATE_TRUNC('week', DATE_TRUNC('month', CURRENT_TIMESTAMP) + INTERVAL '1 day') - INTERVAL '1 day' AS CALENDAR_BEGIN,
	DATE_TRUNC('week', (DATE_TRUNC('month', CURRENT_TIMESTAMP + INTERVAL '1 month') - INTERVAL '1 day') + INTERVAL '1 day' + INTERVAL '1 week') - INTERVAL '2 day' AS CALENDAR_END
