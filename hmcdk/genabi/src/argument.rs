use syn::{GenericArgument, PathArguments, Type};

use super::types;
use std::fmt;
use std::fmt::Display;

pub enum ArgumentError {
    UnrecognizedType,
    UnsupportedType(String),
}

#[derive(Clone, Serialize, Deserialize, Debug)]
pub struct Argument {
    r#type: String,
    name: String,
}

impl Display for Argument {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.r#type)
    }
}

pub fn from_name_and_type(name: &String, ty: &Type) -> Result<Argument, ArgumentError> {
    match ty {
        Type::Path(type_path) => {
            let pair = type_path
                .path
                .segments
                .last()
                .ok_or(ArgumentError::UnrecognizedType)?;
            let ident = &pair.ident;
            match types::to_abi_primitive(&ident.to_string()) {
                Some(x) => {
                    return Ok(Argument {
                        r#type: types::get_primitive_name(x),
                        name: name.clone(),
                    })
                }
                None => match ident.to_string().as_str() {
                    "R" => {
                        return match &pair.arguments {
                            PathArguments::AngleBracketed(x) => match x.args.first() {
                                Some(GenericArgument::Type(ty)) => from_name_and_type(&name, &ty),
                                _ => Err(ArgumentError::UnsupportedType("R".to_string())),
                            },
                            _ => Err(ArgumentError::UnsupportedType("R".to_string())),
                        }
                    }
                    _ => return Err(ArgumentError::UnsupportedType(ident.to_string())),
                },
            }
        }
        _ => Err(ArgumentError::UnrecognizedType),
    }
}
