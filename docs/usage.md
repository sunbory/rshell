# 使用说明

## 版本说明

- Linux：rshell
- Windows：rshell.exe

## 命令行交互执行模式

### 启动

```
#rshell
 ______     ______     __  __     ______     __         __
/\  == \   /\  ___\   /\ \_\ \   /\  ___\   /\ \       /\ \
\ \  __<   \ \___  \  \ \  __ \  \ \  __\   \ \ \____  \ \ \____
 \ \_\ \_\  \/\_____\  \ \_\ \_\  \ \_____\  \ \_____\  \ \_____\
  \/_/ /_/   \/_____/   \/_/\/_/   \/_____/   \/_____/   \/_____/
------ Rshell @7.0 Type "?" or "help" for more information. -----
{The Correct Step: 1.help -> 2.load -> 3.*/sudo/download/upload}
rshell:
```

### 应用场景

#### 获取帮助

```
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
    - encrypt Cloud12#$

decrypt ciphertext_password
    --- Decrypt ciphertext_password with aes 256 cfb

    Examples:
    - decrypt 1c15b86d686758158d5fd9551c0ccca6168a2c80f149c38bca

exit
    --- Exit rshell
?
    --- Help
```

#### 加载环境信息

- 加载或切换主机组

```
rshell: load -Htest-group01
[alpha-env-ttt-key-su@test-group01:22]#load -Htest-group02
[alpha-env-root-pass@test-group02:22]#

```

- 切换至特定IP地址

```
[alpha-env-root-pass@test-group02:22]# load -Aalpha-env-root-pass -H192.168.31.63
[alpha-env-root-pass@192.168.31.63:22]# load -Aalpha-env-root-pass -H192.168.31.63 -P23
[alpha-env-root-pass@192.168.31.63:23]#
```

#### 执行远程Shell命令

- 以普通用户执行单条命令

```
[alpha-env-ttt-key-su@test-group01:22]# whoami;
TASK [do@test-group01                                   ] *********************
HOST [192.168.31.63   ] =======================================================
ttt

```

- 以普通用户执行多条命令

```
[alpha-env-ttt-key-su@test-group01:22]# whoami;date;pwd;
TASK [do@test-group01                                   ] *********************
HOST [192.168.31.63   ] =======================================================
ttt
Thu Jan 24 18:35:14 CST 2019
/home/ttt
```

- 自动切换root用户并执行单条命令

```
[alpha-env-ttt-key-su@test-group01:22]# sudo whoami;
TASK [sudo@test-group01                                 ] *********************
HOST [192.168.31.63   ] =======================================================
root
Last login: Thu Jan 24 18:20:47 CST 2019

STDERR =>
Password: 
```

- 自动切换root用户并执行多条命令

```
[alpha-env-ttt-key-su@test-group01:22]# sudo whoami;pwd;date;
TASK [sudo@test-group01                                 ] *********************
HOST [192.168.31.63   ] =======================================================
root
/root
Thu Jan 24 18:35:42 CST 2019
Last login: Thu Jan 24 18:35:26 CST 2019

STDERR =>
Password: 
```

#### 执行远程上传下载

> 上传下载操作采用SSH通道，不支持自动切换root用户，可以通过执行命令方式实现文件权限及路径调整动作

- 上传文件到用户目录当前目录

> 文件需要已经存在

```
[alpha-env-ttt-key-su@test-group01:22]# upload rshell.exe .
TASK [upload@test-group01                               ] *********************
HOST [192.168.31.63   ] =======================================================
SUCCESS [rshell.exe -> ./] :
rshell.exe
```

- 上传文件到特定目录

> 文件需要已经存在，目录需要已经存在并且具有用户操作权限

```
[alpha-env-ttt-key-su@test-group01:22]# upload rshell.exe /home/ttt
TASK [upload@test-group01                               ] *********************
HOST [192.168.31.63   ] =======================================================
SUCCESS [rshell.exe -> /home/ttt/] :
rshell.exe
```

- 下载文件到当前目录

> 文件需要已经存在

```
[alpha-env-ttt-key-su@test-group01:22]# download /home/ttt/rshell.exe .
TASK [download@test-group01                             ] *********************
HOST [192.168.31.63   ] =======================================================
SUCCESS [/home/ttt/rshell.exe -> ./] :
/home/ttt/rshell.exe
```

- 下载文件到特定目录

> 文件需要已经存在

```
[alpha-env-ttt-key-su@test-group01:22]# download /home/ttt/rshell.exe /tmp
TASK [download@test-group01                             ] *********************
HOST [192.168.31.63   ] =======================================================
SUCCESS [/home/ttt/rshell.exe -> /tmp/] :
/home/ttt/rshell.exe
```

#### 加解密密码

> 需要在cfg.yaml中配置加密类型和key值

