package base

import "jvmgo/ch07/rtda"
import "jvmgo/ch07/rtda/heap"

// jvms 5.5
func InitClass(thread *rtda.Thread, class *heap.Class) {
	class.StartInit()				//置位initStarted
	scheduleClinit(thread, class)
	initSuperClass(thread, class)
}

//子类的<clinit>先压栈
func scheduleClinit(thread *rtda.Thread, class *heap.Class) {
	clinit := class.GetClinitMethod()
	if clinit != nil {
		// exec <clinit>
		newFrame := thread.NewFrame(clinit)
		thread.PushFrame(newFrame)
	}
}

//
func initSuperClass(thread *rtda.Thread, class *heap.Class) {
	if !class.IsInterface() {
		superClass := class.SuperClass()
		if superClass != nil && !superClass.InitStarted() {
			InitClass(thread, superClass)
		}
	}
}
