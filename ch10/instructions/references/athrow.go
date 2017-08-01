package references

import "reflect"
import "jvmgo/ch10/instructions/base"
import "jvmgo/ch10/rtda"
import "jvmgo/ch10/rtda/heap"

// Throw exception or error
type ATHROW struct{ base.NoOperandsInstruction }

func (self *ATHROW) Execute(frame *rtda.Frame) {
	ex := frame.OperandStack().PopRef()
	if ex == nil {
		panic("java.lang.NullPointerException")
	}
	thread := frame.Thread()
	if !findAndGotoExceptionHandler(thread, ex) {
		handleUncaughtException(thread, ex)
	}
	//ex := frame.OperandStack().PopRef()
	//if ex == nil {
	//	panic("java.lang.NullPointerException")
	//}
	//
	//thread := frame.Thread()
	//if !findAndGotoExceptionHandler(thread, ex) {
	//	handleUncaughtException(thread, ex)
	//}
}
func findAndGotoExceptionHandler(thread *rtda.Thread, ex *heap.Object) bool {
	for {
		frame := thread.CurrentFrame()
		pc := frame.NextPC() - 1

		handlerPC := frame.Method().FindExceptionHandler(ex.Class(), pc)
		if handlerPC > 0 {
			stack := frame.OperandStack()
			stack.Clear()
			stack.PushRef(ex)
			frame.SetNextPC(handlerPC)	//如果找到可以处理ex的代码块，把ex放在操作数栈栈顶，设置当前帧也是栈顶帧的pc
										//接下来就执行handlerpc后面的指令啦
			return true
		}
		thread.PopFrame()
		if thread.IsStackEmpty() {
			break
		}
	}
	return false
}

//func findAndGotoExceptionHandler(thread *rtda.Thread, ex *heap.Object) bool {
//	for {
//		frame := thread.CurrentFrame()
//		pc := frame.NextPC() - 1	//因为endPc不包含在内，所以要-1
//
//		handlerPC := frame.Method().FindExceptionHandler(ex.Class(), pc)
//		if handlerPC > 0 {
//			stack := frame.OperandStack()
//			stack.Clear()
//			stack.PushRef(ex)
//			frame.SetNextPC(handlerPC)
//			return true
//		}
//
//		thread.PopFrame()
//		if thread.IsStackEmpty() {
//			break
//		}
//	}
//	return false
//}

// todo 不知道这里是怎么实现的，也就是退栈完了也没有catch符合，那么打印出虚拟机栈调用信息
func handleUncaughtException(thread *rtda.Thread, ex *heap.Object) {
	thread.ClearStack()

	jMsg := ex.GetRefVar("detailMessage", "Ljava/lang/String;")
	goMsg := heap.GoString(jMsg)
	println(ex.Class().JavaName() + ": " + goMsg)

	stes := reflect.ValueOf(ex.Extra())
	for i := 0; i < stes.Len(); i++ {
		ste := stes.Index(i).Interface().(interface {
			String() string
		})
		println("\tat " + ste.String())
	}
}
