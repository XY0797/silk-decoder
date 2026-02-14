package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"silk-decoder/silk"
)

var (
	input        = flag.String("i", "", "")
	dir          = flag.String("d", "", "")
	sampleRate   = flag.Int("sampleRate", 24000, "")
	verbose      = flag.Bool("verbose", false, "")
	output       = flag.String("o", "", "")
	ffmpegPath   = flag.String("ffmpeg", "", "")
	targetFormat = flag.String("format", "", "")
	pattern      *regexp.Regexp
)

func main() {
	flag.Usage = printUsage
	flag.Parse()
	if *verbose {
		silk.Verbose(*verbose)
	}

	if *input == "" {
		printUsage()
		fmt.Println("[错误] 输入文件必填。")
		os.Exit(1)
	}

	if *targetFormat != "" && *ffmpegPath == "" {
		// 检查 PATH 中是否有 ffmpeg
		if _, err := exec.LookPath("ffmpeg"); err != nil {
			fmt.Println("[错误] 未指定 -ffmpeg 参数，且系统 PATH 中找不到 ffmpeg，请安装或指定路径。")
			os.Exit(1)
		}
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
	buf, err := silk.Decode(in, silk.WithSampleRate(*sampleRate))
	if err != nil {
		in.Close()
		return fmt.Errorf("对输入文件 %q 解码失败: %w", path, err)
	}
	in.Close()

	var pcmSuffix = ".pcm"
	var outputName string

	if *targetFormat != "" {
		// 如果指定了目标格式，调用 ffmpeg 转换
		finalOutput := getOutputName(path, "."+*targetFormat, *output, batch)
		ffmpeg := "ffmpeg"
		if *ffmpegPath != "" {
			ffmpeg = *ffmpegPath
		}

		cmd := exec.Command(
			ffmpeg,
			"-y",
			"-f", "s16le", // PCM 是小端格式
			"-ar", fmt.Sprintf("%d", *sampleRate),
			"-i", "-", // 从 stdin 读取
			finalOutput,
		)

		// buf 作为 stdin 输入
		cmd.Stdin = bytes.NewReader(buf)

		if *verbose {
			fmt.Printf("执行命令: %s\n", strings.Join(cmd.Args, " "))
		}

		outBytes, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("执行 ffmpeg 失败: %w\n输出: %s", err, string(outBytes))
		}

		if *verbose {
			fmt.Printf("已生成: %s\n", finalOutput)
		}
	} else {
		// 否则直接输出 .pcm 或用户指定的 -o
		outputName = getOutputName(path, pcmSuffix, *output, batch)
		out, err := os.OpenFile(outputName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			return fmt.Errorf("打开/创建输出文件 %q 失败: %w", outputName, err)
		}

		_, err = out.Write(buf)
		if err != nil {
			out.Close()
			return fmt.Errorf("写入输出文件 %q 失败: %w", outputName, err)
		}
		out.Close()
	}

	return nil
}

func getOutputName(path, suffix, output string, batch bool) string {
	var outputName string
	if batch { // 批量
		if output != "" {
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
		if output != "" {
			outputName = output
		} else {
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
	fmt.Println("将 silk v3 格式的文件解码为 pcm，再可选地通过 ffmpeg 转换为其他音频格式")
	fmt.Println("原作者：youthlin，XY0797 修改")
	fmt.Println("原项目 GitHub: https://github.com/youthlin/silk")
	fmt.Println()
	fmt.Printf("用法：%s -i <输入文件> [选项]\n", name)
	fmt.Println("  -i <输入文件>\t\t\t输入文件或输入文件夹（需配合 -d 使用）")
	fmt.Println("  [选项]")
	fmt.Println("    -d <正则表达式>\t\t指明 -i 是目录，对匹配正则的文件批量处理（递归处理子目录）")
	fmt.Println("    -sampleRate <采样率>\t采样率（Hz），默认 24000")
	fmt.Println("    -o <输出文件>\t\t指定输出文件名（单文件）或后缀（批量）")
	fmt.Println("    -ffmpeg <路径>\t\t指定 ffmpeg 二进制路径")
	fmt.Println("    -format <格式>\t\t目标音频格式（如 mp3、wav、flac 等）")
	fmt.Println("    -verbose\t\t\t输出调试日志（默认 false）")
	fmt.Println()
	fmt.Println("注意：若使用 -format，则必须确保 ffmpeg 可用（通过 -ffmpeg 指定或确保 PATH 中存在）")
	fmt.Println()
	fmt.Println("示例：")
	fmt.Printf("%s -i a.amr -format mp3\n\t将 a.amr 解码并转换为 a.mp3\n", name)
	fmt.Printf("%s -i voice -d \".*\\.amr\" -format mp3\n\t批量将 voice 文件夹中的 .amr 转为 .mp3\n", name)
	fmt.Printf("%s -i a.amr -format mp3 -ffmpeg /usr/local/bin/ffmpeg\n\t使用指定路径的 ffmpeg\n", name)
	fmt.Println()
}
