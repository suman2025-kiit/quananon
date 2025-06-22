
import pennylane as qml
from pennylane import numpy as np
from sklearn.datasets import make_classification
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler
from sklearn.metrics import accuracy_score
import matplotlib.pyplot as plt

# Generate synthetic dataset
X, y = make_classification(n_samples=100, n_features=4, n_informative=3, n_redundant=0, random_state=42)
X = StandardScaler().fit_transform(X)
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.3, random_state=42)

n_qubits = 4
dev = qml.device("default.qubit", wires=n_qubits)

# ZZ-Feature Map (quantum feature embedding)
def zz_feature_map(x):
    for i in range(n_qubits):
        qml.Hadamard(wires=i)
    for i in range(n_qubits):
        qml.RZ(x[i], wires=i)
    for i in range(n_qubits - 1):
        qml.CNOT(wires=[i, i + 1])
        qml.RZ((np.pi - x[i] * x[i+1]), wires=i + 1)
        qml.CNOT(wires=[i, i + 1])

# Ansatz (trainable)
@qml.qnode(dev, interface="autograd")
def circuit(weights, x=None):
    zz_feature_map(x)
    qml.BasicEntanglerLayers(weights, wires=range(n_qubits))
    return qml.expval(qml.PauliZ(0))

# Variational classifier
def variational_classifier(weights, bias, x):
    return circuit(weights, x) + bias

def cost(weights, bias, X, Y):
    predictions = [variational_classifier(weights, bias, x) for x in X]
    return np.mean((np.sign(predictions) - Y) ** 2)

# Label conversion
Y_train = 2 * y_train - 1
Y_test = 2 * y_test - 1

np.random.seed(0)
weights = 0.01 * np.random.randn(3, n_qubits)
bias = 0.0
opt = qml.GradientDescentOptimizer(stepsize=0.4)

# Training loop
losses = []
for it in range(50):
    weights, bias, loss = opt.step_and_cost(lambda w, b: cost(w, b, X_train, Y_train), weights, bias)
    losses.append(loss)

# Evaluation
predictions = [np.sign(variational_classifier(weights, bias, x)) for x in X_test]
accuracy = accuracy_score(Y_test, predictions)
print(f"Test Accuracy: {accuracy * 100:.2f}%")

# Plotting training loss
plt.plot(losses)
plt.xlabel("Iterations")
plt.ylabel("Loss")
plt.title("Training Loss Curve for Anon-tQMLSAM-DeW")
plt.grid(True)
plt.tight_layout()
plt.savefig("training_loss_curve.png")
plt.show()
