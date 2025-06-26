# Dify 插件守护进程

## 概述

Dify 插件守护进程是一个管理插件生命周期的服务。它负责三种类型的运行时：

1. 本地运行时：在与 Dify 服务器相同的机器上运行。
2. 调试运行时：监听一个端口，等待调试插件连接。
3. 无服务器运行时：在 AWS Lambda 等无服务器平台上运行。

Dify API 服务器将与守护进程通信，以获取所有插件的状态，例如哪个插件已安装到哪个工作区，并接收来自 Dify API 服务器的请求以调用插件，如同无服务器函数。

所有来自 Dify API 的请求基于 HTTP 协议，但根据运行时类型，守护进程将以不同方式将请求转发到相应的运行时。

- 对于本地运行时，守护进程将插件作为子进程启动，并通过 STDIN/STDOUT 与插件进行通信。
- 对于调试运行时，守护进程等待插件连接，并以全双工方式进行 TCP 通信。
- 对于无服务器运行时，插件将被打包到 AWS Lambda 等第三方服务中，然后由守护进程通过 HTTP 协议调用。有关更多详细信息，请参阅 SRI 文档。

有关 Dify 插件的更详细介绍，请参阅我们的文档 [https://docs.dify.ai/plugins/introduction](https://docs.dify.ai/plugins/introduction).

## CLI

提供了一个 CLI 工具用于本地环境下的插件开发。

- 通过`brew`安装

Linux 和 MacOS 均支持 arm64 或 amd64 架构。

1. 打开 [ Dify CLI 的 Homebrew tap](https://github.com/langgenius/homebrew-dify)
2. 使用 brew 安装 Dify cli

```bash
brew tap langgenius/dify
brew install dify
```

- 通过二进制文件安装

从发布页面的资产列表中下载二进制文件 [发布页面](https://github.com/langgenius/dify-plugin-daemon/releases).

## 开发

### 运行守护进程

首先，将 `.env.example` 文件复制到 `.env` 并设置正确的环境变量如 `DB_HOST` 等。

```bash
cp .env.example .env
```

如果在版本 0.1.2 之前您使用的是非 AWS S3 存储，则需要手动在 .env 文件中将 S3_USE_AWS 环境变量设置为 false。

请注意，`PYTHON_INTERPRETER_PATH` 是指向 Python 解释器的路径，请根据您的 Python 安装指定正确的路径，并确保 Python 版本为 3.11 或更高，因为 dify-plugin-sdk 要求如此。

我们推荐您使用 `vscode` 来调试守护进程，并在 `.vscode` 目录中提供了一个 `launch.json`.json 文件。



### Python环境
#### UV
守护进程使用 `uv` 来管理插件的依赖项，启动守护进程之前，您需要自行安装 [uv](https://github.com/astral-sh/uv)。

#### 解释器
您的机器上可能安装了多个 Python 版本，提供了一个变量 `PYTHON_INTERPRETER_PATH` 来指定 Python 解释器的路径。

## 部署

目前，守护进程仅支持 Linux 和 MacOS，对于 Windows，需要进行大量适配，欢迎您贡献代码以满足该需求。

### Docker

> **NOTE:** 由于守护进程依赖于共享的 `cwd` 目录来运行插件，因此不建议使用基于网络的卷或从主机外部绑定的挂载。这可能导致性能不佳，例如插件未能及时启动。

使用 Docker 卷与主机共享目录，性能更佳。

### Kubernetes

目前，守护进程社区版不会支持根据副本数平滑扩展，如果您对此功能感兴趣，请与我们联系。我们有一个更适合生产环境的企业版。

## 基准测试

请参阅 [基准测试](https://langgenius.github.io/dify-plugin-daemon/benchmark-data/)

## 许可证

Dify 插件守护进程根据 [Apache-2.0 license](LICENSE)发布。