# 通用键

## 查询

* now： 可选，用于测试时固定时间范围为now-interval
* interval： 监控的周期，以文档的`@timestamp`字段来过滤
* index： 符合ES API规范的索引名，可含多个索引、通配符、日期计算表达式
* paradigm： 统计方式，可填`count`、`percentage`、`spike`之一
* condition： 因paradigm而异，见“各paradigm专用键”
* condition.match： 值格式如`{"gt": 0, "lt": 1}`，表示当统计数字在0到1范围内时，发送消息；gt和lt可以相交或不相交
* detail： 用于将查询所得汇总，格式同es聚合DSL

## 通知

* title： 消息标题
* alarms： 消息发送方式，键可选`stdout`、`email`、`ding`之中的一个或多个
* alarms.stdout： 消息发送到标准输出，值为空对象
* alarms.email： 消息以邮件方式发送，值格式为`{"to": []}`，数组可含多个收件人
* alarms.ding： 消息以钉钉方式发送，值格式为`{"chats": [], "users": [], "robots": []}`，数组分别可含多个钉钉群ID、用户ID、机器人access_token

# 各paradigm专用键

以下用于筛选文档的键所对应的值，格式同es查询DSL

## count

* condition.scope：筛选出参与计数的文档
  
## percentage

* condition.part：筛选出作为分子的文档数
* condition.whole：筛选出作为分母的文档数

## spike

* condition.scope：筛选出参与计数的文档
* condition.reference：当上一周期的文档数为零时，以此值代替，若无设置，则不报警
