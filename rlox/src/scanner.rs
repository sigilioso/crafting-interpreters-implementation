pub struct Scanner {
    source: Vec<u8>,
    start: u64,
    current: u64,
    line: u64,
}

impl Scanner {
    pub fn new(source: Vec<u8>) -> Self {
        Self {
            source,
            start: 0,
            current: 0,
            line: 0,
        }
    }

    pub fn scan_token(&mut self) -> Token {
        self.skip_whitespace_and_comments();
        self.start = self.current;

        if self.is_at_end() {
            return self.build_token(TokenType::Eof);
        }

        let c = self.advance();

        match c {
            b'(' => self.build_token(TokenType::LeftParen),
            b')' => self.build_token(TokenType::RightParen),
            b'{' => self.build_token(TokenType::LeftBrace),
            b'}' => self.build_token(TokenType::RightBrace),
            b';' => self.build_token(TokenType::SemiColon),
            b',' => self.build_token(TokenType::Comma),
            b'.' => self.build_token(TokenType::Dot),
            b'-' => self.build_token(TokenType::Minus),
            b'+' => self.build_token(TokenType::Plus),
            b'/' => self.build_token(TokenType::Slash),
            b'*' => self.build_token(TokenType::Star),
            b'!' => self.matches_or(b'=', TokenType::BangEqual, TokenType::Bang),
            b'=' => self.matches_or(b'=', TokenType::EqualEqual, TokenType::Equal),
            b'<' => self.matches_or(b'=', TokenType::LessEqual, TokenType::Less),
            b'>' => self.matches_or(b'=', TokenType::GreaterEqual, TokenType::Greater),
            b'"' => self.string(),
            _ if c.is_ascii_alphabetic() => self.identifier(),
            _ if c.is_ascii_digit() => self.number(),
            _ => self.err_token("Unexpected character."),
        }
    }

    pub fn build_token(&self, token_type: TokenType) -> Token {
        Token {
            token_type,
            lexeme: self.source[self.start..self.current].to_vec(),
            line: self.line,
        }
    }

    pub fn err_token(&self, message: &str) -> Token {
        Token {
            token_type: TokenType::Error,
            lexeme: message.as_bytes().to_owned(),
            line: self.line,
        }
    }

    fn is_at_end(&self) -> bool {
        self.source[self.current] == b'\0'
    }

    fn skip_whitespace_and_comments(&mut self) {
        loop {
            match self.peek() {
                b' ' | b'\t' | b'r' => self.current += 1,
                b'\n' => {
                    self.line += 1;
                    self.current += 1
                }
                b'/' => {
                    if self.peek_next() == b'/' {
                        // Ignore the rest of the line
                        while (self.peek() != b'\n' && !self.is_at_end()) {
                            self.current += 1
                        }
                    }
                }
                _ => {
                    return;
                }
            }
        }
    }

    fn string(&mut self) -> Token {
        while self.peek() != b'"' && !self.is_at_end() {
            if self.peek() == b'\n' {
                self.line += 1;
            }
            self.current += 1;
        }
        if self.is_at_end() {
            self.err_token("Unterminated string.")
        } else {
            self.current += 1; // the closing quote
            self.build_token(TokenType::String)
        }
    }

    fn number(&mut self) -> Token {
        while self.peek().is_ascii_digit() {
            self.current += 1;
        }
        // fractional part
        if self.peek() == b'.' && self.peek_next().is_ascii_digit() {
            self.current += 1; // consume the dot
            while self.peek().is_ascii_digit() {
                self.current += 1;
            }
        }
        self.build_token(TokenType::Number)
    }

    fn identifier(&mut self) -> Token {
        while self.peek().is_ascii_alphanumeric() {
            self.current += 1
        }
        let identifier_type = self.identifier_type();
        self.build_token(identifier_type)
    }

