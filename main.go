package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var logFile = false

func main() {
	_dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	_path, err := exec.LookPath(os.Args[1])
	if err != nil {
		fmt.Printf("Didn't find executable\n")
	} else {
		fmt.Printf("executable is in path: %s\n", _path)
		for {
			cmd := exec.Command(_path, os.Args[2:]...) // 使用你的可执行文件的路径替换

			cmd.Dir = _dir // 设置工作目录
			var stdoutFile, stderrFile *os.File
			if logFile {
				// 打开日志文件，如果文件不存在则创建
				stdoutFile, err = os.OpenFile("stdout.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}
				stderrFile, err = os.OpenFile("stderr.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}

				cmd.Stdout = stdoutFile // 标准输出重定向到文件
				cmd.Stderr = stderrFile // 标准错误重定向到文件
			} else {
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
			}

			err = cmd.Start()
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("The process has started with pid %d", cmd.Process.Pid)

			err = cmd.Wait()

			if logFile {
				stdoutFile.Close() // 关闭文件
				stderrFile.Close() // 关闭文件
			}

			if err != nil {
				log.Printf("The process has exited with error %s, restarting", err)
			} else {
				log.Print("The process has exited correctly, restarting")
			}

			time.Sleep(5 * time.Second) // 等待1秒以避免过度使用CPU
		}
	}
}
