SELECT
  EXTRACT(YEAR FROM age(current_timestamp AT time ZONE 'Asia/Seoul', TO_TIMESTAMP(BIRTH_DT))) AS AGE
