use std::str;

extern "C" {
    fn __get_arg(idx: usize, value_buf_ptr: *mut u8, value_buf_len: usize) -> i64;

    fn __set_response(msg: *const u8, len: usize);
    fn __log(msg: *const u8, len: usize);

    fn __read_state(msg: *const u8, len: usize, value_buf_ptr: *mut u8, value_buf_len: usize) -> i64;
    fn __write_state(msg1: *const u8, len1: usize, msg2: *const u8, len2: usize) -> i64;
}

struct API {}

impl API {
    fn get_arg_str(idx: usize) -> Result<String, String> {
        let mut buf = [0u8; 64];
        let size = unsafe {
            __get_arg(idx, buf.as_mut_ptr(), 64)
        };
        if size == -1 {
            return Err(format!("argument {} not found", idx));
        }
        match str::from_utf8(&buf[0 .. size as usize]) {
            Ok(v) => Ok(v.to_string()),
            Err(e) => Err(format!("Invalid UTF-8 sequence: {}", e)),
        }
    }

    fn log(b: &[u8]) {
        unsafe {
            __log(b.as_ptr(), b.len());
        }
    }

    fn read_state(key: &[u8]) -> Result<String, String> {
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

    fn write_state(key: &[u8], value: &[u8]) {
        unsafe {
            let ret = __write_state(key.as_ptr(), key.len(), value.as_ptr(), value.len());
            if ret == -1 {
                return;
            }
        }
    }

    fn revert(msg: String) {
        API::log(msg.as_bytes());
        panic!(msg);
    }
}

#[no_mangle]
pub extern "C" fn app_main() -> i32 {
    let name = API::get_arg_str(0).unwrap();
    let amount = API::get_arg_str(1).unwrap().parse::<i64>().unwrap();
    if amount <= 0 {
        API::revert(format!("must specify posotive value, not {}", amount))
    }

    API::log(format!("will incr {}", amount).as_bytes());
    let v: i64 = match API::read_state(name.as_bytes()) {
        Ok(v) => {
            API::log(format!("read {}", v).as_bytes());
            v.parse::<i64>().unwrap()
        },
        Err(m) => {
            API::log(m.as_bytes());
            0
        },
    } + amount;

    API::log(format!("will write {}", v).as_bytes());
    API::write_state(name.as_bytes(), format!("{}", v).as_bytes());

    return 0;
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
