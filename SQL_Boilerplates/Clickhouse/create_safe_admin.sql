-- Drop 불가 admin 생성

-- 1. 역할 생성
CREATE ROLE safe_admin_role;

-- 2. 역할에 권한 부여
-- 2. SELECT/INSERT/UPDATE 권한
GRANT SELECT ON *.* TO safe_admin_role;
GRANT INSERT ON *.* TO safe_admin_role;
GRANT ALTER UPDATE ON *.* TO safe_admin_role;
GRANT ALTER DELETE ON *.* TO safe_admin_role;
GRANT SHOW ON *.* TO safe_admin_role;
GRANT SHOW DATABASES ON *.* TO safe_admin_role;

-- 3. DDL 권한 (DROP 제외)
GRANT CREATE TABLE ON *.* TO safe_admin_role;
GRANT CREATE VIEW ON *.* TO safe_admin_role;
GRANT CREATE DATABASE ON *.* TO safe_admin_role;
GRANT ALTER TABLE ON *.* TO safe_admin_role;
GRANT ALTER VIEW ON *.* TO safe_admin_role;;

-- 4. 사용자 생성 및 역할 할당
CREATE USER safe_admin IDENTIFIED BY 'gqFqyperQqmB8O3D';
GRANT safe_admin_role TO safe_admin;
