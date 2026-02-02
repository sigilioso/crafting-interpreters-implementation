/// Helper macro to implement the Operations enum
macro_rules! byte_enum {
    ($name:ident { $($variant:ident = $val:expr),* $(,)? }) => {
        #[repr(u8)]
        #[derive(Debug, PartialEq, Copy, Clone)]
        pub enum $name {
            $($variant = $val),*
        }

        impl TryFrom<u8> for $name {
            type Error = u8;

            fn try_from(v: u8) -> Result<Self, Self::Error> {
                match v {
                    $($val => Ok($name::$variant),)*
                    _ => Err(v),
                }
            }
        }

        impl From<$name> for u8 {
            fn from(v: $name) -> Self {
                v as u8
            }
        }
    };
}

byte_enum!(
    Operation {
        Return = 0,
        Constant = 1,
        Negate = 2,
        Add = 3,
        Subtract = 4,
        Multiply = 5,
        Divide = 6,
    }
);
