use crate::config::ConfigAble;

pub mod config;
pub mod point;
pub mod plugins;
pub mod server;
pub mod storage;


// == test mod
#[cfg(test)]
mod config_test;


fn main() {
    let config_instance = ConfigAble::new_empty();
    if let Some(s_instance) = storage::new_engine(
        String::from("memory"),
        String::from("v1"),
        &config_instance) {
        println!("Init success")
    }else {
        panic!("Error to init storage instance.")
    }
    println!("Hello, world!");
}
