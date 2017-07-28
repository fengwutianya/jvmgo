package heap

import "unicode/utf16"

//字符串池的实体
//map 代表字符串池的存储方式，key是string 也就是go的字符串 *Object也就是java的字符串
var internedStrings = map[string]*Object{}
//var internedStrings = map[string]*Object{}

// todo
// go string -> java.lang.String 也就是string.intern()的处理过程
func JString (loader *ClassLoader, goStr string) *Object {
	//map的存取方式 如果能找到 就返回
	if internedString, ok := internedStrings[goStr]; ok {
		return internedString
	}
	//找不到的话，存下来，先编码，然后按字符数组存储Object中
	chars := stringToUtf16(goStr) // go string -> java string
	jChars := &Object{loader.LoadClass("[C"), chars}	//构造对象 *Class data

	//最重要存的是java.lang.String，字符数组是value值，jString引用 里面有jChars引用
	jStr :=loader.LoadClass("java/lang/String").NewObject()
	jStr.SetRefVar("value", "[C", jChars)

	//放到字符串数组里面
	internedStrings[goStr] = jStr
	return jStr
}
//func JString(loader *ClassLoader, goStr string) *Object {
//	if internedStr, ok := internedStrings[goStr]; ok {
//		return internedStr
//	}
//
//	chars := stringToUtf16(goStr)
//	jChars := &Object{loader.LoadClass("[C"), chars}
//
//	jStr := loader.LoadClass("java/lang/String").NewObject()
//	jStr.SetRefVar("value", "[C", jChars)
//
//	internedStrings[goStr] = jStr
//	return jStr
//}

// java.lang.String -> go string
func GoString(jStr *Object) string {
	charArr := jStr.GetRefVar("value", "[C")
	return utf16ToString(charArr.Chars())
}

// utf8 -> utf16
func stringToUtf16(s string) []uint16 {
	runes := []rune(s)         // utf32
	return utf16.Encode(runes) // func Encode(s []rune) []uint16
}

// utf16 -> utf8
func utf16ToString(s []uint16) string {
	runes := utf16.Decode(s) // func Decode(s []uint16) []rune
	return string(runes)
}
