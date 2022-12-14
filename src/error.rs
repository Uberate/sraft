use std::error::Error;
use std::fmt::{Debug, Display, Formatter};

pub struct ErrorMessage {
    err: Box<ErrorImpl>,
}

impl Debug for ErrorMessage {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        todo!()
    }
}

impl Display for ErrorMessage {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        todo!()
    }
}

impl Error for ErrorMessage{

}

pub struct ErrorImpl {

}