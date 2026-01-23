use super::{Chunk, operation};

impl Chunk {
    pub fn disassemble(&self, name: &str) {
        println!("== {name} ==");

        let mut offset = 0;
        while offset < self.code.count() {
            offset = self.disassemble_instruction(offset)
        }
    }

    fn disassemble_instruction(&self, offset: usize) -> usize {
        print!("{offset:04} ");
        let op = self.code.get_value(offset);
        if offset > 0 && self.lines.get_value(offset) == self.lines.get_value(offset - 1) {
            print!("   | ")
        } else {
            print!("{:>4} ", self.lines.get_value(offset))
        }
        match op {
            operation::OP_RETURN => self.disassemble_simple_instruction("OP_RETURN", offset),
            operation::OP_CONSTANT => self.disassemble_constant_instruction("OP_CONSTANT", offset),
            _ => {
                println!("Unknown opcode {op}");
                offset + 1
            }
        }
    }

    fn disassemble_constant_instruction(&self, name: &str, offset: usize) -> usize {
        let constant_index = self.code.get_value(offset + 1);
        let value = self.constants.get_ref(constant_index.into());
        println!("{name:<16} {constant_index} '{value}'");
        offset + 2
    }

    fn disassemble_simple_instruction(&self, name: &str, offset: usize) -> usize {
        println!("{name}");
        offset + 1
    }
}
