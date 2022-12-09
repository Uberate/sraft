/**
The point package the trans-protocol for remote process call.
 */
pub mod point {
    /**
    [Point] represents a trans protocol to send requests. All the [Point] has the [Server] and [Client].
    And the [Point] has some different implements. And different implements may has different configs.
    The client and SDK should define the point type and config to connection to the SRaft server.
     */
    pub trait Point{
        /**
        kind define the [Point] id, it used for config read and other SRaft server init client.
         */
        fn kind(&self) -> String;

        /**
        client return the [Point]'s [Client] by specify config.
         */
        fn client(&self) -> Box<dyn Client>;

        /**
        server return the [Point]'s [Server] by specify config.
         */
        fn server(&self) -> Box<dyn Server>;
    }

    /**
    Client is used to request to specify [Server] which type is same of [Client].
     */
    pub trait Client {}

    /**
    [Server] is used for SRaft server to listen requests from specify [Client] which type is same of
    [Server].
     */
    pub trait Server {
        fn handle(self: &Self);
    }
}

