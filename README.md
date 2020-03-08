> yq2blog是一个利用webhook将语雀文档同步到博客的小工具

## 前提准备
1.拥有一个语雀知识库，记录知识库的Namespace

例如你的知识库地址是https://www.yuque.com/xx-xx/xxxxx
，则该知识库的namespace为xx-xx/xxxxx

如果不想其他用户看到你的语雀知识库，请将知识库设置为私有，创建并记录你的token

公开的知识库不需要使用token

2.已经拥有一个博客

## 使用方法
在程序同目录下准备配置文件config.yaml
```yaml
yuque:
  #语雀用户token，仅当知识库私有时需要
  token: 
  #语雀文档知识库Namespace
  repo: 
blog:
  #博客种类，暂仅支持hugo
  type: hugo
  #hugo部署方式，local/git
  deloyment: local
  #hugo博客的路径，类型为local，则为本地路径；类型为git，则为git地址
  path: 
```
运行程序，当知识库内有文档变更，博客会同步变更
