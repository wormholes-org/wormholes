# Multi-Dimensional Blockchain Compatibility

## Scalable and Flexible Architecture for Multiple Chains and Multiple Layers

  The multi-dimensional blockchain mechanism, being inherently flexibile and scalable, supports many kinds of blockchain structures, including single layer, 
  multiple layers, parallel independent chains, core chaining chains, homogenous chains, heterogeneous chains, mixed heterogenoushomogenous chains, or any 
  combination of the above. 

  The multi-dimensional blockchain architecture may consist of an index chain and several core chains and subchains. The index chain may contain general 
  information about the protocol and the current values of its parameters, the set of validators and their stakes, the set of currently active core chains, 
  and, most importantly, the set of hashes relating to all core chains and subchains. The core chains may contain the value-transfer and smart-contract 
  transactions. The multiple core chains can be homogeneous or heterogeneous.

## Example of Homogeneous Multiple Chain System

  One example of homogeneous systems is shown in Figure, in which all core chains may use the same format for transactions, the same virtual machine for 
  executing smart-contract code, share the same cryptocurrency, and so on, and this similarity is explicitly exploited but with different data in each core 
  chain.
  
  <img width="648" alt="图片" src="https://user-images.githubusercontent.com/107660058/189838470-48ffe117-beba-4e2d-ad8a-a9a1cb226dbe.png">

## Example of Heterogeneous Multiple Chain System

  Example of heterogeneous systems is shown in Figure, in which different core chains may have different rules, meaning different formats of account 
  addresses, different formats of transactions, different virtual machines (VMs) for smart contracts, different basic cryptocurrencies, and so on. 
  However, they all must satisfy certain basic interoperability criteria to make the interaction between different core chains possible and relatively 
  simple.
  
  <img width="631" alt="图片" src="https://user-images.githubusercontent.com/107660058/189838568-e58ddae8-4e6b-4d11-933c-b967765e654a.png">

## Example of Multiple Subchains

  The subchains are subdivided by each core chain, have the same rules and block format as the core chain itself, but are responsible only for a subset of 
  data, such as accounts which depend on several first (most significant) bits of the account address. Because all these subchains share a common block 
  format and rules, the multi-dimensional blockchain is homogeneous in this respect. 

  A multi-dimensional blockchain system contains one index chain and can potentially accommodate up to 232 or more core chains; each core chain can be 
  divided into up to 232 or more subchains.
  
  Figure shows one example of the multi-dimensional blockchain system, in which the mainchain uses a single block structure that contains all the 
  transactions from all the core chains. Because the main blockchain is divided into core chains, the existence of the blockchain is “virtual”, meaning 
  that it is a collection of core chains. The blockchain has only a purely virtual or logical existence inside the core chains. It is a coincidence that 
  the potentially huge number of core chains can be grouped into one blockchain.
  
  <img width="613" alt="图片" src="https://user-images.githubusercontent.com/107660058/189825233-63691f51-7fec-416d-ac5d-5bf704b15590.png">

  In this multi-dimensional blockchain architecture, accounts or smart contracts are grouped into core chains. Then there are several core chains, each 
  describing the state and state transitions of only a group of accounts, sending messages to each other to transfer value and information. 

  Practically, there is only one index chain but core chains may be up to 214, with each responsible only for a set of accounts depending on the several 
  first (most significant) bits of the account address, as shown in Figure:
  
  <img width="636" alt="图片" src="https://user-images.githubusercontent.com/107660058/189825419-ac129b81-9c6c-4364-a57b-595926e33b08.png">

  In this case, all the transactions in the block are divided into physical core chains, and each block contains either one or zero specific core chain 
  block, as shown in Figure:
  
  <img width="214" alt="图片" src="https://user-images.githubusercontent.com/107660058/189825567-89b43165-4d4c-4695-a295-0f64a20ca0b6.png">

  In this system, no node of the blockchain network needs to download the full state or the full block. Each node only maintains the state that corresponds
  to the core chains that they validate transactions for, and the whole state of the system is changed in all the core chains.

## Example of Homogeneous Multiple Chain and Multiple Layer System

  An example of a homogeneous multiple chain and multiple layer systems with main chains X and layer chains Y, X $\in$ {1, 2, …} and Y $\in$ {1, 2, …}, is 
  shown in Figure, in which all chains may use the same format of transactions, the same virtual machine for executing smart-contract code, and so on, 
  and this similarity is explicitly exploited but with different data in each chain. Layer chain Y refers to the intended scaling solutions, such as 
  protocols or networks, that operate atop an upper layer blockchain, essentially functioning as different layers of blockchain. This architecture allows 
  several layers of blockchains to be built on top of each other, linking them through a parent-child connection. The parent layer chain assigns and 
  distributes the tasks among its children chains, which in turn execute them and send back the result to the mainchain, relieving the parent chain of its 
  workload and increasing scalability. The multiple layer architecture can make transactions faster and lower gas fees while solving scaling problems.
  
  <img width="945" alt="图片" src="https://user-images.githubusercontent.com/107660058/189837478-34b92d61-8fe6-4d15-8276-58758f186ca3.png">

## Example of Heterogeneous Multiple Chain and Multiple Layer System

  An example of a heterogeneous multiple chain and multiple layer systems with main chains X and layer chains Y, X $\in$ {1, 2, …} and Y $\in$ {1, 2, …}, 
  is shown in Figure, in which different chains may have different rules, meaning different formats of account addresses, different formats of transactions, 
  different virtual machines (VMs) for smart contracts, different basic cryptocurrencies, and so on. Layer chain Y operates atop an upper layer blockchain,
  allowing several layers of blockchains to be built on top of each other heterogeneously, linking them through a parent-child connection. The 
  heterogeneous multiple chain and multiple layer architecture can make transactions faster and lower gas fees.
  
  <img width="944" alt="图片" src="https://user-images.githubusercontent.com/107660058/189838135-01bc359a-f236-46bf-bcba-044f05f0dedd.png">
