use std::fmt::Display;

#[derive(Clone, Copy)]
pub enum Value {
    Number(f64),
}

impl Display for Value {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Value::Number(n) => write!(f, "{n}"),
        }
    }
}

impl Value {
    // TODO: temporal operation implementations for first stages
    pub fn negate(self) -> Self {
        match self {
            Value::Number(n) => Value::Number(-n),
        }
    }

    pub fn add(a: Self, b: Self) -> Self {
        match (a, b) {
            (Value::Number(a), Value::Number(b)) => Value::Number(a + b),
        }
    }

    pub fn subtract(a: Self, b: Self) -> Self {
        match (a, b) {
            (Value::Number(a), Value::Number(b)) => Value::Number(a - b),
        }
    }

    pub fn multiply(a: Self, b: Self) -> Self {
        match (a, b) {
            (Value::Number(a), Value::Number(b)) => Value::Number(a * b),
        }
    }

    pub fn divide(a: Self, b: Self) -> Self {
        match (a, b) {
            (Value::Number(a), Value::Number(b)) => Value::Number(a / b),
        }
    }
}
