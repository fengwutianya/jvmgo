package rtda

// stack frame
type Frame struct {
	//线程执行函数时的栈帧--------------局部变量表，操作数栈  这标志着每个函数调用
	lower        *Frame // stack is implemented as linked list 非必需，用链表实现，这是代价
	localVars    LocalVars		//局部变量表
	operandStack *OperandStack	//操作数栈 这辆是由字节码操纵的
	// todo
}

func NewFrame(maxLocals, maxStack uint) *Frame {
	return &Frame{
		localVars:    newLocalVars(maxLocals),
		operandStack: newOperandStack(maxStack),
	}
}

// getters
func (self *Frame) LocalVars() LocalVars {
	return self.localVars
}
func (self *Frame) OperandStack() *OperandStack {
	return self.operandStack
}
