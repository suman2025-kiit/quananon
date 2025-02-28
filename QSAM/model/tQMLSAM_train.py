import numpy as np
import jax
import jax.numpy as jnp
import optax
import equinox as eqx
import pickle as pkl
import time
import tensorcircuit as tc
from tensorcircuit import shadows

from quantum.dataset import build_dataset
from quantum.embedding.onehot import embedding as onehot_embedding
from quantum.embedding.tf import embedding as tf_embedding
from quantum.model import tQTKSAMClassifier

# Training configurations
dataset_name = 'RP'
train_path = f'./data/train.csv'
test_path = f'./data/test.csv'
model_path = f'./model/tQMLSAM.train.model'

batch_size = 6
epochs = 100
learning_rate = 0.05

# Convert dataset into quantum-compatible format
with open(vocab_path, 'rb') as f:
    vocab = pkl.load(f)
with open(idf_path, 'rb') as f:
    idf = pkl.load(f)

x_train, y_train = np.load(train_path, allow_pickle=True)
x_test, y_test = np.load(test_path, allow_pickle=True)

# Apply embeddings
x_train_we, x_train_se = onehot_embedding(x_train, len(vocab)), tf_embedding(x_train, len(vocab), idf)
x_test_we, x_test_se = onehot_embedding(x_test, len(vocab)), tf_embedding(x_test, len(vocab), idf)

# Initialize model
key = jax.random.PRNGKey(0)
model = tQTKSAMClassifier(embed_dim=len(vocab), n_qubits=6, n_classes=2, key=key)

# Optimizer
optimizer = optax.adam(learning_rate)
state = optimizer.init(model)

# Training loop
best_accuracy = -1e5
for epoch in range(1, epochs + 1):
    start = time.time()
    
    loss, grads = eqx.filter_value_and_grad(lambda m: optax.softmax_cross_entropy_with_integer_labels(m(x_train_we), y_train).mean())(model)
    updates, state = optimizer.update(grads, state, model)
    model = eqx.apply_updates(model, updates)
    
    predictions = model(x_test_we)
    test_accuracy = (jnp.argmax(predictions, axis=-1) == y_test).mean()
    
    end = time.time()
    
    print(f'Epoch {epoch}: Loss = {loss:.4f}, Test Accuracy = {test_accuracy:.4f}, Time = {end-start:.2f}s')
    
    if test_accuracy > best_accuracy:
        best_accuracy = test_accuracy
        eqx.tree_serialise_leaves(model_path, model)

print(f'Best test accuracy achieved: {best_accuracy:.4f}')
