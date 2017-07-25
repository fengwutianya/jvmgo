package references

import "jvmgo/ch06/instructions/base"
import "jvmgo/ch06/rtda"
import "jvmgo/ch06/rtda/heap"

// Set static field in class
//type PUT_STATIC struct{ base.Index16Instruction }
type PUT_STATIC struct {
	base.Index16Instruction
}

func (self *PUT_STATIC) Execute(frame *rtda.Frame) {
	//注意当前类访问其他类的static属性，类1和类2是可以不同的
	currentMethod := frame.Method()							//当前方法
	currentClass := currentMethod.Class()					//当前类1
	cp := currentClass.ConstantPool()						//当前运行时常量池
	fieldRef := cp.GetConstant(self.Index).(*heap.FieldRef)	//从运行时常量池中拿到field的引用
	field := fieldRef.ResolvedField()						//解析field
	class := field.Class()									//field所在类2
	//todo: init class 2
	if !field.IsStatic() {									//非static抛出IncompatibleClassChangeError，不能用类名访问非静态属性
		panic("java.lang.IncompatibleClassChangeError")
	}
	if field.IsFinal() {									//static final属性只能是本类的、类初始化方法<clinit>给他赋值
		if currentClass != class || currentMethod.Name() != "<clinit>" {
			panic("java.lang.IllegalAccessError")
		}
	}
	descriptor := field.Descriptor()	//对于field来说 是ZBCSI/F/J/D/L/[
	slotId := field.SlotId()			//常量池里面从哪一项开始读，用slotId标识
	slots := class.StaticVars()
	stack := frame.OperandStack()
	switch descriptor[0] {
	case 'Z', 'B', 'C', 'S', 'I': slots.SetInt(slotId, stack.PopInt())
	case 'J': slots.SetLong(slotId, stack.PopLong())
	case 'D': slots.SetDouble(slotId, stack.PopDouble())
	case 'L', '[':
		slots.SetRef(slotId, stack.PopRef())
	default:
		//todo
	}
}

//func (self *PUT_STATIC) Execute(frame *rtda.Frame) {
//	currentMethod := frame.Method()
//	currentClass := currentMethod.Class()
//	cp := currentClass.ConstantPool()
//	fieldRef := cp.GetConstant(self.Index).(*heap.FieldRef)
//	field := fieldRef.ResolvedField()
//	class := field.Class()
//	 todo: init class
	//
	//if !field.IsStatic() {
	//	panic("java.lang.IncompatibleClassChangeError")
	//}
	//if field.IsFinal() {
	//	if currentClass != class || currentMethod.Name() != "<clinit>" {
	//		panic("java.lang.IllegalAccessError")
	//	}
	//}
	//
	//descriptor := field.Descriptor()
	//slotId := field.SlotId()
	//slots := class.StaticVars()
	//stack := frame.OperandStack()
	//
	//switch descriptor[0] {
	//case 'Z', 'B', 'C', 'S', 'I':
	//	slots.SetInt(slotId, stack.PopInt())
	//case 'F':
	//	slots.SetFloat(slotId, stack.PopFloat())
	//case 'J':
	//	slots.SetLong(slotId, stack.PopLong())
	//case 'D':
	//	slots.SetDouble(slotId, stack.PopDouble())
	//case 'L', '[':
	//	slots.SetRef(slotId, stack.PopRef())
	//default:
	//	 todo
	//}
//}
