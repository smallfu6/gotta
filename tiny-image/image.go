package main

/*
	Docker实现了容器自身的重用和标准化传播, 使得开发、交付、运维流水线上的
	各个角色真正围绕同一交付物, "test what you write,ship what you test(写什
	么就测什么, 测什么就交付什么)"成为现实;
	Docker镜像构建走到今天, 追求又快又小的镜像已成为云原生开发者的共识;
	结合最新容器构建技术为应用程序其构建出极小的镜像, 使其在云原生生态系统
	中能发挥出更大的优势, 得以更为广泛地应用;

	小镜像打包快、下载快、启动快、占用资源小以及受攻击面小, 对于任何语言的
	开发者来说, 都是十分具有吸引力的; 本节主要讨论如何构建最小的go程序容器
	镜像;

	2008年, Linux Container(LXC)功能被合并到Linux内核中; LXC是一种内核级虚拟
	化技术, 主要基于Namespaces和Cgroups技术, 实现共享一个操作系统内核前提下的
	进程资源隔离, 为进程提供独立的虚拟执行环境, 这样的一个虚拟执行环境就是一
	个容器; 本质上LXC容器与现在Docker所提供的容器是一样的; Docker也是基于
	Namespaces和Cgroups技术实现的,  Docker的创新之处在于其基于Union File
	System技术定义了一套容器打包规范, 真正将容器中的应用及其运行的所有依赖
	都封装到一种特定格式的文件中, 而这种文件就被称为镜像(image);
	镜像是容器的序列化标准, 这一创新为容器的存储、重用和传输奠定了基础, 使得
	镜像可以在全世界快速传播并高速发展;

	使用Dockerfile构建镜像, 代表一种用于编写Dockerfile的领域特定语言, 采用
	Dockerfile方式是镜像构建的标准方法, 其可重复、可自动化、可维护以及分层
	精确控制等特点是传统采用docker commit命令提交镜像所不能比拟的;

	本节主要通过将./httpserver.go源文件编译为httpd程序并通过镜像发布展示如何
	构建最小的go程序镜像;
	- 初始构建(./Dockerfile)
		go build -t repodemo/httpd:latest .
		可以通过 docker inspect imageId 检视构建的镜像的信息
		可以通过 docker history imageId 查看镜像的构建过程详情, 其中大小不为0的
		过程代表属于镜像的一个层

	- builder 模式
		根据初始构建过程中的镜像分层情况, 最终镜像中包含构建环境是多余的,
		只需要在最终镜像中包含足够支撑httpd应用运行的运行环境即可; 而
		base image 就可以满足, 因此可以借助	builder image 构建出应用程序,
		Docker官方为此推出了各种主流编程语言(比如Go、Java、Python及Ruby等)的
		官方基础镜像(base image);
		整个目标镜像的构建分为两个阶段:
		- 构建负责源码编译的构建者镜像 ./Dockerfile.build
			docker build -t repodemo/httpd-builder:latest -f Dockerfile.build .
			构建好的应用程序httpd被放在了镜像repodemo/httpd-builder中的
			/go/src目录下, 需要将 httpd 取出作为下一个阶段的输入
			docker create --name extract-httpserver repodemo/httpd-builder
			docker cp extract-httpserver:/go/src/httpd ./httpd
			docker rm -f extract-httpserver
			docker rmi repodemo/httpd-builder
		- 将第一阶段的输出作为输入, 构建出最终的镜像 ./Dockerfile.target
			docker build -t repodemo/httpd:latest -f Dockerfile.target .

			docker images
				repodemo/httpd latest    d31aa2fcb2f5   5 hours ago   75.4MB
			可以看到目标大小的镜像大小为75.4MB

	- 追求最小镜像
		为了减轻重量将所有不必要的东西都拆掉: 仅保留能支撑我们的应用运行的必要库、
		命令, 其余的一律不纳入目标镜像; 当然这不仅仅是基于尺寸上的考量, 小镜像
		还有额外的好处, 比如:内存占用小, 启动速度快, 更加高效; 不会因其他不必
		要的工具、库的漏洞而被攻击, 减少了攻击面, 更加安全;
		开发者可以挑选合适的基础镜像(base image), 有 busybox 和 alpine, 推荐
		普通开发者选择 alpine, 同时golang base image也提供了alpine版本;
		./Dockerfile.build.alpine, ./Dockerfile.target.alpine
		构建流程同builder模式, 最终构建的目标镜像大小为18.5MB;

	- 多阶段构建
		虽然在之前实现了目标镜像的最小化, 但是构建过程繁琐(需要清理中间产物),
		如果想用一个Dockerfile完成构建, 就需要依赖Docker引擎对多阶段构建
		(multi-stage build)的支持; ./Dockerfile.multistage
		Dockerfile中可以写多个FROM baseimage的语句, 每个FROM语句开启一个构建
		阶段, 并且可以通过as语法为此阶段构建命名(比如builder); 还可以通过
		COPY命令在两个阶段构建产物之间传递数据, 比如传递的httpd程序; 最终
		构建的目标镜像和builder模式构建的镜像在效果上是等价的, 大小也为18.5MB;

*/
