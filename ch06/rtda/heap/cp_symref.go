package heap

// symbolic reference
type SymRef struct {		//常量池类信息cpclassref，方法信息cpmethodref信息，字段信息cpfieldref,接口方法interfacemethodref
							//后面三者又被抽象为memberref
	cp        *ConstantPool	//所在常量池的引用
	className string		//所属class名
	class     *Class		//所属class的引用
}

//func (self *SymRef) ResolvedClass() *Class {
//	if self.class == nil {
//		self.resolveClassRef()
//	}
//	return self.class
//}
func (self *SymRef) ResolvedClass() *Class {
	if self.class != nil {
		self.resolveClassRef()
	}
	return self.class
}

// jvms8 5.4.3.1
//func (self *SymRef) resolveClassRef() {
//	d := self.cp.class							//类的加载 是发生在 当前指令要加载另外一个类
//	c := d.loader.LoadClass(self.className)		//class引用 classloader对象 LoadClass方法加载要加载的类
//	if !c.isAccessibleTo(d) {
//		panic("java.lang.IllegalAccessError")
//	}
//
//	self.class = c
//}
func (self *SymRef) resolveClassRef() {
	d := self.cp.class							//本类去加载另外一个类，拿到要加载c的类，也就是c坐在的常量池所属于的类，，，
	c := d.loader.LoadClass(self.className)		//加载要加载的类c
	if !c.isAccessibleTo(d) {					//检查d能否访问c
		panic("java.lang.IllegalAccessError")
	}
	self.class = c
}
