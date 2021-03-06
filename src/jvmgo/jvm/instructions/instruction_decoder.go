package instructions

type InstructionDecoder struct {
    pc      int
    code    []byte // bytecodes
}

func NewInstructionDecoder() (*InstructionDecoder) {
    return &InstructionDecoder{}
}

func (self *InstructionDecoder) Decode(code []byte, pc int) (uint8, Instruction, int) {
    self.code = code
    self.pc = pc

    opcode := self.readUint8()
    instruction := newInstruction(opcode)
    instruction.fetchOperands(self)
    nextPC := self.pc

    return opcode, instruction, nextPC
}

func (self *InstructionDecoder) readInt8() (int8) {
    return int8(self.readUint8())
}
func (self *InstructionDecoder) readUint8() (uint8) {
    i := self.code[self.pc]
    self.pc++
    return i
}

func (self *InstructionDecoder) readInt16() (int16) {
    return int16(self.readUint16())
}
func (self *InstructionDecoder) readUint16() (uint16) {
    byte1 := uint16(self.readUint8())
    byte2 := uint16(self.readUint8())
    return (byte1 << 8) | byte2
}

func (self *InstructionDecoder) readInt32() (int32) {
    return int32(self.readUint32())
}
func (self *InstructionDecoder) readUint32() (uint32) {
    byte1 := uint32(self.readUint8())
    byte2 := uint32(self.readUint8())
    byte3 := uint32(self.readUint8())
    byte4 := uint32(self.readUint8())
    return (byte1 << 24) | (byte2 << 16) | (byte3 << 8) | byte4
}

func (self *InstructionDecoder) readInt32s(count int32) ([]int32) {
    ints := make([]int32, count)
    for i := range ints {
        ints[i] = self.readInt32()
    }
    return ints
}
