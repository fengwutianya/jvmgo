package rtda

/*
JVM
  Thread
    pc
    Stack
      Frame
        LocalVars
        OperandStack
*/
type Thread struct {
	pc    int // the address of the instruction currently being executed
	stack *Stack	//线程栈 线程私有 线程包括多个函数调用，每个函数调用都会新创建一个函数栈帧，用Frame链表实现
	// todo
}

func NewThread() *Thread {
	return &Thread{
		stack: newStack(1024),
	}
}

func (self *Thread) PC() int {
	return self.pc
}
func (self *Thread) SetPC(pc int) {
	self.pc = pc
}

func (self *Thread) PushFrame(frame *Frame) {
	self.stack.push(frame)
}
func (self *Thread) PopFrame() *Frame {
	return self.stack.pop()
}

func (self *Thread) CurrentFrame() *Frame {
	return self.stack.top()
}
