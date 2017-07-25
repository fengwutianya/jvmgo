package heap

import "jvmgo/ch06/classfile"

type Field struct {
	ClassMember				//accessflags, name, descriptor, *class
	constValueIndex uint
	slotId          uint	//给一个field都有一个编号，因为所有的字段值都存在[]Slot里面，大小又是变化的，因此需要标号
							//来确定哪个变量的值是哪个
							//需要注意的是 静态变量存在方法区里面，在ConstantPool里面用[]Slot存储，存在于每一个Class里面
							//而实例变量是在堆里，也就是每个Object里面用[]Slot存储
}

//func newFields(class *Class, cfFields []*classfile.MemberInfo) []*Field {
//	fields := make([]*Field, len(cfFields))
//	for i, cfField := range cfFields {
//		fields[i] = &Field{}
//		fields[i].class = class
//		fields[i].copyMemberInfo(cfField)
//		fields[i].copyAttributes(cfField)
//	}
//	return fields
//}
func newFields(class *Class, cfFields []*classfile.MemberInfo) []*Field {
	fields := make([]*Field, len(cfFields))
	for i, cfField := range cfFields {
		fields[i] = &Field{}
		fields[i].class = class
		fields[i].copyMemberInfo(cfField)
		fields[i].copyAttributes(cfField)
	}
	return fields
}
func (self *Field) copyAttributes(cfField *classfile.MemberInfo) {
	if valAttr := cfField.ConstantValueAttribute(); valAttr != nil {
		self.constValueIndex = uint(valAttr.ConstantValueIndex())
	}
}

func (self *Field) IsVolatile() bool {
	return 0 != self.accessFlags&ACC_VOLATILE
}
func (self *Field) IsTransient() bool {
	return 0 != self.accessFlags&ACC_TRANSIENT
}
func (self *Field) IsEnum() bool {
	return 0 != self.accessFlags&ACC_ENUM
}

func (self *Field) ConstValueIndex() uint {
	return self.constValueIndex
}
func (self *Field) SlotId() uint {
	return self.slotId
}
func (self *Field) isLongOrDouble() bool {
	return self.descriptor == "J" || self.descriptor == "D"
}
