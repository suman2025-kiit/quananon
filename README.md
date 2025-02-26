# quananon_sam

## Quantum_AnonCreds_Dew_Edge_Quantum_Self_Attention_Mechanism
#### Initial setup for dynamic identity generation, 160-bit address generation , 4-Qubits GHZ Bell states generation via smart contract and chaincode validation 

We utilized Ubuntu Bionic 18.04.4, Pandas, Python 3.6, Jupyter Notebook, and Hyperledger AnonCreds 1.0 to develop and validate Dynamic Identities (DID), 160-bit dynamic addresses, and Verifiable Credentials (VC) for participants in the Healthcare CHoT within a Quantum-Blockchain distributed network. AnonCreds 1.0, integrated with the Hyperledger Indy platform, facilitates the creation of an Indy System Pool. This setup assigns verifier roles such as Trust Anchor (role '101') and Trustee (role '0') using Nym Transactions in our 'DID_WalletAddress_Generator.py' script.

To start Docker, used the command: docker-compose up -d. To stop Docker after the experiment, use: docker-compose down.

For installation, AnonCreds 1.0 works alongside the Hyperledger Indy platform to generate an Indy System Pool. This pool assigns predefined verifiers like Trust Anchors (role '101') and Trustees (role '0'), through Nym Transactions managed by our 'DID_WalletAddress_Generator.py' script. Navigate to the directory with the command: cd fabric.

./ install.sh

during Installation of Fabric, also please confirm that the following commands will run successfully /network down

./network up

Add mychannel:

./network.sh createChannel -c mychannel

Add simple chaincode to channel:

./network.sh deployCC -ccn simple -ccp /home/suman/Downloads/caliper-benchmarks/src/fabric/samples/marbles/go -ccl go

Invoke chain code:

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n simple --peerAddresses localhost:7051 --tlsRootCertFiles

${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles

${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

in test-networkfolder (onetime thing):

export FABRIC_CFG_PATH=$PWD/../config/

export CORE_PEER_TLS_ENABLED=true

export CORE_PEER_LOCALMSPID="Org1MSP"

export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp

export CORE_PEER_ADDRESS=localhost:7051

Installation of Ethereum platform cd ethereum
sudo ./ install.sh

Installation of Hyperledger Indy platform for starting the system pool for AnonCreds and run the AnonCreds funcionality under Hyperledger platform 
a.First clone Indy node repository for starting the repository using the commands-
git clone https://github.com/hyperledger/indy-node.git

Then go inside the directory using the command - cd indy-sdk

b. starting with a pre-configured docker image to build and run it for the pool:

docker build -f ci/indy-pool.dockerfile -t indy_pool.

docker run -itd -p 9701-9708: 9701-9708 ghoshbishakh/indy_pool

This creates an Indy container housing a pool of system-validated authenticators, each with a unique identity, pool number, and assigned port between 9701-9708.

c.Then run- docker ps , to get the container identity forexample in our case it is 351k39691g56. Then go inside the indy pool docker container using the command - docker exec -it 351k39691g56 bash

d. Now go inside the container 351k39691g56 and run the command - cat /var/lib/indy/sandbox/pool_transactions_genesis to get the details information of each validator nodes.

e. Now open a terminal and copy the information of all the information of validators nodes in to an text editotor that is opened using the command -'code.'. and past all the information into the text editor and save it as 'Validator_pool.txn' that is basically a type of JSON file for communication with the AnonCreds main code.

f. Import "DID_WalletAddress_Generator.py" into the folder containing 'Validator_pool.txn' to integrate Hyperledger AnonCreds for generating dynamic IDs, wallet addresses, and verifiable credentials for anonymous marketplace participant validation.

g. Now Import "Registrationcontract.sol" into the editor's folder containing 'Validator_pool.txn', then execute "DID_WalletAddress_Generator.py" to validate dynamic identities and wallet addresses before registering participants.

All such coding related the aforesaid procedure is already mentioned in DID_WalletAddress_Generator.py file

Now you have to consider for QSAM Model and for this model , we need to install pyqpanda==3.8.3.2 and pyvqnet==2.11.0. 

a. We have to use the following code to install 'pyvqnet' platform:

$ pip install pyvqnet --index-url https://pypi.originqc.com.cn

b. During installing pyvqnet, pyqpanda will be installed automatically as a dependency.

Installation of python 3.6, Jupiter Note book, Pandas, pyvqnet and pyvqnet for running the Quantum Machine Learning Algorithms using Qiskit Platforms, predict the valid quantum entangled states with Score Vectore, trace the Frauds using our proposed Tiny Quantum Machine Learning Optimization Model and successfully running the and working with the mechanism of Quantum Self Attantation Mechanism (QSAM), the combined and lightweight model - "tQMLTKSAM"  Model for the probability > 0.75 of state vectors for Frauds with high prediction capability. Such mechanisms are considered using the following steps.

a.The relevant information is uploaded in separate CSV files, linked to dynamically generated DID and Wallet Addresses, 'varifiable,' from 'Schema'information to facilitate efficient data exchange and validation (sample.csv).

b. Initially, install Python 3.6, Jupyter Notebook, and Pandas to efficiently generate valid 4-qubit entangled GHZ states using 'GHZ_Theta-Copy.ipynb' file and θ-protocol based on Generated DID and dynamic 160-bit addresses.

c. Afterwards, importing the file 'ghz_state_Fidelity_generation.py' Fidelity is generated for each 4-Qubit states from the .CSV files to check the Fidelity(F)>=0.50 or >= 50% and almost tends to 99% due to QSAM integration. Accondingly the valid 4-quibit engangled states are traced and fraud/invalid states are discarded.

d. Normalized and diagonally reduced the feature space with valid information using Quantum PCA (QPCA) as implementated in 'QPCA1.ipynb' file.

e. Run the basic QSAM model using 'QSAM.ipynb' file. (If you want to explore the functionality of QSAM model for our scheme)

f. Based on the valid 4-Qubit states for Fidelity(F)>=0.50, Feature related information is generated using ZZfeature Map for to generate Hilbert space for N-qubits of the valid GHZ states and encode the information with >= 50% accuracy.

g. Based on the Fidelity score, we can optimize the prediction mechanism for our proposed novel and lightweight model (tQMLTKSAM) using following two different models used for the tranning (tQMLSAM_train) and testing(tQMLSAM_test) purposes respectively. Additionally, we optimized using the tQMLSAM_test for the prediction mechanism from >=50% to >= 99%.


# Usage 
Take the QSAM/model/ directory for accessing the respective quantum models .

#Use the following code for model testing:

$ python tQMLSAM_test.py model-name

#Use the following code for model training:

$ python tQMLSAM_train.py model-name

The command-line argument model-name can be either tQMLSAM_test' or 'tQMLSAM_train'. The pre-trained models are located in the QSAM/model/ directory. 

# The lightweight tQML model (tQMLTKSAM) is running using the following details :
#Use the following code for model testing:
$ python tQMLSAM_test.py model-name
tQMLTKSAM Model
