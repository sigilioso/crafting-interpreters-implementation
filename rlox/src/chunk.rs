use crate::simple_vec::SimpleVec;

pub mod operation {
    pub const OP_RETURN: u8 = 0;
}

pub struct Chunk {
    code: SimpleVec<u8>,
}

impl Chunk {
    pub fn new() -> Self {
        Self {
            code: SimpleVec::new(),
        }
    }

    pub fn write(&mut self, op: u8) {
        self.code.push(op);
    }

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
        match op {
            operation::OP_RETURN => self.disassemble_simple_instruction("OP_RETURN", offset),
            _ => {
                println!("Unknown opcode {op}");
                offset + 1
            }
        }
    }

    fn disassemble_simple_instruction(&self, name: &str, offset: usize) -> usize {
        println!("{name}");
        offset + 1
    }
}

pub mod vec_impl {
    use crate::value::Value;
    pub enum Operation {
        Constant(usize),
        Return,
    }

    #[derive(Default)]
    pub struct Chunk {
        // Using a rust vectors for now
        pub code: Vec<Operation>,
        pub lines: Vec<usize>,
        pub constants: Vec<Value>,
    }

    impl Chunk {
        pub fn write(&mut self, op: Operation, line: usize) {
            self.code.push(op);
            self.lines.push(line);
        }

        pub fn add_constant(&mut self, v: Value) -> usize {
            let index = self.constants.len();
            self.constants.push(v);
            index
        }

        pub fn disassemble(&self, name: &str) {
            println!("== {name} ==");
            self.code
                .iter()
                .enumerate()
                .for_each(|(i, op)| self.disassemble_op(i, op));
        }

        fn disassemble_op(&self, i: usize, op: &Operation) {
            print!("{i:04}");
            if i > 0 && self.lines[i] == self.lines[i - 1] {
                print!("   | ")
            } else {
                print!("{:>4} ", self.lines[i])
            }
            match op {
                Operation::Constant(i) => {
                    println!("{:<16} {} '{}'", "OP_CONSTANT", i, self.constants[*i])
                }
                Operation::Return => println!("OP_RETURN"),
            }
        }
    }
}
