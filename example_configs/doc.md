# 通用键

## 查询

* now： 可选，用于测试时固定时间范围为now-interval
* time_field： 可选，指定监控周期依据文档的哪个时间字段，默认是`@timestamp`
* interval： 监控的周期
* index： 符合ES API规范的索引名，可含多个索引、通配符、日期计算表达式
* paradigm： 统计方式，可填`count`、`percentage`、`spike`之一
* condition： 因paradigm而异，见“各paradigm专用键及算法”
* condition.match： 值格式如`{"gt": 0, "lt": 1, "not": false, "aggs": false}`，此例表示当统计数字在0到1范围内时，发送消息；`gt`和`lt`可以相交或不相交；`not`可选，默认为`fales`，填`true`时，取`gt`和`lt`所指的范围之外。`aggs`可选，默认为`false`，表示以`condition.query`的统计数来判断是否符合范围，填`true`则会在`condition.aggs`的结果集里有符合范围的组时，才发送消息

## 通知

* title： 消息标题
* alarms： 消息发送方式，键可选`stdout`、`email`、`ding`、`web_hook`之中的一个或多个
* alarms.stdout： 消息发送到标准输出，值为空对象
* alarms.email： 消息以邮件方式发送，值格式为`{"to": []}`，数组可含多个收件人
* alarms.ding： 消息以钉钉方式发送，值格式为`{"chats": [], "users": [], "robots": []}`，数组分别可含多个钉钉群ID、用户ID、机器人access_token
* alarms.web_hook： 消息以JSON形式发送，例如`{"title": "...", "abstract": "...", details: [{"terms": ["value_1", ...], "count": 7, "calculated": 0.12}, ...]}`。配置格式为`{"url": "http://....", "method": "POST"}`，其中`method`可选，默认为`"POST"`

# 各paradigm专用键及算法

各种paradigm均可用`condition.query`筛选出参与计数的文档，用`condition.aggs`对选出的文档进行分组，语法同es查询DSL。每种计算方法对结果集的利用如下

## count

根据`condition.query`查出`interval`内doc的数量，符合范围则报警

若`condition.match.aggs: true`，则只在`condition.aggs`结果集里有任一组doc的数量符合范围时，才报警
  
## percentage

根据`condition.query`查出`interval`内doc的数量，并在此基础上以`condition.partial_query`进一步筛选出部分doc，当该部分占比符合范围则报警

若`condition.match.aggs: true`，则无需配置`condition.partial_query`，程序会在`condition.aggs`结果集里有任一组doc的占比符合范围时，才报警

## spike

根据`condition.query`分别查出最近`interval`及上一`interval`内doc的数量，当两个周期的比值（变化率）符合范围则报警

当上一周期的文档数为零时，以`condition.reference`的值代替，若无设置，则不报警

若`condition.match.aggs: true`，则只在任一分组的doc在两个周期之间的变化率符合范围时才报警。如果分组在本周期有出现而上一周期没出现，则其在上一周期的数量以`condition.reference`的值代替，若无设置，则不报警。如果分组在本周期没出现而上一周期有出现，则其在本周期的数量以`condition.reference`的值代替，若无设置，则本周期数量按“零”计算。因此分组数`size`最好设置一个较大的值，以尽量包含所有分组（当然，性能会有所下降）
