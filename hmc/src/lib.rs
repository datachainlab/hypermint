use std::str;

extern "C" {
    fn __get_arg(idx: usize, value_buf_ptr: *mut u8, value_buf_len: usize) -> i64;
    fn __get_sender(value_buf_ptr: *mut u8, value_buf_len: usize) -> i64;
    fn __call_contract(addr: *const u8, addr_size: usize, entry: *const u8, entry_size: usize, value_buf_ptr: *mut u8, value_buf_len: usize) -> i64;

    fn __set_response(msg: *const u8, len: usize) -> i64;
    fn __log(msg: *const u8, len: usize);

    fn __read_state(msg: *const u8, len: usize, value_buf_ptr: *mut u8, value_buf_len: usize) -> i64;
    fn __write_state(msg1: *const u8, len1: usize, msg2: *const u8, len2: usize) -> i64;
}

pub fn get_arg_str(idx: usize) -> Result<String, String> {
    let mut buf = [0u8; 64];
    match unsafe {
        __get_arg(idx, buf.as_mut_ptr(), 64)
    } {
        -1 => Err(format!("argument {} not found", idx)),
        size => match str::from_utf8(&buf[0 .. size as usize]) {
            Ok(v) => Ok(v.to_string()),
            Err(e) => Err(format!("Invalid UTF-8 sequence: {}", e)),
        }
    }
}

pub fn get_sender() -> Result<[u8; 20], String> {
    let mut buf = [0u8; 20];
    match unsafe {
        __get_sender(buf.as_mut_ptr(), 20)
    } {
        -1 => Err("sender not found".to_string()),
        _ => Ok(buf)
    }
}

pub fn get_sender_str() -> Result<String, String> {
    let sender = try!(get_sender());
    Ok(format!("{:X?}", sender))
}

pub fn call_contract(addr: &[u8], entry: &[u8]) -> Result<Vec<u8>, String> {
    let mut val_buf = [0u8; 64];
    match unsafe {
        __call_contract(addr.as_ptr(), addr.len(), entry.as_ptr(), entry.len(), val_buf.as_mut_ptr(), val_buf.len())
    } {
        -1 => Err("failed to call contract".to_string()),
        size => Ok((&val_buf[0 .. size as usize]).to_vec())
    }
}

pub fn log(b: &[u8]) {
    unsafe {
        __log(b.as_ptr(), b.len());
    }
}

pub fn read_state(key: &[u8]) -> Result<Vec<u8>, String> {
    let mut val_buf = [0u8; 64];
    match unsafe {
        __read_state(key.as_ptr(), key.len(), val_buf.as_mut_ptr(), val_buf.len())
    } {
        -1 => Err("key not found".to_string()),
        size => Ok((&val_buf[0 .. size as usize]).to_vec())
    }
}

pub fn read_state_str(key: &[u8]) -> Result<String, String> {
    let mut val_buf = [0u8; 64];
    match unsafe {
        __read_state(key.as_ptr(), key.len(), val_buf.as_mut_ptr(), val_buf.len())
    } {
        -1 => Err("key not found".to_string()),
        size => match str::from_utf8(&val_buf[0 .. size as usize]) {
            Ok(v) => Ok(v.to_string()),
            Err(e) => Err(format!("Invalid UTF-8 sequence: {}", e)),
        }
    }
}

pub fn write_state(key: &[u8], value: &[u8]) {
    unsafe {
        let ret = __write_state(key.as_ptr(), key.len(), value.as_ptr(), value.len());
        if ret == -1 {
            return;
        }
    }
}

pub fn return_value(v: &[u8]) -> i64 {
    unsafe {
        __set_response(v.as_ptr(), v.len())
    }
}

pub fn revert(msg: String) {
    log(msg.as_bytes());
    panic!(msg);
}

pub fn hex_to_bytes(hex_asm: &str) -> Vec<u8> {
    let bs = if hex_asm.starts_with("0x") {
        &hex_asm[2..].as_bytes()
    } else {
        hex_asm.as_bytes()
    };
    let mut hex_bytes = bs.iter().filter_map(|b| {
        match b {
            b'0'...b'9' => Some(b - b'0'),
            b'a'...b'f' => Some(b - b'a' + 10),
            b'A'...b'F' => Some(b - b'A' + 10),
            _ => None,
        }
    }).fuse();

    let mut bytes = Vec::new();
    while let (Some(h), Some(l)) = (hex_bytes.next(), hex_bytes.next()) {
        bytes.push(h << 4 | l)
    }
    bytes
}