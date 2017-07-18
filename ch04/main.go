package main

import "fmt"
//import "strings"
import (
	"jvmgo/ch04/classpath"
	"jvmgo/ch04/classfile"
	"jvmgo/ch04/rtda"
)

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
	//cp := classpath.Parse(cmd.XjreOption, cmd.cpOption)
	//className := strings.Replace(cmd.class, ".", "/", -1)
	//cf := loadClass(className, cp)
	//fmt.Println(cmd.class)
	//printClassInfo(cf)
	frame := rtda.NewFrame(100, 100)
	testLocalVars(frame.LocalVars())
	testOperandStack(frame.OperandStack())
}
func testOperandStack(stack *rtda.OperandStack) {
	stack.PushInt(100)
	stack.PushInt(-100)
	stack.PushLong(2997924580)
	stack.PushLong(-2997924580)
	stack.PushFloat(3.1415926)
	stack.PushDouble(2.71828182845)
	stack.PushRef(nil)
	println(stack.PopRef())
	println(stack.PopDouble())
	println(stack.PopFloat())
	println(stack.PopLong())
	println(stack.PopLong())
	println(stack.PopInt())
	println(stack.PopInt())
}
func testLocalVars(vars rtda.LocalVars) {
	vars.SetInt(0,100)
	vars.SetInt(1,-100)
	vars.SetLong(2, 100001111)
	vars.SetLong(4, -100001111)
	vars.SetFloat(6, 3.1415926)
	vars.SetDouble(7, 2.71828182845)
	vars.SetRef(9, nil)
	println(vars.GetInt(0))
	println(vars.GetInt(1))
	println(vars.GetLong(2))
	println(vars.GetLong(4))
	println(vars.GetFloat(6))
	println(vars.GetDouble(7))
	println(vars.GetRef(9))
}
func printClassInfo(cf *classfile.ClassFile) {
	fmt.Printf("version: %v.%v\n", cf.MajorVersion(), cf.MinorVersion())
	fmt.Printf("constants count: %v\n", len(cf.ConstantPool()))
	fmt.Printf("access flags: 0x%x\n", cf.AccessFlags())
	fmt.Printf("this class: %v\n", cf.ClassName())
	fmt.Printf("super class: %v\n", cf.SuperClassName())
	fmt.Printf("interfaces: %v\n", cf.InterfaceNames())
	fmt.Printf("fields count: %v\n", len(cf.Fields()))
	for _, f := range cf.Fields() {
		fmt.Printf("		%s\n", f.Name())
	}
	fmt.Printf("methods count: %v\n", len(cf.Methods()))
	for _, m := range cf.Methods() {
		fmt.Printf("		%s\n", m.Name())
	}
}
func loadClass(className string, cp *classpath.Classpath) *classfile.ClassFile {
	classData, _, err := cp.ReadClass(className)
	if err != nil {
		panic(err)
	}
	cf, err := classfile.Parse(classData)
	if err != nil {
		panic(err)
	}
	return cf
}
