# 企业微信机器人客户端

## 简介

一个可以通过命令行发送企业微信机器人所支持的各种类型消息的小程序。

## 使用

### 添加机器人

```
$ wx robot add <name> <hookURL>
```

### 发送消息

首先切换到指定的机器人

```
$ wx robot use <name>
```

#### 文本消息

```
$(robot1): send [--text|-t] '<text message>'
```
#### Markdown 消息

```
$(robot1): send [--markdown|-m] '<markdown file>'
```

#### 图片消息

```
$(robot1): send [--image|-i] '<image file>'
```

#### 文件消息

```
$(robot1): send [--file|-f] '<file path>'
```

#### 图文消息

```
$(robot1): send [--news|-n] '<news file>'
```

这个图片消息比较复杂，因为字段比较多，所以这里的 `<news file>` 文件里面需要按照下面的格式定义好。

```
articles:
- title: <title 1>
  description: <description 1>
  url: <jump to url 1>
  picurl: <picture url 1>
- title: <title 2>
  description: <description 2>
  url: <jump to url 2>
  picurl: <picture url 2>
```

## 备注

没有什么用途的话，假装机器人也很好玩啊！