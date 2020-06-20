# 企业微信机器人客户端

## 简介

一个可以通过命令行发送企业微信机器人所支持的各种类型消息的小程序。

## 使用

### 添加机器人

```
$ wxrobot add <name> <hookURL>
```

### 查看所有机器人
```
$ wxrobot list
```

### 发送消息

首先切换到指定的机器人

```
$ wxrobot use <name>
```

切换到指定的机器人时，命令行会输出如下的信息：

```
Run command `source /var/folders/_x/pwjgfy4n50g6mqhmh39bs59w0000gp/T/xiaomi_bash_profile` to make the robot bash profile effective
```

可以将输出中的 source 命令复制出来放到 Shell 中执行，可以改变当前 Shell 会话窗口的 Prompt 的文字，方便识别正在使用哪个机器人，如果系统不支持（如Windows），则需要在发送消息的时候指定机器人名称。

```
$ source /var/folders/_x/pwjgfy4n50g6mqhmh39bs59w0000gp/T/xiaomi_bash_profile
```

#### 文本消息

```
$ wxrobot send [--robot|-r] <name> [--text|-t] '<text message>'
```
#### Markdown 消息

```
$ wxrobot send [--robot|-r] <name> [--markdown|-m] '<markdown file>'
```

#### 图片消息

```
$ wxrobot send [--robot|-r] <name> [--image|-i] '<image file>'
```

#### 文件消息

```
$ wxrobot send [--robot|-r] <name> [--file|-f] '<file path>'
```

#### 图文消息

```
$ wxrobot send [--robot|-r] <name> [--news|-n] '<news file>'
```

这个图片消息比较复杂，因为字段比较多，所以这里的 `<news file>` 文件里面需要按照下面的格式定义好。

```
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