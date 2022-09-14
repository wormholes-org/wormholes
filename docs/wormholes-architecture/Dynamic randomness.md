# Dynamic randomness 

Many blockchain consensus protocols have to allocate the creation of block creators, whose selection procedure requires a method for distributed randomness. However, computers are based on a Turing machine that is a deterministic device, and the same input seed always produces the same output sequence. Thus, computers are bad at the generation of randomness, and their outputs are pseudo-random. 

Many blockchains rely on a source of randomness for selecting participants that carry out certain actions in the blockchain process, such as consensus, shuffling, voting, and cryptographic algorithm. Distributed randomness is also crucial for many distributed applications built on the blockchain, such as building blockchain games, NFTs, and DeFI and Metaverse applications. If malicious participants are able to influence the source of randomness, they can increase their chances of being selected and getting the payout and possibly compromise the security of the blockchain and related applications. 

Finding randomness on the blockchain is difficult, as the underlying system is both open and deterministic. There are quite a few examples where randomness extracted from existing data sources have been considered, such as the blockchain variables (block difficulty, block number, block timestamp, or block gaslimit, etc.), a blockhash that expects a numeric argument that specifies the number of the block or a private seed, or published financial data, but all these variables can be manipulated by participants, so they cannot be used as a source of randomness because of the participants’ incentives. More importantly, the block variables are obviously 
shared within the same block. If a malicious participant calls the victim contract via an internal message, the same randomness in both contracts will yield the same outcome. Therefore, these schemes all suffer from the unavoidable problem that it is easy for malicious participants to determine how choices they make will affect the randomness generated on-chain and are thus vulnerable to manipulation by malicious participants. 

In all of the above scenarios, it is easy for participants to affect the result of randomness by seeing different inputs. To address this issue, various techniques have been developed, such as publicly-verifiable secret sharing (PVSS) and threshold signature schemes, as well as verifiable random functions (VRFs) such as Algorand and Ouroboros, etc., but they all suffer from requiring a non-colluding honest majority. 

Besides VRFs, there are some other popular approaches to produce randomness in the 
blockchain space. One of them is the commit–reveal scheme, such as RANDAO, in which each of the participants first privately chooses a pseudo-random number, submits a commitment to the chosen number, and all agree on a set of commitments using consensus algorithm; they then all reveal their chosen numbers, reach a consensus on the revealed numbers, and have the XOR of the revealed numbers to be the output of the protocol. It is unpredictable, but biased-able. For instance, a malicious participant can observe the network once others start to reveal their 
numbers and choose to reveal or not to reveal its number based on XOR of the numbers observed so far. This allows a single malicious participant to have one bit of influence on the output, and a malicious participant controlling multiple participants can have as many bits of influence as the number of participants they are controlling. Therefore, being able to influence or predict randomness allows a malicious participant to affect when they will be chosen to mine a block. 

During the randomness sampling, a malicious participant may re-create a block multiple times with adversarially-biased hashes (such as randomness-biasing or randomness-grinding attacks) until it is likely that the participant can create a second block shortly afterwards, which means that the malicious participant is able to bias the nonce that is used to seed the hash since the malicious participant can place arbitrary seeds in the blocks it contributes. 

To overcome this randomness-biasing, verifiable delay function (VDF), which takes a long time to generate the randomness but can quickly verify the randomness without redoing the expensive computation, provides a time delay imposed on the output of randomness. This delay prevents malicious participants from influencing the output of the randomness, since all inputs will be finalized before anyone can finish computing the VDF. VDF requires a moderate amount of sequential computation to evaluate; once a solution is found, it is easy for anyone to verify that it is correct. When used for leader selection, VDF only requires the presence of any honest participant. This added robustness is due to the fact that no amount of parallelism will speed up the VDF, and any non-malicious participant can easily verify the accuracy of anyone else’s claimed VDF output. 

