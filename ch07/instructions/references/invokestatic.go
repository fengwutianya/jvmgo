package references

import "jvmgo/ch07/instructions/base"
import "jvmgo/ch07/rtda"
import "jvmgo/ch07/rtda/heap"

type INVOKE_STATIC struct {
	base.Index16Instruction	//method reference index
}
func (self *INVOKE_STATIC) Execute(frame *rtda.Frame) {
	cp := frame.Method().Class().ConstantPool()
	methodRef := cp.GetConstant(self.Index).(*heap.MethodRef)
	resolvedMethod := methodRef.ResolvedMethod()
	if !resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}
	//resovled* 已经调用了ResolvedClass()把类加载进来了，只是还没调用<clinit>()V进行初始化
	class := resolvedMethod.Class()
	if !class.InitStarted() {	//调用静态方法所在类如果还未初始化，则先调用初始化函数，先把pc从thread保存到frame里面
		frame.RevertNextPC()
		base.InitClass(frame.Thread(), class)
		return
	}
	//不用检查访问权限因为编译的时候已经检查过了，生成了字节码说明访问权限没问题
	//不能调用类初始化方法<clinit>，不过这一点由class文件验证器保证
	base.InvokeMethod(frame, resolvedMethod)
}
// Invoke a class (static) method
//type INVOKE_STATIC struct{ base.Index16Instruction }
//
//func (self *INVOKE_STATIC) Execute(frame *rtda.Frame) {
//	cp := frame.Method().Class().ConstantPool()
//	methodRef := cp.GetConstant(self.Index).(*heap.MethodRef)
//	resolvedMethod := methodRef.ResolvedMethod()
//	if !resolvedMethod.IsStatic() {
//		panic("java.lang.IncompatibleClassChangeError")
//	}
//
//	class := resolvedMethod.Class()
//	if !class.InitStarted() {
//		frame.RevertNextPC()
//		base.InitClass(frame.Thread(), class)
//		return
//	}
//
//	base.InvokeMethod(frame, resolvedMethod)
//}
//