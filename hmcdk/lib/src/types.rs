use crate::error::{from_str, Error};
use std::borrow::Borrow;

pub type Value = Vec<u8>;
pub type R<T> = Result<Option<T>, Error>;
pub type ArgBytes = Vec<u8>;
pub type Address = [u8; 20];

pub trait FromBytes: Sized {
    fn from_bytes<T: Borrow<ArgBytes>>(_: T) -> Result<Self, Error>;
}

impl FromBytes for u32 {
    fn from_bytes<T: Borrow<ArgBytes>>(value: T) -> Result<Self, Error> {
        let v = value.borrow();
        if v.len() != 4 {
            Err(from_str(format!(
                "a length of bytes must be 4, but got {}",
                v.len()
            )))
        } else {
            let mut b: [u8; 4] = Default::default();
            b.copy_from_slice(v);
            Ok(u32::from_be_bytes(b))
        }
    }
}

impl FromBytes for u64 {
    fn from_bytes<T: Borrow<ArgBytes>>(value: T) -> Result<Self, Error> {
        let v = value.borrow();
        if v.len() != 8 {
            Err(from_str(format!(
                "a length of bytes must be 8, but got {}",
                v.len()
            )))
        } else {
            let mut b: [u8; 8] = Default::default();
            b.copy_from_slice(value.borrow());
            Ok(u64::from_be_bytes(b))
        }
    }
}

impl FromBytes for i32 {
    fn from_bytes<T: Borrow<ArgBytes>>(value: T) -> Result<Self, Error> {
        let v = value.borrow();
        if v.len() != 4 {
            Err(from_str(format!(
                "a length of bytes must be 4, but got {}",
                v.len()
            )))
        } else {
            let mut b: [u8; 4] = Default::default();
            b.copy_from_slice(v);
            Ok(i32::from_be_bytes(b))
        }
    }
}

impl FromBytes for i64 {
    fn from_bytes<T: Borrow<ArgBytes>>(value: T) -> Result<Self, Error> {
        let v = value.borrow();
        if v.len() != 8 {
            Err(from_str(format!(
                "a length of bytes must be 8, but got {}",
                v.len()
            )))
        } else {
            let mut b: [u8; 8] = Default::default();
            b.copy_from_slice(value.borrow());
            Ok(i64::from_be_bytes(b))
        }
    }
}

impl FromBytes for String {
    fn from_bytes<T: Borrow<ArgBytes>>(value: T) -> Result<Self, Error> {
        Ok(String::from_utf8(value.borrow().clone())?)
    }
}

impl FromBytes for Vec<u8> {
    fn from_bytes<T: Borrow<ArgBytes>>(value: T) -> Result<Self, Error> {
        Ok(value.borrow().clone())
    }
}

impl FromBytes for Address {
    fn from_bytes<T: Borrow<ArgBytes>>(value: T) -> Result<Self, Error> {
        let v = value.borrow();
        if v.len() != 20 {
            Err(from_str(format!(
                "a length of bytes must be 20, but got {}",
                v.len()
            )))
        } else {
            let mut addr: Address = Default::default();
            addr.copy_from_slice(v);
            Ok(addr)
        }
    }
}

pub trait ToBytes: Sized {
    fn to_bytes(&self) -> Vec<u8>;
}

impl ToBytes for i32 {
    fn to_bytes(&self) -> Vec<u8> {
        self.to_be_bytes().to_vec()
    }
}

impl ToBytes for i64 {
    fn to_bytes(&self) -> Vec<u8> {
        self.to_be_bytes().to_vec()
    }
}

impl ToBytes for u32 {
    fn to_bytes(&self) -> Vec<u8> {
        self.to_be_bytes().to_vec()
    }
}

impl ToBytes for u64 {
    fn to_bytes(&self) -> Vec<u8> {
        self.to_be_bytes().to_vec()
    }
}

impl ToBytes for Vec<u8> {
    fn to_bytes(&self) -> Vec<u8> {
        self.clone()
    }
}

impl ToBytes for String {
    fn to_bytes(&self) -> Vec<u8> {
        self.as_bytes().to_vec()
    }
}

impl ToBytes for Address {
    fn to_bytes(&self) -> Vec<u8> {
        (&self).to_vec()
    }
}

// TODO
impl ToBytes for &Address {
    fn to_bytes(&self) -> Vec<u8> {
        (&self).to_vec()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    struct Args {
        args: Vec<ArgBytes>,
    }

    impl Args {
        fn new() -> Self {
            Args { args: vec![] }
        }

        fn push(&mut self, v: ArgBytes) {
            self.args.push(v)
        }

        pub fn get_arg<T: FromBytes>(&self, idx: usize) -> Result<T, Error> {
            let b = &self.args[idx];
            Ok(T::from_bytes(b)?)
        }
    }

    fn test_type_conversion<T: ToBytes + FromBytes + std::fmt::Debug + Eq>(v: T) {
        let mut a = Args::new();
        let b1 = v.to_bytes();
        a.push(b1.clone());

        let vv: Vec<u8> = a.get_arg(0).unwrap();
        assert_eq!(b1, vv);
        let vi: T = a.get_arg(0).unwrap();
        assert_eq!(v, vi);
    }

    fn test_try_conversion<
        T: ToBytes + FromBytes + std::fmt::Debug + Eq,
        U: ToBytes + FromBytes + std::fmt::Debug + Eq,
    >(
        v: T,
    ) -> Result<(), Error> {
        let mut a = Args::new();
        let b1 = v.to_bytes();
        a.push(b1.clone());

        let _: U = a.get_arg(0)?;
        Ok(())
    }

    #[test]
    fn test_from_bytes() {
        test_type_conversion(1u32);
        test_type_conversion(1u64);
        test_type_conversion(1i32);
        test_type_conversion(1i64);
        test_type_conversion("test".to_string());
        test_type_conversion(Vec::<u8>::new());
        test_type_conversion(b"test".to_vec());

        assert!(test_try_conversion::<i32, i64>(1).is_err());
        assert!(test_try_conversion::<i64, i32>(1).is_err());
        assert!(test_try_conversion::<i32, u32>(1).is_ok());
        assert!(test_try_conversion::<i64, u64>(1).is_ok());
        assert!(test_try_conversion::<String, i64>("test".to_string()).is_err());
        assert!(test_try_conversion::<String, u64>("test".to_string()).is_err());
        assert!(test_try_conversion::<String, Vec<u8>>("test".to_string()).is_ok());
        assert!(test_try_conversion::<String, String>("test".to_string()).is_ok());
        assert!(test_try_conversion::<Vec<u8>, String>(b"test".to_vec()).is_ok());
        assert!(test_try_conversion::<Vec<u8>, Vec<u8>>(b"test".to_vec()).is_ok());
        assert!(test_try_conversion::<Vec<u8>, i64>(b"test".to_vec()).is_err());
        assert!(test_try_conversion::<Vec<u8>, u64>(b"test".to_vec()).is_err());
    }
}