    fn identifier_type(&self) -> TokenType {
        match self.source[self.start] {
            b'a' => self.check_keyword(1, 2, b"nd", TokenType::And),
            b'c' => self.check_keyword(1, 4, b"lass", TokenType::Class),
            b'e' => self.check_keyword(1, 3, b"lse", TokenType::Else),
            b'i' => self.check_keyword(1, 1, b"f", TokenType::If),
            b'n' => self.check_keyword(1, 2, b"il", TokenType::Nil),
            b'o' => self.check_keyword(1, 1, b"r", TokenType::Or),
            b'p' => self.check_keyword(1, 4, b"rint", TokenType::Print),
            b'r' => self.check_keyword(1, 5, b"eturn", TokenType::Return),
            b's' => self.check_keyword(1, 4, b"uper", TokenType::Super),
            b'v' => self.check_keyword(1, 2, b"ar", TokenType::Var),
            b'w' => self.check_keyword(1, 4, b"hile", TokenType::While),
            b'f' => {
                if self.current - self.start > 1 {
                    match self.source[self.start + 1] {
                        b'a' => self.check_keyword(2, 3, b"lse", TokenType::False),
                        b'o' => self.check_keyword(2, 1, b"r", TokenType::For),
                        b'u' => self.check_keyword(2, 1, b"n", TokenType::Fun),
                        _ => TokenType::Identifier,
                    }
                } else {
                    TokenType::Identifier
                }
            }
            b't' => {
                if self.current - self.start > 1 {
                    match self.source[self.current + 1] {
                        b'h' => self.check_keyword(2, 2, b"is", TokenType::This),
                        b'r' => self.check_keyword(2, 2, b"ue", TokenType::True),
                        _ => TokenType::Identifier,
                    }
                } else {
                    TokenType::Identifier
                }
            }
            _ => TokenType::Identifier,
        }
    }

    fn check_keyword(
        &self,
        start: u64,
        length: u64,
        rest: &[u8],
        token_type: TokenType,
    ) -> TokenType {
        if (self.current - self.start == start + length)
            && (self.source[self.start + start..self.start + start + length] == *rest)
        {
            token_type
        } else {
            TokenType::Identifier
        }
    }

    fn advance(&mut self) -> u8 {
        let c = self.source[self.current];
        self.current += 1;
        c
    }

    fn advance_if_matches(&mut self, expected: u8) -> bool {
        if self.is_at_end() || self.source[self.current] != expected {
            false
        } else {
            self.current += 1;
            true
        }
    }

    fn matches_or(&mut self, expected: u8, a: TokenType, b: TokenType) -> Token {
        let t = if self.advance_if_matches(expected) {
            a
        } else {
            b
        };
        self.build_token(t)
    }

    fn peek(&self) -> u8 {
        self.source[self.current]
    }

    fn peek_next(&self) -> u8 {
        if self.is_at_end() {
            b'\0'
        } else {
            self.source[self.current + 1]
        }
    }
}

#[derive(Debug, Clone)]
pub struct Token {
    pub token_type: TokenType,
    pub lexeme: Vec<u8>,
    pub line: u64,
}

#[derive(Debug, Clone)]
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
    And,
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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_identifiers() {
        let mut sc = Scanner {
            source: b"class ! >= and else if print someId42 fun(1+2.0)\0 ".to_vec(),
            start: 0,
            current: 0,
            line: 0,
        };
        assert!(matches!(sc.scan_token().token_type, TokenType::Class));
        assert!(matches!(sc.scan_token().token_type, TokenType::Bang));
        assert!(matches!(
            sc.scan_token().token_type,
            TokenType::GreaterEqual
        ));
        assert!(matches!(sc.scan_token().token_type, TokenType::And));
        assert!(matches!(sc.scan_token().token_type, TokenType::Else));
        assert!(matches!(sc.scan_token().token_type, TokenType::If));
        assert!(matches!(sc.scan_token().token_type, TokenType::Print));
        assert!(matches!(sc.scan_token().token_type, TokenType::Identifier));
        assert!(matches!(sc.scan_token().token_type, TokenType::Fun));
        assert!(matches!(sc.scan_token().token_type, TokenType::LeftParen));
        assert!(matches!(sc.scan_token().token_type, TokenType::Number));
        assert!(matches!(sc.scan_token().token_type, TokenType::Plus));
        assert!(matches!(sc.scan_token().token_type, TokenType::Number));
        assert!(matches!(sc.scan_token().token_type, TokenType::RightParen));
        assert!(matches!(sc.scan_token().token_type, TokenType::Eof));
    }
}
