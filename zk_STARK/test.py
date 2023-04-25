from mk_and_verify_proofs import *
from utils import * 
from constants import *
from por_verify_tool import verify_sum_proofs, verify_all_inclusion_proof
import json
import time

def test_full_process():
    init_user_data(15, 0)
    init_user_data(15, 1)
    init_user_data(15, 2)
    init_user_data(15, 3)
    
    mk_batch_proof(16, "./user_data/batch0.json", "./sum_proof_data/batches/a0/")
    mk_batch_proof(16, "./user_data/batch1.json", "./sum_proof_data/batches/a1/")
    mk_batch_proof(16, "./user_data/batch2.json", "./sum_proof_data/batches/a2/")
    mk_batch_proof(32, "./user_data/batch3.json", "./sum_proof_data/batches/b0/")
    
    mk_trunk_proof("./sum_proof_data/batches/", "./sum_proof_data/trunk/")
    
    verify_batch_proof("./sum_proof_data/batches/a0/")
    verify_batch_proof("./sum_proof_data/batches/a1/")
    verify_batch_proof("./sum_proof_data/batches/a2/")
    verify_batch_proof("./sum_proof_data/batches/b0/")
    
    verify_trunk_proof("./sum_proof_data/trunk/")

    mk_inclusion_proof(0, 16, "./sum_proof_data/batches/a0/", "./sum_proof_data/trunk/", "./inclusion_proof_data/a0/")
    mk_inclusion_proof(1, 16, "./sum_proof_data/batches/a1/", "./sum_proof_data/trunk/", "./inclusion_proof_data/a1/")
    mk_inclusion_proof(2, 16, "./sum_proof_data/batches/a2/", "./sum_proof_data/trunk/", "./inclusion_proof_data/a2/")
    mk_inclusion_proof(3, 32, "./sum_proof_data/batches/b0/", "./sum_proof_data/trunk/", "./inclusion_proof_data/b0/")

    verify_inclusion_proof("./inclusion_proof_data/a0/")
    verify_inclusion_proof("./inclusion_proof_data/a1/")
    verify_inclusion_proof("./inclusion_proof_data/a2/")
    verify_inclusion_proof("./inclusion_proof_data/b0/")
    
    print("Full Process Test Passed")

    return


def test_try_invalid_sum_value():
    with open("./sum_proof_data/batches/a0/sum_values.json", "r") as ff:
        sum_values_json = json.load(ff)
        sum_values_json["total_value"] += 1
        
    with open("./sum_proof_data/batches/a0/sum_values.json", "w") as ff:
        json.dump(sum_values_json, ff)
    try:
        verify_batch_proof("./sum_proof_data/batches/a0/")
    except:
        print("Invalid Sum Value")
    
    return

def test_negative_value_with_positive_net_value():
    with open("./user_data/batch0.json", "r") as ff:
        user_data_json = json.load(ff)
        user_data_json[0]["BTC"] = str(-1)
    
    with open("./user_data/batch0.json", "w") as ff:
        json.dump(user_data_json, ff)
    
    mk_batch_proof(16, "./user_data/batch0.json", "./sum_proof_data/batches/a0/")

    verify_batch_proof("./sum_proof_data/batches/a0/")
    
    print("Negative Value with Positive Net Value Is Allowed")

    init_user_data(15, 0)

    return

def test_negative_net_value():
    with open("./user_data/batch0.json", "r") as ff:
        user_data_json = json.load(ff)
        data_copy = user_data_json
        for coin in COINS:
            user_data_json[0][coin] = str(- int(user_data_json[0][coin]))
    
    with open("./user_data/batch0.json", "w") as ff:
        json.dump(user_data_json, ff)
    
    mk_batch_proof(16, "./user_data/batch0.json", "./sum_proof_data/batches/a0/")

    try:
        verify_batch_proof("./sum_proof_data/batches/a0/")
    except:
        print("Invalid Net Value")
    
    init_user_data(15, 0)
        
    return

def test_invalid_inclusion_proof():
    with open("./inclusion_proof_data/a0/user_0_inclusion_proof.json", "r") as ff:
        inclusion_proof_json = json.load(ff)  
        inclusion_proof_json["batch_inclusion_proof"]["total_value"] += 1
        
    with open("./inclusion_proof_data/a0/user_0_inclusion_proof.json", "w") as ff:
        json.dump(inclusion_proof_json, ff)
    
    try:
        verify_inclusion_proof("./inclusion_proof_data/a0/")
    except:
        print("Invalid Inclusion Proof")
    
    mk_inclusion_proof(0, 16, "./sum_proof_data/batches/a0/", "./sum_proof_data/trunk/", "./inclusion_proof_data/a0/")

    return

if __name__ == '__main__':
    time0 = time.time()
    test_full_process()
    test_try_invalid_sum_value()
    test_negative_value_with_positive_net_value()
    test_negative_net_value()   
    test_invalid_inclusion_proof()
    test_full_process()
    verify_sum_proofs()
    verify_all_inclusion_proof()
    print("all test finished in %d sec" %(time.time()-time0))