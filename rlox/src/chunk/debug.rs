use super::super::Operation;
use super::Chunk;

impl Chunk {
    pub fn disassemble(&self, name: &str) {
        println!("== {name} ==");

        let mut offset = 0;
        while offset < self.code.count() {
            offset = self.disassemble_instruction(offset)
        }
    }

    pub fn disassemble_instruction(&self, offset: usize) -> usize {
        print!("{offset:04} ");
        let op = self.code.get_value(offset);
        if offset > 0 && self.lines.get_value(offset) == self.lines.get_value(offset - 1) {
            print!("   | ")
        } else {
            print!("{:>4} ", self.lines.get_value(offset))
        }
        let Ok(op) = Operation::try_from(op).inspect_err(|err| {
            println!("Unknown opcode {op}");
        }) else {
            return offset + 1;
        };
        match op {
            Operation::Return => self.disassemble_simple_instruction("OP_RETURN", offset),
            Operation::Constant => self.disassemble_constant_instruction("OP_CONSTANT", offset),
            Operation::Negate => self.disassemble_simple_instruction("OP_NEGATE", offset),
            Operation::Add => self.disassemble_simple_instruction("OP_ADD", offset),
            Operation::Subtract => self.disassemble_simple_instruction("OP_SUBTRACT", offset),
            Operation::Multiply => self.disassemble_simple_instruction("OP_MULTIPLY", offset),
            Operation::Divide => self.disassemble_simple_instruction("OP_DIVIDE", offset),
        }
    }

    fn disassemble_constant_instruction(&self, name: &str, offset: usize) -> usize {
        let constant_index = self.code.get_value(offset + 1);
        let value = self.constants.get_value(constant_index.into());
        println!("{name:<16} {constant_index} '{value}'");
        offset + 2
    }

    fn disassemble_simple_instruction(&self, name: &str, offset: usize) -> usize {
        println!("{name}");
        offset + 1
    }
}
