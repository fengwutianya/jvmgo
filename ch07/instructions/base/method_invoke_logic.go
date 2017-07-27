package base

import "fmt"
import "jvmgo/ch07/rtda"
import "jvmgo/ch07/rtda/heap"

func InvokeMethod(invokerFrame *rtda.Frame, method *heap.Method) {
	//thread := invokerFrame.Thread()		//获取当前线程，为了得到jvm栈，jvm栈是每个线程一个，所属于线程
	//newFrame := thread.NewFrame(method)	//每个栈帧对应一次函数调用
	//thread.PushFrame(newFrame)			//栈帧入栈
	thread := invokerFrame.Thread()
	newFrame := thread.NewFrame(method)
	thread.PushFrame(newFrame)

	argSlotSlot := int(method.ArgSlotCount())	//由于参数传递时，参数个数和slot不一一对应的关系，比方说由double long
												//另外就是 第一个参数永远的this，于是传slots
	if argSlotSlot > 0 {
		//参数传递
		for i := argSlotSlot - 1; i >= 0; i-- {
			slot := invokerFrame.OperandStack().PopSlot()
			newFrame.LocalVars().SetSlot(uint(i), slot)
		}
	}

	// hack!
	if method.IsNative() {
		if method.Name() == "registerNatives" {
			thread.PopFrame()
		} else {
			panic(fmt.Sprintf("native method: %v.%v%v\n",
				method.Class().Name(), method.Name(), method.Descriptor()))
		}
	}
}
