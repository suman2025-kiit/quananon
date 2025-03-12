# solver.py adjusted for the Anon-tQMLSAM-DeW Model

import tensorcircuit as tc
import jax
import jax.numpy as jnp
import equinox as eqx
import optax
from sklearn.metrics import balanced_accuracy_score

K = tc.set_backend("jax")
tc.set_dtype("complex128")

class Anon_tQMLSAM_DeW(eqx.Module):
    weights: jax.Array
    n_qubits: int
    n_layers: int
    encoding_type: int
    measurement_type: int

    def __init__(self, n_qubits, n_layers, encoding_type, measurement_type, key):
        self.weights = jax.random.normal(key, (n_qubits//2,))
        self.n_qubits = n_qubits
        self.encoding_type = encoding_type
        self.n_layers = n_layers
        self.measurement_type = measurement_type

    def zz_feature_map(self, c, features, wires):
        for i, wire in enumerate(wires):
            c.rx(wire, theta=features[i])
            if i < len(wires) - 1:
                c.rzz(wire, wires[i+1], theta=features[i] * features[i+1])
        return c

    @eqx.filter_jit
    def quantum_attention(self, x):
        c = tc.Circuit(self.n_qubits)
        for _ in range(self.n_layers):
            c = self.zz_feature_map(c, x, range(self.n_qubits))
        psi = c.state()
        psi_red = tc.quantum.reduced_density_matrix(psi, cut=list(range(self.n_qubits // 2)))

        if self.measurement_type == 0:
            return tc.quantum.renyi_entropy(psi_red, k=2)
        elif self.measurement_type == 1:
            return jnp.real(-jnp.trace(psi_red * jnp.log(psi_red)))
        else:
            return jnp.mean(jnp.abs(psi_red))

    def __call__(self, x):
        return K.vmap(self.quantum_attention, vectorized_argnums=(0))(x)

# Integrate Anon_tQMLSAM_DeW into existing training pipeline
class Solver:
    def __init__(self, args):
        self.args = args
        self.model = Anon_tQMLSAM_DeW(args.n_qubits, args.n_layers, args.encoding_type, args.measurement_type, args.key)
        self.optimizer = optax.adamw(args.lr)
        self.state = self.optimizer.init(self.model)

    @eqx.filter_value_and_grad
    def compute_loss(self, model, x, y):
        logits = model(x)
        loss = optax.softmax_cross_entropy_with_integer_labels(logits, y)
        return jnp.mean(loss)

    @eqx.filter_jit
    def update(self, x, y):
        loss, grads = self.compute_loss(self.model, x, y)
        updates, self.state = self.optimizer.update(grads, self.state, self.model)
        self.model = eqx.apply_updates(self.model, updates)
        return loss

    def train_epoch(self, train_loader):
        epoch_loss = 0
        for x, y in train_loader:
            epoch_loss += self.update(x, y)
        return epoch_loss / len(train_loader)

    def evaluate(self, loader):
        preds, labels = [], []
        for x, y in loader:
            logits = self.model(x)
            predictions = jnp.argmax(logits, axis=1)
            preds.extend(predictions)
            labels.extend(y)
        return balanced_accuracy_score(labels, preds)

    def train(self, train_loader, test_loader):
        for epoch in range(self.args.epochs):
            loss = self.train_epoch(train_loader)
            train_acc = self.evaluate(train_loader)
            test_acc = self.evaluate(test_loader)
            print(f"Epoch {epoch+1}: Train Loss={loss:.4f}, Train Acc={train_acc:.2f}, Test Acc={test_acc:.2f}")