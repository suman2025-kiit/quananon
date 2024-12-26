from qiskit import QuantumCircuit, execute, Aer
from qiskit.quantum_info import state_fidelity, Statevector

# Initialize a Quantum Circuit with 2 qubits and 2 classical bits
qc = QuantumCircuit(2, 2)

# Step 4: Apply a Hadamard gate on the first qubit
qc.h(0)

# Step 5: Apply a CNOT gate, with the first qubit as control and the second as target
qc.cx(0, 1)

# Simulate the state after the gates to check the GHZ state
backend = Aer.get_backend('statevector_simulator')
out_state = execute(qc, backend).result().get_statevector()
print("Statevector after applying Hadamard and CNOT:", out_state)

# Measure both qubits and store the results in the classical bits
qc.measure([0, 1], [0, 1])

# Execute the circuit on the qasm simulator
simulator = Aer.get_backend('qasm_simulator')
result = execute(qc, simulator, shots=1024).result()
counts = result.get_counts(qc)
print("Measurement results:", counts)

# Calculate and print the fidelity, assuming ideal |Phi+> state
ideal_state = [1 / (2**0.5), 0, 0, 1 / (2**0.5)]  # |Phi+> state
fidelity = state_fidelity(Statevector(ideal_state), out_state)
print("Fidelity with ideal Bell state |Phi+>:", fidelity)

# Determine if the state is valid
validity = "valid" if fidelity > 0.5 else "invalid"
print("The generated state is", validity)
