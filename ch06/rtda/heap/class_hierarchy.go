package heap

// jvms8 6.5.instanceof
// jvms8 6.5.checkcast
func (self *Class) isAssignableFrom(other *Class) bool {
	s, t := other, self

	if s == t {		//自己类是可以的
		return true
	}

	if !t.IsInterface() {
		return s.isSubClassOf(t)	//栈上所示类是不是操作码指示类的子类？
	} else {
		return s.isImplements(t)	//栈上是不是实现了操作码所对应的接口？
	}
}
func (self *Class) isImplements(other *Class) bool {
	for c := self; c != nil; c = c.superClass {
		for _, i := range c.interfaces	{
			if i == other || i.isSubInterfaceOf(other) {
				return true
			}
		}
	}
	return false
}
func (self *Class) isSubClassOf(other *Class) bool {
	for c := self.superClass; c != nil; c = c.superClass {
		if c == other {
			return true
		}
	}
	return false
}

// self extends c
//func (self *Class) isSubClassOf(other *Class) bool {
//	for c := self.superClass; c != nil; c = c.superClass {
//		if c == other {
//			return true
//		}
//	}
//	return false
//}

// self implements iface
//func (self *Class) isImplements(iface *Class) bool {
//	for c := self; c != nil; c = c.superClass {
//		for _, i := range c.interfaces {
//			if i == iface || i.isSubInterfaceOf(iface) {
//				return true
//			}
//		}
//	}
//	return false
//}

// self extends iface
//func (self *Class) isSubInterfaceOf(iface *Class) bool {
//	for _, superInterface := range self.interfaces {
//		if superInterface == iface || superInterface.isSubInterfaceOf(iface) {
//			return true
//		}
//	}
//	return false
//}
func (self *Class) isSubInterfaceOf(iface *Class) bool {
	for _, superInterface := range self.interfaces {
		if superInterface == iface || superInterface.isSubInterfaceOf(iface) {
			return true
		}
	}
	return false
}
