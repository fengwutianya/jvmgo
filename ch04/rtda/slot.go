package rtda

//本地变量表和操作数栈都要求能存放一个int & 能存放一个引用 so
type Slot struct {
	num int32
	ref *Object
}
