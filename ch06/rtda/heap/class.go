package heap

import "strings"
import "jvmgo/ch06/classfile"

//常量池 存放 类信息Class，字段信息，方法信息，类变量，类加载器

// name, superClassName and interfaceNames are all binary names(jvms8-4.2.1)
//type Class struct {
//	accessFlags       uint16
//	name              string // thisClassName
//	superClassName    string
//	interfaceNames    []string
//	constantPool      *ConstantPool
//	fields            []*Field
//	methods           []*Method
//	loader            *ClassLoader
//	superClass        *Class
//	interfaces        []*Class
//	instanceSlotCount uint
//	staticSlotCount   uint
//	staticVars        Slots
//}
//运行时常量池 也是属于某个类的 和class文件的常量池对应
type Class struct {
	accessFlags 		uint16			//ClassFile.AccessFlags()

	//字面量
	name 				string			//ClassFile.ClassName()
	superClassName		string			//ClassFile.SuperClassName()
	interfaceNames 		[]string		//ClassFile.InterfaceNames()

	constantPool 		*ConstantPool	//newConstantPool(class, ClassFile.ConstantPool())

	fields 				[]*Field		//newFields(class, ClassFile.Fields())
	methods 			[]*Method		//newMethods(class, ClassFile.Methods())

	loader 				*ClassLoader	//ClassLoader.LoadClass(defineClass(class.loader = self))

	//区别于符号引用，这里是真实引用
	superClass			*Class			//resolveSuperClass(Class)
	interfaces 			[]*Class		//resolveInterfaces(Class)

	instanceSlotCount	uint			//所有实例变量的大小
	staticSlotCount		uint			//所有静态变量的大小

	staticVars			Slots	//静态变量 方法区
}

//func newClass(cf *classfile.ClassFile) *Class {
//	class := &Class{}
//	class.accessFlags = cf.AccessFlags()
//	class.name = cf.ClassName()
//	class.superClassName = cf.SuperClassName()
//	class.interfaceNames = cf.InterfaceNames()
//	class.constantPool = newConstantPool(class, cf.ConstantPool())
//	class.fields = newFields(class, cf.Fields())
//	class.methods = newMethods(class, cf.Methods())
//	return class
//}
func newClass(cf *classfile.ClassFile) *Class {
	class := &Class{}
	class.accessFlags = cf.AccessFlags()
	class.name = cf.ClassName()
	class.superClassName = cf.SuperClassName()
	class.interfaceNames = cf.InterfaceNames()
	class.constantPool = newConstantPool(class, cf.ConstantPool())
	class.fields = newFields(class, cf.Fields())
	class.methods = newMethods(class, cf.Methods())
	return class
}

//func (self *Class) IsPublic() bool {
//	return 0 != self.accessFlags&ACC_PUBLIC
//}
func (self *Class) IsPublic() bool {
	return (self.accessFlags & ACC_PUBLIC) != 0
}
func (self *Class) IsFinal() bool {
	return 0 != self.accessFlags&ACC_FINAL
}
func (self *Class) IsSuper() bool {
	return 0 != self.accessFlags&ACC_SUPER
}
func (self *Class) IsInterface() bool {
	return 0 != self.accessFlags&ACC_INTERFACE
}
func (self *Class) IsAbstract() bool {
	return 0 != self.accessFlags&ACC_ABSTRACT
}
func (self *Class) IsSynthetic() bool {
	return 0 != self.accessFlags&ACC_SYNTHETIC
}
func (self *Class) IsAnnotation() bool {
	return 0 != self.accessFlags&ACC_ANNOTATION
}
func (self *Class) IsEnum() bool {
	return 0 != self.accessFlags&ACC_ENUM
}

// getters
func (self *Class) ConstantPool() *ConstantPool {
	return self.constantPool
}
func (self *Class) StaticVars() Slots {
	return self.staticVars
}

// jvms 5.4.4
func (self *Class) isAccessibleTo(other *Class) bool {
	return self.IsPublic() ||
		self.getPackageName() == other.getPackageName()
}

func (self *Class) getPackageName() string {
	if i := strings.LastIndex(self.name, "/"); i >= 0 {
		return self.name[:i]
	}
	return ""
}

func (self *Class) GetMainMethod() *Method {
	return self.getStaticMethod("main", "([Ljava/lang/String;)V")
}

func (self *Class) getStaticMethod(name, descriptor string) *Method {
	for _, method := range self.methods {
		if method.IsStatic() &&
			method.name == name &&
			method.descriptor == descriptor {

			return method
		}
	}
	return nil
}

func (self *Class) NewObject() *Object {
	return newObject(self)
}
