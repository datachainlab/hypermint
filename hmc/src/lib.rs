use std::str;

const BUF_SIZE: usize = 128;

extern "C" {
    fn __get_arg(idx: usize, offset: usize, value_buf_ptr: *mut u8, value_buf_len: usize) -> i32;
    fn __get_sender(value_buf_ptr: *mut u8, value_buf_len: usize) -> i32;
    fn __get_contract_address(value_buf_ptr: *mut u8, value_buf_len: usize) -> i32;
    fn __read(id: usize, offset: usize, value_buf_ptr: *mut u8, value_buf_len: usize) -> i32;
    fn __call_contract(
        addr: *const u8,
        addr_size: usize,
        entry: *const u8,
        entry_size: usize,
        args: *const u8,
        args_size: usize,
    ) -> i32;

    fn __set_response(msg: *const u8, len: usize) -> i32;
    fn __log(msg: *const u8, len: usize) -> i32;

    fn __read_state(
        key_ptr: *const u8,
        key_len: usize,
        offset: usize,
        value_buf_ptr: *mut u8,
        value_buf_len: usize,
    ) -> i32;
    fn __write_state(key: *const u8, key_len: usize, value: *const u8, value_len: usize) -> i32;

    fn __keccak256(
        msg: *const u8,
        msg_len: usize,
        value_buf_ptr: *mut u8,
        value_buf_len: usize,
    ) -> i32;

    fn __sha256(
        msg: *const u8,
        msg_len: usize,
        value_buf_ptr: *mut u8,
        value_buf_len: usize,
    ) -> i32;

    fn __ecrecover(
        h: *const u8,
        h_len: usize,
        v: *const u8,
        v_len: usize,
        r: *const u8,
        r_len: usize,
        s: *const u8,
        s_len: usize,
        ret: *mut u8,
        ret_len: usize,
    ) -> i32;
    fn __ecrecover_address(
        h: *const u8,
        h_len: usize,
        v: *const u8,
        v_len: usize,
        r: *const u8,
        r_len: usize,
        s: *const u8,
        s_len: usize,
        ret: *mut u8,
        ret_len: usize,
    ) -> i32;
    fn __emit_event(
        ev: *const u8,
        ev_len: usize,
        data: *const u8,
        data_len: usize,
    ) -> i32;
}

pub fn keccak256(msg: &[u8]) -> Result<[u8; 32], String> {
    let mut buf = [0u8; 32];
    match unsafe { __keccak256(msg.as_ptr(), msg.len(), buf.as_mut_ptr(), buf.len()) } {
        -1 => Err(format!("fail to call keccak256")),
        _ => Ok(buf),
    }
}

pub fn sha256(msg: &[u8]) -> Result<[u8; 32], String> {
    let mut buf = [0u8; 32];
    match unsafe { __sha256(msg.as_ptr(), msg.len(), buf.as_mut_ptr(), buf.len()) } {
        -1 => Err(format!("fail to call sha256")),
        _ => Ok(buf),
    }
}

pub fn emit_event(name: &str, value: &[u8]) -> Result<(), String> {
    match unsafe {
        __emit_event(name.as_ptr(), name.len(), value.as_ptr(), value.len())
    } {
        -1 => Err(format!("fail to emit event")),
        _ => Ok(())
    }
}

pub fn ecrecover(h: &[u8], v: &[u8], r: &[u8], s: &[u8]) -> Result<[u8; 65], String> {
    if h.len() != 32 {
        return Err(format!("length of h should be 32, got {}", h.len()));
    } else if v.len() != 1 {
        return Err(format!("length of v should be 1, got {}", v.len()));
    } else if r.len() != 32 {
        return Err(format!("length of r should be 32, got {}", r.len()));
    } else if s.len() != 32 {
        return Err(format!("length of s should be 32, got {}", s.len()));
    }
    let mut buf = [0u8; 65];
    match unsafe {
        __ecrecover(
            h.as_ptr(),
            h.len(),
            v.as_ptr(),
            v.len(),
            r.as_ptr(),
            r.len(),
            s.as_ptr(),
            s.len(),
            buf.as_mut_ptr(),
            buf.len(),
        )
    } {
        -1 => Err(format!("fail to ecrecover")),
        _ => Ok(buf),
    }
}

pub fn ecrecover_address(h: &[u8], v: &[u8], r: &[u8], s: &[u8]) -> Result<[u8; 20], String> {
    if h.len() != 32 {
        return Err(format!("length of h should be 32, got {}", h.len()));
    } else if v.len() != 1 {
        return Err(format!("length of v should be 1, got {}", v.len()));
    } else if r.len() != 32 {
        return Err(format!("length of r should be 32, got {}", r.len()));
    } else if s.len() != 32 {
        return Err(format!("length of s should be 32, got {}", s.len()));
    }
    let mut buf = [0u8; 20];
    match unsafe {
        __ecrecover_address(
            h.as_ptr(),
            h.len(),
            v.as_ptr(),
            v.len(),
            r.as_ptr(),
            r.len(),
            s.as_ptr(),
            s.len(),
            buf.as_mut_ptr(),
            buf.len(),
        )
    } {
        -1 => Err(format!("fail to ecrecover")),
        _ => Ok(buf),
    }
}

pub fn get_arg(idx: usize) -> Result<Vec<u8>, String> {
    let mut buf = [0u8; BUF_SIZE];
    let mut offset = 0;
    let mut val: Vec<u8> = Vec::new();
    loop {
        match unsafe { __get_arg(idx, offset, buf.as_mut_ptr(), buf.len()) } {
            -1 => return Err("read_state: key not found".to_string()),
            0 => break,
            n => {
                val.extend_from_slice(&buf[0..n as usize]);
                if n < BUF_SIZE as i32 {
                    break;
                }
                offset += n as usize;
            }
        }
    }
    Ok(val)
}

