#[cfg(test)]
pub mod config_test {
    use std::collections::HashMap;
    use serde_json::{Value, json};
    use crate::config::ConfigAble;

    #[test]
    fn test_to_json_string() {

        let mut cf = ConfigAble::new_empty();
        let json_str = cf.to_json();
        assert_eq!(json_str, "{}");

        cf.value.insert(String::from("test"), json!("test"));
        let json_str = cf.to_json();
        assert_eq!(json_str, "{\"test\":\"test\"}")
    }
}