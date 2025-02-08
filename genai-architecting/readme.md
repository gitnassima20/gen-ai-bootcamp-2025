## 1. Functional Requirements:

The company wants to own its infrastructure to control costs and ensure data privacy.

*Budget*: $10-15K investment in an AI PC for model hosting.

*Privacy Concern*: Avoid sending student data to external cloud providers.

*Users*: 300 active students, all located in Osaka.

*Cloud Consideration*: Future adoption of cloud for deployment and performance scaling.


## 2. Assumptions:

- We are assuming that the Open-source LLMs that we choose will be powerful enough to run on hardware with an investment of 10-15K.

- We're just going to hook up a single server in our office to the internet and we should have enough bandwidth to serve the 300 students.

- We excpect that our infrastructure should handle latency applications without immediate cloud dependency.


## 3. Data Strategy

*Data Collection*: We are considering synthetic data for pretraining.

*Copyright Compliance*: We want to purchase and store legally approved materials.

*Storage Strategy*: Implement a cost-effective local database for fast retrieval.

*Inference Optimization*: Consider vector databases for efficient model querying.


## 4. Model Selection

*Primary Choice*: IBM Granite (Open-source with transparent training data: https://huggingface.co/ibm-granite).


## 5. Deployment & Maintenance Plan

+ Initial Setup:

    - Install AI PC with necessary dependencies.

    - Deploy IBM Granite LLM.

+ Monitoring:

    - Implement logging for performance tracking.

    - Set up alerts for hardware resource usage.

+ Backup & Recovery:

    - Automate backups.

+ Security Measures:

    - Implement encryption for sensitive data.

+ Cost Optimization:

    - Use power-efficient hardware to minimize electricity costs.

    - Optimize batch processing to reduce computational load.

