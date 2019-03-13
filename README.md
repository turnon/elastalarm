# ElastAlarm

## 运行

```
elastalarm -host=http://192.168.0.60:9200 -configs=/path/to/config_dir
```

或者

```
docker run -v /path/to/config_dir:/configs -d --env-file env.list \
  daocloud.io/shutdown/elastalarm:latest -host http://192.168.0.60:9200
```

## 统计方法和配置方式

配置例子可见`./example_configs/*.json`，配置文件的各字段意义见`./example_configs/doc.md`，统计方式解释如下

* `count`： `interval`内根据`scope`条件所查出doc数量，是否符合范围
* `percentage`： `interval`内`part`除以`whole`所得的百分比，是否符合范围
* `spike`： 最近`interval`的`scope`与上一`interval`的`scope`相比，变化幅度是否符合范围

## 环境变量

用于邮件通知

* `ESALARM_MAIL_SERVER=smtp.qq.com:465`
* `ESALARM_MAIL_SKIP_VERIFY=true`（跳过SSL）
* `ESALARM_MAIL_FROM=someone@qq.com`
* `ESALARM_MAIL_PASSWD=klkjjfkj5645`

用于钉钉通知

* `ESALARM_DING_CORPID=ding134nbcvbmn`
* `ESALARM_DING_SECRET=S-mbvfhbfjhgjh657nvnjvhf74`
* `ESALARM_DING_AGENT=12345`
