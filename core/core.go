package core

import (
	"bytes"
	"github.com/robfig/cron/v3"
	"github.com/soxft/gitea-backuper/backup"
	"github.com/soxft/gitea-backuper/config"
	"github.com/soxft/gitea-backuper/tool"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func Run() {

	// check work dir and Local exists
	if err := tool.PathExists(config.Gitea.WorkDir); err != nil {
		log.Fatalf("WorkDir %s not exists: %v", config.Gitea.WorkDir, err)
		return
	}
	if err := tool.PathExists(config.Local.Dir); err != nil {
		log.Fatalf("LocalDir %s not exists: %v", config.Local.Dir, err)
		return
	}

	// 启动时立即执行一次
	coreFunc()

	// start cron
	c := cron.New()

	if _, err := c.AddFunc(config.Gitea.Cron, coreFunc); err != nil {
		log.Fatalf("Add Cron error: %v", err)
	} else {
		log.Printf("Cron added: %s", config.Gitea.Cron)
	}
	c.Start()

	// wait for interrupt signal to gracefully shut down the server with
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	c.Stop()
	log.Println("Bye! :)")
}

// The CoreFunc
func coreFunc() {
	defer func() {
		// recover
		if err := recover(); err != nil {
			log.Printf("error %v", err)
		}
	}()

	log.Println("Start backup")

	// execute gitea dump
	_giteaBin := config.Gitea.BinPath

	cmd := exec.Command(_giteaBin, "dump")
	cmd.Dir = config.Gitea.WorkDir // select working directory

	// get stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Println(stdout.String(), stderr.String())
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	_ = cmd.Wait()

	// start process backup logic

	var fList []os.DirEntry
	if fList, err = os.ReadDir(config.Gitea.WorkDir); err != nil {
		log.Printf("error when read dir: %v", err)
		return
	}

	// 遍历fList, 匹配 gitea-dump-*.zip
	var _dumpFileName string
	for _, f := range fList {
		if f.IsDir() {
			continue
		}
		if tool.GetDumpFileName(f.Name()) != "" {
			_dumpFileName = f.Name()
		}
	}

	if _dumpFileName == "" {
		log.Println("dump file not found")
		log.Println(stdout.String(), stderr.String())
		return
	}

	dumpFileAbsPath := config.Gitea.WorkDir + _dumpFileName
	localFileAbsPath := config.Local.Dir + _dumpFileName
	log.Printf("dump file path: %s", dumpFileAbsPath)

	// check dump file exists
	if _, err := os.Stat(dumpFileAbsPath); err != nil {
		log.Printf("dump file %s not exists, %v", dumpFileAbsPath, err)
		return
	}

	// move dump file to local
	err = tool.MoveFile(dumpFileAbsPath, localFileAbsPath)
	if err != nil {
		log.Println("error when move dump file to local", err)
		return
	}

	log.Printf("Move dump file to local: %s", localFileAbsPath)

	// remove local backup files if max file num is set
	err = tool.DeleteLocal(config.Local.Dir, config.Local.MaxFileNum)
	if err != nil {
		log.Println("error when clear local backup files", err)
		return
	}
	log.Printf("Clear local backup files success, max file num: %d", config.Local.MaxFileNum)

	log.Println("Start upload to cos")
	// Upload to remote
	remotePath, err := backup.ToCos(localFileAbsPath)
	if err != nil {
		log.Println("error when upload to cos", err)
		return
	}

	log.Printf("Upload to cos success, remote path: %s", remotePath)
	log.Println("Backup success")
}
