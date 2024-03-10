# logagent
这是一个基于etcd的日志收集热部署模块
用户只需要修改所需日志项的配置文件，该模块会自动将配置文件放入etcd中并在程序运行时实现热部署，基本流程如下：
![logagent_lv1](https://github.com/Sillyhumans/logagent/assets/107688428/c01a16c7-f71a-444f-8a74-e6d5384b1667)


模块细分
* ### Kafka
  连接kafka，并开启一个goroutine从日志通道中循环读取日志
* ### Taillog
  开启一个goroutine循环获取etcd中的日志项配置，正对每一个配置项创建一个tailObj，并开启一个goroutine将日志读入日志通道内
* ### Etcd
  加载用户日志项的配置文件，并且设置哨兵监控配置项实现热部署
* ### Consumer
  开启一个goroutine从kafka中读取日志放入Consumer日志通道内
* ### Es
  开启一个goroutine从Consumer日志通道中读取日志
![logagent_lv2](https://github.com/Sillyhumans/logagent/assets/107688428/2f6208d4-9943-4614-bfbe-8b5d2a282f63)
