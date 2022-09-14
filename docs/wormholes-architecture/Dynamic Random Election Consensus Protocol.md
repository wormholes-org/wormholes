  In a consensus protocol, all participant nodes distributed in a blockchain network share transactions and agree on the integrity of the shared 
  transactions. For instance, Proof of Work (PoW) requires exhaustive computational power from participants for block generation, and Proof of Stake (PoS) 
  uses participants’ stakes for generating blocks. With extensive research, many consensus algorithms have been developed to improve the consensus 
  confirmation time and power consumption of blockchain-powered distributed ledgers. Our DRE protocol describes a general model of asynchronous BFT (aBFT) 
  consensus protocols.
  
 # Dynamic Random Election
  
  Dynamic Random Election (DRE) consensus is a set of practical Byzantine fault-tolerant (BFT) protocols for completely asynchronous environments. A 
  synchronous BFT system utilizes broadcast voting and asks all nodes to vote on the validity of each block while an asynchronous BFT system achieves a 
  local view with a high probability of being a consistent global view according to the concepts of distributed common knowledge. The definitions of the 
  DRE protocol are given below:
  
  Node: a machine participating in the DRE protocol. Each node has a local state consisting of local histories, messages, proposal blocks, and peer 
  information.
  
  Related-to: the relationship between nodes which have proposal blocks. If there is a path from a proposal block x to y, then x related-to y. 
  “x related-to y” means that the node creating y knows proposal block x.
  
  Lamport timestamp: for topological ordering, a Lamport timestamp algorithm uses the relatedto relation to determine a partial order of the whole proposal
  block based on logical clocks.
  
  DAG: the history of proposal blocks forming a Directed Acyclic Graph (DAG). DAG G = (V, E) consists of a set of vertices V and a set of edges E. A path 
  in G is a sequence of vertices (v1, v2, …, vk) by following the edges in E such that it uses no edge more than once. Each vertex vi V is a proposal block.
  An edge (vi, vj) E refers to a hashing reference from vi to vj; that is, vi ↪ vj. 
  A DAG is held by each node to identify topological ordering between proposal blocks, to select pilot block candidates, and to compute consensus time of 
  avowal blocks and proposal blocks under its subgraph.
  
  Subgraph: for a vertex v in a DAG G, let G[v] = (Vv, Ev) denote an induced-subgraph of G such that Vv consists of all ancestors of v including v, and Ev 
  is the induced edges of Vv in G.
  
  Proposal block: a tuple s, a, s’ consisting of a state, an action, and a state. The j-th proposal block in history hi of process i is si j-1, a, 
  sij, denoted by vij. Nodes can create proposal blocks. A proposal block includes the generation time, signature, transaction history, and reference to 
  parent proposal blocks. 
  
  Frame: contains a disjoint set of trunk blocks and proposal blocks, in which the history of proposal blocks is divided into frames.
  
  State: a (local) state of node i is denoted by sij consisting of a sequence of proposal blocks sij=vi0, vi1, …, vij. In a DAG-based protocol, each 
  proposal block vij is valid only if the reference blocks exist before it. A local state sij is corresponding to a unique DAG. In a DAG, we simply denote 
  the j-th local state of a node i by the DAG gij. Let Gi denote the current local DAG of a process i.
  
  Local history: a sequence of local states starting with an initial state, denoted by hi. A set Hi of possible local histories for each process i. A 
  process’s state can be obtained from its initial state and the sequence of actions or proposal blocks that have occurred up to the current state. The 
  DRE protocol uses append-only semantics. In DRE, a local history is equivalently expressed as: hi = gi0, gi1, gi2, gi 3 … where gij is the j-th local 
  DAG (local state) of the process i.
  
  Discharge: each asynchronous discharge is a vector of local histories, denoted by d = h1, h2, h3, …hN. Let D denote the set of asynchronous discharges.
  A global state of discharge d is an nvector of prefixes of local histories of d, one prefix per process. The related-to relation can be used to define a 
  consistent global state, often termed a consistent cut.
  
  Trunk block: a proposal block if either it is the first proposal block of a node, or it can reach more than 2/3 of the blockchain network’s validating 
  power from other trunk blocks. A trunk block set Ts contains all the trunk blocks of a frame. A frame f is a natural number assigned to trunk block sets 
  and its dependent proposal blocks. The set of all first proposal blocks of all nodes forms the first trunk block set T1 (|T1| = n). The trunk block set 
  Tk consists of all trunk blocks ti such that ti $\notin$ Ti, $\forall$ i $\in$ {1, ..., (k-1)} and ti can reach more than 2/3 validating power from other 
  trunk blocks in the current frame, i $\in$ {1, ..., (k-1)}.
  
  Pilot block: a trunk block at layer i that is known by a trunk block of higher frames (i + j), j {(i + 1), (i +2), ….}
  
  Avowal block: a pilot block assigned with a consensus time.
  
  Groupuscule: a list of avowal blocks and the subgraphs reachable from those avowal blocks.
  
