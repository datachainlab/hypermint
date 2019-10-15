extern crate proc_macro;

use crate::proc_macro::TokenStream;
use quote::quote;
use std::ops::Deref;
use syn::{parse_macro_input, parse_quote, FnArg, Ident, ItemFn, Pat, Stmt, Type};

fn get_assignment_from_name_and_type(
    name: &Ident,
    ty: &Type,
    index: usize,
) -> Result<Stmt, String> {
    match ty {
        Type::Path(type_path) => {
            let pair = type_path.path.segments.last().ok_or("internal error")?;
            let ident = &pair.ident;
            let a = parse_quote! {
                let #name: #ident = get_arg(#index).unwrap();
            };
            Ok(a)
        }
        _ => Err("invalid arg type".to_string()),
    }
}

fn get_assignment_from_arg(arg: &FnArg, index: usize) -> Result<(&Ident, Stmt), String> {
    match arg {
        FnArg::Typed(pat_type) => match &pat_type.pat.deref() {
            Pat::Ident(pat) => {
                let pat = &pat.ident;
                match get_assignment_from_name_and_type(pat, pat_type.ty.deref(), index) {
                    Ok(stmt) => Ok((pat, stmt)),
                    Err(err) => Err(err),
                }
            }
            _ => Err("no parameter name".to_string()),
        },
        FnArg::Receiver(_receiver) => Err("receiver not supported".to_string()),
    }
}

#[proc_macro_attribute]
pub fn contract(_attr: TokenStream, item: TokenStream) -> TokenStream {
    let mut ast = parse_macro_input!(item as ItemFn);
    let org_name = &ast.sig.ident;
    let export_name = format!("{}", org_name);
    let fn_name = syn::Ident::new(format!("__{}", org_name).as_str(), org_name.span());

    let decl = &ast.sig;
    let inputs = &decl.inputs;
    let mut c = 0;
    let mut assignments = Vec::new();
    let mut arguments = Vec::new();
    for arg in inputs {
        let a = get_assignment_from_arg(&arg, c);
        c += 1;
        match a {
            Ok((ident, stmt)) => {
                assignments.push(stmt);
                arguments.push(ident);
            }
            Err(e) => {
                eprintln!("{:?}", e);
            }
        }
    }

    let pre = quote! {
        #[cfg_attr(not(feature = "emulation"), export_name = #export_name)]
        pub fn #fn_name() -> i32 {
            use hmcdk::prelude::*;
            use hmcdk::api::{return_value, log};
            #(#assignments);*
            match #org_name(#(#arguments),*) {
                Ok(Some(v)) => {
                    match return_value(&v.to_bytes()) {
                        0 => 0,
                        _ => -1,
                    }
                },
                Ok(None) => 0,
                Err(e) => match return_value(&format!("{:?}", e).as_str().as_bytes()) {
                    _ => -1,
                }
            }
        }
    };
    let pre: TokenStream = pre.into();

    let mut stmts = ast.block.stmts.clone();

    ast.block.stmts.clear();
    ast.block.stmts.append(&mut stmts);

    let gen = quote! {
        #ast
    };
    let t: TokenStream = gen.into();
    format!("{} {}", pre.to_string(), t.to_string())
        .parse::<TokenStream>()
        .unwrap()
}
