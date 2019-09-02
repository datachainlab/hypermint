extern crate hmc;

#[no_mangle]
pub fn test_plus() -> i32 {
    let x = hmc::get_arg_str(0).unwrap().parse::<i64>().unwrap();
    let y = hmc::get_arg_str(1).unwrap().parse::<i64>().unwrap();
    let msg = format!("{}", x + y);
    hmc::return_value(msg.as_bytes())
}

#[no_mangle]
pub fn who_am_i() -> i32 {
    let sender = hmc::get_sender().unwrap();
    hmc::return_value(&sender)
}

#[no_mangle]
pub fn get_contract_address() -> i32 {
    let address = hmc::get_contract_address().unwrap();
    hmc::return_value(&address)
}

#[no_mangle]
pub fn init() -> i32 {
    0
}
