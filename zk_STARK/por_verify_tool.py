from mk_and_verify_proofs import *
from utils import * 
from constants import *
import json
import sys
import time
import re

def verify_sum_proofs(config_path):
    with open(config_path, "r") as ff:
        config_json = json.load(ff)
        coins = config_json["coins"]
    abs_dir = os.path.dirname(os.path.realpath(sys.argv[0]))
    start_time = time.time()
    batches_proof_path = abs_dir + "/sum_proof_data/batches/"
    sum_values = [0] * (len(coins) + 1)
    for root, dirs, files in os.walk(batches_proof_path):
        for dir in dirs:
            result = verify_batch_proof(batches_proof_path + dir + "/", config_json)
            for i in range(len(sum_values)):
                sum_values[i] = (sum_values[i] + result[i]) % MODULUS
            print("Sum Proof of Batch %s Verified" %dir)

    trunk_proof_path = abs_dir + "/sum_proof_data/trunk/"
    result = verify_trunk_proof(trunk_proof_path, config_json)
    for i in range(len(sum_values)):
        assert sum_values[i] == (result[i] % MODULUS)
    print("Sum Proof of Trunk Verified")

    print("All Proofs Verified in %.4f secs!" %(time.time() - start_time))

def verify_all_inclusion_proof():
    start_time = time.time()
    abs_dir = os.path.dirname(os.path.realpath(sys.argv[0]))
    for root, dirs, files in os.walk(abs_dir+ "/inclusion_proof_data"):
        for dir in dirs:
            verify_inclusion_proof(abs_dir + "/" + dir + "/")
            print("Inclusion Proof %s Verified" %dir)

    print("All Proofs Verified in %.4f secs!" %(time.time() - start_time))

if __name__ == '__main__':
    verify_sum_proofs(CONFIG_PATH)
    verify_all_inclusion_proof()