The core idea of DRE is the DAG and dynamic randomness for an election. In the DRE protocol, a node can create a new proposal block, which has a set of 2 
to k parents. Nodes generate and propagate proposal blocks asynchronously, and the DRE algorithm achieves consensus by confirming how many nodes verify 
the proposal blocks.

The DAG is used to compute special proposal blocks, such as trunk block, prepared blocks, and verified blocks. The Groupuscule consists of ordered verified
block proposal blocks, which can maintain reliable information between proposal blocks. The DAG and Groupuscule are updated with newly generated proposal 
blocks frequently and can quickly and forcefully respond to attacks.

In the DRE protocol, each node can create messages and send messages to and receive messages from other nodes. The communication between nodes is 
asynchronous. Each node stores a DAG, which is the DAG of proposal blocks. A block has some edges to its parent proposal blocks. A node can create a 
proposal block after the node communicates the current status of its DAG with its peers.

The DAG is a graph structure stored on each node. The DAG consists of proposal blocks and references between them as edges.

The local DAG is updated quickly as each node creates and synchronizes proposal blocks with each other. For high-speed transaction processing, proposal 
blocks are assumed to arrive at very high speeds asynchronously. Let G = (V, E) be the current DAG and G’ = (V’, E’) denote the diff graph, which consists 
of the changes to G at a time, either at proposal block creation or arrival. The vertex sets V and V’ are disjoint, similar to the edge sets E and E’. At 
each graph update, the updated DAG becomes Gnew = (V ∪ V’, E ∪ E’).

Each node uses a local view of the DAG to identify the trunk block, prepared block, and verified block vertices and to determine the topological ordering 
of the proposal blocks.

Figure shows a general framework of the DRE consensus protocol.

<img width="199" alt="图片" src="https://user-images.githubusercontent.com/107660058/190090220-b05d46be-906e-499a-a020-fddcd42d11b0.png">

Based on the DRE consensus protocol, each node contains a DAG consisting of proposal blocks and stores the information of accounts and their stakes.

The key steps in the DRE consensus protocol include the following: 

  1. Proposal block creation 
  2. Computing validation score 
  3. Selecting trunk blocks and updating the trunk block sets 
  4. Assigning weights to new trunk blocks 
  5. Deciding frames 
  6. Deciding prepared/verified blocks 
  7. Ordering the final blocks

The DRE protocol supports dynamic participation so that all participants can join the blockchain network.

# PoS + DRE Protocol Fairness and Security

  The fairness and security benefits of the Wormholes Blockchain PoS + DRE protocol (hereinafter referred to as the “Wormholes Consensus Protocol”) are 
  highlighted here.
  