The biggest problem with VDF, however, is that an attacker with very expensive specialized hardware can compute the VDF before the solution is found. This means that an attacker who has a specialized ASIC that runs the VDF faster than the time allocated to reveal commitments can still compute the final output with and without the shared information of the other parties and choose to reveal or not to reveal based on those outputs. For instance, Ethereum uses RANDAO supported by the VDF approach as their randomness beacon, but this technology is not as viable for other protocols that cannot invest in designing their own very fast ASICs that have to be 
much faster than any potential attacker’s specialized hardware. 

To address these scenarios, for instance, randomness from trusted third parties, such as the random beacon may be used. However, the additional trust assumptions and reliance on a central randomness provider, which may know the beacon values well in advance of publishing or could even manipulate the produced values without being detected, is undesirable. 

Besides common security threats, some potential risks can be foreseen. One of the most severe is the advent of quantum computing, as a functional quantum computer could easily undermine the security of the most distributed ledgers, making them practically useless. 

An unpredictable, verifiable, reliable, and fast source of randomness, therefore, is what many blockchain applications are looking for; it means the randomness must have the following properties: 
1. Unpredictable: no one should be able to predict the randomness before it is generated. 
2. Unbiased-able: the process of generating the randomness should not be biased-able by any participant. 
3. Verifiable: the validity of the generated randomness should be verifiable by any observer. 
4. Fast: generating randomness should be fast and require low computational
resources. 
5. Scalable: the protocol of randomness generation should scale to a large number of 
participants and should tolerate some percentage of participants going offline or trying intentionally to stall the protocol. 


The multi-dimensional blockchain mechanism realizes an unpredictable, verifiable, reliable, scalable, and fast source of randomness, which is what many blockchain functions, including consensus, voting, and cryptographic algorithm, are looking for. 

In the multi-dimensional blockchain mechanism, a randomness for the i-th chain at h height can be obtained by:
 
RN<sub>i</sub>(h<sub>i</sub>) = {Merkle(C<sub>i</sub>(h<sub>j</sub>)) | $\forall$ C<sub>j</sub> $\in$ G(C<sub>j</sub>(h<sub>j</sub>)), j $\in$ {0, 1, 2, …, N<sub>t</sub>-1}}, h<sub>i</sub> – 1 ≤ h<sub>j</sub> < h<sub>i</sub>, I $\in$ {0, 1, 2, …, N<sub>t</sub>-1}, h $\in$ {0, 1, 2, … } (24) 

where RN<sub>i</sub>(h<sub>i</sub>) is a randomness for the i-th chain at the height of h<sub>i</sub>, C<sub>j</sub> (h<sub>j</sub>) $\in$ {C<sub>0</sub>, C<sub>1</sub>, …} is the jth chain at h<sub>j</sub> height, and G(C<sub>j</sub>(h<sub>j</sub>)) is the set of C<sub>j</sub>(h<sub>j</sub>) for coupling with other peers, Merkle is a Merkle root function, which means that the Merkle root of the chains coupled at the height of h<sub>j</sub> – 1 can be used as a dynamic randomness for the current chain. 

The RN is unpredictable, verifiable, reliable, and fast. 
Or, it can be obtained by: 

RN<sub>i</sub>(h) = {Hash(C<sub>j</sub>(h)) | j ≠ i, C<sub>j</sub> $\in$ G<sub>b</sub>(C<sub>i</sub>(h)), j $\in$ {0, 1, 2, …, N<sub>t</sub>-1}}, i $\in$ {0, 1, 2, …, N<sub>t</sub>-1}, h $\in$ {0, 1, 2, … } (25) Or, 

RN<sub>i</sub>(h) = {Hash(C<sub>j</sub>(h) + C<sub>k</sub>(h)) | j ≠ k ≠ i, C<sub>j</sub>,C<sub>k</sub> $\in$ G<sub>b</sub>(C<sub>i</sub>(h)), j,k $\in$ {0, 1, 2, …, N<sub>t</sub>-1}}, i $\in$ {0, 1, 2, …, N<sub>t</sub>-1}, h $\in$ {0, 1, 2, … } (26) 

The dynamic randomness supports not only multi-dimensional blockchains but also other kinds of blockchains, such as sharded chains and single chain.
