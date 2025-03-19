package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zhangshanwen/html2go/parse"
)

var pkg = flag.String("pkg", "", "generated htmlgo pkg name")
var vuetifyPkg = flag.String("vpkg", "", "generated vuetify pkg name")
var vuetifyxPkg = flag.String("vxpkg", "", "generated vuetifyx pkg name")
var childrenMode = flag.Bool("c", false, "children mode")

func main() {
	flag.Parse()

	fmt.Println(parse.GenerateHTMLGo(*pkg, *vuetifyPkg, *vuetifyxPkg, *childrenMode, os.Stdin))
}
