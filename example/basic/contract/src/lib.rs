use std::str;

extern "C" {
    fn __get_arg_str(idx: usize, value_buf_ptr: *mut u8, value_buf_len: usize) -> i64;
    fn __set_response(msg: *const u8, len: usize);
    fn __log(msg: *const u8, len: usize);

    fn __read_state(msg: *const u8, len: usize, value_buf_ptr: *mut u8, value_buf_len: usize) -> i64;
    fn __write_state(msg1: *const u8, len1: usize, msg2: *const u8, len2: usize) -> i64;
}

struct API {}

impl API {
    fn get_arg_str(idx: usize) -> String {
        let mut buf = [0u8; 64];
        let size = unsafe {
            __get_arg_str(idx, buf.as_mut_ptr(), 64)
        };
        if size == -1 {
            panic!("argument not found")
        }
        let s = match str::from_utf8(&buf[0 .. size as usize]) {
            Ok(v) => v,
            Err(e) => panic!("Invalid UTF-8 sequence: {}", e),
        };
        return s.to_string();
    }

    fn log(b: &[u8]) {
        unsafe {
            __log(b.as_ptr(), b.len());
        }
    }

    fn read_state(key: &[u8]) -> String {
        let mut val_buf = [0u8; 64];
        let size = unsafe {
            __read_state(key.as_ptr(), key.len(), val_buf.as_mut_ptr(), val_buf.len())
        };
        if size == -1 {
            panic!("panic")
        }
        let s = match str::from_utf8(&val_buf[0 .. size as usize]) {
            Ok(v) => v,
            Err(e) => panic!("Invalid UTF-8 sequence: {}", e),
        };
        return s.to_string();
    }

    fn write_state(key: &[u8], value: &[u8]) {
        unsafe {
            let ret = __write_state(key.as_ptr(), key.len(), value.as_ptr(), value.len());
            if ret == -1 {
                return;
            }
        }
    }
}

#[no_mangle]
pub extern "C" fn app_main() -> i32 {
    let key = API::get_arg_str(0);
    API::log(key.as_bytes());
    let value = API::get_arg_str(1);
    API::log(value.as_bytes());

    API::write_state(key.as_bytes(), value.as_bytes());

    return 0;
}

#[no_mangle]
pub extern "C" fn app_read() -> i32 {
    let key = API::get_arg_str(0);
    API::log(key.as_bytes());

    let value = API::read_state(key.as_bytes());
    API::log(value.as_bytes());

    return 0;
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
