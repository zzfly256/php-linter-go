# php-linter-go

## 简介
一个高性能 PHP 语法检测工具

## 安装

### 方法一：直接下载二进制文件(文件)
点击 [此处](https://github.com/zzfly256/php-linter-go/releases) 前往下载页，直接下载最新版本的二进制文件，赋予可执行权限即可使用。

```shell
## 赋予执行权限(仅Linux/Mac)
chmod +x [文件路径]/php-linter-go
```

### 方法二：从 go 源码编译

```shell
## 下载源码
git clone https://github.com/zzfly256/php-linter-go.git

## 编译(Linux/Mac)
go build -o php-linter-go php-linter.go
## 赋予执行权限
chmod +x ./php-linter-go
## 复制到系统路径(可选)
cp ./php-linter-go /usr/bin/php-linter-go

## 编译(Windows)
go build -o php-linter-go.exe php-linter.go
```

## 使用方法

| 选项 | 作用 | 说明 |
|:---|:---|:---|
| --help | 显示帮助菜单 |  |
| --path | 检测指定目录下的所有 PHP 文件语法 | 若不传该参数，且无 `--git` `--svn` 选项，则默认检测当前文件夹 |
| --recursive | 递归检测全部目录下的所有 PHP 文件语法 |  |
| --git | 检测 Git 中的变更 | 增量检测；若当前目录为 Git 仓库，则将**当前分支/当前变更**与 `origin/master` 作比较，对变动过的 PHP 文件做语法检查 |
| --svn | 检测 SVN 中的变更 | 增量检测；若当前目录为 SVN 仓库，则将当前文件夹与**上一次提交**作比较，对变动过的 PHP 文件做语法检查 |
| --php-executor | 指定 PHP 解释器 | 若本地有多个 PHP 版本，可使用该参数指定 PHP 解释器文件路径 |


### 场景一：检测当前文件夹全部文件

1. 进入要检测的文件夹

2. 直接执行命令 `php-linter-go`

![img](https://github.com/zzfly256/php-linter-go/raw/master/doc/images/sense1.png)

### 场景二：检测当前 Git 仓库中所有变动过的 PHP 文件

1. 进入 Git 仓库下的任意一个地方

2. 直接执行命令 `php-linter-go --git`

![img](https://github.com/zzfly256/php-linter-go/raw/master/doc/images/sense2.png)

### 场景三：递归监测某个目录内全部 PHP 文件

1. 直接执行命令 `php-linter-go  --recursive --path=/home/www/demo`

### 场景四：使用自定义的 PHP 版本做语法检测

假设在我的机器上安装的 PHP 5.6 位于 `/opt/homebrew/Cellar/php@5.6/5.6.40` 目录，bin 可执行文件地址为 `/opt/homebrew/Cellar/php@5.6/5.6.40/bin/php`

1. 直接执行命令 `php-linter-go --php-executor=/opt/homebrew/Cellar/php@5.6/5.6.40/bin/php`

![img](https://github.com/zzfly256/php-linter-go/raw/master/doc/images/sense4.png)

## 其他

### 关于性能
程序使用了 Go 语言调度 PHP 解释器执行 `php -l` 命令，并控制最大 8 条协程同时进行。 在 Macbook Air M1 2020 中使用 PHP 7.4 实测一秒约检测 100 + 个 PHP 文件，相比直接编写 shell 执行语法检测快了约 **66.69%**。

检测结果脱敏后如下：

```shell
#  ls -l | grep .php | wc -l
136

#  time php-linter

Simple PHP Linter (PHP 7.4.21)

Result:
Check: 136 / Errors: 1

php-linter  1.97s user 0.97s system 88% cpu 3.329 total

#  time php-linter-go

Simple PHP Linter (PHP 7.4.21)

Result:
Check: 136 / Errors: 1

php-linter-go  3.62s user 2.92s system 589% cpu 1.109 total
```

### TODO List
1. 打包二进制文件，接入 composer
2. 支持 `--git-path` 、`--svn-path` 等代码仓库位置指定