- 加密

```
[alpha-env-ttt-key-su@test-group01:22]# encrypt Huawei@123
766b9f12486f29f925626fb65a60da0fd5f6827859ac102938d7
```

- 解密

```
[alpha-env-ttt-key-su@test-group01:22]# decrypt 766b9f12486f29f925626fb65a60da0fd5f6827859ac102938d7
Huawei@123
```

#### 退出

```
rshell: exit

```

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

## 脚本任务编排执行模式

> 如果脚本含有变量，需要提前准备包含变量配置的yaml

- 执行无变量脚本

```
# rshell.exe -f examples\test.yaml
TASK [test all/test do@test-group01                     ] *********************
HOST [192.168.31.37   ] =======================================================
Last failed login: Mon Dec  3 19:05:33 CST 2018 from 192.168.1.22 on ssh:notty
There were 2 failed login attempts since the last successful login.
ttt
Mon Dec  3 21:15:46 CST 2018
{{.group1.abc}}
{{.group2.abc}}

HOST [192.168.31.63   ] =======================================================
Last failed login: Mon Dec  3 19:02:29 CST 2018 from 192.168.1.22 on ssh:notty
There were 2 failed login attempts since the last successful login.
ttt
Mon Dec  3 21:12:43 CST 2018
{{.group1.abc}}
{{.group2.abc}}

TASK [test all/test sudo@test-group01                   ] *********************
HOST [192.168.31.37   ] =======================================================
root
Mon Dec  3 21:15:47 CST 2018
Last login: Mon Dec  3 21:13:34 CST 2018

STDERR =>
Password:
HOST [192.168.31.63   ] =======================================================
root
Mon Dec  3 21:12:44 CST 2018
Last login: Mon Dec  3 21:10:31 CST 2018

STDERR =>
Password:
TASK [test all/test upload@test-group01                 ] *********************
HOST [192.168.31.37   ] =======================================================
UPLOAD Success [examples/test.yaml -> ./]
HOST [192.168.31.63   ] =======================================================
UPLOAD Success [examples/test.yaml -> ./]
TASK [test all/check upload@test-group01                ] *********************
HOST [192.168.31.37   ] =======================================================
-rw-rw-r--. 1 ttt ttt 542 Dec  3 21:15 test.yaml

HOST [192.168.31.63   ] =======================================================
-rw-rw-r--. 1 ttt ttt 542 Dec  3 21:12 test.yaml

TASK [test all/test download@test-group01               ] *********************
HOST [192.168.31.37   ] =======================================================
DOWNLOAD Success [test.yaml -> test-group01]
HOST [192.168.31.63   ] =======================================================
DOWNLOAD Success [test.yaml -> test-group01]

```

- 执行有变量脚本

```
# rshell.exe -f examples\test.yaml -v examples\values.yaml
TASK [test all/test do@test-group01                     ] *********************
HOST [192.168.31.63   ] =======================================================
Last failed login: Mon Dec  3 19:02:29 CST 2018 from 192.168.1.22 on ssh:notty
There were 2 failed login attempts since the last successful login.
ttt
Mon Dec  3 21:14:06 CST 2018
123
efg

HOST [192.168.31.37   ] =======================================================
Last failed login: Mon Dec  3 19:05:33 CST 2018 from 192.168.1.22 on ssh:notty
There were 2 failed login attempts since the last successful login.
ttt
Mon Dec  3 21:17:09 CST 2018
123
efg

TASK [test all/test sudo@test-group01                   ] *********************
HOST [192.168.31.37   ] =======================================================
root
Mon Dec  3 21:17:10 CST 2018
Last login: Mon Dec  3 21:15:47 CST 2018

STDERR =>
Password:
HOST [192.168.31.63   ] =======================================================
root
Mon Dec  3 21:14:07 CST 2018
Last login: Mon Dec  3 21:12:44 CST 2018

STDERR =>
Password:
TASK [test all/test upload@test-group01                 ] *********************
HOST [192.168.31.37   ] =======================================================
UPLOAD Success [examples/test.yaml -> ./]
HOST [192.168.31.63   ] =======================================================
UPLOAD Success [examples/test.yaml -> ./]
TASK [test all/check upload@test-group01                ] *********************
HOST [192.168.31.37   ] =======================================================
-rw-rw-r--. 1 ttt ttt 542 Dec  3 21:17 test.yaml

HOST [192.168.31.63   ] =======================================================
-rw-rw-r--. 1 ttt ttt 542 Dec  3 21:14 test.yaml

TASK [test all/test download@test-group01               ] *********************
HOST [192.168.31.37   ] =======================================================
DOWNLOAD Success [test.yaml -> test-group01]
HOST [192.168.31.63   ] =======================================================
DOWNLOAD Success [test.yaml -> test-group01]

```
