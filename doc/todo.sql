--
CREATE TABLE IF NOT EXISTS public.user_register_info

(
    id               serial8      NOT NULL PRIMARY KEY,
    version          varchar(64)  NOT NULL,
    os_name          varchar(64)  NOT NULL,
    os_version       varchar(64)  NOT NULL,
    device_name      varchar(64)  NOT NULL,
    device_info      varchar(64)  NOT NULL,
    birthday         integer      NOT NULL    DEFAULT 0,
    secure_question1 varchar(256) NOT NULL,
    secure_question2 varchar(256) NOT NULL,
    secure_question3 varchar(256) NOT NULL,
    create_time      timestamp with time zone,
    update_time      timestamp with time zone DEFAULT current_timestamp
);

-- 您的出生地是哪里？
-- 您最喜欢的食物是什么？
-- 您的小学校名是什么？
-- 您的父亲的姓名是什么？
-- 您的母亲的生日是哪天？
-- 您最喜欢的电影是哪部？
-- 您最喜欢的颜色是什么？
-- 您的第一个宠物的名字是什么？
-- 您最喜欢的体育运动是什么？
-- 您的幸运数字是多少？
-- 您的第一个手机号是多少？
--------



