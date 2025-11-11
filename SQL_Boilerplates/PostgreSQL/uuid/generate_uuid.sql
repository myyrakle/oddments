SELECT uuid_in(md5(random()::text || clock_timestamp()::text)::cstring);
