package heap

import "fmt"
import "jvmgo/ch06/classfile"
import "jvmgo/ch06/classpath"

/*
class names:
    - primitive types: boolean, byte, int ...
    - primitive arrays: [Z, [B, [I ...
    - non-array classes: java/lang/Object ...
    - array classes: [Ljava/lang/Object; ...
*/
//type ClassLoader struct {
//	cp       *classpath.Classpath
//	classMap map[string]*Class // loaded classes
//}
type ClassLoader struct {
	//ClassLoader from cp to classMap
	cp 			*classpath.Classpath
	classMap 	map[string]*Class 	//loaded classes，方法区可以看作是一个classMap，加载前为nil，加载后类的全限定名作为key
									//对应着*Class引用的value为已经加载的类
}

//func NewClassLoader(cp *classpath.Classpath) *ClassLoader {
//	return &ClassLoader{
//		cp:       cp,
//		classMap: make(map[string]*Class),
//	}
//}
func NewClassLoader(cp *classpath.Classpath) *ClassLoader {
	return &ClassLoader{
		cp:cp,
		classMap:make(map[string]*Class),
	}
}

//func (self *ClassLoader) LoadClass(name string) *Class {
//	if class, ok := self.classMap[name]; ok {
//		 already loaded
		//return class
	//}
	//
	//return self.loadNonArrayClass(name)
