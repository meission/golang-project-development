## Go 项目目录


一个小型项目用不到很多目录。对于小型项目，可以考虑先包含 cmd、pkg、internal 3 个目录，其他目录后面按需创建

项目目录样例

```

.
├── api
│   └── api.go
├── changelog.md
├── cmd
│   └── main.go
├── configs
│   └── config.toml
├── internal
│   ├── config
│   ├── dao
│   │   └── dao.go
│   ├── data
│   │   └── mysql.go
│   └── service
│       └── service.go
├── LICENSE
├── Makefile
├── pkg
│   └── log
│       └── log.go
├── README.md
|____go.mod
|____go.sum
└── server
    ├── grpc.go
    └── http.go

``` 

### 项目目录说明

#### 根据需要可灵活调整

1. api 目录中存放的是当前项目对外提供的各种不同类型的 API 接口定义文件，其中可能包含类似 protobuf，api、swagger 的目录，这些目录包含了当前项目对外提供和依赖的所有 API 文件。
 

2. cmd一个项目有很多组件，可以把组件 main 函数所在的文件夹统一放在/cmd 目录下，如果有多个组建可以在cmd下创建不同的目录


3. configs这个目录用来配置文件模板或默认配置。例如，可以在这里存放 xml,json,yaml,toml等配置文件样例模板。注意，配置中不能携带敏感信息，这些敏感信息，我们可以用占位符来替代


4. internal
业务代码以及本项目私有代码不希望在其他应用和库中被导入，可以将这部分代码放在/internal 目录下。
可以通过 Go 语言本身的机制来约束其他项目 import 项目内部的包。      
    /internal 目录建议包含如下目录：   
        /internal/dao: 该目录中存放访问数据的逻辑代码，包括数据库，缓存等数据操作，在internal 范围内保持私有共享。       
        /internal/data: 该目录中存放访问数据的结构代码，包括数据库，缓存等自定义的数据，在internal 范围内保持私有共享，项目外不共享的包。
        /internal/serivce: 该目录中存放业务逻辑实现代码。     
        /internal/config: 该目录中存放处理配置文件数据实现代码。       



5. pkg 一般存放本代码仓库较为通用的代码。 目录是 Go 语言项目中非常常见的目录，我们几乎能够在所有知名的开源项目（非框架）中找到它的身影，例如 Kubernetes、Prometheus、Moby、Knative 等。这些包提供了比较基础、通用的功能。  
该目录中存放可以被外部应用使用的代码库，其他项目可以直接通过 import 导入这里的代码。所以，我们在将代码库放入该目录时一定要慎重。

6. server :http和grpc实例的创建和配置

7. CHANGELOG
当项目有更新时，为了方便了解当前版本的更新内容或者历史更新内容，需要将更新记录存放到 CHANGELOG 目录。编写 CHANGELOG 是一个复杂、繁琐的工作，我们可以结合 Angular 规范 和 git-chglog 来自动生成 CHANGELOG。

8. LICENSE
版权文件可以是私有的，也可以是开源的。常用的开源协议有：Apache 2.0、MIT、BSD、GPL、Mozilla、LGPL。有时候，公有云产品为了打造品牌影响力，会对外发布一个本产品的开源版本，所以在项目规划初期最好就能规划下未来产品的走向，选择合适的 LICENSE。   
为了声明版权，你可能会需要将 LICENSE 头添加到源码文件或者其他文件中，这部分工作可以通过工具实现自动化，推荐工具： addlicense 。   
当代码中引用了其它开源代码时，需要在 LICENSE 中说明对其它源码的引用，这就需要知道代码引用了哪些源码，以及这些源码的开源协议，可以借助工具来进行检查，推荐工具： glice 。至于如何说明对其它源码的引用，大家可以参考下 IAM 项目的 LICENSE 文件。   

9. README.md
项目的 README 文件一般包含了项目的介绍、功能、快速安装和使用指引、详细的文档链接以及开发指引等。有时候 README 文档会比较长，为了能够快速定位到所需内容，需要添加 markdown toc 索引，可以借助工具 tocenize 来完成索引的添加。

10. hack  // 编译、构建及校验的工具类

11. build // 构建和测试脚本

12. test 测试

13. third_party  //第三方代码，protobuf、golang-reflect等