## Fairness
  
  PoW protocol is fair because a miner with pi fraction of the total computational power can create a block with the probability pi. PoS protocol is fair 
  too because an individual node with si fraction of the total stake or tokens can create a new block with si probability. However, initial holders of 
  tokens in PoS systems tend to keep their tokens in order to gain more rewards.
  
  The Wormholes Consensus Protocol is fair because every node has a chance to create a proposal block equally. Like other PoS protocols, any node in the 
  Wormholes Consensus protocol can create a new proposal block with a stake-based probability. Unlike PoW protocols, nodes in the Wormholes Consensus 
  Protocol do not require expensive hardware.
  
  Like a PoS blockchain system, it is a possible concern that the initial holders of tokens will not have an incentive to release their tokens to third 
  parties, as the token balance directly contributes to their wealth. Unlike existing PoS protocols, each node in the Wormholes Consensus Protocol is 
  required to validate parent proposal blocks before it can create a new block. Thus, the economic rewards a node earns through proposal block creation 
  is, in fact, to compensate for their contribution to the onchain validation of past proposal blocks and its new proposal block. Remarkably, the Wormholes
  Consensus Protocol is more intuitive because our reward model used in stake-based validation can lead to a more reliable and sustainable network.
  
  <img width="675" alt="图片" src="https://user-images.githubusercontent.com/107660058/190091087-9b9d864e-9127-4129-b041-b5a480788c4c.png">

  The table above provides a comparison of PoW, PoS, and the Wormholes Consensus Protocols, where pi is the computation power of a node and P is the total 
  computation power of the blockchain network, si is the stake of a node, S is the total stake of the whole network, and n is the number of nodes.
  
