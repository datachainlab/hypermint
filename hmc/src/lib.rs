use std::str;

extern "C" {
    fn __get_arg(idx: usize, value_buf_ptr: *mut u8, value_buf_len: usize) -> i64;
    fn __get_sender(value_buf_ptr: *mut u8, value_buf_len: usize) -> i64;

    fn __set_response(msg: *const u8, len: usize);
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

pub fn log(b: &[u8]) {
    unsafe {
        __log(b.as_ptr(), b.len());
    }
}

pub fn read_state(key: &[u8]) -> Result<String, String> {
    let mut val_buf = [0u8; 64];
    let size = unsafe {
        __read_state(key.as_ptr(), key.len(), val_buf.as_mut_ptr(), val_buf.len())
    };
    if size == -1 {
        return Err("key not found".to_string())
    }
    match str::from_utf8(&val_buf[0 .. size as usize]) {
        Ok(v) => Ok(v.to_string()),
        Err(e) => Err(format!("Invalid UTF-8 sequence: {}", e)),
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

pub fn revert(msg: String) {
    log(msg.as_bytes());
    panic!(msg);
}