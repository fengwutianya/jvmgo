package main

import "fmt"
import "strings"
import "jvmgo/ch09/classpath"
import "jvmgo/ch09/rtda/heap"

func main() {
	cmd := parseCmd()

	if cmd.versionFlag {
		fmt.Println("version 0.0.1")
	} else if cmd.helpFlag || cmd.class == "" {
		printUsage()
	} else {
		startJVM(cmd)
	}
}

func startJVM(cmd *Cmd) {
	cp := classpath.Parse(cmd.XjreOption, cmd.cpOption)				//路径模块 解析jre/lib路径 classpath路径
	classLoader := heap.NewClassLoader(cp, cmd.verboseClassFlag)	//产生类加载器模块,已经加载了primitive类们

	className := strings.Replace(cmd.class, ".", "/", -1)
	mainClass := classLoader.LoadClass(className)					//加载主类，获得class对象
	mainMethod := mainClass.GetMainMethod()							//从主类得到main方法，开始执行字节码
	if mainMethod != nil {
		interpret(mainMethod, cmd.verboseInstFlag, cmd.args)
	} else {
		fmt.Printf("Main method not found in class %s\n", cmd.class)
	}
}
