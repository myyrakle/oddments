select 
	j.name
from (
	select CONCAT('["', REPLACE('셔츠,블라우스,데님셔츠,옥스포드셔츠,옥스포드,스트라이프셔츠,체크셔츠,남방', ',', '","'), '"]') as name
) t
join json_table(
  t.name,
  '$[*]' columns (name varchar(50) path '$')
) j
