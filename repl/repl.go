package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tanlian/rulego/program"
)

func Start() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("RuleGo REPL")
	fmt.Println("Type 'exit' to quit")
	pro := program.New()

	for {
		fmt.Print(">> ")
		// 读取用户输入
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		input = strings.TrimSpace(input)

		// 检查用户是否输入了退出命令
		if input == "exit" {
			break
		}

		pro.Run(input)
	}
	fmt.Println("Goodbye!")
}
