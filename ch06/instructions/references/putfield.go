package references

import "jvmgo/ch06/instructions/base"
import "jvmgo/ch06/rtda"
import "jvmgo/ch06/rtda/heap"

// Set field in object
//type PUT_FIELD struct{ base.Index16Instruction }
type PUT_FIELD struct {
	base.Index16Instruction
}

func (self *PUT_FIELD) Execute(frame *rtda.Frame) {
	currentMethod := frame.Method()
	currentClass := currentMethod.Class()
	constPool := currentClass.ConstantPool()
	fieldRef := constPool.GetConstant(self.Index).(*heap.FieldRef)
	field := fieldRef.ResolvedField()

	if field.IsStatic() {	//此指令不能操作static变量
		panic("java.lang.IncompatibleClassChangeError")
	}
	if field.IsFinal() {	//非static的final变量只能在本类的普通构造函数中赋值<init>
		if currentClass != field.Class() || currentMethod.Name() != "<init>" {
			panic("java.lang.IllegalAccessError")
		}
	}

	//PUT_FIELD指令 需要三个操作数，第一个是找fieldref的index，第二个是要付的值，第三个是变量所在的类的引用
	descriptor := field.Descriptor()
	slotId := field.SlotId()	//和putstatic指令不同，slotid是实例中的slotid，
	stack := frame.OperandStack()
	switch descriptor[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		val := stack.PopInt()
		ref := stack.PopRef()	//此ref指的是Object对象
		if ref == nil {
			panic("java.lang.NullPointerException")
		}
		ref.Fields().SetInt(slotId, val)	//为甚恶魔不使用self.index，是因为index是为了找引用的，而slotid是值的索引
			case 'F':
				val := stack.PopFloat()
				ref := stack.PopRef()
				if ref == nil {
					panic("java.lang.NullPointerException")
				}
				ref.Fields().SetFloat(slotId, val)
			case 'J':
				val := stack.PopLong()
				ref := stack.PopRef()
				if ref == nil {
					panic("java.lang.NullPointerException")
				}
				ref.Fields().SetLong(slotId, val)
			case 'D':
				val := stack.PopDouble()
				ref := stack.PopRef()
				if ref == nil {
					panic("java.lang.NullPointerException")
				}
				ref.Fields().SetDouble(slotId, val)
			case 'L', '[':
				val := stack.PopRef()
				ref := stack.PopRef()
				if ref == nil {
					panic("java.lang.NullPointerException")
				}
				ref.Fields().SetRef(slotId, val)
			default:
		//		 todo
	}
}
//func (self *PUT_FIELD) Execute(frame *rtda.Frame) {
//	currentMethod := frame.Method()
//	currentClass := currentMethod.Class()
//	cp := currentClass.ConstantPool()
//	fieldRef := cp.GetConstant(self.Index).(*heap.FieldRef)
//	field := fieldRef.ResolvedField()
//
//	if field.IsStatic() {
//		panic("java.lang.IncompatibleClassChangeError")
//	}
//	if field.IsFinal() {
//		if currentClass != field.Class() || currentMethod.Name() != "<init>" {
//			panic("java.lang.IllegalAccessError")
//		}
//	}
//
//	descriptor := field.Descriptor()
//	slotId := field.SlotId()
//	stack := frame.OperandStack()
//
//	switch descriptor[0] {
//	case 'Z', 'B', 'C', 'S', 'I':
//		val := stack.PopInt()
//		ref := stack.PopRef()
//		if ref == nil {
//			panic("java.lang.NullPointerException")
//		}
//		ref.Fields().SetInt(slotId, val)
//	case 'F':
//		val := stack.PopFloat()
//		ref := stack.PopRef()
//		if ref == nil {
//			panic("java.lang.NullPointerException")
//		}
//		ref.Fields().SetFloat(slotId, val)
//	case 'J':
//		val := stack.PopLong()
//		ref := stack.PopRef()
//		if ref == nil {
//			panic("java.lang.NullPointerException")
//		}
//		ref.Fields().SetLong(slotId, val)
//	case 'D':
//		val := stack.PopDouble()
//		ref := stack.PopRef()
//		if ref == nil {
//			panic("java.lang.NullPointerException")
//		}
//		ref.Fields().SetDouble(slotId, val)
//	case 'L', '[':
//		val := stack.PopRef()
//		ref := stack.PopRef()
//		if ref == nil {
//			panic("java.lang.NullPointerException")
//		}
//		ref.Fields().SetRef(slotId, val)
//	default:
//		 todo
	//}
//}
