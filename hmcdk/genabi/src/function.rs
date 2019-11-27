use std::ops::Deref;
use syn::visit::{self, Visit};
use syn::{FnArg, ItemFn, Pat, ReturnType};

use super::argument;
use super::argument::Argument;
use super::argument::ArgumentError;

#[derive(Clone, Serialize, Deserialize, Debug)]
pub struct Function {
    r#type: String,
    name: String,
    simulate: bool,
    inputs: Vec<Argument>,
    outputs: Vec<Argument>,
}

pub struct FnVisitor {
    pub functions: Vec<Function>,
    pub error: Option<ArgumentError>,
}

impl<'ast> Visit<'ast> for FnVisitor {
    fn visit_item_fn(&mut self, node: &'ast ItemFn) {
        match from_item_fn(node) {
            Ok(Some(f)) => self.functions.push(f),
            Ok(None) => {}
            Err(x) => self.error = Some(x),
        }
        visit::visit_item_fn(self, node);
    }
}

fn get_argument_from_arg(arg: &FnArg) -> Result<Argument, ArgumentError> {
    match arg {
        FnArg::Typed(pat_type) => match pat_type.pat.deref() {
            Pat::Ident(pat) => {
                let pat = &pat.ident;
                return argument::from_name_and_type(&pat.to_string(), &pat_type.ty);
            }
            _ => Err(ArgumentError::UnrecognizedType),
        },
        FnArg::Receiver(_receiver) => Err(ArgumentError::UnrecognizedType),
    }
}

fn get_return_type(rtype: &ReturnType) -> Result<Option<Argument>, ArgumentError> {
    match rtype {
        ReturnType::Default => Ok(None),
        ReturnType::Type(_token, b) => {
            match argument::from_name_and_type(&"".to_string(), b.deref()) {
                Ok(a) => Ok(Some(a)),
                Err(e) => Err(e),
            }
        }
    }
}

fn from_item_fn(f: &ItemFn) -> Result<Option<Function>, ArgumentError> {
    let mut is_contract = false;
    let mut is_readonly = false;
    let fn_name = f.sig.ident.to_string();
    for a in &f.attrs {
        if a.path.segments.len() == 1 {
            let s = a.path.segments.last();
            match s {
                Some(s) => {
                    if s.ident == "contract" {
                        is_contract = true;
                        let ident: syn::Ident = match a.parse_args() {
                            Ok(ident) => ident,
                            Err(err) => return Err(ArgumentError::InvalidFunction(err.to_string()))
                        };
                        if ident == "readonly" {
                            is_readonly = true
                        }
                    }
                }
                None => {}
            }
        }
    }
    let func = if is_contract {
        let maybe_args = f
            .sig
            .inputs
            .iter()
            .map(|input| get_argument_from_arg(input));
        let rt = get_return_type(&f.sig.output)?;
        let args = maybe_args.filter_map(Result::ok).collect::<Vec<Argument>>();
        Some(Function {
            r#type: "function".to_string(),
            name: fn_name,
            simulate: is_readonly,
            inputs: args,
            outputs: match rt {
                Some(rt) => [rt].to_vec(),
                None => [].to_vec(),
            },
        })
    } else {
        None
    };
    Ok(func)
}
