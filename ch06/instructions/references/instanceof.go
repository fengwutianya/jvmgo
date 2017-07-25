package references

import "jvmgo/ch06/instructions/base"
import "jvmgo/ch06/rtda"
import "jvmgo/ch06/rtda/heap"

// Determine if object is of given type
type INSTANCE_OF struct{ base.Index16Instruction }

func (self *INSTANCE_OF) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	ref := stack.PopRef()
	if ref == nil {
		stack.PushInt(0)	//null不是任何类的子类
		return
	}
	classRef := frame.Method().Class().ConstantPool().GetConstant(self.Index).(*heap.ClassRef)
	class := classRef.ResolvedClass()	//指令所带操作数对应的常量池索引索引到的常量转化为classref类型
	if ref.IsInstanceOf(class) {	//也就是说 操作数栈上的引用类是不是字节码索引到的引用类的子类
		stack.PushInt(1)
	} else {
		stack.PushInt(0)
	}
}
//func (self *INSTANCE_OF) Execute(frame *rtda.Frame) {
//	stack := frame.OperandStack()
//	ref := stack.PopRef()
//	if ref == nil {
//		stack.PushInt(0)
//		return
//	}
//
//	cp := frame.Method().Class().ConstantPool()
//	classRef := cp.GetConstant(self.Index).(*heap.ClassRef)
//	class := classRef.ResolvedClass()
//	if ref.IsInstanceOf(class) {
//		stack.PushInt(1)
//	} else {
//		stack.PushInt(0)
//	}
//}
