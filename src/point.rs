use std::time::{SystemTime, UNIX_EPOCH};
use std::result::Result as sResult;
use std::error::Error as sErr;
use serde::{Deserialize, Serialize};
use serde_json::Result;
use crate::config::ConfigAble;

/**
The mod point abstract the RPC(Remote process call). The sraft not care the implements of the RPC.
 */

pub trait PointEngine {
    fn server(id: String, config_able: ConfigAble) -> Box<dyn Server>;
    fn client(id: String, config_able: ConfigAble) -> Box<dyn Client>;
}

/**
SendMessage package the message from client to server.
 */
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

    fn generator_receive_message(&self) -> ResponseMessage {
        ResponseMessage {
            source: self.target.clone(),
            target: self.source.clone(),
            // Get now nano second. From UNIX_EPOCH.
            create_at: SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_nanos(),
            payload: "".to_string(),
            send_message: self.clone(),
        }
    }
}

/**
ResponseMessage package the message from server to client.
 */
#[derive(Serialize, Deserialize)]
pub struct ResponseMessage {
    pub source: String,
    pub target: String,

    pub create_at: u128,
    pub payload: String,

    pub send_message: SendMessage,
    // todo: add the error message
}

// todo its a error
pub struct ErrMessage {

}

impl ResponseMessage {
    pub fn to_any<'a, T>(&'a self) -> Result<T> where T: Deserialize<'a> {
        serde_json::from_str(self.payload.as_str())
    }

    pub fn send_message(&self) -> SendMessage {
        self.send_message.clone()
    }

    pub fn duration(&self) -> u128 {
        self.create_at - self.send_message.create_at
    }
}

/**
Client send message to [Server], and receive the [ResponseMessage] from [Server]. Any server should
bind with a [Server]. But [Server] does not, one [Server] can bind more then one [Client].

The [Client]'s kind and version should equals of [Server], if not, the request may be errors.
 */
pub trait Client {
    fn id(&self) -> &String;
    fn server_point(&self) -> &String;

    // todo: add send_message function
    // todo: add send_any function
}

pub trait Server {}