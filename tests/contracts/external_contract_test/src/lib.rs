extern crate hmcdk;
use hmcdk::api;
use hmcdk::prelude::*;

#[contract]
pub fn test_plus() -> R<i64> {
    let x: i64 = api::get_arg(0)?;
    let y: i64 = api::get_arg(1)?;
    Ok(Some(x+y))
}

#[contract]
pub fn who_am_i() -> R<Address> {
    let sender = api::get_sender()?;
    Ok(Some(sender))
}

#[contract]
pub fn get_contract_address() -> R<Address> {
    let address = api::get_contract_address()?;
    Ok(Some(address))
}

#[contract]
pub fn init() -> R<i32> {
    Ok(None)
}
