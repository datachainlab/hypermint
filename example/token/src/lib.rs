extern crate hmcdk;
use hmcdk::api::{emit_event, get_arg, get_sender, read_state, write_state};
use hmcdk::error;
use hmcdk::prelude::*;

static TOTAL: i64 = 10000;

#[contract]
pub fn get_balance() -> R<i64> {
    Ok(Some(read_state(&get_sender()?)?))
}

fn get_balance_from_addr(addr: &Address) -> Result<i64, Error> {
    read_state::<i64>(addr)
}

#[contract]
pub fn transfer() -> R<i64> {
    let to: Address = get_arg(0)?;
    let amount: i64 = get_arg(1)?;
    let sender = get_sender()?;

    let from_balance = get_balance_from_addr(&sender)?;
    if from_balance < amount {
        return Err(error::from_str(format!(
            "error: {} < {}",
            from_balance, amount
        )));
    }
    let to_balance = get_balance_from_addr(&to).unwrap_or(0);
    write_state(&sender, &(from_balance - amount).to_bytes());
    let to_amount = to_balance + amount;
    write_state(&to, &to_amount.to_bytes());
    emit_event(
        "Transfer",
        format!("from={:X?} to={:X?} amount={}", sender, to, amount).as_bytes(),
    )?;

    Ok(Some(to_amount))
}

#[contract]
pub fn init() -> R<Vec<u8>> {
    let sender = get_sender()?;
    write_state(&sender, &TOTAL.to_bytes());
    Ok(None)
}
