/**
Mod storage define the interface to deal the data storage problem. In different project, the
storage ability is different. Such like high change with less data requirement, the memory is a
good idea. But the Big-Data, the disk is only way.

Every one who use this repo can custom themselves storage implements. For the repo organized and
layout, I suggest to implements them in `plugins` dir.
 */
pub mod storage {}

use std::result::Iter;
use crate::config::ConfigAble;

/**
StorageEngine is a interface to define the storage ability sraft need. The storage should
provide the Key-Value save. The key was the path.
 */
pub trait StorageEngine {
    fn new(&self, config: ConfigAble) -> Box<dyn StorageEngine>;

    fn set(&mut self, key: String, value: String);

    /**
    Returns a reference to the value corresponding to the key.
     */
    fn get(&self, path: String) -> Option<&String>;

    /**
    Return a mutable reference to the value corresponding to the path.
     */
    fn get_mut(&self, path: String) -> Option<&mut String>;

    /**
    Removes a path from the storage engine, returning the value at the path if the path was
    previously in the storage engine.
     */
    fn remove(&mut self, path: String) -> Option<String>;

    /**
    Removes a path from the map, returning the stored path and value if the path was previously in
    the storage engine.
     */
    fn remove_entry(&mut self, path: String) -> Option<(String, String)>;

    /**
    Returns the number of elements in the storage engine.
     */
    fn count_path(&self) -> usize;

    /**
    Return true when count_path is zero.
     */
    fn is_empty(&self) -> bool {
        self.count_path() == 0
    }

    /**
    Return the bytes of storage engine used. It was the size of values.
     */
    fn size(&self) -> usize;

    fn keys(&self, pattern: String) -> Vec<String>;

    fn entries(&self, key_pattern: String) -> Vec<(String, String)>;
}