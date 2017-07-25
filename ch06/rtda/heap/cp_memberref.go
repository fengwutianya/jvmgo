package heap

import "jvmgo/ch06/classfile"

//对cpfieldref，cpmethodref，cpinterfacere的抽象
type MemberRef struct {
	SymRef
	name       string
	descriptor string
}

//class文件常量池对应的ConstantXXXXInfo也是这么抽象的
func (self *MemberRef) copyMemberRefInfo(refInfo *classfile.ConstantMemberrefInfo) {
	self.className = refInfo.ClassName()
	self.name, self.descriptor = refInfo.NameAndDescriptor()
}

func (self *MemberRef) Name() string {
	return self.name
}
func (self *MemberRef) Descriptor() string {
	return self.descriptor
}
