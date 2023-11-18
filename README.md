## Blob Merger

Introducing a solution to optimize blob data usage on Ethereum, aligned with EIP-4844. As blobs maintain a fixed size, not every user fills the entire capacity. Our solution processes submitted blob data, crafting efficient blobs that maximize available space.

# Description
After EIP-4844, rollups send data to Layer 1 in blobs. They can transmit up to 128 kB of data in one blob, but often, particularly when referring to rollup blockchains with less traffic, the data volume is lower. Rollups are required to pay fees for the entire blob, even if they utilize only a portion of it.
Therefore, it is more effective and economical to merge blobs from different rollups before sending them to L1. We have developed an app that receives blobs from various rollups and merges them to create a new blob containing data from different rollups. For example, if Optimism wishes to execute a transaction containing data worth 60 kB, Arbitrum One worth 24 kB, and Boba Network worth 33 kB, these could be combined into one blob with a calldata size of 117 kB. Instead of incurring fees three times for separate blobs and utilizing three separate blobs, the data will now be consolidated into a single blob, with the fees distributed accordingly among the participants.
This solution will enhance the entire Ethereum L2 ecosystem in terms of cost efficiency and speed, by reducing the time rollups have to wait to include their data into a blob.

# Technical description
