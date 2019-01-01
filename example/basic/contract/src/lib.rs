extern "C" {
    fn __get_arg_str(idx: usize, value_buf_ptr: *mut u8, value_buf_len: usize) -> i64;

    fn __read_state_str(msg: *const u8, len: usize, value_buf_ptr: *mut u8, value_buf_len: usize) -> i64;
    fn __write_state_str(msg1: *const u8, len1: usize, msg2: *const u8, len2: usize) -> i64;

    fn __log(msg: *const u8, len: usize);
}

#[no_mangle]
pub extern "C" fn app_main() -> i32 {
    let mut key_buf = [0u8; 64];
    let mut val_buf = [0u8; 64];

    unsafe {
        let key_size = __get_arg_str(0, key_buf.as_mut_ptr(), 64);
        if (key_size == -1) {
            return -1;
        }
        let val_size = __get_arg_str(1, val_buf.as_mut_ptr(), 64);
        if (val_size == -1) {
            return -1;
        }

        __log(key_buf.as_ptr(), key_size as usize);
        __log(val_buf.as_ptr(), val_size as usize);

        let ret = __write_state_str(key_buf.as_ptr(), key_size as usize, val_buf.as_ptr(), val_size as usize);
        if (ret == -1) {
            return -1;
        }
    }

    return 0;
}

#[no_mangle]
pub extern "C" fn app_read() -> i32 {
    let mut key_buf = [0u8; 64];
    let mut val_buf = [0u8; 64];

    unsafe {
        let key_size = __get_arg_str(0, key_buf.as_mut_ptr(), 64);
        if (key_size == -1) {
            return -1;
        }

        __log(key_buf.as_ptr(), key_size as usize);
        let size = __read_state_str(key_buf.as_mut_ptr(), key_size as usize, val_buf.as_mut_ptr(), 64);
        if (size == -1) {
            return -1;
        }

        __log(val_buf.as_ptr(), size as usize);
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
