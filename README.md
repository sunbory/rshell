# rshell

## 功能说明

多Linux主机远程批量执行Shell命令和上传下载文件（跨平台，无依赖，免安装）

- 简单化，单文件运行，无外部依赖
- 跨平台，运行支持Win和Linux平台
- 三模式，支持文件编排、命令行交互和单命令行操作
- 双类型，支持ssh命令和ftp上传下载文件
- 双认证，支持密码和key认证
- 自切换，支持自动切换root用户
- 高安全，支持高危命令黑名单，密码支持加密
- 智能化，支持自动提示补全，历史搜索
- 定制化，支持提示符、分隔符、超时等定制
- 模板化，文件编排支持变量自定义
- 多样化，支持text、json、yaml格式输出

## 应用安装

```
go get github.com/luckywinds/rshell
```

## 应用构建

```
go build rshell.go
```

- rshell：    Linux版本
- rshell.exe：Windows版本

## 目录

> Linux和Windows用法相同，注意基本的命令和路径分割符区别即可

- [配置说明](#配置说明)
  - [配置文件路径](#配置文件路径)
  - [系统配置](#系统配置)
  - [认证信息配置](#认证信息配置)
  - [主机信息配置](#主机信息配置)
  - [任务编排配置](#任务编排配置)
  - [任务变量配置](#任务变量配置)
- [使用说明](#使用说明)
  - [交互式命令行模式](#交互式命令行模式)
    - [加载或切换当前环境信息](#加载或切换当前环境信息)
    - [以正常登陆用户执行shell命令](#以正常登陆用户执行shell命令)
    - [以正常登陆用户并自动切换root用户执行shell命令](#以正常登陆用户并自动切换root用户执行shell命令)
    - [以正常登陆用户执行上传下载动作](#以正常登陆用户执行上传下载动作)
    - [加解密](#加解密)
    - [快捷键](#快捷键)
  - [脚本任务编排模式](#脚本任务编排模式)
  - [普通命令行模式](#普通命令行模式)
  - [输出格式自定义](#输出格式自定义)
- [FAQ](#faq)

## 配置说明

### 配置文件路径

> 当前目录下.rshell文件夹

> 配置文件中的路径分隔符统一采用正斜杠/

### 系统配置

> .rshell/cfg.yaml

```
concurrency: 10                                         #The num of concurrency goroutine for hosts, Default: 10
tasktimeout: 300                                        #The total timeout [second] of everyone goroutine, Default: 300
connecttimeout: 3600                                    #The total timeout [second] of everyone ssh connection, Default: 3600
cmdseparator: ";"                                       #The command separator, Default: ";"
promptstring: "rshell: "                                #The system prompt string, Default: "rshell: "
hostgroupsize: 200                                      #The max number of ip in hostgroup, Default: 200
passcrypttype: "aes"                                    #The password crypt type, support [aes], Default: ""
passcryptkey: "i`a!M@a#H$a%P^p&Y*c(R)y_P+t,K.eY"        #The aes crypt key, length must be 32. Not be empty if passcrypttype = aes.
outputtype: "yaml"                                      #The result output format, must be in [text, json, yaml]. Default: "text"
blackcmdlist:                                           #The dangerous black command list, Default if =cmd or ~=cmdprefix can not run.
- cmd: rm -rf /
  cmdprefix: rm -rf /root
mostusedcmds:                                           #The most used commands
- pwd
- date
- whoami
updateserver:                                           #The auto update server address
- https://github.com/luckywinds/rshell/raw/master/releases
```

### 认证信息配置

> .rshell/auth.yaml

```
authmethods:
- name: default_user_key
  username: ttt
  privatekey: id_rsa
  passphrase: 81fed25324dfe10da02dca3b4a7ba098e9698a8feb510d0138

- name: default_user_key_root
  username: ttt
  privatekey: id_rsa
  passphrase: 81fed25324dfe10da02dca3b4a7ba098e9698a8feb510d0138
  sudotype: "su -"
  sudopass: df61be1184a01047cde6e3ef2d5d2a26b2c20b4d3b8094733e
  
- name: default_user_pass
  username: ttt
  password: 96b12173e0f2a5d5f52257bcc2e651b0cf00ca09d894465a9db4

- name: default_user_pass_root
  username: ttt
  password: 96b12173e0f2a5d5f52257bcc2e651b0cf00ca09d894465a9db4
  sudotype: "su -"
  sudopass: df61be1184a01047cde6e3ef2d5d2a26b2c20b4d3b8094733e
  
- name: default_root_key
  username: root
  privatekey: id_rsa
  passphrase: 81fed25324dfe10da02dca3b4a7ba098e9698a8feb510d0138

- name: default_root_pass
  username: root
  password: df61be1184a01047cde6e3ef2d5d2a26b2c20b4d3b8094733e

```

> 样例的密码都是随机的，自行根据需求更改即可，通过encrypt和decrypt进行加解密

> 如果想直接填写明文不需要加密的话，系统配置项passcrypttype留空即可

### 主机信息配置

> .rshell/hosts.yaml

```
hostgroups:
- groupname: dev-centos_user_key
  authmethod: default_user_key
  hosts:
  - 192.168.31.63
  - 192.168.31.37

- groupname: dev-centos_user_key_root
  authmethod: default_user_key_root
  hosts:
  - 192.168.31.63
  - 192.168.31.37

- groupname: dev-centos_user_pass
  authmethod: default_user_pass
  hosts:
  - 192.168.31.63
  - 192.168.31.37

- groupname: dev-centos_user_pass_root
  authmethod: default_user_pass_root
  hosts:
  - 192.168.31.63
  - 192.168.31.37

- groupname: dev-centos_root_key
  authmethod: default_root_key
  hosts:
  - 192.168.31.63
  - 192.168.31.37

- groupname: dev-centos_root_pass
  authmethod: default_root_pass
  hosts:
  - 192.168.31.63
  - 192.168.31.37

- groupname: dev-suse_root_pass
  authmethod: default_root_pass
  hosts:
  - 192.168.1.126

- groupname: dev-centos-suse_root_pass
  authmethod: default_root_pass
  groups:
  - dev-centos_root_pass
  - dev-suse_root_pass
  hosts:
  - 192.168.31.197
  hostranges:
  - from: 192.168.1.126
    to: 192.168.1.127
```

> 主机组支持组合模式，groups/hosts/hostranges至少选1项即可

### 任务编排配置

> 示例：example\test.yaml

```
tasks:
- name: test ssh
  hostgroup: dev-centos_root_pass
  subtasks:
  - name: test do
    mode: ssh
    cmds:
    - whoami
    - date
    - echo {{.group1.abc}}
    - echo {{.group2.abc}}

  - name: test sudo
    mode: ssh
    sudo: true
    cmds:
    - whoami
    - date

- name: test sftp
  hostgroup: dev-centos_root_pass
  subtasks:
  - name: test upload
    mode: sftp
    ftptype: upload
    srcfile: examples/test.yaml
    desdir: .

  - name: check upload
    mode: ssh
    cmds:
    - ls -alh test.yaml

  - name: test download
    mode: sftp
    ftptype: download
    srcfile: test.yaml
    desdir: . 

```

> 变量替换语法为[golang template](https://golang.org/pkg/text/template/)

### 任务变量配置

> 示例：example\values.yaml

```
group1:               
  abc: 123
group2:
  abc: efg
```

## 使用说明

> 首字符串load/sudo/download/upload/encrypt/decrypt/exit将被识别为rshell关键字，操作中如有冲突请发散思考。

## 交互式命令行模式

### 语法

```
# rshell.exe
rshell: help
load -A<auth> -H<host> -P<port>
    --- Load current env(auth@host:port)

    Tips:
    - host，must in hosts.yaml or ipv4
    - auth, optional if config with host in hosts.yaml, must in auth.yaml
    - port, optional if config with host in hosts.yaml, default: 22

    Examples:
    - load -Hhostgroup03
    - load -Aalpha-env-root-pass -H192.168.31.63
    - load -Aalpha-env-root-pass -H192.168.31.63 -P23

cmd1;cmd2;cmd3
    --- Run cmds on TARGETS as normal user

    Examples:
    - pwd
    - pwd;whoami;date

sudo cmd1;cmd2;cmd3
    --- Run cmds on TARGETS as root which auto change from normal user

    Examples:
    - sudo pwd
    - sudo pwd;whoami;date

download srcFile desDir
    --- Download srcFile from TARGETS to LOCAL desDir as normal user

    Examples:
    - download .bashrc .

upload srcFile desDir
    --- Upload srcFile from LOCAL to TARGETS desDir as normal user

    Examples:
    - upload README.md .

encrypt cleartext_password
    --- Encrypt cleartext_password with aes 256 cfb

    Examples:
    - encrypt rshell@123

decrypt ciphertext_password
    --- Decrypt ciphertext_password with aes 256 cfb

    Examples:
    - decrypt 20d59477070d93439a5390e9f425ec4afe1c781775727c35817f

exit
    --- Exit rshell
?
    --- Help
```

### 样例

#### 加载或切换当前环境信息

```
rshell: load -Hdev-centos_user_key
[default_user_key@dev-centos_user_key:22]# pwd;date;whoami;
TASK  [do                  ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.63       ] --------------------------------------------------
/home/ttt
Thu Jan 31 16:27:03 CST 2019
ttt

HOST  [192.168.31.37       ] --------------------------------------------------
Last failed login: Thu Jan 24 14:17:09 CST 2019 from 192.168.1.22 on ssh:notty
There were 10 failed login attempts since the last successful login.
/home/ttt
Thu Jan 31 16:29:41 CST 2019
ttt

[default_user_key@dev-centos_user_key:22]# load -Hdev-centos_root_pass
[default_root_pass@dev-centos_root_pass:22]# pwd;date;whoami;
TASK  [do                  ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.63       ] --------------------------------------------------
/root
Thu Jan 31 16:27:05 CST 2019
root

HOST  [192.168.31.37       ] --------------------------------------------------
/root
Thu Jan 31 16:29:46 CST 2019
root

[default_root_pass@dev-centos_root_pass:22]# load -Adefault_root_key -H192.168.31.197
[default_root_key@192.168.31.197:22]# pwd;date;whoami;
TASK  [do                  ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.197      ] --------------------------------------------------
/root
Thu Jan 31 16:33:39 CST 2019
root

[default_root_key@192.168.31.197:22]# load -Hdev-centos-suse_root_pass
[default_root_pass@dev-centos-suse_root_pass:22]# pwd;date;whoami;
TASK  [do                  ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.197      ] --------------------------------------------------
/root
Thu Jan 31 16:33:40 CST 2019
root

HOST  [192.168.31.63       ] --------------------------------------------------
/root
Thu Jan 31 16:27:09 CST 2019
root

HOST  [192.168.31.37       ] --------------------------------------------------
/root
Thu Jan 31 16:29:50 CST 2019
root

HOST  [192.168.1.126       ] --------------------------------------------------
/root
Thu Jan 31 16:25:10 CST 2019
root

```

> 默认下次打开，会自动加载上次退出前记录的环境信息

#### 以正常登陆用户执行shell命令

```
[default_user_key_root@dev-centos_user_key_root:22]# date;pwd;whoami;
TASK  [do                  ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.63       ] --------------------------------------------------
Thu Jan 31 16:30:07 CST 2019
/home/ttt
ttt

HOST  [192.168.31.37       ] --------------------------------------------------
Thu Jan 31 16:32:48 CST 2019
/home/ttt
ttt

```

#### 以正常登陆用户并自动切换root用户执行shell命令

```
[default_user_key_root@dev-centos_user_key_root:22]# sudo date;pwd;whoami;
TASK  [sudo                ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.63       ] --------------------------------------------------
Thu Jan 31 16:31:41 CST 2019
/root
root
Last login: Thu Jan 31 16:21:22 CST 2019

STDERR =>
Password:
HOST  [192.168.31.37       ] --------------------------------------------------
Thu Jan 31 16:34:22 CST 2019
/root
root
Last login: Thu Jan 31 16:24:02 CST 2019

STDERR =>
Password:
```

#### 以正常登陆用户执行上传下载动作

```
[default_user_key_root@dev-centos_user_key_root:22]# upload README.md /tmp
TASK  [upload              ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.63       ] --------------------------------------------------
SUCCESS [README.md -> /tmp/] :
README.md
HOST  [192.168.31.37       ] --------------------------------------------------
SUCCESS [README.md -> /tmp/] :
README.md
[default_user_key_root@dev-centos_user_key_root:22]# ls /tmp/*.md
TASK  [do                  ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.63       ] --------------------------------------------------
/tmp/README.md

HOST  [192.168.31.37       ] --------------------------------------------------
/tmp/README.md

[default_user_key_root@dev-centos_user_key_root:22]# download /tmp/README.md /home
TASK  [download            ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.63       ] --------------------------------------------------
SUCCESS [/tmp/README.md -> /home/] :
/tmp/README.md
HOST  [192.168.31.37       ] --------------------------------------------------
SUCCESS [/tmp/README.md -> /home/] :
/tmp/README.md
```

> 上传下载不支持自动切换root用户，如有需要可以配合执行shell命令满足诉求

#### 加解密

```
[default_user_key_root@dev-centos_user_key_root:22]# encrypt rshell@123
20d59477070d93439a5390e9f425ec4afe1c781775727c35817f
[default_user_key_root@dev-centos_user_key_root:22]# decrypt 20d59477070d93439a5390e9f425ec4afe1c781775727c35817f
rshell@123
```

> 加解密不要求提前加载环境信息

#### 快捷键

Keystroke    | Action
---------    | ------
Ctrl-A, Home | Move cursor to beginning of line
Ctrl-E, End  | Move cursor to end of line
Ctrl-B, Left | Move cursor one character left
Ctrl-F, Right| Move cursor one character right
Ctrl-Left, Alt-B    | Move cursor to previous word
Ctrl-Right, Alt-F   | Move cursor to next word
Ctrl-D, Del  | (if line is *not* empty) Delete character under cursor
Ctrl-D       | (if line *is* empty) End of File - usually quits application
Ctrl-C       | Reset input (create new empty prompt)
Ctrl-L       | Clear screen (line is unmodified)
Ctrl-T       | Transpose previous character with current character
Ctrl-H, BackSpace | Delete character before cursor
Ctrl-W       | Delete word leading up to cursor
Ctrl-K       | Delete from cursor to end of line
Ctrl-U       | Delete from start of line to cursor
Ctrl-P, Up   | Previous match from history
Ctrl-N, Down | Next match from history
Ctrl-R       | Reverse Search history (Ctrl-G cancel)
Ctrl-Y       | Paste from Yank buffer
Tab          | Next completion

> Ctrl-R，可以通过关键字快速从历史命令记录中快速搜索想要的命令

## 脚本任务编排模式

### 语法

```
# rshell.exe -h
Usage: rshell.exe [<options>]

Options:
  -f string
        The script yaml.
  -v string
        The script values yaml.
```

> -v 可选，当不需要变量时

### 样例

```
# rshell.exe -f examples\test.yaml -v examples\values.yaml
TASK  [test ssh            ] ++++++++++++++++++++++++++++++++++++++++++++++++++
STASK [test do             ] ==================================================
HOST  [192.168.31.63       ] --------------------------------------------------
root
Thu Jan 31 16:21:21 CST 2019
123
efg

HOST  [192.168.31.37       ] --------------------------------------------------
root
Thu Jan 31 16:24:02 CST 2019
123
efg

STASK [test sudo           ] ==================================================
HOST  [192.168.31.63       ] --------------------------------------------------
root
Thu Jan 31 16:21:22 CST 2019

HOST  [192.168.31.37       ] --------------------------------------------------
root
Thu Jan 31 16:24:02 CST 2019

TASK  [test sftp           ] ++++++++++++++++++++++++++++++++++++++++++++++++++
STASK [test upload         ] ==================================================
HOST  [192.168.31.63       ] --------------------------------------------------
SUCCESS [examples/test.yaml -> ./] :
examples/test.yaml
HOST  [192.168.31.37       ] --------------------------------------------------
SUCCESS [examples/test.yaml -> ./] :
examples/test.yaml
STASK [check upload        ] ==================================================
HOST  [192.168.31.63       ] --------------------------------------------------
-rw-r--r--. 1 root root 614 Jan 31 16:21 test.yaml

HOST  [192.168.31.37       ] --------------------------------------------------
-rw-r--r--. 1 root root 614 Jan 31 16:24 test.yaml

STASK [test download       ] ==================================================
HOST  [192.168.31.63       ] --------------------------------------------------
SUCCESS [test.yaml -> ./] :
test.yaml
HOST  [192.168.31.37       ] --------------------------------------------------
SUCCESS [test.yaml -> ./] :
test.yaml

```

## 普通命令行模式

### 语法

```
# rshell.exe -h
Usage: rshell.exe [<options>]

Options:
  -A string
        The auth method name.
  -H string
        The host name.
  -L string
        The command line.
  -P int
        The ssh port. (default 22)
```

> -H 必选，可以指定主机组名字或IPv4地址

> -A 可选，当-H指定IPv4地址时必选

> -P 可选，默认22

> -L 支持和交互式命令行模式同样的命令

### 样例

```
# rshell.exe -H dev-centos_user_key_root -L "pwd;date;whoami;"
TASK  [do                  ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.63       ] --------------------------------------------------
/home/ttt
Thu Jan 31 16:53:55 CST 2019
ttt

HOST  [192.168.31.37       ] --------------------------------------------------
Last failed login: Thu Jan 24 14:17:09 CST 2019 from 192.168.1.22 on ssh:notty
There were 10 failed login attempts since the last successful login.
/home/ttt
Thu Jan 31 16:56:36 CST 2019
ttt


# rshell.exe -H dev-centos_user_key_root -L "sudo pwd;date;whoami;"
TASK  [sudo                ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.63       ] --------------------------------------------------
/root
Thu Jan 31 16:53:57 CST 2019
root
Last login: Thu Jan 31 16:31:41 CST 2019

STDERR =>
Password:

# rshell.exe -A default_root_key -H 192.168.31.197 -L "upload README.md /tmp"
TASK  [upload              ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.197      ] --------------------------------------------------
SUCCESS [README.md -> /tmp/] :
README.md

# rshell.exe -A default_root_key -H 192.168.31.197 -L "ls /tmp/*.md"
TASK  [do                  ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.197      ] --------------------------------------------------
/tmp/README.md


# rshell.exe -A default_root_key -H 192.168.31.197 -L "download /tmp/README.md /home"
TASK  [download            ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.197      ] --------------------------------------------------
SUCCESS [/tmp/README.md -> /home/] :
/tmp/README.md

```

## 输出格式自定义

> 配置文件.rshell/cfg.yaml，配置项：outputtype

> 支持text、json、yaml三种格式

> 支持交互式命令行模式、脚本任务编排模式、命令行模式

### 样例

#### text格式

```
# rshell.exe -H dev-centos_user_key_root -L "pwd;date;whoami"
TASK  [do                  ] ++++++++++++++++++++++++++++++++++++++++++++++++++
HOST  [192.168.31.63       ] --------------------------------------------------
/home/ttt
Thu Jan 31 17:02:06 CST 2019
ttt

HOST  [192.168.31.37       ] --------------------------------------------------
Last failed login: Thu Jan 24 14:17:09 CST 2019 from 192.168.1.22 on ssh:notty
There were 10 failed login attempts since the last successful login.
/home/ttt
Thu Jan 31 17:04:47 CST 2019
ttt
```

### json格式

```
# rshell.exe -H dev-centos_user_key_root -L "pwd;date;whoami"
[
  {
    "Actiontype": "do",
    "Groupname": "dev-centos_user_key_root",
    "Hostaddr": "192.168.31.63",
    "Error": "",
    "Stdout": "/home/ttt\nThu Jan 31 17:02:29 CST 2019\nttt\n",
    "Stderr": ""
  },
  {
    "Actiontype": "do",
    "Groupname": "dev-centos_user_key_root",
    "Hostaddr": "192.168.31.37",
    "Error": "",
    "Stdout": "Last failed login: Thu Jan 24 14:17:09 CST 2019 from 192.168.1.22 on ssh:notty\nThere were 10 failed login attempts since the last successful login.\n/home/ttt\nThu Jan 31 17:05:09 CST 2019\nttt\n",
    "Stderr": ""
  }
]
```

### yaml格式

```
# rshell.exe -H dev-centos_user_key_root -L "pwd;date;whoami"
- actiontype: do
  groupname: dev-centos_user_key_root
  hostaddr: 192.168.31.63
  stdout: |
    /home/ttt
    Thu Jan 31 17:03:10 CST 2019
    ttt
- actiontype: do
  groupname: dev-centos_user_key_root
  hostaddr: 192.168.31.37
  stdout: |
    Last failed login: Thu Jan 24 14:17:09 CST 2019 from 192.168.1.22 on ssh:notty
    There were 10 failed login attempts since the last successful login.
    /home/ttt
    Thu Jan 31 17:05:51 CST 2019
    ttt

```

## FAQ

- 执行sudo相关命令时，STDERR报sudo: no tty present and no askpass program specified

可能原因：自动切换用户时默认需要经过tty获取密码

建议措施：增加-S选项，sudotype: "sudo -S su -"
