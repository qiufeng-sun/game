工程依赖管理
================

* 使用的工具: gvp + gpm + gpm-git
---------------------------------

* [gvp](https://github.com/pote/gvp)用于设置GOPATH

```bash
$ git clone https://github.com/pote/gvp.git && cd gvp
$ git checkout v0.2.0 # You can ignore this part if you want to install HEAD.
$ ./configure
$ make install
```

* [gpm](https://github.com/pote/gpm)用于下载依赖

```bash
$ git clone https://github.com/pote/gpm.git && cd gpm
$ git checkout v1.3.1 # You can ignore this part if you want to install HEAD.
$ ./configure
$ make install
```

* [gpm-git](https://github.com/technosophos/gpm-git)gpm下载依赖的插件

```sh
$ git clone https://github.com/technosophos/gpm-git.git
$ git checkout v1.0.1
$ make install
```

* 获取工程, 并编译(以goproject/mining为例)
-----------------------------------------

```sh
$ cd ~ #进入work目录
$ mkdir -p goproject/src #创建目录
$ cd goproject # GOPATH目录
$ . gvp # 设置$GOPATH(go env可查看). 当前的为"/home/sxg/goproject/.godeps:/home/sxg/goproject"
$ cd src && git clone xxxx #克隆工程代码
$ cd xxxx # 进入工程目录
$ gpm-git # 读取当前目录下Godeps-Git下载相关依赖
$ go install # 生成可执行文件(会生成在$GOBIN中)
$ cp $GOBIN/xxxx ~/goproject/bin/ # 从$GOBIN拷贝可执行文件到bin目录下
$ cp -r conf/ ~/goproject/bin/ # 拷贝配置文件到bin目录下
```
