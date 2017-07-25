package references

import "jvmgo/ch06/instructions/base"
import "jvmgo/ch06/rtda"
import "jvmgo/ch06/rtda/heap"

// Create new object
//type NEW struct{ base.Index16Instruction }
//
//func (self *NEW) Execute(frame *rtda.Frame) {
//	cp := frame.Method().Class().ConstantPool()
//	classRef := cp.GetConstant(self.Index).(*heap.ClassRef)
//	class := classRef.ResolvedClass()
//	 todo: init class
	//
	//if class.IsInterface() || class.IsAbstract() {
	//	panic("java.lang.InstantiationError")
	//}
	//
	//ref := class.NewObject()
	//frame.OperandStack().PushRef(ref)
//}
type NEW struct {
	base.Index16Instruction	//new 加载类之后，new + 运行时常量池中classref的索引uint
}

func (self *NEW) Execute(frame *rtda.Frame) {
	cp := frame.Method().Class().ConstantPool()	//栈帧所对应的方法，所在类对应的运行时常量池
	classRef := cp.GetConstant(self.Index).(*heap.ClassRef)	//常量池中的classref
	//根据classref拿到class
	class := classRef.ResolvedClass()
	if class.IsInterface() || class.IsAbstract() {
		panic("java.lang.InstantiationError")
	}
	//class 是模板，产生对象
	ref := class.NewObject()
	// 把对象引用推入操作数栈栈顶
	frame.OperandStack().PushRef(ref)
}
