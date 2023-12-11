# Test and experiments to verify the clients, servers and methodology of Iris

The experiments to test and evaluate Iris consisted of many different parts, from evaluating the algorithms to evaluate the thresholds use for the clients and servers.

# Experiment on Optimal Epidemic Spreading Strategy
The Iris P2P system proposes new ways to deal with the security concerns of threat intelligence sharing. However, many of these ideas need to be verified and explored in simulated experiments in order to understand how a network of peers may behave under different conditions. We simulate and evaluate different spreading strategies in an epidemic sense. As a result of this experiment, we choose the optimal spreading strategy for each different condition and for the protocols for spreading Alerts and for spreading file metadata with different levels of severity.

We decided to run one very large experiment with many parameters instead of a large number of small experiments. In this way it was possible to validate how the conditions relate to each other.

## Goal
The goal of this experiment is to evaluate different spreading strategies based on how fast they spread the message in networks and how much they flood networks with messages. To define a spreading strategy, we need to answer five questions: 

- What algorithm of spreading to use?

  In our case, only the Push algorithm makes sense. The others assume that gossip is constantly updated in the system, and thus peers proactively ask other peers for updates. That is not true in our environment - it might happen that no gossip is spread for a long time. In that case, the proactive update messages in Pull and Push/Pull algorithms are unnecessary.

- To how many nodes should infected nodes spread the message at once?

  We will call this value a spreading factor. Finding the optimal spreading factor is one of the goals of the experiment.

- How long should the infected nodes wait until they ask/spread again?

  We will call this value spreading period. Finding the optimal spreading period is one of the goals of the experiment.

- How do infected nodes choose recipients for spreading from a list of candidates?

  By choosing the recipients completely randomly, we can use the service trust of candidates to design a better heuristic. We propose 3 options:

  - Choose recipients with a uniform probability. Most of the Epidemic Protocols employ this option.
  - Sort recipients based on their service trust and choose first the most trusted ones. This technique assumes that peers with bigger service trust will more likely follow the protocol and thus contribute more to the dissemination of messages.
  - Choose recipients with a probability that exponentially grows with the service trust of each recipient. This technique was already described when we talked about choosing recipients in the Network Opinion Protocol. It works firstly by transforming service trust values of all candidates using an exponential function $f(st) =\frac{a^{st}-1}{a-1}$ with a constant value $a=10$. After that, Iris normalizes the transformed values which represent the probabilities of the candidates being chosen. This technique still favors the trusted peers but also allows one to choose less trusted peers for a bigger variety. The reason is that the previous option (sort by service trust) might result in flooding the trusted peers and not spreading the message to less trusted parts of the network which may also contain benign peers.
- After how much time should the infected node move to a Removed state?

  We will call this value spreading expiration. Optimally, the spreading should stop after the gossip reaches the entire network. Nonetheless, our network is decentralized and does not provide complete information about the status of all peers. For this reason, the spreading expiration should be set beforehand to a duration that guarantees the full convergence. Note that the spreading expiration depends on the previous points because they essentially define the speed of spreading.

As a consequence, we define a _spreading strategy_ as a triplet of:
1. Spreading factor $\in \mathbb{N}$
2. Spreading period $\in \mathbb{N}$
3. Algorithm to choose recipients - we define three options:
  - Choose recipients with uniform probability
  - Choose first the peers with the biggest service trust
  - Choose recipients with a probability that exponentially grows with the service trust of recipients

Finally, after we evaluate the spreading strategies, we should be able to choose an optimal strategy for following types of messages:

- Alerts: Alert messages should be spread rather quickly to their recipients. The reason is that the event that has caused the alert is usually relevant to the given moment and can help other peers to prepare for an incident.
- File Metadata: Each metadata message has either NORMAL or CRITICAL severity. Each severity level should guarantee different spreading properties. Metadata messages with CRITICAL severity should spread faster to their recipients even for the cost of flooding the network with messages. On the other hand, metadata messages with NORMAL severity do not have to be spread as fast, and thus the peers should choose a spreading strategy that minimises the magnitude of flooding the network with messages.

# Assumptions
Our experiment makes some important assumptions. First, we make assumptions about the Fides trust model. We assume that, in the long run, the average service trust of all malicious peers is smaller than the average service trust of all benign peers. This assumption essentially means that Fides behaves correctly.

Also, we think it is rational to assume that Fides does not provide to peers estimates of service trust values that 100\% match the ground-truth values. We assume that an error of the estimation has a mean $\mu=0$ and a standard deviation $\sigma=0.25$.

Other assumptions are about the network itself. We assume that the churn rate of the network is practically zero. It means we assume that the network is static - no peers join or leave the network. Real networks probably never have a churn rate of zero, but we think that our network will not be too dynamic. The reason is that once a peer joins, it is in his best interest to stay connected for the benefit of collective defense.

Lastly, we assume that the underlying graph of our network is connected. However, the connectivity of the underlying graph is difficult to verify in the real world because we designed our network as an almost unstructured (DHT k-buckets form a small structure) and pure P2P network.





























  

