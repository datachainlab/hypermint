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

    fn check_conversion_hex_and_bytes(s: &str, b: &[u8]) {
        let v = hex_to_bytes(s);
        assert_eq!(v, b.to_vec());
    }

    #[test]
    fn test_hex_conversion() {
        check_conversion_hex_and_bytes("0x0000", &[0, 0]);
        check_conversion_hex_and_bytes("0x0100", &[1, 0]);
        check_conversion_hex_and_bytes("0x0101", &[1, 1]);
        check_conversion_hex_and_bytes("0xFF01", &[255, 1]);
        check_conversion_hex_and_bytes("0xFFFF", &[255, 255]);
        check_conversion_hex_and_bytes("FFFF", &[255, 255]);
        check_conversion_hex_and_bytes("0xffffffff", &[255, 255, 255, 255]);
        check_conversion_hex_and_bytes("0xFffFfFFf", &[255, 255, 255, 255]);
    }
}
