mod chunk;
mod operation;
mod simple_vec;
mod value;
mod vm;

use chunk::Chunk;
use operation::Operation;
use value::Value;

use crate::vm::VM;

fn main() {
    let mut vm = VM::new();

    let mut c = Chunk::new();

    let constant_index = c.add_constant(Value::Number(1.2));
    c.write(Operation::Constant.into(), 123);
    c.write(constant_index, 123);

    let constant_index = c.add_constant(Value::Number(3.4));
    c.write(Operation::Constant.into(), 123);
    c.write(constant_index, 123);

    c.write(Operation::Add.into(), 123);

    let constant_index = c.add_constant(Value::Number(5.6));
    c.write(Operation::Constant.into(), 123);
    c.write(constant_index, 123);

    c.write(Operation::Divide.into(), 123);

    c.write(Operation::Negate.into(), 123);
    c.write(Operation::Return.into(), 123);
    c.disassemble("test chunk");

    let _ = vm.interpret(&c);
}
