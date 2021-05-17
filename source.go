package main

import (
	"github.com/lilly/lilly"
	"os"
)

func main() {
	tree := lilly.NewDOMTreeFromURL("https://biz.chosun.com/site/data/html_dir/2021/04/04/2021040401477.html")
	file, _ := os.Create("./extract.txt")
	file.WriteString(tree.ExtractContent())
	//fmt.Println("=====================================================================================")
	accurate, _ := os.Create("./accurate.txt")
	accurate.WriteString(tree.ExtractAccurateContent())
}