## Security

  The table below shows a comparison between existing PoW, PoS, DPoS, and the Wormholes Consensus Protocols, corresponding to the effects of common types 
  of attacks. Note that, existing PoS and DPos protocols have addressed some of the known vulnerabilities in one way or another. Obviously, the Wormholes 
  Consensus Protocol is more secure than PoW, PoS, and DPoS.
  
  <img width="789" alt="图片" src="https://user-images.githubusercontent.com/107660058/190091535-af11b38f-9e12-428a-bd40-e77bcf42b195.png">
  
  PoW-based systems are facing selfish mining attack, where an attacker reveals mined blocks selectively in order to waste computational resources of 
  honest miners. 

  Both PoW and PoS share some common vulnerabilities, such as Sybil attack and DoS attack. In a Sybil attack, the attacker creates multiple fake nodes to 
  disrupt the blockchain network. In a DoS attack, the attacker disrupts the blockchain network by flooding the nodes. In a Bribe attack, which is another 
  shared vulnerability, the attacker obtains the majority of computational power or stake through bribing for a limited duration. PoS is more vulnerable 
  because a PoS Bribe attack costs much less than a PoW Bribe attack. 

  Furthermore, PoS has more weaknesses that are not relevant in PoW, including the following:
  
  Sabotage: an attacker owning W/3 + smin of the stakes can appear offline by not voting and hence checkpoints and transactions cannot be finalized, where 
  smin is minimum number of tokens that can be staked by an account. Users are expected to coordinate outside of the blockchain network to censor the 
  malicious validators. 

  Grinding Attack: in the blockchain validation process, if validators are able to manipulate the random election process of a consensus algorithm, 
  grinding attacks happen. For example, when the election process depends on the block hash, the node elected as a leader for adding a block can manipulate
  the block-hash by adding or removing certain transactions or trying many parameters to form different block hashes. This node, hence, can grind through 
  many different combinations and choose the one that reelects itself with a high likelihood for the next block. 
  This is a serious potential source of threat to a blockchain consensus protocol based on ‘proof-ofstake’ coins. Nothing at Stake Attack: a malicious node
  can mine on an alternative chain in PoS at no cost, whereas it would lose CPU time if working on an alternative chain in PoW.
  
  Double-Spending: this is the first and the most common attack introduced in the field of cryptocurrencies. In such a scenario, an attacker may spend a 
  given set of coins in more than one transaction. There are several ways to perform a double-spending attack, including the following: 

  Pre-mine one transaction into a block and spend the same coins before releasing the block to invalidate that transaction (called a Finney attack). 

  Send two conflicting transactions in rapid succession into the blockchain network (called a Race attack). 

  Response to Attacks
  
  Like other decentralized blockchain technologies, the Wormholes Blockchain network may also face potential attacks by attackers. Here, we describe 
  several possible attack scenarios and show how the Wormholes Consensus Protocol can be used to avoid such attacks.
  
  Attack cost: due to the scarcity of native tokens on a PoS blockchain relative to a PoW blockchain, it will cost more for an attack on the Wormholes 
  Consensus Protocol since it is based on PoS mechanism. In order to gain more stake for an attack in the Wormholes Consensus Protocol, it will be very 
  costly for an outside attacker. The Wormholes Consensus Protocol can effectively avoid potential attacks as attackers will need to gather 2W/3 (or W/3 
  for certain attacks) tokens to control the validation process, where W is the total number of all tokens staked on the blockchain, regardless of the 
  token price. Acquiring more tokens will increase the price of tokens, leading to a massive cost. Another risk to an attacker is that all its tokens will 
  be burned following any detected attempt.
  
  Denial of Service: the Wormholes Consensus Protocol is leaderless, requiring 2n/3 participation. An attacker would have to attack more than n/3 (or more 
  than W/3) validators to be able to successfully mount a DDoS attack. 

  Bribery Attack: since 2n/3 participating nodes (or nodes with a total of at least 2W/3) are required to validate in the Wormholes Consensus Protocol, 
  this would make it necessary for an attacker to bribe more than n/3 of all nodes (or more than W/3) to launch a bribery attack, which is extremely 
  expensive. 

  Sybil: Each participating node must stake a minimum amount of ERB tokens to take part in the Wormholes Consensus Protocol. However, staking 2/3 of the 
  total stake would be prohibitively expensive. 

  Transaction Flooding: An attacker may discharge a large number of valid transactions from their accounts with the purpose of overloading the blockchain 
  network. To avoid such attempts, the Wormholes Blockchain network requires a minimal transaction fee that is reasonable for normal users, but will cost 
  the attackers a significant amount if they send a large number of transactions. 

  Parasite Chain Attack: A malicious node can make a parasite chain and attempt to make a malicious proposal block. In the Wormholes Consensus Protocol, 
  when consensus is reached on a finalised proposal block, it will be verified before it is added into the Groupuscule.
  
  The Wormholes Consensus Protocol is based on the BFT mechanism, which is secure unless 1/3 or more of the nodes are dishonest or malicious. The malicious
  nodes may create a parasite chain. As trunk blocks are known by nodes with more than 2W/3 validating power, a parasite chain can only be shared between 
  malicious nodes, which account for less than one-third of participating nodes. Nodes with a parasite chain are unable to generate trunk blocks, and they 
  will be detected by the avowal blocks.
  
  Double spending: a double spending attack is when malicious attackers attempt to spend their funds twice. A blockchain fork opens the door to double 
  spending attacks. With the Wormholes Consensus Protocol, any honest node can structurally detect the fork. Because honest validators are incentivized to 
  detect and avoid such misbehavior, attackers will be caught and their stakes will be burned. 

  Sabotage Attack: consider a Sabotage Attack, when attackers may refuse to vote for a period of time. In the Wormholes Consensus Protocol, the absence of
  validators for a period will be discovered by honest validators. The absent validators will be pruned from the blockchain network and will not be counted 
  further, allowing the blockchain network to continue normal operation with the rest of the nodes. 

  Grinding Attack: a grinding attack happens when a validator is able to manipulate the random election process of a consensus algorithm. The Wormholes 
  Consensus Protocol ensures that the randomness of its random election process is unpredictable and unbias-able at commitment time; therefore, a stake 
  grinding attack becomes unfeasible.
  
  Long-Range Attack: in some blockchains, if the fork chain is longer than the original, the longer chain will be accepted as the main chain since the 
  longer chain will have done more work (or have a large stake) since its creation; then in these circumstances, an attacker may take the chance to create 
  a separate fork chain. This kind of attack is not possible in the Wormholes Blockchain network because the Wormholes Consensus Protocol is fork-free, and
  adding a new block into the block chain requires an agreement of 2W/3. To accomplish a long-range attack, attackers will have to gain validating power of
  2W/3 to create the new chain.
  
  
