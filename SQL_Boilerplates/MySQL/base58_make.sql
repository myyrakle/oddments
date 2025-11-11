WITH RECURSIVE
BASE58_LIST AS
(
  SELECT
   Json_array
   ('1', '2', '3', '4', '5', '6', '7', '8', '9',
    'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i',
    'j', 'k', 'm', 'n', 'o', 'p', 'q', 'r', 's',
    't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B',
    'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L',
    'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V',
    'W', 'X', 'Y', 'Z'
  ) AS list
)
, NUMS AS
(
select 1 as num
 union
 select
   num + 1
 from NUMS where num<10
)
 select
    concat(
        json_unquote(json_extract(
            (select list from BASE58_LIST limit 1)
            , concat('$[', floor(rand()*58), ']')
        ))
        ,json_unquote(json_extract(
            (select list from BASE58_LIST limit 1)
            , concat('$[', floor(rand()*58), ']')
        ))
        ,json_unquote(json_extract(
            (select list from BASE58_LIST limit 1)
            , concat('$[', floor(rand()*58), ']')
        ))
        ,json_unquote(json_extract(
            (select list from BASE58_LIST limit 1)
            , concat('$[', floor(rand()*58), ']')
        )),json_unquote(json_extract(
            (select list from BASE58_LIST limit 1)
            , concat('$[', floor(rand()*58), ']')
        ))
        ,json_unquote(json_extract(
            (select list from BASE58_LIST limit 1)
            , concat('$[', floor(rand()*58), ']')
        ))
    )
   as code
 from NUMS
