package heap

//type Object struct {
//	class  *Class
//	fields Slots
//}
type Object struct {
	class 	*Class		//该对象所属类在方法区对应的Class引用
	fields	Slots		//所有的实例变量存储
}

func newObject(class *Class) *Object {
	return &Object{
		class:  class,
		fields: newSlots(class.instanceSlotCount),
	}
}

// getters
func (self *Object) Class() *Class {
	return self.class
}
func (self *Object) Fields() Slots {
	return self.fields
}

func (self *Object) IsInstanceOf(class *Class) bool {
	return class.isAssignableFrom(self.class)
}
