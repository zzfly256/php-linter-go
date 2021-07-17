package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
)

// php-linter-go 一个高性能 PHP 语法检测工具
// 用 Golang 给 PHP 做嫁衣的感觉真不错
// by Rytia <admin@zzfly.net>, 2021-07-17 周六


var wg sync.WaitGroup
var result sync.Map
var totalFile uint32
var errFile uint32

// PHP 可执行文件路径
var phpExecutor string = "php"

// 控制最大运行协程个数
var ch = make(chan bool, 8)

func main() {

	// 获取运行路径
	wd, _ := os.Getwd()

	isRecursive := flag.Bool("recursive", false, "是否递归执行")
	isGit := flag.Bool("git", false, "检测 Git 中的变更")
	isSvn := flag.Bool("svn", false, "检测 SVN 中的变更")
	inputPath := flag.String("path", wd, "检测地址(默认为当前目录)")
	inputPhpExecFile := flag.String("php-executor", "php", "检测地址(默认为当前目录)")
	flag.Parse()

	phpExecutor = *inputPhpExecFile
	printWelcome()

	if *isGit {
		lintGit()
	} else if *isSvn {
		lintSvn()
	} else {
		if inputPath != nil {
			lintPath(*inputPath, *isRecursive)
		} else {
			color.HiRed("Unknown error for getting path")
		}
	}

	wg.Wait()
	printResult()
	os.Exit(0)

}

func printWelcome() {
	getPhpVersionCmd := phpExecutor + " -r\"echo PHP_VERSION;\""

	stdout, err := execCommand(getPhpVersionCmd)
	if err != nil {
		color.HiRed("PHP execute error. Pleas ensure that php is installed and \"-php-executor\" parameter is correct")
		os.Exit(255)
	}

	fmt.Printf("\nSimple PHP Linter (PHP %s) \n\n", stdout)
}

func printResult() {
	fmt.Println("\n---------------------------")
	fmt.Println("Result:")
	fmt.Printf("Check: %d / Errors: %d \n", totalFile, errFile)

	if errFile > 0 {
		result.Range(processResultPrint)
	}

	fmt.Println("\n")
}

// 获取执行目录下的文件，执行检测
func lintPath(path string, recursive bool) {

	// 获取文件
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}

	// 逐个检查
	for _, file := range dir {
		if recursive && file.IsDir() {
			lintPath(path+"/"+file.Name(), recursive)
		} else if strings.Contains(file.Name(), ".php") {
			wg.Add(1)
			go processPhpLint(path, file.Name())
		}
	}
}

// 获取 Git 文件差异，执行检测
func lintGit() {
	getGitRootPathCmd := "git rev-parse --show-toplevel"
	gitRootPath, err := execCommand(getGitRootPathCmd)
	if err != nil {
		color.HiRed("Git is not installed or work directory is not a git repository")
		os.Exit(255)
	}
	gitRootPath = strings.Trim(gitRootPath, "\n")
	fmt.Printf("Git Root:\t%s\n", gitRootPath)

	getGitBranchCmd := "git symbolic-ref --short -q HEAD"
	gitBranch, err := execCommand(getGitBranchCmd)
	if err != nil {
		color.HiRed("Git branch error")
	}
	fmt.Printf("Git Branch:\t%s\n", gitBranch)

	getGitDiffFilesCmd := "git diff origin/master --name-only|grep .php"
	gitDiffFileString, _ := execCommand(getGitDiffFilesCmd)
	gitDiffFiles := strings.Split(gitDiffFileString, "\n")

	for _, file :=range gitDiffFiles{
		// 防止空串影响
		if len(file) != 0 {
			wg.Add(1)
			go processPhpLint(gitRootPath, file)
		}
	}
}

// 获取 SVN 文件差异，执行检测
func lintSvn() {
	getSvnUrlCmd := "svn info --show-item url"
	svnUrl, err := execCommand(getSvnUrlCmd)
	if err != nil {
		color.HiRed("SVN is not installed or work directory is not a svn repository")
		os.Exit(255)
	}
	svnUrl = strings.Trim(svnUrl, "\n")
	fmt.Printf("SVN Url:\t%s\n", svnUrl)

	svnPath, _ := os.Getwd()
	fmt.Printf("SVN Path:\t%s\n\n", svnPath)

	getSvnDiffFilesCmd := "svn diff --summarize|grep .php|awk '{print $2}'"
	svnDiffFileString, _ := execCommand(getSvnDiffFilesCmd)
	svnDiffFiles := strings.Split(svnDiffFileString, "\n")

	for _, file :=range svnDiffFiles{
		// 防止空串影响
		if len(file) != 0 {
			wg.Add(1)
			go processPhpLint(svnPath, file)
		}
	}
}

// 处理 PHP 语法检测
func processPhpLint(path string, file string) {
	defer wg.Done()

	// 控制最大协程数
	ch <- true

	// 拼接文件真实路径
	realFilePath := file
	if len(path) != 0 {
		realFilePath = path + "/" + file
	}

	lintPhpCmd := phpExecutor + " -l " + realFilePath
	stdout, err := execCommand(lintPhpCmd)

	if err != nil {
		fmt.Println(color.HiRedString("[ERR]"), file, err.Error())
		result.Store(path+"/"+file, string(stdout))
		atomic.AddUint32(&errFile, 1)
	} else {
		fmt.Println(color.HiGreenString("[OK] "), file)
	}

	atomic.AddUint32(&totalFile, 1)

	<-ch
}

// 遍历输出检测结果
func processResultPrint(key, value interface{}) bool {
	file := strings.Trim(fmt.Sprintf("%v", key), " ")
	fileString := color.CyanString("[%s]", file)
	infoString := strings.Trim(fmt.Sprintf("%v", value), "\n")
	fmt.Printf("\n%s\n%s\n", fileString, infoString)
	return true
}

// 简易执行 shell 命令
func execCommand(command string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", command)
	stdout, err := cmd.Output()

	return string(stdout), err
}
