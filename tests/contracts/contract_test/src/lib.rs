extern crate hmc;

fn call_check_signature() -> Result<i32, String> {
    let sender = hmc::get_sender()?;
    let msg_hash = hmc::get_arg(0)?;
    let sig = hmc::get_arg(1)?;
    let addr = hmc::ecrecover_address(&msg_hash, &sig[64..65], &sig[0..32], &sig[32..64])?;

    if sender == addr {
        Ok(0)
    } else {
        Err(format!("invalid signer"))
    }
}

#[no_mangle]
pub fn check_signature() -> i32 {
    match call_check_signature() {
        Ok(v) => v,
        Err(e) => {
            hmc::log(e.as_bytes());
            -1
        }
    }
}

#[no_mangle]
pub fn test_read_uncommitted_state() -> i32 {
    let b = [0u8; 255];
    hmc::write_state("key".as_bytes(), &b);
    match hmc::read_state("key".as_bytes()) {
        Ok(v) => {
            let r: &[u8] = &b;
            if v == r {
                0
            } else {
                hmc::log("not match".as_bytes());
                -1
            }
        }
        Err(e) => {
            hmc::log(e.as_bytes());
            -1
        }
    }
}

#[no_mangle]
pub fn test_write_state() -> i32 {
    let key = hmc::get_arg(0).unwrap();
    let value = hmc::get_arg(1).unwrap();
    hmc::write_state(&key, &value);
    0
}

#[no_mangle]
pub fn test_read_state() -> i32 {
    let key = hmc::get_arg(0).unwrap();
    let value = match hmc::read_state(&key) {
        Ok(v) => v,
        Err(_) => vec![],
    };
    hmc::return_value(&value)
}

#[no_mangle]
pub fn test_read_write_state() -> i32 {
    let key = hmc::get_arg(0).unwrap();
    let value = hmc::get_arg(1).unwrap();

    // read a value, but nop
    match hmc::read_state(&key) {
        Ok(v) => v,
        Err(_) => vec![],
    };

    hmc::write_state(&key, &value);
    0
}

#[no_mangle]
pub fn test_keccak256() -> i32 {
    let msg = hmc::get_arg(0).unwrap();
    let h = hmc::keccak256(&msg).unwrap();
    hmc::return_value(&h)
}

#[no_mangle]
pub fn test_sha256() -> i32 {
    let msg = hmc::get_arg(0).unwrap();
    let h = hmc::sha256(&msg).unwrap();
    hmc::return_value(&h)
}

#[no_mangle]
pub fn test_emit_event() -> i32 {
    let msg0 = hmc::get_arg(0).unwrap();
    let msg1 = hmc::get_arg(1).unwrap();
    let name0 = "test-event-name-0";
    let name1 = "test-event-name-1";

    hmc::emit_event(&name0, &msg0).unwrap();
    hmc::emit_event(&name1, &msg1).unwrap();
    0
}

#[no_mangle]
pub fn test_call_external_contract() -> i32 {
    let addr = hmc::hex_to_bytes(hmc::get_arg_str(0).unwrap().as_ref());
    let x = hmc::get_arg(1).unwrap();
    let y = hmc::get_arg(2).unwrap();
    match hmc::call_contract(&addr, "test_plus".as_bytes(), vec![&x, &y]) {
        Ok(v) => {
            hmc::return_value(&v)
        },
        Err(m) => {
            hmc::revert(format!("{}", m));
            -1
        }
    }
}

#[no_mangle]
pub fn test_call_who_am_i_on_external_contract() -> i32 {
    let addr = hmc::hex_to_bytes(hmc::get_arg_str(0).unwrap().as_ref());
    match hmc::call_contract(&addr, "who_am_i".as_bytes(), vec![]) {
        Ok(v) => {
            hmc::return_value(&v)
        },
        Err(m) => {
            hmc::revert(format!("{}", m));
            -1
        }
    }
}

#[no_mangle]
pub fn init() -> i32 {
    0
}
