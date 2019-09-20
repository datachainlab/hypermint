
pub enum AbiInteger {
    I8, I16, I32, I64,
    U8, U16, U32, U64,
}

pub enum AbiPrimitive {
    Integer(AbiInteger),
    Bytes,
    Address,
    Hash,
}

pub fn get_primitive_name(prim: AbiPrimitive) -> String {
    use AbiPrimitive::*;
    use AbiInteger::*;

    match prim {
        Integer (x) => match x {
            I8 => "i8",
            I16 => "i16",
            I32 => "i32",
            I64 => "i64",
            U8 => "u8",
            U16 => "u16",
            U32 => "u32",
            U64 => "u64",
        },
        Bytes => "bytes",
        Address => "address",
        Hash => "hash",
    }.to_string()
}

pub fn to_abi_primitive(rust_type: &String) -> Option<AbiPrimitive> {
    use AbiPrimitive::*;
    use AbiInteger::*;

    let t = rust_type.as_str();
    match t {
        "i8" => Some(Integer(I8)),
        "i16" => Some(Integer(I16)),
        "i32" => Some(Integer(I32)),
        "i64" => Some(Integer(I64)),
        "u8" => Some(Integer(U8)),
        "u16" => Some(Integer(U16)),
        "u32" => Some(Integer(U32)),
        "u64" => Some(Integer(U64)),
        "Address" => Some(Address),
        "Hash" => Some(Hash),
        "Vec" => Some(Bytes),
        _ => None,
    }
}
