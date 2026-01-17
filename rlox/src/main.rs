mod chunk;
mod simple_vec;
mod value;

use chunk::{Chunk, operation};
use value::Value;

use crate::chunk::vec_impl;

fn main() {
    // *mut u8 version
    let mut c = Chunk::new();
    c.write(operation::OP_RETURN);
    c.disassemble("test chunk");

    // vec version
    let mut c = vec_impl::Chunk::default();
    let constant_index = c.add_constant(Value::Number(1.2));
    c.write(vec_impl::Operation::Constant(constant_index), 123);
    c.write(vec_impl::Operation::Return, 123);
    c.disassemble("test chunk");
}
