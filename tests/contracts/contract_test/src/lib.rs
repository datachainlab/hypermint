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

#[contract(readonly)]
pub fn test_get_sender() -> R<Address> {
    Ok(Some(get_sender()?))
}

#[contract(readonly)]
pub fn test_get_contract_address() -> R<Address> {
    Ok(Some(get_contract_address()?))
}

#[contract(readonly)]
pub fn test_get_arguments(arg_idx: i32) -> R<Vec<u8>> {
    let arg: Vec<u8> = get_arg(arg_idx as usize)?;
    Ok(Some(arg))
}

#[contract(readonly)]
pub fn check_signature(msg_hash: Vec<u8>, sig: Vec<u8>) -> R<i32> {
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
pub fn test_write_state(key: Vec<u8>, value: Vec<u8>) -> R<i32> {
    write_state(&key, &value);
    Ok(None)
}

#[contract(readonly)]
pub fn test_read_state(key: Vec<u8>) -> R<Vec<u8>> {
    let value: Vec<u8> = match read_state(&key) {
        Ok(v) => v,
        Err(_) => vec![],
    };
    Ok(Some(value))
}

#[contract]
pub fn test_read_write_state(key: Vec<u8>, value: Vec<u8>) -> R<i32> {
    // read a value, but nop
    read_state::<Vec<u8>>(&key);

    write_state(&key, &value);
    Ok(None)
}

#[contract]
pub fn test_write_to_same_key(key: Vec<u8>, value1: Vec<u8>, value2: Vec<u8>) -> R<i32> {
    // read a value, but nop
    read_state::<Vec<u8>>(&key);

    write_state(&key, &value1);
    write_state(&key, &value2);

    Ok(None)
}

#[contract]
pub fn test_write_to_multiple_key(key1: Vec<u8>, value1: Vec<u8>, key2: Vec<u8>, value2: Vec<u8>) -> R<i32> {
    // read a value, but nop
    read_state::<Vec<u8>>(&key1);
    read_state::<Vec<u8>>(&key2);

    write_state(&key1, &value1);
    write_state(&key2, &value2);

    Ok(None)
}

#[contract(readonly)]
pub fn test_keccak256(msg: Vec<u8>) -> R<Vec<u8>> {
    Ok(Some(keccak256(&msg)?.to_vec()))
}

#[contract(readonly)]
pub fn test_sha256(msg: Vec<u8>) -> R<Vec<u8>> {
    Ok(Some(sha256(&msg)?.to_vec()))
}

#[contract]
pub fn test_emit_event(msg0: Vec<u8>, msg1: Vec<u8>) -> R<Vec<u8>> {
    let name0 = "test-event-name-0";
    let name1 = "test-event-name-1";

    emit_event(name0, &msg0)?;
    emit_event(name1, &msg1)?;
    Ok(None)
}

#[contract]
pub fn test_external_emit_event(org_msg: Vec<u8>, addr: Address, ext_msg: Vec<u8>) -> R<Vec<u8>> {
    let _: Vec<u8> = call_contract(&addr, "test_emit_event".as_bytes(), vec![&ext_msg])?;
    let name = "test-org-event-name";
    emit_event(name, &org_msg)?;
    Ok(None)
}

#[contract]
pub fn test_call_external_contract(addr: Address, x: Vec<u8>, y: Vec<u8>) -> R<Vec<u8>> {
    Ok(Some(call_contract(&addr, "test_plus".as_bytes(), vec![&x, &y])?))
}

#[contract]
pub fn test_call_who_am_i_on_external_contract(addr: Address) -> R<Vec<u8>> {
    Ok(Some(call_contract(&addr, "who_am_i".as_bytes(), vec![])?))
}

#[contract]
pub fn test_call_get_contract_address_on_external_contract(addr: Address) -> R<Vec<u8>> {
    Ok(Some(call_contract(&addr, "get_contract_address".as_bytes(), vec![])?))
}

#[contract]
pub fn init() -> R<i32> {
    Ok(None)
}
