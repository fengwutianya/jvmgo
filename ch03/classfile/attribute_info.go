package classfile

/*
和constant pool一样个结构各不相同，无法用统一结构来描述
不同是常量是规定14种，可以用u1大小的tag来严格定义
而属性是可扩展的，因此要属性名来区分，属性名也存在于常量池里面，用u2 uint16来索引
attribute_info {
    u2 attribute_name_index;
    u4 attribute_length;
    u1 info[attribute_length];
}
*/
type AttributeInfo interface {
	readInfo(reader *ClassReader)
}

func readAttributes(reader *ClassReader, cp ConstantPool) []AttributeInfo {
	attributesCount := reader.readUint16()
	attributes := make([]AttributeInfo, attributesCount)
	for i := range attributes {
		attributes[i] = readAttribute(reader, cp)
	}
	return attributes
}

func readAttribute(reader *ClassReader, cp ConstantPool) AttributeInfo {
	attrNameIndex := reader.readUint16()				//根据index去常量池寻找Attribute Name
	attrName := cp.getUtf8(attrNameIndex)				//在cp里面找到了 都出来
	attrLen := reader.readUint32()						//得到attribute的length，确定attribute的大小
	attrInfo := newAttributeInfo(attrName, attrLen, cp)	//根据名字和大小创建attribute info
	attrInfo.readInfo(reader)							//填充attributeInfo的内容
	return attrInfo
}

func newAttributeInfo(attrName string, attrLen uint32, cp ConstantPool) AttributeInfo {
	switch attrName {
	case "Code":
		return &CodeAttribute{cp: cp}
	case "ConstantValue":
		return &ConstantValueAttribute{}
	case "Deprecated":
		return &DeprecatedAttribute{}
	case "Exceptions":
		return &ExceptionsAttribute{}
	case "LineNumberTable":
		return &LineNumberTableAttribute{}
	case "LocalVariableTable":
		return &LocalVariableTableAttribute{}
	case "SourceFile":
		return &SourceFileAttribute{cp: cp}
	case "Synthetic":
		return &SyntheticAttribute{}
	default:
		return &UnparsedAttribute{attrName, attrLen, nil}
	}
}
