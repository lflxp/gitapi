package gitapi

import (
	"os"
	"runtime"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Git struct {
	Url		string
	Version 	string
	Env 		[]string
	Os 		string
	Arch 		string
	GoVersion 	string
	CmdString	string
	lock 		sync.Mutex
	Args 		string
	Kargs 		[]string
	Debug 		bool
}
//初始化 检测git命令安装和url是否为git项目
func (this *Git) Init() error {
	if this.Url == "" {
		if this.Debug {
			Log.Error("Url is empty")
		}
		return errors.New("Url is empty")
	}
	version,err := ExecCommand("git version")
	if err != nil {
		if this.Debug {
			Log.Error("GIT IS NOT INSTALL")
		}
		return errors.New("GIT IS NOT INSTALL")
	}

	this.Version = version
	this.Env = os.Environ()
	this.Os = runtime.GOOS
	this.Arch = runtime.GOARCH
	this.GoVersion = runtime.Version()
	this.CmdString = "git -C "+this.Url+" "
	_,err = ExecCommand(this.CmdString+" status")
	if err != nil {
		if this.Debug {
			Log.Error(fmt.Sprintf("%s IS NOT a git repository (or any of the parent directories): .git",this.Url))
		}
		return errors.New(fmt.Sprintf("%s IS NOT a git repository (or any of the parent directories): .git",this.Url))
	}
	return nil
}
//查看状态
//-v, --verbose         be verbose
//-s, --short           show status concisely
//-b, --branch          show branch information
//--porcelain           machine-readable output
//--long                show status in long format (default)
//-z, --null            terminate entries with NUL
//-u, --untracked-files[=<mode>]
//		  show untracked files, optional modes: all, normal, no. (Default: all)
//--ignored             show ignored files
//--ignore-submodules[=<when>]
//		  ignore changes to submodules, optional when: all, dirty, untracked. (Default: all)
//--column[=<style>]    list untracked files in columns
//返回未被管理的文件列表 git status -s |grep '??'
func (this *Git) Status(args ...string) (string,error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.CmdString = "git -C "+this.Url+" "

	var cmd string
	if len(args) != 0 {
		cmd = this.CmdString+"status "+strings.Join(args," ")
	} else {
		cmd = this.CmdString+" status"
	}

	rs,err := ExecCommand(cmd)
	if this.Debug {
		Log.Debug(cmd)
		if err != nil {
			Log.Error(err.Error())
		}
	}
	return rs,err
}
//创建裸库 init
func (this *Git) Bare(path string) (string,error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	rs,err := ExecCommand("git init "+path)
	if this.Debug {
		Log.Debug("git init "+path)
		if err != nil {
			Log.Error(err.Error())
		}
	}
	return rs,err
}
//添加操作 git add .
func (this *Git) Add(args ...string) (string,error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.CmdString = "git -C "+this.Url+" "

	var cmd string
	if len(args) == 0 {
		cmd = this.CmdString+" add ."
	} else {
		cmd = this.CmdString+" add "+strings.Join(args," ")
	}
	rs,err := ExecCommand(cmd)
	if this.Debug {
		Log.Debug(cmd)
		if err != nil {
			Log.Error(err.Error())
		}
	}
	return rs,err
}
//添加本地库操作 git commit -m "123"
func (this *Git) Commit(common string) (string,error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.CmdString = "git -C "+this.Url+" "

	rs,err := ExecCommand(this.CmdString+" commit -m \""+common+"\"")
	if this.Debug {
		Log.Debug(this.CmdString+" commit -m \""+common+"\"")
		if err != nil {
			Log.Error(err.Error())
		}
	}
	return rs,err
}
//查看当前分支
func (this *Git) Branch() (string,error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.CmdString = "git -C "+this.Url+" "

	rs,err := ExecCommand(this.CmdString+" branch")
	if this.Debug {
		Log.Debug(this.CmdString+" branch")
		if err != nil {
			Log.Error(err.Error())
		}
	}
	return rs,err
}
//查看git log --grep=ok
//http://www.cnblogs.com/gbyukg/archive/2011/12/12/2285419.html
//git log --pretty=oneline -g   操作记录 获取文件列表
//git log --pretty=format:"%h|%an|%cn|%ce|%cd|%cr|%s"  提交历史 获取commit记录 tag
//git log --pretty=format:"%h|%an|%ae|%ar|%cn|%ce|%cr|%s"
//git log --pretty=format:"%h|%an|%ae|%ar|%cn|%ce|%cr|%s" --graph
//git log --pretty=oneline --graph --stat
func (this *Git) Log(args ...string) (string,error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.CmdString = "git -C "+this.Url+" "

	var cmd string
	if len(args) == 0 {
		cmd = this.CmdString+" log"
	} else {
		cmd = this.CmdString+" log "+strings.Join(args," ")
	}

	rs,err := ExecCommand(cmd)
	if this.Debug {
		Log.Debug(cmd)
		if err != nil {
			Log.Error(err.Error())
		}
	}
	return rs,err
}
//切换分支 没有分支就直接创建
func (this *Git) CheckOut(branch string,args ...string) (string,error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.CmdString = "git -C "+this.Url+" "

	var cmd string
	if len(args) == 1 && args[0] == "-b" {
		cmd = this.CmdString+"checkout -b "+branch
	} else {
		cmd = this.CmdString+"checkout "+branch
	}
	rs,err := ExecCommand(cmd)
	if this.Debug {
		Log.Debug(cmd)
		if err != nil {
			Log.Error(err.Error())
		}
	}
	return rs,err
}
//pull操作 origin name is unchanged
func (this *Git) Pull(branch ...string) (string,error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.CmdString = "git -C "+this.Url+" "

	var cmd string
	if len(branch) == 0 {
		cmd = this.CmdString+"pull origin master"
	} else {
		cmd = this.CmdString+"pull origin "+branch[0]
	}
	rs,err := ExecCommand(cmd)
	if this.Debug {
		Log.Debug(cmd)
		if err != nil {
			Log.Error(err.Error())
		}
	}
	return rs,err
}
//push操作 origin name is unchanged
func (this *Git) Push(branch ...string) (string,error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.CmdString = "git -C "+this.Url+" "

	var cmd string
	if len(branch) == 0 {
		cmd = this.CmdString+"push origin master"
	} else {
		cmd = this.CmdString+"push origin "+branch[0]
	}
	rs,err := ExecCommand(cmd)
	if this.Debug {
		Log.Debug(cmd)
		if err != nil {
			Log.Error(err.Error())
		}
	}
	return rs,err
}
//获取活动分支、未被管理的文件和判断是否有变更
//false 干净的
//true 有新增
func (this *Git) Is_dirty() (bool,error) {
	data,err := this.Status("-s|wc -l")
	rs := strings.Split(data,"\n")
	if rs[0] == "0" {
		if this.Debug {
			Log.Info(fmt.Sprintf("Is_dirty is %s : %s","false",err.Error()))
		}
		return false,err
	} else {
		if this.Debug {
			if this.Debug {
				Log.Info(fmt.Sprintf("Is_dirty is %s : %s","true",err.Error()))
			}
		}
		return true,err
	}
}
//Clone 克隆和初始化一个新的仓库
func (this *Git) Clone(path string,args ...string) (string,error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	var cmd string
	if len(args) == 0 {
		cmd = "git clone "+path
	} else {
		cmd = "git clone "+path+" "+strings.Join(args," ")
	}

	rs,err := ExecCommand(cmd)
	if this.Debug {
		Log.Debug(cmd)
		if err != nil {
			Log.Error(err.Error())
		}
	}
	return rs,err
}
//回退 reset
func (this *Git) Reset(tags string,args ...string) (string,error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.CmdString = "git -C "+this.Url+" "

	var cmd string
	if len(args) == 0 {
		cmd = this.CmdString+" reset "+tags
	} else {
		cmd = this.CmdString+" reset "+tags+" "+strings.Join(args," ")
	}

	rs,err := ExecCommand(cmd)
	if this.Debug {
		Log.Debug(cmd)
		if err != nil {
			Log.Error(err.Error())
		}
	}
	return rs,err
}
//查看tag show detail
func (this *Git) Show(tags string) (string,error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.CmdString = "git -C "+this.Url+" "

	rs,err := ExecCommand(this.CmdString+"show "+tags)

	if this.Debug {
		Log.Debug(this.CmdString+"show ")
		if err != nil {
			Log.Error(err.Error())
		}
	}
	return rs,err
}
//执行shell命令
func (this *Git) UnsafeCmd(cmd string) (string,error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	rs,err := ExecCommand(cmd)
	if this.Debug {
		Log.Debug(cmd)
		if err != nil {
			Log.Error(err.Error())
		}
	}
	return rs,err
}
//# 启动log日志查看
//    def exposed_checklog(self,server_ip,client_port,command):
//        #判断文件是否存在
//        if not os.path.exists(r'/tmp/gotty'):
//            urls = server_ip+'/static/gotty'
//            cmd = 'wget -P /tmp '+urls+';chmod +x /tmp/gotty'
//            os.system(cmd)
//        # 启动gotty docker logs -f odoo
//        cmd1 = 'chmod +x /tmp/gotty && export TERM=xterm && /tmp/gotty -c cst:online -A gotty.log --title-format "{{ .RemoteAddr }} {{ .Hostname }}" -p "' + client_port + '" --once -w '+command+'> /tmp/gotty.log &'
//        status = os.system(cmd1)
//        return status
func (this *Git) Gotty(server_url string,Kargs ...string) (string,error) {
	port,command := Kargs[0],Kargs[1]
	if IsExistFile("/tmp/gotty") {
		rs,err := this.UnsafeCmd(`chmod +x /tmp/gotty && export TERM=xterm && /tmp/gotty -c cst:online -A gotty.log --title-format "{{ .RemoteAddr }} {{ .Hostname }}" -p "`+port+`" --once -w `+command+` >/tmp/gotty.log &`)
		return rs,err
	} else {
		err := Download("http://"+server_url+"/static/software/gotty","/tmp/gotty")
		if err != nil {
			return "http://"+server_url+"/static/software/gotty is bad",err
		}
		rs,err := this.UnsafeCmd(`chmod +x /tmp/gotty && export TERM=xterm && /tmp/gotty -c cst:online -A gotty.log --title-format "{{ .RemoteAddr }} {{ .Hostname }}" -p "`+port+`" --once -w `+command+` >/tmp/gotty.log &`)
		return rs,err
	}
}