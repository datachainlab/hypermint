extern crate hmc;

#[no_mangle]
pub fn app_main() -> i32 {
    hmc::log(format!("sender is {:?}", hmc::get_sender().unwrap()).as_bytes());

    let sender = hmc::get_sender().unwrap();
    let amount = hmc::get_arg_str(0).unwrap().parse::<i64>().unwrap();
    if amount <= 0 {
        hmc::revert(format!("must specify posotive value, not {}", amount))
    }

    hmc::log(format!("will incr {}", amount).as_bytes());
    let v: i64 = match hmc::read_state_str(&sender) {
        Ok(v) => {
            hmc::log(format!("read {}", v).as_bytes());
            v.parse::<i64>().unwrap()
        },
        Err(m) => {
            hmc::log(m.as_bytes());
            0
        },
    } + amount;

    hmc::log(format!("will write {}", v).as_bytes());
    hmc::write_state(&sender, v.to_string().as_bytes());

    return 0;
}

#[no_mangle]
pub fn call_contract() -> i64 {
    let token_addr = hmc::hex_to_bytes(hmc::get_arg_str(0).unwrap().as_ref());
    match hmc::call_contract(&token_addr, "get_balance".as_bytes(), vec![]) {
        Ok(v) => {
            hmc::return_value(&v)
        },
        Err(m) => {
            hmc::revert(format!("{}", m));
            -1
        }
    }
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
