package native

import "jvmgo/ch09/rtda"

//type NativeMethod func(frame *rtda.Frame)
type NativeMethod func(frame *rtda.Frame)	//这个才是go里面函数的interface吧 没有函数体 或者叫函数指针

//var registry = map[string]NativeMethod{}
var registry = map[string]NativeMethod{}

func emptyNativeMethod(frame *rtda.Frame) {
	// do nothing
}

//唯一确定一个函数 类名 函数名 描述符（包括参数表和返回值）
func Register(className, methodName, methodDescriptor string, method NativeMethod) {
	key := className + "~" + methodName + "~" + methodDescriptor
	registry[key] = method
}
//func Register(className, methodName, methodDescriptor string, method NativeMethod) {
//	key := className + "~" + methodName + "~" + methodDescriptor
//	registry[key] = method
//}

//查找本地函数
func FindNativeMethod(className, methodName, methodDescriptor string) NativeMethod {
	key := className + "~" + methodName + "~" + methodDescriptor
	if method, ok := registry[key]; ok {
		return method
	}
	//Object.registerNatives()被忽略，因为没实现JNI
	if methodName == "registerNatives" && methodDescriptor == "()V" {
		return emptyNativeMethod
	}
	return nil
}
//func FindNativeMethod(className, methodName, methodDescriptor string) NativeMethod {
//	key := className + "~" + methodName + "~" + methodDescriptor
//	if method, ok := registry[key]; ok {
//		return method
//	}
//	if methodDescriptor == "()V" && methodName == "registerNatives" {
//		return emptyNativeMethod
//	}
//	return nil
//}
