# HTML2Go

一个双向转换工具，用于在 HTML 和使用 HTMLGo 组件的 Go 代码之间进行转换。

## 概述

HTML2Go 是一个功能强大的工具，设计用于在以下两个方向上实现无缝转换：
1. Go 代码（使用 htmlgo/vuetify/vuetifyx 组件）到 HTML
2. HTML 到 Go 代码（使用 htmlgo/vuetify/vuetifyx 组件）

这个工具对于使用基于 Go 的 Web 应用程序并利用 Vuetify 和 VuetifyX 等 HTML 组件库的开发人员特别有用。

## 功能特点

- **双向转换**：支持两个方向的转换 - Go 代码到 HTML 和 HTML 到 Go 代码
- **组件支持**：
  - 标准 HTML 元素
  - Vuetify 组件（`v-btn`, `v-card` 等）
  - VuetifyX 组件（`v-xbtn`, `v-xcard` 等）
  - 自定义组件
- **属性处理**：正确映射 Go 方法调用和 HTML 属性
- **布尔属性**：特殊处理 `disabled`, `required` 等布尔属性
- **嵌套结构**：支持复杂的嵌套组件结构

## 安装

### 从源代码构建

```bash
git clone https://github.com/zhangshanwen/html2go.git
cd html2go
go build
```

### 使用 Go Install 安装

```bash
go install github.com/zhangshanwen/html2go@latest
```

## 使用方法

### HTML 转换为 Go 代码

```bash
# 基本用法
echo '<div>Hello, world!</div>' | html2go

# 指定包名
echo '<div>Hello, world!</div>' | html2go -p mypackage

# 指定 vuetify 包名
echo '<v-btn>Submit</v-btn>' | html2go -v vuetify

# 指定 vuetifyX 包名
echo '<v-xbtn>Submit</v-xbtn>' | html2go -vx vuetifyx

# 使用 Children 模式生成
echo '<div><span>Text</span></div>' | html2go -c

# 保存到文件
echo '<div>Hello, world!</div>' | html2go > output.go

# 处理 HTML 文件
cat input.html | html2go > output.go
```

### Go 代码转换为 HTML

```bash
# 使用 -r 标志将 Go 代码转换为 HTML
echo 'Div(Text("Hello, world!"))' | html2go -r

# 指定 HTML 缩进大小
echo 'Div(Text("Hello, world!"))' | html2go -r -indent 4

# 禁用 HTML 格式化
echo 'Div(Text("Hello, world!"))' | html2go -r -format=false

# 处理 Go 文件
cat input.go | html2go -r > output.html
```

### 命令行参数

#### HTML 到 Go 转换参数
- `-p <package>`: 指定生成代码的包名
- `-v <package>`: 指定 Vuetify 组件的包名
- `-vx <package>`: 指定 VuetifyX 组件的包名
- `-c`: 启用 Children 模式生成代码

#### Go 到 HTML 转换参数
- `-r`: 启用反向模式（Go 代码到 HTML）
- `-indent <size>`: 设置 HTML 缩进大小（默认：2）
- `-format <bool>`: 是否格式化生成的 HTML（默认：true）

## 示例

### HTML 到 Go 代码转换

输入（HTML）:
```html
<div class="container">
  <h1>Hello, World!</h1>
  <v-btn color="primary" text="Submit"></v-btn>
</div>
```

输出（Go）:
```go
package main

var n = Body(
    Div(
        Class("container"),
        H1(
            Text("Hello, World!"),
        ),
        VBtn(
            Color("primary"),
            Text("Submit"),
        ),
    ),
)
```

### Go 到 HTML 转换

输入（Go）:
```go
Div(
    Class("container"),
    H1(
        Text("Hello, World!"),
    ),
    VBtn(
        Color("primary"),
        Text("Submit"),
    ),
)
```

输出（HTML）:
```html
<div class="container">
  <h1>Hello, World!</h1>
  <v-btn color="primary" text="Submit"></v-btn>
</div>
```

## 支持的组件

- **HTML 元素**: div, span, p, h1-h6, input, button, form 等
- **Vuetify 组件**: v-btn, v-card, v-text-field, v-select 等
- **VuetifyX 组件**: v-xbtn, v-xcard, v-xdialog 等
- **自定义组件**: 支持自定义定义的组件

## 布尔属性

布尔属性（如 `disabled`, `required` 等）会被特殊处理：
- 在 Go 代码中: `Disabled(true)` 
- 在 HTML 中: `disabled`（没有值）

## 高级用法

### 自定义组件处理

对于不直接映射的自定义组件，使用 `Tag` 函数：

```go
Tag("custom-element", 
    Attr("custom-attribute", "value"),
    Text("Content"),
)
```

### 使用 Children 方法

对于接受子元素的组件：

```go
Div(
    Children(
        Span(Text("子元素 1")),
        Span(Text("子元素 2")),
    ),
)
```

### 处理 Alpine.js 属性

HTML2Go 支持 Alpine.js 属性的处理：

```html
<div x-data="{open: false}" x-on:click="open = true">
  点击打开
</div>
```

转换为 Go 代码：

```go
Div(
    Attr("x-data", "{open: false}"),
    Attr("x-on:click", "open = true"),
    Text("点击打开"),
)
```

## 贡献

欢迎贡献！请随时提交 Pull Request。

1. Fork 仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 许可证

本项目使用 MIT 许可证 - 详情请查看 LICENSE 文件。
