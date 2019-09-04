use failure::{Error as FError, Fail};

pub type Error = FError;

#[derive(Debug, Fail)]
#[fail(display = "failed to call contract: {:?}", msg)]
pub struct ErrorWithMsg {
    msg: String,
}

pub fn from_str<T: Into<String>>(msg: T) -> Error {
    Error::from(ErrorWithMsg { msg: msg.into() })
}
