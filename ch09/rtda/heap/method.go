package heap

import "jvmgo/ch09/classfile"

type Method struct {
	ClassMember
	maxStack     uint
	maxLocals    uint
	code         []byte
	argSlotCount uint	//局部变量所占slot数目
}

func newMethods(class *Class, cfMethods []*classfile.MemberInfo) []*Method {
	methods := make([]*Method, len(cfMethods))
	for i, cfMethod := range cfMethods {
		methods[i] = newMethod(class, cfMethod)
	}
	return methods
}

func newMethod(class *Class, cfMethod *classfile.MemberInfo) *Method {
	method := &Method{}
	method.class = class
	method.copyMemberInfo(cfMethod)
	method.copyAttributes(cfMethod)
	md := parseMethodDescriptor(method.descriptor)
	method.calcArgSlotCount(md.parameterTypes)
	if method.IsNative() {
		method.injectCodeAttribute(md.returnType)
	}
	return method
}

func (self *Method) copyAttributes(cfMethod *classfile.MemberInfo) {
	if codeAttr := cfMethod.CodeAttribute(); codeAttr != nil {
		self.maxStack = codeAttr.MaxStack()
		self.maxLocals = codeAttr.MaxLocals()
		self.code = codeAttr.Code()
	}
}

func (self *Method) calcArgSlotCount(paramTypes []string) {
	for _, paramType := range paramTypes {
		self.argSlotCount++
		if paramType == "J" || paramType == "D" {
			self.argSlotCount++
		}
	}
	if !self.IsStatic() {
		self.argSlotCount++ // `this` reference
	}
}

//Native方法要injectCodeAttribute没有字节码code
func (self *Method) injectCodeAttribute(returnType string) {
	self.maxStack = 4 // todo
	self.maxLocals = self.argSlotCount
	switch returnType[0] {
	case 'V':
		//0Xfe使用保留指令来标识native方法，后面一个标识返回值，没有字节码就没有指令后面的操作数，函数形参都在栈里
		self.code = []byte{0xfe, 0xb1} // return
	case 'L', '[':
		self.code = []byte{0xfe, 0xb0} // areturn
	case 'D':
		self.code = []byte{0xfe, 0xaf} // dreturn
	case 'F':
		self.code = []byte{0xfe, 0xae} // freturn
	case 'J':
		self.code = []byte{0xfe, 0xad} // lreturn
	default:
		self.code = []byte{0xfe, 0xac} // ireturn
	}
}

func (self *Method) IsSynchronized() bool {
	return 0 != self.accessFlags&ACC_SYNCHRONIZED
}
func (self *Method) IsBridge() bool {
	return 0 != self.accessFlags&ACC_BRIDGE
}
func (self *Method) IsVarargs() bool {
	return 0 != self.accessFlags&ACC_VARARGS
}
func (self *Method) IsNative() bool {
	return 0 != self.accessFlags&ACC_NATIVE
}
func (self *Method) IsAbstract() bool {
	return 0 != self.accessFlags&ACC_ABSTRACT
}
func (self *Method) IsStrict() bool {
	return 0 != self.accessFlags&ACC_STRICT
}

// getters
func (self *Method) MaxStack() uint {
	return self.maxStack
}
func (self *Method) MaxLocals() uint {
	return self.maxLocals
}
func (self *Method) Code() []byte {
	return self.code
}
func (self *Method) ArgSlotCount() uint {
	return self.argSlotCount
}
