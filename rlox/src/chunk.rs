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

pub mod unsafe_impl {
    use std::{
        alloc::{Layout, alloc, dealloc, handle_alloc_error},
        ptr,
    };

    pub mod operation {
        pub const OP_RETURN: u8 = 0;
    }

    fn grow_capacity(current: usize) -> usize {
        if current < 8 { 8 } else { current * 2 }
    }

    pub struct Chunk {
        code: *mut u8,
        capacity: usize,
        count: usize,
    }

    impl Chunk {
        pub fn new() -> Self {
            Self {
                code: ptr::null_mut(),
                capacity: 0,
                count: 0,
            }
        }

        pub fn write(&mut self, op: u8) {
            if self.capacity < self.count + 1 {
                self.capacity = grow_capacity(self.capacity);

                // TODO: extract to a generic helper
                let layout = Layout::array::<u8>(self.capacity).expect("failure reserving memory");
                self.code = unsafe { alloc(layout) };
                if self.code.is_null() {
                    handle_alloc_error(layout)
                }
            }
            unsafe {
                ptr::write(self.code.add(self.count), op);
            }
            self.count += 1
        }

        pub fn disassemble(&self, name: &str) {
            println!("== {name} ==");

            let mut offset = 0;
            while offset < self.count {
                offset = self.disassemble_instruction(offset)
            }
        }

        fn disassemble_instruction(&self, offset: usize) -> usize {
            print!("{offset:04} ");
            let op = unsafe { ptr::read(self.code.add(offset)) };
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

    impl Drop for Chunk {
        /// dealloc on drop
        fn drop(&mut self) {
            // TODO: also extract to a generic helper
            if !self.code.is_null() && self.capacity > 0 {
                let layout = Layout::array::<u8>(self.capacity).expect("failure freeing memory");
                unsafe {
                    dealloc(self.code, layout);
                }
            }
        }
    }
}
