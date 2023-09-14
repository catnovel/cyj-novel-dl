# 项目名称

该项目是一个用于下载次元姬小说的命令行工具。

## 依赖

该项目依赖以下第三方库：

- `github.com/catnovelapi/cyj`
- `github.com/tidwall/gjson`

## 安装

在项目根目录下执行以下命令安装依赖：

```shell
go mod tidy
```

## 使用方法

运行以下命令来下载小说：

```shell
go run main.go -d <bookId> -o <saveDir>
```

参数说明：

```  
  -a string
        login account
  -d string
        add bookid to download book
  -o string
        save file name (default "books")
  -s string
        Search book name to download book
  -t string
        set login token
```

- `-d`：要下载的书籍的ID。
- `-o`：保存文件的目录，默认为`books`。
- `-a`：登入账号。
- `-t`：设置登录令牌。

或者，你可以使用以下命令来搜索并下载小说：

```shell
go run main.go -s <bookName> -o <saveDir> 
```

参数说明：

- `-s`：要搜索并下载的书籍的名称。

## 示例

下载指定ID的书籍：

```shell
go run main.go -d 12345 -o books
```

搜索并下载指定名称的书籍：

```shell
go run main.go -s "小说名称" -o books 
```

## 注意事项

- 请确保提供正确的账号和登录令牌。
- 下载的小说将保存在指定的目录中。 
 