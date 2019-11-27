#[macro_use]
extern crate serde;
extern crate serde_json;

use crate::Error::{ArgError, InternalError};
use std::env;
use std::fmt::{self, Display};
use std::fs;
use std::io::{self, Write};
use std::path::PathBuf;
use std::process;
use syn::visit::Visit;
use syn::File;

use abi::types;

mod argument;
mod function;

enum Error {
    IncorrectUsage,
    ArgError(argument::ArgumentError),
    ReadFile(io::Error),
    InternalError,
}

impl Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        use self::Error::*;
        use argument::ArgumentError::*;

        match self {
            IncorrectUsage => write!(f, "Usage: dump-syntax path/to/filename.rs"),
            ArgError(UnrecognizedType) => write!(f, "unrecognized arg type"),
            ArgError(UnsupportedType(x)) => write!(f, "unsupported arg type: {}", x),
            ArgError(InvalidFunction(x)) => write!(f, "invalid function: {}", x),
            ReadFile(error) => write!(f, "Unable to read file: {}", error),
            InternalError => write!(f, "internal error"),
        }
    }
}

fn main() {
    if let Err(error) = try_main() {
        let _ = writeln!(io::stderr(), "{}", error);
        process::exit(1);
    }
}

fn try_main() -> Result<(), Error> {
    let mut args = env::args_os();
    let _ = args.next(); // executable name

    let filepath = match (args.next(), args.next()) {
        (Some(arg), None) => PathBuf::from(arg),
        _ => return Err(Error::IncorrectUsage),
    };

    let code = fs::read_to_string(&filepath).map_err(Error::ReadFile)?;
    let syntax_tree: File = syn::parse_file(&code).unwrap();
    let mut visitor = function::FnVisitor {
        functions: Vec::new(),
        error: None,
    };
    visitor.visit_file(&syntax_tree);
    match visitor.error.map(|e| ArgError(e)) {
        Some(e) => return Err(e),
        None => {}
    }
    let json = serde_json::to_string(&visitor.functions).map_err(|_x| InternalError)?;
    println!("{}", json);
    Ok(())
}
