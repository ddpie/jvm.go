package classfile

//import "errors"

/*
ClassFile {
    u4             magic;
    u2             minor_version;
    u2             major_version;
    u2             constant_pool_count;
    cp_info        constant_pool[constant_pool_count-1];
    u2             access_flags;
    u2             this_class;
    u2             super_class;
    u2             interfaces_count;
    u2             interfaces[interfaces_count];
    u2             fields_count;
    field_info     fields[fields_count];
    u2             methods_count;
    method_info    methods[methods_count];
    u2             attributes_count;
    attribute_info attributes[attributes_count];
}
*/
type ClassFile struct {
  //magic           uint32
    minorVersion    uint16
    majorVersion    uint16
    constantPool    *ConstantPool
    accessFlags     uint16
    thisClass       uint16
    superClass      uint16
    interfaces      []uint16
    fields          []*FieldInfo
    methods         []*MethodInfo
    AttributeTable
}

func (self *ClassFile) read(reader *ClassReader) {
    self.readAndCheckMagic(reader)
    self.readVersions(reader)
    self.readConstantPool(reader)
    self.accessFlags = reader.readUint16()
    self.thisClass = reader.readUint16()
    self.superClass = reader.readUint16()
    self.readInterfaces(reader)
    self.readFields(reader)
    self.readMethods(reader)
    self.attributes = readAttributes(reader, self.constantPool)
}

func (self *ClassFile) readAndCheckMagic(reader *ClassReader) {
    magic := reader.readUint32()
    if magic != 0xCAFEBABE {
        panic("Bad magic!")
    }
}

func (self *ClassFile) readVersions(reader *ClassReader) {
    self.minorVersion = reader.readUint16()
    self.majorVersion = reader.readUint16()
    // todo check versions
}

func (self *ClassFile) readConstantPool(reader *ClassReader) {
    self.constantPool = &ConstantPool{}
    self.constantPool.read(reader)
}

func (self *ClassFile) readInterfaces(reader *ClassReader) {
    interfacesCount := reader.readUint16()
    self.interfaces = make([]uint16, interfacesCount)
    for i := range self.interfaces {
        self.interfaces[i] = reader.readUint16()
    }
}

func (self *ClassFile) readFields(reader *ClassReader) {
    fieldsCount := reader.readUint16()
    self.fields = make([]*FieldInfo, fieldsCount)
    for i := range self.fields {
        self.fields[i] = &FieldInfo{}
        self.fields[i].cp = self.constantPool
        self.fields[i].read(reader)
    }
}

func (self *ClassFile) readMethods(reader *ClassReader) {
    methodsCount := reader.readUint16()
    self.methods = make([]*MethodInfo, methodsCount)
    for i := range self.methods {
        self.methods[i] = &MethodInfo{}
        self.methods[i].cp = self.constantPool
        self.methods[i].read(reader)
    }
}

func (self *ClassFile) ConstantPool() (*ConstantPool) {
    return self.constantPool
}
func (self *ClassFile) AccessFlags() (uint16) {
    return self.accessFlags
}
func (self *ClassFile) Fields() ([]*FieldInfo) {
    return self.fields
}
func (self *ClassFile) Methods() ([]*MethodInfo) {
    return self.methods
}

func (self *ClassFile) ClassName() (string) {
    return self.constantPool.getClassName(self.thisClass)
}

func (self *ClassFile) SuperClassName() (string) {
    if self.superClass != 0 {
        return self.constantPool.getClassName(self.superClass)
    } else {
        // todo Object
        return ""
    }
}

func (self *ClassFile) InterfaceNames() ([]string) {
    interfaceNames := make([]string, len(self.interfaces))
    for i, cpIndex := range self.interfaces {
        interfaceNames[i] = self.constantPool.getClassName(cpIndex)
    }
    return interfaceNames
}

func (self *ClassFile) FileName() string {
    sfAttr := self.SourceFileAttribute()
    if sfAttr != nil {
        return self.constantPool.getUtf8(sfAttr.sourceFileIndex)
    } else {
        return "Unknown" // todo
    }
}
