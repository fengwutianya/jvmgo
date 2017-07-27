package references

import "jvmgo/ch07/instructions/base"
import "jvmgo/ch07/rtda"
import "jvmgo/ch07/rtda/heap"

// Invoke instance method;
// special handling for superclass, private, and instance initialization method invocations
//编译器确定函数调用的一些函数 super.method(), private void method(), <init>()V，全部是实例方法，非static
type INVOKE_SPECIAL struct{ base.Index16Instruction }

func (self *INVOKE_SPECIAL) Execute(frame *rtda.Frame) {
	//currentClass当前栈帧所属函数所属类的class resolvedClass调用函数所属class，ref this变量
	currentClass := frame.Method().Class()
	cp := currentClass.ConstantPool()
	methodRef := cp.GetConstant(self.Index).(*heap.MethodRef)
	resolvedClass := methodRef.ResolvedClass()
	resolvedMethod := methodRef.ResolvedMethod()
	if resolvedMethod.Name() == "<init>" && resolvedMethod.Class() != resolvedClass {	//本类调用自己的构造函数
		panic("java.lang.NoSuchMethodError")
	}
	if resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	//拿到this引用，不能破坏栈，用getRefFromTop 从栈顶n-1个变量时this
	ref := frame.OperandStack().GetRefFromTop(resolvedMethod.ArgSlotCount() - 1)
	if ref == nil {
		panic("java.lang.NullPointerException")
	}
	//访问权限。protected 子类 本类
	//这里学到一个访问权限的知识点就是 protected 同一包也可以访问
	//对于当前栈所属方法所在类currentClass，调用父类resolvedClass中的方法，调用者ref不是自己，也不是子类，无访问protected权限
	if resolvedMethod.IsProtected() &&
		resolvedMethod.Class().IsSuperClassOf(currentClass) && //子类调用父类方法protected
		resolvedMethod.Class().GetPackageName() != currentClass.GetPackageName() && //同一包的类也可以
		ref.Class() != currentClass &&
		!ref.Class().IsSubClassOf(currentClass) {

		panic("java.lang.IllegalAccessError")
	}

	methodToBeInvoked := resolvedMethod
	if currentClass.IsSuper() &&
		resolvedClass.IsSuperClassOf(currentClass) &&
		resolvedMethod.Name() != "<init>" {

		//对invokenovirtual的改正，invokespecial也要找到继承树中最低父节点拥有此函数的父类的方法，super.method()也是动态绑定了？
		methodToBeInvoked = heap.LookupMethodInClass(currentClass.SuperClass(),
			methodRef.Name(), methodRef.Descriptor())
	}

	if methodToBeInvoked == nil || methodToBeInvoked.IsAbstract() {
		panic("java.lang.AbstractMethodError")
	}

	base.InvokeMethod(frame, methodToBeInvoked)
}
