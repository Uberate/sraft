/**
Mod storage define the interface to deal the data storage problem. In different project, the
storage ability is different. Such like high change with less data requirement, the memory is a
good idea. But the Big-Data, the disk is only way.

Every one who use this repo can custom themselves storage implements. For the repo organized and
layout, I suggest to implements them in `plugins` dir.
 */

use crate::config::ConfigAble;
use crate::plugins;

/**
Get a storage engine by specify type and version.

If not found specify type and version, return [None].

### Implementations list
| Type | Version| Instance |
| :--- | :---: | :--- |
| `memory` | `v1` | [plugins::storage::memory_v1::Engine] |


 */
pub fn new_engine(typ: String, version: String, config_able: &ConfigAble) -> Option<Box<dyn StorageEngine>> {
    let id = (typ.as_str(), version.as_str());

    match id {
        ("memory", "v1") => { Some(plugins::storage::memory_v1::Engine::new(config_able)) }
        _ => { None }
    }
}

/**
StorageEngine is a interface to define the storage ability sraft need. The storage should
provide the Key-Value save. The key was the path.
 */
pub trait StorageEngine {
    fn set(&mut self, path: String, value: String);

    /**
    Returns a reference to the value corresponding to the key.
     */
    fn get(&self, path: &String) -> Option<&String>;

    /**
    Return a clone value to the value corresponding to the path.
     */
    fn get_mut(&mut self, path: &String) -> Option<& mut String>;

    /**
    Removes a path from the storage engine, returning the value at the path if the path was
    previously in the storage engine.
     */
    fn remove(&mut self, path: &String) -> Option<String>;

    /**
    Removes a path from the map, returning the stored path and value if the path was previously in
    the storage engine.
     */
    fn remove_entry(&mut self, path: &String) -> Option<(String, String)>;

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

    /**
    Return the key list.
    // TODO: It may be a iterator batter.
     */
    fn keys(&self, pattern: String) -> Vec<String>;

    /**
    Return the entry list.
    // TODO: It may be a iterator batter.
     */
    fn entries(&self, key_pattern: String) -> Vec<(String, String)>;
}