package heap

type Object struct {
	class *Class
	//原来是slots []Slot运来存放对象的各个实例变量值，现在改成interface{}也就是void*类型
	//可以存放任何变量，对于普通类，存放[]Slot，对于数组类，存放各种数组[]int8 []int16 []int32 []int64 []float32 []float64 []*Object []uint16
	data  interface{} // Slots for Object, []int32 for int[] ...
}

// create normal (non-array) object
func newObject(class *Class) *Object {
	return &Object{
		class: class,
		data:  newSlots(class.instanceSlotCount),
	}
}

// getters
func (self *Object) Class() *Class {
	return self.class
}
func (self *Object) Fields() Slots {
	return self.data.(Slots)
}

func (self *Object) IsInstanceOf(class *Class) bool {
	return class.isAssignableFrom(self.class)
}

// reflection
func (self *Object) GetRefVar(name, descriptor string) *Object {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	return slots.GetRef(field.slotId)
}
func (self *Object) SetRefVar(name, descriptor string, ref *Object) {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	slots.SetRef(field.slotId, ref)
}