pub fn get_arg_str(idx: usize) -> Result<String, String> {
    let v = get_arg(idx)?;
    match str::from_utf8(&v) {
        Ok(v) => Ok(v.to_string()),
        Err(e) => Err(format!("Invalid UTF-8 sequence: {}", e)),
    }
}

pub fn get_sender() -> Result<[u8; 20], String> {
    let mut buf = [0u8; 20];
    match unsafe { __get_sender(buf.as_mut_ptr(), 20) } {
        -1 => Err("sender not found".to_string()),
        _ => Ok(buf),
    }
}

pub fn get_sender_str() -> Result<String, String> {
    let sender = get_sender()?;
    Ok(format!("{:X?}", sender))
}

pub fn get_contract_address() -> Result<[u8; 20], String> {
    let mut buf = [0u8; 20];
    match unsafe { __get_contract_address(buf.as_mut_ptr(), 20) } {
        -1 => Err("contract address not found".to_string()),
        _ => Ok(buf),
    }
}

pub fn call_contract(addr: &[u8], entry: &[u8], args: Vec<&[u8]>) -> Result<Vec<u8>, String> {
    let a = serialize_args(&args);
    let id = match unsafe {
        __call_contract(
            addr.as_ptr(),
            addr.len(),
            entry.as_ptr(),
            entry.len(),
            a.as_ptr(),
            a.len(),
        )
    } {
        -1 => return Err("failed to call contract".to_string()),
        id => id as usize,
    };

    let mut buf = [0u8; BUF_SIZE];
    let mut offset = 0;
    let mut val: Vec<u8> = Vec::new();

    loop {
        match unsafe { __read(id, offset, buf.as_mut_ptr(), buf.len()) } {
            -1 => return Err("read_state: key not found".to_string()),
            0 => break,
            n => {
                val.extend_from_slice(&buf[0..n as usize]);
                if n < BUF_SIZE as i32 {
                    break;
                }
                offset += n as usize;
            }
        }
    }
    Ok(val)
}

// format: <elem_num: 4byte>|<elem1_size: 4byte>|<elem1_data>|<elem2_size: 4byte>|<elem2_data>|...
fn serialize_args(args: &Vec<&[u8]>) -> Vec<u8> {
    let mut bs: Vec<u8> = vec![];
    bs.extend_from_slice(&(args.len() as u32).to_be_bytes());
    for arg in args {
        bs.extend_from_slice(&(arg.len() as u32).to_be_bytes());
        bs.extend_from_slice(arg);
    }
    bs
}

pub fn log(b: &[u8]) -> i32 {
    unsafe { __log(b.as_ptr(), b.len()) }
}

pub fn read_state(key: &[u8]) -> Result<Vec<u8>, String> {
    let mut val_buf = [0u8; BUF_SIZE];
    let mut offset = 0;
    let mut val: Vec<u8> = Vec::new();
    loop {
        match unsafe {
            __read_state(
                key.as_ptr(),
                key.len(),
                offset,
                val_buf.as_mut_ptr(),
                val_buf.len(),
            )
        } {
            -1 => return Err("read_state: key not found".to_string()),
            0 => break,
            n => {
                val.extend_from_slice(&val_buf[0..n as usize]);
                if n < BUF_SIZE as i32 {
                    break;
                }
                offset += n as usize;
            }
        }
    }
    Ok(val)
}

pub fn read_state_str(key: &[u8]) -> Result<String, String> {
    let v = read_state(key)?;
    match str::from_utf8(&v) {
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

pub fn return_value(v: &[u8]) -> i32 {
    unsafe { __set_response(v.as_ptr(), v.len()) }
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
    let mut hex_bytes = bs
        .iter()
        .filter_map(|b| match b {
            b'0'...b'9' => Some(b - b'0'),
            b'a'...b'f' => Some(b - b'a' + 10),
            b'A'...b'F' => Some(b - b'A' + 10),
            _ => None,
        })
        .fuse();

    let mut bytes = Vec::new();
    while let (Some(h), Some(l)) = (hex_bytes.next(), hex_bytes.next()) {
        bytes.push(h << 4 | l)
    }
    bytes
}

#[cfg(test)]
mod tests {
    use super::*;

    fn deserialize_args(bs: &Vec<u8>) -> Result<Vec<Vec<u8>>, String> {
        let mut args: Vec<Vec<u8>> = vec![];
        let mut num_bs = [0u8; 4];
        num_bs.copy_from_slice(&bs[0..4]);
        let num = u32::from_be_bytes(num_bs);
        let mut offset: usize = 4;
        for _ in 0..num {
            let mut b = [0u8; 4];
            b.copy_from_slice(&bs[offset..offset + 4]);
            let size = u32::from_be_bytes(b) as usize;
            let mut arg: Vec<u8> = vec![];
            arg.extend_from_slice(&bs[offset + 4..offset + 4 + size]);
            offset += 4 + size;
            args.push(arg);
        }
        Ok(args)
    }

    fn vec_str_to_vec_u8(vs: Vec<&str>) -> Vec<&[u8]> {
        let mut ret: Vec<&[u8]> = vec![];
        for v in vs.iter() {
            ret.push(v.as_bytes());
        }
        ret
    }

    #[test]
    fn encoding_test() {
        let raw = vec!["first", "second", "third"];
        let args = vec_str_to_vec_u8(raw);
        let s = serialize_args(&args);
        let a = deserialize_args(&s).unwrap();

        for (n, m) in args.iter().zip(a.iter()) {
            assert_eq!(m, n);
        }
    }
}
