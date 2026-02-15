pub struct Scanner {
    start: usize,
    current: usize,
    line: usize,
}

pub struct Token {
    token_type: TokenType,
    lexeme: &'static str,
    line: usize,
}

impl Token {
    pub fn new(token_type: TokenType, line: usize) -> Self {
        Self {
            token_type,
            line,
            lexeme: "", // TODO
        }
    }

    fn example(source: &'static str) -> Token {
        let x = &source[..2];
        Token {
            token_type: TokenType::Add,
            line: 0,
            lexeme: x,
        }
    }
}

pub enum TokenType {
    // Single character
    LeftParen,
    RightParen,
    LeftBrace,
    RightBrace,
    Comma,
    Dot,
    Minus,
    Plus,
    SemiColon,
    Slash,
    Star,
    // One or two characters
    Bang,
    BangEqual,
    Equal,
    EqualEqual,
    Greater,
    GreaterEqual,
    Less,
    LessEqual,
    // Literals
    Identifier,
    String,
    Number,
    // Keywords
    Add,
    Class,
    Else,
    False,
    For,
    Fun,
    If,
    Nil,
    Or,
    Print,
    Return,
    Super,
    This,
    True,
    Var,
    While,

    Error,
    Eof,
}
