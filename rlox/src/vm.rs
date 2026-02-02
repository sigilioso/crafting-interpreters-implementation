use crate::{chunk::Chunk, operation::Operation, simple_vec::SimpleVec, value::Value};

const STACK_MAX: usize = 265;

pub enum InterpretError {
    CompileError,
    RuntimeError,
}

type InterpretResult = Result<(), InterpretError>;

pub struct VM {
    stack: [Value; STACK_MAX],
    stack_top: usize,
}

// TODO: impl drop when needed
impl VM {
    pub fn new() -> Self {
        Self {
            stack: [Value::Number(0.0); STACK_MAX], // TODO: better initial value maybe?
            stack_top: 0,
        }
    }

    pub fn interpret(&mut self, c: &Chunk) -> InterpretResult {
        let ip: usize = 0;
        self.run(c, ip)
    }

    pub fn run(&mut self, c: &Chunk, ip: usize) -> InterpretResult {
        let mut ip = ip;
        loop {
            if debug_stack_trace() {
                self.debug_stack();
                c.disassemble_instruction(ip);
            }
            match Operation::try_from(c.instruction(ip)).expect("Unknown opcode in VM") {
                Operation::Return => {
                    println!("{}", self.pop());
                    return Ok(());
                }
                Operation::Constant => {
                    let index = c.instruction(ip + 1);
                    ip += 1;
                    let value = c.constant(index as usize);
                    self.push(value);
                }
                Operation::Negate => {
                    let v = self.pop();
                    self.push(v.negate());
                }
                Operation::Add => self.binary_operation(Value::add),
                Operation::Subtract => self.binary_operation(Value::subtract),
                Operation::Multiply => self.binary_operation(Value::multiply),
                Operation::Divide => self.binary_operation(Value::divide),
            }
            ip += 1;
        }
    }

    fn pop(&mut self) -> Value {
        self.stack_top -= 1;
        self.stack[self.stack_top]
    }

    fn push(&mut self, v: Value) {
        self.stack[self.stack_top] = v;
        self.stack_top += 1
    }

    // TODO: this might not be enough for other operations or when the Value types changes
    fn binary_operation<F>(&mut self, f: F)
    where
        F: Fn(Value, Value) -> Value,
    {
        let a = self.pop();
        let b = self.pop();
        self.push(f(a, b));
    }

    fn debug_stack(&self) {
        print!("          ");
        for i in 0..self.stack_top {
            print!("[{}]", self.stack[i])
        }
        println!();
    }
}

// TODO: improve
pub fn debug_stack_trace() -> bool {
    std::env::var("DEBUG_STACK_TRACE").is_ok()
}