//}
func (self *ClassLoader) LoadClass(name string) *Class {
	if class, ok := self.classMap[name]; ok {
		//already loaded
		return class
	}
	return self.loadNonArrayClass(name)
}
func (self *ClassLoader) loadNonArrayClass(name string) *Class {
	//读取到内存
	data, entry := self.readClass(name)
	//解析成Class对象，放到方法区
	class := self.defineClass(data)
	//进行链接
	link(class)
	fmt.Printf("[Loaded %s from %s]\n", name, entry)
	return class
}
func link(class *Class) {
	//前面是加载
	//接下来是验证
	verify(class)
	//准备
	prepare(class)

}
func prepare(class *Class) {
	calcInstanceFieldSlotIds(class)	//给Class.Fields里面每一个Field编号
	calcStaticFieldSlotIds(class)	//给Class.StaticVars加标号
	allocAndInitStaticVars(class)	//类变量初始化
}
func allocAndInitStaticVars(class *Class) {
	class.staticVars = newSlots(class.staticSlotCount)
	for _, field := range class.fields {
		if field.IsStatic() && field.IsFinal() {	//只有static final才会出现在class文件的常量池中
			initStaticFinalVar(class, field)
		}
	}
}
func initStaticFinalVar(class *Class, field *Field) {
	vars := class.staticVars
	cp := class.constantPool
	cpIndex := field.ConstValueIndex() //第几个constvalue
	slotId := field.SlotId()
	if cpIndex > 0 {					//是constvalue
		switch field.Descriptor() {
		case "Z", "B", "C", "S", "I":				//boolean, byte, char, short, int 5个可以用int存储的类型
			val := cp.GetConstant(cpIndex).(int32)
			vars.SetInt(slotId, val)
		case "J":
			val := cp.GetConstant(cpIndex).(int64)
			vars.SetLong(slotId, val)
		case "F":
			val := cp.GetConstant(cpIndex).(float32)
			vars.SetFloat(slotId, val)
		case "D":
			val := cp.GetConstant(cpIndex).(float64)
			vars.SetDouble(slotId, val)
		case "Ljava/lang/String;":
			panic("todo")	//todo
		}
	}
}
func calcStaticFieldSlotIds(class *Class) {
	slotId := uint(0)
	//这里没有对父类静态变量进行操作 已经说明了一个问题，静态变量随着类，已经属于子类所有了
	//如果子类对静态函数或者静态变量进行覆盖，那么从子类里就没有父类的对应静态信息了
	for _, field := range class.fields {
		if field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.staticSlotCount = slotId
}
func calcInstanceFieldSlotIds(class *Class) {
	slotId := uint(0)
	if class.superClass != nil {
		slotId = class.superClass.instanceSlotCount	//由于加载每个类都需要做这个工作，因此吧整个继承体系的都算进去了
	}
	for _, field := range class.fields {
		if !field.IsStatic() {		//静态变量和实例变量分开标号
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.instanceSlotCount = slotId
}
func verify(class *Class) {
	//todo
}

func (self *ClassLoader) readClass(name string) ([]byte, classpath.Entry) {
	//使用路径对象加载类到内存形成字节数组
	data, entry, err := self.cp.ReadClass(name)
	if err != nil {
		panic("java.lang.ClassNotFoundException: " + name)
	}
	return data, entry
}

func (self *ClassLoader) defineClass(data []byte) *Class {
	//parseClass将字节数组解析成Class对象
	class := parseClass(data)
	//自己加载的class，标注上
	class.loader = self
	//解析符号引用到真实引用
	resolveSuperClass(class)
	resolveInterfaces(class)
	//这就叫放入方法区
	self.classMap[class.name] = class
	return class
}
func resolveInterfaces(class *Class) {
	interfaceCount := len(class.interfaceNames)
	if interfaceCount > 0 {
		class.interfaces = make([]*Class, interfaceCount)
		for i, interfaceName := range class.interfaceNames {
			class.interfaces[i] = class.loader.LoadClass(interfaceName)
		}
	}
}
func resolveSuperClass(class *Class) {
	//除非是Object类，所有其他类都有超类，用本身classloader来加载超类即可，这里算是递归加载类
	if class.name != "java/lang/Object" {
		class.superClass = class.loader.LoadClass(class.superClassName)
	}
}
func parseClass(data []byte) *Class {
	cf, err := classfile.Parse(data)
	if err != nil {
		panic("java.lang.ClassFormatError")
	}
	return newClass(cf)
}


//func (self *ClassLoader) loadNonArrayClass(name string) *Class {
//	data, entry := self.readClass(name)
//	class := self.defineClass(data)
//	link(class)
//	fmt.Printf("[Loaded %s from %s]\n", name, entry)
//	return class
//}
//
//func (self *ClassLoader) readClass(name string) ([]byte, classpath.Entry) {
//	data, entry, err := self.cp.ReadClass(name)
//	if err != nil {
//		panic("java.lang.ClassNotFoundException: " + name)
//	}
//	return data, entry
//}
//
// jvms 5.3.5
//func (self *ClassLoader) defineClass(data []byte) *Class {
//	class := parseClass(data)
//	class.loader = self
//	resolveSuperClass(class)
//	resolveInterfaces(class)
//	self.classMap[class.name] = class
//	return class
//}
//
//func parseClass(data []byte) *Class {
//	cf, err := classfile.Parse(data)
//	if err != nil {
//		panic("java.lang.ClassFormatError")
		//panic(err)
	//}
	//return newClass(cf)
//}
//
// jvms 5.4.3.1
//func resolveSuperClass(class *Class) {
//	if class.name != "java/lang/Object" {
//		class.superClass = class.loader.LoadClass(class.superClassName)
//	}
//}
//func resolveInterfaces(class *Class) {
//	interfaceCount := len(class.interfaceNames)
//	if interfaceCount > 0 {
//		class.interfaces = make([]*Class, interfaceCount)
//		for i, interfaceName := range class.interfaceNames {
//			class.interfaces[i] = class.loader.LoadClass(interfaceName)
//		}
//	}
//}
//
//func link(class *Class) {
//	verify(class)
//	prepare(class)
//}
//
//func verify(class *Class) {
//	 todo
//}
//
// jvms 5.4.2
//func prepare(class *Class) {
//	calcInstanceFieldSlotIds(class)
//	calcStaticFieldSlotIds(class)
//	allocAndInitStaticVars(class)
//}
//
//func calcInstanceFieldSlotIds(class *Class) {
//	slotId := uint(0)
//	if class.superClass != nil {
//		slotId = class.superClass.instanceSlotCount
//	}
//	for _, field := range class.fields {
//		if !field.IsStatic() {
//			field.slotId = slotId
//			slotId++
//			if field.isLongOrDouble() {
//				slotId++
//			}
//		}
//	}
//	class.instanceSlotCount = slotId
//}
//
//func calcStaticFieldSlotIds(class *Class) {
//	slotId := uint(0)
//	for _, field := range class.fields {
//		if field.IsStatic() {
//			field.slotId = slotId
//			slotId++
//			if field.isLongOrDouble() {
//				slotId++
//			}
//		}
//	}
//	class.staticSlotCount = slotId
//}
//
//func allocAndInitStaticVars(class *Class) {
//	class.staticVars = newSlots(class.staticSlotCount)
//	for _, field := range class.fields {
//		if field.IsStatic() && field.IsFinal() {
//			initStaticFinalVar(class, field)
//		}
//	}
//}
//
//func initStaticFinalVar(class *Class, field *Field) {
//	vars := class.staticVars
//	cp := class.constantPool
//	cpIndex := field.ConstValueIndex()
//	slotId := field.SlotId()
//
//	if cpIndex > 0 {
//		switch field.Descriptor() {
//		case "Z", "B", "C", "S", "I":
//			val := cp.GetConstant(cpIndex).(int32)
//			vars.SetInt(slotId, val)
//		case "J":
//			val := cp.GetConstant(cpIndex).(int64)
//			vars.SetLong(slotId, val)
//		case "F":
//			val := cp.GetConstant(cpIndex).(float32)
//			vars.SetFloat(slotId, val)
//		case "D":
//			val := cp.GetConstant(cpIndex).(float64)
//			vars.SetDouble(slotId, val)
//		case "Ljava/lang/String;":
//			panic("todo")
//		}
//	}
//}
//