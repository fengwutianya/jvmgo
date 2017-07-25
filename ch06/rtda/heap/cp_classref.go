package heap

import "jvmgo/ch06/classfile"

/*
classref
	symref
		constantpool*
		classname
		class*
 */
type ClassRef struct {
	SymRef	//只有三个信息所属常量池，所属类类名和缓存的class引用，第一次解析，以后直接用，不同于class常量池中只存在符号引用，不存在直接饮用的缓存
}

func newClassRef(cp *ConstantPool, classInfo *classfile.ConstantClassInfo) *ClassRef {
	ref := &ClassRef{}
	ref.cp = cp
	ref.className = classInfo.Name()
	return ref
}
