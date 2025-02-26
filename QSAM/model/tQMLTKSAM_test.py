#!/usr/bin/env python
import argparse
import numpy as np
import jax
import jax.numpy as jnp
import equinox as eqx
import tensorcircuit as tc

# Custom imports (make sure these modules are available in your project)
from quantum.dataset import build_dataset
from quantum.embedding.onehot import embedding as onehot_embedding
from quantum.embedding.tf import embedding as tf_embedding
from quantum.tokenizer.ngram import get_tokenizer
from quantum.filter.stw import get_stw, stw_filter

# Import the tQMLTKSAMClassifier from your model definition file
# Adjust the import path as needed.
from tqmltksam_model import tQMLTKSAMClassifier


def test_model(args):
    # Build dataset using your tokenizer and filtering utilities
    tokenizer = get_tokenizer(ngram=[1], token_filter=stw_filter, la='en')
    vocab, idf, train_data, test_data = build_dataset(
        args.train_path,
        args.test_path,
        tokenizer,
        args.seq_len,
        need_pad=True
    )
    # Convert test data into numpy arrays
    x_test, y_test = np.asarray(test_data[0], dtype=int), np.asarray(test_data[1], dtype=int)
    
    # Generate embeddings for the test data (using one-hot & TF embeddings)
    x_test_we = onehot_embedding(x_test, len(vocab))
    x_test_se = tf_embedding(x_test, len(vocab), idf)
    
    # Initialize the model (the embed_dim is set to len(vocab))
    model = tQMLTKSAMClassifier(
        embed_dim=len(vocab),
        n_qubits=args.n_qubits,
        n_classes=args.n_classes,
        key=eqx.random.PRNGKey(args.seed)
    )
    
    # Load the trained model parameters
    model_params = eqx.tree_deserialise_leaves(args.model_path, model)
    model = model.replace(**model_params)
    
    # Run inference
    predictions = model(x_test_we)
    test_accuracy = jnp.mean(jnp.argmax(predictions, axis=-1) == y_test)
    
    print(f"Test Accuracy: {test_accuracy:.4f}")


def main():
    parser = argparse.ArgumentParser(description="Testing for tQMLTKSAM Model")
    parser.add_argument('--train_path', type=str, default='./data/train.npy',
                        help="Path to training data (used for building vocabulary)")
    parser.add_argument('--test_path', type=str, default='./data/test.npy',
                        help="Path to test data")
    parser.add_argument('--vocab_path', type=str, default='./data/vocab.quantum.pkl',
                        help="Path to vocabulary pickle file")
    parser.add_argument('--idf_path', type=str, default='./data/RP/idf.quantum.pkl',
                        help="Path to IDF pickle file")
    parser.add_argument('--model_path', type=str, default='./model/tQMLTKSAM.model',
                        help="Path to the saved tQMLTKSAM model")
    parser.add_argument('--n_qubits', type=int, default=6,
                        help="Number of qubits used in the quantum circuits")
    parser.add_argument('--n_classes', type=int, default=2,
                        help="Number of output classes")
    parser.add_argument('--seed', type=int, default=0,
                        help="Random seed for reproducibility")
    parser.add_argument('--seq_len', type=int, default=5,
                        help="Sequence length for tokenization (used in test mode)")
    
    args = parser.parse_args()
    test_model(args)


if __name__ == '__main__':
    main()
