package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/zhangshanwen/html2go/parse"
)

func main() {
	// HTML 到 Go 的相关参数
	pkg := flag.String("p", "", "package name")
	vuetifyPkg := flag.String("v", "", "vuetify package name")
	vuetifyxPkg := flag.String("vx", "", "vuetify-x package name")
	childrenMode := flag.Bool("c", false, "generate children pattern")

	// Go 到 HTML 的相关参数

	indentSize := flag.Int("indent", 2, "HTML缩进大小")
	reverseMode := flag.Bool("r", false, "启用反向模式 (Go 代码到 HTML)")
	formatHTML := flag.Bool("format", true, "是否格式化生成的 HTML")

	flag.Parse()

	// 从标准输入读取数据
	inputBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取输入时出错: %v\n", err)
		os.Exit(1)
	}
	inputString := string(inputBytes)

	if *reverseMode {
		// 反向模式：Go 代码到 HTML
		html, err := parse.GenerateGoHTML(inputString, *formatHTML, *indentSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "生成 HTML 时出错: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(html)
	} else {
		// 默认模式：HTML 到 Go 代码
		goCode := parse.GenerateHTMLGo(*pkg, *vuetifyPkg, *vuetifyxPkg, *childrenMode, strings.NewReader(inputString))
		fmt.Println(goCode)
	}
}
