use serde::{Deserialize, Serialize};
use crate::config::ConfigAble;
use crate::error::ErrorMessage;
use crate::point::{Client, PointEngine, ResponseMessage, SendMessage, Server};

pub mod http_v1;

pub struct Engine {}

impl PointEngine for Engine {
    fn server(config_able: ConfigAble) -> Box<dyn Server> {
        todo!()
    }

    fn client(config_able: ConfigAble) -> Box<dyn Client> {
        todo!()
    }
}

#[derive(Serialize, Deserialize)]
pub struct HttpV1Server {
    id: String,
    listen_point: String,
}

impl HttpV1Server {
    fn new(config_able: ConfigAble) -> Box<dyn Server> {
        let config_json = config_able.to_json();
        let http_v1_server = serde_json::from_str(&config_json).unwrap();
        http_v1_server
    }
}

impl Server for HttpV1Server {
    fn id(&self) -> &String {
        todo!()
    }

    fn point(&self) -> &String {
        todo!()
    }

    fn handler(&self, path: String, Fun: fn(SendMessage) -> Result<ResponseMessage, ErrorMessage>) {
        todo!()
    }

    fn run(&self) -> Result<(), ErrorMessage> {
        todo!()
    }
}

pub struct HttpV1Client {}

impl Client for HttpV1Client {
    fn id(&self) -> &String {
        todo!()
    }

    fn server_point(&self) -> &String {
        todo!()
    }

    fn send_message(&self, send_message: SendMessage) -> Result<ResponseMessage, ErrorMessage> {
        todo!()
    }

    fn ping(&self) {
        todo!()
    }
}