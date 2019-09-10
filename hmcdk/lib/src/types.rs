use crate::error::{from_str, Error};
use std::borrow::Borrow;

pub type Value = Vec<u8>;
pub type R<T> = Result<Option<T>, Error>;
pub type ArgBytes = Vec<u8>;
pub type Address = [u8; 20];

pub trait FromBytes: Sized {
    fn from_bytes<T: Borrow<ArgBytes>>(_: T) -> Result<Self, Error>;
}

macro_rules! num_from_bytes_impl {
    ($(($t:ty,$u:tt))*) => ($(
        impl FromBytes for $t {
            fn from_bytes<T: Borrow<ArgBytes>>(value: T) -> Result<Self, Error> {
                let v = value.borrow();
                if v.len() != $u {
                    Err(from_str(format!(
                        "a length of bytes must be $u, but got {}",
                        v.len()
                    )))
                } else {
                    let mut b: [u8; $u] = Default::default();
                    b.copy_from_slice(v);
                    Ok(<$t>::from_be_bytes(b))
                }
            }
        }
    )*)
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

impl FromBytes for bool {
    fn from_bytes<T: Borrow<ArgBytes>>(value: T) -> Result<Self, Error> {
        let v = value.borrow();
        if v.len() != 1 {
            Err(from_str(format!(
                "a length of bytes must be 1, but got {}",
                v.len()
            )))
        } else {
            match v[0] {
                0 => Ok(false),
                1 => Ok(true),
                n => Err(from_str(format!("expected a boolean value, but got {}", n))),
            }
        }
    }
}

num_from_bytes_impl! {
    (u8,1) (u16,2) (u32,4) (u64,8)
    (i8,1) (i16,2) (i32,4) (i64,8)
}

pub trait ToBytes: Sized {
    fn to_bytes(&self) -> Vec<u8>;
}

macro_rules! num_to_bytes_impl {
    ($($t:ty)*) => ($(
        impl ToBytes for $t {
            fn to_bytes(&self) -> Vec<u8> {
                self.to_be_bytes().to_vec()
            }
        }
    )*)
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

impl ToBytes for bool {
    fn to_bytes(&self) -> Vec<u8> {
        match self {
            true => vec![1u8],
            false => vec![0u8],
        }
    }
}

num_to_bytes_impl! {
    u8 u16 u32 u64
    i8 i16 i32 i64
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
        test_type_conversion(1u8);
        test_type_conversion(1u16);
        test_type_conversion(1u32);
        test_type_conversion(1u64);
        test_type_conversion(1i8);
        test_type_conversion(1i16);
        test_type_conversion(1i32);
        test_type_conversion(1i64);
        test_type_conversion("test".to_string());
        test_type_conversion(Vec::<u8>::new());
        test_type_conversion(b"test".to_vec());
        test_type_conversion(true);
        test_type_conversion(false);

        assert!(test_try_conversion::<i32, i32>(-1).is_ok());
        assert!(test_try_conversion::<i32, u32>(-1).is_ok());
        assert!(test_try_conversion::<i32, i64>(1).is_err());
        assert!(test_try_conversion::<i64, i32>(1).is_err());
        assert!(test_try_conversion::<i16, u16>(1).is_ok());
        assert!(test_try_conversion::<i32, u32>(1).is_ok());
        assert!(test_try_conversion::<i64, u64>(1).is_ok());
        assert!(test_try_conversion::<String, i64>("test".to_string()).is_err());
        assert!(test_try_conversion::<String, u64>("test".to_string()).is_err());
        assert!(test_try_conversion::<String, Vec<u8>>("test".to_string()).is_ok());
        assert!(test_try_conversion::<String, String>("test".to_string()).is_ok());
        assert!(test_try_conversion::<Vec<u8>, String>(b"test".to_vec()).is_ok());
        assert!(test_try_conversion::<Vec<u8>, Vec<u8>>(b"test".to_vec()).is_ok());
        assert!(test_try_conversion::<Vec<u8>, i8>(b"1".to_vec()).is_ok());
        assert!(test_try_conversion::<Vec<u8>, u8>(b"1".to_vec()).is_ok());
        assert!(test_try_conversion::<u8, Vec<u8>>('1' as u8).is_ok());
        assert!(test_try_conversion::<i8, Vec<u8>>('1' as i8).is_ok());
        assert!(test_try_conversion::<Vec<u8>, i64>(b"test".to_vec()).is_err());
        assert!(test_try_conversion::<Vec<u8>, u64>(b"test".to_vec()).is_err());
        assert!(test_try_conversion::<i8, bool>(1).is_ok());
        assert!(test_try_conversion::<i8, bool>(0).is_ok());
        assert!(test_try_conversion::<i8, bool>(-1).is_err());
        assert!(test_try_conversion::<u8, bool>(2).is_err());
        assert!(test_try_conversion::<u8, bool>(1).is_ok());
        assert!(test_try_conversion::<u8, bool>(0).is_ok());
        assert!(test_try_conversion::<u8, bool>(2).is_err());
        assert!(test_try_conversion::<bool, u8>(true).is_ok());
        assert!(test_try_conversion::<bool, u8>(false).is_ok());
        assert!(test_try_conversion::<bool, i8>(true).is_ok());
        assert!(test_try_conversion::<bool, i8>(false).is_ok());
    }
}
