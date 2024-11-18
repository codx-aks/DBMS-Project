# Wallet System using Golang, Echo Framework, PGX, and CockroachDB

## Project Overview

This wallet system is built using **Golang**, **Echo Framework**, and **CockroachDB** for distributed storage. The system allows users to transfer wallet money (not real money) to specific vendors. Vendors can also verify the transactions through a dedicated verification route. This ensures a secure, scalable, and consistent transaction process in a distributed environment.

## Features

- **Wallet Transfers**: Users can transfer money from their wallet to a specific vendor.
- **Transaction Verification**: Vendors can verify the transaction through a separate route.
- **Distributed Database**: CockroachDB is used to store transactions, ensuring scalability, high availability, and fault tolerance.

---

## Architecture

### CockroachDB Architecture

CockroachDB is the distributed SQL database that powers this wallet system. It is designed for high scalability, high availability, and strong consistency across distributed clusters. Here's a breakdown of how CockroachDB operates:

### **1. Distributed and Horizontally Scalable**
- **Cluster**: A CockroachDB cluster consists of multiple nodes (servers) that work together to store and manage data. As your data grows, additional nodes can be added to the cluster to handle the increased load, providing seamless horizontal scalability.
- **Sharding**: Data is split into **ranges**, and each range is automatically distributed across different nodes. Each range is replicated to ensure fault tolerance. As the system scales, data is automatically rebalanced across the nodes.

### **2. Strong Consistency with ACID Transactions**
- CockroachDB guarantees **strong consistency** across all nodes in the cluster using the **Raft consensus algorithm**. This ensures that all nodes maintain the same data at any given time, even in the event of node failures.
- The system supports **ACID** transactions (Atomicity, Consistency, Isolation, Durability), which ensures that transactions like money transfers in the wallet system are fully reliable and consistent.

### **3. Transaction Handling and Parallelism**
- **Distributed Transactions**: CockroachDB supports distributed transactions, meaning a single transaction can span multiple nodes. This is ideal for systems that need to operate across multiple regions or clusters, as is the case for wallet systems managing money transfers between users and vendors.
- **Parallelism**: CockroachDB is optimized for **parallelism**. It can execute different queries and operations simultaneously across nodes to improve throughput and reduce response times. This is essential for handling large numbers of concurrent transactions in real-time.
- **Two-Phase Commit Protocol**: CockroachDB uses the **two-phase commit protocol** for distributed transactions. This ensures that a transaction either commits successfully on all involved nodes or rolls back entirely if something fails, maintaining data consistency and preventing partial commits.
- **Serializable Isolation**: CockroachDB provides **serializable isolation** for transactions, the highest level of isolation. This means that transactions are executed in a way that they behave as if they were happening one after another, preventing anomalies such as lost updates, temporary inconsistency, or phantom reads.

### **4. Fault Tolerance and Self-Healing**
- **Fault Tolerance**: CockroachDB is designed to automatically recover from node failures. If a node in the cluster goes down, the system continues to operate normally by re-replicating data to other healthy nodes.
- **Self-Healing**: In case of node failure or when a node is added to the cluster, CockroachDB automatically redistributes data and ensures it remains consistent across the entire cluster.

---