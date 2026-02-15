mod chunk;
mod compiler;
mod operation;
mod scanner;
mod simple_vec;
mod value;
mod vm;

use std::{fs::read_to_string, io::stdin, process::exit};

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

fn repl() {
    loop {
        print!("> ");
        let mut line = String::new();
        stdin().read_line(&mut line).unwrap_or_else(|err| {
            eprintln!("error reading input: {err}");
            exit(74);
        });
        interpret(line);
    }
}

fn run_file(path: &str) {
    let source = read_to_string(path).unwrap_or_else(|err| {
        eprintln!("Could not read file '{path}': {err}");
        exit(74);
    });
    interpret(source);
}

fn interpret(source: String) {
    todo!("to be implemented in the VM")
}
