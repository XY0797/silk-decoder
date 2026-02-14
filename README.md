# silk-decoder

Go 语言版本的 Silk v3 解码程序。

可用于解码国内通信软件(微信/QQ)语音文件，得到`.pcm`文件。

如果机器中存在`ffmpeg`，可以指定`-format`自定义输出格式。

## 编译

编译当前平台的版本：

```sh
go build -ldflags="-s -w" -trimpath
```

## 用法

```
silk-decoder -i <输入文件> [选项]
  -i <输入文件>                 输入文件或输入文件夹（需配合 -d 使用）
  [选项]
    -d <正则表达式>             指明 -i 是目录，对匹配正则的文件批量处理（递归处理子目录）
    -sampleRate <采样率>        采样率（Hz），默认 24000
    -o <输出文件>               指定输出文件名（单文件）或后缀（批量）
    -ffmpeg <路径>              指定 ffmpeg 二进制路径
    -format <格式>              目标音频格式（如 mp3、wav、flac 等）
    -verbose                    输出调试日志（默认 false）

注意：若使用 -format，则必须确保 ffmpeg 可用（通过 -ffmpeg 指定或确保 PATH 中存在）

示例：
silk-decoder -i a.amr -format mp3
        将 a.amr 解码并转换为 a.mp3
silk-decoder -i voice -d ".*\.amr" -format mp3
        批量将 voice 文件夹中的 .amr 转为 .mp3
silk-decoder -i a.amr -format mp3 -ffmpeg /usr/local/bin/ffmpeg
        使用指定路径的 ffmpeg
```

### 转mp3示例

首先确保已经下载ffmpeg并且成功添加到环境变量中。

执行下面命令：

```sh
silk-decoder -i test.amr -format mp3
```

如果需要接近无损的音质，可以使用下面命令(对于说话录音来说没必要)：

```sh
silk-decoder -i test.amr -sampleRate 44100 -format mp3
```

## 致谢

- https://github.comyouthlin/silk 原项目
- https://github.com/gaozehua/SILKCodec 源码
- https://github.com/kn007/silk-v3-decoder 兼容国内软件的版本
- https://github.com/wdvxdr1123/go-silk ccgo 转写为 go 的版本
- https://github.com/zxfishhack/go-silk 可直接转 wav 的版本
- [Go语言高级编程 - 第 2 章 CGO 编程](https://chai2010.cn/advanced-go-programming-book/ch2-cgo/index.html)

## LICENSE

MIT

C 源码开源协议见每个文件头部注释。
