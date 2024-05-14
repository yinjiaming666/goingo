//go:build ignore

package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"go_udp/utils"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"time"
)

func main() {
	proName := utils.GetConfig("server", "server", "pro_name")
	userName := utils.GetConfig("server", "server", "username")
	password := utils.GetConfig("server", "server", "password")
	ip := utils.GetConfig("server", "server", "ip")

	cgo := utils.GetConfig("build", "build", "CGO_ENABLED")
	goos := utils.GetConfig("build", "build", "GOOS")
	goarch := utils.GetConfig("build", "build", "GOARCH")

	cmd := exec.Command("/bin/bash", "deploy.sh", proName, goos, goarch, cgo)
	res, err := cmd.Output()
	if err != nil {
		//fmt.Printf("  out:\n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err.Error())
	}
	fmt.Printf(string(res))

	if ip == "" || userName == "" || password == "" {
		return
	}

	// 压缩包上传服务器
	localPath := proName + ".tar.gz"
	remotePath := "/data/" // 服务器路径
	start := time.Now()
	sftpClient, err := sftpConnect(userName, password, ip, 22)
	if err != nil {
		log.Fatal(err)
	}
	defer func(sftpClient *sftp.Client) {
		err := sftpClient.Close()
		if err != nil {
			panic(err.Error())
		}
	}(sftpClient)
	uploadFile(sftpClient, localPath, remotePath)
	elapsed := time.Since(start)
	fmt.Println("elapsed time : ", elapsed)

	// 执行shell
	// kill 进程 kill $(ps -ef | grep go_udp)
	client, err := sshConnect(userName, password, ip, 22)
	if err != nil {
		log.Fatal(err)
	}
	session, _ := client.NewSession()
	defer session.Close()

	// 解压缩文件 cd /data/ && rm -f /data/blog/blog && tar -zxvf blog.tar.gz && rm -rf blog.tar.gz && cd blog && touch blog.log
	session, _ = client.NewSession()
	var str = ""
	str += "cd " + remotePath + " && "
	str += "rm -f " + remotePath + proName + "/" + proName + " && "
	str += "tar -zxvf " + localPath + " && "
	str += "rm -rf " + localPath + " && "
	str += "cd " + proName + " && "
	str += "touch " + proName + ".log"
	var buf []byte
	buf, err = session.CombinedOutput(str)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
	fmt.Println("解压缩文件：" + string(buf))

	// 执行 nohup cd /data/blog && ./run.sh blog
	session, err = client.NewSession()
	if err != nil {
		panic(err)
	}

	str = "cd " + remotePath + proName + " && chmod 777 run.sh && ./run.sh " + proName
	fmt.Println(str)
	buf, err = session.CombinedOutput(str)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("结束：" + string(buf))
}

func sftpConnect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}
	// connect to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func uploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string) {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		log.Fatal(err)
	}
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {
			panic(err.Error())
		}
	}(srcFile)
	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		log.Fatal(err)
	}
	defer func(dstFile *sftp.File) {
		err := dstFile.Close()
		if err != nil {
			panic(err.Error())
		}
	}(dstFile)
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		log.Fatal(err)
	}
	_, err = dstFile.Write(ff)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(localFilePath + " copy file to remote server finished!")
}

func sshConnect(user, password, host string, port int) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	return client, nil
}
