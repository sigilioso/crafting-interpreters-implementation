mod chunk;
mod simple_vec;
mod value;

use chunk::{Chunk, operation};
use value::Value;

fn main() {
    let mut c = Chunk::new();
    let constant_index = c.add_constant(Value::Number(1.2));
    c.write(operation::OP_CONSTANT, 123);
    c.write(constant_index, 123);
    c.write(operation::OP_RETURN, 123);
    c.disassemble("test chunk");
}
