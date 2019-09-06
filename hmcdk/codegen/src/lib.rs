extern crate proc_macro;

use crate::proc_macro::TokenStream;
use quote::quote;
use syn::{parse_macro_input, ItemFn};

#[proc_macro_attribute]
pub fn contract(_attr: TokenStream, item: TokenStream) -> TokenStream {
    let mut ast = parse_macro_input!(item as ItemFn);
    let org_name = &ast.ident;
    let export_name = format!("{}", org_name);
    let fn_name = syn::Ident::new(format!("__{}", org_name).as_str(), org_name.span());

    let pre = quote! {
        #[cfg_attr(not(feature = "emulation"), export_name = #export_name)]
        pub fn #fn_name() -> i32 {
            use hmcdk::prelude::*;
            use hmcdk::api::{return_value, log};
            match #org_name() {
                Ok(Some(v)) => {
                    match return_value(&v) {
                        0 => 0,
                        _ => -1,
                    }
                },
                Ok(None) => 0,
                Err(e) => {
                    log(&format!("{:?}", e).as_str().as_bytes());
                    -1
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
