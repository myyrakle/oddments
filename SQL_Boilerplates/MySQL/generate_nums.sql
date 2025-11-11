with recursive nums as
(
 select 1 as num
 union all
 select num + 1 from nums where num<10
)
select * from nums;
