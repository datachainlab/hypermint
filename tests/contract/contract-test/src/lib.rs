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
    let b = [0u8; 255];
    hmc::write_state("key".as_bytes(), &b);
    0
}

#[no_mangle]
pub fn test_read_state() -> i32 {
    let b = [0u8; 255];
    hmc::write_state("key".as_bytes(), &b);
    0
}

#[no_mangle]
pub fn test_keccak256() -> i32 {
    let msg = hmc::get_arg(0).unwrap();
    let h = hmc::keccak256(&msg).unwrap();
    hmc::return_value(&h)
}

#[no_mangle]
pub fn init() -> i32 {
    0
}
