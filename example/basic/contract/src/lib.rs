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

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
