# yaml

YAML 文件使用 .yml 或 .yaml 扩展名，并遵循特定的语法规则。

YAML 最常见的用途之一是创建配置文件。相比 JSON，因为 YAML 有更好的可读性，对用户更友好，所以通常建议用 YAML 来编写配置文件，尽管它们在大多数情况下可以互换使用。

除了在 Ansible 中使用之外，YAML 还用于 Kubernetes 资源和部署。

使用 YAML 的一大好处是，YAML 文件可以添加到源代码控制中，比如 Github，这样就可以跟踪和审计变更。


YAML 语法
YAML 的语法如下：
使用缩进表示层级关系。
缩进时不允许使用 Tab 键，只允许使用空格。
缩进的空格数不重要，但同级元素必须左侧对齐。
大小写敏感。



安装

```
go get gopkg.in/yaml.v3

```

```
package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type T struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

func main() {
	t := T{}

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- t:\n%+v\n\n", t)
}
```