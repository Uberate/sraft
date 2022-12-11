use crate::config::ConfigAble;
use crate::storage::StorageEngine;

struct Engine {

}

impl StorageEngine for Engine{
    fn new(&self, config: ConfigAble) -> Box<dyn StorageEngine> {
        Box::new(Engine{})
    }

    fn set(&mut self, key: String, value: String) {
        todo!()
    }

    fn get(&self, path: String) -> Option<&String> {
        todo!()
    }

    fn get_clone(&self, path: String) -> Option<String> {
        todo!()
    }

    fn remove(&mut self, path: String) -> Option<String> {
        todo!()
    }

    fn remove_entry(&mut self, path: String) -> Option<(String, String)> {
        todo!()
    }

    fn count_path(&self) -> usize {
        todo!()
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