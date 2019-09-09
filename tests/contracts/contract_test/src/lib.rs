extern crate hmcdk;
use hmcdk::api::{emit_event, get_arg, get_sender, read_state, write_state, ecrecover_address, get_contract_address, keccak256, sha256, call_contract};
use hmcdk::error;
use hmcdk::prelude::*;

fn call_check_signature() -> Result<i32, Error> {
    let sender = get_sender()?;
    let msg_hash: Vec<u8> = get_arg(0)?;
    let sig: Vec<u8> = get_arg(1)?;
    let addr: Address = ecrecover_address(&msg_hash, &sig[64..65], &sig[0..32], &sig[32..64])?;

    if sender == addr {
        Ok(0)
    } else {
        Err(error::from_str("invalid signer"))
    }
}

#[contract]
pub fn test_get_sender() -> R<Address> {
    Ok(Some(get_sender()?))
}

#[contract]
pub fn test_get_contract_address() -> R<Address> {
    Ok(Some(get_contract_address()?))
}

#[contract]
pub fn test_get_arguments() -> R<Vec<u8>> {
    let argIdx: i32 = get_arg(0)?;
    let arg: Vec<u8> = get_arg(argIdx as usize)?;
    Ok(Some(arg))
}

#[contract]
pub fn check_signature() -> R<i32> {
    Ok(Some(call_check_signature()?))
}

#[contract]
pub fn test_read_uncommitted_state() -> R<i32> {
    let b = [0u8; 255];
    write_state("key".as_bytes(), &b);
    let v: Vec<u8> = read_state("key".as_bytes())?;
    let r: &[u8] = &b;
    if v == r {
        Ok(None)
    } else {
        Err(error::from_str("not match"))
    }
}

#[contract]
pub fn test_write_state() -> R<i32> {
    let key: Vec<u8> = get_arg(0)?;
    let value: Vec<u8> = get_arg(1)?;
    write_state(&key, &value);
    Ok(None)
}

#[contract]
pub fn test_read_state() -> R<Vec<u8>> {
    let key: Vec<u8> = get_arg(0)?;
    let value: Vec<u8> = match read_state(&key) {
        Ok(v) => v,
        Err(_) => vec![],
    };
    Ok(Some(value))
}

#[contract]
pub fn test_read_write_state() -> R<i32> {
    let key: Vec<u8> = get_arg(0)?;
    let value: Vec<u8> = get_arg(1)?;

    // read a value, but nop
    read_state::<Vec<u8>>(&key);

    write_state(&key, &value);
    Ok(None)
}

#[contract]
pub fn test_write_to_same_key() -> R<i32> {
    let key: Vec<u8> = get_arg(0)?;
    let value1: Vec<u8> = get_arg(1)?;
    let value2: Vec<u8> = get_arg(2)?;

    // read a value, but nop
    read_state::<Vec<u8>>(&key);

    write_state(&key, &value1);
    write_state(&key, &value2);

    Ok(None)
}

#[contract]
pub fn test_write_to_multiple_key() -> R<i32> {
    let key1: Vec<u8> = get_arg(0)?;
    let value1: Vec<u8> = get_arg(1)?;
    let key2: Vec<u8> = get_arg(2)?;
    let value2: Vec<u8> = get_arg(3)?;

    // read a value, but nop
    read_state::<Vec<u8>>(&key1);
    read_state::<Vec<u8>>(&key2);

    write_state(&key1, &value1);
    write_state(&key2, &value2);

    Ok(None)
}

#[contract]
pub fn test_keccak256() -> R<Vec<u8>> {
    let msg: Vec<u8> = get_arg(0)?;
    Ok(Some(keccak256(&msg)?.to_vec()))
}

#[contract]
pub fn test_sha256() -> R<Vec<u8>> {
    let msg: Vec<u8> = get_arg(0)?;
    Ok(Some(sha256(&msg)?.to_vec()))
}

#[contract]
pub fn test_emit_event() -> R<Vec<u8>> {
    let msg0: Vec<u8> = get_arg(0)?;
    let msg1: Vec<u8> = get_arg(1)?;
    let name0 = "test-event-name-0";
    let name1 = "test-event-name-1";

    emit_event(name0, &msg0)?;
    emit_event(name1, &msg1)?;
    Ok(None)
}

#[contract]
pub fn test_call_external_contract() -> R<Vec<u8>> {
    let addr: Address = get_arg(0)?;
    let x: Vec<u8> = get_arg(1)?;
    let y: Vec<u8> = get_arg(2)?;
    Ok(Some(call_contract(&addr, "test_plus".as_bytes(), vec![&x, &y])?))
}

#[contract]
pub fn test_call_who_am_i_on_external_contract() -> R<Vec<u8>> {
    let addr: Address = get_arg(0)?;
    Ok(Some(call_contract(&addr, "who_am_i".as_bytes(), vec![])?))
}

#[contract]
pub fn test_call_get_contract_address_on_external_contract() -> R<Vec<u8>> {
    let addr: Address = get_arg(0)?;
    Ok(Some(call_contract(&addr, "get_contract_address".as_bytes(), vec![])?))
}

#[contract]
pub fn init() -> R<i32> {
    Ok(None)
}
