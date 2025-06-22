
import pennylane as qml
from pennylane import numpy as np
import matplotlib.pyplot as plt

# Define the quantum device
dev = qml.device("default.qubit", wires=4)

# Define quantum circuits for different models
def discocat_circuit(params):
    for i in range(4):
        qml.RY(params[i], wires=i)
    return qml.expval(qml.PauliZ(0))

def qsann_circuit(params):
    for i in range(4):
        qml.RX(params[i], wires=i)
        qml.CNOT(wires=[i, (i+1)%4])
    return qml.expval(qml.PauliZ(1))

def qsam_circuit(params):
    qml.templates.StronglyEntanglingLayers(params, wires=[0, 1, 2, 3])
    return qml.expval(qml.PauliZ(2))

def anon_tqmlsam_dew_circuit(params):
    qml.templates.AngleEmbedding(params, wires=range(4))
    qml.templates.BasicEntanglerLayers(params, wires=range(4))
    return qml.expval(qml.PauliZ(3))

# Define cost function
def cost_fn(model, params, target=1.0):
    return (model(params) - target)**2

# Optimizer
opt = qml.GradientDescentOptimizer(stepsize=0.2)

# Initialize
iterations = 140
losses_discocat = []
losses_qsann = []
losses_qsam = []
losses_anon = []

params_discocat = np.random.randn(4)
params_qsann = np.random.randn(4)
params_qsam = np.random.randn(3, 4, 3)
params_anon = np.random.randn(4)

# Training loop
for _ in range(iterations):
    params_discocat = opt.step(lambda p: cost_fn(discocat_circuit, p), params_discocat)
    params_qsann = opt.step(lambda p: cost_fn(qsann_circuit, p), params_qsann)
    params_qsam = opt.step(lambda p: cost_fn(qsam_circuit, p), params_qsam)
    params_anon = opt.step(lambda p: cost_fn(anon_tqmlsam_dew_circuit, p), params_anon)

    losses_discocat.append(cost_fn(discocat_circuit, params_discocat))
    losses_qsann.append(cost_fn(qsann_circuit, params_qsann))
    losses_qsam.append(cost_fn(qsam_circuit, params_qsam))
    losses_anon.append(cost_fn(anon_tqmlsam_dew_circuit, params_anon))

# Plotting
plt.figure(figsize=(12, 8))
plt.plot(losses_discocat, label='DisCoCat', color='green')
plt.plot(losses_qsann, label='QSANN', color='black')
plt.plot(losses_qsam, label='QSAM', color='red')
plt.plot(losses_anon, label='Anon-tQMLSAM-DeW', color='blue')
plt.xlabel("Iterations")
plt.ylabel("Training Loss")
plt.title("Training Loss Comparison on MC Dataset")
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.savefig("comparison_loss_mc.png")
plt.show()
