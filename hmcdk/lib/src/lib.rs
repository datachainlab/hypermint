pub mod api;
pub mod error;
pub mod types;
pub mod utils;

pub mod prelude {
    pub use crate::error::Error;
    pub use crate::types::{Address, FromBytes, ToBytes, R};
    pub use hmcdk_codegen::*;
}
