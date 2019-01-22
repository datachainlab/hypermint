extern crate hmc;

static TOTAL: i64 = 10000;

#[no_mangle]
pub fn get_balance() -> i64 {
    let sender = hmc::get_sender().unwrap();
    match hmc::read_state(&sender) {
        Ok(v) => {
            hmc::log(format!("read {:?}", v).as_bytes());
            hmc::return_value(&v)
        },
        Err(m) => {
            hmc::log(m.as_bytes());
            -1
        },
    }
}

fn get_balance_from_addr(addr: &[u8]) -> i64 {
    match hmc::read_state_str(addr) {
        Ok(v) => {
            v.parse::<i64>().unwrap()
        },
        Err(m) => {
            hmc::log(m.as_bytes());
            0
        },
    }
}

#[no_mangle]
pub fn transfer() -> i64 {
    let to = hmc::hex_to_bytes(hmc::get_arg_str(0).unwrap().as_ref());
    let amount = hmc::get_arg_str(1).unwrap().parse::<i64>().unwrap();
    let sender = hmc::get_sender().unwrap();

    let from_balance = get_balance_from_addr(&sender);
    if from_balance <= amount {
        hmc::log(format!("error: {} <= {}", from_balance, amount).as_bytes());
        return -1;
    }
    let to_balance = get_balance_from_addr(&to);

    hmc::write_state(&sender, format!("{}", from_balance-amount).as_bytes());
    let to_amount = format!("{}", to_balance+amount);
    hmc::write_state(&to, to_amount.as_bytes());

    return hmc::return_value(to_amount.as_bytes());
}

#[no_mangle]
pub fn init() -> i64 {
    // TODO add initializer once
    let sender = hmc::get_sender().unwrap();
    hmc::write_state(&sender, format!("{}", TOTAL).as_bytes());
    0
}
