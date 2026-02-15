mod chunk;
mod compiler;
mod operation;
mod scanner;
mod simple_vec;
mod value;
mod vm;

use std::{
    env,
    fs::{File, read_to_string},
    io::{BufRead, Read, stdin},
    process::exit,
};

use chunk::Chunk;
use operation::Operation;
use value::Value;

use crate::{
    scanner::{Scanner, TokenType},
    vm::VM,
};

fn main() {
    let args = env::args();
    match args.len() {
        1 => repl(),
        2 => run_file(args.into_iter().nth(1).unwrap().as_str()),
        _ => {
            eprintln!("usage: rlox [path]");
            exit(64);
        }
    }
}

fn repl() {
    loop {
        print!("> ");
        interpret(read_stdin_line_bytes());
    }
}

fn read_stdin_line_bytes() -> Vec<u8> {
    let mut line: Vec<u8> = Vec::new();
    let mut stdin_handle = stdin().lock();
    stdin_handle
        .read_until(b'\n', &mut line)
        .unwrap_or_else(|err| {
            eprintln!("error reading input: {err}");
            exit(74);
        });
    line
}

fn run_file(path: &str) {
    let mut f = File::open(path).unwrap_or_else(|err| {
        eprintln!("Could not read file '{path}': {err}");
        exit(74);
    });
    let mut source: Vec<u8> = Vec::new();
    f.read_to_end(&mut source).unwrap_or_else(|err| {
        eprintln!("Could not read file '{path}': {err}");
        exit(74);
    });
    interpret(source);
}

// TODO: wire it up properly
fn interpret(source: Vec<u8>) {
    let mut scanner = Scanner::new(source);
    let mut line: usize = 0;
    let mut first = true;
    loop {
        let token = scanner.scan_token();
        if token.line != line || first {
            first = false;
            print!("{line}");
        } else {
            print!("   | ");
        }
        println!(
            "{:?} {}",
            token.token_type,
            String::from_utf8(token.lexeme).unwrap()
        );
        if matches!(token.token_type, TokenType::Eof) {
            break;
        }
    }
}

// TODO: remove
fn _old_wire_up() {
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
