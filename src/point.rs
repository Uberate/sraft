use std::fmt::{Debug, Error};
use std::time::{SystemTime, UNIX_EPOCH};
use serde::{Deserialize, Serialize};
use serde_json::Result;
use crate::config::ConfigAble;

/**
The mod point abstract the RPC(Remote process call). The sraft not care the implements of the RPC.
 */

pub trait PointEngine {
    fn server(name: String, config_able: ConfigAble) -> Box<dyn Server>;
    fn client(name: String, config_able: ConfigAble) -> Box<dyn Client>;
}

#[derive(Serialize, Deserialize, Clone)]
pub struct SendMessage {
    pub source: String,
    pub target: String,

    pub create_at: u128,
    pub payload: String,
}

impl SendMessage {
    fn from_string(source: String, target: String, payload: String) -> SendMessage {
        SendMessage {
            source,
            target,
            payload,
            // Get now nano second. From UNIX_EPOCH.
            create_at: SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_nanos(),
        }
    }

    fn new_empty_payload(source: String, target: String) -> SendMessage {
        SendMessage::from_string(source, target, "".to_string())
    }

    fn from_any<T>(source: String, target: String, value: &T) -> Result<SendMessage>
        where T: Sized + Serialize {
        match serde_json::to_string(value) {
            Ok(v) => {
                Ok(SendMessage::from_string(source, target, v))
            }
            Err(err) => {
                Err(err)
            }
        }
    }

    fn generator_receive_message(&self) -> ReceiveMessage {
        ReceiveMessage {
            source: self.target.clone(),
            target: self.source.clone(),
            // Get now nano second. From UNIX_EPOCH.
            create_at: SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_nanos(),
            payload: "".to_string(),
            send_message: self.clone(),
        }
    }
}

#[derive(Serialize, Deserialize)]
pub struct ReceiveMessage {
    pub source: String,
    pub target: String,

    pub create_at: u128,
    pub payload: String,

    pub send_message: SendMessage,
}

impl ReceiveMessage {
    pub fn to_any<'a, T>(&'a self) -> Result<T> where T: Deserialize<'a> {
        serde_json::from_str(self.payload.as_str())
    }
}

pub trait Client {}

pub trait Server {}