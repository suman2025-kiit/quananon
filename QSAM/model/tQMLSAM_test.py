import numpy as np
import jax.numpy as jnp
import optax
import equinox as eqx
from quantum.test import test_loop
from quantum.utils import stats
from quantum.dataset import build_dataset
from quantum.tokenizer.ngram import get_tokenizer
from quantum.filter.stw import get_stw, stw_filter
from quantum.embedding.onehot import embedding as onehot_embedding
from quantum.embedding.tf import embedding as tf_embedding
from quantum.model import tQTKSAMClassifier

# Load the dataset
train_path = './data/train.csv'
test_path = './data/test.csv'
vocab_path = './data/vocab.quantum.pkl'
idf_path = './data/RP/idf.quantum.pkl'

# Tokenizer
seq_len = 5
tokenizer = get_tokenizer(ngram=[1], token_filter=stw_filter, la='en')
vocab, idf, train_data, test_data = build_dataset(train_path, test_path, tokenizer, seq_len, need_pad=True)

# Convert data into quantum-compatible format
x_test, y_test = np.asarray(test_data[0], dtype=int), np.asarray(test_data[1], dtype=int)

# Embedding
x_test_we, x_test_se = onehot_embedding(x_test, len(vocab)), tf_embedding(x_test, len(vocab), idf)

# Load trained model
model_path = './model/tQMLSAM_test.model'
model = tQTKSAMClassifier(embed_dim=len(vocab), n_qubits=6, n_classes=2, key=eqx.random.PRNGKey(0))
model_para = eqx.tree_deserialise_leaves(model_path, model)
model = model.replace(**model_para)

# Define test function
def evaluate_model(model, x_test_we, x_test_se, y_test):
    test_loss, test_acc = test_loop(
        model,
        optax.softmax_cross_entropy_with_integer_labels,
        (x_test_we, x_test_se, y_test, 6),
        metric_func=stats.accuracy
    )
    print(f'Test Accuracy: {test_acc:.4f}')

# Run the evaluation
evaluate_model(model, x_test_we, x_test_se, y_test)
