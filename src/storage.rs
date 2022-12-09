/**
Mod storage define the interface to deal the data storage problem. In different project, the
storage ability is different. Such like high change with less data requirement, the memory is a
good idea. But the Big-Data, the disk is only way.
 */
pub mod storage {
    use crate::config::ConfigAble;

    pub trait StorageEngine {
        fn new(&self) -> Box<dyn StorageEngine>;

        fn config(&mut self, config: ConfigAble);
    }
}