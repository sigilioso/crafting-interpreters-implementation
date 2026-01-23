use crate::{simple_vec::SimpleVec, value::Value};

mod debug;

pub mod operation {
    pub const OP_RETURN: u8 = 0;
    pub const OP_CONSTANT: u8 = 1;
}

pub struct Chunk {
    code: SimpleVec<u8>,
    constants: SimpleVec<Value>,
    lines: SimpleVec<u64>,
}

impl Chunk {
    pub fn new() -> Self {
        Self {
            code: SimpleVec::new(),
            constants: SimpleVec::new(),
            lines: SimpleVec::new(),
        }
    }

    pub fn write(&mut self, op: u8, line: u64) {
        self.code.push(op);
        self.lines.push(line);
    }

    pub fn add_constant(&mut self, value: Value) -> u8 {
        let index = self.constants.count();
        self.constants.push(value);
        u8::try_from(index).expect("the number of constants in a chunk must fit in a single byte")
    }
}
