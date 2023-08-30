## 接口设计说明
### 路径
  - 路径统一格式： `/{projectName}/api/{version}/{模块名}/{接口功能名}`
    - 例如 `/myoption/api/v1/user/queryUserInfo`
    - 对个别路径可以对接口功能更进一步分类
  - 路径统一采用驼峰大小写样式，例如 `feedBack`
### 请求方法
  - 查询类请求采用GET，数据修改用POST
### 参数
  - 参数用驼峰样式
  - GET请求的参数值如果为中文或含有特殊字符，需要进行UrlEncode
  - POST请求统一用采用body方式传入为json字符串
### 响应
  - http status code统一为200。使用业务码code作为请求成功或失败的识别符。如下结构。
  - 响应结果统一结构 {"code":0,"message":"success","data":Object}
    - code为0代表请求成功，其他code代表失败，或者进行特殊约定。
    - message代表提示信息。data为实体数据结构，是一个对象。


## 默认值格式
### 时间
- 时间统一采用int格式，为毫秒时间戳