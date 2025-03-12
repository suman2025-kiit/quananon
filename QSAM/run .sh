# Customized run.sh for Anon-tQMLSAM-DeW Model

attn_layers_list=(2 3)
n_layers_list=(4 5)
use_bias_list=(1 0)
dataset_list=(RP MC)
measurement_type_list=(0 1 2)

for dataset in "${dataset_list[@]}"; do
    for attn_layers in "${attn_layers_list[@]}"; do
        for n_layers in "${n_layers_list[@]}"; do
            for use_bias in "${use_bias_list[@]}"; do
                for measurement_type in "${measurement_type_list[@]}"; do
                    echo "Running dataset=$dataset with attn_layers=$attn_layers, n_layers=$n_layers, use_bias=$use_bias, measurement_type=$measurement_type"
                    time python main.py \
                        --seed 42 \
                        --dataset $dataset \
                        --n_features 12 \
                        --batch_size 32 \
                        --q_lr 5e-3 \
                        --c_lr 5e-3 \
                        --epochs 300 \
                        --runs 5 \
                        --num_samples_infer 300 \
                        --n_attention_layers $attn_layers \
                        --n_qubits 8 \
                        --n_layers $n_layers \
                        --encoding_type 2 \
                        --measurement_type $measurement_type \
                        --use_bias $use_bias \
                        --n_classes 2 \
                        --noise 1
                done
            done
    done
done
