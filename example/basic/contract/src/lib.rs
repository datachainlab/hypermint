extern "C" {
    fn __read_str(msg: *const u8, len: usize, value_buf_ptr: *mut u8, value_buf_len: usize) -> i64;
    fn __write_str(msg1: *const u8, len1: usize, msg2: *const u8, len2: usize) -> i64;
}

#[no_mangle]
pub extern "C" fn app_main() -> i32 {
    let key = "key".as_bytes();
    let value = "value".as_bytes();

    let mut val_buf = [0u8; 64];

    unsafe {
        let ret = __read_str(key.as_ptr(), key.len(), val_buf.as_mut_ptr(), 64);
        if (ret == -1) {
            __write_str(key.as_ptr(), key.len(), value.as_ptr(), value.len());
        }
    }

    return 0;
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
