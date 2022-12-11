use std::collections::HashMap;
use crate::config::ConfigAble;
use crate::storage::StorageEngine;

/**
The memory engine is a [StorageEngine] implementation, it storage value in memory. It uses at less
data, and operator quick scenarios.

---
### Metadata
- Version: `v1`
- Type: `memory`

---
### Init param

*NONE*

---
### Usage

#### Get an Engine:
```rust
// Memory does not need any args.
let config_instance = ConfigAble::new_empty();
if let Some(s_instance) = storage::new_engine(
    String::from("memory"),
    String::from("v1"),
    &config_instance) {
    // your logic here
    println!("Init success")
}else {
    panic!("Error to init storage instance.")
}
```
 */
pub struct Engine {
    map: HashMap<String, String>,
}

impl Engine {
    pub fn new(config: &ConfigAble) -> Box<dyn StorageEngine> {
        Box::new(Engine {
            map: HashMap::new()
        })
    }
}

impl StorageEngine for Engine {
    fn set(&mut self, path: String, value: String) {
        self.map.insert(path, value);
    }

    fn get(&self, path: &String) -> Option<&String> {
        self.map.get(path)
    }

    fn get_mut(&mut self, path: &String) -> Option<&mut String> {
        self.map.get_mut(path)
    }

    fn remove(&mut self, path: &String) -> Option<String> {
        self.map.remove(path)
    }

    fn remove_entry(&mut self, path: &String) -> Option<(String, String)> {
        if let Some(value) = self.map.remove(path) {
            Some((path.clone(), value))
        } else {
            None
        }
    }

    fn count_path(&self) -> usize {
        self.map.keys().len()
    }

    fn size(&self) -> usize {
        todo!()
    }

    fn keys(&self, pattern: String) -> Vec<String> {
        todo!()
    }

    fn entries(&self, key_pattern: String) -> Vec<(String, String)> {
        todo!()
    }
}