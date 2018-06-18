--登录 psql -Upgmxo -d pgmxodb
--导入 psql -U pgmxo -d pgmxodb -f pg.sql 
--导出 pg_dump -U postgres -f pgdump.sql pgmxodb
--更改用户 ALTER TABLE public.kloud_entry OWNER TO pgmxo

-- postgresql10-contrib
create extension pgcrypto;

CREATE TABLE kloud_user (
  id VARCHAR primary key not null,
  name VARCHAR NOT NULL,
  email VARCHAR unique NOT NULL,
  role int2 not null default 2, -- '012 admin/normal/guest'
  password VARCHAR NOT NULL,
  avatar_url VARCHAR, --头像
  create_time timestamp not null default current_timestamp,
  disabled bool default FALSE not null, -- '帐号状态：禁用'
  max_size int8 not null default 4294967296, -- 1024*1024*1024*4=4G
  cur_size int8 not null default 0
);

create index kloud_user_idx_name on kloud_user(name);
create index kloud_user_idx_email on kloud_user(email);

CREATE TABLE kloud_session (
  id VARCHAR primary key not null,
  user_id VARCHAR references kloud_user (id) ON DELETE CASCADE not null,
  user_agent VARCHAR not null, -- 创建该session的UA
  ip VARCHAR not null, -- 创建该session的IP
  create_time  timestamp not null default current_timestamp,
  expire_time  timestamp not null -- 有效期至
);

create index kloud_session_idx_user_id on kloud_session(user_id);

CREATE TABLE kloud_site (
  name VARCHAR default 'kloud' not null,
  domain VARCHAR not null, -- null表示没有
  favicon_url VARCHAR,
  logo_url VARCHAR
);

CREATE TABLE kloud_file (
  id VARCHAR primary key not null,
  sha1 VARCHAR unique not null,
  size int8 not null,
  create_time  timestamp not null default current_timestamp,
  name VARCHAR not null,  
  path VARCHAR unique not null, --储存路径
  rc int4 default 0 not null, --引用计数
  disabled bool default FALSE not null --管理员标记文件为禁用，仅管理员可以下载
);

create index kloud_file_idx_sha1 on kloud_file(sha1);

CREATE TABLE kloud_entry (
  id VARCHAR primary key not null,
  user_id VARCHAR references kloud_user (id) not null, --所属用户
  file_id VARCHAR references kloud_file (id),--文件id, null表示目录
  parent_id VARCHAR references kloud_entry (id), --上级目录id, /的为 null
  name VARCHAR not NULL, 
  create_time timestamp not null default current_timestamp,
  modify_time timestamp not null default current_timestamp,
  delete_time timestamp , --如果不为null则在回收站，从回收站删除就彻底没了
  UNIQUE (user_id, parent_id, name) --防止同一级目录建立了一样的名字的
);

create index kloud_entry_idx_user_id on kloud_entry(user_id);
create index kloud_entry_idx_parent_id on kloud_entry(parent_id);
create index kloud_entry_idx_name on kloud_entry(name);

-- 分享
CREATE TABLE kloud_share (
  id VARCHAR primary key not null,
  name VARCHAR not null, 
  user_id VARCHAR references kloud_user (id) ON DELETE CASCADE not null, --所属用户
  entry_id VARCHAR references kloud_entry (id) ON DELETE SET NULL, --条目id，从这个token这里只能下载其中的内容
  create_time timestamp not null default current_timestamp,
  expire_time timestamp, -- null表示不限
  password VARCHAR -- null表示没有密码。并禁掉其它的 oringin Referer（一般来说是域名，没有几个站果跑ip的，不管referer是ip的，如此简单一些）
);
create index kloud_share_idx_user_id on kloud_share(user_id);

-- 公开，作为静态资源，主要是为了和分享分开
CREATE TABLE kloud_public (
  id VARCHAR primary key not null,
  name VARCHAR not null,  
  user_id VARCHAR references kloud_user (id) ON DELETE CASCADE not null, --所属用户
  entry_id VARCHAR references kloud_entry (id) ON DELETE SET NULL,--条目id，从这个token这里只能下载其中的内容
  create_time timestamp not null default current_timestamp,
  expire_time timestamp,
  referer VARCHAR -- 允许的域名正则, 默认空表示不限
);
create index kloud_public_idx_user_id on kloud_public(user_id);

-- 以后再研究。。
-- 用户给第三方网站（或未认证用户）的
-- 删除的话，第三方网站服务器就可以直接删除，还是要个delete-token？
-- CREATE TABLE kloud_upload (
--   id uuid primary key default gen_random_uuid(), --新的条目id（完成后插入entry表）
--   user_id uuid not null, --所属用户
--   parent_uuid not null,--父级条目id 
--   name VARCHAR not null, -- 上传后的名字
--   sha1 VARCHAR, -- null表示目录
--   length int8,
--   create_time timestamp not null default current_timestamp,
--   expire_time timestamp not null,
--   allow_ip VARCHAR -- null表示不限
-- );


-- 注意 -n 在前面，否则就多个换行符
--  echo -n  pd |sha256sum

-- admin.kloud.000
insert into kloud_user (id, name, email, role, password) values(
  'KU4lv0o47ompYWorvYSmW6',
  'Wspsxing', 'biluohc@qq.com', 0, encode(digest(encode(digest('admin.kloud.000'||'kloud', 'sha256'), 'hex')||'kloud', 'sha256'), 'hex') );

-- demo.kloud.000
insert into kloud_user (id, name, email, role, password) values('Vm2ZjtC517ObGJd3y8JJUX','demo', 'dustdawn@qq.com', 1, encode(digest(encode(digest('demo.kloud.000'||'kloud', 'sha256'), 'hex')||'kloud', 'sha256'), 'hex'));

-- drop table kloud_public;
-- drop table kloud_share;
-- drop table kloud_entry;
-- drop table kloud_file;
-- drop table kloud_session;
-- drop table kloud_user;
-- drop table kloud_site;



-- https://github.com/go-pg/pg/wiki/Writing-Queries#executing-custom-queries
-- 还是可以建表的，以后来做，就是差个表在不在的 api
-- qs := []string{
--     "CREATE TEMP TABLE users (id int, name text)",
--     "CREATE TEMP TABLE profiles (id int, lang text, active bool, user_id int)",
--     "INSERT INTO users VALUES (1, 'user 1')",
--     "INSERT INTO profiles VALUES (1, 'en', TRUE, 1), (2, 'ru', TRUE, 1), (3, 'md', FALSE, 1)",
-- }
-- for _, q := range qs {
--     _, err := db.Exec(q)
--     if err != nil {
--         panic(err)
--     }
-- }