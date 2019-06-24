# ElastAlarm

## 运行

```
elastalarm -configs=/path/to/config_dir
```

或者

```
docker run -v /path/to/config_dir:/configs -d --env-file env.list \
  daocloud.io/shutdown/elastalarm:latest
```

## 统计方法和配置方式

配置例子可见`./example_configs/*.json`，配置文件的各字段意义及统计方式，详见`./example_configs/doc.md`

## 环境变量

ElasticSearch服务地址

* `ESALARM_HOST=http://192.168.0.60:9200`

用于邮件通知

* `ESALARM_MAIL_SERVER=smtp.qq.com:465`
* `ESALARM_MAIL_SKIP_VERIFY=true`（跳过SSL）
* `ESALARM_MAIL_FROM=someone@qq.com`
* `ESALARM_MAIL_PASSWD=klkjjfkj5645`

用于钉钉通知

* `ESALARM_DING_CORPID=ding134nbcvbmn`
* `ESALARM_DING_SECRET=S-mbvfhbfjhgjh657nvnjvhf74`
* `ESALARM_DING_AGENT=12345`
