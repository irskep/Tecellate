Game
----

* Tests
* More accurate agent communication and perception
* Food
* Better death (instead of dying just deny move request)
* Resurrection

Infrastructure
--------------

* Set up coordinators as grid
    * Requires that coordinators have the ability to pass responsibility for an agent
        * Serialize all of an agent's data (requires new API for agent)
        * Stop agent
        * Send serialized data
        * New coordinator starts new agent process
    * Some smarts or config file loading for determining partitions
* Large-scale testing

Nice To Have
------------

* Granular logging that can be configured with command line arguments
* Agents in multiple languages
