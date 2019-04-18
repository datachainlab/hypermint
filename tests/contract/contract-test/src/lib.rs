extern crate hmc;

fn call_check_signature() -> Result<i64, String> {
    let sender = hmc::get_sender()?;
    let msgHash = hmc::get_arg(0)?;
    let sig = hmc::get_arg(1)?;
    let addr = hmc::ecrecover_address(&msgHash, &sig[64..65], &sig[0..32], &sig[32..64])?;

    if sender == addr {
        Ok(0)
    } else {
        Err(format!("invalid signer"))
    }
}

#[no_mangle]
pub fn check_signature() -> i64 {
    match call_check_signature() {
        Ok(v) => v,
        Err(e) => {
            hmc::log(e.as_bytes());
            1
        }
    }
}

#[no_mangle]
pub fn init() -> i64 {
    0
}
