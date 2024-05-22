# Vector Indexes

This repository contains a Vector Indexes for efficient ANN search.
Embedding models represent their vectors as float32.

## Implementations
- Flat Index
- Index with Product Quantization
- IVF Index
- IVF PQ Index
- HNSW
- Vamana


## Related Publications
* Malkov, Y.A., Yashunin, D.A.. (2016). [Efficient and robust approximate nearest neighbor search using Hierarchical Navigable Small World graphs. CoRR](http://arxiv.org/abs/1603.09320)
* Malkov, Y., Ponomarenko, A., Logvinov, A., & Krylov, V., 2014. [Approximate nearest neighbor algorithm based on navigable small world graphs.](http://www.sciencedirect.com/science/article/pii/S0306437913001300)
* A. Ponomarenko, Y. Malkov, A. Logvinov, and V. Krylov  [Approximate nearest neighbor search small world approach.](http://www.iiis.org/CDs2011/CD2011IDI/ICTA_2011/Abstract.asp?myurl=CT175ON.pdf)
* H. Jegou, M. Douze, and C. Schmid [Product Quantization for Nearest Neighbor Search](https://ieeexplore.ieee.org/document/5432202/)
* T. Ge, K. He, Q. Ke, and J. Sun [Optimized Product Quantization](https://ieeexplore.ieee.org/document/6678503/)
* Y. Matsui, Y. Uchida, H. Jegou, and S. Satoh [A Survey of Product Quantization](https://www.jstage.jst.go.jp/article/mta/6/1/6_2/_pdf/)
* Suhas Jayaram Subramanya, Devvrit, Rohan Kadekodi, Ravishankar Krishnaswamy, and Harsha Simhadri. 2019 [DiskANN: Fast Accurate Billion-point Nearest Neighbor Search on a Single Node](https://proceedings.neurips.cc/paper_files/paper/2019/file/09853c7fb1d3f8ee67a61b6bf4a7f8e6-Paper.pdf)