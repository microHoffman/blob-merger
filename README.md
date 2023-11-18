# Blob Merger

Introducing a solution to optimize blob data usage on Ethereum, aligned with EIP-4844. As blobs maintain a fixed size, not every user fills the entire capacity. Our solution processes submitted blob data, crafting efficient blobs that maximize available space.

## Description
After EIP-4844, rollups send data to Layer 1 in blobs. They can transmit up to 128 kB of data in one blob, but often, particularly when referring to rollup blockchains with less traffic, the data volume is lower. Rollups are required to pay fees for the entire blob, even if they utilize only a portion of it.
Therefore, it is more effective and economical to merge blobs from different rollups before sending them to L1. We have developed an app that receives blobs from various rollups and merges them to create a new blob containing data from different rollups. For example, if Optimism wishes to execute a transaction containing data worth 60 kB, Arbitrum One worth 24 kB, and Boba Network worth 33 kB, these could be combined into one blob with a calldata size of 117 kB. Instead of incurring fees three times for separate blobs and utilizing three separate blobs, the data will now be consolidated into a single blob, with the fees distributed accordingly among the participants.
This solution will enhance the entire Ethereum L2 ecosystem in terms of cost efficiency and speed, by reducing the time rollups have to wait to include their data into a blob.

## Technical description
We are building on the Suave blockchain, and our SUAPP comprises two components. First, a precompiled contract that is responsible for merging the supplied data. Second, a regular app contract that manages the connection between users (typically rollups) and block builders. Users can submit their blob data to the app contract on the Suave blockchain, and builders can utilize the contract to receive optimized bundles.

Since determining the optimal blob composition is an NP problem (formally known as the bin packing problem), we have opted for a simple algorithm to showcase the concept: When blobs are sent to the precompile, they are sorted in descending order based on the calldata size. The algorithm selects the largest blob and then attempts to find the next largest blob that can still fit in, repeating this process until it traverses the entire stack. Once the stack has been searched, the merged blob is created, and the algorithm iterates through the stack again, creating another blob. This process continues until the stack is empty.

## Future plans
This is a concept of the blob merger that will be completed in the future to bring the best value to the entire L2s ecosystem. To operate at its best, the next steps need to be taken:

### Fees
Sending merged blobs to the builder requires paying the transaction fee. To ensure fairness, we propose that the fees be split among entities that have data in the merged blob, in a ratio corresponding to the data volume. For example, if Optimism wishes to execute a transaction containing data worth 60 kB, Arbitrum One worth 24 kB, and Boba Network worth 33 kB, these could be combined into one blob with a data size of 117 kB. There will still be 11 kB of free space remaining in the merged blob that cannot be filled by any blob. We suggest Optimism pay 60/117, meaning 51.28% of fees; Arbitrum One pays 24/117, meaning 20.52% of fees, and Boba Network pays 33/117, meaning 28.2% of fees.

### Splitting the data
Currently, the merger takes blobs and tries to fit them into one merged blob without splitting it into more blobs. So, when we have a merged blob of 117 kB, there is still 11 kB remaining. However, if there isn’t any blob sent by the rollup with the calldata size of 11 or smaller, it remains empty. For greater effectiveness, we suggest splitting the calldata between merged blobs. For example, if there are merged blobs of sizes 117 kB and 123 kB, and in the stack is still the rollup blob of size 13 kB, 11 kB will be added to the first blob, making it 128 kB in size, and the remaining 2 kB to the second blob, making it 125 kB in size.

### Communication back to the user
The Future Blob Merger will notify the rollup that is sending the data to the smart contract that it has been successfully sent to the builder.
