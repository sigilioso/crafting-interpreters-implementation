mod chunk;
mod value;

use chunk::{Chunk, Operation};
use value::Value;

use crate::chunk::unsafe_impl;

fn main() {
    let mut c = Chunk::default();
    let constant_index = c.add_constant(Value::Number(1.2));
    c.write(Operation::Constant(constant_index), 123);
    c.write(Operation::Return, 123);
    c.disassemble("test chunk");

    // unsafe version
    let mut c = unsafe_impl::Chunk::new();
    c.write(unsafe_impl::operation::OP_RETURN);
    c.disassemble("test unsafe chunk");
}
