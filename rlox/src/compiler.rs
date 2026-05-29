use std::cell::{RefCell, RefMut};

use crate::{
    chunk::Chunk,
    operation::Operation,
    scanner::{Scanner, Token, TokenType},
};

pub struct Compiler {
    parser: Parser,
    compiling_chunk: RefCell<Chunk>, // TODO: check if RefCell is the best approach here
}

impl Compiler {
    pub fn emit_byte(&self, byte: u8) {
        self.current_chunk().write(byte, self.parser.previous.line);
    }

    pub fn emit_bytes(&self, byte1: u8, byte2: u8) {
        self.emit_byte(byte1);
        self.emit_byte(byte2);
    }

    pub fn emit_return(&self) {
        self.emit_byte(Operation::Return.into());
    }

    pub fn current_chunk(&self) -> RefMut<'_, Chunk> {
        self.compiling_chunk.borrow_mut()
    }
}

pub struct Parser {
    current: Token,
    previous: Token,
    scanner: Scanner,
    had_error: bool,
    panic_mode: bool,
}

impl Parser {
    pub fn advance(&mut self) {
        self.previous = std::mem::replace(&mut self.current, self.scanner.scan_token());
        // advance reporting errors until a non-error token is found
        loop {
            if !matches!(self.current.token_type, TokenType::Error) {
                break;
            }
            let err_token = std::mem::replace(&mut self.current, self.scanner.scan_token());
            self.error_at(&err_token, &err_token.lexeme);
        }
    }

    pub fn consume(&mut self, token_type: TokenType, message: &[u8]) {
        if matches!(&self.current.token_type, token_type) {
            self.advance();
        } else {
            let err_token = self.current.clone();
            self.error_at(&err_token, message);
        }
    }

    pub fn error_at(&mut self, token: &Token, message: &[u8]) {
        if self.panic_mode {
            return;
        }
        self.panic_mode = true;
        eprint!("[{}] Error", token.line);
        match token.token_type {
            TokenType::Eof => eprint!(" at end"),
            TokenType::Error => {}
            _ => eprint!(" at '{}'", String::from_utf8_lossy(&token.lexeme)),
        }
        eprintln!(": {}", String::from_utf8_lossy(message));
        self.had_error = true;
    }
}
