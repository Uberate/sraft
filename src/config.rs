use std::collections::HashMap;
use serde::{Serialize, Deserialize};
use serde_json::map::{Map};
use serde_json::Value;

///
///ConfigAble struct represent the json object from serde_json. The user can operate the inner value
///directly. The ConfigAble is impl the [Serialize] and [Deserialize].
///
///It abstract the method for different client.
#[derive(Serialize, Deserialize, Debug)]
pub struct ConfigAble {
    pub value: Map<String, Value>,
}

impl ConfigAble {
    /**
    new function will create a ConfigAble by specify value.

    */
    pub fn new(value: Map<String, Value>) -> Self {
        Self { value }
    }

    /// new_empty return an empty ConfigAble.
    pub fn new_empty() -> Self {
        Self { value: Map::new() }
    }

    /// to_json return the json str for use. In the client and implements, use it to deserialize
    /// data quickly.
    pub fn to_json(&self) -> String {
        serde_json::to_string(&self.value).unwrap()
    }
}
