# 通用键

* title 消息标题
* now 可选，用于测试时固定时间范围为now-interval
* interval 监控的周期
* index 索引名或通配符
* paradigm 统计方式，可填`count`、`percentage`、`spike`之一
* condition 因paradigm而异，见“各paradigm专用键”
* condition.match 值格式如`{"gt": 0, "lt": 1}`，表示当统计数字在0到1范围内时，发送消息；gt和lt可以不相交
* detail 因paradigm而异，见“各paradigm专用键”
* alarms 消息发送方式，键可选`stdout`、`email`、`ding`之中的一个或多个
* alarms.stdout 消息发送到标准输出，值为空对象
* alarms.email 消息以邮件方式发送，值格式为`{"to": []}`，数组可含多个收件人
* alarms.ding 消息以钉钉方式发送，值格式为`{"chats": [], "users": []}`，数组分别可含多个钉钉群ID和用户ID

# 各paradigm专用键

