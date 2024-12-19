This is project 3 in the course COP5615 Distributed Operating System Principles.

In this project I implemented the Chord protocol for a peer-to-peer distributed hash table. The project is exclusively written in Pony with the help of the actor model functionality.

The simulation models a peer-to-peer network with nodes organized in a ring, where each node can lookup keys, join, and manage routing information. The Chord protocol is known for its scalability, with logarithmic performance on node lookups and hops.

Overview:
The simulation initializes a Chord network with a specified number of nodes, each responsible for locating a portion of the distributed key space. Each node maintains a finger table, which helps in efficiently routing requests to find the node responsible for a given key. Nodes periodically issue requests to locate successors and resolve keys, with the average hop count being reported as a measure of performance.

The specification of the Chord protocol can be found in the paperChord: A Scalable Peer-to-peer Lookup Service for Internet Applicationsby Ion Stoica,  Robert  Morris,  David  Karger,  M.  Frans  Kaashoek,  Hari  Balakrishnan. https://pdos.csail.mit.edu/papers/ton:chord/paper-ton.pdfLinks to an external site..  You can also refer to the Wikipedia page: https://en.wikipedia.org/wiki/Chord(peer-to-peer)Links to an external site. The paper above, in section 2.3 contains a specification of the Chord API and of the API to be implemented by the application.

Requirements
You have to implement the network join and routing as described in the Chord paper (Section 4) and encode the simple application that associates a key (same as the ids used in Chord) with a string.  You can change the message type sent and the specific activity as long as you implement it using a similar API to the one described in the paper.

Input: The input provided will be of the form:

project3 numNodes  numRequests

Where numNodesis the number of peers to be created in the peer-to-peer system and numRequests is the number of requests each peer has to make.  When all peers performed that many requests, the program can exit.  Each peer should send a request/second.

Output: Print the average number of hops (node connections) that have to be traversed to deliver a message.

Actor modeling: In this project, you have to use exclusively the  actor framework (projects that do not use multiple actors or use any other form of parallelism will receive no credit).  You should have one actor for each of the peers modeled.
