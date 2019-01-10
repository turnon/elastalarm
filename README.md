# ElastAlarm

## 运行

```
elastalarm -host=http://192.168.0.60:9200 -configs=/path/to/config_dir
```

## 环境变量

* `ESALARM_MAIL_SERVER=smtp.qq.com:465`
* `ESALARM_MAIL_SKIP_VERIFY="true"`
* `ESALARM_MAIL_FROM="someone@qq.com"`
* `ESALARM_MAIL_PASSWD="klkjjfkj5645"`
* `ESALARM_MAIL_TO="another_one@qq.com"`
* `ESALARM_DING_CORPID="ding134nbcvbmn"`
* `ESALARM_DING_SECRET="S-mbvfhbfjhgjh657nvnjvhf74"`
* `ESALARM_DING_CHATID="chatbjbnjghgukh"`

## 统计方法和配置方式

详见`./example_configs`

* `count`： `interval`内根据`scope`条件所查出doc数量是否符合范围
* `percentage`： `interval`内`part`/`whole` 所得的百分比是否符合范围
* `spike`： 最近`interval`的`scope`与上一`interval`的`scope`相比，变化幅度是否符合范围
