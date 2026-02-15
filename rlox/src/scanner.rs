pub struct Scanner<'a> {
    source: &'a [char],
    start: usize,
    current: usize,
    line: usize,
}

impl<'a> Scanner<'a> {
    pub fn token(&self, token_type: TokenType) -> Token<'a> {
        Token {
            token_type,
            lexeme: &self.source[self.start..self.current],
            line: self.line,
        }
    }

    pub fn err_token(&self, message: &'a [char]) -> Token<'a> {
        Token {
            token_type: TokenType::Error,
            lexeme: message,
            line: self.line,
        }
    }

    pub fn is_at_end(&self) -> bool {
        self.source[self.current] == '\0'
    }
}

pub struct Token<'a> {
    token_type: TokenType,
    lexeme: &'a [char],
    line: usize,
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
