原项目为单体应用，改造为是使用golang进行重构的，目前使用的go版本为1.13和go-micro框架。

目前这个简单的重构考虑了如下的内容：
采用了微服务的架构：采用了使用http的BFF，后台的service采用了grpc进行通信。
在服务中增加了crontab job来进行job处理。

API采用了proto格式进行管理
使用gorm框架提供数据库的通信。
在服务中加入了限流和熔断组件。
在服务中增加了Kafka ELK的支持，支持了opentracing，也集成了zap进行集中log处理。
在框架中增加了缓存的使用。

还缺少的内容为：
Admin模块
增加k8s支持
增加新的缓存种类
升级go版本
迁移到kratos或者go-zero