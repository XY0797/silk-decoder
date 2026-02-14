package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"silk-decoder/silk"
)

var (
	input      = flag.String("i", "", "")
	dir        = flag.String("d", "", "")
	sampleRate = flag.Int("sampleRate", 24000, "")
	verboe     = flag.Bool("verbose", false, "")
	output     = flag.String("o", "", "")
	pattern    *regexp.Regexp
)

func main() {
	flag.Usage = printUsage
	flag.Parse()
	if *verboe {
		silk.Verbose(*verboe)
	}

	if *input == "" {
		printUsage()
		fmt.Println("[错误] 输入文件必填。")
		os.Exit(1)
	}

	if *dir == "" { // input file
		if err := decodeOneFile(*input, false); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		return
	}

	// input dir
	exp, err := regexp.Compile(*dir)
	if err != nil {
		fmt.Printf("[错误] 输入文件的正则表达式 %s 无法识别：%+v\n", *dir, err)
		os.Exit(1)
	}
	pattern = exp

	filepath.WalkDir(*input, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !pattern.MatchString(d.Name()) {
			return nil // ignore
		}
		if err = decodeOneFile(path, true); err != nil {
			fmt.Println(err.Error()) // ignore error
		}
		return nil
	})

}

func decodeOneFile(path string, batch bool) error {
	in, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("打开输入文件 %q 失败: %w", path, err)
	}
	defer in.Close()

	buf, err := silk.Decode(in, silk.WithSampleRate(*sampleRate))
	if err != nil {
		return fmt.Errorf("对输入文件 %q 解码失败: %w", path, err)
	}

	var suffix = ".pcm"

	var outputName = getOutputName(path, suffix, *output, batch)

	out, err := os.OpenFile(outputName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("打开/创建输入文件 %q 失败: %w", outputName, err)

	}
	defer out.Close()

	_, err = out.Write(buf)
	if err != nil {
		return fmt.Errorf("写入输出文件 %q 失败: %w", outputName, err)
	}
	return nil
}

func getOutputName(path, suffix, output string, batch bool) string {
	var outputName string
	if batch { // 批量
		if output != "" { // 指定输出后缀名
			if strings.HasPrefix(output, ".") {
				suffix = output
			} else {
				suffix = "." + output
			}
		}
		if i := strings.LastIndex(path, "."); i > 0 {
			outputName = path[:i] + suffix
		} else {
			outputName = path + suffix
		}
	} else { // 单个
		if output != "" { // 指定了输出名
			outputName = output
		} else { // 没指定输出文件名，默认为 pcm
			if i := strings.LastIndex(path, "."); i > 0 {
				outputName = path[:i] + suffix
			} else {
				outputName = path + suffix
			}
		}
	}
	return outputName
}

func printUsage() {
	base := filepath.Base(os.Args[0])
	name := strings.TrimSuffix(base, filepath.Ext(base))
	fmt.Println()
	fmt.Println("Silk 解码器，Go 语言版本，基于 v1.0.9 的 C 语言版本")
	fmt.Println("将 silk v3 格式的文件解码为 pcm, 原作者：youthlin, XY0797修改")
	fmt.Println("原项目GitHub: https://github.comyouthlin/silk")
	fmt.Println()
	fmt.Printf("用法：%s -i <输入文件> [选项]\n", name)
	fmt.Println("  -i <输入文件>\t\t\t输入文件或输入文件夹(需要和 -d 连用)")
	fmt.Println("  [选项]")
	fmt.Println("    -d <正则表达式>\t\t指明 -i 的参数是文件夹，对输入文件夹(及子文件夹中)中，文件名符合正则表达式的文件进行解码")
	fmt.Println("    -sampleRate <采样率>\t单位为赫兹，默认值为 24000")
	fmt.Println("    -o <输出文件>\t\t指定输出文件名，或指定输出文件后缀名（当使用-d 时）。\n\t\t\t\t如果为空则自动推断")
	fmt.Println("    -verbose\t\t\t输出调试日志(默认值为 false)")
	fmt.Println()
	fmt.Println("示例：")
	fmt.Printf("%s -i a.amr\n\t将 a.amr 解码为 a.pcm\n", name)
	fmt.Printf("%s -i amr.1\n\t将 amr.1 解码为 amr.pcm\n", name)
	fmt.Printf("%s -i file\n\t将 file 解码为 file.pcm\n", name)
	fmt.Printf("%s -i a.amr -o b.pcm\n\t将 a.amr 解码为 b.pcm\n", name)
	fmt.Printf("%s -i voice -d \".*\\.amr\"\n\t将当前文件夹中的所有 .amr 文件转换为 .pcm 文件\n\t  例如：voice 文件夹下有如下文件：\n\t\tvoice/a.amr\n\t\tvoice/other.txt\n\t\tvoice/sub/b.amr\n\t  转换结果：\n\t\tvoice/a.pcm\n\t\tvoice/sub/b.pcm\n", name)
	fmt.Println()
}
