package references

import "jvmgo/ch08/instructions/base"
import "jvmgo/ch08/rtda"
import "jvmgo/ch08/rtda/heap"

const (
	//数组内部元素的类型 基本元素8中 boolean byte short char int float long double
	//Array Type  atype
	AT_BOOLEAN = 4
	AT_CHAR    = 5
	AT_FLOAT   = 6
	AT_DOUBLE  = 7
	AT_BYTE    = 8
	AT_SHORT   = 9
	AT_INT     = 10
	AT_LONG    = 11
)

// Create new array
//两个操作数 一个在字节码里 直接在后面找到的byte型表述数组元素类型 另外一个是操作数栈里面弹出，表示数组长度
//type NEW_ARRAY struct {
//	atype uint8
//}
type NEW_ARRAY struct {
	atype uint8
}

func (self *NEW_ARRAY) FetchOperands(reader *base.BytecodeReader) {
	self.atype = reader.ReadUint8()
}

func (self *NEW_ARRAY) Execute(frame *rtda.Frame) {
	//数组长度 在操作数栈中
	stack := frame.OperandStack()
	count := stack.PopInt()
	if count < 0 {
		panic("java.lang.NegativeArraySizeException")
	}
	classLoader := frame.Method().Class().Loader()
	arrClass := getPrimitiveArrayClass(classLoader, self.atype)	//newarray只能产生基本类型的数组
	arr := arrClass.NewArray(uint(count))
	stack.PushRef(arr)
}
//func (self *NEW_ARRAY) FetchOperands(reader *base.BytecodeReader) {
//	self.atype = reader.ReadUint8()
//}
//func (self *NEW_ARRAY) Execute(frame *rtda.Frame) {
//	stack := frame.OperandStack()
//	count := stack.PopInt()
//	if count < 0 {
//		panic("java.lang.NegativeArraySizeException")
//	}
//
//	classLoader := frame.Method().Class().Loader()
//	arrClass := getPrimitiveArrayClass(classLoader, self.atype)
//	arr := arrClass.NewArray(uint(count))
//	stack.PushRef(arr)
//}

func getPrimitiveArrayClass(loader *heap.ClassLoader, atype uint8) *heap.Class {
	switch atype {
	case AT_BOOLEAN:
		return loader.LoadClass("[Z")
	case AT_BYTE:
		return loader.LoadClass("[B")
	case AT_CHAR:
		return loader.LoadClass("[C")
	case AT_SHORT:
		return loader.LoadClass("[S")
	case AT_INT:
		return loader.LoadClass("[I")
	case AT_LONG:
		return loader.LoadClass("[J")
	case AT_FLOAT:
		return loader.LoadClass("[F")
	case AT_DOUBLE:
		return loader.LoadClass("[D")
	default:
		panic("Invalid atype!")
	}
}
