# silk-decoder

Go 语言版本的 Silk v3 解码程序。

可用于解码国内通信软件(微信/QQ)语音文件，得到`.pcm`文件。

使用ffmpeg即可把`.pcm`转换为任意的其它格式(`ffmpeg -y -f s16le -ar 24000 -i test.pcm test.mp3`)。

## 编译

编译当前平台的版本：

```sh
go build -ldflags="-s -w" -trimpath
```

## 用法

```
silk-decoder -i <输入文件> [选项]
  -i <输入文件>                 输入文件或输入文件夹(需要和 -d 连用)
  [选项]
    -d <正则表达式>             指明 -i 的参数是文件夹，对输入文件夹(及子文件夹中)中，文件名符合正则表达式的文件进行解码
    -sampleRate <采样率>        单位为赫兹，默认值为 24000
    -o <输出文件>               指定输出文件名，或指定输出文件后缀名（当使用-d 时）。
                                如果为空则自动推断
    -verbose                    输出调试日志(默认值为 false)

示例：
silk-decoder -i a.amr
        将 a.amr 解码为 a.pcm
silk-decoder -i amr.1
        将 amr.1 解码为 amr.pcm
silk-decoder -i file
        将 file 解码为 file.pcm
silk-decoder -i a.amr -o b.pcm
        将 a.amr 解码为 b.pcm
silk-decoder -i voice -d ".*\.amr"
        将当前文件夹中的所有 .amr 文件转换为 .pcm 文件
          例如：voice 文件夹下有如下文件：
                voice/a.amr
                voice/other.txt
                voice/sub/b.amr
          转换结果：
                voice/a.pcm
                voice/sub/b.pcm
```

### 转mp3示例

首先确保已经下载ffmpeg并且成功添加到环境变量中。

执行下面命令：

```sh
silk-decoder -i test.amr -o test.pcm
ffmpeg -y -f s16le -ar 24000 -i test.pcm test.mp3
```

如果需要接近无损的音质，可以使用下面命令(对于说话录音来说没必要)：

```sh
silk-decoder -i test.amr -sampleRate 44100 -o test.pcm
ffmpeg -y -f s16le -ar 44100 -i test.pcm test.mp3
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
