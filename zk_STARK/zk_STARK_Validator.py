from mk_and_verify_proofs import verify_batch_proof, verify_trunk_proof, verify_single_inclusion_proof
from constants import MODULUS, COINS
import json
import sys
import re
import os

SUM_PROOF_PASS_INFO = "Total sum and non-negative constraint validation passed"
SUM_PROOF_FAIL_INFO = "Total sum and non-negative constraint validation failed"
INCLUSION_PROOF_PASS_INFO = "Inclusion constraint validation passed"
INCLUSION_PROOF_FAIL_INFO = "Inclusion constraint validation failed"



def por_user_verify_proofs():
    print("============ Validation started ==============")
    abs_dir = os.path.dirname(os.path.realpath(sys.argv[0]))
    for files in os.listdir(abs_dir):
        if (files == "sum_proof_data"):
            batches_proof_path = abs_dir + "/sum_proof_data/batches/"
            sum_values = [0] * (len(COINS) + 1)
            result = [0] * (len(COINS) + 1)
            success = True
            # verify batch proofs
            for root, dirs, subfiles in os.walk(batches_proof_path):
                for dir in dirs:
                    try:
                        result = verify_batch_proof(batches_proof_path + dir + "/")
                    except:
                        success = False
                        break
                    for i in range(len(sum_values)):
                        sum_values[i] = (sum_values[i] + result[i]) % MODULUS
            
            # verify trunk proofs
            if success:
                trunk_proof_path = abs_dir + "/sum_proof_data/trunk/"
                try:
                    result = verify_trunk_proof(trunk_proof_path)
                except:
                    success = False
            
            # check the consistence between batches and trunk
            if success:
                for i in range(len(sum_values)):
                    if sum_values[i] != (result[i] % MODULUS):
                        success = False
            # sum and non-negativity verification pass
            if success:
                print(SUM_PROOF_PASS_INFO)
            else:
                print(SUM_PROOF_FAIL_INFO)



        if re.search("inclusion_proof.json", files):
            try:
                verify_single_inclusion_proof(abs_dir + "/" + files)
                print(INCLUSION_PROOF_PASS_INFO)

            except:
                print(INCLUSION_PROOF_FAIL_INFO)
    
    print("============ Validation finished =============")
    input()



if __name__ == '__main__':
    por_user_verify_proofs()